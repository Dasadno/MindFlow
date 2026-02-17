# AI Agent Society — Backend Server

Autonomous AI agents with personalities, emotions, memory, and social relationships.

**Stack:** Go 1.22 | SQLite (modernc.org/sqlite) | Ollama (Gemma 3 4B)

**Repository:** [Milk-IslandAI](https://github.com/Milk-IslandAI)

---

## Quick Start

```bash
# 1. Установить и запустить Ollama
ollama pull gemma3:4b
ollama serve

# 2. Запустить сервер (из корня проекта)
go run ./server/cmd/server/

# Или с кастомной конфигурацией
OLLAMA_URL=http://localhost:11434 OLLAMA_MODEL=gemma3:4b DB_PATH=server/data/society.db go run ./server/cmd/server/
```

Default port: `:8080`

```bash
# 3. Создать двух агентов
curl -X POST localhost:8080/control/spawn \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","personality":{"openness":0.9,"extraversion":0.7,"agreeableness":0.6,"conscientiousness":0.5,"neuroticism":0.2,"coreValues":["curiosity"],"quirks":["talks to herself"]}}'

curl -X POST localhost:8080/control/spawn \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob","personality":{"openness":0.4,"extraversion":0.3,"agreeableness":0.8,"conscientiousness":0.9,"neuroticism":0.3,"coreValues":["loyalty"],"quirks":["always speaks in metaphors"]}}'

# 4. Подключиться к стриму диалогов (SSE)
curl -N http://localhost:8080/events/stream

# 5. Вписаться в разговор
curl -X POST localhost:8080/agents/{alice_id}/inject \
  -H "Content-Type: application/json" \
  -d '{"type":"message","content":"Alice, do you believe in fate?"}'
```

### Конфигурация

| Переменная | Default | Описание |
|------------|---------|----------|
| `OLLAMA_URL` | `http://localhost:11434` | Адрес Ollama сервера |
| `OLLAMA_MODEL` | `gemma3:4b` | Модель для диалогов |
| `DB_PATH` | `server/data/society.db` | Путь к SQLite базе |

**Захардкоженные параметры** (пока не читаются из env):
- Tick interval: `15s` (в `orchestrator.go`)
- Conversation turns: `4` (в `orchestrator.go`)
- CORS origin: `http://localhost:5173` (в `main.go`)

---

## Agent Intelligence Flow (Gemma via Ollama)

Этот раздел объясняет как агенты «живут»: от создания до диалога, стрима мыслей и подключения человека.

### Общая картина

```
┌─────────────────────────────────────────────────────────────────────┐
│                          SERVER (Go)                                │
│                                                                     │
│  ┌──────────┐    spawn     ┌─────────────────────────────────────┐ │
│  │  Client  │ ──────────► │           agents (SQLite)           │ │
│  │ (HTTP)   │             │  id, name, personality, mood, goals  │ │
│  └──────────┘             └─────────────────┬───────────────────┘ │
│       ▲                                     │ load active agents   │
│       │ SSE                                 ▼                      │
│  ┌────┴─────┐             ┌─────────────────────────────────────┐ │
│  │  /events │◄────────────│         Orchestrator (ticker)       │ │
│  │  /stream │             │  every 15s: pick 2 agents → talk    │ │
│  └──────────┘             └─────────────────┬───────────────────┘ │
│                                             │ runConversation()    │
│  ┌──────────┐  inject      ▼                ▼                      │
│  │  Client  │ ──────► ┌────────┐    ┌─────────────────────────┐  │
│  │ (HTTP)   │         │ Brain  │    │      Ollama (Gemma)      │  │
│  └──────────┘         │Think() │◄──►│  POST /api/chat          │  │
│                        └────────┘    │  model: gemma3:4b        │  │
│                                      └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

---

### 1. Создание агента (`POST /control/spawn`)

Клиент отправляет имя и черты личности (Big Five). Сервер:

1. Генерирует UUID
2. Сохраняет агента в SQLite (`agents` таблица)
3. Возвращает профиль агента

На этом этапе агент — просто запись в БД. LLM не вызывается.

```
Client → POST /control/spawn { name, personality }
       ← 201 { id, name, personality, currentMood: "neutral" }
```

---

### 2. Оркестратор — сердце симуляции

При старте сервера запускается фоновая горутина — **Orchestrator**. Каждые **15 секунд** она:

```
tick:
  1. Загружает список активных агентов из БД
  2. Если агентов < 2 → пропускает тик
  3. Случайно выбирает агента-инициатора (Agent A) и собеседника (Agent B)
  4. Вызывает runConversation(A, B)
  5. Инкрементирует счётчик тиков в world_state
```

---

### 3. Диалог двух агентов (`runConversation`)

Это ключевой процесс. Каждый диалог — 4 реплики (2 хода каждого агента).

```
Шаг 1: Формируем системный промпт для Agent A
───────────────────────────────────────────────
You are Alice. You are curious and open-minded (openness: 0.85).
You value honesty and creativity.
Your quirk: you talk to yourself.
Your current mood: excited.
Your goals: make a new friend, explore the eastern zone.
Keep responses concise (2-3 sentences). Stay in character.

Шаг 2: Формируем начальный контекст (user message)
───────────────────────────────────────────────────
You see Bob nearby. Bob seems calm today.
What do you say to start a conversation?

Шаг 3: Отправляем в Ollama → получаем реплику Alice
────────────────────────────────────────────────────
POST http://localhost:11434/api/chat
{ "model": "gemma3:4b", "messages": [system, user], "stream": false }
← "Hey Bob! You seem thoughtful today. What's on your mind?"

Шаг 4: Сохраняем реплику в events таблицу
──────────────────────────────────────────
INSERT INTO events (topic="interaction", type="conversation",
  source=alice_id, affected_agents=[alice_id, bob_id],
  payload={ speaker:"Alice", content:"Hey Bob!...", tick:42 })

Шаг 5: Пушим реплику в SSE-канал (все подключённые клиенты получают её мгновенно)

Шаг 6: Формируем промпт для Bob (история диалога передаётся как контекст)
──────────────────────────────────────────────────────────────────────────
System: [Bob's personality prompt]
User: Alice says to you: "Hey Bob! You seem thoughtful today. What's on your mind?"
      What do you reply?

← Bob: "Actually, I've been thinking about our last conversation about trust..."

Шаг 7-8: Повторяем для Alice (2-й ход) и Bob (2-й ход)
```

Итого: 4 LLM-вызова на диалог, ~15–30 секунд суммарно при локальном Gemma.

---

### 4. SSE стрим — как фронтенд видит диалоги

Клиент подключается один раз:

```
GET /events/stream
Accept: text/event-stream

← data: {"type":"conversation","speaker":"Alice","target":"Bob","content":"Hey Bob!...","tick":42}
← data: {"type":"conversation","speaker":"Bob","target":"Alice","content":"Actually, I've been thinking...","tick":42}
← data: {"type":"conversation","speaker":"Alice","target":"Bob","content":"Trust is everything...","tick":42}
...
```

Соединение остаётся открытым. Каждая новая реплика приходит по мере генерации.

**Стрим мыслей конкретного агента** (`GET /agents/{id}/thoughts`) — то же самое, но только события типа `thought` от этого агента. Используется для панели «внутренний монолог» на дашборде.

---

### 5. Человек присоединяется к диалогу

Если пользователь хочет написать агенту:

```
POST /agents/{alice_id}/inject
{ "type": "message", "content": "Alice, what do you think about the concept of free will?" }
```

Сервер сохраняет инъекцию. При следующем тике оркестратор:
1. Видит входящее сообщение в очереди агента
2. Добавляет его в контекст: `User (human) says: "Alice, what do you think about..."`
3. Gemma генерирует ответ Alice с учётом вопроса
4. Ответ уходит в SSE-стрим

Таким образом человек органично вписывается в диалог, не ломая его логику.

---

### 6. Полный data flow

```
[Human] POST /control/spawn
           │
           ▼
    [SQLite: agents] ←─────────────────────────────────────────┐
           │                                                    │
           │ (every 15s)                                        │
           ▼                                                    │
    [Orchestrator]                                              │
    pickTwoAgents() ──► loadFromDB()                           │
           │                                                    │
           ▼                                                    │
    runConversation(A, B)                                       │
           │                                                    │
           ├──► Brain.buildSystemPrompt(A) ──► string           │
           │                                                    │
           ├──► POST ollama/api/chat ──► reply_A               │
           │                                                    │
           ├──► saveEvent(reply_A) ─────────────────────────────┘
           │
           ├──► sseChannel ◄── GET /events/stream ◄── [Frontend]
           │
           ├──► Brain.buildSystemPrompt(B) ──► string
           │
           └──► POST ollama/api/chat ──► reply_B ──► saveEvent ──► sseChannel

[Human] POST /agents/{id}/inject ──► injectionQueue[agentId]
                                            │
                                     (next tick)
                                            │
                                            ▼
                                   added to conversation context
```

---

## API Endpoints

Base URL: `http://localhost:8080` (no `/api/v1` prefix — routes are at root level)

### Agents

| Method | Path | Status | Description | Response Schema |
|--------|------|--------|-------------|-----------------|
| `GET` | `/agents` | Done | List all agents (paginated, filterable) | `{ "agents": [AgentSummary], "pagination": {...} }` |
| `GET` | `/agents/{id}` | Done | Get detailed agent profile | `AgentDetail` (see below) |
| `GET` | `/agents/{id}/memory` | TODO | Get agent memories | `{ "memories": [MemoryEntry], "summary": "..." }` |
| `GET` | `/agents/{id}/thoughts` | TODO | Stream agent's live thought process | SSE |
| `POST` | `/agents/{id}/inject` | Done | Inject message into agent's conversation | Request: `{ "type": "message", "content": "..." }` |

### Relationships

| Method | Path | Status | Description | Response Schema |
|--------|------|--------|-------------|-----------------|
| `GET` | `/relationships` | TODO | Get full relationship graph | `RelationshipGraph` |
| `GET` | `/relationships/{agentId}` | TODO | Get one agent's relationships | `{ "relationships": [...] }` |
| `POST` | `/relationships` | TODO | Force-create/modify a relationship | Request: `{ "agent1": "id", "agent2": "id", "type": "friend" }` |

### Events

| Method | Path | Status | Description | Response Schema |
|--------|------|--------|-------------|-----------------|
| `GET` | `/events` | TODO | List world events (filterable) | `{ "events": [Event] }` |
| `POST` | `/events` | TODO | Inject global event | Request: `{ "type": "disaster", "description": "...", "affectedAgents": ["all"] }` |
| `GET` | `/events/stream` | Done | **SSE** — real-time conversation stream | `data: {"type":"conversation","speaker":"Alice","target":"Bob","content":"...","tick":42}` |

### World

| Method | Path | Status | Description | Response Schema |
|--------|------|--------|-------------|-----------------|
| `GET` | `/world/status` | Done | Current simulation state | `WorldStatus` (see below) |
| `POST` | `/world/control` | TODO | Pause/resume/step/set speed | Request: `{ "action": "pause" }` |
| `GET` | `/world/statistics` | TODO | Aggregate statistics | `WorldStatistics` |

### Control Panel (Admin)

| Method | Path | Status | Description | Response Schema |
|--------|------|--------|-------------|-----------------|
| `POST` | `/control/spawn` | Done | Spawn new agent | Request: `{ "name": "Eve", "personality": {...} }` |
| `DELETE` | `/control/agents/{id}` | Done | Remove agent (soft delete) | `{ "success": true, "message": "agent deactivated" }` |
| `POST` | `/control/reset` | TODO | Reset world state | Request: `{ "confirm": true, "preserveAgents": false }` |

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

---

## Database Schema

**Engine:** SQLite | **File:** `server/data/society.db`

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

## Project Structure

```
server/
├── cmd/server/main.go              # Entry point, DI, graceful shutdown
├── data/
│   ├── db.go                       # Global SQLite connection (modernc.org/sqlite)
│   └── society.db                  # SQLite database file
├── internal/
│   ├── api/
│   │   ├── router.go               # http.ServeMux router, route registration
│   │   ├── handler.go              # Handler struct (repo + hub DI)
│   │   ├── agent_handlers.go       # GET /agents, GET /agents/{id}
│   │   ├── control_handlers.go     # POST /control/spawn, DELETE /control/agents/{id}
│   │   ├── world_handlers.go       # GET /world/status
│   │   ├── sse_handlers.go         # GET /events/stream (SSE), POST /agents/{id}/inject
│   │   ├── hub.go                  # SSE broadcast hub + human injection queue
│   │   ├── helpers.go              # writeJSON, writeError utilities
│   │   ├── middleware.go           # APIError type, error codes
│   │   └── dto.go                  # Request/response JSON schemas
│   ├── agent/
│   │   ├── agent.go                # Core agent entity, types (Agent, Personality, Goal, Stimulus)
│   │   ├── brain.go                # LLM-driven cognition, system prompt builder
│   │   ├── memory.go               # Three-tier memory model (working/episodic/semantic)
│   │   └── emotions.go             # PAD + discrete emotion engine, appraisal types
│   ├── world/
│   │   ├── orchestrator.go         # Simulation ticker: 15s interval, 4-turn conversations
│   │   ├── time.go                 # Simulation clock types (not yet integrated)
│   │   └── eventbus.go             # Pub-sub event types (not yet integrated)
│   └── storage/
│       ├── sqlite.go               # SQLite Repository: CRUD agents, events, counts
│       └── vectordb.go             # Vector store types for episodic memory (not yet integrated)
├── pkg/llm/
│   └── client.go                   # Ollama HTTP client (POST /api/chat)
└── server                          # Compiled binary
```

---

# Client side of the project
