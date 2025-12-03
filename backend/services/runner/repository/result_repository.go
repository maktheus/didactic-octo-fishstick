package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// ResultRepository stores submission results.
type ResultRepository struct {
	store      storage.Repository[models.Submission]
	traceStore storage.Repository[models.TraceEvent]
}

// New creates repository.
func New(store storage.Repository[models.Submission], traceStore storage.Repository[models.TraceEvent]) *ResultRepository {
	return &ResultRepository{store: store, traceStore: traceStore}
}

// Save stores submission.
func (r *ResultRepository) Save(sub models.Submission) {
	r.store.Save(sub.ID, sub)
}

// SaveTrace stores a trace event.
func (r *ResultRepository) SaveTrace(trace models.TraceEvent) {
	if r.traceStore != nil {
		r.traceStore.Save(trace.ID, trace)
	}
}

// List returns submissions.
func (r *ResultRepository) List() []models.Submission {
	return r.store.List()
}
