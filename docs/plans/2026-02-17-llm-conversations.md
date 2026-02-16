# LLM Agent Conversations Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Подключить Gemma через Ollama — агенты автономно разговаривают друг с другом, диалоги стримятся через SSE, человек может вписаться в разговор.

**Architecture:** Ollama HTTP-клиент вызывает Gemma для каждой реплики. Orchestrator-тикер каждые 15 секунд берёт двух случайных агентов и запускает диалог из 4 реплик. Каждая реплика сохраняется в `events` и пушится в SSE Hub. Инъекция человека хранится в памяти Hub и добавляется в контекст следующего тика.

**Tech Stack:** Go 1.24, `net/http` (stdlib), `encoding/json`, `modernc.org/sqlite`, Ollama REST API (`POST /api/chat`)

---

## Контекст кодовой базы

```
server/
├── cmd/server/main.go                   ← добавить запуск Orchestrator
├── cmd/server/middleware/mux.go          ← подключить SSE и inject handlers
├── data/db.go                           ← global data.Db
├── internal/
│   ├── api/
│   │   ├── handler.go                   ← Handler{repo, hub} — добавить hub
│   │   ├── agent_handlers.go            ← не трогать
│   │   ├── control_handlers.go          ← не трогать
│   │   ├── world_handlers.go            ← не трогать
│   │   ├── helpers.go                   ← writeJSON, writeError
│   │   ├── dto.go                       ← все DTO
│   │   └── middleware.go                ← APIError
│   ├── agent/
│   │   ├── agent.go                     ← Agent, Personality, Goal, Mood типы
│   │   ├── brain.go                     ← добавить NewBrain, Think, buildSystemPrompt
│   │   └── emotions.go                  ← Mood, PADState типы
│   ├── world/
│   │   └── orchestrator.go              ← добавить NewOrchestrator, Start, runConversation
│   └── storage/
│       └── sqlite.go                    ← добавить GetRandomActiveAgents
└── pkg/gigachat/
    └── client.go                        ← реализовать Complete()
```

**Важно:**
- Модуль: `milk`
- Ollama URL: env `OLLAMA_URL` (default: `http://localhost:11434`)
- Модель: env `OLLAMA_MODEL` (default: `gemma3`)
- Тик: 15 секунд
- Реплик в диалоге: 4 (по 2 каждому агенту)

---

### Task 1: Реализовать Ollama HTTP-клиент

**Files:**
- Modify: `server/pkg/gigachat/client.go`

**Step 1: Добавить реализацию в конец `server/pkg/gigachat/client.go`**

Файл уже содержит типы `Client`, `CompletionRequest`, `CompletionResponse`, `Message`, `ClientMode`, `ClientConfig`. Нужно добавить конструктор и метод `Complete`.

```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

// NewClient создаёт Ollama-клиент из env-переменных.
// OLLAMA_URL — адрес сервера (default: http://localhost:11434)
// OLLAMA_MODEL — модель (default: gemma3)
func NewClient() *Client {
    baseURL := os.Getenv("OLLAMA_URL")
    if baseURL == "" {
        baseURL = "http://localhost:11434"
    }
    model := os.Getenv("OLLAMA_MODEL")
    if model == "" {
        model = "gemma3"
    }
    return &Client{
        BaseURL: baseURL,
        Model:   model,
        Mode:    ModeOllama,
        HTTPClient: &http.Client{
            Timeout: 60 * time.Second,
        },
        Config: ClientConfig{
            DefaultTemperature: 0.7,
            MaxTokens:          512,
            Timeout:            60 * time.Second,
        },
    }
}

// ollamaChatRequest — тело запроса к Ollama /api/chat.
type ollamaChatRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
    Stream   bool      `json:"stream"`
    Options  map[string]any `json:"options,omitempty"`
}

// ollamaChatResponse — ответ от Ollama /api/chat (stream: false).
type ollamaChatResponse struct {
    Message struct {
        Content string `json:"content"`
    } `json:"message"`
    Model string `json:"model"`
    Done  bool   `json:"done"`
}

// Complete отправляет запрос в Ollama и возвращает ответ.
func (c *Client) Complete(req CompletionRequest) (CompletionResponse, error) {
    start := time.Now()

    messages := req.Messages
    if req.SystemPrompt != "" {
        messages = append([]Message{{Role: "system", Content: req.SystemPrompt}}, messages...)
    }

    temp := c.Config.DefaultTemperature
    if req.Temperature != nil {
        temp = *req.Temperature
    }

    body := ollamaChatRequest{
        Model:    c.Model,
        Messages: messages,
        Stream:   false,
        Options:  map[string]any{"temperature": temp},
    }

    data, err := json.Marshal(body)
    if err != nil {
        return CompletionResponse{}, fmt.Errorf("Complete marshal: %w", err)
    }

    resp, err := c.HTTPClient.Post(
        c.BaseURL+"/api/chat",
        "application/json",
        bytes.NewReader(data),
    )
    if err != nil {
        return CompletionResponse{}, fmt.Errorf("Complete http: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return CompletionResponse{}, fmt.Errorf("Complete status %d", resp.StatusCode)
    }

    var ollResp ollamaChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&ollResp); err != nil {
        return CompletionResponse{}, fmt.Errorf("Complete decode: %w", err)
    }

    return CompletionResponse{
        Content:  ollResp.Message.Content,
        Model:    ollResp.Model,
        Duration: time.Since(start),
    }, nil
}
```

**Step 2: Добавить `"bytes"` и `"os"` в imports**

Текущие imports в файле: `"net/http"`, `"time"`. Нужно добавить `"bytes"`, `"encoding/json"`, `"fmt"`, `"os"`.

**Step 3: Компиляция**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI && go build ./server/...
```
Expected: без ошибок.

**Step 4: Commit**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI
git add server/pkg/gigachat/client.go
git commit -m "feat(gigachat): implement Ollama HTTP client"
```

---

### Task 2: Brain — системный промпт + Think()

**Files:**
- Modify: `server/internal/agent/brain.go`

**Step 1: Добавить в конец `server/internal/agent/brain.go`**

Файл содержит типы `Brain`, `BrainConfig`, `Thought`, `CognitiveContext`, `CognitiveOutput`. Нужно добавить конструктор и методы.

```go
import (
    "context"
    "fmt"
    "strings"
    "time"

    "milk/server/pkg/gigachat"
)

// llmClient — интерфейс LLM-клиента для Brain.
// Позволяет подменить реальный клиент на мок в тестах.
type llmClient interface {
    Complete(req gigachat.CompletionRequest) (gigachat.CompletionResponse, error)
}

// NewBrain создаёт Brain для агента.
func NewBrain(personality *Personality, client llmClient) *Brain {
    creativity := 0.5 + personality.Openness*0.5 // 0.5–1.0
    return &Brain{
        Personality:   personality,
        ThoughtBuffer: make([]Thought, 0, 10),
        ThoughtStream: make(chan Thought, 32),
        Config: BrainConfig{
            MaxThoughts:      10,
            CreativityFactor: creativity,
            ResponseTimeout:  60 * time.Second,
        },
    }
}

// buildSystemPrompt строит системный промпт из личности агента.
func buildSystemPrompt(name string, p *Personality, mood Mood, goals []Goal) string {
    var sb strings.Builder

    sb.WriteString(fmt.Sprintf("You are %s, an autonomous AI agent in a social simulation.\n\n", name))

    // Личность в словах
    sb.WriteString("Your personality:\n")
    if p.Openness > 0.6 {
        sb.WriteString("- You are curious and open to new ideas.\n")
    } else {
        sb.WriteString("- You are practical and prefer familiar things.\n")
    }
    if p.Extraversion > 0.6 {
        sb.WriteString("- You are outgoing and enjoy socializing.\n")
    } else {
        sb.WriteString("- You are introverted and prefer quiet reflection.\n")
    }
    if p.Agreeableness > 0.6 {
        sb.WriteString("- You are warm, cooperative, and empathetic.\n")
    } else {
        sb.WriteString("- You are direct, competitive, and skeptical.\n")
    }
    if p.Conscientiousness > 0.6 {
        sb.WriteString("- You are organized and disciplined.\n")
    } else {
        sb.WriteString("- You are spontaneous and flexible.\n")
    }
    if p.Neuroticism > 0.6 {
        sb.WriteString("- You tend to be anxious and emotionally reactive.\n")
    } else {
        sb.WriteString("- You are emotionally stable and calm under pressure.\n")
    }

    if len(p.CoreValues) > 0 {
        sb.WriteString(fmt.Sprintf("\nYour core values: %s.\n", strings.Join(p.CoreValues, ", ")))
    }
    if len(p.Quirks) > 0 {
        sb.WriteString(fmt.Sprintf("Your quirks: %s.\n", strings.Join(p.Quirks, ", ")))
    }

    sb.WriteString(fmt.Sprintf("\nYour current mood: %s.\n", string(mood)))

    if len(goals) > 0 {
        sb.WriteString("\nYour current goals:\n")
        for _, g := range goals {
            if !g.IsCompleted {
                sb.WriteString(fmt.Sprintf("- %s\n", g.Description))
            }
        }
    }

    sb.WriteString("\nIMPORTANT: Keep responses concise (2-3 sentences max). Stay in character. Be natural and conversational.")

    return sb.String()
}

// Think вызывает LLM с историей диалога и возвращает следующую реплику.
// history — чередование реплик [{"role":"user",...},{"role":"assistant",...}].
func (b *Brain) Think(
    ctx context.Context,
    client llmClient,
    name string,
    mood Mood,
    goals []Goal,
    history []gigachat.Message,
) (string, error) {
    sysPrompt := buildSystemPrompt(name, b.Personality, mood, goals)

    req := gigachat.CompletionRequest{
        SystemPrompt: sysPrompt,
        Messages:     history,
    }

    if b.Config.CreativityFactor > 0 {
        t := b.Config.CreativityFactor * 0.9
        req.Temperature = &t
    }

    resp, err := client.Complete(req)
    if err != nil {
        return "", fmt.Errorf("Brain.Think: %w", err)
    }

    thought := Thought{
        Content:   resp.Content,
        Type:      ThoughtDecision,
        Timestamp: time.Now(),
    }

    // Добавляем в буфер (кольцо)
    b.ThoughtBuffer = append(b.ThoughtBuffer, thought)
    if len(b.ThoughtBuffer) > b.Config.MaxThoughts {
        b.ThoughtBuffer = b.ThoughtBuffer[1:]
    }

    // Пушим в стрим (non-blocking)
    select {
    case b.ThoughtStream <- thought:
    default:
    }

    return resp.Content, nil
}
```

**Step 2: Компиляция**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI && go build ./server/...
```

**Step 3: Commit**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI
git add server/internal/agent/brain.go
git commit -m "feat(brain): add NewBrain, buildSystemPrompt, Think"
```

---

### Task 3: SSE Hub

**Files:**
- Create: `server/internal/api/hub.go`

**Step 1: Создать `server/internal/api/hub.go`**

```go
package api

import (
    "encoding/json"
    "sync"
)

// SSEEvent — событие, отправляемое клиентам через SSE.
type SSEEvent struct {
    Type    string `json:"type"`    // "conversation", "thought", "system"
    Speaker string `json:"speaker,omitempty"`
    Target  string `json:"target,omitempty"`
    Content string `json:"content"`
    AgentID string `json:"agentId,omitempty"`
    Tick    int64  `json:"tick"`
}

// Hub — in-memory SSE broadcast hub.
// Хранит подписчиков (каналы) и очередь инъекций от человека.
type Hub struct {
    mu          sync.RWMutex
    subscribers map[chan []byte]struct{}

    injMu      sync.Mutex
    injections map[string][]string // agentID → []message
}

// NewHub создаёт Hub.
func NewHub() *Hub {
    return &Hub{
        subscribers: make(map[chan []byte]struct{}),
        injections:  make(map[string][]string),
    }
}

// Subscribe регистрирует нового SSE-клиента и возвращает его канал.
func (h *Hub) Subscribe() chan []byte {
    ch := make(chan []byte, 32)
    h.mu.Lock()
    h.subscribers[ch] = struct{}{}
    h.mu.Unlock()
    return ch
}

// Unsubscribe удаляет клиента и закрывает его канал.
func (h *Hub) Unsubscribe(ch chan []byte) {
    h.mu.Lock()
    delete(h.subscribers, ch)
    h.mu.Unlock()
    close(ch)
}

// Broadcast пушит событие всем подключённым клиентам.
func (h *Hub) Broadcast(evt SSEEvent) {
    data, err := json.Marshal(evt)
    if err != nil {
        return
    }
    // Формат SSE: "data: {json}\n\n"
    msg := append([]byte("data: "), data...)
    msg = append(msg, '\n', '\n')

    h.mu.RLock()
    defer h.mu.RUnlock()
    for ch := range h.subscribers {
        select {
        case ch <- msg:
        default: // клиент не успевает читать — пропускаем
        }
    }
}

// Inject добавляет сообщение человека в очередь агента.
func (h *Hub) Inject(agentID, message string) {
    h.injMu.Lock()
    h.injections[agentID] = append(h.injections[agentID], message)
    h.injMu.Unlock()
}

// DrainInjections забирает и очищает очередь инъекций для агента.
func (h *Hub) DrainInjections(agentID string) []string {
    h.injMu.Lock()
    defer h.injMu.Unlock()
    msgs := h.injections[agentID]
    delete(h.injections, agentID)
    return msgs
}
```

**Step 2: Компиляция**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI && go build ./server/...
```

**Step 3: Commit**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI
git add server/internal/api/hub.go
git commit -m "feat(api): add SSE Hub for broadcast and injection queue"
```

---

### Task 4: SSE и Inject handlers

**Files:**
- Create: `server/internal/api/sse_handlers.go`
- Modify: `server/internal/api/handler.go`

**Step 1: Обновить `server/internal/api/handler.go` — добавить Hub**

Текущее содержимое handler.go:
```go
package api

import "milk/server/internal/storage"

type Handler struct {
    repo *storage.Repository
}

func NewHandler(repo *storage.Repository) *Handler {
    return &Handler{repo: repo}
}
```

Заменить на:
```go
package api

import "milk/server/internal/storage"

type Handler struct {
    repo *storage.Repository
    hub  *Hub
}

func NewHandler(repo *storage.Repository, hub *Hub) *Handler {
    return &Handler{repo: repo, hub: hub}
}
```

**Step 2: Создать `server/internal/api/sse_handlers.go`**

```go
package api

import (
    "encoding/json"
    "net/http"
)

// EventsStream — GET /events/stream
// SSE-эндпоинт. Держит соединение открытым и пушит события по мере поступления.
func (h *Handler) EventsStream(w http.ResponseWriter, r *http.Request) {
    // Проверяем поддержку flushing
    flusher, ok := w.(http.Flusher)
    if !ok {
        writeError(w, http.StatusInternalServerError, ErrCodeInternalError, "SSE not supported")
        return
    }

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    ch := h.hub.Subscribe()
    defer h.hub.Unsubscribe(ch)

    // Отправляем connected event
    w.Write([]byte("data: {\"type\":\"connected\"}\n\n"))
    flusher.Flush()

    for {
        select {
        case msg, ok := <-ch:
            if !ok {
                return
            }
            w.Write(msg)
            flusher.Flush()
        case <-r.Context().Done():
            return
        }
    }
}

// InjectMessage — POST /agents/{id}/inject
// Добавляет сообщение человека в очередь агента.
func (h *Handler) InjectMessage(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "missing agent id")
        return
    }

    var req InjectThoughtRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid JSON body")
        return
    }
    if req.Content == "" {
        writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "content is required")
        return
    }

    // Проверяем что агент существует
    rec, err := h.repo.GetAgentByID(id)
    if err != nil || rec == nil {
        writeError(w, http.StatusNotFound, ErrCodeNotFound, "agent not found")
        return
    }

    h.hub.Inject(id, req.Content)

    writeJSON(w, http.StatusOK, SuccessResponse{
        Success: true,
        Message: "message injected — agent will respond on next tick",
    })
}
```

**Step 3: Исправить `main.go` — `NewHandler` теперь принимает 2 аргумента**

В `server/cmd/server/main.go` строку:
```go
handler := api.NewHandler(repo)
```
заменить на:
```go
hub := api.NewHub()
handler := api.NewHandler(repo, hub)
```

**Step 4: Обновить `mux.go` — подключить реальные handlers**

В `server/cmd/server/middleware/mux.go`:
- `mux.HandleFunc("GET /events/stream", TODO)` → `mux.HandleFunc("GET /events/stream", h.EventsStream)`
- `mux.HandleFunc("POST /agents/{id}/inject", TODO)` → `mux.HandleFunc("POST /agents/{id}/inject", h.InjectMessage)`

**Step 5: Компиляция**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI && go build ./server/...
```

**Step 6: Commit**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI
git add server/internal/api/sse_handlers.go server/internal/api/handler.go \
        server/cmd/server/main.go server/cmd/server/middleware/mux.go
git commit -m "feat(api): SSE stream and inject handlers"
```

---

### Task 5: `GetRandomActiveAgents` в storage

**Files:**
- Modify: `server/internal/storage/sqlite.go`

**Step 1: Добавить метод в конец `server/internal/storage/sqlite.go`**

```go
// GetRandomActiveAgents возвращает до n случайных активных агентов.
// Используется оркестратором для выбора участников диалога.
func (r *Repository) GetRandomActiveAgents(n int) ([]AgentRecord, error) {
    query := `SELECT id, name, personality, mood_state, goals, state, is_active, created_at, last_active, snapshot
              FROM agents WHERE is_active = true ORDER BY RANDOM() LIMIT ?`
    rows, err := r.DB.Query(query, n)
    if err != nil {
        return nil, fmt.Errorf("GetRandomActiveAgents: %w", err)
    }
    defer rows.Close()

    var agents []AgentRecord
    for rows.Next() {
        var a AgentRecord
        if err := rows.Scan(
            &a.ID, &a.Name, &a.Personality, &a.MoodState, &a.Goals,
            &a.State, &a.IsActive, &a.CreatedAt, &a.LastActive, &a.Snapshot,
        ); err != nil {
            return nil, fmt.Errorf("GetRandomActiveAgents scan: %w", err)
        }
        agents = append(agents, a)
    }
    return agents, rows.Err()
}

// SaveConversationEvent сохраняет одну реплику диалога в таблицу events.
func (r *Repository) SaveConversationEvent(speakerID, targetID, content string, tick int64) error {
    payload := fmt.Sprintf(`{"speakerId":%q,"targetId":%q,"content":%q}`,
        speakerID, targetID, content)
    affectedJSON := fmt.Sprintf(`[%q,%q]`, speakerID, targetID)

    _, err := r.DB.Exec(
        `INSERT INTO events (id, topic, type, source, affected_agents, payload, status, tick, created_at)
         VALUES (?, 'interaction', 'conversation', ?, ?, ?, 'completed', ?, ?)`,
        newUUID(), speakerID, affectedJSON, payload, tick, timeNow(),
    )
    return err
}

// newUUID генерирует UUID v4.
func newUUID() string {
    return uuid.New().String()
}

// timeNow возвращает текущее UTC время.
func timeNow() time.Time {
    return time.Now().UTC()
}
```

**Step 2: Добавить import `"github.com/google/uuid"` в sqlite.go**

В блок imports добавить `"github.com/google/uuid"`.

**Step 3: Компиляция**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI && go build ./server/...
```

**Step 4: Commit**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI
git add server/internal/storage/sqlite.go
git commit -m "feat(storage): add GetRandomActiveAgents and SaveConversationEvent"
```

---

### Task 6: Orchestrator — тикер и диалоги

**Files:**
- Modify: `server/internal/world/orchestrator.go`

**Step 1: Добавить в конец `server/internal/world/orchestrator.go`**

Файл содержит типы `Orchestrator`, `WorldState`, `InteractionResult`, `DialogueTurn`. Нужно добавить реализацию.

```go
import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"

    "milk/server/internal/agent"
    "milk/server/internal/api"
    "milk/server/internal/storage"
    "milk/server/pkg/gigachat"
)

// OrchestratorV2 — реализация тикера диалогов.
// Использует другое имя чтобы не конфликтовать с существующим Orchestrator.
type OrchestratorV2 struct {
    repo      *storage.Repository
    llm       *gigachat.Client
    hub       *api.Hub
    tickInterval time.Duration
    turns     int
    currentTick int64
    mu        sync.Mutex
    cancel    context.CancelFunc
}

// NewOrchestrator создаёт OrchestratorV2.
func NewOrchestrator(repo *storage.Repository, llm *gigachat.Client, hub *api.Hub) *OrchestratorV2 {
    return &OrchestratorV2{
        repo:         repo,
        llm:          llm,
        hub:          hub,
        tickInterval: 15 * time.Second,
        turns:        4,
    }
}

// Start запускает фоновый тикер. Блокирует до отмены ctx.
func (o *OrchestratorV2) Start(ctx context.Context) {
    ctx, cancel := context.WithCancel(ctx)
    o.cancel = cancel

    log.Println("orchestrator: started, tick interval", o.tickInterval)

    ticker := time.NewTicker(o.tickInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            o.mu.Lock()
            o.currentTick++
            tick := o.currentTick
            o.mu.Unlock()

            go o.runTick(ctx, tick)

        case <-ctx.Done():
            log.Println("orchestrator: stopped")
            return
        }
    }
}

// Stop останавливает тикер.
func (o *OrchestratorV2) Stop() {
    if o.cancel != nil {
        o.cancel()
    }
}

func (o *OrchestratorV2) runTick(ctx context.Context, tick int64) {
    agents, err := o.repo.GetRandomActiveAgents(2)
    if err != nil {
        log.Printf("orchestrator tick %d: getAgents error: %v", tick, err)
        return
    }
    if len(agents) < 2 {
        log.Printf("orchestrator tick %d: not enough agents (%d)", tick, len(agents))
        return
    }

    log.Printf("orchestrator tick %d: %s <-> %s", tick, agents[0].Name, agents[1].Name)
    o.runConversation(ctx, agents[0], agents[1], tick)
}

func (o *OrchestratorV2) runConversation(
    ctx context.Context,
    a1, a2 storage.AgentRecord,
    tick int64,
) {
    p1 := parsePersonalityRec(a1.Personality)
    p2 := parsePersonalityRec(a2.Personality)

    brain1 := agent.NewBrain(&p1, o.llm)
    brain2 := agent.NewBrain(&p2, o.llm)

    mood1 := agent.MoodNeutral
    mood2 := agent.MoodNeutral

    goals1 := parseGoalsRec(a1.Goals.String)
    goals2 := parseGoalsRec(a2.Goals.String)

    // История диалога (формат OpenAI/Ollama): чередуем user/assistant
    // Для agent1: его реплики = assistant, реплики agent2 = user
    // Для agent2: наоборот
    var history1, history2 []gigachat.Message

    // Начальное сообщение: агент 1 видит агента 2
    opener := fmt.Sprintf(
        "You notice %s nearby. Start a conversation naturally based on your personality and current mood.",
        a2.Name,
    )
    history1 = append(history1, gigachat.Message{Role: "user", Content: opener})

    for i := 0; i < o.turns; i++ {
        select {
        case <-ctx.Done():
            return
        default:
        }

        // Ход агента 1
        injections1 := o.hub.DrainInjections(a1.ID)
        if len(injections1) > 0 {
            for _, inj := range injections1 {
                history1 = append(history1, gigachat.Message{
                    Role:    "user",
                    Content: fmt.Sprintf("[Human says to you]: %s", inj),
                })
            }
        }

        reply1, err := brain1.Think(ctx, o.llm, a1.Name, mood1, goals1, history1)
        if err != nil {
            log.Printf("orchestrator: brain1 error: %v", err)
            return
        }

        // Сохраняем и бродкастим
        o.repo.SaveConversationEvent(a1.ID, a2.ID, reply1, tick) //nolint:errcheck
        o.hub.Broadcast(api.SSEEvent{
            Type:    "conversation",
            Speaker: a1.Name,
            Target:  a2.Name,
            Content: reply1,
            AgentID: a1.ID,
            Tick:    tick,
        })

        // Обновляем истории
        history1 = append(history1, gigachat.Message{Role: "assistant", Content: reply1})
        history2 = append(history2, gigachat.Message{Role: "user",
            Content: fmt.Sprintf("%s says: %s", a1.Name, reply1)})

        // Ход агента 2 (только чётные итерации)
        if i%2 == 0 {
            injections2 := o.hub.DrainInjections(a2.ID)
            if len(injections2) > 0 {
                for _, inj := range injections2 {
                    history2 = append(history2, gigachat.Message{
                        Role:    "user",
                        Content: fmt.Sprintf("[Human says to you]: %s", inj),
                    })
                }
            }

            // Добавляем контекст если история пустая
            if len(history2) == 1 {
                history2 = append([]gigachat.Message{{
                    Role:    "user",
                    Content: fmt.Sprintf("%s is talking to you. Context: %s", a1.Name, reply1),
                }}, history2...)
            }

            reply2, err := brain2.Think(ctx, o.llm, a2.Name, mood2, goals2, history2)
            if err != nil {
                log.Printf("orchestrator: brain2 error: %v", err)
                return
            }

            o.repo.SaveConversationEvent(a2.ID, a1.ID, reply2, tick) //nolint:errcheck
            o.hub.Broadcast(api.SSEEvent{
                Type:    "conversation",
                Speaker: a2.Name,
                Target:  a1.Name,
                Content: reply2,
                AgentID: a2.ID,
                Tick:    tick,
            })

            history2 = append(history2, gigachat.Message{Role: "assistant", Content: reply2})
            history1 = append(history1, gigachat.Message{Role: "user",
                Content: fmt.Sprintf("%s says: %s", a2.Name, reply2)})
        }
    }
}

// parsePersonalityRec десериализует JSON-строку личности в Personality.
func parsePersonalityRec(raw string) agent.Personality {
    var p agent.Personality
    if raw != "" {
        json.Unmarshal([]byte(raw), &p) //nolint:errcheck
    }
    return p
}

// parseGoalsRec десериализует JSON-строку целей в []Goal.
func parseGoalsRec(raw string) []agent.Goal {
    if raw == "" {
        return nil
    }
    var goals []agent.Goal
    json.Unmarshal([]byte(raw), &goals) //nolint:errcheck
    return goals
}
```

**Step 2: Компиляция**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI && go build ./server/...
```

**Step 3: Commit**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI
git add server/internal/world/orchestrator.go
git commit -m "feat(orchestrator): implement ticker and conversation loop"
```

---

### Task 7: Wire всё в main.go

**Files:**
- Modify: `server/cmd/server/main.go`

**Step 1: Заменить содержимое `server/cmd/server/main.go`**

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/rs/cors"
    "milk/server/cmd/server/middleware"
    "milk/server/data"
    "milk/server/internal/api"
    "milk/server/internal/storage"
    "milk/server/internal/world"
    "milk/server/pkg/gigachat"
)

func main() {
    // 1. База данных
    data.DbConnection()
    repo := storage.NewRepository(data.Db)

    // 2. LLM клиент (Ollama)
    llmClient := gigachat.NewClient()
    fmt.Printf("LLM: %s @ %s\n", llmClient.Model, llmClient.BaseURL)

    // 3. SSE Hub
    hub := api.NewHub()

    // 4. HTTP Handler
    handler := api.NewHandler(repo, hub)
    mux := middleware.NewMux(handler)

    h := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:5173"},
        AllowedMethods: []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"},
        AllowedHeaders: []string{"Content-type"},
        Debug:          false,
    }).Handler(mux)

    // 5. Оркестратор — запускаем в фоне
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    orch := world.NewOrchestrator(repo, llmClient, hub)
    go orch.Start(ctx)

    // 6. Graceful shutdown
    go func() {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
        <-sig
        fmt.Println("\nshutting down...")
        cancel()
    }()

    fmt.Println("server starting on :8080")
    if err := http.ListenAndServe(":8080", h); err != nil {
        fmt.Println("server error:", err)
    }
}
```

**Step 2: Компиляция**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI && go build ./server/...
```

**Step 3: Commit**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI
git add server/cmd/server/main.go
git commit -m "feat: wire Orchestrator, LLM client, SSE hub in main"
```

---

### Task 8: Smoke-test

**Step 1: Убедиться что Ollama запущена с Gemma**

```bash
curl -s http://localhost:11434/api/tags | python3 -m json.tool
```
Expected: список моделей включает `gemma3` (или другую — тогда установить `OLLAMA_MODEL=имя`).

Если модели нет:
```bash
ollama pull gemma3
```

**Step 2: Запустить сервер**

```bash
cd /home/dasadno/Documents/Einstein-IslandAI
DB_PATH=server/data/society.db go run ./server/cmd/server/
```
Expected:
```
connected succesfully
LLM: gemma3 @ http://localhost:11434
orchestrator: started, tick interval 15s
server starting on :8080
```

**Step 3: Создать двух агентов (если БД пуста)**

```bash
curl -s -X POST localhost:8080/control/spawn \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","personality":{"openness":0.9,"extraversion":0.7,"agreeableness":0.6,"conscientiousness":0.5,"neuroticism":0.2,"coreValues":["curiosity","creativity"],"quirks":["talks to herself","loves philosophical questions"]}}'

curl -s -X POST localhost:8080/control/spawn \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob","personality":{"openness":0.4,"extraversion":0.3,"agreeableness":0.8,"conscientiousness":0.9,"neuroticism":0.3,"coreValues":["loyalty","honesty"],"quirks":["always speaks in metaphors"]}}'
```

**Step 4: Подключиться к SSE-стриму**

```bash
curl -N http://localhost:8080/events/stream
```
Expected через ~15 секунд:
```
data: {"type":"connected"}

data: {"type":"conversation","speaker":"Alice","target":"Bob","content":"Hey Bob! I've been wondering...","agentId":"...","tick":1}

data: {"type":"conversation","speaker":"Bob","target":"Alice","content":"That's like asking whether the river...","agentId":"...","tick":1}
```

**Step 5: Вписаться в диалог**

В другом терминале:
```bash
ALICE_ID=$(curl -s localhost:8080/agents | python3 -c "import sys,json; print([a['id'] for a in json.load(sys.stdin)['agents'] if a['name']=='Alice'][0])")
curl -X POST localhost:8080/agents/$ALICE_ID/inject \
  -H "Content-Type: application/json" \
  -d '{"type":"message","content":"Alice, do you believe in free will?"}'
```
Expected: в SSE-стриме следующая реплика Alice учтёт вопрос человека.

**Step 6: Проверить что events сохраняются**

```bash
curl -s localhost:8080/events | python3 -m json.tool
```
Expected: список событий типа `conversation`.

---

## Итоговые файлы

| Файл | Изменение |
|------|-----------|
| `server/pkg/gigachat/client.go` | Реализован Ollama HTTP-клиент |
| `server/internal/agent/brain.go` | NewBrain, buildSystemPrompt, Think |
| `server/internal/api/hub.go` | SSE Hub (broadcast + inject queue) |
| `server/internal/api/handler.go` | Добавлен hub |
| `server/internal/api/sse_handlers.go` | EventsStream, InjectMessage |
| `server/internal/storage/sqlite.go` | GetRandomActiveAgents, SaveConversationEvent |
| `server/internal/world/orchestrator.go` | OrchestratorV2, Start, runConversation |
| `server/cmd/server/main.go` | LLM + Hub + Orchestrator |
| `server/cmd/server/middleware/mux.go` | SSE и inject маршруты |
