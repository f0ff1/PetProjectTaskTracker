CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255) DEFAULT '',
    first_name VARCHAR(255) DEFAULT '',
    last_name VARCHAR(255) DEFAULT '',
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    user_task_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    tags TEXT[],
    CONSTRAINT fk_tasks_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_user_task_id ON tasks(user_id, user_task_id);
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);
CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(completed);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
CREATE INDEX IF NOT EXISTS idx_tasks_tags ON tasks USING GIN(tags);


CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_user_task_id ON tasks(user_id, user_task_id);


CREATE OR REPLACE FUNCTION set_user_task_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.user_task_id := COALESCE(
        (SELECT MAX(user_task_id) FROM tasks WHERE user_id = NEW.user_id), 0
    ) + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS trigger_set_user_task_id ON tasks;
CREATE TRIGGER trigger_set_user_task_id
    BEFORE INSERT ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION set_user_task_id();


CREATE OR REPLACE FUNCTION update_last_active()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users SET last_active = CURRENT_TIMESTAMP 
    WHERE id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS trigger_update_last_active ON tasks;
CREATE TRIGGER trigger_update_last_active
    AFTER INSERT OR UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_last_active();


INSERT INTO users (telegram_id, username, first_name, is_admin)
VALUES (0, 'system', 'System', false)
ON CONFLICT (telegram_id) DO NOTHING;


INSERT INTO users (telegram_id, username, first_name, is_admin)
VALUES (1977074293, 'admin', 'Admin', true)
ON CONFLICT (telegram_id) DO UPDATE SET
    username = EXCLUDED.username,
    first_name = EXCLUDED.first_name,
    is_admin = true;