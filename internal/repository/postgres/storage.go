package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
)

type PostgresRepo struct {
	dbPool *pgxpool.Pool
}

func NewPostgresRepo(connStr string) (*PostgresRepo, error) {
	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к БД: %w | %w", myErrors.ErrCantConnectToDB, err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("Ошибка при проверке подключения к БД: %w | %w", myErrors.ErrPingEx, err)
	}

	if err := createTableIfNotExists(dbPool); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("Ошибка при создании таблицы: %w | %w", myErrors.ErrCantCreateDBTable, err)
	}

	return &PostgresRepo{dbPool: dbPool}, nil
}

func createTableIfNotExists(pool *pgxpool.Pool) error {
	createTableSQL := ` 
	CREATE TABLE IF NOT EXISTS tasks  (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    tags TEXT[]
	);`
	_, err := pool.Exec(context.Background(), createTableSQL)
	return err
}

func (s *PostgresRepo) Close() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
}

func (s *PostgresRepo) Add(ctx context.Context, task *model.Task) (*model.Task, error) {
	taskTitle := task.Title
	taskDesc := task.Description
	taskTags := task.Tags
	task.CompletedAt = nil
	addQuery := `INSERT INTO tasks (title, description, tags)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	err := s.dbPool.QueryRow(ctx, addQuery, taskTitle, taskDesc, taskTags).Scan(
		&task.ID,
		&task.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Ошибка во время добавления задачи: %w | %w", myErrors.ErrCantSaveTaskToDB, err)
	}

	return task, nil
}

func (s *PostgresRepo) DeleteByID(ctx context.Context, id int) error {
	deleteQuery := `DELETE FROM tasks where id = $1`
	row, err := s.dbPool.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении задачи: %w | %w", myErrors.ErrCantDeleteTask, err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("Ошибка при удалении задачи: %w | %w", myErrors.ErrIdNotExists, err)
	}

	return nil
}

func (s *PostgresRepo) GetAll(ctx context.Context) ([]*model.Task, error) {
	getAllQuery := `SELECT * FROM tasks ORDER BY id`

	rows, err := s.dbPool.Query(ctx, getAllQuery)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}
	defer rows.Close()

	var tasks []*model.Task

	for rows.Next() {
		task := &model.Task{}

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.Tags,
		)

		if err != nil {
			return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}

	return tasks, nil
}

func (s *PostgresRepo) GetByID(ctx context.Context, id int) (*model.Task, error) {
	getByIdQuery := `select * from tasks where id = $1`

	var task model.Task

	err := s.dbPool.QueryRow(ctx, getByIdQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.Tags,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrTableIsEmpty, err)
	}
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}
	return &task, nil
}

func (s *PostgresRepo) GetByTag(ctx context.Context, tag string) ([]*model.Task, error) {
	getByTagQuery := `select * from tasks where tags @> array[$1]`

	var tasksWithTag []*model.Task
	rows, err := s.dbPool.Query(ctx, getByTagQuery, tag)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.Tags,
		)

		if err != nil {
			return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
		}

		tasksWithTag = append(tasksWithTag, &task)
	}

	return tasksWithTag, nil
}

func (s *PostgresRepo) Complete(ctx context.Context, id int) (*model.Task, error) {
	completeQuery := `update tasks set completed = true, completed_at = CURRENT_TIMESTAMP where id = $1 RETURNING *`

	var task model.Task
	err := s.dbPool.QueryRow(ctx, completeQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.Tags,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrTableIsEmpty, err)
	}
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}
	return &task, nil
}

func (s *PostgresRepo) GetStats(ctx context.Context) (*model.TaskStats, error) {
	// getStatsQuery := `SELECT tag, COUNT(*) as usage_count FROM (SELECT unnest(tags) as tag FROM tasks) WHERE tag IS NOT NULL GROUP BY tag ORDER BY usage_count DESC LIMIT 3`

	// rows, err := s.dbPool.Query(ctx, getStatsQuery)
	// if err != nil {
	// 	return nil, fmt.Errorf("Ошибка при чтении статистики: %w | %w", myErrors.ErrCantReadTable, err)
	// }
	// defer rows.Close()

	// var dataStats []string

	// for rows.Next() {
	// 	var tag string
	// 	var usageCount int

	// 	if err := rows.Scan(&tag, &usageCount); err != nil {
	// 		return nil, err
	// 	}
	// 	dataStats = append(dataStats, fmt.Sprintf("Тег: %s | Количество: %d", tag, usageCount))
	// }

	// return dataStats, nil

	query := `
        -- 1. Базовая статистика по задачам
        WITH task_counts AS (
            SELECT 
                COUNT(*) as total,
                COUNT(CASE WHEN completed = true THEN 1 END) as completed,
                COUNT(CASE WHEN completed = false THEN 1 END) as pending,
                ROUND(
                    COUNT(CASE WHEN completed = true THEN 1 END)::numeric / 
                    NULLIF(COUNT(*), 0) * 100, 
                    2
                ) as completion_rate
            FROM tasks
        ),
        
        -- 2. Статистика по тегам
        tag_stats AS (
            SELECT 
                tag, 
                COUNT(*) as usage_count
            FROM (
                SELECT unnest(tags) as tag 
                FROM tasks 
                WHERE tags IS NOT NULL AND array_length(tags, 1) > 0
            ) t
            WHERE tag IS NOT NULL AND tag != ''
            GROUP BY tag
        ),
        
        -- 3. Общее количество уникальных тегов (отдельный CTE)
        unique_tags_count AS (
            SELECT COUNT(*) as total_unique_tags
            FROM tag_stats
        ),
        
        -- 4. Топ-3 тега
        top_tags AS (
            SELECT 
                tag, 
                usage_count
            FROM tag_stats
            ORDER BY usage_count DESC
            LIMIT 3
        ),
        
        -- 5. Ежедневная статистика (последние 30 дней)
        daily_stats AS (
            SELECT 
                DATE(created_at) as day,
                COUNT(*) as total_count,
                COUNT(CASE WHEN completed = true THEN 1 END) as completed_count
            FROM tasks
            WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
            GROUP BY DATE(created_at)
            ORDER BY day DESC
        ),
        
        -- 6. Среднее время выполнения задачи
        completion_time AS (
            SELECT 
                ROUND(AVG(EXTRACT(EPOCH FROM (completed_at - created_at))/3600)::numeric, 2) as avg_hours
            FROM tasks
            WHERE completed = true 
                AND completed_at IS NOT NULL
                AND completed_at > created_at
        ),
        
        -- 7. Распределение по часам
        hourly_distribution AS (
            SELECT 
                EXTRACT(HOUR FROM created_at) as hour,
                COUNT(*) as count
            FROM tasks
            GROUP BY EXTRACT(HOUR FROM created_at)
            ORDER BY hour
        ),
        
        -- 8. Распределение по дням недели
        weekday_distribution AS (
            SELECT 
                TRIM(TO_CHAR(created_at, 'Day')) as weekday,
                COUNT(*) as count
            FROM tasks
            GROUP BY TRIM(TO_CHAR(created_at, 'Day'))
            ORDER BY 
                CASE TRIM(TO_CHAR(created_at, 'Day'))
                    WHEN 'Monday'    THEN 1
                    WHEN 'Tuesday'   THEN 2
                    WHEN 'Wednesday' THEN 3
                    WHEN 'Thursday'  THEN 4
                    WHEN 'Friday'    THEN 5
                    WHEN 'Saturday'  THEN 6
                    WHEN 'Sunday'    THEN 7
                END
        )
        
        -- 9. Финальная выборка
        SELECT 
            -- Основная статистика
            tc.total,
            tc.completed,
            tc.pending,
            tc.completion_rate,
            
            -- Уникальные теги
            COALESCE(utc.total_unique_tags, 0) as total_unique_tags,
            
            -- Топ-теги в JSON
            COALESCE(
                (SELECT json_agg(json_build_object('tag', tt.tag, 'usage_count', tt.usage_count))
                 FROM top_tags tt),
                '[]'::json
            ) as top_tags,
            
            -- Ежедневная статистика в JSON
            COALESCE(
                (SELECT json_agg(json_build_object('date', ds.day, 'count', ds.total_count))
                 FROM daily_stats ds),
                '[]'::json
            ) as daily_stats,
            
            -- Статистика по завершенным задачам по дням
            COALESCE(
                (SELECT json_agg(json_build_object('date', ds.day, 'count', ds.completed_count))
                 FROM daily_stats ds
                 WHERE ds.completed_count > 0),
                '[]'::json
            ) as completion_by_day,
            
            -- Среднее время выполнения
            COALESCE(ct.avg_hours, 0) as avg_completion_time,
            
            -- Распределение по часам
            COALESCE(
                (SELECT json_agg(json_build_object('hour', hd.hour, 'count', hd.count))
                 FROM hourly_distribution hd),
                '[]'::json
            ) as hourly_stats,
            
            -- Распределение по дням недели
            COALESCE(
                (SELECT json_agg(json_build_object('weekday', wd.weekday, 'count', wd.count))
                 FROM weekday_distribution wd),
                '[]'::json
            ) as weekday_stats
            
        FROM task_counts tc
        CROSS JOIN unique_tags_count utc
        CROSS JOIN completion_time ct
    `

	var stats model.TaskStats
	var topTagsJSON, dailyStatsJSON, completionByDayJSON, hourlyStatsJSON, weekdayStatsJSON []byte

	err := s.dbPool.QueryRow(ctx, query).Scan(
		&stats.TotalTasks,
		&stats.CompletedTasks,
		&stats.PendingTasks,
		&stats.CompletionRate,
		&stats.TotalUniqueTags,
		&topTagsJSON,
		&dailyStatsJSON,
		&completionByDayJSON,
		&stats.AvgCompletionTime,
		&hourlyStatsJSON,
		&weekdayStatsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении статистики: %w | %w", myErrors.ErrCantReadTable, err)
	}

	if err := json.Unmarshal(topTagsJSON, &stats.TopTags); err != nil {
		slog.WarnContext(ctx, "Ошибка парсинга топ-тегов", "error", err)
		stats.TopTags = []model.TagStat{}
	}

	if err := json.Unmarshal(dailyStatsJSON, &stats.TasksByDay); err != nil {
		slog.WarnContext(ctx, "Ошибка парсинга дневной статистики", "error", err)
		stats.TasksByDay = []model.DailyStat{}
	}

	if err := json.Unmarshal(completionByDayJSON, &stats.CompletionByDay); err != nil {
		stats.CompletionByDay = []model.DailyStat{}
	}

	if err := json.Unmarshal(hourlyStatsJSON, &stats.TasksByHour); err != nil {
		stats.TasksByHour = []model.HourlyStat{}
	}

	if err := json.Unmarshal(weekdayStatsJSON, &stats.TasksByWeekday); err != nil {
		stats.TasksByWeekday = []model.WeekdayStat{}
	}

	// Находим самый продуктивный день
	if len(stats.TasksByDay) > 0 {
		maxCount := 0
		for _, day := range stats.TasksByDay {
			if day.Count > maxCount {
				maxCount = day.Count
				stats.MostProductiveDay = day.Date
			}
		}
	}
	return &stats, nil

}
