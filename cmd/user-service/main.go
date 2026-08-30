package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/GitAlex9/go-microservice-order/internal/application/factory"
	"github.com/GitAlex9/go-microservice-order/internal/infrastructure/database/postgres"
	"github.com/GitAlex9/go-microservice-order/internal/interfaces/http/routes"
	"github.com/GitAlex9/go-microservice-order/internal/pkg/logger"
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
	services := factory.NewServiceFactory(connection.Pool(), appLogger)

	router := routes.NewUserServiceRouter(services, services.TokenManager, appLogger)

	addr := ":8084"
	log.Println("✓ customer-service listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
