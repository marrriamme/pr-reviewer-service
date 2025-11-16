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
	PR *PullRequestWithReviewersDTO `json:"pr"`
}

type ReassignResponseDTO struct {
	PR         *PullRequestWithReviewersDTO `json:"pr"`
	ReplacedBy string                       `json:"replaced_by"`
}

type PullRequestWithReviewersDTO struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	CreatedAt         *time.Time `json:"created_at"`
	MergedAt          *time.Time `json:"merged_at"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
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
	prDTO := &PullRequestWithReviewersDTO{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		AssignedReviewers: pr.AssignedReviewers,
	}
	return PRResponseDTO{PR: prDTO}
}

func ToReassignResponseDTO(pr *models.PullRequestWithReviewers, newUserID string) ReassignResponseDTO {
	prDTO := &PullRequestWithReviewersDTO{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		AssignedReviewers: pr.AssignedReviewers,
	}
	return ReassignResponseDTO{
		PR:         prDTO,
		ReplacedBy: newUserID,
	}
}
