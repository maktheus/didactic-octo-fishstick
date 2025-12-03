package main

import (
	"fmt"
	"net/http"

	"github.com/example/back-end-tcc/pkg/config"
	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/storage"
	orchestratorhandlers "github.com/example/back-end-tcc/services/orchestrator/handlers"
	orchestratorrepository "github.com/example/back-end-tcc/services/orchestrator/repository"
	orchestratorservice "github.com/example/back-end-tcc/services/orchestrator/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	log := logger.New(logger.WithPrefix("orchestrator "))
	meter := metrics.NewInMemory()

	bus := queue.NewBus(queue.WithLogger(log), queue.WithMetrics(meter))

	store := storage.NewMemoryRepository[models.Submission]()
	repo := orchestratorrepository.New(store)
	srv := orchestratorservice.New(
		repo,
		bus,
		orchestratorservice.WithLogger(log),
		orchestratorservice.WithMetrics(meter),
	)
	handlers := orchestratorhandlers.New(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/submissions", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: handlers.Submit,
		http.MethodGet:  handlers.List,
	}))

	log.Printf("orchestrator service listening on :%d", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), mux); err != nil {
		log.Println("server error:", err)
	}
}

func withMethod(handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := handlers[r.Method]; ok {
			h(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
