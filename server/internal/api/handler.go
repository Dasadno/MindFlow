// server/internal/api/handler.go
package api

import "milk/server/internal/storage"

// Handler — HTTP-обработчики с доступом к хранилищу.
type Handler struct {
	repo *storage.Repository
}

// NewHandler создаёт Handler с инъекцией зависимости Repository.
func NewHandler(repo *storage.Repository) *Handler {
	return &Handler{repo: repo}
}
