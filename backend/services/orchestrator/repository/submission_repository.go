package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// SubmissionRepository stores submission state.
type SubmissionRepository struct {
	store storage.Repository[models.Submission]
}

// New creates repository.
func New(store storage.Repository[models.Submission]) *SubmissionRepository {
	return &SubmissionRepository{store: store}
}

// Save updates submission.
func (r *SubmissionRepository) Save(sub models.Submission) {
	r.store.Save(sub.ID, sub)
}

// List returns submissions.
func (r *SubmissionRepository) List() []models.Submission {
	return r.store.List()
}
