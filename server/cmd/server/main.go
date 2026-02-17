package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"milk/server/data"
	"milk/server/internal/api"
	"milk/server/internal/storage"
	"milk/server/internal/world"
	"milk/server/pkg/llm"

	"github.com/rs/cors"
)

func main() {
	// База данных
	data.DbConnection()
	repo := storage.NewRepository(data.Db)

	// LLM клиент (Ollama)
	llmClient := llm.NewClient()
	fmt.Printf("LLM: %s @ %s\n", llmClient.Model, llmClient.BaseURL)

	// SSE Hub
	hub := api.NewHub()

	// HTTP Handler
	handler := api.NewHandler(repo, hub)
	mux := api.NewMux(handler)

	origin := os.Getenv("ALLOWED_ORIGIN")
	h := cors.New(cors.Options{
		AllowedOrigins: []string{origin},
		AllowedMethods: []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{"Content-type"},
		Debug:          false,
	}).Handler(mux)

	// Оркестратор
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := world.NewOrchestrator(repo, llmClient, hub)
	go orch.Start(ctx)

	// Graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		fmt.Println("\nshutting down...")
		cancel()
		os.Exit(0)
	}()

	fmt.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", h); err != nil {
		fmt.Println("server error:", err)
	}
}
