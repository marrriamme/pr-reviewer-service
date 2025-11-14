package pr

const (
	queryCreatePR = `
        INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
        VALUES ($1, $2, $3, $4)
        RETURNING pull_request_id, pull_request_name, author_id, status, created_at, merged_at`

	queryGetPR = `
        SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at
        FROM pull_requests pr
        WHERE pr.pull_request_id = $1`

	queryMergePR = `
        UPDATE pull_requests 
        SET status = 'MERGED', merged_at = NOW()
        WHERE pull_request_id = $1
        RETURNING pull_request_id, pull_request_name, author_id, status, created_at, merged_at`

	queryGetUserReviewPRs = `
        SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
        FROM pull_requests pr
        JOIN pull_request_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        WHERE prr.reviewer_id = $1 AND pr.status = 'OPEN'`

	queryPRExists = `
        SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`

	queryAddReviewer = `
        INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
        VALUES ($1, $2)
        ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING`

	queryDeletePRReviewers = `
        DELETE FROM pull_request_reviewers 
        WHERE pull_request_id = $1`

	queryGetPRReviewers = `
        SELECT reviewer_id 
        FROM pull_request_reviewers 
        WHERE pull_request_id = $1`

	queryCheckPRReviewer = `
        SELECT EXISTS(
            SELECT 1 FROM pull_request_reviewers 
            WHERE pull_request_id = $1 AND reviewer_id = $2
        )`
)
