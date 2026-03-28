package service

import (
	"context"

	"TaskTracker/internal/model"
)

type ExtendedTaskService interface {
	DeleteTask(ctx context.Context, id int) error
	GetStats(ctx context.Context) (*model.TaskStats, error)
}
