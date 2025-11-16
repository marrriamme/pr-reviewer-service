package stats

const (
	queryGetStatsUserAssignments = `
        SELECT reviewer_id, COUNT(*) as assignment_count
        FROM pull_request_reviewers prr
        JOIN pull_requests pr ON prr.pull_request_id = pr.pull_request_id
        WHERE pr.status = 'OPEN'
        GROUP BY reviewer_id
        ORDER BY assignment_count DESC`

	queryGetStatsPRAssignments = `
        SELECT pr.pull_request_id, pr.pull_request_name, COUNT(prr.reviewer_id) as reviewer_count
        FROM pull_requests pr
        LEFT JOIN pull_request_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        WHERE pr.status = 'OPEN'
        GROUP BY pr.pull_request_id, pr.pull_request_name
        ORDER BY reviewer_count DESC`

	queryGetTotalOpenPRs = `
        SELECT COUNT(*) FROM pull_requests WHERE status = 'OPEN'`
)
