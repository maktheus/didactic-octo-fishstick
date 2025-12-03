package integration

import (
	"context"
	"testing"
	"time"

	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/storage"
	agentrepository "github.com/example/back-end-tcc/services/agent/repository"
	benchmarkrepository "github.com/example/back-end-tcc/services/benchmark/repository"
	orchestratorrepository "github.com/example/back-end-tcc/services/orchestrator/repository"
	orchestratorservice "github.com/example/back-end-tcc/services/orchestrator/service"
	runnerrepository "github.com/example/back-end-tcc/services/runner/repository"
	runnerservice "github.com/example/back-end-tcc/services/runner/service"
	scoringrepository "github.com/example/back-end-tcc/services/scoring/repository"
	scoringservice "github.com/example/back-end-tcc/services/scoring/service"
)

func TestEndToEndBenchmarkFlow(t *testing.T) {
	bus := queue.NewBus()

	submissionStore := storage.NewMemoryRepository[models.Submission]()
	orchestratorRepo := orchestratorrepository.New(submissionStore)
	orchestratorSvc := orchestratorservice.New(orchestratorRepo, bus)

	runnerRepo := runnerrepository.New(storage.NewMemoryRepository[models.Submission]())
	agentRepo := agentrepository.NewAgentRepository(storage.NewMemoryRepository[models.User]())
	agentRepo.Save(models.User{ID: "agent", Name: "Test Agent"})

	benchmarkRepo := benchmarkrepository.New(storage.NewMemoryRepository[models.Benchmark]())
	benchmarkRepo.Save(models.Benchmark{ID: "bench", Name: "Test Benchmark", Tasks: []models.Task{{Prompt: "Test Prompt"}}})

	runnerSvc := runnerservice.New(runnerRepo, agentRepo, benchmarkRepo, bus, bus)
	runnerSvc.Start()

	scoringRepo := scoringrepository.New(storage.NewMemoryRepository[models.ScoreSummary]())
	scoringSvc := scoringservice.New(scoringRepo, bus, bus)
	scoringSvc.Start()

	if _, err := orchestratorSvc.Submit(context.Background(), "bench", "agent", "payload"); err != nil {
		t.Fatalf("submit error: %v", err)
	}

	// allow asynchronous handlers to run synchronously by reusing same goroutine
	// (handlers execute inline in the in-memory bus).
	if len(runnerSvc.Results()) == 0 {
		t.Fatal("expected runner results to be available")
	}

	eventually := time.After(10 * time.Millisecond)
	<-eventually

	summaries := scoringSvc.Summaries()
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
}
