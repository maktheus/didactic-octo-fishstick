package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// BenchmarkRepository stores benchmark definitions.
type BenchmarkRepository struct {
	store storage.Repository[models.Benchmark]
}

// New creates repo.
func New(store storage.Repository[models.Benchmark]) *BenchmarkRepository {
	return &BenchmarkRepository{store: store}
}

// Save persists a benchmark.
func (r *BenchmarkRepository) Save(benchmark models.Benchmark) {
	r.store.Save(benchmark.ID, benchmark)
}

// List returns benchmarks.
func (r *BenchmarkRepository) List() []models.Benchmark {
	return r.store.List()
}

// Get returns a benchmark by ID.
func (r *BenchmarkRepository) Get(id string) (models.Benchmark, bool) {
	return r.store.Get(id)
}
