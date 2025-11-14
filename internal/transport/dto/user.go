package dto

import "github.com/marrria_mme/pr-reviewer-service/internal/models"

type UserActivityRequestDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponseDTO struct {
	User *models.User `json:"user"`
}

type UserReviewResponseDTO struct {
	UserID       string                    `json:"user_id"`
	PullRequests []models.PullRequestShort `json:"pull_requests"`
}

func ToUserResponseDTO(user *models.User) UserResponseDTO {
	return UserResponseDTO{User: user}
}
