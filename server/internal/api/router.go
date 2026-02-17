// Package api provides HTTP handlers and routing for the AI Agent Society dashboard.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This file configures the main HTTP router and middleware stack.
// It acts as the central routing configuration for all API endpoints.
//
// =============================================================================
// ROUTER SETUP (to be implemented):
// =============================================================================
//
// func NewRouter(deps *Dependencies) *gin.Engine
//
// MIDDLEWARE STACK:
// 1. Recovery middleware (panic recovery with logging)
// 2. Logger middleware (request/response logging)
// 3. CORS middleware (allow frontend dashboard access)
// 4. Request ID middleware (trace requests across logs)
//
// =============================================================================
// ROUTE GROUPS:
// =============================================================================
//
// /api/v1
// ├── /agents              → AgentHandler (CRUD, state, actions)
// ├── /relationships       → RelationshipHandler (graph, connections)
// ├── /events              → EventHandler (world events, history)
// ├── /world               → WorldHandler (time, simulation control)
// └── /control             → ControlHandler (admin operations)
//
// =============================================================================
// DEPENDENCIES INJECTION:
// =============================================================================
//
// type Dependencies struct {
//     AgentService       *agent.Service
//     WorldOrchestrator  *world.Orchestrator
//     Storage            *storage.Repository
//     GigaChatClient     *gigachat.Client
// }
//
// All handlers receive dependencies through constructor injection,
// enabling testability and loose coupling.

package api

import (
	"net/http"
)

func TODO(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("endpoint in work"))
}

func NewMux(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// AGENTS
	mux.HandleFunc("GET /agents", h.ListAgents)
	mux.HandleFunc("GET /agents/{id}/memory", h.GetAgentMemories)
	mux.HandleFunc("GET /agents/{id}/thoughts", TODO)
	mux.HandleFunc("GET /agents/{id}", h.GetAgent)
	mux.HandleFunc("POST /agents/{id}/inject", TODO)

	// RELATIONSHIPS
	mux.HandleFunc("GET /relationships", TODO)
	mux.HandleFunc("GET /relationships/{agentId}", TODO)
	mux.HandleFunc("POST /relationships", TODO)

	// EVENTS
	mux.HandleFunc("GET /events", TODO)
	mux.HandleFunc("POST /events", TODO)
	mux.HandleFunc("GET /events/stream", TODO)

	// WORLD
	mux.HandleFunc("GET /world/status", h.GetWorldStatus)
	mux.HandleFunc("POST /world/control", TODO)
	mux.HandleFunc("GET /world/statistics", TODO)

	// CONTROL PANEL
	mux.HandleFunc("POST /control/spawn", h.SpawnAgent)
	mux.HandleFunc("DELETE /control/agents/{id}", h.DeactivateAgentHandler)
	mux.HandleFunc("POST /control/reset", TODO)

	return mux
}
