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
	scoringhandlers "github.com/example/back-end-tcc/services/scoring/handlers"
	scoringrepository "github.com/example/back-end-tcc/services/scoring/repository"
	scoringservice "github.com/example/back-end-tcc/services/scoring/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	log := logger.New(logger.WithPrefix("scoring "))
	meter := metrics.NewInMemory()

	bus := queue.NewBus(queue.WithLogger(log), queue.WithMetrics(meter))
	repo := scoringrepository.New(storage.NewMemoryRepository[models.ScoreSummary]())
	srv := scoringservice.New(
		repo,
		bus,
		bus,
		scoringservice.WithLogger(log),
		scoringservice.WithMetrics(meter),
	)
	srv.Start()
	handlers := scoringhandlers.New(srv)

	// Seed with a completed submission event to demonstrate scoring.
	seed := models.Submission{ID: "seed", ScoreSummary: &models.ScoreSummary{Score: 1.0, Metrics: map[string]float64{"latency": 50}, Calculated: time.Now()}}
	_ = bus.Publish(context.Background(), queue.Message{Type: scoringservice.ScoreCalculated, Data: seed})

	mux := http.NewServeMux()
	mux.HandleFunc("/scores", handlers.List)

	log.Printf("scoring service listening on :%d", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), mux); err != nil {
		log.Println("server error:", err)
	}
}
