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
	tracehandlers "github.com/example/back-end-tcc/services/trace/handlers"
	tracerepository "github.com/example/back-end-tcc/services/trace/repository"
	traceservice "github.com/example/back-end-tcc/services/trace/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	log := logger.New(logger.WithPrefix("trace "))
	meter := metrics.NewInMemory()

	bus := queue.NewBus(queue.WithLogger(log), queue.WithMetrics(meter))
	repo := tracerepository.New(storage.NewMemoryRepository[models.TraceEvent]())
	srv := traceservice.New(
		repo,
		bus,
		traceservice.WithLogger(log),
		traceservice.WithMetrics(meter),
	)
	handlers := tracehandlers.New(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/traces", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: handlers.Record,
		http.MethodGet:  handlers.List,
	}))

	log.Printf("trace service listening on :%d", cfg.HTTPPort)
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
