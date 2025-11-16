package team

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	if err = tx.QueryRowContext(ctx, queryTeamExists, team.TeamName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check team existence: %w", err)
	}
	if exists {
		return errs.ErrTeamExists
	}

	if _, err = tx.ExecContext(ctx, queryCreateTeam, team.TeamName); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	for _, member := range team.Members {
		if _, err = tx.ExecContext(ctx, queryAddTeamMembers,
			member.UserID, member.Username, team.TeamName, member.IsActive); err != nil {
			return fmt.Errorf("failed to add team member: %w", err)
		}
	}

	return tx.Commit()
}

func (r *TeamRepository) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	var team models.Team
	team.TeamName = teamName

	rows, err := r.db.QueryContext(ctx, queryGetTeamMembers, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		if err = rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan team member: %w", err)
		}
		members = append(members, member)
	}

	if len(members) == 0 {
		return nil, errs.ErrNotFound
	}

	team.Members = members

	return &team, nil
}

func (r *TeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, queryTeamExists, teamName).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}

	return exists, nil
}
