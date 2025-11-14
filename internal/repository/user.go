package repository

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IUserRepository interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUserActivity(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetActiveTeamMembers(ctx context.Context, teamName, excludeUserID string) ([]models.User, error)
	GetRandomActiveTeamMember(ctx context.Context, teamName, excludeUserID string) (string, error)
	UserExists(ctx context.Context, userID string) (bool, error)
	GetRandomActiveTeamMembers(ctx context.Context, teamName, excludeUserID string, limit int) ([]string, error)
}
