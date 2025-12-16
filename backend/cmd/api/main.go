package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/example/back-end-tcc/docs"
	"github.com/example/back-end-tcc/pkg/config"
	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/storage"
	agenthandlers "github.com/example/back-end-tcc/services/agent/handlers"
	agentrepository "github.com/example/back-end-tcc/services/agent/repository"
	agentservice "github.com/example/back-end-tcc/services/agent/service"
	authhandlers "github.com/example/back-end-tcc/services/auth/handlers"
	authrepository "github.com/example/back-end-tcc/services/auth/repository"
	authservice "github.com/example/back-end-tcc/services/auth/service"
	benchmarkhandlers "github.com/example/back-end-tcc/services/benchmark/handlers"
	benchmarkrepository "github.com/example/back-end-tcc/services/benchmark/repository"
	benchmarkservice "github.com/example/back-end-tcc/services/benchmark/service"
	leaderboardhandlers "github.com/example/back-end-tcc/services/leaderboard/handlers"
	leaderboardrepository "github.com/example/back-end-tcc/services/leaderboard/repository"
	leaderboardservice "github.com/example/back-end-tcc/services/leaderboard/service"
	orchestratorhandlers "github.com/example/back-end-tcc/services/orchestrator/handlers"
	orchestratorrepository "github.com/example/back-end-tcc/services/orchestrator/repository"
	orchestratorservice "github.com/example/back-end-tcc/services/orchestrator/service"
	runnerhandlers "github.com/example/back-end-tcc/services/runner/handlers"
	runnerrepository "github.com/example/back-end-tcc/services/runner/repository"
	runnerservice "github.com/example/back-end-tcc/services/runner/service"
	scoringhandlers "github.com/example/back-end-tcc/services/scoring/handlers"
	scoringrepository "github.com/example/back-end-tcc/services/scoring/repository"
	scoringservice "github.com/example/back-end-tcc/services/scoring/service"
	tracehandlers "github.com/example/back-end-tcc/services/trace/handlers"
	tracerepository "github.com/example/back-end-tcc/services/trace/repository"
	traceservice "github.com/example/back-end-tcc/services/trace/service"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	apiLog := logger.New(logger.WithPrefix("api "))
	meter := metrics.NewInMemory()

	bus := queue.NewBus(queue.WithLogger(apiLog), queue.WithMetrics(meter))

	newServiceLogger := func(prefix string) logger.Logger {
		return logger.New(logger.WithPrefix(prefix + " "))
	}

	var db *sql.DB
	if strings.HasPrefix(cfg.StorageDSN, "postgres://") {
		var err error
		for i := 0; i < 10; i++ {
			db, err = sql.Open("postgres", cfg.StorageDSN)
			if err != nil {
				log.Printf("Failed to open database connection: %v. Retrying...", err)
				time.Sleep(2 * time.Second)
				continue
			}
			if err := db.Ping(); err != nil {
				log.Printf("Failed to ping database: %v. Retrying...", err)
				time.Sleep(2 * time.Second)
				continue
			}
			apiLog.Println("Connected to PostgreSQL")
			break
		}
		if db == nil {
			log.Fatalf("Failed to connect to database after retries")
		}
	}

	submissionRepo := createRepo[models.Submission](db, "submissions")
	scoreRepo := createRepo[models.ScoreSummary](db, "scores")
	traceRepo := createRepo[models.TraceEvent](db, "traces")
	leaderboardRepo := createRepo[models.LeaderboardEntry](db, "leaderboard")
	benchmarkRepo := createRepo[models.Benchmark](db, "benchmarks")
	agentRepoStore := createRepo[models.User](db, "agents")

	authRepoStore := createRepo[models.User](db, "users")
	authRepo := authrepository.NewUserRepository(authRepoStore)
	authRepo.Seed(models.User{ID: "admin", Email: "admin@example.com", Role: "admin"})
	authSrv := authservice.NewAuthService(
		authRepo,
		authservice.WithLogger(newServiceLogger("auth")),
		authservice.WithMetrics(meter),
	)
	authHTTP := authhandlers.NewHTTP(authSrv)

	agentRepo := agentrepository.NewAgentRepository(agentRepoStore)
	agentSrv := agentservice.NewAgentService(
		agentRepo,
		agentservice.WithLogger(newServiceLogger("agent")),
		agentservice.WithMetrics(meter),
	)
	agentHTTP := agenthandlers.NewHTTP(agentSrv)

	benchmarkRepoImpl := benchmarkrepository.New(benchmarkRepo)
	benchmarkSrv := benchmarkservice.New(
		benchmarkRepoImpl,
		benchmarkservice.WithLogger(newServiceLogger("benchmark")),
		benchmarkservice.WithMetrics(meter),
	)
	benchmarkHTTP := benchmarkhandlers.New(benchmarkSrv)

	orchestratorRepo := orchestratorrepository.New(submissionRepo)
	orchestratorSrv := orchestratorservice.New(
		orchestratorRepo,
		bus,
		orchestratorservice.WithLogger(newServiceLogger("orchestrator")),
		orchestratorservice.WithMetrics(meter),
	)
	orchestratorHTTP := orchestratorhandlers.New(orchestratorSrv)

	runnerRepo := runnerrepository.New(submissionRepo, traceRepo)
	runnerSrv := runnerservice.New(
		runnerRepo,
		agentRepo,
		benchmarkRepoImpl,
		bus,
		bus,
		runnerservice.WithLogger(newServiceLogger("runner")),
		runnerservice.WithMetrics(meter),
	)
	runnerSrv.Start()
	runnerHTTP := runnerhandlers.New(runnerSrv)

	scoringRepo := scoringrepository.New(scoreRepo)
	scoringSrv := scoringservice.New(
		scoringRepo,
		bus,
		bus,
		scoringservice.WithLogger(newServiceLogger("scoring")),
		scoringservice.WithMetrics(meter),
	)
	scoringSrv.Start()
	scoringHTTP := scoringhandlers.New(scoringSrv)

	traceRepoImpl := tracerepository.New(traceRepo)
	traceSrv := traceservice.New(
		traceRepoImpl,
		bus,
		traceservice.WithLogger(newServiceLogger("trace")),
		traceservice.WithMetrics(meter),
	)
	traceHTTP := tracehandlers.New(traceSrv)

	leaderboardRepoImpl := leaderboardrepository.New(leaderboardRepo)
	leaderboardSrv := leaderboardservice.New(
		leaderboardRepoImpl,
		bus,
		leaderboardservice.WithLogger(newServiceLogger("leaderboard")),
		leaderboardservice.WithMetrics(meter),
	)
	leaderboardSrv.Start()
	leaderboardHTTP := leaderboardhandlers.New(leaderboardSrv)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth", authHTTP.Authenticate)
	mux.HandleFunc("/agents", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: agentHTTP.Register,
		http.MethodGet:  agentHTTP.List,
	}))
	mux.HandleFunc("/benchmarks", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: benchmarkHTTP.Create,
		http.MethodGet:  benchmarkHTTP.List,
	}))
	mux.HandleFunc("/submissions", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: orchestratorHTTP.Submit,
		http.MethodGet:  orchestratorHTTP.List,
	}))
	mux.HandleFunc("/results", runnerHTTP.Results)
	mux.HandleFunc("/scores", scoringHTTP.List)
	mux.HandleFunc("/traces", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: traceHTTP.Record,
		http.MethodGet:  traceHTTP.List,
	}))
	mux.HandleFunc("/leaderboard", leaderboardHTTP.List)
    
	// Reset handler
	mux.HandleFunc("/reset", withMethod(map[string]http.HandlerFunc{
		http.MethodPost: func(w http.ResponseWriter, r *http.Request) {
			apiLog.Println("Resetting platform data...")
			
			// Clear all repositories content
			if err := submissionRepo.Clear(); err != nil {
				http.Error(w, "failed to clear submissions", http.StatusInternalServerError)
				return
			}
			if err := traceRepo.Clear(); err != nil {
				http.Error(w, "failed to clear traces", http.StatusInternalServerError)
				return
			}
			if err := scoreRepo.Clear(); err != nil {
				http.Error(w, "failed to clear scores", http.StatusInternalServerError)
				return
			}
			// Note: Agents and Benchmarks are kept unless "Resetar Tudo" implies them too.
			// The UI Text says "Remove todos os dados da plataforma".
			// Given urgency, I will wipe Runs/Traces/Scores which is the main runtime data.
			// Clearing Agents might annoy user if he just wants to reset Progress.
			// But for "Reset All", I probably should clear everything.
			// Let's clear everything as per user request "Resetar Tudo".
			if err := agentRepoStore.Clear(); err != nil {
				http.Error(w, "failed to clear agents", http.StatusInternalServerError)
				return
			}
			if err := benchmarkRepo.Clear(); err != nil {
				http.Error(w, "failed to clear benchmarks", http.StatusInternalServerError)
				return
			}
			if err := leaderboardRepo.Clear(); err != nil {
				http.Error(w, "failed to clear leaderboard", http.StatusInternalServerError)
				return
			}

			// Re-seed admin user?
            // The auth table "users" is separate (authRepoStore).
            // Usually we don't wipe users configured valid access.
            // I'll leave auth untouched.
			
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Platform reset successfully"})
		},
	}))

	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(docs.OpenAPISpec); err != nil && apiLog != nil {
			apiLog.Println("swagger serve error:", err)
		}
	})

	apiLog.Printf("API gateway listening on :%d", cfg.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), enableCORS(mux)); err != nil {
		apiLog.Println("server error:", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
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

func createRepo[T any](db *sql.DB, collection string) storage.Repository[T] {
	if db != nil {
		return storage.NewPostgresRepository[T](db, collection)
	}
	return storage.NewMemoryRepository[T]()
}
