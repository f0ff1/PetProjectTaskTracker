package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	customError "TaskTracker/errors"
	"TaskTracker/internal/model"
	"TaskTracker/internal/repository"
)

type TaskService struct {
	repo repository.Repository
}

func NewTaskService(repo repository.Repository) *TaskService {
	return &TaskService{repo: repo}
}

func createTask(title, desc string, tags []string, dueDate *time.Time, reminderOffset string) *model.Task {
	if title == "" {
		title = generateDefaultName()
	}

	var reminderOffsetPtr *string
	if reminderOffset != "" {
		reminderOffsetPtr = &reminderOffset
	}

	return &model.Task{
		Title:          title,
		Description:    desc,
		Tags:           tags,
		DueDate:        dueDate,
		ReminderOffset: reminderOffsetPtr,
	}
}

func generateDefaultName() string {
	randNumbers := rand.Intn(10000000)
	return fmt.Sprintf("def-name-exr-%07d", randNumbers)
}

func (s *TaskService) AddTask(ctx context.Context, userID int, title, desc string, tags []string) (*model.Task, error) {
	task := createTask(title, desc, tags, nil, "")
	return s.repo.Add(ctx, userID, task)
}

// AddTaskWithReminder добавляет задачу с датой и напоминанием
func (s *TaskService) AddTaskWithReminder(ctx context.Context, userID int, title, desc string, tags []string, dueDateStr *string, reminderOffset string) (*model.Task, error) {
	var dueDate *time.Time
	if dueDateStr != nil {
		parsedDate, err := time.ParseInLocation("02.01.2006 15:04", *dueDateStr, time.Local)
		if err != nil {
			return nil, fmt.Errorf("Ошибка сервиса: %w", customError.ErrWrongDate)
		}
		dueDate = &parsedDate
	}
	task := createTask(title, desc, tags, dueDate, reminderOffset)
	return s.repo.Add(ctx, userID, task)
}

func (s *TaskService) GetAllTasks(ctx context.Context, userID int) ([]*model.Task, error) {
	return s.repo.GetAllTasksByUser(ctx, userID)
}

func (s *TaskService) GetTaskById(ctx context.Context, userID, taskID int) (*model.Task, error) {
	if taskID < 1 {
		return nil, fmt.Errorf("Ошибка сервиса: %w", customError.ErrIdNotExists)
	}
	return s.repo.GetByID(ctx, userID, taskID)
}

func (s *TaskService) GetTasksByTag(ctx context.Context, userID int, tag string) ([]*model.Task, error) {
	if tag == "" {
		return nil, fmt.Errorf("Ошибка сервиса: %w", customError.ErrWrongTag)
	}
	return s.repo.GetByTag(ctx, userID, tag)
}

func (s *TaskService) CompleteTask(ctx context.Context, userID, taskID int) (*model.Task, error) {
	if taskID < 1 {
		return nil, fmt.Errorf("Ошибка сервиса: %w", customError.ErrIdNotExists)
	}
	return s.repo.Complete(ctx, userID, taskID)
}

/* Ебучее расширение для Postgres, чтобы не плодить кучу сервисов и не дублировать код */
/* Ебучее расширение для Postgres, чтобы не плодить кучу сервисов и не дублировать код */
/* Ебучее расширение для Postgres, чтобы не плодить кучу сервисов и не дублировать код */

type PostgresTaskService struct {
	*TaskService
	repo repository.PostgresRepository
}

func NewPostgresTaskService(repo repository.PostgresRepository) *PostgresTaskService {
	return &PostgresTaskService{TaskService: NewTaskService(repo), repo: repo}
}

func (s *PostgresTaskService) DeleteTask(ctx context.Context, userID, taskID int) error {
	if taskID < 1 {
		return fmt.Errorf("Ошибка сервиса: %w", customError.ErrIdNotExists)
	}
	return s.repo.DeleteByID(ctx, userID, taskID)
}

func (s *PostgresTaskService) GetStatsWithInfo(ctx context.Context, userID int) (*model.TaskStats, time.Time, bool, error) {
	return s.repo.GetStatsWithInfo(ctx, userID)
}

// GetStats получает статистику (с проверкой свежести)
func (s *PostgresTaskService) GetStats(ctx context.Context, userID int) (*model.TaskStats, error) {
	return s.repo.GetStatsWithRefresh(ctx, userID, false)
}

// GetStatsForce принудительно обновляет статистику
func (s *PostgresTaskService) GetStatsForce(ctx context.Context, userID int) (*model.TaskStats, error) {
	return s.repo.GetStatsWithRefresh(ctx, userID, true)
}

func (s *PostgresTaskService) GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*model.User, error) {
	return s.repo.GetOrCreateUser(ctx, telegramID, username, firstName, lastName)
}

func (s *PostgresTaskService) GetUserByID(ctx context.Context, userID int) (*model.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *PostgresTaskService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *PostgresTaskService) GetAllTasksForAdmin(ctx context.Context) ([]*model.Task, error) {
	return s.repo.GetAllTasksForAdmin(ctx)
}

func (s *PostgresTaskService) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	return s.repo.GetUserByTelegramID(ctx, telegramID)
}

func (s *PostgresTaskService) GetTasksForReminder(ctx context.Context) ([]*model.Task, error) {
	return s.repo.GetTasksForReminder(ctx)
}

func (s *PostgresTaskService) MarkReminderSent(ctx context.Context, taskID int) error {
	return s.repo.MarkReminderSent(ctx, taskID)
}
