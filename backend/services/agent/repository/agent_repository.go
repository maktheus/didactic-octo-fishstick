package repository

import (
	"github.com/example/back-end-tcc/pkg/models"
	"github.com/example/back-end-tcc/pkg/storage"
)

// AgentRepository stores agent definitions.
type AgentRepository struct {
	store storage.Repository[models.User]
}

// NewAgentRepository creates repo.
func NewAgentRepository(store storage.Repository[models.User]) *AgentRepository {
	return &AgentRepository{store: store}
}

// Save stores agent data.
func (r *AgentRepository) Save(agent models.User) {
	r.store.Save(agent.ID, agent)
}

// List returns agents.
func (r *AgentRepository) List() []models.User {
	return r.store.List()
}

// Get returns an agent by ID.
func (r *AgentRepository) Get(id string) (models.User, bool) {
	return r.store.Get(id)
}
