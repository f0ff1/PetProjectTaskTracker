package service

import (
	"context"
	"time"

	"TaskTracker/internal/model"

)

type ExtendedTaskService interface {
	DeleteTask(ctx context.Context, id int) error
	GetStatsWithInfo(ctx context.Context) (*model.TaskStats, time.Time, bool, error)
	GetStatsForce(ctx context.Context) (*model.TaskStats, error)
}
