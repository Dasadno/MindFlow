# AI Agent Society — Backend Server

Autonomous AI agents with personalities, emotions, memory, and social relationships.

**Stack:** Go 1.22 | SQLite | Chromem-go (Vector DB) | GigaChat via Ollama

**Repository:** [Milk-IslandAI](https://github.com/Milk-IslandAI)

## Project Structure

```
server/
├── cmd/server/main.go              # Entry point, bootstrap sequence
├── internal/
│   ├── api/
│   │   ├── router.go               # Gin router, middleware stack, route groups
│   │   ├── handlers.go             # REST endpoint implementations
│   │   ├── middleware.go            # CORS, logging, rate limiting, recovery
│   │   └── dto.go                  # Request/response JSON schemas
│   ├── agent/
│   │   ├── agent.go                # Core agent entity, lifecycle (tick/reflect)
│   │   ├── brain.go                # LLM-driven cognition, prompt engineering
│   │   ├── memory.go               # Three-tier memory (working/episodic/semantic)
│   │   └── emotions.go             # PAD + discrete emotion engine, appraisal
│   ├── world/
│   │   ├── orchestrator.go         # Simulation coordinator, action resolution
│   │   ├── time.go                 # Simulation clock, tick management
│   │   └── eventbus.go             # Pub-sub event system for agent communication
│   └── storage/
│       ├── sqlite.go               # SQLite schema, migrations, CRUD operations
│       └── vectordb.go             # Chromem-go vector store for episodic memory
├── pkg/gigachat/
│   └── client.go                   # GigaChat/Ollama LLM client with streaming
├── go.mod
└── README.md
```

## Running

```bash
# Start Ollama with a model
ollama run gigachat

# Run the server
go run ./cmd/server
```

Default port: `:8080`

---

## API Endpoints

Base URL: `/api/v1`

### Agents

| Method | Path | Description | Response Schema |
|--------|------|-------------|-----------------|
| `GET` | `/agents` | List all agents (paginated, filterable) | `{ "agents": [AgentSummary], "pagination": { "page": 1, "limit": 20, "total": 42 } }` |
| `GET` | `/agents/:id` | Get detailed agent profile | `AgentDetail` (see below) |
| `GET` | `/agents/:id/memory` | Get agent memories | `{ "memories": [MemoryEntry], "summary": "AI-generated summary" }` |
| `GET` | `/agents/:id/thoughts` | **SSE** — stream agent's live thought process | `data: { "content": "I wonder if...", "type": "reasoning", "timestamp": "..." }` |
| `POST` | `/agents/:id/inject` | Inject thought/memory/goal into agent | Request: `{ "type": "thought", "content": "..." }` |

### Relationships

| Method | Path | Description | Response Schema |
|--------|------|-------------|-----------------|
| `GET` | `/relationships` | Get full relationship graph | `RelationshipGraph` (see below) |
| `GET` | `/relationships/:agentId` | Get one agent's relationships | `{ "relationships": [RelationshipEdge] }` |
| `POST` | `/relationships` | Force-create/modify a relationship | Request: `{ "agent1": "id", "agent2": "id", "type": "friend" }` |

### Events

| Method | Path | Description | Response Schema |
|--------|------|-------------|-----------------|
| `GET` | `/events` | List world events (filterable) | `{ "events": [Event] }` |
| `POST` | `/events` | Inject global event | Request: `{ "type": "disaster", "description": "...", "affectedAgents": ["all"] }` |
| `GET` | `/events/stream` | **SSE** — real-time event stream | `data: { "id": "...", "type": "interaction", "description": "..." }` |

### World

| Method | Path | Description | Response Schema |
|--------|------|-------------|-----------------|
| `GET` | `/world/status` | Current simulation state | `WorldStatus` (see below) |
| `POST` | `/world/control` | Pause/resume/step/set speed | Request: `{ "action": "pause", "value": null }` |
| `GET` | `/world/statistics` | Aggregate statistics | `WorldStatistics` (see below) |

### Control Panel (Admin)

| Method | Path | Description | Response Schema |
|--------|------|-------------|-----------------|
| `POST` | `/control/spawn` | Spawn new agent | Request: `{ "name": "Eve", "personality": { "openness": 0.9, ... } }` |
| `DELETE` | `/control/agents/:id` | Remove agent (soft delete) | `{ "success": true, "message": "Agent deactivated" }` |
| `POST` | `/control/reset` | Reset world state | Request: `{ "confirm": true, "preserveAgents": false }` |

---

## Response Schemas

### AgentSummary

```json
{
  "id": "a1b2c3d4-...",
  "name": "Alice",
  "personalityType": "explorer",
  "currentMood": "excited",
  "moodIntensity": 0.75,
  "isActive": true
}
```

### AgentDetail

```json
{
  "id": "a1b2c3d4-...",
  "name": "Alice",
  "personality": {
    "openness": 0.85,
    "conscientiousness": 0.60,
    "extraversion": 0.70,
    "agreeableness": 0.55,
    "neuroticism": 0.30,
    "coreValues": ["curiosity", "honesty"],
    "quirks": ["talks to herself", "collects stones"]
  },
  "currentMood": {
    "label": "excited",
    "pad": { "pleasure": 0.6, "arousal": 0.8, "dominance": 0.4 },
    "activeEmotions": [
      { "type": "joy", "intensity": 0.7, "trigger": "discovered a new area" },
      { "type": "anticipation", "intensity": 0.5, "trigger": "upcoming meeting" }
    ]
  },
  "goals": [
    { "id": "g1", "description": "Make a new friend", "priority": 0.8, "progress": 0.3 },
    { "id": "g2", "description": "Explore the eastern zone", "priority": 0.5, "progress": 0.0 }
  ],
  "stats": {
    "totalInteractions": 47,
    "memoriesCount": 128,
    "relationshipsCount": 5,
    "daysSinceCreation": 3
  },
  "createdAt": "2026-02-15T10:00:00Z"
}
```

### MemoryEntry

```json
{
  "id": "m1e2f3...",
  "type": "episodic",
  "content": "Had a deep conversation with Bob about the meaning of trust",
  "emotionalTag": "trust",
  "importance": 0.85,
  "timestamp": "2026-02-15T14:30:00Z",
  "relatedAgents": ["bob-uuid"]
}
```

### RelationshipGraph

```json
{
  "nodes": [
    { "id": "alice-uuid", "label": "Alice", "type": "explorer", "size": 5 },
    { "id": "bob-uuid", "label": "Bob", "type": "guardian", "size": 3 }
  ],
  "edges": [
    {
      "source": "alice-uuid",
      "target": "bob-uuid",
      "type": "friend",
      "strength": 0.72,
      "label": "bonded over shared values"
    }
  ]
}
```

### Event

```json
{
  "id": "evt-001",
  "type": "global",
  "category": "discovery",
  "description": "A hidden cave was discovered in the northern mountains",
  "affectedAgents": ["alice-uuid", "bob-uuid"],
  "timestamp": "2026-02-15T16:00:00Z",
  "status": "active"
}
```

### WorldStatus

```json
{
  "currentTick": 1547,
  "simulationSpeed": 1.0,
  "isPaused": false,
  "activeAgents": 8,
  "totalEvents": 234,
  "uptime": "2h15m30s"
}
```

### WorldStatistics

```json
{
  "moodDistribution": {
    "happy": 3,
    "calm": 2,
    "anxious": 1,
    "excited": 2
  },
  "relationshipStats": {
    "totalConnections": 12,
    "averageStrength": 0.45,
    "strongestBond": { "agents": ["alice-uuid", "bob-uuid"], "strength": 0.92 },
    "rivalries": 2
  },
  "activityMetrics": {
    "interactionsLastHour": 23,
    "memoriesFormedLastHour": 67,
    "eventsLastHour": 5
  },
  "topInteractingAgents": [
    { "id": "alice-uuid", "name": "Alice", "personalityType": "explorer", "currentMood": "excited", "moodIntensity": 0.75, "isActive": true }
  ]
}
```

### Error Response

All errors follow this format:

```json
{
  "code": "AGENT_NOT_FOUND",
  "message": "Agent with ID 'xyz' does not exist",
  "details": null
}
```

| HTTP Status | Code | When |
|-------------|------|------|
| 400 | `BAD_REQUEST` | Invalid JSON, missing required fields |
| 404 | `NOT_FOUND` | Agent/event/relationship not found |
| 429 | `RATE_LIMITED` | Too many requests |
| 500 | `INTERNAL_ERROR` | Server error (details hidden in production) |
