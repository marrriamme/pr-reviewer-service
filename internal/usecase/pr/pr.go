package pr

import (
	"context"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
	"github.com/marrria_mme/pr-reviewer-service/internal/repository"
)

type PRUsecase struct {
	prRepo   repository.IPRRepository
	userRepo repository.IUserRepository
}

func NewPRUsecase(prRepo repository.IPRRepository, userRepo repository.IUserRepository) *PRUsecase {
	return &PRUsecase{
		prRepo:   prRepo,
		userRepo: userRepo,
	}
}

func (u *PRUsecase) CreatePR(ctx context.Context, pr *models.PullRequestWithReviewers) (*models.PullRequestWithReviewers, error) {
	exists, err := u.prRepo.PRExists(ctx, pr.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to check PR existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("PR already exists: %w", errs.ErrPRExists)
	}

	author, err := u.userRepo.GetUser(ctx, pr.AuthorID)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, fmt.Errorf("author not found: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	if !author.IsActive {
		return nil, fmt.Errorf("author is not active: %w", errs.ErrUserNotActive)
	}

	if author.TeamName == "" {
		return nil, fmt.Errorf("author does not belong to any team: %w", errs.ErrUserNoTeam)
	}

	reviewers, err := u.userRepo.GetRandomActiveTeamMembers(ctx, author.TeamName, pr.AuthorID, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to get random reviewers: %w", err)
	}

	pr.AssignedReviewers = reviewers

	createdPR, err := u.prRepo.CreatePR(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	return createdPR, nil
}

func (u *PRUsecase) MergePR(ctx context.Context, prID string) (*models.PullRequestWithReviewers, error) {
	pr, err := u.prRepo.GetPR(ctx, prID)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, fmt.Errorf("PR not found: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	mergedPR, err := u.prRepo.MergePR(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	return mergedPR, nil
}

func (u *PRUsecase) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequestWithReviewers, string, error) {
	pr, err := u.prRepo.GetPR(ctx, prID)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, "", fmt.Errorf("PR not found: %w", errs.ErrNotFound)
		}
		return nil, "", fmt.Errorf("failed to get PR: %w", err)
	}

	if pr.Status == "MERGED" {
		return nil, "", fmt.Errorf("PR is already merged: %w", errs.ErrPRMerged)
	}

	found := false
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", fmt.Errorf("user not assigned as reviewer: %w", errs.ErrNotAssigned)
	}

	oldUser, err := u.userRepo.GetUser(ctx, oldUserID)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, "", fmt.Errorf("old reviewer not found: %w", errs.ErrNotFound)
		}
		return nil, "", fmt.Errorf("failed to get old reviewer: %w", err)
	}

	newUserID, err := u.userRepo.GetRandomActiveTeamMember(ctx, oldUser.TeamName, oldUserID, pr.AuthorID)
	if err != nil {
		if err == errs.ErrNoCandidate {
			return nil, "", fmt.Errorf("no available reviewers found: %w", errs.ErrNoCandidate)
		}
		return nil, "", fmt.Errorf("failed to get new reviewer: %w", err)
	}

	newReviewers := make([]string, len(pr.AssignedReviewers))
	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUserID {
			newReviewers[i] = newUserID
		} else {
			newReviewers[i] = reviewer
		}
	}

	updatedPR, err := u.prRepo.UpdatePRReviewers(ctx, prID, newReviewers)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update PR reviewers: %w", err)
	}

	return updatedPR, newUserID, nil
}
