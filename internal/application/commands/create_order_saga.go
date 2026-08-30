package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-microservice-order/internal/application/contracts"
	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
	"github.com/GitAlex9/go-microservice-order/internal/application/mapper"
	"github.com/GitAlex9/go-microservice-order/internal/application/validation"
	"github.com/GitAlex9/go-microservice-order/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-microservice-order/internal/domain/errors"
	domainevents "github.com/GitAlex9/go-microservice-order/internal/domain/events"
	"github.com/GitAlex9/go-microservice-order/internal/domain/valueobjects"
	"github.com/GitAlex9/go-microservice-order/internal/integration/httpclient"
)

type CreateOrderSagaHandler struct {
	uow            contracts.UnitOfWork
	dispatcher     domainevents.Dispatcher
	productClient  *httpclient.ProductClient
	customerClient *httpclient.CustomerClient
}

func NewCreateOrderSagaHandler(
	uow contracts.UnitOfWork,
	dispatcher domainevents.Dispatcher,
	productClient *httpclient.ProductClient,
	customerClient *httpclient.CustomerClient,
) *CreateOrderSagaHandler {
	return &CreateOrderSagaHandler{
		uow: uow, dispatcher: dispatcher,
		productClient: productClient, customerClient: customerClient,
	}
}

type reservation struct {
	productID string
	quantity  int
}

func (h *CreateOrderSagaHandler) Handle(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	customerID, verr := validation.ValidateCreateOrder(req)
	if verr.HasErrors() {
		return nil, verr
	}

	if _, err := h.customerClient.Get(customerID.String()); err != nil {
		return nil, fmt.Errorf("saga: customer validation failed: %w", err)
	}

	var reservations []reservation
	var items []entities.OrderItem

	compensate := func() {
		for _, r := range reservations {
			_ = h.productClient.IncreaseStock(r.productID, r.quantity)
		}
	}

	for _, itemReq := range req.Items {
		product, err := h.productClient.Get(itemReq.ProductID)
		if err != nil {
			compensate()
			return nil, fmt.Errorf("saga: product lookup failed: %w", err)
		}
		if !product.Active {
			compensate()
			return nil, domainerrors.ErrInactiveProduct
		}

		if err := h.productClient.DecreaseStock(itemReq.ProductID, itemReq.Quantity); err != nil {
			compensate()
			return nil, fmt.Errorf("saga: stock reservation failed: %w", err)
		}
		reservations = append(reservations, reservation{productID: itemReq.ProductID, quantity: itemReq.Quantity})

		productID, _ := uuid.Parse(itemReq.ProductID)
		price, err := valueobjects.NewMoneyFromFloat(product.Price)
		if err != nil {
			compensate()
			return nil, err
		}

		item, err := entities.NewOrderItem(productID, product.Name, price, itemReq.Quantity)
		if err != nil {
			compensate()
			return nil, err
		}
		items = append(items, *item)
	}

	order, err := entities.NewOrder(customerID, items)
	if err != nil {
		compensate()
		return nil, err
	}

	saveErr := h.uow.Execute(ctx, func(repos contracts.Repositories) error {
		return repos.Order.Save(ctx, order)
	})
	if saveErr != nil {
		compensate()
		return nil, saveErr
	}

	h.dispatcher.Dispatch(ctx, order.Events())
	order.ClearEvents()

	response := mapper.OrderToResponse(order)
	return &response, nil
}
