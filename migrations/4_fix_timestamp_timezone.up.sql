-- Convert due_date to TIMESTAMP WITH TIME ZONE for proper timezone handling
ALTER TABLE tasks 
ALTER COLUMN due_date TYPE TIMESTAMP WITH TIME ZONE USING due_date AT TIME ZONE 'Europe/Moscow';

-- Also update created_at and completed_at to ensure consistent timezone handling
ALTER TABLE tasks
ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING created_at AT TIME ZONE 'Europe/Moscow';

ALTER TABLE tasks
ALTER COLUMN completed_at TYPE TIMESTAMP WITH TIME ZONE USING completed_at AT TIME ZONE 'Europe/Moscow';

-- Update user table timestamps too
ALTER TABLE users
ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE USING created_at AT TIME ZONE 'Europe/Moscow';

ALTER TABLE users
ALTER COLUMN last_active TYPE TIMESTAMP WITH TIME ZONE USING last_active AT TIME ZONE 'Europe/Moscow';
