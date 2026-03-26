package service

import (
	"fmt"
	"math/rand"

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

func generateDefaultName() string {
	randNumbers := rand.Intn(10000000)
	return fmt.Sprintf("def-name-exr-%07d", randNumbers)
}

func (s *TaskService) AddTask(title, desc string, tags []string) (*model.Task, error) {
	if title == "" {
		title = generateDefaultName()
	}
	return s.repo.Add(title, desc, tags)
}

func (s *TaskService) GetAllTasks() ([]*model.Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) GetTaskById(id int) (*model.Task, error) {
	if id < 1 {
		return nil, customError.ErrIdNotExists
	}
	return s.repo.GetByID(id)
}

func (s *TaskService) GetTasksByTag(tag string) ([]*model.Task, error) {
	if tag == "" {
		return nil, customError.ErrWrongTag
	}
	return s.repo.GetByTag(tag)
}

func (s *TaskService) CompleteTask(id int) (*model.Task, error) {
	if id < 1 {
		return nil, customError.ErrIdNotExists
	}
	return s.repo.Complete(id)
}

func (s *TaskService) DeleteTask(id int) error {
	if id < 1 {
		return customError.ErrIdNotExists
	}
	return s.repo.DeleteByID(id)
}

func (s *TaskService) GetStats() ([]string,error) {
	return s.repo.GetStats()
}
