// Package world provides the simulation orchestrator and world management.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// The Orchestrator is the central coordinator of the AI Agent Society simulation.
// It manages the world tick loop, coordinates agent actions, resolves conflicts,
// and maintains global world state. Think of it as the "game master".
//
// =============================================================================
// ORCHESTRATOR STRUCTURE:
// =============================================================================
//
// type Orchestrator struct {
//     agents       map[string]*agent.Agent   // All registered agents
//     eventBus     *EventBus                 // Inter-agent and world event system
//     timeManager  *TimeManager              // Simulation clock
//     storage      *storage.Repository       // Persistence layer
//     llmClient    *gigachat.Client           // For world narration and event generation
//     state        WorldState                 // Current simulation state
//     mu           sync.RWMutex               // Concurrency protection
//     ctx          context.Context            // Cancellation support
//     cancel       context.CancelFunc
// }
//
// type WorldState struct {
//     IsPaused     bool
//     CurrentTick  int64
//     Speed        float64      // Simulation speed multiplier (1.0 = real-time)
//     ActiveAgents int
//     StartedAt    time.Time
// }
//
// =============================================================================
// ORCHESTRATOR LIFECYCLE:
// =============================================================================
//
// func NewOrchestrator(deps OrchestratorDeps) *Orchestrator
//   - Initialize with dependencies (storage, event bus, GigaChat client)
//   - Set default world state (paused, tick 0, speed 1.0)
//   - Create internal communication channels
//
// func (o *Orchestrator) Start(ctx context.Context)
//   - Begin the main simulation loop in a goroutine
//   - Each tick:
//     1. Advance time via TimeManager
//     2. Collect world context (nearby agents, active events)
//     3. Call agent.Tick() for each active agent (concurrently)
//     4. Resolve agent actions (handle conflicts, interactions)
//     5. Process pending events from EventBus
//     6. Trigger periodic tasks (memory consolidation, mood decay)
//     7. Persist state snapshots at configured intervals
//
// func (o *Orchestrator) Stop()
//   - Cancel context to stop simulation loop
//   - Wait for in-flight agent ticks to complete
//   - Flush all pending state to storage
//
// =============================================================================
// AGENT MANAGEMENT:
// =============================================================================
//
// func (o *Orchestrator) RegisterAgent(a *agent.Agent)
//   - Add agent to the simulation
//   - Subscribe agent to relevant EventBus topics
//   - Emit "agent_joined" world event
//
// func (o *Orchestrator) RemoveAgent(agentID string)
//   - Soft-remove agent (mark inactive, preserve history)
//   - Unsubscribe from EventBus
//   - Emit "agent_left" world event
//
// func (o *Orchestrator) GetAgent(agentID string) *agent.Agent
//   - Retrieve agent by ID (thread-safe)
//
// func (o *Orchestrator) ListAgents() []*agent.Agent
//   - Return all active agents
//
// =============================================================================
// ACTION RESOLUTION:
// =============================================================================
//
// func (o *Orchestrator) resolveActions(actions []AgentAction) []ActionResult
//   - Process all agent actions from current tick
//   - Handle interaction requests: match initiator with target
//   - Resolve conflicts: two agents competing for same resource/action
//   - Generate interaction contexts for paired agents
//   - Return results for each action (success, failure, modified)
//
// func (o *Orchestrator) processInteraction(a1, a2 *agent.Agent, intent InteractionIntent) InteractionResult
//   - Facilitate multi-turn conversation between two agents
//   - Each agent generates responses through their Brain
//   - Track emotional shifts during interaction
//   - Store interaction as memory for both agents
//   - Update relationship graph based on outcome
//
// =============================================================================
// WORLD EVENTS:
// =============================================================================
//
// func (o *Orchestrator) InjectEvent(event WorldEvent)
//   - Inject external event (from API or scheduled)
//   - Broadcast to affected agents via EventBus
//   - Each agent processes event through ReceiveStimulus()
//
// func (o *Orchestrator) GenerateRandomEvent() WorldEvent
//   - Use LLM to generate contextually appropriate random events
//   - Consider current world state, agent moods, relationship tensions
//   - Add variety: weather changes, discoveries, conflicts, celebrations
//
// =============================================================================
// SIMULATION CONTROL:
// =============================================================================
//
// func (o *Orchestrator) Pause()
// func (o *Orchestrator) Resume()
// func (o *Orchestrator) SetSpeed(multiplier float64)
// func (o *Orchestrator) Step()    // Execute single tick while paused
// func (o *Orchestrator) GetStatus() WorldState

package world
