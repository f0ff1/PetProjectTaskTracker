package repository

import (
	"context"

	"TaskTracker/internal/model"
)

type Repository interface {
	Add(ctx context.Context, task *model.Task) (*model.Task, error)
	GetAll(ctx context.Context) ([]*model.Task, error)
	GetByID(ctx context.Context, id int) (*model.Task, error)
	Complete(ctx context.Context, id int) (*model.Task, error)
	GetByTag(ctx context.Context, tag string) ([]*model.Task, error)
}

type PostgresRepository interface {
	Repository
	DeleteByID(ctx context.Context, id int) error
	GetStats(ctx context.Context) (*model.TaskStats, error)
}
