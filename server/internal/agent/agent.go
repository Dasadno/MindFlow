// Package agent provides the core AI agent implementation for the society simulation.
//
// This file defines the Agent entity — the central actor in the simulation.
// An Agent is an autonomous AI entity with personality, emotions, memory,
// and the ability to make decisions and interact with other agents.

package agent

import "time"

// -----------------------------------------------------------------------------
// Agent — центральная сущность симуляции
// -----------------------------------------------------------------------------
// Каждый агент — автономная «личность» с собственным мозгом, памятью, эмоциями
// и целями. Оркестратор (world.Orchestrator) вызывает Tick() каждый тик
// симуляции, запуская когнитивный цикл: восприятие → мышление → решение → действие.

type Agent struct {
	// ID — уникальный идентификатор агента (UUID v4).
	ID string

	// Name — отображаемое имя агента ("Alice", "Bob").
	Name string

	// Personality — черты личности по модели Big Five (OCEAN) + ценности и причуды.
	// Определяет, КАК агент думает и реагирует на стимулы.
	Personality *Personality

	// Brain — когнитивное ядро, обёртка над LLM (GigaChat).
	// Отвечает за мышление, принятие решений и генерацию речи.
	Brain *Brain

	// Emotions — движок эмоций (PAD-модель + дискретные эмоции).
	// Влияет на принятие решений, формирование памяти и стиль взаимодействия.
	Emotions *EmotionEngine

	// Goals — менеджер целей: текущие задачи и мотивации агента.
	// Цели формируются при создании и обновляются в процессе рефлексии.
	Goals []Goal

	// State — текущее состояние активности агента в симуляции.
	State AgentState

	// CreatedAt — момент создания агента.
	CreatedAt time.Time

	// LastActive — время последнего тика, в котором агент участвовал.
	LastActive time.Time
}

// AgentState — перечисление возможных состояний агента в симуляции.
// Состояние определяет, какие действия доступны и как агент отображается на дашборде.
type AgentState string

const (
	StateIdle        AgentState = "idle"        // Бездействует, ожидает следующего тика
	StateThinking    AgentState = "thinking"    // Обрабатывает когнитивный цикл (LLM-вызов)
	StateActing      AgentState = "acting"      // Выполняет выбранное действие
	StateInteracting AgentState = "interacting" // Ведёт диалог с другим агентом
	StateReflecting  AgentState = "reflecting"  // Проводит мета-когнитивную рефлексию
	StateSleeping    AgentState = "sleeping"    // «Спит» — период консолидации памяти
)

// -----------------------------------------------------------------------------
// Personality — модель личности Big Five (OCEAN)
// -----------------------------------------------------------------------------
// Каждый trait — число от 0.0 до 1.0. Определяет поведенческие паттерны агента:
// как он реагирует на стимулы, какие цели ставит, как общается.

type Personality struct {
	// Openness — открытость опыту.
	// Высокая: любопытство, креативность, склонность к экспериментам.
	// Низкая: практичность, предпочтение привычного.
	Openness float64

	// Conscientiousness — добросовестность.
	// Высокая: организованность, дисциплина, планирование.
	// Низкая: спонтанность, гибкость, импульсивность.
	Conscientiousness float64

	// Extraversion — экстраверсия.
	// Высокая: общительность, энергичность, инициирует взаимодействия.
	// Низкая: интроверсия, предпочитает одиночество и наблюдение.
	Extraversion float64

	// Agreeableness — доброжелательность.
	// Высокая: кооперативность, доверчивость, эмпатия.
	// Низкая: конкурентность, скептицизм, прямолинейность.
	Agreeableness float64

	// Neuroticism — нейротизм.
	// Высокая: эмоциональная нестабильность, тревожность, чувствительность к стрессу.
	// Низкая: эмоциональная устойчивость, спокойствие, стрессоустойчивость.
	Neuroticism float64

	// CoreValues — базовые ценности агента (например, "honesty", "creativity", "justice").
	// Влияют на AppraisalCheck() в EmotionEngine — стимулы оцениваются через призму ценностей.
	CoreValues []string

	// Quirks — уникальные поведенческие особенности ("talks to herself", "collects stones").
	// Добавляют в system prompt для LLM, придавая агенту индивидуальность.
	Quirks []string
}

// -----------------------------------------------------------------------------
// Goal — цель агента
// -----------------------------------------------------------------------------
// Цели формируются при создании (на основе Personality) и обновляются
// во время рефлексии (Reflect). Используются при принятии решений —
// Brain оценивает, какое действие приближает к цели.

type Goal struct {
	// ID — уникальный идентификатор цели.
	ID string

	// Description — текстовое описание цели ("Make a new friend", "Explore the eastern zone").
	Description string

	// Priority — приоритет от 0.0 до 1.0. Высокий приоритет → больше влияния на решения.
	Priority float64

	// Progress — прогресс выполнения от 0.0 до 1.0.
	Progress float64

	// IsCompleted — завершена ли цель. Завершённые цели сохраняются для рефлексии.
	IsCompleted bool

	// CreatedAt — когда цель была поставлена.
	CreatedAt time.Time
}

// -----------------------------------------------------------------------------
// Stimulus — внешнее воздействие на агента
// -----------------------------------------------------------------------------
// Стимулы приходят от мировых событий, других агентов или API-инъекций.
// Каждый стимул проходит через EmotionEngine.AppraisalCheck() для оценки
// эмоционального воздействия, затем может быть сохранён в памяти.

type Stimulus struct {
	// Type — тип стимула: событие, сообщение от агента, инъекция извне.
	Type StimulusType

	// Source — откуда пришёл стимул (agent ID, "system", "api").
	Source string

	// Content — содержимое стимула в свободной текстовой форме.
	Content string

	// Intensity — сила воздействия от 0.0 до 1.0. Влияет на эмоциональный отклик.
	Intensity float64

	// Timestamp — когда стимул произошёл.
	Timestamp time.Time
}

// StimulusType — классификация входящих стимулов.
type StimulusType string

const (
	StimulusEvent       StimulusType = "event"       // Мировое или локальное событие
	StimulusMessage     StimulusType = "message"     // Сообщение от другого агента
	StimulusInjection   StimulusType = "injection"   // Инъекция мысли/памяти через API
	StimulusEnvironment StimulusType = "environment" // Изменение окружения
)

// -----------------------------------------------------------------------------
// AgentAction — действие, выбранное агентом в текущем тике
// -----------------------------------------------------------------------------
// Возвращается из Agent.Tick(). Оркестратор обрабатывает действия всех агентов,
// разрешает конфликты и сопоставляет запросы на взаимодействие.

type AgentAction struct {
	// AgentID — кто совершает действие.
	AgentID string

	// Type — тип действия.
	Type ActionType

	// TargetAgentID — ID целевого агента (для взаимодействий). Пустой, если действие одиночное.
	TargetAgentID string

	// Description — текстовое описание действия, сгенерированное Brain.
	Description string

	// Intent — намерение взаимодействия (только для Interact-действий).
	Intent InteractionIntent
}

// ActionType — классификация действий агента.
type ActionType string

const (
	ActionIdle     ActionType = "idle"     // Ничего не делать
	ActionThink    ActionType = "think"    // Внутреннее размышление
	ActionInteract ActionType = "interact" // Начать взаимодействие с другим агентом
	ActionExplore  ActionType = "explore"  // Исследовать окружение
	ActionReflect  ActionType = "reflect"  // Провести рефлексию
)

// InteractionIntent — намерение агента при инициации взаимодействия.
type InteractionIntent string

const (
	IntentChat     InteractionIntent = "chat"     // Дружеский разговор
	IntentDebate   InteractionIntent = "debate"   // Спор, обсуждение
	IntentHelp     InteractionIntent = "help"     // Предложить помощь
	IntentAsk      InteractionIntent = "ask"      // Попросить о чём-то
	IntentConflict InteractionIntent = "conflict" // Конфликт, конфронтация
)

// -----------------------------------------------------------------------------
// AgentSnapshot — сериализуемое состояние агента для персистенции
// -----------------------------------------------------------------------------
// Сохраняется в БД (agents.snapshot). Позволяет полностью восстановить агента
// после перезапуска сервера с сохранением личности, настроения и целей.

type AgentSnapshot struct {
	// ID и Name — идентификация.
	ID   string `json:"id"`
	Name string `json:"name"`

	// Personality — полный набор черт личности (сериализуется в JSON).
	Personality *Personality `json:"personality"`

	// MoodState — текущее эмоциональное состояние (PAD).
	MoodState *PADState `json:"moodState"`

	// Goals — список активных целей.
	Goals []Goal `json:"goals"`

	// State — последнее состояние агента.
	State AgentState `json:"state"`

	// MemorySummary — краткое описание ключевых воспоминаний для быстрого восстановления.
	MemorySummary string `json:"memorySummary"`
}

// -----------------------------------------------------------------------------
// WorldContext — контекст мира, передаваемый агенту каждый тик
// -----------------------------------------------------------------------------
// Оркестратор формирует WorldContext и передаёт в Agent.Tick().
// Содержит всё, что агент «видит» в текущий момент симуляции.

type WorldContext struct {
	// CurrentTick — номер текущего тика симуляции.
	CurrentTick int64

	// NearbyAgents — список агентов, доступных для взаимодействия.
	NearbyAgents []AgentSummary

	// ActiveEvents — мировые события, происходящие в данный момент.
	ActiveEvents []string

	// SimTime — текущее время симуляции.
	SimTime time.Time
}

// AgentSummary — краткая информация об агенте для WorldContext.
// Другие агенты видят только публичную информацию, не внутреннее состояние.
type AgentSummary struct {
	ID          string
	Name        string
	CurrentMood string
	State       AgentState
}
