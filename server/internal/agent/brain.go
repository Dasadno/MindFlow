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
	Memory        *MemorySystem
	Personality   *Personality
	Emotions      *EmotionEngine
	ThoughtBuffer []Thought
	ThoughtStream chan Thought
	Config        BrainConfig
}

// BrainConfig — конфигурация когнитивного процесса.
type BrainConfig struct {
	Memories         []string
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
		ThoughtBuffer: make([]Thought, 0, 5),
		ThoughtStream: make(chan Thought, 16),
		Config: BrainConfig{
			Memories:         make([]string, 0, 10),
			MaxThoughts:      5,
			CreativityFactor: creativity,
			ResponseTimeout:  5 * time.Minute,
		},
	}
}

// BuildSystemPrompt строит системный промпт из личности агента.
func (b *Brain) BuildSystemPrompt(name string, p *Personality, mood Mood, goals []Goal) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Ты русский %s, разговаривай только на русском\n\n", name))
	sb.WriteString(fmt.Sprintf("You are %s, an autonomous AI agent in a social simulation.\n\n ", name))

	sb.WriteString("Your personality:\n")
	if p.Openness > 0.6 {
		sb.WriteString("- You are curious and open to new ideas.\n")
	} else {
		sb.WriteString("- You are practical and prefer familiar things.\n")
	}
	if p.Extraversion > 0.6 {
		sb.WriteString("- You are outgoing and enjoy socializing.\n")
	} else {
		sb.WriteString("- You are introverted and prefer quiet reflection.\n")
	}
	if p.Agreeableness > 0.6 {
		sb.WriteString("- You are warm, cooperative, and empathetic.\n")
	} else {
		sb.WriteString("- You are direct, competitive, and skeptical.\n")
	}
	if p.Conscientiousness > 0.6 {
		sb.WriteString("- You are organized and disciplined.\n")
	} else {
		sb.WriteString("- You are spontaneous and flexible.\n")
	}
	if p.Neuroticism > 0.6 {
		sb.WriteString("- You tend to be anxious and emotionally reactive.\n")
	} else {
		sb.WriteString("- You are emotionally stable and calm under pressure.\n")
	}

	if len(p.CoreValues) > 0 {
		sb.WriteString(fmt.Sprintf("\nYour core values: %s.\n", strings.Join(p.CoreValues, ", ")))
	}
	if len(p.Quirks) > 0 {
		sb.WriteString(fmt.Sprintf("Your quirks: %s.\n", strings.Join(p.Quirks, ", ")))
	}

	sb.WriteString(fmt.Sprintf("\nYour current mood: %s.\n", string(mood)))

	if len(goals) > 0 {
		sb.WriteString("\nYour current goals:\n")
		for _, g := range goals {
			if !g.IsCompleted {
				sb.WriteString(fmt.Sprintf("- %s\n", g.Description))
			}
		}
	}

	sb.WriteString("\nIMPORTANT: Keep responses concise (2-3 sentences max). Stay in character. Be natural and conversational.")
	b.Config.Memories = append(b.Config.Memories, sb.String())
	return sb.String()
}

// Think вызывает LLM с историей диалога и возвращает следующую реплику.
func (a *Agent) Think(
	ctx context.Context,
	client LLMClient,
	name string,
	mood Mood,
	goals []Goal,
	history []llm.Message,
) (string, error) {
	sysPrompt := a.Brain.BuildSystemPrompt(name, a.Brain.Personality, mood, goals)

	req := llm.CompletionRequest{
		SystemPrompt: sysPrompt,
		Messages:     history,
	}

	if a.Brain.Config.CreativityFactor > 0 {
		t := a.Brain.Config.CreativityFactor * 0.9
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

	a.Brain.Config.Memories = append(a.Brain.Config.Memories, thought.Content)

	a.Brain.ThoughtBuffer = append(a.Brain.ThoughtBuffer, thought)
	if len(a.Brain.ThoughtBuffer) > a.Brain.Config.MaxThoughts {
		a.Brain.ThoughtBuffer = a.Brain.ThoughtBuffer[1:]
	}

	select {
	case a.Brain.ThoughtStream <- thought:
	default:
	}

	return resp.Content, nil
}
