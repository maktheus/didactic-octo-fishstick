package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// UserRepository provides access to user records.
type UserRepository struct {
	store storage.Repository[models.User]
}

// NewUserRepository constructs a new repository using the provided storage.
func NewUserRepository(store storage.Repository[models.User]) *UserRepository {
	return &UserRepository{store: store}
}

// Seed adds a user to the repository.
func (r *UserRepository) Seed(user models.User) {
	r.store.Save(user.ID, user)
}

// FindByID returns a user with the given id.
func (r *UserRepository) FindByID(id string) (models.User, bool) {
	return r.store.Get(id)
}
