package usecase

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IUserUsecase interface {
	SetUserActivity(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserReviewPRs(ctx context.Context, userID string) ([]models.PullRequestShort, error)
	GetRandomActiveTeamMember(ctx context.Context, teamName, excludeUserID string) (string, error)
	GetActiveTeamMembers(ctx context.Context, teamName, excludeUserID string) ([]models.User, error)
}
