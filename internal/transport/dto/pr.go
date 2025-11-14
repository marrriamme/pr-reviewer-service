package dto

import (
	"time"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
)

type CreatePRRequestDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequestDTO struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignPRRequestDTO struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type PRResponseDTO struct {
	PR *models.PullRequestWithReviewers `json:"pr"`
}

type ReassignResponseDTO struct {
	PR         *models.PullRequestWithReviewers `json:"pr"`
	ReplacedBy string                           `json:"replaced_by"`
}

func ToPRModel(dto CreatePRRequestDTO) *models.PullRequestWithReviewers {
	now := time.Now()
	return &models.PullRequestWithReviewers{
		PullRequest: models.PullRequest{
			PullRequestID:   dto.PullRequestID,
			PullRequestName: dto.PullRequestName,
			AuthorID:        dto.AuthorID,
			Status:          "OPEN",
			CreatedAt:       &now,
		},
		AssignedReviewers: []string{},
	}
}

func ToPRResponseDTO(pr *models.PullRequestWithReviewers) PRResponseDTO {
	return PRResponseDTO{PR: pr}
}

func ToReassignResponseDTO(pr *models.PullRequestWithReviewers, newUserID string) ReassignResponseDTO {
	return ReassignResponseDTO{
		PR:         pr,
		ReplacedBy: newUserID,
	}
}
