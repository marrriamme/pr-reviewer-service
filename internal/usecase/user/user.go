package user

import (
	"context"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
	"github.com/marrria_mme/pr-reviewer-service/internal/repository"
)

type UserUsecase struct {
	userRepo repository.IUserRepository
}

func NewUserUsecase(userRepo repository.IUserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) SetUserActivity(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	user, err := u.userRepo.UpdateUserActivity(ctx, userID, isActive)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, fmt.Errorf("user not found: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to update user activity: %w", err)
	}

	return user, nil
}

func (u *UserUsecase) GetUserReviewPRs(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	exists, err := u.userRepo.UserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("user not found: %w", errs.ErrNotFound)
	}

	return u.userRepo.GetUserReviewPRs(ctx, userID)
}
