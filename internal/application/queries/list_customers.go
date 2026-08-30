package queries

import (
	"context"

	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
	"github.com/GitAlex9/go-microservice-order/internal/application/mapper"
	"github.com/GitAlex9/go-microservice-order/internal/domain/repositories"
)

type ListCustomersHandler struct {
	repo repositories.CustomerRepository
}

func NewListCustomersHandler(repo repositories.CustomerRepository) *ListCustomersHandler {
	return &ListCustomersHandler{repo: repo}
}

func (h *ListCustomersHandler) Handle(ctx context.Context, offset, limit int) ([]dto.CustomerResponse, error) {
	customers, err := h.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return mapper.CustomersToResponse(customers), nil
}
