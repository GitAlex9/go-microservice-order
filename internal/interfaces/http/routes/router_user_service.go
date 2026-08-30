package routes

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/GitAlex9/go-microservice-order/internal/application/factory"
	"github.com/GitAlex9/go-microservice-order/internal/interfaces/http/handlers"
	appmiddleware "github.com/GitAlex9/go-microservice-order/internal/interfaces/http/middleware"
	"github.com/GitAlex9/go-microservice-order/internal/pkg/jwt"
	"github.com/GitAlex9/go-microservice-order/internal/pkg/logger"
)

func NewUserServiceRouter(services *factory.ServiceFactory, tokenManager *jwt.TokenManager, log logger.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(appmiddleware.Logging(log))
	r.Use(appmiddleware.Recovery(log))

	authHandler := handlers.NewAuthHandler(services.AuthService)
	userHandler := handlers.NewUserHandler(services.UserService)

	r.Route("/api/v1", func(r chi.Router) {
		RegisterAuthRoutes(r, authHandler)
		RegisterUserRoutes(r, userHandler, tokenManager)
	})
	return r
}
