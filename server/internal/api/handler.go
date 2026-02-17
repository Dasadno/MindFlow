// server/internal/api/handler.go
package api

import "milk/server/internal/storage"

// Handler — HTTP-обработчики с доступом к хранилищу и SSE hub.
type Handler struct {
	repo *storage.Repository
	hub  *Hub
}

// NewHandler создаёт Handler с инъекцией зависимости Repository и SSE Hub.
func NewHandler(repo *storage.Repository, hub *Hub) *Handler {
	return &Handler{repo: repo, hub: hub}
}
