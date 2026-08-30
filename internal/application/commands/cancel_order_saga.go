package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-microservice-order/internal/application/contracts"
	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
	"github.com/GitAlex9/go-microservice-order/internal/application/mapper"
	"github.com/GitAlex9/go-microservice-order/internal/domain/entities"
	domainevents "github.com/GitAlex9/go-microservice-order/internal/domain/events"
	"github.com/GitAlex9/go-microservice-order/internal/domain/repositories"
	"github.com/GitAlex9/go-microservice-order/internal/integration/httpclient"
)

type CancelOrderSagaHandler struct {
	orderRepo     repositories.OrderRepository
	uow           contracts.UnitOfWork
	dispatcher    domainevents.Dispatcher
	productClient *httpclient.ProductClient
}

func NewCancelOrderSagaHandler(
	orderRepo repositories.OrderRepository,
	uow contracts.UnitOfWork,
	dispatcher domainevents.Dispatcher,
	productClient *httpclient.ProductClient,
) *CancelOrderSagaHandler {
	return &CancelOrderSagaHandler{orderRepo: orderRepo, uow: uow, dispatcher: dispatcher, productClient: productClient}
}

func (h *CancelOrderSagaHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	order, err := h.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	for _, item := range order.Items() {
		if err := h.productClient.IncreaseStock(item.ProductID().String(), item.Quantity()); err != nil {
			return nil, err
		}
	}

	if err := h.uow.Execute(ctx, func(repos contracts.Repositories) error {
		return repos.Order.Save(ctx, order)
	}); err != nil {
		return nil, err
	}

	h.dispatcher.Dispatch(ctx, order.Events())
	order.ClearEvents()

	response := mapper.OrderToResponse(order)
	return &response, nil
}

var _ = entities.OrderStatusPending
