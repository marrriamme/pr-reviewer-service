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

type PullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type UserReviewResponseDTO struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
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

func ToPullRequestShortDTO(pr models.PullRequestShort) PullRequestShortDTO {
	return PullRequestShortDTO{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status,
	}
}

func NewUserReviewResponseDTO(userID string, prs []models.PullRequestShort) UserReviewResponseDTO {
	prDTOs := make([]PullRequestShortDTO, len(prs))
	for i, pr := range prs {
		prDTOs[i] = ToPullRequestShortDTO(pr)
	}

	return UserReviewResponseDTO{
		UserID:       userID,
		PullRequests: prDTOs,
	}
}
