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

func createTask(title, desc string, tags []string) *model.Task {
	if title == "" {
		title = generateDefaultName()
	}
	return &model.Task{
		Title:       title,
		Description: desc,
		Tags:        tags,
	}
}

func generateDefaultName() string {
	randNumbers := rand.Intn(10000000)
	return fmt.Sprintf("def-name-exr-%07d", randNumbers)
}

func (s *TaskService) AddTask(ctx context.Context, title, desc string, tags []string) (*model.Task, error) {
	task := createTask(title, desc, tags)
	return s.repo.Add(ctx, task)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) GetTaskById(ctx context.Context, id int) (*model.Task, error) {
	if id < 1 {
		return nil, fmt.Errorf("Ошибка сервиса: %w", customError.ErrIdNotExists)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) GetTasksByTag(ctx context.Context, tag string) ([]*model.Task, error) {
	if tag == "" {
		return nil, fmt.Errorf("Ошибка сервиса: %w", customError.ErrWrongTag)
	}
	return s.repo.GetByTag(ctx, tag)
}

func (s *TaskService) CompleteTask(ctx context.Context, id int) (*model.Task, error) {
	if id < 1 {
		return nil, fmt.Errorf("Ошибка сервиса: %w", customError.ErrIdNotExists)
	}
	return s.repo.Complete(ctx, id)
}

type PostgresTaskService struct {
	*TaskService
	repo repository.PostgresRepository
}

func NewPostgresTaskService(repo repository.PostgresRepository) *PostgresTaskService {
	return &PostgresTaskService{TaskService: NewTaskService(repo), repo: repo}
}

func (s *PostgresTaskService) DeleteTask(ctx context.Context, id int) error {
	if id < 1 {
		return fmt.Errorf("Ошибка сервиса: %w", customError.ErrIdNotExists)
	}
	return s.repo.DeleteByID(ctx, id)
}

func (s *PostgresTaskService) GetStatsWithInfo(ctx context.Context) (*model.TaskStats, time.Time, bool, error) {
	return s.repo.GetStatsWithInfo(ctx)
}
