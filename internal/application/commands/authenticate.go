package commands

import (
	"context"

	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
	"github.com/GitAlex9/go-microservice-order/internal/application/mapper"
	domainerrors "github.com/GitAlex9/go-microservice-order/internal/domain/errors"
	"github.com/GitAlex9/go-microservice-order/internal/domain/repositories"
	"github.com/GitAlex9/go-microservice-order/internal/domain/valueobjects"
	"github.com/GitAlex9/go-microservice-order/internal/pkg/jwt"
)

type AuthenticateHandler struct {
	repo         repositories.UserRepository
	tokenManager *jwt.TokenManager
}

func NewAuthenticateHandler(repo repositories.UserRepository, tokenManager *jwt.TokenManager) *AuthenticateHandler {
	return &AuthenticateHandler{repo: repo, tokenManager: tokenManager}
}

func (h *AuthenticateHandler) Handle(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	email, err := valueobjects.NewEmail(req.Email)
	if err != nil {
		return nil, domainerrors.ErrInvalidCredentials
	}

	user, err := h.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, domainerrors.ErrInvalidCredentials
	}

	if !user.Active() {
		return nil, domainerrors.ErrInvalidCredentials
	}

	if !user.CheckPassword(req.Password) {
		return nil, domainerrors.ErrInvalidCredentials
	}

	token, err := h.tokenManager.Generate(user.ID(), user.Email().String(), user.Role().String())
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  mapper.UserToResponse(user),
	}, nil
}
