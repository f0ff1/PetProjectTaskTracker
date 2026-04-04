package service

import (
	"context"
	"time"

	"TaskTracker/internal/model"
)

type ExtendedTaskService interface {
	DeleteTask(ctx context.Context, userID int, id int) error

	GetStats(ctx context.Context, userID int) (*model.TaskStats, error)
	GetStatsForce(ctx context.Context, userID int) (*model.TaskStats, error)

	GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*model.User, error)
	GetUserByID(ctx context.Context, userID int) (*model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)

	GetAllTasksForAdmin(ctx context.Context) ([]*model.Task, error)

	GetStatsWithInfo(ctx context.Context, userID int) (*model.TaskStats, time.Time, bool, error)

	GetTasksForReminder(ctx context.Context) ([]*model.Task, error)
	MarkReminderSent(ctx context.Context, taskID int) error
}
