// Package main is the entry point for the AI Agent Society server.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This file bootstraps the entire AI Agent Society application. It initializes
// all core subsystems and starts the HTTP server for the dashboard API.
//
// =============================================================================
// INITIALIZATION SEQUENCE (to be implemented):
// =============================================================================
//
// 1. CONFIGURATION LOADING
//    - Load environment variables or config file
//    - Parse command-line flags (port, db path, GigaChat endpoint)
//    - Validate configuration completeness
//
// 2. DATABASE INITIALIZATION
//    - Initialize SQLite connection pool via internal/storage
//    - Run database migrations (create tables if not exist)
//    - Initialize Vector DB (Chromem-go) for episodic memory
//
// 3. CORE SERVICES BOOTSTRAP
//    - Create GigaChat client (pkg/gigachat) with configured endpoint
//    - Initialize the World Orchestrator (internal/world)
//    - Create Agent Factory with personality templates
//
// 4. AGENT SPAWNING
//    - Load existing agents from database OR
//    - Create initial agent population with diverse personalities
//    - Register agents with the World Orchestrator
//
// 5. API SERVER STARTUP
//    - Initialize Gin router with middleware (CORS, logging, recovery)
//    - Register all API handlers (agents, relationships, events, control)
//    - Start HTTP server on configured port
//
// 6. WORLD SIMULATION LOOP
//    - Start background goroutine for world tick processing
//    - Each tick: process agent actions, update moods, handle interactions
//    - Graceful shutdown handling with context cancellation
//
// =============================================================================
// GRACEFUL SHUTDOWN:
// =============================================================================
// - Listen for SIGINT/SIGTERM signals
// - Stop world simulation loop
// - Flush pending memory writes to database
// - Close database connections
// - Shutdown HTTP server with timeout
//
// =============================================================================
// DEPENDENCIES:
// =============================================================================
// - internal/api:     HTTP handlers and routing
// - internal/agent:   Agent brain, memory, emotions
// - internal/world:   Orchestrator and event bus
// - internal/storage: Database layer (SQLite + Vector DB)
// - pkg/gigachat:     LLM client for GigaChat/Ollama

package main

import (
	"fmt"
	"milk/server/cmd/server/middleware"
	"milk/server/data"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	data.DbConnection()

	mux := middleware.NewMux()
	_ = mux //gtg rn, i'll comeback soon

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{"Content-type"},
		Debug:          true,
	}).Handler(mux)

	fmt.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("failed to connect server: ", err)
	}

}
