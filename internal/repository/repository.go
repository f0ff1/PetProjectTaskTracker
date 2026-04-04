package repository

import (
	"context"
	"time"

	"TaskTracker/internal/model"
)

// База
type Repository interface {
	Add(ctx context.Context, userID int, task *model.Task) (*model.Task, error)
	GetAllTasksByUser(ctx context.Context, userID int) ([]*model.Task, error)
	GetByID(ctx context.Context, userID int, taskID int) (*model.Task, error)
	Complete(ctx context.Context, userID int, taskID int) (*model.Task, error)
	GetByTag(ctx context.Context, userID int, tag string) ([]*model.Task, error)
}

// Расширение для Postgres
type PostgresRepository interface {
	Repository
	DeleteByID(ctx context.Context, userID int, taskID int) error
	GetAllTasksForAdmin(ctx context.Context) ([]*model.Task, error)
	GetUserByID(ctx context.Context, userID int) (*model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
	GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetStatsWithInfo(ctx context.Context, userID int) (*model.TaskStats, time.Time, bool, error)
	GetStatsWithRefresh(ctx context.Context, userID int, forceRefresh bool) (*model.TaskStats, error)
}
