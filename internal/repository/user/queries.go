package user

const (
	queryGetUser = `
        SELECT user_id, username, team_name, is_active
        FROM users
        WHERE user_id = $1`

	queryUpdateUserActivity = `
        UPDATE users 
        SET is_active = $1, updated_at = NOW()
        WHERE user_id = $2
        RETURNING user_id, username, team_name, is_active`

	queryGetActiveTeamMembers = `
        SELECT user_id, username, team_name, is_active
        FROM users
        WHERE team_name = $1 AND is_active = true AND user_id != $2`

	queryGetRandomActiveTeamMember = `
        SELECT user_id
        FROM users
        WHERE team_name = $1 
          AND is_active = true 
          AND user_id != $2 
          AND user_id != $3 
        ORDER BY RANDOM()
        LIMIT 1`

	queryGetRandomActiveTeamMembers = `
        SELECT user_id
        FROM users
        WHERE team_name = $1 
          AND is_active = true 
          AND user_id != $2 
        ORDER BY RANDOM()
        LIMIT $3`

	queryUserExists = `
        SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`

	queryGetUserReviewPRs = `
        SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
        FROM pull_requests pr
        JOIN pull_request_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        WHERE prr.reviewer_id = $1 AND pr.status = 'OPEN'`
)
