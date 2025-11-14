package repository

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IPRRepository interface {
	CreatePR(ctx context.Context, pr *models.PullRequestWithReviewers) (*models.PullRequestWithReviewers, error)
	GetPR(ctx context.Context, prID string) (*models.PullRequestWithReviewers, error)
	UpdatePRReviewers(ctx context.Context, prID string, reviewers []string) (*models.PullRequestWithReviewers, error)
	MergePR(ctx context.Context, prID string) (*models.PullRequestWithReviewers, error)
	GetUserReviewPRs(ctx context.Context, userID string) ([]models.PullRequestShort, error)
	PRExists(ctx context.Context, prID string) (bool, error)
	IsPRReviewer(ctx context.Context, prID, userID string) (bool, error)
}
