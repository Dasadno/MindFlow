// Package agent provides the Memory System for agent cognition.
//
// The Memory System manages three tiers of memory:
//   Working Memory  → short-term circular buffer (last few events)
//   Episodic Memory → vector database for similarity-based retrieval (experiences)
//   Semantic Memory → SQLite for persistent generalized knowledge (facts, beliefs)
//
// Memories flow: Working → Episodic → Semantic (via consolidation).

package agent

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// -----------------------------------------------------------------------------
// MemorySystem — трёхуровневая система памяти агента
// -----------------------------------------------------------------------------
// Архитектура повторяет когнитивную модель человеческой памяти:
//   Working Memory  — «что я только что видел» (буфер последних N событий)
//   Episodic Memory — «что со мной было» (конкретные эпизоды, vector search)
//   Semantic Memory — «что я знаю» (обобщённые знания, факты)
//
// MemoryConsolidator периодически (при рефлексии / «сне») объединяет
// повторяющиеся эпизодические воспоминания в семантические знания.

type MemorySystem struct {
	// WorkingMem — кратковременная память, кольцевой буфер фиксированного размера.
	// Хранит последние N событий/стимулов. Переполнение → запись в Episodic.
	WorkingMem *WorkingMemory

	// Config — настройки системы памяти.
	Config MemoryConfig
}

// MemoryConfig — конфигурация системы памяти.
type MemoryConfig struct {
	// WorkingMemoryCapacity — размер буфера кратковременной памяти.
	// Типичное значение: 7±2 (как у человека, закон Миллера).
	WorkingMemoryCapacity int

	// EpisodicSearchLimit — сколько результатов возвращать из vector search.
	EpisodicSearchLimit int

	// ConsolidationThreshold — минимальное количество связанных эпизодических
	// воспоминаний для запуска консолидации в семантическую память.
	ConsolidationThreshold int

	// ForgetThreshold — порог важности, ниже которого воспоминания «забываются».
	// Применяется вместе с access_count и recency.
	ForgetThreshold float64

	// ImportanceDecayRate — скорость снижения importance для редко вспоминаемых записей.
	ImportanceDecayRate float64
}

// -----------------------------------------------------------------------------
// MemoryEntry — единица памяти (хранится в SQLite и/или Vector DB)
// -----------------------------------------------------------------------------
// Каждое воспоминание проходит путь:
//   1. Encode() — стимул → MemoryEntry (вычисляются Embedding, Importance)
//   2. Хранение в Episodic (vector DB) и/или Semantic (SQLite)
//   3. Recall() — поиск по similarity + recency + importance + access frequency
//   4. Forget() — удаление, если importance * recency * access < threshold

type MemoryEntry struct {
	// ID — уникальный идентификатор (UUID).
	ID string

	// Type — тип памяти: эпизодическая, семантическая или процедурная.
	Type MemoryType

	// Content — текстовое содержимое воспоминания на естественном языке.
	// Пример: "Had a deep conversation with Bob about the meaning of trust"
	Content string

	// Embedding — векторное представление Content для similarity search.
	// Генерируется через GigaChat.Embed() или LocalEmbedder.
	Embedding []float32

	// EmotionalTag — эмоция в момент формирования воспоминания.
	// Эмоционально окрашенные воспоминания имеют повышенный Importance
	// и медленнее забываются (как в реальной психологии).
	EmotionalTag EmotionType

	// Importance — вычисленная значимость от 0.0 до 1.0.
	// Формула: f(emotional_intensity, novelty, goal_relevance).
	// Высокая важность → долгое хранение, приоритет при Recall.
	Importance float64

	// Timestamp — когда воспоминание было сформировано.
	Timestamp time.Time

	// AccessCount — сколько раз воспоминание было извлечено через Recall.
	// Часто вспоминаемые воспоминания укрепляются (importance не падает).
	AccessCount int

	// LastAccessed — когда воспоминание было последний раз извлечено.
	// Вместе с AccessCount определяет «свежесть» при Recall.
	LastAccessed time.Time

	// RelatedAgents — ID агентов, участвовавших в событии.
	// Используется для RecallRelated() — построения контекста отношений.
	RelatedAgents []string

	// Metadata — произвольные дополнительные данные (место, тип взаимодействия и т.д.)
	Metadata map[string]any
}

// MemoryType — классификация типов памяти.
type MemoryType string

const (
	// MemoryEpisodic — конкретные события и опыт.
	// Пример: "Met Alice at the park on tick 150, she was happy"
	MemoryEpisodic MemoryType = "episodic"

	// MemorySemantic — обобщённые знания и убеждения.
	// Пример: "Alice is generally a cheerful person who likes parks"
	// Формируется из нескольких эпизодических воспоминаний через консолидацию.
	MemorySemantic MemoryType = "semantic"

	// MemoryProcedural — знания о том, «как делать» (паттерны поведения).
	// Пример: "When someone is sad, offering help improves the relationship"
	MemoryProcedural MemoryType = "procedural"
)

// -----------------------------------------------------------------------------
// WorkingMemory — кратковременная память (кольцевой буфер)
// -----------------------------------------------------------------------------
// Хранит последние N событий. При переполнении старые записи вытесняются.
// Используется Brain.Think() для немедленного контекста —
// «что произошло только что» (последние 1-3 тика).

type WorkingMemory struct {
	// Buffer — массив последних записей, кольцевой буфер.
	Buffer []MemoryEntry // Нужно подумать насчет этого

	// Capacity — максимальный размер буфера (обычно 5–10).
	Capacity int

	// Current — индекс следующей позиции для записи (циклически).
	Current int

	// Count — текущее количество записей (≤ Capacity).
	Count int
}

// -----------------------------------------------------------------------------
// Experience — сырой опыт до преобразования в MemoryEntry
// -----------------------------------------------------------------------------
// Передаётся в MemorySystem.Encode(). Содержит исходные данные о событии
// до расчёта importance и генерации embedding.

type Experience struct {
	// Content — описание опыта на естественном языке.
	Content string

	// Source — откуда пришёл опыт (agent ID, "world", "self").
	Source string

	// EmotionalContext — эмоциональное состояние агента в момент опыта.
	// Влияет на importance и emotional_tag сформированного воспоминания.
	EmotionalContext PADState

	// RelatedAgents — ID агентов, вовлечённых в событие.
	RelatedAgents []string

	// Timestamp — когда произошло событие.
	Timestamp time.Time
}

func (b *Brain) Memoring(ctx context.Context, memory string) error { // TODO
	b.Memory.WorkingMem.Buffer = append(b.Memory.WorkingMem.Buffer, MemoryEntry{
		ID:           uuid.New().String(),
		Type:         MemoryEpisodic,
		Content:      memory,
		Importance:   rand.Float64(),
		Timestamp:    time.Now(),
		AccessCount:  0,
		LastAccessed: time.Now(),
	})
	return nil
}
