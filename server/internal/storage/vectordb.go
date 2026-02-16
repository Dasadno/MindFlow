// Package storage provides the Vector Database layer for episodic memory.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This file wraps the Chromem-go (or custom) vector database used for
// similarity-based memory retrieval. Episodic memories are stored as vector
// embeddings, enabling agents to recall experiences based on semantic
// similarity rather than exact keyword matches.
//
// =============================================================================
// WHY A VECTOR DB:
// =============================================================================
// When an agent encounters a situation, it needs to recall RELEVANT past
// experiences — not just recent ones. Vector similarity search enables:
//   - "I'm at a party" → recalls all party-related memories
//   - "Agent X is angry" → recalls past conflicts with Agent X
//   - "I feel lonely" → recalls times of connection and isolation
//
// =============================================================================
// VECTOR STORE STRUCTURE:
// =============================================================================
//
// type VectorStore struct {
//     db         *chromem.DB              // Chromem-go instance
//     collection *chromem.Collection      // One collection per agent
//     embedder   EmbeddingProvider        // Generates vector embeddings
// }
//
// type EmbeddingProvider interface {
//     Embed(text string) ([]float32, error)
//     EmbedBatch(texts []string) ([][]float32, error)
// }
//
// =============================================================================
// EMBEDDING STRATEGY:
// =============================================================================
//
// Two options for generating embeddings:
//
// Option A — GigaChat/Ollama Embeddings:
//   - Use the LLM's embedding endpoint
//   - Higher quality, but requires API calls
//   - Recommended for production
//
// Option B — Simple Local Embeddings:
//   - Bag-of-words or TF-IDF based vectors
//   - No external dependencies, faster
//   - Good enough for hackathon / demo
//
// type LocalEmbedder struct {
//     vocabulary map[string]int    // word → dimension index
//     idfScores  map[string]float64
//     dimensions int
// }
//
// =============================================================================
// VECTOR STORE OPERATIONS:
// =============================================================================
//
// func NewVectorStore(dbPath string, embedder EmbeddingProvider) (*VectorStore, error)
//   - Initialize Chromem-go with persistent storage at dbPath
//   - Create default collection if not exists
//
// func (vs *VectorStore) CreateCollection(agentID string) error
//   - Create a separate collection per agent
//   - Enables isolated memory spaces
//
// func (vs *VectorStore) AddMemory(agentID string, memory VectorMemory) error
//   - Generate embedding for memory content
//   - Store with metadata (timestamp, emotional tag, importance)
//   - Assign unique document ID
//
// func (vs *VectorStore) Search(agentID string, query string, limit int) ([]VectorSearchResult, error)
//   - Generate embedding for query
//   - Find top-K most similar memories
//   - Return with similarity scores
//   - Filter by agent's collection
//
// func (vs *VectorStore) SearchWithFilter(agentID string, query string, filter MemoryFilter, limit int) ([]VectorSearchResult, error)
//   - Similarity search with additional metadata filters
//   - Filter by: time range, emotional tag, importance threshold, related agents
//
// func (vs *VectorStore) DeleteMemory(agentID string, memoryID string) error
//   - Remove specific memory from vector store
//   - Used during memory consolidation (replace episodic with semantic)
//
// func (vs *VectorStore) DeleteCollection(agentID string) error
//   - Remove all memories for an agent
//   - Used when removing agent from simulation
//
// =============================================================================
// DATA STRUCTURES:
// =============================================================================
//
// type VectorMemory struct {
//     ID            string
//     Content       string            // Text to embed
//     EmotionalTag  string
//     Importance    float64
//     Timestamp     time.Time
//     RelatedAgents []string
//     Metadata      map[string]string // Chromem metadata (string values only)
// }
//
// type VectorSearchResult struct {
//     Memory     VectorMemory
//     Similarity float32           // Cosine similarity score (0-1)
// }
//
// type MemoryFilter struct {
//     After       *time.Time        // Only memories after this time
//     Before      *time.Time        // Only memories before this time
//     MinImportance float64         // Minimum importance threshold
//     EmotionalTag  string          // Filter by emotion
//     RelatedAgent  string          // Must involve this agent
// }

package storage
