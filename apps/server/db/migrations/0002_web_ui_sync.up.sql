-- 0002_web_ui_sync.up.sql
-- Adds UI required states to progress statuses and enrollment statuses

SET search_path TO coderz, public;

-- Modify assignment_problem_status
ALTER TYPE assignment_problem_status ADD VALUE IF NOT EXISTS 'not_started';
ALTER TYPE assignment_problem_status ADD VALUE IF NOT EXISTS 'discussion_needed';
ALTER TYPE assignment_problem_status ADD VALUE IF NOT EXISTS 'revision_needed';

-- Modify enrollment_status
ALTER TYPE enrollment_status ADD VALUE IF NOT EXISTS 'pending';

