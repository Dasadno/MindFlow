package api

import (
	"net/http"
	"strconv"
	"time"
)

// serverStartTime хранит время запуска для вычисления uptime.
var serverStartTime = time.Now()

// GetWorldStatus — GET /world/status
func (h *Handler) GetWorldStatus(w http.ResponseWriter, r *http.Request) {
	activeAgents, err := h.repo.CountActiveAgents()
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, "failed to count agents")
		return
	}

	totalEvents, err := h.repo.CountTotalEvents()
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, "failed to count events")
		return
	}

	state, err := h.repo.GetWorldState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, "failed to get world state")
		return
	}

	currentTick, _ := strconv.ParseInt(state["current_tick"], 10, 64)
	simSpeed := 1.0
	if s, ok := state["simulation_speed"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			simSpeed = v
		}
	}
	isPaused := state["is_paused"] == "true"

	writeJSON(w, http.StatusOK, WorldStatusResponse{
		CurrentTick:     currentTick,
		SimulationSpeed: simSpeed,
		IsPaused:        isPaused,
		ActiveAgents:    activeAgents,
		TotalEvents:     totalEvents,
		Uptime:          formatUptime(time.Since(serverStartTime)),
	})
}

// formatUptime форматирует duration в строку "2h15m30s".
func formatUptime(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return strconv.Itoa(h) + "h" + strconv.Itoa(m) + "m" + strconv.Itoa(s) + "s"
	}
	if m > 0 {
		return strconv.Itoa(m) + "m" + strconv.Itoa(s) + "s"
	}
	return strconv.Itoa(s) + "s"
}
