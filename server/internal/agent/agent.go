// Package agent provides the core AI agent implementation for the society simulation.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This file defines the Agent entity - the central actor in the simulation.
// An Agent is an autonomous AI entity with personality, emotions, memory,
// and the ability to make decisions and interact with other agents.
//
// =============================================================================
// AGENT STRUCTURE:
// =============================================================================
//
// type Agent struct {
//     ID          string
//     Name        string
//     Personality *Personality      // Big Five traits + values
//     Brain       *Brain            // Decision-making and reflection
//     Memory      *MemorySystem     // Episodic + Semantic memory
//     Emotions    *EmotionEngine    // Mood tracking and emotional responses
//     Goals       *GoalManager      // Current objectives and motivations
//     Social      *SocialModule     // Relationship tracking
//     State       AgentState        // Current activity state
//     CreatedAt   time.Time
//     LastActive  time.Time
// }
//
// type AgentState string
// const (
//     StateIdle       AgentState = "idle"
//     StateThinking   AgentState = "thinking"
//     StateActing     AgentState = "acting"
//     StateInteracting AgentState = "interacting"
//     StateSleeping   AgentState = "sleeping"
// )
//
// =============================================================================
// AGENT LIFECYCLE METHODS:
// =============================================================================
//
// func NewAgent(name string, personality *Personality, deps AgentDeps) *Agent
//   - Create new agent with initial state
//   - Initialize all subsystems (brain, memory, emotions, goals)
//   - Generate initial goals based on personality
//
// func (a *Agent) Tick(worldContext WorldContext) AgentAction
//   - Main update loop called by world orchestrator
//   - Process sensory input from world context
//   - Run cognitive cycle: perceive → think → decide → act
//   - Return chosen action for execution
//
// func (a *Agent) ReceiveStimulus(stimulus Stimulus) Response
//   - Handle external stimuli (events, messages, interactions)
//   - Update emotional state based on stimulus
//   - Store in episodic memory if significant
//   - Generate appropriate response
//
// func (a *Agent) Reflect() ReflectionResult
//   - Periodic self-reflection (meta-cognition)
//   - Review recent memories and experiences
//   - Update goals based on experiences
//   - Consolidate memories (episodic → semantic)
//
// =============================================================================
// INTER-AGENT INTERACTION:
// =============================================================================
//
// func (a *Agent) InitiateInteraction(target *Agent, intent InteractionIntent) Interaction
//   - Start conversation or action with another agent
//   - Consider relationship history and current moods
//
// func (a *Agent) RespondToInteraction(interaction Interaction) Response
//   - Process incoming interaction from another agent
//   - Generate contextually appropriate response
//   - Update relationship based on interaction outcome
//
// =============================================================================
// SERIALIZATION:
// =============================================================================
//
// func (a *Agent) ToSnapshot() AgentSnapshot
//   - Serialize agent state for persistence
//   - Include: personality, current mood, goals, memory summary
//
// func AgentFromSnapshot(snapshot AgentSnapshot, deps AgentDeps) *Agent
//   - Restore agent from persisted state
//   - Reinitialize runtime components

package agent
