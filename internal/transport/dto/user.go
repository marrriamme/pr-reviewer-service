package dto

import "github.com/marrria_mme/pr-reviewer-service/internal/models"

type UserActivityRequestDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponseDTO struct {
	User *UserDTO `json:"user"`
}

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserReviewResponseDTO struct {
	UserID       string                    `json:"user_id"`
	PullRequests []models.PullRequestShort `json:"pull_requests"`
}

func ToUserResponseDTO(user *models.User) UserResponseDTO {
	userDTO := &UserDTO{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
	return UserResponseDTO{User: userDTO}
}
