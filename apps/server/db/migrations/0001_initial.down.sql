-- ============================================================
-- 0001_initial.down.sql
-- Drop all tables and types in reverse dependency order
-- ============================================================
-- Analytics
DROP SCHEMA IF EXISTS coderz CASCADE;

SET search_path TO coderz, public;

DROP TABLE IF EXISTS poll_votes;

DROP TABLE IF EXISTS polls;

DROP TABLE IF EXISTS leaderboard_entries;

-- Progress Tracking
DROP TABLE IF EXISTS doubts;

-- Assignment Layer
DROP TABLE IF EXISTS assignment_problems;

DROP TABLE IF EXISTS assignments;

DROP TABLE IF EXISTS assignment_group_problems;

DROP TABLE IF EXISTS assignment_groups;

-- Problem Content
DROP TABLE IF EXISTS problem_resources;

DROP TABLE IF EXISTS problem_tags;

DROP TABLE IF EXISTS tags;

DROP TABLE IF EXISTS problems;

-- Bootcamp
DROP TABLE IF EXISTS bootcamp_enrollments;

DROP TABLE IF EXISTS bootcamps;

-- Organization
DROP TABLE IF EXISTS organization_members;

DROP TABLE IF EXISTS organizations;

-- Auth
DROP TABLE IF EXISTS refresh_tokens;

DROP TABLE IF EXISTS users;

-- Trigger function
DROP FUNCTION IF EXISTS update_updated_at_column ();

-- Enum types
DROP TYPE IF EXISTS poll_vote_value;

DROP TYPE IF EXISTS assignment_problem_status;

DROP TYPE IF EXISTS assignment_status;

DROP TYPE IF EXISTS difficulty_level;

DROP TYPE IF EXISTS bootcamp_enrollment_role;

DROP TYPE IF EXISTS enrollment_status;

DROP TYPE IF EXISTS org_member_role;

DROP TYPE IF EXISTS org_status;

DROP TYPE IF EXISTS user_role;