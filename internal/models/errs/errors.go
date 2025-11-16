package errs

import (
	"errors"
)

var (
	ErrTeamExists    = errors.New("team_name already exists")
	ErrPRExists      = errors.New("PR already exists")
	ErrPRMerged      = errors.New("PR is merged")
	ErrNotAssigned   = errors.New("reviewer not assigned")
	ErrNoCandidate   = errors.New("no active replacement candidate")
	ErrNotFound      = errors.New("resource not found")
	ErrUserNotActive = errors.New("user is not active")
	ErrUserNoTeam    = errors.New("user does not belong to any team")
)
