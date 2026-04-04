
DROP TRIGGER IF EXISTS trigger_update_last_active ON tasks;
DROP FUNCTION IF EXISTS update_last_active();


ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_user_id;


ALTER TABLE tasks DROP COLUMN IF EXISTS user_id;


DROP INDEX IF EXISTS idx_tasks_user_id;
DROP INDEX IF EXISTS idx_users_telegram_id;


DROP TABLE IF EXISTS users;