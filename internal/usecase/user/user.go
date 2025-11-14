package user

import (
	"context"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
	"github.com/marrria_mme/pr-reviewer-service/internal/repository"
)

type UserUsecase struct {
	repo   repository.IUserRepository
	prRepo repository.IPRRepository
}

func NewUserUsecase(repo repository.IUserRepository, prRepo repository.IPRRepository) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		prRepo: prRepo,
	}
}

func (u *UserUsecase) SetUserActivity(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	user, err := u.repo.UpdateUserActivity(ctx, userID, isActive)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, fmt.Errorf("user not found: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to update user activity: %w", err)
	}

	return user, nil
}

func (u *UserUsecase) GetUserReviewPRs(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	exists, err := u.repo.UserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("user not found: %w", errs.ErrNotFound)
	}

	return u.prRepo.GetUserReviewPRs(ctx, userID)
}

func (u *UserUsecase) GetRandomActiveTeamMember(ctx context.Context, teamName, excludeUserID string) (string, error) {
	userID, err := u.repo.GetRandomActiveTeamMember(ctx, teamName, excludeUserID)
	if err != nil {
		if err == errs.ErrNoCandidate {
			return "", fmt.Errorf("no active team members available: %w", errs.ErrNoCandidate)
		}
		return "", fmt.Errorf("failed to get random team member: %w", err)
	}

	return userID, nil
}

func (u *UserUsecase) GetActiveTeamMembers(ctx context.Context, teamName, excludeUserID string) ([]models.User, error) {
	users, err := u.repo.GetActiveTeamMembers(ctx, teamName, excludeUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active team members: %w", err)
	}

	return users, nil
}
