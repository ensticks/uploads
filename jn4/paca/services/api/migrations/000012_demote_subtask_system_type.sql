-- Demote the legacy "Subtask" system type to a regular user-visible type.
-- Projects created before v0.5.0 were seeded with a system-managed Subtask type
-- that can no longer be managed by users. Making it non-system lets users rename
-- or delete it while preserving all existing task data.

UPDATE task_types
SET is_system = false
WHERE name = 'Subtask'
  AND is_system = true;
