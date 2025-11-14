package pr

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
)

type PRRepository struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) *PRRepository {
	return &PRRepository{db: db}
}

func (r *PRRepository) CreatePR(ctx context.Context, pr *models.PullRequestWithReviewers) (*models.PullRequestWithReviewers, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var createdPR models.PullRequest
	if err = tx.QueryRowContext(ctx, queryCreatePR,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status).Scan(
		&createdPR.PullRequestID, &createdPR.PullRequestName, &createdPR.AuthorID, &createdPR.Status,
		&createdPR.CreatedAt, &createdPR.MergedAt); err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	if len(pr.AssignedReviewers) > 0 {
		for _, reviewer := range pr.AssignedReviewers {
			if _, err = tx.ExecContext(ctx, queryAddReviewer, pr.PullRequestID, reviewer); err != nil {
				return nil, fmt.Errorf("failed to add reviewer %s: %w", reviewer, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetPR(ctx, pr.PullRequestID)
}

func (r *PRRepository) GetPR(ctx context.Context, prID string) (*models.PullRequestWithReviewers, error) {
	var prDB models.PullRequest
	if err := r.db.QueryRowContext(ctx, queryGetPR, prID).Scan(
		&prDB.PullRequestID, &prDB.PullRequestName, &prDB.AuthorID, &prDB.Status,
		&prDB.CreatedAt, &prDB.MergedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	reviewers, err := r.getPRReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	return &models.PullRequestWithReviewers{
		PullRequest:       prDB,
		AssignedReviewers: reviewers,
	}, nil
}

func (r *PRRepository) UpdatePRReviewers(ctx context.Context, prID string, reviewers []string) (*models.PullRequestWithReviewers, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, queryDeletePRReviewers, prID); err != nil {
		return nil, fmt.Errorf("failed to delete old reviewers: %w", err)
	}

	for _, reviewer := range reviewers {
		if _, err = tx.ExecContext(ctx, queryAddReviewer, prID, reviewer); err != nil {
			return nil, fmt.Errorf("failed to add reviewer %s: %w", reviewer, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetPR(ctx, prID)
}

func (r *PRRepository) MergePR(ctx context.Context, prID string) (*models.PullRequestWithReviewers, error) {
	var prDB models.PullRequest
	if err := r.db.QueryRowContext(ctx, queryMergePR, prID).Scan(
		&prDB.PullRequestID, &prDB.PullRequestName, &prDB.AuthorID, &prDB.Status,
		&prDB.CreatedAt, &prDB.MergedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	reviewers, err := r.getPRReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	return &models.PullRequestWithReviewers{
		PullRequest:       prDB,
		AssignedReviewers: reviewers,
	}, nil
}

func (r *PRRepository) GetUserReviewPRs(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	rows, err := r.db.QueryContext(ctx, queryGetUserReviewPRs, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user review PRs: %w", err)
	}
	defer rows.Close()

	var prs []models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		if err = rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, fmt.Errorf("failed to scan PR: %w", err)
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *PRRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, queryPRExists, prID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check PR existence: %w", err)
	}

	return exists, nil
}

func (r *PRRepository) getPRReviewers(ctx context.Context, prID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, queryGetPRReviewers, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR reviewers: %w", err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var reviewer string
		if err := rows.Scan(&reviewer); err != nil {
			return nil, fmt.Errorf("failed to scan reviewer: %w", err)
		}
		reviewers = append(reviewers, reviewer)
	}

	return reviewers, nil
}

// Новый метод для проверки, является ли пользователь ревьювером PR
func (r *PRRepository) IsPRReviewer(ctx context.Context, prID, userID string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, queryCheckPRReviewer, prID, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check PR reviewer: %w", err)
	}
	return exists, nil
}
