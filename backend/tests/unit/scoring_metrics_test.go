package unit

import (
	"context"
	"testing"
	"time"

	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/queue"
	"github.com/example/back-end-tcc/pkg/storage"
	scoringrepository "github.com/example/back-end-tcc/services/scoring/repository"
	scoringservice "github.com/example/back-end-tcc/services/scoring/service"
)

func TestScoringServiceStoresSummary(t *testing.T) {
	bus := queue.NewBus()
	repo := scoringrepository.New(storage.NewMemoryRepository[models.ScoreSummary]())
	service := scoringservice.New(repo, bus, bus)
	service.Start()

	submission := models.Submission{ID: "id", ScoreSummary: &models.ScoreSummary{Score: 0.5, Metrics: map[string]float64{"precision": 0.5}, Calculated: time.Now()}}
	if err := bus.Publish(context.Background(), queue.Message{Type: scoringservice.ScoreCalculated, Data: submission}); err != nil {
		t.Fatalf("publish error: %v", err)
	}

	summaries := service.Summaries()
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	if summaries[0].Score != 0.5 {
		t.Fatalf("expected score 0.5, got %f", summaries[0].Score)
	}
}
