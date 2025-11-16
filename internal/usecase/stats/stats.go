package stats

import (
	"context"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/repository"
)

type StatsUsecase struct {
	statsRepo repository.IStatsRepository
}

func NewStatsUsecase(statsRepo repository.IStatsRepository) *StatsUsecase {
	return &StatsUsecase{
		statsRepo: statsRepo,
	}
}

func (u *StatsUsecase) GetStats(ctx context.Context) (*models.StatsResponse, error) {
	userStats, err := u.statsRepo.GetUserAssignmentStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	prStats, err := u.statsRepo.GetPRAssignmentStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR stats: %w", err)
	}

	totalOpenPRs, err := u.statsRepo.GetTotalOpenPRs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total open PRs: %w", err)
	}
	
	totalAssignments := 0
	for _, stat := range userStats {
		totalAssignments += stat.AssignmentCount
	}

	return &models.StatsResponse{
		UserAssignments:  userStats,
		PRAssignments:    prStats,
		TotalOpenPRs:     totalOpenPRs,
		TotalAssignments: totalAssignments,
	}, nil
}
