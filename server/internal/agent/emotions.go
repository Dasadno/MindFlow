// Package agent provides the Emotion Engine for affective computing.
//
// The Emotion Engine models the agent's emotional state using a dimensional
// approach (PAD model: Pleasure-Arousal-Dominance) combined with discrete
// emotions. It influences decision-making, memory encoding, and interactions.
//
// Hybrid model:
//   1. PAD (Pleasure-Arousal-Dominance) — непрерывное пространство эмоций
//   2. Discrete emotions (Ekman + social emotions) — человеко-читаемые метки
//   3. Mood — долгосрочный эмоциональный фон (медленнее меняется, чем эмоции)

package agent

import "time"

// -----------------------------------------------------------------------------
// EmotionEngine — движок аффективных вычислений агента
// -----------------------------------------------------------------------------
// Три слоя PAD-состояния создают эмоциональную инерцию:
//   currentState  — мгновенное состояние (быстро меняется от стимулов)
//   moodBaseline  — настроение (медленно смещается к среднему от recent emotions)
//   personalityBias — личностный аттрактор (определяется Personality, не меняется)
//
// Стимулы быстро сдвигают currentState, но настроение (mood) реагирует медленно,
// а personalityBias задаёт долгосрочную точку притяжения — как в реальной психологии.

type EmotionEngine struct {
	// currentState — текущее мгновенное эмоциональное состояние в пространстве PAD.
	// Быстро меняется в ответ на стимулы, затем затухает обратно к moodBaseline.
	CurrentState PADState

	// moodBaseline — долгосрочное настроение. Смещается медленно, как скользящее
	// среднее от recent emotional states. Определяет «фоновое» самочувствие.
	MoodBaseline PADState

	// personalityBias — PAD-значения, к которым настроение стремится в покое.
	// Вычисляется из Personality при создании агента. Например, высокий Neuroticism
	// сдвигает bias в сторону низкого Pleasure и высокого Arousal.
	PersonalityBias PADState

	// activeEmotions — список дискретных эмоций, активных в данный момент.
	// Каждая имеет интенсивность и время жизни. Используются для:
	//   - человеко-читаемых меток на дашборде
	//   - контекста в LLM-промптах ("I'm feeling joy because...")
	ActiveEmotions []DiscreteEmotion

	// history — история эмоциональных снапшотов для визуализации трендов на дашборде.
	History []EmotionSnapshot

	// decayRate — скорость затухания эмоций обратно к baseline (0.0–1.0 за тик).
	// Высокий нейротизм → медленное затухание (эмоции дольше держатся).
	DecayRate float64

	// config — настройки движка.
	Config EmotionConfig
}

// EmotionConfig — конфигурация движка эмоций.
type EmotionConfig struct {
	// MaxActiveEmotions — максимальное количество одновременно активных дискретных эмоций.
	// При превышении самые слабые удаляются.
	MaxActiveEmotions int

	// HistorySize — сколько снапшотов хранить для графика трендов.
	HistorySize int

	// MoodInertia — инерция настроения: чем выше, тем медленнее mood следует за emotions.
	// Диапазон 0.0 (мгновенная реакция) — 1.0 (настроение практически не меняется).
	MoodInertia float64

	// MinIntensityThreshold — минимальная интенсивность, ниже которой эмоция удаляется.
	MinIntensityThreshold float64
}

// -----------------------------------------------------------------------------
// PADState — непрерывное эмоциональное пространство (Pleasure-Arousal-Dominance)
// -----------------------------------------------------------------------------
// Модель Мехрабяна-Рассела. Три оси описывают любое эмоциональное состояние:
//   Pleasure  — валентность (насколько хорошо/плохо себя чувствует)
//   Arousal   — активация (спокоен или возбуждён)
//   Dominance — доминирование (чувствует контроль или подчинение)

type PADState struct {
	// Pleasure — ось удовольствия. -1.0 (несчастье, страдание) до +1.0 (счастье, эйфория).
	Pleasure float64 `json:"pleasure"`

	// Arousal — ось возбуждения. -1.0 (сонливость, апатия) до +1.0 (ажитация, паника).
	Arousal float64 `json:"arousal"`

	// Dominance — ось доминирования. -1.0 (бессилие, подчинение) до +1.0 (контроль, власть).
	Dominance float64 `json:"dominance"`
}

// -----------------------------------------------------------------------------
// DiscreteEmotion — конкретная именованная эмоция с интенсивностью
// -----------------------------------------------------------------------------
// Дискретные эмоции вычисляются из PAD-состояния и результатов AppraisalCheck.
// Они более понятны человеку, чем числа PAD, и используются в промптах и UI.

type DiscreteEmotion struct {
	// Type — тип эмоции (joy, sadness, anger и т.д.)
	Type EmotionType

	// Intensity — сила эмоции от 0.0 (едва заметна) до 1.0 (подавляющая).
	Intensity float64

	// Trigger — что вызвало эту эмоцию ("met a new friend", "lost a debate").
	// Сохраняется в памяти и используется в промптах Brain.
	Trigger string

	// StartTime — когда эмоция возникла. Используется для расчёта длительности.
	StartTime time.Time

	// Duration — ожидаемая длительность эмоции. По истечении — затухание.
	Duration time.Duration
}

// EmotionType — перечисление типов дискретных эмоций.
// Включает 6 базовых эмоций Экмана + социальные и самосознательные эмоции.
type EmotionType string

const (
	EmotionJoy          EmotionType = "joy"          // Радость — высокий Pleasure, средний Arousal
	EmotionSadness      EmotionType = "sadness"      // Грусть — низкий Pleasure, низкий Arousal
	EmotionAnger        EmotionType = "anger"        // Гнев — низкий Pleasure, высокий Arousal, высокий Dominance
	EmotionFear         EmotionType = "fear"         // Страх — низкий Pleasure, высокий Arousal, низкий Dominance
	EmotionSurprise     EmotionType = "surprise"     // Удивление — нейтральный Pleasure, высокий Arousal
	EmotionDisgust      EmotionType = "disgust"      // Отвращение — низкий Pleasure, средний Arousal
	EmotionTrust        EmotionType = "trust"        // Доверие — социальная эмоция, укрепляет связи
	EmotionAnticipation EmotionType = "anticipation" // Предвкушение — средний Pleasure, средний Arousal
	EmotionLoneliness   EmotionType = "loneliness"   // Одиночество — социальная, низкий Pleasure, низкий Arousal
	EmotionPride        EmotionType = "pride"        // Гордость — самосознательная, высокий Dominance
	EmotionShame        EmotionType = "shame"        // Стыд — самосознательная, низкий Dominance
)

// -----------------------------------------------------------------------------
// Mood — дискретная метка настроения для UI
// -----------------------------------------------------------------------------
// Вычисляется из PADState путём маппинга областей PAD-пространства
// на человеко-читаемые метки. Используется на дашборде для быстрого
// понимания состояния агента без погружения в числа PAD.

type Mood string

const (
	MoodHappy   Mood = "happy"   // P > 0.3, A > 0.0
	MoodSad     Mood = "sad"     // P < -0.3, A < 0.0
	MoodAnxious Mood = "anxious" // P < 0.0, A > 0.5, D < 0.0
	MoodCalm    Mood = "calm"    // P > 0.0, A < -0.3
	MoodAngry   Mood = "angry"   // P < -0.3, A > 0.3, D > 0.3
	MoodExcited Mood = "excited" // P > 0.3, A > 0.5
	MoodBored   Mood = "bored"   // P ~ 0, A < -0.5
	MoodContent Mood = "content" // P > 0.2, A ~ 0, D > 0.0
	MoodNeutral Mood = "neutral" // Все оси близки к 0
)

// -----------------------------------------------------------------------------
// EmotionSnapshot — снимок эмоционального состояния для истории
// -----------------------------------------------------------------------------
// Сохраняется периодически (каждый тик или каждые N тиков).
// Массив снапшотов отправляется в API для отрисовки графиков настроения.

type EmotionSnapshot struct {
	// State — PAD-состояние в момент снимка.
	State PADState

	// DominantEmotion — самая сильная дискретная эмоция в момент снимка.
	DominantEmotion EmotionType

	// Mood — дискретная метка настроения.
	Mood Mood

	// Tick — номер тика симуляции.
	Tick int64

	// Timestamp — время снимка.
	Timestamp time.Time
}

// -----------------------------------------------------------------------------
// MoodInfluence — как текущее настроение влияет на поведение агента
// -----------------------------------------------------------------------------
// Возвращается из GetMoodInfluence(). Brain использует эти модификаторы
// при принятии решений и генерации ответов.

type MoodInfluence struct {
	// ImpulsivityModifier — модификатор импульсивности (-1.0 .. +1.0).
	// Высокий Arousal → больше импульсивности → спонтанные решения.
	ImpulsivityModifier float64

	// SociabilityModifier — модификатор общительности (-1.0 .. +1.0).
	// Высокий Pleasure → больше желания общаться.
	SociabilityModifier float64

	// RiskModifier — модификатор склонности к риску (-1.0 .. +1.0).
	// Низкий Pleasure → избегание риска, осторожность.
	RiskModifier float64

	// AssertivenessModifier — модификатор напористости (-1.0 .. +1.0).
	// Высокий Dominance → более напористые взаимодействия.
	AssertivenessModifier float64
}

// -----------------------------------------------------------------------------
// AppraisalResult — результат когнитивной оценки стимула (теория Лазаруса)
// -----------------------------------------------------------------------------
// EmotionEngine.AppraisalCheck() оценивает входящий стимул по пяти критериям.
// Результат определяет, КАКИЕ эмоции возникнут и с какой интенсивностью.
// Это делает эмоциональные реакции зависимыми от личности и контекста.

type AppraisalResult struct {
	// Novelty — насколько стимул неожиданный (0.0 = ожидаемый, 1.0 = полностью новый).
	// Высокая новизна → Surprise, повышенный Arousal.
	Novelty float64

	// GoalRelevance — насколько стимул связан с целями агента (0.0 = нерелевантен, 1.0 = критичен).
	// Нерелевантные стимулы вызывают слабый эмоциональный отклик.
	GoalRelevance float64

	// GoalCongruence — помогает (+1.0) или мешает (-1.0) стимул достижению целей.
	// Конгруэнтные → Joy/Trust, неконгруэнтные → Anger/Fear/Sadness.
	GoalCongruence float64

	// Agency — кто вызвал стимул: "self", "other", "situation".
	// Определяет тип эмоции: self → Pride/Shame, other → Anger/Trust.
	Agency string

	// NormCompatibility — совместимость с ценностями агента (-1.0 .. +1.0).
	// Несовместимость с CoreValues → Disgust, совместимость → Trust.
	NormCompatibility float64
}
