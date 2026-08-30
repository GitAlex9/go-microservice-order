package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-microservice-order/internal/domain/repositories"
)

type DeleteProductHandler struct {
	repo repositories.ProductRepository
}

func NewDeleteProductHandler(repo repositories.ProductRepository) *DeleteProductHandler {
	return &DeleteProductHandler{repo: repo}
}

func (h *DeleteProductHandler) Handle(ctx context.Context, id uuid.UUID) error {
	return h.repo.Delete(ctx, id)
}
