-- 000014_soft_delete_projects.sql
-- Adds deleted_at to projects, converting hard-delete to soft-delete.

BEGIN;

ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects (deleted_at) WHERE deleted_at IS NOT NULL;

COMMIT;
