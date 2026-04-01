SET search_path TO coderz, public;

ALTER TABLE users
    ADD COLUMN username VARCHAR(80),
    ADD COLUMN bio TEXT,
    ADD COLUMN github_url TEXT,
    ADD COLUMN linkedin_url TEXT;

ALTER TABLE users
    ALTER COLUMN username SET DEFAULT ('user_' || REPLACE(LEFT(uuidv7()::text, 8), '-', ''));

WITH prepared AS (
    SELECT
        id,
        COALESCE(
            NULLIF(
                LOWER(REGEXP_REPLACE(SPLIT_PART(COALESCE(email, ''), '@', 1), '[^a-zA-Z0-9_]+', '', 'g')),
                ''
            ),
            NULLIF(
                LOWER(REGEXP_REPLACE(COALESCE(name, ''), '[^a-zA-Z0-9_]+', '', 'g')),
                ''
            ),
            'user'
        ) AS base_username
    FROM users
),
ranked AS (
    SELECT
        id,
        base_username,
        ROW_NUMBER() OVER (PARTITION BY base_username ORDER BY id) AS seq
    FROM prepared
)
UPDATE users u
SET username = CASE
    WHEN ranked.seq = 1 THEN ranked.base_username
    ELSE ranked.base_username || ranked.seq::text
END
FROM ranked
WHERE ranked.id = u.id
  AND u.username IS NULL;

ALTER TABLE users
    ALTER COLUMN username SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT uq_users_username UNIQUE (username),
    ADD CONSTRAINT chk_users_username_format CHECK (username ~ '^[a-z0-9_]+$');

ALTER TABLE bootcamp_enrollments
    ADD COLUMN assigned_sheet_key VARCHAR(64);

ALTER TABLE bootcamp_enrollments
    ADD CONSTRAINT chk_bootcamp_enrollments_assigned_sheet_key
    CHECK (
        assigned_sheet_key IS NULL
        OR assigned_sheet_key IN ('gfg-dsa-360', 'strivers-dsa-sheet')
    );

ALTER TABLE assignment_problems
    ADD COLUMN resources TEXT,
    ADD COLUMN app_progress_status VARCHAR(32) NOT NULL DEFAULT 'not_started';

ALTER TABLE assignment_problems
    ADD CONSTRAINT chk_assignment_problems_app_progress_status
    CHECK (
        app_progress_status IN (
            'not_started',
            'discussion_needed',
            'revision_needed',
            'completed'
        )
    );

CREATE TABLE mentee_requests (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID        NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    bootcamp_id     UUID        NOT NULL REFERENCES bootcamps(id) ON DELETE CASCADE,
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',
    sheet_key       VARCHAR(64),
    reviewed_by     UUID        REFERENCES organization_members(id) ON DELETE SET NULL,
    reviewed_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_mentee_requests_user_bootcamp UNIQUE (user_id, bootcamp_id),
    CONSTRAINT chk_mentee_requests_status CHECK (status IN ('pending', 'approved', 'rejected')),
    CONSTRAINT chk_mentee_requests_sheet_key CHECK (
        sheet_key IS NULL
        OR sheet_key IN ('gfg-dsa-360', 'strivers-dsa-sheet')
    )
);

CREATE TRIGGER trg_mentee_requests_updated_at
    BEFORE UPDATE ON mentee_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_mentee_requests_bootcamp_status ON mentee_requests(bootcamp_id, status);
CREATE INDEX idx_mentee_requests_user_id ON mentee_requests(user_id);

CREATE TABLE mentee_day_assignments (
    id                     UUID PRIMARY KEY DEFAULT uuidv7(),
    bootcamp_enrollment_id UUID        NOT NULL REFERENCES bootcamp_enrollments(id) ON DELETE CASCADE,
    weekday                VARCHAR(16) NOT NULL,
    created_by             UUID        REFERENCES organization_members(id) ON DELETE SET NULL,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_mentee_day_assignments UNIQUE (bootcamp_enrollment_id, weekday),
    CONSTRAINT chk_mentee_day_assignments_weekday CHECK (
        weekday IN (
            'monday',
            'tuesday',
            'wednesday',
            'thursday',
            'friday',
            'saturday',
            'sunday'
        )
    )
);

CREATE TRIGGER trg_mentee_day_assignments_updated_at
    BEFORE UPDATE ON mentee_day_assignments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_mentee_day_assignments_weekday ON mentee_day_assignments(weekday);
