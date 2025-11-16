package repository

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type ITeamRepository interface {
	CreateTeam(context.Context, *models.Team) error
	GetTeam(context.Context, string) (*models.Team, error)
	TeamExists(context.Context, string) (bool, error)
}
