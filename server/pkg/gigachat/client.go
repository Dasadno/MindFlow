// Package gigachat provides a client for the GigaChat LLM API via Ollama.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This package encapsulates all communication with the GigaChat language model.
// It can operate in two modes:
//   1. Via Ollama — local inference server (recommended for hackathon)
//   2. Direct GigaChat API — cloud-based (requires API key)
//
// The client provides high-level methods tailored to agent cognition needs:
// thinking, reflection, conversation, and embedding generation.
//
// =============================================================================
// CLIENT STRUCTURE:
// =============================================================================
//
// type Client struct {
//     baseURL    string            // Ollama: "http://localhost:11434", GigaChat: API URL
//     model      string            // Model name: "gigachat" or specific Ollama model
//     apiKey     string            // Only for direct GigaChat API
//     httpClient *http.Client
//     mode       ClientMode        // Ollama or Direct
//     config     ClientConfig
// }
//
// type ClientMode string
// const (
//     ModeOllama   ClientMode = "ollama"   // Local Ollama server
//     ModeDirect   ClientMode = "direct"   // Direct GigaChat API
// )
//
// type ClientConfig struct {
//     DefaultTemperature float64       // Base LLM temperature (0.0 - 2.0)
//     MaxTokens          int           // Max response tokens
//     Timeout            time.Duration // Request timeout
//     RetryAttempts      int           // Retry on transient failures
//     RetryDelay         time.Duration // Delay between retries
// }
//
// =============================================================================
// CLIENT INITIALIZATION:
// =============================================================================
//
// func NewClient(baseURL, model string, opts ...ClientOption) *Client
//   - Create client with Ollama endpoint and model name
//   - Apply functional options for configuration
//   - Default: temperature 0.7, max tokens 1024, timeout 30s
//
// func NewDirectClient(apiURL, apiKey, model string, opts ...ClientOption) *Client
//   - Create client for direct GigaChat API access
//   - Handle GigaChat-specific authentication (OAuth token)
//
// type ClientOption func(*Client)
// func WithTemperature(t float64) ClientOption
// func WithMaxTokens(n int) ClientOption
// func WithTimeout(d time.Duration) ClientOption
//
// =============================================================================
// CORE LLM OPERATIONS:
// =============================================================================
//
// func (c *Client) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
//   - Send chat completion request to LLM
//   - Handle Ollama API format: POST /api/chat
//   - Handle retry logic and timeout
//   - Return parsed response with token usage stats
//
// type CompletionRequest struct {
//     SystemPrompt string            // Agent personality and role
//     Messages     []Message         // Conversation history
//     Temperature  *float64          // Override default temperature
//     MaxTokens    *int              // Override default max tokens
//     Stream       bool              // Enable streaming response
// }
//
// type Message struct {
//     Role    string `json:"role"`    // "system", "user", "assistant"
//     Content string `json:"content"`
// }
//
// type CompletionResponse struct {
//     Content    string  // Generated text
//     TokensUsed int     // Total tokens consumed
//     Model      string  // Model that responded
//     Duration   time.Duration
// }
//
// =============================================================================
// AGENT-SPECIFIC METHODS:
// =============================================================================
//
// func (c *Client) Think(ctx context.Context, personality string, context string, memories []string) (string, error)
//   - High-level thinking method for agent cognition
//   - Builds system prompt from personality
//   - Includes memories and current context
//   - Returns structured thought output
//
// func (c *Client) Converse(ctx context.Context, agent1Prompt string, agent2Prompt string, history []Message) (string, error)
//   - Generate next dialogue turn in agent conversation
//   - Consider both agents' contexts
//
// func (c *Client) Reflect(ctx context.Context, personality string, memories []string) (string, error)
//   - Deep reflection prompt for meta-cognition
//   - Ask LLM to find patterns, insights, goal updates
//   - Return structured reflection output
//
// func (c *Client) Summarize(ctx context.Context, texts []string) (string, error)
//   - Summarize multiple memory entries into semantic knowledge
//   - Used by MemoryConsolidator
//
// =============================================================================
// EMBEDDING OPERATIONS:
// =============================================================================
//
// func (c *Client) Embed(ctx context.Context, text string) ([]float32, error)
//   - Generate vector embedding for text
//   - Ollama endpoint: POST /api/embeddings
//   - Used by VectorStore for memory encoding and similarity search
//
// func (c *Client) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
//   - Batch embedding generation
//   - More efficient for bulk operations (initial memory loading)
//
// =============================================================================
// STREAMING SUPPORT:
// =============================================================================
//
// func (c *Client) CompleteStream(ctx context.Context, req CompletionRequest) (<-chan StreamChunk, error)
//   - Return channel of streaming response chunks
//   - Used by Brain.StreamThoughts() for real-time dashboard
//
// type StreamChunk struct {
//     Content string // Partial response text
//     Done    bool   // Is this the final chunk?
//     Error   error  // Non-nil if stream errored
// }

package gigachat
