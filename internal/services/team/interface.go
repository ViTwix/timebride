package team

import (
	"context"

	"timebride/internal/models"
	"timebride/internal/types"

	"github.com/google/uuid"
)

// ITeamService визначає інтерфейс сервісу команди
type ITeamService interface {
	// CreateMember створює нового члена команди
	CreateMember(ctx context.Context, userID uuid.UUID, member *models.TeamMember) error

	// UpdateMember оновлює дані члена команди
	UpdateMember(ctx context.Context, member *models.TeamMember) error

	// DeleteMember видаляє члена команди
	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	// GetMember отримує члена команди за ID
	GetMember(ctx context.Context, memberID uuid.UUID) (*models.TeamMember, error)

	// ListMembers повертає список членів команди
	ListMembers(ctx context.Context, userID uuid.UUID) ([]*models.TeamMember, error)

	// UpdateRole оновлює роль члена команди
	UpdateRole(ctx context.Context, memberID uuid.UUID, role types.TeamRole) error
}
