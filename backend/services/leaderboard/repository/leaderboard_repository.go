package repository

import (
	"sort"

	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// Repository stores leaderboard entries.
type Repository struct {
	store storage.Repository[models.LeaderboardEntry]
}

// New creates repository.
func New(store storage.Repository[models.LeaderboardEntry]) *Repository {
	return &Repository{store: store}
}

// Save stores entry.
func (r *Repository) Save(entry models.LeaderboardEntry) {
	r.store.Save(entry.SubmissionID, entry)
}

// List returns entries sorted by rank.
func (r *Repository) List() []models.LeaderboardEntry {
	entries := r.store.List()
	sort.Slice(entries, func(i, j int) bool { return entries[i].Rank < entries[j].Rank })
	return entries
}
