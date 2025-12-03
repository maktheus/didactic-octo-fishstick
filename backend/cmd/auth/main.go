package main

import (
	"fmt"
	"net/http"

	"github.com/example/back-end-tcc/pkg/config"
	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/storage"
	authhandlers "github.com/example/back-end-tcc/services/auth/handlers"
	authrepository "github.com/example/back-end-tcc/services/auth/repository"
	authservice "github.com/example/back-end-tcc/services/auth/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	log := logger.New(logger.WithPrefix("auth "))
	meter := metrics.NewInMemory()

	repoStore := storage.NewMemoryRepository[models.User]()
	repo := authrepository.NewUserRepository(repoStore)
	repo.Seed(models.User{ID: "admin", Email: "admin@example.com", Role: "admin"})

	srv := authservice.NewAuthService(
		repo,
		authservice.WithLogger(log),
		authservice.WithMetrics(meter),
	)
	httpHandlers := authhandlers.NewHTTP(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth", httpHandlers.Authenticate)

	log.Printf("auth service listening on :%d", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), mux); err != nil {
		log.Println("server error:", err)
	}
}
