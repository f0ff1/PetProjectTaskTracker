package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
)

type DataBaseRepo struct {
	dbPool *pgxpool.Pool
	cache  *statsCache
}

func NewPostgresRepo(connStr string) (*DataBaseRepo, error) {
	if err := runMigrations(connStr); err != nil {
		log.Printf("Ошибка миграций: %v", err)
	}
	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к БД: %w | %w", myErrors.ErrCantConnectToDB, err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("Ошибка при проверке подключения к БД: %w | %w", myErrors.ErrPingEx, err)
	}

	return &DataBaseRepo{
		dbPool: dbPool,
		cache:  newStatsCache(),
	}, nil
}

func (s *DataBaseRepo) Close() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
}

// calculateReminderTime вычисляет время, когда должно быть отправлено напоминание
func calculateReminderTime(dueDate time.Time, reminderOffset string) (time.Time, error) {
	reminderOffset = strings.TrimSpace(strings.ToLower(reminderOffset))
	if len(reminderOffset) < 2 {
		return time.Time{}, fmt.Errorf("invalid reminder offset format: %s", reminderOffset)
	}

	lastChar := reminderOffset[len(reminderOffset)-1]
	numStr := reminderOffset[:len(reminderOffset)-1]

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid reminder offset number: %s", numStr)
	}

	var duration time.Duration
	switch lastChar {
	case 'm':
		duration = time.Duration(num) * time.Minute
	case 'h':
		duration = time.Duration(num) * time.Hour
	case 'd':
		duration = time.Duration(num) * 24 * time.Hour
	default:
		return time.Time{}, fmt.Errorf("invalid reminder offset unit: %c", lastChar)
	}

	result := dueDate.Add(-duration)
	log.Printf("[REMINDER CALC] DueDate=%s, Offset=%s(num=%d, unit=%c, duration=%v) => ReminderTime=%s",
		dueDate.Format("02.01.2006 15:04:05"),
		reminderOffset,
		num, lastChar, duration,
		result.Format("02.01.2006 15:04:05"))

	return result, nil
}

func (s *DataBaseRepo) Add(ctx context.Context, userID int, task *model.Task) (*model.Task, error) {
	taskTitle := task.Title
	taskDesc := task.Description
	taskTags := task.Tags
	taskDueDate := task.DueDate

	// Если reminderOffset пустой, вставляем NULL, иначе строку
	var taskReminderOffset *string
	if task.ReminderOffset != nil && *task.ReminderOffset != "" {
		taskReminderOffset = task.ReminderOffset
	}

	task.CompletedAt = nil
	task.CreatedAt = time.Now()

	log.Printf("[DB ADD] Saving task: Title=%s, DueDate=%s, ReminderOffset=%s",
		taskTitle,
		func() string {
			if taskDueDate != nil {
				return taskDueDate.Format("02.01.2006 15:04:05")
			}
			return "nil"
		}(),
		func() string {
			if taskReminderOffset != nil {
				return *taskReminderOffset
			}
			return "nil"
		}())

	addQuery := `INSERT INTO tasks (user_id, title, description, tags, due_date, reminder_offset, created_at)
	VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7)
	RETURNING id, user_task_id, created_at`

	err := s.dbPool.QueryRow(ctx, addQuery, userID, taskTitle, taskDesc, taskTags, taskDueDate, taskReminderOffset, task.CreatedAt).Scan(
		&task.ID,
		&task.UserTaskID,
		&task.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Ошибка во время добавления задачи: %w | %w", myErrors.ErrCantSaveTaskToDB, err)
	}

	log.Printf("[DB ADD] Task saved successfully: ID=%d, UserTaskID=%d, ReminderOffset=%s",
		task.ID, task.UserTaskID,
		func() string {
			if taskReminderOffset != nil {
				return *taskReminderOffset
			}
			return "nil"
		}())

	task.UserID = userID
	log.Printf("Title bytes: % x", []byte(task.Title))
	log.Printf("Title valid UTF-8: %v", utf8.ValidString(task.Title))

	return task, nil
}

func (s *DataBaseRepo) DeleteByID(ctx context.Context, userID int, taskID int) error {
	deleteQuery := `DELETE FROM tasks where user_task_id = $1 and user_id = $2`
	row, err := s.dbPool.Exec(ctx, deleteQuery, taskID, userID)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении задачи: %w | %w", myErrors.ErrCantDeleteTask, err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("Ошибка при удалении задачи: %w | %w", myErrors.ErrIdNotExists, err)
	}

	return nil
}

func (s *DataBaseRepo) GetAllTasksByUser(ctx context.Context, userID int) ([]*model.Task, error) {
	getAllQuery := `SELECT id, user_id, user_task_id, title, description, completed, created_at, completed_at, tags, due_date, reminder_sent, reminder_offset FROM tasks WHERE user_id = $1 ORDER BY id`

	rows, err := s.dbPool.Query(ctx, getAllQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}
	defer rows.Close()

	var tasks []*model.Task

	for rows.Next() {
		task := &model.Task{}

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.UserTaskID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.Tags,
			&task.DueDate,
			&task.ReminderSent,
			&task.ReminderOffset,
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

func (s *DataBaseRepo) GetByID(ctx context.Context, userID int, taskID int) (*model.Task, error) {
	getByIdQuery := `select id, user_id, user_task_id, title, description, completed, created_at, completed_at, tags, due_date, reminder_sent, reminder_offset from tasks where user_task_id = $1 and user_id = $2`

	var task model.Task

	err := s.dbPool.QueryRow(ctx, getByIdQuery, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.UserTaskID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.Tags,
		&task.DueDate,
		&task.ReminderSent,
		&task.ReminderOffset,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrTableIsEmpty, err)
	}
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}
	return &task, nil
}

func (s *DataBaseRepo) GetByTag(ctx context.Context, userID int, tag string) ([]*model.Task, error) {
	getByTagQuery := `select id, user_id, user_task_id, title, description, completed, created_at, completed_at, tags, due_date, reminder_sent, reminder_offset from tasks where tags @> array[$1] and user_id = $2`

	var tasksWithTag []*model.Task
	rows, err := s.dbPool.Query(ctx, getByTagQuery, tag, userID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.UserTaskID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.Tags,
			&task.DueDate,
			&task.ReminderSent,
			&task.ReminderOffset,
		)

		if err != nil {
			return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
		}

		tasksWithTag = append(tasksWithTag, &task)
	}

	return tasksWithTag, nil
}

func (s *DataBaseRepo) Complete(ctx context.Context, userID int, taskID int) (*model.Task, error) {
	completeTime := time.Now()
	completeQuery := `update tasks set completed = true, completed_at = $1 where user_task_id = $2 and user_id = $3 RETURNING *`

	var task model.Task
	err := s.dbPool.QueryRow(ctx, completeQuery, completeTime, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.UserTaskID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.Tags,
		&task.DueDate,
		&task.ReminderSent,
		&task.ReminderOffset,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrTableIsEmpty, err)
	}
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задач: %w | %w", myErrors.ErrCantReadTable, err)
	}
	return &task, nil
}

func (s *DataBaseRepo) GetAllTasksForAdmin(ctx context.Context) ([]*model.Task, error) {
	query := `
        SELECT id, user_id, user_task_id, title, description, completed, created_at, completed_at, tags, due_date, reminder_sent, reminder_offset
        FROM tasks
        ORDER BY user_id, id
    `

	rows, err := s.dbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения задач: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.UserTaskID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.Tags,
			&task.DueDate,
			&task.ReminderSent,
			&task.ReminderOffset,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *DataBaseRepo) GetTasksForReminder(ctx context.Context) ([]*model.Task, error) {
	query := `SELECT id, user_id, user_task_id, title, description, completed, created_at, completed_at, tags, due_date, reminder_sent, reminder_offset from tasks WHERE due_date IS NOT NULL
	AND reminder_sent = false
	AND completed = false
	AND reminder_offset IS NOT NULL
	AND reminder_offset != ''`

	rows, err := s.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []*model.Task
	var pendingTasks []*model.Task

	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.UserTaskID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
			&task.Tags,
			&task.DueDate,
			&task.ReminderSent,
			&task.ReminderOffset,
		)
		if err != nil {
			return nil, err
		}

		pendingTasks = append(pendingTasks, task)
	}

	log.Printf("[REMINDER DB] Found %d pending tasks from DB", len(pendingTasks))

	// Filter tasks based on reminder time in Go to handle timezone correctly
	now := time.Now()

	// Get local timezone location
	loc := time.Local
	log.Printf("[REMINDER DB] Local timezone: %s", loc.String())
	log.Printf("[REMINDER DB] Current time: %s (Unix: %d)", now.Format("02.01.2006 15:04:05"), now.Unix())

	for _, task := range pendingTasks {
		if task.DueDate == nil || task.ReminderOffset == nil {
			log.Printf("[REMINDER DB] Task %d (ID: %d): Skipping - missing DueDate or ReminderOffset", task.UserTaskID, task.ID)
			continue
		}

		// Ensure DueDate has timezone info for proper comparison
		dueDateWithTz := task.DueDate
		if dueDateWithTz.Location() == time.UTC || dueDateWithTz.Location().String() == "UTC" {
			// If time is in UTC, convert to local time
			dueDateWithTz = dueDateWithTz.In(time.Local)
			log.Printf("[REMINDER DB] Task %d (ID: %d): Converted DueDate from UTC to local: %s",
				task.UserTaskID, task.ID, dueDateWithTz.Format("02.01.2006 15:04:05"))
		} else if dueDateWithTz.Location() == nil || dueDateWithTz.Location().String() == "" {
			// If timezone is unknown, treat as local time
			log.Printf("[REMINDER DB] Task %d (ID: %d): DueDate has no timezone info, treating as local", task.UserTaskID, task.ID)
		}

		reminderTime, err := calculateReminderTime(dueDateWithTz, *task.ReminderOffset)
		if err != nil {
			log.Printf("[REMINDER DB] Task %d (ID: %d): Error calculating reminder time: %v", task.UserTaskID, task.ID, err)
			continue
		}

		log.Printf("[REMINDER DB] Task %d (ID: %d): DueDate=%s (Unix:%d), ReminderOffset=%s, ReminderTime=%s (Unix:%d), Now=%s (Unix:%d)",
			task.UserTaskID, task.ID,
			dueDateWithTz.Format("02.01.2006 15:04:05"), dueDateWithTz.Unix(),
			*task.ReminderOffset,
			reminderTime.Format("02.01.2006 15:04:05"), reminderTime.Unix(),
			now.Format("02.01.2006 15:04:05"), now.Unix())

		if now.After(reminderTime) || now.Equal(reminderTime) {
			log.Printf("[REMINDER DB] Task %d (ID: %d): MATCH! Adding to send list (now=%d >= reminderTime=%d)",
				task.UserTaskID, task.ID, now.Unix(), reminderTime.Unix())
			tasks = append(tasks, task)
		} else {
			log.Printf("[REMINDER DB] Task %d (ID: %d): NOT YET (now=%d < reminderTime=%d, diff=%d seconds)",
				task.UserTaskID, task.ID, now.Unix(), reminderTime.Unix(), reminderTime.Unix()-now.Unix())
		}
	}

	log.Printf("[REMINDER DB] Filtered to %d tasks ready for reminder", len(tasks))
	return tasks, nil
}

func (s *DataBaseRepo) MarkReminderSent(ctx context.Context, taskID int) error {
	query := `UPDATE tasks SET reminder_sent = true where id = $1`
	_, err := s.dbPool.Exec(ctx, query, taskID)
	return err
}
