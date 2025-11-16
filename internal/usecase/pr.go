package usecase

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IPRUsecase interface {
	CreatePR(context.Context, *models.PullRequestWithReviewers) (*models.PullRequestWithReviewers, error)
	MergePR(context.Context, string) (*models.PullRequestWithReviewers, error)
	ReassignReviewer(context.Context, string, string) (*models.PullRequestWithReviewers, string, error)
}
