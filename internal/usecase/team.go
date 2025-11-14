package usecase

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type ITeamUsecase interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}
