// Package agent provides the Memory System for agent cognition.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// The Memory System manages both episodic (event-based) and semantic
// (knowledge-based) memory. It uses vector embeddings for similarity search,
// enabling agents to recall relevant memories based on current context.
//
// =============================================================================
// MEMORY ARCHITECTURE:
// =============================================================================
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                        MEMORY SYSTEM                                 │
// ├─────────────────────────────────────────────────────────────────────┤
// │  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐  │
// │  │  Working Memory │    │ Episodic Memory │    │ Semantic Memory │  │
// │  │  (Short-term)   │───►│  (Vector DB)    │───►│  (SQLite)       │  │
// │  │  Recent events  │    │  Experiences    │    │  Knowledge      │  │
// │  └─────────────────┘    └─────────────────┘    └─────────────────┘  │
// │           │                      │                      │           │
// │           └──────────────────────┼──────────────────────┘           │
// │                                  ▼                                  │
// │                    ┌─────────────────────────┐                      │
// │                    │   Memory Consolidator   │                      │
// │                    │   (Summarization via    │                      │
// │                    │    GigaChat LLM)        │                      │
// │                    └─────────────────────────┘                      │
// └─────────────────────────────────────────────────────────────────────┘
//
// =============================================================================
// MEMORY STRUCTURES:
// =============================================================================
//
// type MemorySystem struct {
//     workingMemory  *WorkingMemory    // Short-term buffer
//     episodicStore  *VectorStore      // Chromem-go for embeddings
//     semanticStore  *SQLiteStore      // Persistent knowledge
//     consolidator   *MemoryConsolidator
//     config         MemoryConfig
// }
//
// type MemoryEntry struct {
//     ID            string
//     Type          MemoryType        // Episodic, Semantic, Procedural
//     Content       string            // The actual memory content
//     Embedding     []float32         // Vector embedding for similarity
//     EmotionalTag  EmotionTag        // Emotion during encoding
//     Importance    float64           // Calculated salience (0-1)
//     Timestamp     time.Time
//     AccessCount   int               // How often recalled
//     LastAccessed  time.Time
//     RelatedAgents []string          // Agents involved in this memory
//     Metadata      map[string]any    // Flexible additional data
// }
//
// type MemoryType string
// const (
//     MemoryEpisodic   MemoryType = "episodic"   // Specific events
//     MemorySemantic   MemoryType = "semantic"   // General knowledge
//     MemoryProcedural MemoryType = "procedural" // How to do things
// )
//
// =============================================================================
// CORE MEMORY OPERATIONS:
// =============================================================================
//
// func (m *MemorySystem) Encode(experience Experience) *MemoryEntry
//   - Convert raw experience into memory entry
//   - Generate embedding via LLM
//   - Calculate importance based on emotional intensity and novelty
//   - Store in appropriate memory store
//
// func (m *MemorySystem) Recall(query string, limit int) []MemoryEntry
//   - Vector similarity search for relevant memories
//   - Rank by: similarity score, recency, importance, access frequency
//   - Update access metadata on recalled memories
//
// func (m *MemorySystem) RecallRelated(agentID string, limit int) []MemoryEntry
//   - Retrieve memories involving specific agent
//   - Used for relationship context building
//
// func (m *MemorySystem) Forget(threshold float64)
//   - Remove low-importance, rarely-accessed memories
//   - Simulate natural forgetting curve
//   - Preserve emotionally significant memories longer
//
// =============================================================================
// MEMORY CONSOLIDATION:
// =============================================================================
//
// func (c *MemoryConsolidator) Consolidate(memories []MemoryEntry) SemanticMemory
//   - Periodically run (simulating sleep/reflection)
//   - Group related episodic memories
//   - Use LLM to extract patterns and general knowledge
//   - Create semantic memories from episodic clusters
//   - Example: Multiple "met Agent X at cafe" → "Agent X frequents the cafe"
//
// func (c *MemoryConsolidator) Summarize(memories []MemoryEntry) string
//   - Generate natural language summary of memories
//   - Used for context building and API responses
//
// =============================================================================
// WORKING MEMORY:
// =============================================================================
//
// type WorkingMemory struct {
//     buffer   []MemoryEntry  // Fixed-size circular buffer
//     capacity int
//     current  int
// }
//
// func (w *WorkingMemory) Add(entry MemoryEntry)
//   - Add to short-term buffer
//   - Overflow triggers consolidation check
//
// func (w *WorkingMemory) GetRecent(n int) []MemoryEntry
//   - Get n most recent memories
//   - Used for immediate context in thinking

package agent
