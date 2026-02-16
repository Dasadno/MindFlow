// Package agent provides the Brain component for agent cognition.
//
// The Brain is the cognitive core of an agent. It interfaces with the LLM
// (GigaChat) to perform reasoning, decision-making, and natural language
// generation. It orchestrates the thinking process by combining personality,
// memory, and emotional context.

package agent

import "time"

// -----------------------------------------------------------------------------
// Brain — когнитивное ядро агента, обёртка над LLM
// -----------------------------------------------------------------------------
// Brain не хранит данные напрямую — он ссылается на Personality, Memory и Emotions
// агента-владельца. При каждом вызове Think() он:
//   1. Собирает системный промпт из Personality (buildSystemPrompt)
//   2. Извлекает релевантные воспоминания из Memory (vector search)
//   3. Добавляет текущее эмоциональное состояние из Emotions
//   4. Отправляет промпт в GigaChat через llmClient
//   5. Парсит структурированный ответ в CognitiveOutput
//
// CreativityFactor в BrainConfig напрямую влияет на temperature LLM —
// агент с высоким Openness будет генерировать более неожиданные мысли.

type Brain struct {
	// personality — ссылка на личность агента-владельца.
	// Используется для построения системного промпта (buildSystemPrompt).
	Personality *Personality

	// memory — ссылка на систему памяти для извлечения контекста.
	// Think() делает vector search по релевантным воспоминаниям.
	//	Memory *MemorySystem

	// emotions — ссылка на движок эмоций.
	// Текущее эмоциональное состояние включается в контекстный промпт.
	Emotions *EmotionEngine

	// thoughtBuffer — буфер недавних мыслей (рабочая память мышления).
	// Ограничен Config.MaxThoughts. Новые мысли вытесняют старые.
	ThoughtBuffer []Thought

	// thoughtStream — канал для стриминга мыслей на дашборд в реальном времени.
	// Подключается к SSE-эндпоинту GET /api/v1/agents/:id/thoughts.
	ThoughtStream chan Thought

	// config — параметры когнитивного процесса.
	Config BrainConfig
}

// BrainConfig — конфигурация когнитивного процесса.
type BrainConfig struct {
	// MaxThoughts — ёмкость буфера рабочей памяти мышления.
	// Определяет, сколько предыдущих мыслей учитывается при следующем Think().
	MaxThoughts int

	// ReflectionDepth — глубина мета-когниции при рефлексии.
	// 1 = простая рефлексия, 2+ = рефлексия о рефлексии (мета-мета-...).
	ReflectionDepth int

	// CreativityFactor — модификатор температуры LLM (0.0–2.0).
	// Вычисляется из Personality.Openness при создании.
	// 0.3 = консервативные, предсказуемые мысли.
	// 1.5 = креативные, неожиданные связи.
	CreativityFactor float64

	// ResponseTimeout — максимальное время ожидания ответа от LLM.
	ResponseTimeout time.Duration

	// MemoryQueryLimit — сколько воспоминаний извлекать из vector store для контекста.
	MemoryQueryLimit int
}

// -----------------------------------------------------------------------------
// Thought — единица мышления агента
// -----------------------------------------------------------------------------
// Генерируется Brain.Think() и Brain.Reflect(). Каждая мысль попадает:
//   1. В thoughtBuffer — для контекста следующих мыслей
//   2. В thoughtStream — для real-time отображения на дашборде (SSE)
//   3. Опционально в Memory — если достаточно значима

type Thought struct {
	// Content — текстовое содержимое мысли ("Кажется, Алиса сегодня грустит...").
	Content string

	// Type — классификация мысли.
	Type ThoughtType

	// Triggers — что вызвало эту мысль (IDs стимулов, воспоминаний или предыдущих мыслей).
	Triggers []string

	// Timestamp — когда мысль была сгенерирована.
	Timestamp time.Time
}

// ThoughtType — классификация мыслей для фильтрации и отображения.
type ThoughtType string

const (
	ThoughtObservation ThoughtType = "observation" // Наблюдение за окружением: "I see Agent X nearby"
	ThoughtReasoning   ThoughtType = "reasoning"   // Цепочка рассуждений: "If X then Y, so I should..."
	ThoughtDecision    ThoughtType = "decision"    // Принятое решение: "I will talk to Agent X"
	ThoughtEmotion     ThoughtType = "emotion"     // Эмоциональная реакция: "This makes me feel..."
	ThoughtReflection  ThoughtType = "reflection"  // Мета-когниция: "I notice I've been anxious lately"
	ThoughtMemory      ThoughtType = "memory"      // Вспоминание: "This reminds me of the time when..."
)

// -----------------------------------------------------------------------------
// CognitiveContext — входной контекст для когнитивного цикла
// -----------------------------------------------------------------------------
// Формируется перед вызовом Brain.Think(). Объединяет всё, что агент
// «видит, чувствует и помнит» в данный момент.

type CognitiveContext struct {
	// WorldContext — контекст мира (другие агенты, события, тик).
	WorldContext WorldContext

	// RelevantMemories — воспоминания, извлечённые vector search по текущей ситуации.
	//TODO	RelevantMemories []MemoryEntry

	// CurrentEmotions — активные дискретные эмоции в данный момент.
	CurrentEmotions []DiscreteEmotion

	// CurrentMood — текущее настроение (дискретная метка).
	CurrentMood Mood

	// ActiveGoals — текущие цели агента, отсортированные по приоритету.
	ActiveGoals []Goal

	// RecentThoughts — последние N мыслей из thoughtBuffer.
	RecentThoughts []Thought
}

// -----------------------------------------------------------------------------
// CognitiveOutput — структурированный результат когнитивного цикла
// -----------------------------------------------------------------------------
// Brain парсит ответ LLM в эту структуру. Она содержит:
//   - новые мысли для thoughtBuffer
//   - решение о действии
//   - эмоциональную реакцию (если есть)

type CognitiveOutput struct {
	// Thoughts — сгенерированные мысли (может быть несколько за один цикл).
	Thoughts []Thought

	// ChosenAction — действие, которое агент решил предпринять.
	ChosenAction *AgentAction

	// EmotionalShift — изменение эмоционального состояния в результате мышления.
	// nil, если мышление не вызвало эмоциональной реакции.
	EmotionalShift *PADState

	// GoalUpdates — изменения в целях (новые цели, обновление прогресса).
	GoalUpdates []Goal

	// RawResponse — сырой ответ LLM для отладки.
	RawResponse string
}

// -----------------------------------------------------------------------------
// ReflectionInsights — результат рефлексии (мета-когниции)
// -----------------------------------------------------------------------------
// Генерируется Brain.Reflect(). Рефлексия — периодический процесс (не каждый тик),
// во время которого агент анализирует накопленный опыт и обновляет самомодель.

type ReflectionInsights struct {
	// Insights — выявленные паттерны и наблюдения о себе и мире.
	// Пример: "I seem to get anxious when Agent X is around"
	Insights []string

	// NewGoals — новые цели, сформированные в результате рефлексии.
	NewGoals []Goal

	// UpdatedGoals — существующие цели с обновлённым прогрессом или приоритетом.
	UpdatedGoals []Goal

	// MemoriesToConsolidate — ID эпизодических воспоминаний, которые стоит
	// объединить в семантическую память.
	MemoriesToConsolidate []string

	// SelfAssessment — свободный текст самооценки ("I've been too reclusive lately").
	SelfAssessment string
}
