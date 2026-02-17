package api

import (
	"encoding/json"
	"sync"
)

// SSEEvent — событие, отправляемое клиентам через SSE.
type SSEEvent struct {
	Type    string `json:"type"`
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
	msg := append([]byte("data: "), data...)
	msg = append(msg, '\n', '\n')

	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subscribers {
		select {
		case ch <- msg:
		default:
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
