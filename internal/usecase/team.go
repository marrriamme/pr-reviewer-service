package usecase

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type ITeamUsecase interface {
	CreateTeam(context.Context, *models.Team) error
	GetTeam(context.Context, string) (*models.Team, error)
}
