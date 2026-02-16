package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"milk/server/internal/storage"
)

// SpawnAgent — POST /control/spawn
func (h *Handler) SpawnAgent(w http.ResponseWriter, r *http.Request) {
	var req SpawnAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "name is required")
		return
	}

	personalityJSON, err := json.Marshal(req.Personality)
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, "failed to encode personality")
		return
	}

	rec := storage.AgentRecord{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Personality: string(personalityJSON),
		State:       "idle",
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}

	if err := h.repo.CreateAgent(rec); err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, "failed to create agent")
		return
	}

	writeJSON(w, http.StatusCreated, AgentDetailResponse{
		ID:          rec.ID,
		Name:        rec.Name,
		Personality: req.Personality,
		CurrentMood: MoodDTO{Label: "neutral"},
		Goals:       []GoalDTO{},
		CreatedAt:   rec.CreatedAt,
	})
}

// DeactivateAgentHandler — DELETE /control/agents/{id}
func (h *Handler) DeactivateAgentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "missing agent id")
		return
	}

	if err := h.repo.DeactivateAgent(id); err != nil {
		writeError(w, http.StatusNotFound, ErrCodeNotFound, "agent not found")
		return
	}

	writeJSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Message: "agent deactivated",
	})
}
