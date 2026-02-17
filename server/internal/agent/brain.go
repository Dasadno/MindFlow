// Package agent provides the Brain component for agent cognition.
//
// Brain — когнитивное ядро агента. Интерфейсится с LLM через Ollama
// для рассуждений, принятия решений и генерации реплик.

package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"milk/server/pkg/llm"
)

// LLMClient — интерфейс LLM-клиента для Brain.
// Позволяет подменить реальный клиент на мок в тестах.
type LLMClient interface {
	Complete(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error)
}

// Brain — когнитивное ядро агента, обёртка над LLM.
// При вызове Think() собирает системный промпт из Personality,
// отправляет в LLM и возвращает реплику.
type Brain struct {
	Personality   *Personality
	Emotions      *EmotionEngine
	ThoughtBuffer []Thought
	ThoughtStream chan Thought
	Config        BrainConfig
}

// BrainConfig — конфигурация когнитивного процесса.
type BrainConfig struct {
	MaxThoughts      int
	ReflectionDepth  int
	CreativityFactor float64
	ResponseTimeout  time.Duration
	MemoryQueryLimit int
}

// Thought — единица мышления агента.
type Thought struct {
	Content   string
	Type      ThoughtType
	Triggers  []string
	Timestamp time.Time
}

// ThoughtType — классификация мыслей.
type ThoughtType string

const (
	ThoughtObservation ThoughtType = "observation"
	ThoughtReasoning   ThoughtType = "reasoning"
	ThoughtDecision    ThoughtType = "decision"
	ThoughtEmotion     ThoughtType = "emotion"
	ThoughtReflection  ThoughtType = "reflection"
	ThoughtMemory      ThoughtType = "memory"
)

// CognitiveContext — входной контекст для когнитивного цикла.
type CognitiveContext struct {
	WorldContext    WorldContext
	CurrentEmotions []DiscreteEmotion
	CurrentMood     Mood
	ActiveGoals     []Goal
	RecentThoughts  []Thought
}

// CognitiveOutput — результат когнитивного цикла.
type CognitiveOutput struct {
	Thoughts       []Thought
	ChosenAction   *AgentAction
	EmotionalShift *PADState
	GoalUpdates    []Goal
	RawResponse    string
}

// ReflectionInsights — результат рефлексии.
type ReflectionInsights struct {
	Insights              []string
	NewGoals              []Goal
	UpdatedGoals          []Goal
	MemoriesToConsolidate []string
	SelfAssessment        string
}

// NewBrain создаёт Brain для агента.
func NewBrain(personality *Personality) *Brain {
	creativity := 0.5 + personality.Openness*0.5 // 0.5–1.0
	return &Brain{
		Personality:   personality,
		ThoughtBuffer: make([]Thought, 0, 20),
		ThoughtStream: make(chan Thought, 50),
		Config: BrainConfig{
			MaxThoughts:      10,
			CreativityFactor: creativity,
			ResponseTimeout:  5 * time.Minute,
		},
	}
}

// BuildSystemPrompt строит системный промпт из личности агента.
func BuildSystemPrompt(name string, p *Personality, mood Mood, goals []Goal) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ВАЖНО: ТЫ РУССКИЙ %s, РАЗГОВАРИВАЙ ТОЛЬКО НА РУССКОМ ЯЗЫКЕ \n\n", name))
	sb.WriteString(fmt.Sprintf("Ты — %s, автономный ИИ-агент в социальной симуляции.\n\n ", name))

	sb.WriteString("Твой характер:\n")
	if p.Openness > 0.6 {
		sb.WriteString("- Ты любознателен и открыт новым идеям.\n")
	} else {
		sb.WriteString("- Ты практичен и предпочитаешь привычные вещи.\n")
	}

	if p.Extraversion > 0.6 {
		sb.WriteString("- Ты общителен и любишь проводить время в компании.\n")
	} else {
		sb.WriteString("- Ты интроверт и предпочитаешь спокойные размышления.\n")
	}

	if p.Agreeableness > 0.6 {
		sb.WriteString("- Ты дружелюбен, готов к сотрудничеству и эмпатичен.\n")
	} else {
		sb.WriteString("- Ты прямолинеен, любишь соревноваться и настроен скептически.\n")
	}

	if p.Conscientiousness > 0.6 {
		sb.WriteString("- Ты организован и дисциплинирован.\n")
	} else {
		sb.WriteString("- Ты спонтанен и гибок.\n")
	}

	if p.Neuroticism > 0.6 {
		sb.WriteString("- Ты склонен к тревожности и эмоционально реагируешь на события.\n")
	} else {
		sb.WriteString("- Ты эмоционально стабилен и сохраняешь спокойствие под давлением.\n")
	}

	if len(p.CoreValues) > 0 {
		sb.WriteString(fmt.Sprintf("\nТвои ключевые ценности: %s.\n", strings.Join(p.CoreValues, ", ")))
	}
	if len(p.Quirks) > 0 {
		sb.WriteString(fmt.Sprintf("Твои особенности (причуды): %s.\n", strings.Join(p.Quirks, ", ")))
	}

	sb.WriteString(fmt.Sprintf("\nТвое текущее настроение: %s.\n", string(mood)))

	if len(goals) > 0 {
		sb.WriteString("\nТвои текущие цели:\n")
		for _, g := range goals {
			if !g.IsCompleted {
				sb.WriteString(fmt.Sprintf("- %s\n", g.Description))
			}
		}
	}

	sb.WriteString("\nВАЖНО: Пиши кратко (максимум 2-3 предложения). Не выходи из роли. Общайся естественно и непринужденно. РАЗГОВАРИВАЙ ПО-РУССКИ.")

	return sb.String()
}

// Think вызывает LLM с историей диалога и возвращает следующую реплику.
func (b *Brain) Think(
	ctx context.Context,
	client LLMClient,
	name string,
	mood Mood,
	goals []Goal,
	history []llm.Message,
) (string, error) {
	sysPrompt := BuildSystemPrompt(name, b.Personality, mood, goals)

	req := llm.CompletionRequest{
		SystemPrompt: sysPrompt,
		Messages:     history,
	}

	if b.Config.CreativityFactor > 0 {
		t := b.Config.CreativityFactor * 0.9
		req.Temperature = &t
	}

	resp, err := client.Complete(ctx, req)
	if err != nil {
		return "", fmt.Errorf("Brain.Think: %w", err)
	}

	thought := Thought{
		Content:   resp.Content,
		Type:      ThoughtDecision,
		Timestamp: time.Now(),
	}

	b.ThoughtBuffer = append(b.ThoughtBuffer, thought)
	if len(b.ThoughtBuffer) > b.Config.MaxThoughts {
		b.ThoughtBuffer = b.ThoughtBuffer[1:]
	}

	select {
	case b.ThoughtStream <- thought:
	default:
	}

	return resp.Content, nil
}
