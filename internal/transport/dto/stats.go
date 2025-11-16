package dto

import "github.com/marrria_mme/pr-reviewer-service/internal/models"

type UserAssignmentStatsDTO struct {
	UserID          string `json:"user_id"`
	AssignmentCount int    `json:"assignment_count"`
}

type PRAssignmentStatsDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	ReviewerCount   int    `json:"reviewer_count"`
}

type StatsResponseDTO struct {
	UserAssignments  []UserAssignmentStatsDTO `json:"user_assignments"`
	PRAssignments    []PRAssignmentStatsDTO   `json:"pr_assignments"`
	TotalOpenPRs     int                      `json:"total_open_prs"`
	TotalAssignments int                      `json:"total_assignments"`
}

func ToUserAssignmentStatsDTO(model models.UserAssignmentStats) UserAssignmentStatsDTO {
	return UserAssignmentStatsDTO{
		UserID:          model.UserID,
		AssignmentCount: model.AssignmentCount,
	}
}

func ToPRAssignmentStatsDTO(model models.PRAssignmentStats) PRAssignmentStatsDTO {
	return PRAssignmentStatsDTO{
		PullRequestID:   model.PullRequestID,
		PullRequestName: model.PullRequestName,
		ReviewerCount:   model.ReviewerCount,
	}
}

func ToStatsResponseDTO(model *models.StatsResponse) StatsResponseDTO {
	userStatsDTO := make([]UserAssignmentStatsDTO, len(model.UserAssignments))
	for i, stat := range model.UserAssignments {
		userStatsDTO[i] = ToUserAssignmentStatsDTO(stat)
	}

	prStatsDTO := make([]PRAssignmentStatsDTO, len(model.PRAssignments))
	for i, stat := range model.PRAssignments {
		prStatsDTO[i] = ToPRAssignmentStatsDTO(stat)
	}

	return StatsResponseDTO{
		UserAssignments:  userStatsDTO,
		PRAssignments:    prStatsDTO,
		TotalOpenPRs:     model.TotalOpenPRs,
		TotalAssignments: model.TotalAssignments,
	}
}
