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

## Database Schema

**Engine:** SQLite (WAL mode) | **File:** `server/data/society.db`

### ER Diagram

```
┌─────────────────────────────────┐       ┌──────────────────────────────────────┐
│            agents               │       │           relationships              │
├─────────────────────────────────┤       ├──────────────────────────────────────┤
│ PK  id              TEXT        │◄──┐   │ PK  id                TEXT          │
│     name            TEXT    [NN]│   │   │ FK  agent1_id          TEXT     [NN] │──┐
│     personality     TEXT    [NN]│   │   │ FK  agent2_id          TEXT     [NN] │──┤
│     mood_state      TEXT        │   │   │     type               TEXT     [NN] │  │
│     goals           TEXT        │   │   │     strength           REAL   [0.0]  │  │
│     state           TEXT [idle] │   │   │     interaction_count  INTEGER  [0]  │  │
│     is_active       BOOL  [1]  │   │   │     last_interaction   DATETIME      │  │
│     created_at      DATETIME[NN]│   │   │     metadata           TEXT          │  │
│     last_active     DATETIME    │   │   │                                      │  │
│     snapshot        TEXT        │   │   │ UQ (agent1_id, agent2_id)            │  │
└─────────────────────────────────┘   │   └──────────────────────────────────────┘  │
              ▲                       │                                             │
              │                       └─────────────────────────────────────────────┘
              │
              │   ┌──────────────────────────────────────┐
              │   │            memories                   │
              │   ├──────────────────────────────────────┤
              │   │ PK  id              TEXT              │
              ├───│ FK  agent_id        TEXT         [NN] │
              │   │     type            TEXT         [NN] │
              │   │     content         TEXT         [NN] │
              │   │     emotional_tag   TEXT              │
              │   │     importance      REAL       [0.5]  │
              │   │     access_count    INTEGER      [0]  │
              │   │     last_accessed   DATETIME          │
              │   │     related_agents  TEXT              │
              │   │     metadata        TEXT              │
              │   │     created_at      DATETIME    [NN]  │
              │   └──────────────────────────────────────┘
              │
              │   ┌──────────────────────────────────────┐
              │   │             events                    │
              │   ├──────────────────────────────────────┤
              │   │ PK  id              TEXT              │
              │   │     topic           TEXT         [NN] │
              │   │     type            TEXT         [NN] │
              │   │     source          TEXT         [NN] │
              │   │     affected_agents TEXT              │
              │   │     payload         TEXT              │
              │   │     status          TEXT   [pending]  │
              │   │     tick            INTEGER           │
              │   │     created_at      DATETIME    [NN]  │
              │   └──────────────────────────────────────┘

┌──────────────────────────────────────┐
│          world_state (KV)            │
├──────────────────────────────────────┤
│ PK  key             TEXT             │
│     value           TEXT        [NN] │
│     updated_at      DATETIME   [NN]  │
└──────────────────────────────────────┘
```

**Legend:** `PK` — Primary Key | `FK` — Foreign Key | `UQ` — Unique Constraint | `[NN]` — NOT NULL | `[value]` — DEFAULT

### Column Details

#### `agents` — autonomous AI entities

| Column | Type | Default | Description |
|--------|------|---------|-------------|
| `id` | TEXT | — | UUID, primary key |
| `name` | TEXT | — | Display name |
| `personality` | TEXT (JSON) | — | `{ "openness": 0.85, "conscientiousness": 0.6, "extraversion": 0.7, "agreeableness": 0.55, "neuroticism": 0.3, "coreValues": [...], "quirks": [...] }` |
| `mood_state` | TEXT (JSON) | NULL | `{ "pleasure": 0.6, "arousal": 0.8, "dominance": 0.4 }` — PAD model |
| `goals` | TEXT (JSON) | NULL | `[{ "description": "...", "priority": 0.8, "progress": 0.3 }]` |
| `state` | TEXT | `idle` | One of: `idle`, `thinking`, `acting`, `interacting`, `sleeping` |
| `is_active` | BOOLEAN | `1` | `0` = soft-deleted |
| `snapshot` | TEXT (JSON) | NULL | Full serialized state for server restarts |

#### `relationships` — connections between agents

| Column | Type | Default | Description |
|--------|------|---------|-------------|
| `id` | TEXT | — | UUID |
| `agent1_id` | TEXT (FK) | — | First agent (ON DELETE CASCADE) |
| `agent2_id` | TEXT (FK) | — | Second agent (ON DELETE CASCADE) |
| `type` | TEXT | — | `friend` \| `rival` \| `neutral` \| `romantic` |
| `strength` | REAL | `0.0` | Range: `-1.0` (hostile) to `1.0` (close bond) |
| `interaction_count` | INTEGER | `0` | Total interactions between the pair |
| `metadata` | TEXT (JSON) | NULL | `{ "firstMet": "...", "sharedMemories": 5 }` |

#### `memories` — episodic and semantic memory entries

| Column | Type | Default | Description |
|--------|------|---------|-------------|
| `id` | TEXT | — | UUID |
| `agent_id` | TEXT (FK) | — | Owner agent (ON DELETE CASCADE) |
| `type` | TEXT | — | `episodic` \| `semantic` \| `procedural` |
| `content` | TEXT | — | Natural language memory content |
| `emotional_tag` | TEXT | NULL | Emotion at encoding time: `joy`, `fear`, `trust`, etc. |
| `importance` | REAL | `0.5` | Salience score `0.0`–`1.0` |
| `access_count` | INTEGER | `0` | Recall frequency (affects retention) |
| `related_agents` | TEXT (JSON) | NULL | `["agent-uuid-1", "agent-uuid-2"]` |

#### `events` — world and agent events

| Column | Type | Default | Description |
|--------|------|---------|-------------|
| `id` | TEXT | — | UUID |
| `topic` | TEXT | — | `global` \| `interaction` \| `mood_change` \| `goal_update` \| `memory` \| `relationship` \| `system` |
| `type` | TEXT | — | Specific event type: `disaster`, `celebration`, `discovery`, etc. |
| `source` | TEXT | — | Agent ID, `"system"`, or `"api"` |
| `affected_agents` | TEXT (JSON) | NULL | `["all"]` or `["uuid-1", "uuid-2"]` |
| `payload` | TEXT (JSON) | NULL | Event-specific data |
| `status` | TEXT | `pending` | `pending` \| `active` \| `completed` |
| `tick` | INTEGER | NULL | Simulation tick when event occurred |

#### `world_state` — key-value simulation state

| Key | Initial Value | Description |
|-----|---------------|-------------|
| `current_tick` | `0` | Monotonically increasing tick counter |
| `simulation_speed` | `1.0` | Speed multiplier (0.1x – 10.0x) |
| `is_paused` | `true` | Simulation starts paused |

### Indexes

| Index | Table | Column(s) | Purpose |
|-------|-------|-----------|---------|
| `idx_agents_is_active` | agents | `is_active` | Filter active agents |
| `idx_agents_state` | agents | `state` | Filter by state |
| `idx_relationships_agent1` | relationships | `agent1_id` | Lookup by first agent |
| `idx_relationships_agent2` | relationships | `agent2_id` | Lookup by second agent |
| `idx_relationships_type` | relationships | `type` | Filter by type |
| `idx_memories_agent_id` | memories | `agent_id` | Agent's memories |
| `idx_memories_type` | memories | `agent_id, type` | Agent's memories by type |
| `idx_memories_importance` | memories | `importance` | Sort by salience |
| `idx_memories_created_at` | memories | `created_at` | Chronological queries |
| `idx_events_topic` | events | `topic` | Filter by topic |
| `idx_events_source` | events | `source` | Filter by source |
| `idx_events_status` | events | `status` | Filter by status |
| `idx_events_tick` | events | `tick` | Tick-based queries |

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


# Client side of the project 