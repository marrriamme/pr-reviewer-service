package team

const (
	queryCreateTeam = `
        INSERT INTO teams (team_name)
        VALUES ($1)
        RETURNING team_name`

	queryTeamExists = `
        SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`

	queryAddTeamMembers = `
        INSERT INTO users (user_id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) 
        DO UPDATE SET 
            username = EXCLUDED.username,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active,
            updated_at = NOW()`

	queryGetTeamMembers = `
        SELECT user_id, username, is_active
        FROM users
        WHERE team_name = $1`
)
