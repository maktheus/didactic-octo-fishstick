package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/example/back-end-tcc/pkg/config"
	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/storage"
	agentrepository "github.com/example/back-end-tcc/services/agent/repository"
	benchmarkrepository "github.com/example/back-end-tcc/services/benchmark/repository"
	runnerhandlers "github.com/example/back-end-tcc/services/runner/handlers"
	runnerrepository "github.com/example/back-end-tcc/services/runner/repository"
	runnerservice "github.com/example/back-end-tcc/services/runner/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	log := logger.New(logger.WithPrefix("runner "))
	meter := metrics.NewInMemory()

	bus := queue.NewBus(queue.WithLogger(log), queue.WithMetrics(meter))
	repo := runnerrepository.New(storage.NewMemoryRepository[models.Submission]())
	agentRepo := agentrepository.NewAgentRepository(storage.NewMemoryRepository[models.User]())
	benchmarkRepo := benchmarkrepository.New(storage.NewMemoryRepository[models.Benchmark]())
	srv := runnerservice.New(
		repo,
		agentRepo,
		benchmarkRepo,
		bus,
		bus,
		runnerservice.WithLogger(log),
		runnerservice.WithMetrics(meter),
	)
	srv.Start()
	handlers := runnerhandlers.New(srv)

	// Example seeding of submission processing for demonstration.
	_ = bus.Publish(context.Background(), queue.Message{Type: "submission.created", Data: models.Submission{ID: "sample", AgentID: "agent", BenchmarkID: "bench"}})

	mux := http.NewServeMux()
	mux.HandleFunc("/results", handlers.Results)

	log.Printf("runner service listening on :%d", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), mux); err != nil {
		log.Println("server error:", err)
	}
}
