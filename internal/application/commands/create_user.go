package commands

import (
	"context"
	"errors"

	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
	"github.com/GitAlex9/go-microservice-order/internal/application/mapper"
	"github.com/GitAlex9/go-microservice-order/internal/application/validation"
	"github.com/GitAlex9/go-microservice-order/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-microservice-order/internal/domain/errors"
	"github.com/GitAlex9/go-microservice-order/internal/domain/repositories"
)

type CreateUserHandler struct {
	repo repositories.UserRepository
}

func NewCreateUserHandler(repo repositories.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

func (h *CreateUserHandler) Handle(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	email, role, verr := validation.ValidateCreateUser(req)
	if verr.HasErrors() {
		return nil, verr
	}

	existing, err := h.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domainerrors.ErrDuplicateEmail
	}

	user, err := entities.NewUser(email, req.Password, role)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	response := mapper.UserToResponse(user)
	return &response, nil
}
