package main

import (
	"fmt"
	"net/http"

	"github.com/example/back-end-tcc/pkg/config"
	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/storage"
	benchmarkhandlers "github.com/example/back-end-tcc/services/benchmark/handlers"
	benchmarkrepository "github.com/example/back-end-tcc/services/benchmark/repository"
	benchmarkservice "github.com/example/back-end-tcc/services/benchmark/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	log := logger.New(logger.WithPrefix("benchmark "))
	meter := metrics.NewInMemory()

	store := storage.NewMemoryRepository[models.Benchmark]()
	repo := benchmarkrepository.New(store)
	srv := benchmarkservice.New(
		repo,
		benchmarkservice.WithLogger(log),
		benchmarkservice.WithMetrics(meter),
	)
	handlers := benchmarkhandlers.New(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/benchmarks", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: handlers.Create,
		http.MethodGet:  handlers.List,
	}))

	log.Printf("benchmark service listening on :%d", cfg.HTTPPort)
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
