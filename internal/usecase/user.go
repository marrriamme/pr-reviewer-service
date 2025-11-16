package usecase

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IUserUsecase interface {
	SetUserActivity(context.Context, string, bool) (*models.User, error)
	GetUserReviewPRs(context.Context, string) ([]models.PullRequestShort, error)
}
