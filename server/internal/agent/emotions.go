// Package agent provides the Emotion Engine for affective computing.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// The Emotion Engine models the agent's emotional state using a dimensional
// approach (PAD model: Pleasure-Arousal-Dominance) combined with discrete
// emotions. It influences decision-making, memory encoding, and interactions.
//
// =============================================================================
// EMOTIONAL MODEL:
// =============================================================================
//
// We use a hybrid model combining:
// 1. PAD (Pleasure-Arousal-Dominance) continuous dimensions
// 2. Discrete emotions (Ekman's basic emotions + social emotions)
// 3. Mood (longer-term emotional baseline)
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                      EMOTION ENGINE                                  │
// ├─────────────────────────────────────────────────────────────────────┤
// │                                                                      │
// │   CURRENT STATE          MOOD BASELINE          PERSONALITY BIAS    │
// │   ┌───────────┐          ┌───────────┐          ┌───────────┐       │
// │   │ Pleasure  │◄────────►│ Baseline  │◄────────►│ Neuroticism│      │
// │   │ Arousal   │          │ Mood      │          │ Extraversion│     │
// │   │ Dominance │          │           │          │             │     │
// │   └───────────┘          └───────────┘          └───────────────┘   │
// │         │                      │                        │           │
// │         └──────────────────────┼────────────────────────┘           │
// │                                ▼                                    │
// │                    ┌───────────────────────┐                        │
// │                    │   Discrete Emotions   │                        │
// │                    │ joy, sadness, anger,  │                        │
// │                    │ fear, surprise, etc.  │                        │
// │                    └───────────────────────┘                        │
// └─────────────────────────────────────────────────────────────────────┘
//
// =============================================================================
// DATA STRUCTURES:
// =============================================================================
//
// type EmotionEngine struct {
//     currentState   PADState           // Current emotional state
//     moodBaseline   PADState           // Longer-term mood (slower to change)
//     activeEmotions []DiscreteEmotion  // Currently active discrete emotions
//     personalityBias PADState          // Personality influence on baseline
//     history        []EmotionSnapshot  // Emotional history for trending
//     decayRate      float64            // How fast emotions return to baseline
//     config         EmotionConfig
// }
//
// type PADState struct {
//     Pleasure  float64 // -1 (unhappy) to +1 (happy)
//     Arousal   float64 // -1 (calm) to +1 (excited)
//     Dominance float64 // -1 (submissive) to +1 (dominant)
// }
//
// type DiscreteEmotion struct {
//     Type      EmotionType // joy, sadness, anger, fear, surprise, disgust, etc.
//     Intensity float64     // 0 to 1
//     Trigger   string      // What caused this emotion
//     StartTime time.Time
//     Duration  time.Duration
// }
//
// type EmotionType string
// const (
//     EmotionJoy       EmotionType = "joy"
//     EmotionSadness   EmotionType = "sadness"
//     EmotionAnger     EmotionType = "anger"
//     EmotionFear      EmotionType = "fear"
//     EmotionSurprise  EmotionType = "surprise"
//     EmotionDisgust   EmotionType = "disgust"
//     EmotionTrust     EmotionType = "trust"      // Social emotion
//     EmotionAnticipation EmotionType = "anticipation"
//     EmotionLoneliness EmotionType = "loneliness" // Social emotion
//     EmotionPride     EmotionType = "pride"      // Self-conscious emotion
//     EmotionShame     EmotionType = "shame"      // Self-conscious emotion
// )
//
// =============================================================================
// EMOTION PROCESSING:
// =============================================================================
//
// func (e *EmotionEngine) ProcessStimulus(stimulus Stimulus) EmotionalResponse
//   - Evaluate stimulus based on agent's goals, values, personality
//   - Calculate emotional impact (appraisal theory)
//   - Update PAD state and trigger discrete emotions
//   - Return emotional response for behavior influence
//
// func (e *EmotionEngine) AppraisalCheck(stimulus Stimulus) AppraisalResult
//   - Novelty check: Is this expected or surprising?
//   - Goal relevance: Does this affect my goals?
//   - Goal congruence: Does this help or hinder my goals?
//   - Agency: Who caused this? (self, other, situation)
//   - Norm compatibility: Does this align with my values?
//
// func (e *EmotionEngine) Update(deltaTime time.Duration)
//   - Called each tick
//   - Decay active emotions toward baseline
//   - Update mood based on recent emotional average
//   - Expire short-lived emotions
//
// =============================================================================
// MOOD MANAGEMENT:
// =============================================================================
//
// func (e *EmotionEngine) GetMood() Mood
//   - Calculate current mood from PAD state
//   - Map to discrete mood labels for easy interpretation
//
// type Mood string
// const (
//     MoodHappy     Mood = "happy"
//     MoodSad       Mood = "sad"
//     MoodAnxious   Mood = "anxious"
//     MoodCalm      Mood = "calm"
//     MoodAngry     Mood = "angry"
//     MoodExcited   Mood = "excited"
//     MoodBored     Mood = "bored"
//     MoodContent   Mood = "content"
//     MoodNeutral   Mood = "neutral"
// )
//
// func (e *EmotionEngine) GetMoodInfluence() MoodInfluence
//   - Return how current mood should influence behavior
//   - High arousal → more impulsive decisions
//   - Low pleasure → risk-averse, withdrawal behavior
//   - High dominance → assertive interactions
//
// =============================================================================
// EMOTIONAL HISTORY:
// =============================================================================
//
// func (e *EmotionEngine) GetEmotionalTrend(duration time.Duration) []EmotionSnapshot
//   - Return emotional history for visualization
//   - Used by dashboard for mood graphs
//
// func (e *EmotionEngine) GetDominantEmotions(n int) []DiscreteEmotion
//   - Return n strongest active emotions
//   - Used for decision-making context

package agent
