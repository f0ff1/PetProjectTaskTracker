CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE tasks ADD COLUMN IF NOT EXISTS user_id INTEGER;


ALTER TABLE tasks ADD CONSTRAINT fk_tasks_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;


CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);


INSERT INTO users (telegram_id, username, first_name, is_admin, created_at)
VALUES (0, 'system', 'System User', false, CURRENT_TIMESTAMP)
ON CONFLICT (telegram_id) DO NOTHING;


UPDATE tasks SET user_id = (SELECT id FROM users WHERE telegram_id = 0) 
WHERE user_id IS NULL;


ALTER TABLE tasks ALTER COLUMN user_id SET NOT NULL;


CREATE OR REPLACE FUNCTION update_last_active()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users SET last_active = CURRENT_TIMESTAMP 
    WHERE id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_last_active
    AFTER INSERT OR UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_last_active();