package stats

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetUserAssignmentStats(ctx context.Context) ([]models.UserAssignmentStats, error) {
	rows, err := r.db.QueryContext(ctx, queryGetStatsUserAssignments)
	if err != nil {
		return nil, fmt.Errorf("failed to get user assignment stats: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var stats []models.UserAssignmentStats
	for rows.Next() {
		var stat models.UserAssignmentStats
		if err = rows.Scan(&stat.UserID, &stat.AssignmentCount); err != nil {
			return nil, fmt.Errorf("failed to scan user assignment stat: %w", err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (r *StatsRepository) GetPRAssignmentStats(ctx context.Context) ([]models.PRAssignmentStats, error) {
	rows, err := r.db.QueryContext(ctx, queryGetStatsPRAssignments)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR assignment stats: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var stats []models.PRAssignmentStats
	for rows.Next() {
		var stat models.PRAssignmentStats
		if err = rows.Scan(&stat.PullRequestID, &stat.PullRequestName, &stat.ReviewerCount); err != nil {
			return nil, fmt.Errorf("failed to scan PR assignment stat: %w", err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (r *StatsRepository) GetTotalOpenPRs(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, queryGetTotalOpenPRs).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total open PRs: %w", err)
	}
	return count, nil
}
