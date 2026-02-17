package api

import (
	"encoding/json"
	"net/http"
)

// EventsStream — GET /events/stream
// SSE-эндпоинт, держит соединение и пушит события.
func (h *Handler) EventsStream(w http.ResponseWriter, r *http.Request) {
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
