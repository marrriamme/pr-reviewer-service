package usecase

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IStatsUsecase interface {
	GetStats(ctx context.Context) (*models.StatsResponse, error)
}
