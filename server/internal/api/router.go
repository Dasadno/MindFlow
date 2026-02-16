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
