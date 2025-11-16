package repository

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type IStatsRepository interface {
	GetUserAssignmentStats(context.Context) ([]models.UserAssignmentStats, error)
	GetPRAssignmentStats(context.Context) ([]models.PRAssignmentStats, error)
	GetTotalOpenPRs(context.Context) (int, error)
}
