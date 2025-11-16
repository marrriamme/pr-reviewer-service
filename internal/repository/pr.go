package repository

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IPRRepository interface {
	CreatePR(context.Context, *models.PullRequestWithReviewers) (*models.PullRequestWithReviewers, error)
	GetPR(context.Context, string) (*models.PullRequestWithReviewers, error)
	UpdatePRReviewers(context.Context, string, []string) (*models.PullRequestWithReviewers, error)
	MergePR(context.Context, string) (*models.PullRequestWithReviewers, error)
	GetUserReviewPRs(context.Context, string) ([]models.PullRequestShort, error)
	PRExists(context.Context, string) (bool, error)
	IsPRReviewer(context.Context, string, string) (bool, error)
}
