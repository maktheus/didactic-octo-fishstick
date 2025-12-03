package main

import (
	"fmt"
	"net/http"

	"github.com/example/back-end-tcc/pkg/config"
	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/storage"
	agenthandlers "github.com/example/back-end-tcc/services/agent/handlers"
	agentrepository "github.com/example/back-end-tcc/services/agent/repository"
	agentservice "github.com/example/back-end-tcc/services/agent/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	log := logger.New(logger.WithPrefix("agent "))
	meter := metrics.NewInMemory()

	store := storage.NewMemoryRepository[models.User]()
	repo := agentrepository.NewAgentRepository(store)
	srv := agentservice.NewAgentService(
		repo,
		agentservice.WithLogger(log),
		agentservice.WithMetrics(meter),
	)
	handlers := agenthandlers.NewHTTP(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/agents", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: handlers.Register,
		http.MethodGet:  handlers.List,
	}))

	log.Printf("agent service listening on :%d", cfg.HTTPPort)
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
