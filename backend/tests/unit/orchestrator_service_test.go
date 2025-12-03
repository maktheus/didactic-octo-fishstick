package unit

import (
	"context"
	"testing"

	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/storage"
	orchestratorrepository "github.com/example/back-end-tcc/services/orchestrator/repository"
	orchestratorservice "github.com/example/back-end-tcc/services/orchestrator/service"
)

func TestOrchestratorSubmitPublishesMessage(t *testing.T) {
	bus := queue.NewBus()
	repo := orchestratorrepository.New(storage.NewMemoryRepository[models.Submission]())
	service := orchestratorservice.New(repo, bus)

	called := false
	bus.Subscribe(orchestratorservice.SubmissionCreated, func(ctx context.Context, msg queue.Message) error {
		called = true
		return nil
	})

	_, err := service.Submit(context.Background(), "benchmark", "agent", "payload")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected submission.created message to be published")
	}
}
