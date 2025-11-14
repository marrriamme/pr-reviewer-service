package usecase

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IPRUsecase interface {
	CreatePR(ctx context.Context, pr *models.PullRequestWithReviewers) (*models.PullRequestWithReviewers, error)
	MergePR(ctx context.Context, prID string) (*models.PullRequestWithReviewers, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequestWithReviewers, string, error)
	GetUserReviewPRs(ctx context.Context, userID string) ([]models.PullRequestShort, error)
}
