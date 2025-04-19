package repositories

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the basic CRUD operations for all repositories
type Repository[T any] interface {
	// Create adds a new entity to the database
	Create(ctx context.Context, entity *T) error

	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)

	// Update modifies an existing entity
	Update(ctx context.Context, entity *T) error

	// Delete removes an entity by its ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a list of entities with optional filtering
	List(ctx context.Context, filter map[string]interface{}) ([]*T, error)
} 