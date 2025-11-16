package repository

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IUserRepository interface {
	GetUser(context.Context, string) (*models.User, error)
	UpdateUserActivity(context.Context, string, bool) (*models.User, error)
	GetRandomActiveTeamMember(context.Context, string, string, string) (string, error)
	GetRandomActiveTeamMembers(context.Context, string, string, int) ([]string, error)
	UserExists(context.Context, string) (bool, error)
	GetUserReviewPRs(context.Context, string) ([]models.PullRequestShort, error)
}
