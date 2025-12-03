package storage

import "sync"

// Repository defines CRUD operations for simple key-based storage.
type Repository[T any] interface {
	Save(id string, value T)
	Get(id string) (T, bool)
	List() []T
}

// MemoryRepository is a thread-safe in-memory repository.
type MemoryRepository[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

// NewMemoryRepository creates a MemoryRepository.
func NewMemoryRepository[T any]() *MemoryRepository[T] {
	return &MemoryRepository[T]{items: make(map[string]T)}
}

// Save stores the value.
func (r *MemoryRepository[T]) Save(id string, value T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[id] = value
}

// Get retrieves an item by id.
func (r *MemoryRepository[T]) Get(id string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.items[id]
	return v, ok
}

// List returns all values.
func (r *MemoryRepository[T]) List() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	values := make([]T, 0, len(r.items))
	for _, v := range r.items {
		values = append(values, v)
	}
	return values
}
