package queries

import (
	"context"

	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
	"github.com/GitAlex9/go-microservice-order/internal/application/mapper"
	"github.com/GitAlex9/go-microservice-order/internal/domain/repositories"
)

type ListProductsHandler struct {
	repo repositories.ProductRepository
}

func NewListProductsHandler(repo repositories.ProductRepository) *ListProductsHandler {
	return &ListProductsHandler{repo: repo}
}

func (h *ListProductsHandler) Handle(ctx context.Context, offset, limit int) ([]dto.ProductResponse, error) {
	products, err := h.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return mapper.ProductsToResponse(products), nil
}
