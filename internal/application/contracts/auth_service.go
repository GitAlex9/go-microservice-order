package contracts

import (
	"context"

	"github.com/GitAlex9/go-microservice-order/internal/application/dto"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}
