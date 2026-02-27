package api

import (
	"net/http"
	"os"
	"strings"
)

func TODO(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("endpoint in work"))
}

func NewMux(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// =============================================================================
	// API ENDPOINTS
	// =============================================================================

	// AGENTS
	mux.HandleFunc("GET /agents", h.ListAgents)
	mux.HandleFunc("GET /agents/{id}/memory", TODO)
	mux.HandleFunc("GET /agents/{id}/thoughts", TODO)
	mux.HandleFunc("GET /agents/{id}", h.GetAgent)
	mux.HandleFunc("POST /agents/{id}/inject", h.InjectMessage)

	// RELATIONSHIPS
	mux.HandleFunc("GET /relationships", TODO)
	mux.HandleFunc("GET /relationships/{agentId}", TODO)
	mux.HandleFunc("POST /relationships", TODO)

	// EVENTS
	mux.HandleFunc("GET /events", TODO)
	mux.HandleFunc("POST /events", TODO)
	mux.HandleFunc("GET /events/stream", h.EventsStream)

	// WORLD
	mux.HandleFunc("GET /world/status", h.GetWorldStatus)
	mux.HandleFunc("POST /world/control", TODO)
	mux.HandleFunc("GET /world/statistics", TODO)

	// CONTROL PANEL
	mux.HandleFunc("POST /control/spawn", h.SpawnAgent)
	mux.HandleFunc("DELETE /control/agents/{id}", h.DeactivateAgentHandler)
	mux.HandleFunc("POST /control/reset", TODO)

	// =============================================================================
	// FRONTEND STATIC SERVING (Vite/React SPA support)
	// =============================================================================

	distPath := "web/dist"
	fs := http.FileServer(http.Dir(distPath))

	// Этот обработчик перехватывает всё, что не подошло под API выше
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		// Проверяем наличие физического файла (статические ассеты, картинки)
		fpath := distPath + r.URL.Path
		info, err := os.Stat(fpath)

		// Если файла не существует или это папка — отдаем index.html (SPA Fallback)
		if os.IsNotExist(err) || info.IsDir() {
			http.ServeFile(w, r, distPath+"/index.html")
			return
		}

		// В остальных случаях отдаем файл через стандартный FileServer
		fs.ServeHTTP(w, r)
	})

	return mux
}
