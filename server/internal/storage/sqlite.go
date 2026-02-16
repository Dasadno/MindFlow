// Package storage provides the persistence layer for the AI Agent Society.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This file manages the SQLite database connection, schema migrations, and
// all CRUD operations for persistent data (agents, relationships, events,
// semantic memories, world state).
//
// =============================================================================
// DATABASE SCHEMA:
// =============================================================================
//
// TABLE agents:
//   id              TEXT PRIMARY KEY     -- UUID
//   name            TEXT NOT NULL
//   personality     TEXT NOT NULL         -- JSON blob of Big Five traits + values
//   mood_state      TEXT                  -- JSON blob of current PAD state
//   goals           TEXT                  -- JSON array of active goals
//   state           TEXT DEFAULT 'idle'   -- Current AgentState
//   is_active       BOOLEAN DEFAULT TRUE
//   created_at      DATETIME NOT NULL
//   last_active     DATETIME
//   snapshot        TEXT                  -- Full agent state for restore
//
// TABLE relationships:
//   id              TEXT PRIMARY KEY
//   agent1_id       TEXT NOT NULL REFERENCES agents(id)
//   agent2_id       TEXT NOT NULL REFERENCES agents(id)
//   type            TEXT NOT NULL         -- friend, rival, neutral, romantic
//   strength        REAL DEFAULT 0.0     -- -1.0 to 1.0
//   interaction_count INTEGER DEFAULT 0
//   last_interaction DATETIME
//   metadata        TEXT                  -- JSON blob for extra context
//   UNIQUE(agent1_id, agent2_id)
//
// TABLE memories:
//   id              TEXT PRIMARY KEY
//   agent_id        TEXT NOT NULL REFERENCES agents(id)
//   type            TEXT NOT NULL         -- episodic, semantic, procedural
//   content         TEXT NOT NULL
//   emotional_tag   TEXT
//   importance      REAL DEFAULT 0.5
//   access_count    INTEGER DEFAULT 0
//   last_accessed   DATETIME
//   related_agents  TEXT                  -- JSON array of agent IDs
//   metadata        TEXT                  -- JSON blob
//   created_at      DATETIME NOT NULL
//
// TABLE events:
//   id              TEXT PRIMARY KEY
//   topic           TEXT NOT NULL
//   type            TEXT NOT NULL
//   source          TEXT NOT NULL
//   affected_agents TEXT                  -- JSON array
//   payload         TEXT                  -- JSON blob
//   status          TEXT DEFAULT 'pending'
//   tick            INTEGER
//   created_at      DATETIME NOT NULL
//
// TABLE world_state:
//   key             TEXT PRIMARY KEY
//   value           TEXT NOT NULL
//   updated_at      DATETIME NOT NULL
//
// =============================================================================
// REPOSITORY STRUCTURE:
// =============================================================================
//
// type Repository struct {
//     db *sql.DB
// }
//
// func NewRepository(dbPath string) (*Repository, error)
//   - Open SQLite database at given path
//   - Enable WAL mode for concurrent reads
//   - Run Migrate() to ensure schema is up to date
//   - Return initialized repository
//
// func (r *Repository) Migrate() error
//   - Create all tables if they don't exist
//   - Create indexes for common queries
//   - Idempotent: safe to run on every startup
//
// func (r *Repository) Close() error
//   - Close database connection cleanly
//
// =============================================================================
// AGENT OPERATIONS:
// =============================================================================
//
// func (r *Repository) SaveAgent(agent AgentRecord) error
// func (r *Repository) GetAgent(id string) (*AgentRecord, error)
// func (r *Repository) ListAgents(filter AgentFilter) ([]AgentRecord, error)
// func (r *Repository) UpdateAgent(id string, update AgentUpdate) error
// func (r *Repository) DeleteAgent(id string) error  // Soft delete
//
// =============================================================================
// RELATIONSHIP OPERATIONS:
// =============================================================================
//
// func (r *Repository) SaveRelationship(rel RelationshipRecord) error
// func (r *Repository) GetRelationships(agentID string) ([]RelationshipRecord, error)
// func (r *Repository) GetRelationshipGraph() (*GraphData, error)
//   - Returns all nodes and edges for graph visualization
// func (r *Repository) UpdateRelationship(id string, update RelUpdate) error
//
// =============================================================================
// MEMORY OPERATIONS:
// =============================================================================
//
// func (r *Repository) SaveMemory(memory MemoryRecord) error
// func (r *Repository) GetMemories(agentID string, memType string, limit int) ([]MemoryRecord, error)
// func (r *Repository) UpdateMemoryAccess(id string) error  // Increment access count
// func (r *Repository) DeleteOldMemories(agentID string, threshold float64) (int, error)
//
// =============================================================================
// EVENT OPERATIONS:
// =============================================================================
//
// func (r *Repository) SaveEvent(event EventRecord) error
// func (r *Repository) GetEvents(filter EventFilter) ([]EventRecord, error)
// func (r *Repository) UpdateEventStatus(id string, status string) error
//
// =============================================================================
// WORLD STATE OPERATIONS:
// =============================================================================
//
// func (r *Repository) SaveWorldState(key string, value string) error
// func (r *Repository) GetWorldState(key string) (string, error)
//   - Key-value store for: current_tick, simulation_speed, is_paused, etc.

package storage
