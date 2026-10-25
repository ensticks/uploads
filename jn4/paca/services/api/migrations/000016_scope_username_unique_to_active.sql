-- 000016_scope_username_unique_to_active.sql
-- Scopes the username uniqueness constraint to active (non-deleted) users
-- so a username freed up by a soft-deleted account can be reused, matching
-- the uq_project_members_active pattern used elsewhere.

BEGIN;

DROP INDEX IF EXISTS uni_users_username;

CREATE UNIQUE INDEX IF NOT EXISTS uni_users_username_active
    ON users (username)
    WHERE deleted_at IS NULL;

COMMIT;
