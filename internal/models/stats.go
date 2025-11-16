package models

type UserAssignmentStats struct {
	UserID          string
	AssignmentCount int
}

type PRAssignmentStats struct {
	PullRequestID   string
	PullRequestName string
	ReviewerCount   int
}

type StatsResponse struct {
	UserAssignments  []UserAssignmentStats
	PRAssignments    []PRAssignmentStats
	TotalOpenPRs     int
	TotalAssignments int
}
