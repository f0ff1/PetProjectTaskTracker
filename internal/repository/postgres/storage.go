package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"

)

type PostgresStorage struct {
	dbPool *pgxpool.Pool
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, myErrors.ErrCantConnectToDB
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, myErrors.ErrPingEx
	}

	if err := createTableIfNotExists(dbPool); err != nil {
		dbPool.Close()
		return nil, myErrors.ErrCantCreateDBTable
	}

	return &PostgresStorage{dbPool: dbPool}, nil
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

func (s *PostgresStorage) Close() {
	if s.dbPool != nil {
		s.dbPool.Close()
	}
}

func (s *PostgresStorage) Add(title, desc string, tags []string) (*model.Task, error) {
	addQuery := `INSERT INTO tasks (title, description, tags)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	var task model.Task
	task.Title = title
	task.Description = desc
	task.Completed = false
	task.CompletedAt = nil
	task.Tags = tags

	err := s.dbPool.QueryRow(context.Background(), addQuery, title, desc, tags).Scan(
		&task.ID,
		&task.CreatedAt,
	)
	if err != nil {
		return nil, myErrors.ErrCantSaveTaskToDB
	}

	return &task, nil

}

func (s *PostgresStorage) DeleteByID(id int) error {
	deleteQuery := `DELETE FROM tasks where id = $1`
	row, err := s.dbPool.Exec(context.Background(), deleteQuery, id)
	if err != nil {
		return err
	}

	if row.RowsAffected() == 0 {
		fmt.Println(myErrors.ErrIdNotExists)
	}

	return nil
}

func (s *PostgresStorage) GetAll() ([]*model.Task, error) {
	getAllQuery := `SELECT * FROM tasks ORDER BY id`

	rows, err := s.dbPool.Query(context.Background(), getAllQuery)
	if err != nil {
		return nil, myErrors.ErrCantReadTable
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
			return nil, myErrors.ErrCantReadTable
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, myErrors.ErrCantReadTable
	}

	return tasks, nil
}

func (s *PostgresStorage) GetByID(id int) (*model.Task, error) {
	getByIdQuery := `select * from tasks where id = $1`

	var task model.Task

	err := s.dbPool.QueryRow(context.Background(), getByIdQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.Tags,
	)

	if err == pgx.ErrNoRows {
		return nil, myErrors.ErrTableIsEmpty
	}
	if err != nil {
		return nil, myErrors.ErrCantReadTable
	}
	return &task, nil
}

func (s *PostgresStorage) GetByTag(tag string) ([]*model.Task, error) {
	getByTagQuery := `select * from tasks where tags @> array[$1]`

	var tasksWithTag []*model.Task
	rows, err := s.dbPool.Query(context.Background(), getByTagQuery, tag)
	if err != nil {
		return nil, myErrors.ErrCantReadTable
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
			return nil, myErrors.ErrCantReadTable
		}

		tasksWithTag = append(tasksWithTag, &task)
	}

	return tasksWithTag, nil
}

func (s *PostgresStorage) Complete(id int) (*model.Task, error) {
	completeQuery := `update tasks set completed = true, completed_at = CURRENT_TIMESTAMP where id = $1 RETURNING *`

	var task model.Task
	err := s.dbPool.QueryRow(context.Background(), completeQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.Tags,
	)
	if err == pgx.ErrNoRows {
		return nil, myErrors.ErrTableIsEmpty
	}
	if err != nil {
		return nil, myErrors.ErrCantReadTable
	}
	return &task, nil
}

func (s *PostgresStorage) GetStats() ([]string, error) {
	getStatsQuery := `SELECT tag, COUNT(*) as usage_count FROM (SELECT unnest(tags) as tag FROM tasks) WHERE tag IS NOT NULL GROUP BY tag ORDER BY usage_count DESC LIMIT 3`

	rows, err := s.dbPool.Query(context.Background(), getStatsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataStats []string

	for rows.Next() {
		var tag string
		var usageCount int

		if err := rows.Scan(&tag, &usageCount); err != nil {
			return nil, err
		}
		dataStats = append(dataStats, fmt.Sprintf("Тег: %s | Количество: %d", tag, usageCount))
	}

	return dataStats, nil

}
