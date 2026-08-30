package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/GitAlex9/go-microservice-order/internal/application/commands"
	appevents "github.com/GitAlex9/go-microservice-order/internal/application/events"
	"github.com/GitAlex9/go-microservice-order/internal/application/queries"
	"github.com/GitAlex9/go-microservice-order/internal/application/services"
	"github.com/GitAlex9/go-microservice-order/internal/infrastructure/database/postgres"
	repository "github.com/GitAlex9/go-microservice-order/internal/infrastructure/repositories/postgres"
	"github.com/GitAlex9/go-microservice-order/internal/integration/httpclient"
	"github.com/GitAlex9/go-microservice-order/internal/interfaces/http/routes"
	"github.com/GitAlex9/go-microservice-order/internal/pkg/jwt"
	"github.com/GitAlex9/go-microservice-order/internal/pkg/logger"
	"github.com/google/uuid"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ no .env file found, using system environment variables")
	}

	cfg := postgres.NewConfig()
	connection, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	migrator := postgres.NewMigrator(connection.Pool())
	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}

	appLogger := logger.New(os.Getenv("APP_ENV"))
	orderRepo := repository.NewOrderRepository(connection.Pool())
	uow := repository.NewUnitOfWork(connection.Pool(), appLogger)
	dispatcher := appevents.NewDefaultDispatcher(appLogger)

	tokenManager := jwt.NewTokenManager(getEnv("JWT_SECRET", "dev-secret-change-me"), 24*time.Hour)

	systemToken, err := tokenManager.Generate(uuid.Nil, "microservice-order@internal", "admin")
	if err != nil {
		log.Fatal("failed to generate internal service token:", err)
	}

	productClient := httpclient.NewProductClient(getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082"), systemToken)
	customerClient := httpclient.NewCustomerClient(getEnv("CUSTOMER_SERVICE_URL", "http://localhost:8081"), systemToken)

	orderService := services.NewOrderService(
		commands.NewCreateOrderSagaHandler(uow, dispatcher, productClient, customerClient),
		commands.NewPayOrderHandler(uow, dispatcher),
		commands.NewCancelOrderSagaHandler(orderRepo, uow, dispatcher, productClient),
		commands.NewDeleteOrderHandler(uow),
		queries.NewGetOrderHandler(orderRepo),
		queries.NewListOrdersHandler(orderRepo),
	)

	router := routes.NewOrderServiceRouter(orderService, tokenManager, appLogger)

	addr := ":8083"
	log.Println("✓ microservice-order (saga orchestrator) listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
