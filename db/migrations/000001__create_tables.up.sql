DROP INDEX IF EXISTS idx_users_team_active;
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_pr_author_status;
DROP INDEX IF EXISTS idx_pr_status;
DROP INDEX IF EXISTS idx_pr_reviewers_gin;
DROP INDEX IF EXISTS idx_pr_created_at;

DROP TABLE IF EXISTS pull_request_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;

CREATE TABLE teams (
                       team_name TEXT PRIMARY KEY,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE users (
                       user_id TEXT PRIMARY KEY,
                       username TEXT NOT NULL,
                       team_name TEXT NOT NULL,
                       is_active BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                       CONSTRAINT fk_users_team
                           FOREIGN KEY (team_name)
                               REFERENCES teams(team_name)
                               ON DELETE CASCADE
);

CREATE TABLE pull_requests (
                               pull_request_id TEXT PRIMARY KEY,
                               pull_request_name TEXT NOT NULL,
                               author_id TEXT NOT NULL,
                               status TEXT NOT NULL DEFAULT 'OPEN',
                               created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                               merged_at TIMESTAMP WITH TIME ZONE,
                               CONSTRAINT fk_pr_author
                                   FOREIGN KEY (author_id)
                                       REFERENCES users(user_id),
                               CONSTRAINT chk_pr_status
                                   CHECK (status IN ('OPEN', 'MERGED'))
);

CREATE TABLE pull_request_reviewers (
                                        pull_request_id TEXT NOT NULL,
                                        reviewer_id TEXT NOT NULL,
                                        assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                        PRIMARY KEY (pull_request_id, reviewer_id),
                                        CONSTRAINT fk_pr_reviewer_pr
                                            FOREIGN KEY (pull_request_id)
                                                REFERENCES pull_requests(pull_request_id)
                                                ON DELETE CASCADE,
                                        CONSTRAINT fk_pr_reviewer_user
                                            FOREIGN KEY (reviewer_id)
                                                REFERENCES users(user_id)
                                                ON DELETE CASCADE
);

CREATE INDEX idx_users_team_active ON users(team_name, is_active);
CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = true;
CREATE INDEX idx_pr_author_status ON pull_requests(author_id, status);
CREATE INDEX idx_pr_status ON pull_requests(status);
CREATE INDEX idx_pr_created_at ON pull_requests(created_at);
CREATE INDEX idx_pr_reviewers_reviewer_id ON pull_request_reviewers(reviewer_id);
CREATE INDEX idx_pr_reviewers_pr_id ON pull_request_reviewers(pull_request_id);