package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	if err := r.db.QueryRowContext(ctx, queryGetUser, userID).Scan(
		&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserActivity(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	var user models.User
	if err := r.db.QueryRowContext(ctx, queryUpdateUserActivity, isActive, userID).Scan(
		&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update user activity: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetActiveTeamMembers(ctx context.Context, teamName, excludeUserID string) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, queryGetActiveTeamMembers, teamName, excludeUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active team members: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetRandomActiveTeamMember(ctx context.Context, teamName, excludeUserID, excludeAuthorID string) (string, error) {
	var userID string
	if err := r.db.QueryRowContext(ctx, queryGetRandomActiveTeamMember, teamName, excludeUserID, excludeAuthorID).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return "", errs.ErrNoCandidate
		}
		return "", fmt.Errorf("failed to get random team member: %w", err)
	}

	return userID, nil
}

func (r *UserRepository) GetRandomActiveTeamMembers(ctx context.Context, teamName, excludeAuthorID string, limit int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, queryGetRandomActiveTeamMembers, teamName, excludeAuthorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get random active team members: %w", err)
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err = rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return userIDs, nil
}

func (r *UserRepository) UserExists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, queryUserExists, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}
