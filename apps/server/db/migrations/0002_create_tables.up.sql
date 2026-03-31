-- ============================================================
-- 0002_create_tables.up.sql
-- Full schema for Coderz Space Bootcamp platform
-- ============================================================

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- ENUM TYPES
-- ============================================================

CREATE TYPE user_role AS ENUM ('user', 'super_admin');

CREATE TYPE org_status AS ENUM ('pending_approval', 'approved', 'suspended');

CREATE TYPE org_member_role AS ENUM ('admin', 'mentor', 'mentee');

CREATE TYPE enrollment_status AS ENUM ('active', 'inactive');

CREATE TYPE bootcamp_enrollment_role AS ENUM ('mentor', 'mentee');

CREATE TYPE difficulty_level AS ENUM ('easy', 'medium', 'hard');

CREATE TYPE assignment_status AS ENUM ('active', 'completed', 'expired');

CREATE TYPE assignment_problem_status AS ENUM ('pending', 'attempted', 'completed');

CREATE TYPE poll_vote_value AS ENUM ('easy', 'medium', 'hard');

-- ============================================================
-- HELPER: trigger function to auto-update updated_at
-- ============================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================
-- 1. AUTH MODULE — users & refresh_tokens
-- ============================================================

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100)  NOT NULL,
    email         VARCHAR(255)  UNIQUE,
    email_verified BOOLEAN       NOT NULL DEFAULT FALSE,
    password_hash TEXT,
    role          user_role     NOT NULL DEFAULT 'user',
    google_id     VARCHAR(255)  UNIQUE,
    avatar_url    TEXT,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- At least one auth method must be present
    CONSTRAINT chk_auth_method CHECK (
        password_hash IS NOT NULL OR google_id IS NOT NULL
    )
);

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Refresh tokens for sessions
CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT        NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER trg_refresh_tokens_updated_at
    BEFORE UPDATE ON refresh_tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- ============================================================
-- 2. ORGANIZATION MODULE
-- ============================================================

-- 2a. organizations
CREATE TABLE organizations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255)  NOT NULL,
    slug        VARCHAR(255)  NOT NULL UNIQUE,
    description TEXT,
    status      org_status    NOT NULL DEFAULT 'pending_approval',
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER trg_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 2b. organization_members
CREATE TABLE organization_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID            NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         UUID            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            org_member_role NOT NULL,
    joined_at       TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_org_member UNIQUE (organization_id, user_id)
);

CREATE INDEX idx_org_members_org_id  ON organization_members(organization_id);
CREATE INDEX idx_org_members_user_id ON organization_members(user_id);

-- ============================================================
-- 3. BOOTCAMP MODULE
-- ============================================================

-- 3a. bootcamps
CREATE TABLE bootcamps (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_by      UUID         NOT NULL REFERENCES organization_members(id) ON DELETE RESTRICT,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    start_date      DATE,
    end_date        DATE,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    archived_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_bootcamp_dates CHECK (
        start_date IS NULL OR end_date IS NULL OR start_date <= end_date
    )
);

CREATE TRIGGER trg_bootcamps_updated_at
    BEFORE UPDATE ON bootcamps
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_bootcamps_org_id ON bootcamps(organization_id);

-- 3b. bootcamp_enrollments
CREATE TABLE bootcamp_enrollments (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bootcamp_id            UUID                    NOT NULL REFERENCES bootcamps(id) ON DELETE CASCADE,
    organization_member_id UUID                    NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    role                   bootcamp_enrollment_role NOT NULL,
    status                 enrollment_status        NOT NULL DEFAULT 'active',
    enrolled_at            TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_bootcamp_enrollment UNIQUE (bootcamp_id, organization_member_id)
);

CREATE INDEX idx_bootcamp_enrollments_bootcamp ON bootcamp_enrollments(bootcamp_id);
CREATE INDEX idx_bootcamp_enrollments_member   ON bootcamp_enrollments(organization_member_id);

-- ============================================================
-- 4. PROBLEM CONTENT MODULE
-- ============================================================

-- 4a. problems
CREATE TABLE problems (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID             NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_by      UUID             NOT NULL REFERENCES organization_members(id) ON DELETE RESTRICT,
    title           VARCHAR(255)     NOT NULL,
    description     TEXT,
    difficulty      difficulty_level NOT NULL,
    external_link   TEXT,
    archived_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER trg_problems_updated_at
    BEFORE UPDATE ON problems
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_problems_org_id ON problems(organization_id);

-- 4b. tags
CREATE TABLE tags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    created_by      UUID         NOT NULL REFERENCES organization_members(id) ON DELETE RESTRICT,
    name            VARCHAR(100) NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_tag_per_org UNIQUE (organization_id, name)
);

CREATE INDEX idx_tags_org_id ON tags(organization_id);

-- 4c. problem_tags (join table — composite PK)
CREATE TABLE problem_tags (
    problem_id UUID        NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    tag_id     UUID        NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (problem_id, tag_id)
);

-- 4d. problem_resources
CREATE TABLE problem_resources (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID         NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    title      VARCHAR(255),
    url        TEXT,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_problem_resources_problem ON problem_resources(problem_id);

-- ============================================================
-- 5. ASSIGNMENT LAYER
-- ============================================================

-- 5a. assignment_groups (templates)
CREATE TABLE assignment_groups (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bootcamp_id   UUID         NOT NULL REFERENCES bootcamps(id) ON DELETE CASCADE,
    created_by    UUID         NOT NULL REFERENCES organization_members(id) ON DELETE RESTRICT,
    title         VARCHAR(255) NOT NULL,
    description   TEXT,
    deadline_days INTEGER,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER trg_assignment_groups_updated_at
    BEFORE UPDATE ON assignment_groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_assignment_groups_bootcamp ON assignment_groups(bootcamp_id);

-- 5b. assignment_group_problems (join table — composite PK)
CREATE TABLE assignment_group_problems (
    assignment_group_id UUID        NOT NULL REFERENCES assignment_groups(id) ON DELETE CASCADE,
    problem_id          UUID        NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    position            INTEGER,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (assignment_group_id, problem_id)
);

-- 5c. assignments (per-mentee instances)
CREATE TABLE assignments (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_group_id    UUID              NOT NULL REFERENCES assignment_groups(id) ON DELETE CASCADE,
    bootcamp_enrollment_id UUID              NOT NULL REFERENCES bootcamp_enrollments(id) ON DELETE CASCADE,
    assigned_by            UUID              NOT NULL REFERENCES organization_members(id) ON DELETE RESTRICT,
    assigned_at            TIMESTAMPTZ       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deadline_at            TIMESTAMPTZ,
    status                 assignment_status NOT NULL DEFAULT 'active',
    archived_at            TIMESTAMPTZ,
    created_at             TIMESTAMPTZ       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMPTZ       NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER trg_assignments_updated_at
    BEFORE UPDATE ON assignments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_assignments_group      ON assignments(assignment_group_id);
CREATE INDEX idx_assignments_enrollment ON assignments(bootcamp_enrollment_id);

-- 5d. assignment_problems (per-problem progress tracking)
CREATE TABLE assignment_problems (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id UUID                      NOT NULL REFERENCES assignments(id) ON DELETE CASCADE,
    problem_id    UUID                      NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    status        assignment_problem_status NOT NULL DEFAULT 'pending',
    solution_link TEXT,
    notes         TEXT,
    completed_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ               NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_assignment_problem UNIQUE (assignment_id, problem_id)
);

CREATE TRIGGER trg_assignment_problems_updated_at
    BEFORE UPDATE ON assignment_problems
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================
-- 6. PROGRESS TRACKING — doubts
-- ============================================================

CREATE TABLE doubts (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_problem_id UUID        NOT NULL REFERENCES assignment_problems(id) ON DELETE CASCADE,
    raised_by             UUID        NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    message               TEXT        NOT NULL,
    resolved              BOOLEAN     NOT NULL DEFAULT FALSE,
    resolved_by           UUID        REFERENCES organization_members(id) ON DELETE SET NULL,
    resolved_at           TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_doubts_assignment_problem ON doubts(assignment_problem_id);
CREATE INDEX idx_doubts_raised_by          ON doubts(raised_by);

-- ============================================================
-- 7. ANALYTICS LAYER
-- ============================================================

-- 7a. leaderboard_entries (snapshot table)
CREATE TABLE leaderboard_entries (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bootcamp_id            UUID        NOT NULL REFERENCES bootcamps(id) ON DELETE CASCADE,
    bootcamp_enrollment_id UUID        NOT NULL REFERENCES bootcamp_enrollments(id) ON DELETE CASCADE,
    problems_completed     INTEGER     NOT NULL DEFAULT 0,
    problems_attempted     INTEGER     NOT NULL DEFAULT 0,
    completion_rate        REAL        NOT NULL DEFAULT 0.0,
    streak_days            INTEGER     NOT NULL DEFAULT 0,
    score                  INTEGER     NOT NULL DEFAULT 0,
    rank                   INTEGER     NOT NULL DEFAULT 0,
    calculated_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_leaderboard_entry UNIQUE (bootcamp_id, bootcamp_enrollment_id)
);

CREATE INDEX idx_leaderboard_bootcamp ON leaderboard_entries(bootcamp_id);

-- 7b. polls
CREATE TABLE polls (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bootcamp_id UUID         NOT NULL REFERENCES bootcamps(id) ON DELETE CASCADE,
    problem_id  UUID         NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    question    VARCHAR(500) NOT NULL,
    created_by  UUID         NOT NULL REFERENCES organization_members(id) ON DELETE RESTRICT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_polls_bootcamp ON polls(bootcamp_id);

-- 7c. poll_votes
CREATE TABLE poll_votes (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    poll_id    UUID            NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    voter_id   UUID            NOT NULL REFERENCES bootcamp_enrollments(id) ON DELETE CASCADE,
    vote       poll_vote_value NOT NULL,
    created_at TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_poll_vote UNIQUE (poll_id, voter_id)
);

CREATE INDEX idx_poll_votes_poll ON poll_votes(poll_id);
