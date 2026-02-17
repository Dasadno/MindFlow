package api

import (
	"encoding/json"
	"milk/server/internal/storage"
	"net/http"
	"strconv"
	"time"
)

// ListAgents — GET /agents
// Query params: ?page=1&limit=20&active=true
func (h *Handler) ListAgents(w http.ResponseWriter, r *http.Request) {
	filter := storage.AgentFilter{
		Page:  parseIntQuery(r, "page", 1),
		Limit: parseIntQuery(r, "limit", 20),
	}

	if activeStr := r.URL.Query().Get("active"); activeStr != "" {
		v, err := strconv.ParseBool(activeStr)
		if err == nil {
			filter.IsActive = &v
		}
	}

	records, total, err := h.repo.ListAgents(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, "failed to list agents")
		return
	}

	summaries := make([]AgentSummary, 0, len(records))
	for _, rec := range records {
		summaries = append(summaries, recordToSummary(rec))
	}

	writeJSON(w, http.StatusOK, AgentListResponse{
		Agents: summaries,
		Pagination: PaginationMeta{
			Page:  filter.Page,
			Limit: filter.Limit,
			Total: total,
		},
	})
}

// GetMemory - GET /agents/{id}/memory
func (h *Handler) GetAgentMemories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	memType := query.Get("type") // episodic | semantic | procedural
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 50
	}

	agentID := r.PathValue("id")

	memories, err := h.repo.MemoriesByAgent(agentID, memType, limit)
	if err != nil {
		http.Error(w, "Failed to fetch memories", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, memories)
}

// GetAgent — GET /agents/{id}
func (h *Handler) GetAgent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, "missing agent id")
		return
	}

	rec, err := h.repo.GetAgentByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, "failed to get agent")
		return
	}
	if rec == nil {
		writeError(w, http.StatusNotFound, ErrCodeNotFound, "agent not found")
		return
	}

	memories, _ := h.repo.CountMemoriesByAgent(id)
	relations, _ := h.repo.CountRelationshipsByAgent(id)
	interactions, _ := h.repo.CountInteractionsByAgent(id)
	days := int(time.Since(rec.CreatedAt).Hours() / 24)

	personality := parsePersonality(rec.Personality)
	mood := parseMood(rec.MoodState.String)
	goals := parseGoals(rec.Goals.String)

	detail := AgentDetailResponse{
		ID:          rec.ID,
		Name:        rec.Name,
		Personality: personality,
		CurrentMood: mood,
		Goals:       goals,
		Stats: AgentStatsDTO{
			TotalInteractions:  interactions,
			MemoriesCount:      memories,
			RelationshipsCount: relations,
			DaysSinceCreation:  days,
		},
		CreatedAt: rec.CreatedAt,
	}
	writeJSON(w, http.StatusOK, detail)
}

// ── Вспомогательные функции ──────────────────────────────────────────────────

func recordToSummary(rec storage.AgentRecord) AgentSummary {
	mood, intensity := moodFromJSON(rec.MoodState.String)
	personality := parsePersonality(rec.Personality)
	return AgentSummary{
		ID:              rec.ID,
		Name:            rec.Name,
		PersonalityType: personalityType(personality),
		CurrentMood:     mood,
		MoodIntensity:   intensity,
		IsActive:        rec.IsActive,
	}
}

func personalityType(p PersonalityDTO) string {
	if p.Openness > 0.7 {
		return "explorer"
	}
	if p.Extraversion < 0.3 {
		return "introvert"
	}
	return "guardian"
}

func moodFromJSON(moodJSON string) (string, float64) {
	if moodJSON == "" {
		return "neutral", 0.5
	}
	var pad struct {
		Label    string  `json:"label"`
		Pleasure float64 `json:"pleasure"`
	}
	if err := json.Unmarshal([]byte(moodJSON), &pad); err != nil {
		return "neutral", 0.5
	}
	label := pad.Label
	if label == "" {
		label = "neutral"
	}
	intensity := pad.Pleasure
	if intensity < 0 {
		intensity = -intensity
	}
	if intensity == 0 {
		intensity = 0.5
	}
	return label, intensity
}

func parsePersonality(raw string) PersonalityDTO {
	var p PersonalityDTO
	if raw == "" {
		return p
	}
	_ = json.Unmarshal([]byte(raw), &p)
	return p
}

func parseMood(raw string) MoodDTO {
	if raw == "" {
		return MoodDTO{Label: "neutral"}
	}
	var m MoodDTO
	_ = json.Unmarshal([]byte(raw), &m)
	if m.Label == "" {
		m.Label = "neutral"
	}
	return m
}

func parseGoals(raw string) []GoalDTO {
	if raw == "" {
		return []GoalDTO{}
	}
	var goals []GoalDTO
	_ = json.Unmarshal([]byte(raw), &goals)
	if goals == nil {
		return []GoalDTO{}
	}
	return goals
}

func parseIntQuery(r *http.Request, key string, def int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil || v <= 0 {
		return def
	}
	return v
}
