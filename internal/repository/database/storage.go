package database

import (
	"context"
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/golang-migrate/migrate/v4"
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

func runMigrations(connStr string) error {
	migration, err := migrate.New(
		"file://../migrations",
		connStr,
	)

	if err != nil {
		return fmt.Errorf("%v : %w", myErrors.ErrCreateMigration, err)
	}

	defer migration.Close()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%v : %w", myErrors.ErrCantUseMigration, err)
	}
	return nil
}

func NewPostgresRepo(connStr string) (*DataBaseRepo, error) {
	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к БД: %w | %w", myErrors.ErrCantConnectToDB, err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("Ошибка при проверке подключения к БД: %w | %w", myErrors.ErrPingEx, err)
	}

	if err := runMigrations(connStr); err != nil {
		dbPool.Close()
		return nil, err
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

func (s *DataBaseRepo) Add(ctx context.Context, task *model.Task) (*model.Task, error) {
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

	log.Printf("Title bytes: % x", []byte(task.Title))
	log.Printf("Title valid UTF-8: %v", utf8.ValidString(task.Title))

	return task, nil
}

func (s *DataBaseRepo) DeleteByID(ctx context.Context, id int) error {
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

func (s *DataBaseRepo) GetAll(ctx context.Context) ([]*model.Task, error) {
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

func (s *DataBaseRepo) GetByID(ctx context.Context, id int) (*model.Task, error) {
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

func (s *DataBaseRepo) GetByTag(ctx context.Context, tag string) ([]*model.Task, error) {
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

func (s *DataBaseRepo) Complete(ctx context.Context, id int) (*model.Task, error) {
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
