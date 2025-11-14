package repository

import (
	"context"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type ITeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
}
