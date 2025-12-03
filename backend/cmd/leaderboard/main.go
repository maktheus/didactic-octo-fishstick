package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/example/back-end-tcc/pkg/config"
	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/storage"
	leaderboardhandlers "github.com/example/back-end-tcc/services/leaderboard/handlers"
	leaderboardrepository "github.com/example/back-end-tcc/services/leaderboard/repository"
	leaderboardservice "github.com/example/back-end-tcc/services/leaderboard/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	log := logger.New(logger.WithPrefix("leaderboard "))
	meter := metrics.NewInMemory()

	bus := queue.NewBus(queue.WithLogger(log), queue.WithMetrics(meter))
	repo := leaderboardrepository.New(storage.NewMemoryRepository[models.LeaderboardEntry]())
	srv := leaderboardservice.New(
		repo,
		bus,
		leaderboardservice.WithLogger(log),
		leaderboardservice.WithMetrics(meter),
	)
	srv.Start()
	handlers := leaderboardhandlers.New(srv)

	// Seed leaderboard event for demonstration.
	summary := models.ScoreSummary{Score: 1.0, Calculated: time.Now()}
	_ = bus.Publish(context.Background(), queue.Message{Type: "leaderboard.updated", Data: summary})

	mux := http.NewServeMux()
	mux.HandleFunc("/leaderboard", handlers.List)

	log.Printf("leaderboard service listening on :%d", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), mux); err != nil {
		log.Println("server error:", err)
	}
}
