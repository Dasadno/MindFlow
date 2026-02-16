// Package agent provides the Brain component for agent cognition.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// The Brain is the cognitive core of an agent. It interfaces with the LLM
// (GigaChat) to perform reasoning, decision-making, and natural language
// generation. It orchestrates the thinking process by combining personality,
// memory, and emotional context.
//
// =============================================================================
// BRAIN STRUCTURE:
// =============================================================================
//
// type Brain struct {
//     llmClient     *gigachat.Client  // GigaChat/Ollama client
//     personality   *Personality       // Reference to agent's personality
//     memory        *MemorySystem      // Access to memories for context
//     emotions      *EmotionEngine     // Current emotional state
//     thoughtBuffer []Thought          // Recent thoughts (working memory)
//     config        BrainConfig        // Thinking parameters
// }
//
// type BrainConfig struct {
//     MaxThoughts       int     // Working memory capacity
//     ReflectionDepth   int     // How many layers of meta-cognition
//     CreativityFactor  float64 // LLM temperature modifier
//     ResponseTimeout   time.Duration
// }
//
// type Thought struct {
//     Content   string
//     Type      ThoughtType  // observation, reasoning, decision, emotion
//     Timestamp time.Time
//     Triggers  []string     // What caused this thought
// }
//
// =============================================================================
// COGNITIVE FUNCTIONS:
// =============================================================================
//
// func (b *Brain) Think(context CognitiveContext) ThinkingResult
//   - Main reasoning function
//   - Builds prompt with: personality traits, relevant memories, current mood
//   - Calls LLM for reasoning
//   - Parses structured response (thoughts, conclusions, actions)
//
// func (b *Brain) DecideAction(options []PossibleAction, context CognitiveContext) Decision
//   - Evaluate possible actions given current context
//   - Consider: personality alignment, emotional state, goal progress
//   - Return chosen action with reasoning
//
// func (b *Brain) GenerateResponse(interaction Interaction) string
//   - Generate natural language response for interactions
//   - Maintain consistent personality voice
//   - Consider relationship with interacting agent
//
// func (b *Brain) Reflect(memories []MemoryEntry) ReflectionInsights
//   - Deep reflection on experiences
//   - Identify patterns in behavior and outcomes
//   - Generate insights about self and others
//   - May trigger goal updates or personality adjustments
//
// =============================================================================
// PROMPT ENGINEERING:
// =============================================================================
//
// func (b *Brain) buildSystemPrompt() string
//   - Construct agent's "inner voice" prompt
//   - Include: personality description, core values, quirks
//   - Define response format expectations
//
// func (b *Brain) buildContextPrompt(context CognitiveContext) string
//   - Include relevant memories (vector search for related content)
//   - Current emotional state and recent mood trajectory
//   - Active goals and their progress
//   - Recent interactions summary
//
// func (b *Brain) parseStructuredResponse(response string) CognitiveOutput
//   - Parse LLM response into structured data
//   - Extract: thoughts, decisions, emotional reactions
//   - Handle parsing errors gracefully
//
// =============================================================================
// INTERNAL MONOLOGUE:
// =============================================================================
//
// func (b *Brain) StreamThoughts() <-chan Thought
//   - Channel for real-time thought streaming (for dashboard)
//   - Emit thoughts as they are generated
//   - Include meta-thoughts about the thinking process

package agent
