ALTER TABLE tasks
ADD COLUMN IF NOT EXISTS due_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS reminder_sent BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS reminder_offset VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
CREATE INDEX IF NOT EXISTS idx_tasks_reminder_sent ON tasks(reminder_sent);
