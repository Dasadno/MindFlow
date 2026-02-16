package main

import (
	"fmt"
	"net/http"

	"github.com/einstein-islandai/ai-agent-society/server/data"
	"github.com/einstein-islandai/ai-agent-society/server/internal/api"
	"github.com/einstein-islandai/ai-agent-society/server/internal/storage"

	"github.com/rs/cors"
)

func main() {
	data.DbConnection()

	repo := storage.NewRepository(data.Db)
	handler := api.NewHandler(repo)

	mux := api.NewMux(handler)

	h := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{"Content-type"},
		Debug:          true,
	}).Handler(mux)

	fmt.Println("server starting on :8080")
	if err := http.ListenAndServe(":8080", h); err != nil {
		fmt.Println("failed to connect server: ", err)
	}
}
