package middleware

import (
	"net/http"
)

func TODO(w http.ResponseWriter, r *http.Request) { // just TODO function
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("endpoint in work"))
}

func NewMux() *http.ServeMux {

	mux := http.NewServeMux()

	// AGENTS
	mux.HandleFunc("GET /agents", TODO)
	mux.HandleFunc("GET /agents/{id}/memory", TODO)
	mux.HandleFunc("GET /agents/{id}/thoughts", TODO)
	mux.HandleFunc("GET /agents/{id}", TODO)
	mux.HandleFunc("POST /agents/{id}/inject", TODO)

	//RELATIONSHIPS
	mux.HandleFunc("GET /relationships", TODO)
	mux.HandleFunc("GET /relationships/{agentId}", TODO)
	mux.HandleFunc("POST /relationships", TODO)

	//EVENTS
	mux.HandleFunc("GET /events", TODO)
	mux.HandleFunc("POST /events", TODO)
	mux.HandleFunc("GET /events/stream", TODO)

	// WORLD
	mux.HandleFunc("GET /world/status", TODO)
	mux.HandleFunc("POST  /world/control", TODO)
	mux.HandleFunc("GET /world/statistics", TODO)

	//CONTROL PANEL
	mux.HandleFunc("POST /control/spawn", TODO)
	mux.HandleFunc("DELETE /control/agents/{id}", TODO)
	mux.HandleFunc("POST /control/reset", TODO)

	return mux

}
