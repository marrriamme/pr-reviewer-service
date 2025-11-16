package models

import "time"

type PullRequest struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          string
	CreatedAt       *time.Time
	MergedAt        *time.Time
}

type PullRequestWithReviewers struct {
	PullRequest
	AssignedReviewers []string
}

type PullRequestShort struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          string
}
