package team

import (
	"context"

	"github.com/google/uuid"

	"timebride/internal/models"
	"timebride/internal/repositories"
	"timebride/internal/types"
)

type teamService struct {
	teamRepo repositories.Repository[models.TeamMember]
}

// NewTeamService creates a new team service instance
func NewTeamService(teamRepo repositories.Repository[models.TeamMember]) ITeamService {
	return &teamService{
		teamRepo: teamRepo,
	}
}

// CreateMember creates a new team member
func (s *teamService) CreateMember(ctx context.Context, userID uuid.UUID, member *models.TeamMember) error {
	member.UserID = userID
	return s.teamRepo.Create(ctx, member)
}

// UpdateMember оновлює дані члена команди
func (s *teamService) UpdateMember(ctx context.Context, member *models.TeamMember) error {
	return s.teamRepo.Update(ctx, member)
}

// DeleteMember видаляє члена команди
func (s *teamService) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	return s.teamRepo.Delete(ctx, memberID)
}

// GetMember отримує члена команди за ID
func (s *teamService) GetMember(ctx context.Context, memberID uuid.UUID) (*models.TeamMember, error) {
	return s.teamRepo.GetByID(ctx, memberID)
}

// ListMembers повертає список членів команди
func (s *teamService) ListMembers(ctx context.Context, userID uuid.UUID) ([]*models.TeamMember, error) {
	return s.teamRepo.List(ctx, map[string]interface{}{"user_id": userID})
}

// UpdateRole оновлює роль члена команди
func (s *teamService) UpdateRole(ctx context.Context, memberID uuid.UUID, role types.TeamRole) error {
	member, err := s.GetMember(ctx, memberID)
	if err != nil {
		return err
	}

	member.Role = string(role)
	return s.teamRepo.Update(ctx, member)
}
