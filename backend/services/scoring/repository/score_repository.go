package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// ScoreRepository stores score summaries.
type ScoreRepository struct {
	store storage.Repository[models.ScoreSummary]
}

// New creates repository.
func New(store storage.Repository[models.ScoreSummary]) *ScoreRepository {
	return &ScoreRepository{store: store}
}

// Save stores summary.
func (r *ScoreRepository) Save(id string, summary models.ScoreSummary) {
	r.store.Save(id, summary)
}

// List returns all summaries.
func (r *ScoreRepository) List() []models.ScoreSummary {
	return r.store.List()
}
