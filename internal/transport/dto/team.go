package dto

import "github.com/marrria_mme/pr-reviewer-service/internal/models"

type TeamRequestDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamResponseDTO struct {
	Team *models.Team `json:"team"`
}

func ToTeamModel(dto TeamRequestDTO) *models.Team {
	members := make([]models.TeamMember, len(dto.Members))
	for i, member := range dto.Members {
		members[i] = models.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return &models.Team{
		TeamName: dto.TeamName,
		Members:  members,
	}
}

func ToTeamResponseDTO(team *models.Team) TeamResponseDTO {
	return TeamResponseDTO{Team: team}
}
