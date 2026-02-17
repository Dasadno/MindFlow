// Package api provides Data Transfer Objects for API request/response handling.
//
// DTOs деcouple API-контракт от внутренних доменных моделей.
// Фронтенд-команда работает ТОЛЬКО с этими структурами — изменения
// в internal/agent не ломают API, пока маппинг обновляется.

package api

import "time"

// =============================================================================
// PAGINATION
// =============================================================================

// PaginationMeta — метаданные пагинации, включаемые в список-ответы.
type PaginationMeta struct {
	// Page — текущая страница (начинается с 1).
	Page int `json:"page"`

	// Limit — количество элементов на странице.
	Limit int `json:"limit"`

	// Total — общее количество элементов (для расчёта страниц на фронте).
	Total int `json:"total"`
}

// =============================================================================
// AGENT DTOs
// =============================================================================

// AgentListResponse — ответ на GET /api/v1/agents.
type AgentListResponse struct {
	Agents     []AgentSummary `json:"agents"`
	Pagination PaginationMeta `json:"pagination"`
}

// AgentSummary — краткая карточка агента для списков и графов.
// Содержит только публичную информацию, без деталей внутреннего состояния.
type AgentSummary struct {
	// ID — UUID агента.
	ID string `json:"id"`

	// Name — отображаемое имя.
	Name string `json:"name"`

	// PersonalityType — обобщённый тип личности ("explorer", "guardian", "introvert").
	// Вычисляется из доминантных черт Big Five для UI-категоризации.
	PersonalityType string `json:"personalityType"`

	// CurrentMood — текущая дискретная метка настроения ("happy", "anxious").
	CurrentMood string `json:"currentMood"`

	// MoodIntensity — интенсивность настроения от 0.0 до 1.0.
	// Используется для размера/яркости иконки на дашборде.
	MoodIntensity float64 `json:"moodIntensity"`

	// IsActive — активен ли агент в симуляции (false = soft-deleted).
	IsActive bool `json:"isActive"`
}

// AgentDetailResponse — полный профиль агента, ответ на GET /api/v1/agents/:id.
type AgentDetailResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Personality PersonalityDTO `json:"personality"`
	CurrentMood MoodDTO        `json:"currentMood"`
	Goals       []GoalDTO      `json:"goals"`
	Stats       AgentStatsDTO  `json:"stats"`
	CreatedAt   time.Time      `json:"createdAt"`
}

// PersonalityDTO — черты личности Big Five для отображения на дашборде.
// Каждое значение 0.0–1.0 может отображаться как radar chart.
type PersonalityDTO struct {
	Openness          float64  `json:"openness"`
	Conscientiousness float64  `json:"conscientiousness"`
	Extraversion      float64  `json:"extraversion"`
	Agreeableness     float64  `json:"agreeableness"`
	Neuroticism       float64  `json:"neuroticism"`
	CoreValues        []string `json:"coreValues"`
	Quirks            []string `json:"quirks"`
}

// MoodDTO — текущее эмоциональное состояние агента.
type MoodDTO struct {
	// Label — дискретная метка настроения ("happy", "anxious").
	Label string `json:"label"`

	// PAD — непрерывное состояние Pleasure-Arousal-Dominance.
	PAD PADDTO `json:"pad"`

	// ActiveEmotions — список активных дискретных эмоций с причинами.
	ActiveEmotions []EmotionDTO `json:"activeEmotions"`
}

// PADDTO — JSON-представление PAD-вектора для фронтенда.
type PADDTO struct {
	Pleasure  float64 `json:"pleasure"`
	Arousal   float64 `json:"arousal"`
	Dominance float64 `json:"dominance"`
}

// EmotionDTO — активная дискретная эмоция.
type EmotionDTO struct {
	// Type — тип эмоции ("joy", "fear", "trust").
	Type string `json:"type"`

	// Intensity — сила от 0.0 до 1.0.
	Intensity float64 `json:"intensity"`

	// Trigger — причина возникновения ("met a new friend").
	Trigger string `json:"trigger"`
}

// GoalDTO — цель агента.
type GoalDTO struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Priority    float64 `json:"priority"`
	Progress    float64 `json:"progress"`
}

// AgentStatsDTO — агрегированная статистика агента.
type AgentStatsDTO struct {
	// TotalInteractions — общее число взаимодействий с другими агентами.
	TotalInteractions int `json:"totalInteractions"`

	// MemoriesCount — общее число воспоминаний (episodic + semantic).
	MemoriesCount int `json:"memoriesCount"`

	// RelationshipsCount — число связей с другими агентами.
	RelationshipsCount int `json:"relationshipsCount"`

	// DaysSinceCreation — возраст агента в днях симуляции.
	DaysSinceCreation int `json:"daysSinceCreation"`
}

// SpawnAgentRequest — запрос на создание агента, POST /api/v1/control/spawn.
type SpawnAgentRequest struct {
	Name        string         `json:"name" binding:"required"`
	Personality PersonalityDTO `json:"personality" binding:"required"`
}

// InjectThoughtRequest — инъекция мысли/памяти/цели в агента, POST /api/v1/agents/:id/inject.
type InjectThoughtRequest struct {
	// Type — тип инъекции: "thought", "memory", "goal".
	Type string `json:"type" binding:"required"`

	// Content — содержимое инъекции.
	Content string `json:"content" binding:"required"`
}

// =============================================================================
// MEMORY DTOs
// =============================================================================

// MemoryResponse — ответ на GET /api/v1/agents/:id/memory.
type MemoryResponse struct {
	Memories []MemoryEntryDTO `json:"memories"`

	// Summary — AI-сгенерированная сводка воспоминаний агента.
	Summary string `json:"summary"`
}

// MemoryEntryDTO — одно воспоминание.
type MemoryEntryDTO struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"` // "episodic" | "semantic" | "procedural"
	Content       string    `json:"content"`
	EmotionalTag  string    `json:"emotionalTag"`
	Importance    float64   `json:"importance"`
	Timestamp     time.Time `json:"timestamp"`
	RelatedAgents []string  `json:"relatedAgents"`
}

// =============================================================================
// RELATIONSHIP DTOs
// =============================================================================

// RelationshipGraphResponse — граф отношений, ответ на GET /api/v1/relationships.
// Формат совместим с библиотеками визуализации графов (D3.js, vis.js, Sigma.js).
type RelationshipGraphResponse struct {
	Nodes []GraphNodeDTO `json:"nodes"`
	Edges []GraphEdgeDTO `json:"edges"`
}

// GraphNodeDTO — узел графа (агент).
type GraphNodeDTO struct {
	// ID — UUID агента.
	ID string `json:"id"`

	// Label — отображаемое имя.
	Label string `json:"label"`

	// Type — тип личности для цветовой кодировки узла.
	Type string `json:"type"`

	// Size — размер узла, пропорциональный количеству связей.
	Size int `json:"size"`
}

// GraphEdgeDTO — ребро графа (связь между двумя агентами).
type GraphEdgeDTO struct {
	// Source — UUID агента-источника.
	Source string `json:"source"`

	// Target — UUID целевого агента.
	Target string `json:"target"`

	// Type — тип отношений: "friend", "rival", "neutral", "romantic".
	Type string `json:"type"`

	// Strength — сила связи: -1.0 (враждебная) до +1.0 (тесная).
	// Знак определяет цвет ребра, абсолютное значение — толщину.
	Strength float64 `json:"strength"`

	// Label — описание связи ("bonded over shared values").
	Label string `json:"label,omitempty"`
}

// CreateRelationshipRequest — создание/изменение связи, POST /api/v1/relationships.
type CreateRelationshipRequest struct {
	Agent1ID string `json:"agent1" binding:"required"`
	Agent2ID string `json:"agent2" binding:"required"`
	Type     string `json:"type" binding:"required"`
}

// =============================================================================
// EVENT DTOs
// =============================================================================

// EventDTO — мировое или агентное событие.
type EventDTO struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`     // "global", "agent", "interaction"
	Category       string    `json:"category"` // "disaster", "celebration", "discovery"
	Description    string    `json:"description"`
	AffectedAgents []string  `json:"affectedAgents"`
	Timestamp      time.Time `json:"timestamp"`
	Status         string    `json:"status"` // "pending", "active", "completed"
}

// InjectEventRequest — инъекция глобального события, POST /api/v1/events.
type InjectEventRequest struct {
	// Type — категория события: "disaster", "celebration", "discovery" и т.д.
	Type string `json:"type" binding:"required"`

	// Description — текстовое описание события.
	Description string `json:"description" binding:"required"`

	// AffectedAgents — кого затрагивает. Пустой массив или ["all"] = все агенты.
	AffectedAgents []string `json:"affectedAgents"`
}

// =============================================================================
// WORLD DTOs
// =============================================================================

// WorldStatusResponse — статус симуляции, ответ на GET /api/v1/world/status.
type WorldStatusResponse struct {
	// CurrentTick — текущий тик симуляции (монотонно возрастающий счётчик).
	CurrentTick int64 `json:"currentTick"`

	// SimulationSpeed — множитель скорости (1.0 = реальное время).
	SimulationSpeed float64 `json:"simulationSpeed"`

	// IsPaused — на паузе ли симуляция.
	IsPaused bool `json:"isPaused"`

	// ActiveAgents — количество активных агентов.
	ActiveAgents int `json:"activeAgents"`

	// TotalEvents — общее количество произошедших событий.
	TotalEvents int `json:"totalEvents"`

	// Uptime — время работы симуляции в формате "2h15m30s".
	Uptime string `json:"uptime"`
}

// WorldControlRequest — управление симуляцией, POST /api/v1/world/control.
type WorldControlRequest struct {
	// Action — действие: "pause", "resume", "step", "setSpeed".
	Action string `json:"action" binding:"required"`

	// Value — параметр действия (скорость для setSpeed).
	Value *float64 `json:"value,omitempty"`
}

// WorldStatisticsResponse — агрегированная статистика мира, GET /api/v1/world/statistics.
type WorldStatisticsResponse struct {
	// MoodDistribution — распределение настроений: {"happy": 3, "anxious": 1, ...}.
	MoodDistribution map[string]int `json:"moodDistribution"`

	// RelationshipStats — статистика графа отношений.
	RelationshipStats RelationshipStatsDTO `json:"relationshipStats"`

	// ActivityMetrics — метрики активности за последний период.
	ActivityMetrics ActivityMetricsDTO `json:"activityMetrics"`

	// TopInteractingAgents — самые общительные агенты.
	TopInteractingAgents []AgentSummary `json:"topInteractingAgents"`
}

// RelationshipStatsDTO — агрегированная статистика отношений.
type RelationshipStatsDTO struct {
	TotalConnections int     `json:"totalConnections"`
	AverageStrength  float64 `json:"averageStrength"`
	Rivalries        int     `json:"rivalries"`
}

// ActivityMetricsDTO — метрики активности.
type ActivityMetricsDTO struct {
	InteractionsLastHour   int `json:"interactionsLastHour"`
	MemoriesFormedLastHour int `json:"memoriesFormedLastHour"`
	EventsLastHour         int `json:"eventsLastHour"`
}

// ResetRequest — сброс мира, POST /api/v1/control/reset.
type ResetRequest struct {
	Confirm        bool `json:"confirm" binding:"required"`
	PreserveAgents bool `json:"preserveAgents"`
}

// SuccessResponse — универсальный ответ на мутирующие операции.
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
