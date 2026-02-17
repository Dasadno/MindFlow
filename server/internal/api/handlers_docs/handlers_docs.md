// Package api provides HTTP handlers for the AI Agent Society dashboard.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This file contains all HTTP handler implementations for the REST API.
// Handlers translate HTTP requests into service calls and format responses.
//
// =============================================================================
// AGENT HANDLERS:
// =============================================================================
//
// GET  /api/v1/agents
//   - List all agents with basic info (id, name, personality type, current mood)
//   - Support pagination: ?page=1&limit=20
//   - Support filtering: ?personality=introvert&mood=happy
//
// GET  /api/v1/agents/:id
//   - Get detailed agent profile
//   - Include: personality traits, current goals, mood history, memory summary
//
// GET  /api/v1/agents/:id/memory
//   - Retrieve agent's episodic memories
//   - Support: ?type=episodic|semantic&limit=50
//   - Return memory entries with timestamps and emotional tags
//
// GET  /api/v1/agents/:id/thoughts
//   - Stream agent's current thought process (SSE endpoint)
//   - Real-time reflection and decision-making visibility
//
// POST /api/v1/agents/:id/inject
//   - Inject a thought or stimulus into agent's mind
//   - Body: { "type": "thought|memory|goal", "content": "..." }
//
// =============================================================================
// RELATIONSHIP HANDLERS:
// =============================================================================
//
// GET  /api/v1/relationships
//   - Get full relationship graph (nodes = agents, edges = relationships)
//   - Response format suitable for graph visualization libraries
//
// GET  /api/v1/relationships/:agentId
//   - Get relationships for specific agent
//   - Include: relationship type, strength, sentiment, interaction history
//
// POST /api/v1/relationships
//   - Force create/modify relationship between two agents
//   - Body: { "agent1": "id", "agent2": "id", "type": "friend|rival|neutral" }
//
// =============================================================================
// EVENT HANDLERS:
// =============================================================================
//
// GET  /api/v1/events
//   - List world events (past and scheduled)
//   - Support: ?status=pending|completed&type=global|agent
//
// POST /api/v1/events
//   - Inject a global event into the world
//   - Body: { "type": "disaster|celebration|discovery", "description": "...", "affectedAgents": ["all"|"id1,id2"] }
//
// GET  /api/v1/events/stream
//   - SSE endpoint for real-time event notifications
//
// =============================================================================
// WORLD HANDLERS:
// =============================================================================
//
// GET  /api/v1/world/status
//   - Get world state: current tick, simulation speed, active agents count
//
// POST /api/v1/world/control
//   - Control simulation: pause, resume, speed adjustment
//   - Body: { "action": "pause|resume|step|setSpeed", "value": 2.0 }
//
// GET  /api/v1/world/statistics
//   - Aggregate statistics: mood distribution, relationship density, activity metrics
//
// =============================================================================
// CONTROL HANDLERS (Admin):
// =============================================================================
//
// POST /api/v1/control/spawn
//   - Spawn new agent with specified personality
//   - Body: { "name": "AgentX", "personality": { "openness": 0.8, ... } }
//
// DELETE /api/v1/control/agents/:id
//   - Remove agent from simulation (soft delete, preserve history)
//
// POST /api/v1/control/reset
//   - Reset world state (confirmation required)
//   - Body: { "confirm": true, "preserveAgents": false }

package api
