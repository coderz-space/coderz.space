-- 0002_web_ui_sync.down.sql
-- Removing values from ENUM types is not supported in PostgreSQL without recreating the type and rewriting data.
-- We leave this intentionally a no-op to prevent data destruction.
-- To rollback, you would ideally:
-- 1. DELETE FROM assignments/records using the new states or UPDATE to 'pending'.
-- 2. CREATE TYPE xxx_new AS ENUM (...)
-- 3. ALTER TABLE ... ALTER COLUMN ... TYPE xxx_new
-- 4. DROP TYPE xxx
-- 5. ALTER TYPE xxx_new RENAME TO xxx

SELECT 1;
