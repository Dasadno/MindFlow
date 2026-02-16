// Package world provides the EventBus for inter-agent and world communication.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// The EventBus is a publish-subscribe message system that enables decoupled
// communication between agents, the orchestrator, and external systems (API).
// It is the backbone of multi-agent interaction — agents don't call each other
// directly; they publish and subscribe to events.
//
// =============================================================================
// EVENT BUS STRUCTURE:
// =============================================================================
//
// type EventBus struct {
//     subscribers map[EventTopic][]Subscriber  // Topic → list of subscribers
//     eventLog    []WorldEvent                  // History of all events
//     queue       chan WorldEvent               // Buffered event queue
//     mu          sync.RWMutex
// }
//
// type Subscriber struct {
//     ID       string                          // Subscriber identifier (agent ID)
//     Channel  chan WorldEvent                  // Delivery channel
//     Filter   func(WorldEvent) bool           // Optional event filter
// }
//
// type EventTopic string
// const (
//     TopicGlobal        EventTopic = "global"          // World-wide events
//     TopicInteraction   EventTopic = "interaction"     // Agent-to-agent interactions
//     TopicMoodChange    EventTopic = "mood_change"     // Agent mood transitions
//     TopicGoalUpdate    EventTopic = "goal_update"     // Agent goal changes
//     TopicMemory        EventTopic = "memory"          // Memory formation events
//     TopicRelationship  EventTopic = "relationship"    // Relationship changes
//     TopicSystem        EventTopic = "system"          // System-level events (pause, resume)
// )
//
// =============================================================================
// WORLD EVENT STRUCTURE:
// =============================================================================
//
// type WorldEvent struct {
//     ID             string            // Unique event identifier
//     Topic          EventTopic        // Routing topic
//     Type           string            // Specific event type within topic
//     Source         string            // Who/what generated (agent ID, "system", "api")
//     AffectedAgents []string          // Target agents (empty = broadcast to all)
//     Payload        map[string]any    // Event-specific data
//     Timestamp      time.Time
//     Tick           int64             // Simulation tick when event occurred
//     Priority       int               // 0 = normal, higher = more urgent
// }
//
// =============================================================================
// EVENT BUS OPERATIONS:
// =============================================================================
//
// func NewEventBus(bufferSize int) *EventBus
//   - Create event bus with buffered channel of given size
//   - Initialize subscriber maps for all topics
//
// func (eb *EventBus) Subscribe(topic EventTopic, subscriberID string, filter func(WorldEvent) bool) <-chan WorldEvent
//   - Register subscriber for a topic
//   - Optional filter function to receive only matching events
//   - Return receive-only channel for event delivery
//   - Thread-safe: multiple goroutines can subscribe concurrently
//
// func (eb *EventBus) Unsubscribe(topic EventTopic, subscriberID string)
//   - Remove subscriber from topic
//   - Close subscriber's delivery channel
//
// func (eb *EventBus) Publish(event WorldEvent)
//   - Queue event for delivery
//   - Non-blocking: drops events if queue is full (with warning log)
//
// func (eb *EventBus) processEvents(ctx context.Context)
//   - Background goroutine consuming from queue
//   - Route each event to matching subscribers
//   - Respect AffectedAgents filter (deliver only to listed agents)
//   - Append to event log for history
//
// =============================================================================
// EVENT HISTORY & QUERIES:
// =============================================================================
//
// func (eb *EventBus) GetHistory(limit int) []WorldEvent
//   - Return most recent events from log
//   - Used by API for event listing
//
// func (eb *EventBus) GetEventsByTopic(topic EventTopic, limit int) []WorldEvent
//   - Filter event history by topic
//
// func (eb *EventBus) GetEventsByAgent(agentID string, limit int) []WorldEvent
//   - Get events where agent was source or affected
//
// =============================================================================
// SSE BRIDGE:
// =============================================================================
//
// func (eb *EventBus) StreamEvents() <-chan WorldEvent
//   - Special subscriber for the SSE endpoint (GET /api/v1/events/stream)
//   - Receives ALL events for real-time dashboard updates
//   - Separate from agent subscriptions to avoid interference

package world
