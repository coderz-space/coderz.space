SET search_path TO coderz, public;

DROP TRIGGER IF EXISTS trg_mentee_day_assignments_updated_at ON mentee_day_assignments;
DROP TABLE IF EXISTS mentee_day_assignments;

DROP TRIGGER IF EXISTS trg_mentee_requests_updated_at ON mentee_requests;
DROP TABLE IF EXISTS mentee_requests;

ALTER TABLE assignment_problems
    DROP CONSTRAINT IF EXISTS chk_assignment_problems_app_progress_status;

ALTER TABLE assignment_problems
    DROP COLUMN IF EXISTS app_progress_status,
    DROP COLUMN IF EXISTS resources;

ALTER TABLE bootcamp_enrollments
    DROP CONSTRAINT IF EXISTS chk_bootcamp_enrollments_assigned_sheet_key;

ALTER TABLE bootcamp_enrollments
    DROP COLUMN IF EXISTS assigned_sheet_key;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS chk_users_username_format,
    DROP CONSTRAINT IF EXISTS uq_users_username;

ALTER TABLE users
    DROP COLUMN IF EXISTS linkedin_url,
    DROP COLUMN IF EXISTS github_url,
    DROP COLUMN IF EXISTS bio,
    DROP COLUMN IF EXISTS username;
