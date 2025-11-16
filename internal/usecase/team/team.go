package team

import (
	"context"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
	"github.com/marrria_mme/pr-reviewer-service/internal/repository"
)

type TeamUsecase struct {
	repo repository.ITeamRepository
}

func NewTeamUsecase(repo repository.ITeamRepository) *TeamUsecase {
	return &TeamUsecase{repo: repo}
}

func (u *TeamUsecase) CreateTeam(ctx context.Context, team *models.Team) error {
	exists, err := u.repo.TeamExists(ctx, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to check team existence: %w", err)
	}

	if exists {
		return fmt.Errorf("team_name already exists: %w", errs.ErrTeamExists)
	}

	if err = u.repo.CreateTeam(ctx, team); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	return nil
}

func (u *TeamUsecase) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := u.repo.GetTeam(ctx, teamName)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, fmt.Errorf("team not found: %w", errs.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}
