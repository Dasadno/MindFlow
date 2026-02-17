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
	// 1. База данных
	data.DbConnection()
	repo := storage.NewRepository(data.Db)

	// 2. LLM клиент (Ollama)
	llmClient := llm.NewClient()
	fmt.Printf("LLM: %s @ %s\n", llmClient.Model, llmClient.BaseURL)

	// 3. SSE Hub
	hub := api.NewHub()

	// 4. HTTP Handler
	handler := api.NewHandler(repo, hub)
	mux := api.NewMux(handler)

	h := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{"Content-type"},
		Debug:          false,
	}).Handler(mux)

	// 5. Оркестратор — запускаем в фоне
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := world.NewOrchestrator(repo, llmClient, hub)
	go orch.Start(ctx)

	// 6. Graceful shutdown
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
