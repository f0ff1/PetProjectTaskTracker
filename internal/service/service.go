package service

import (
	customError "TaskTracker/errors"
	"TaskTracker/internal/model"
	"TaskTracker/internal/repository"
	"fmt"
	"math/rand"
)

type TaskService struct {
	repo repository.Repository
}

func NewNewTaskService(repo repository.Repository) *TaskService {
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
		return nil, customError.ErrWrongTypeID
	}
	return s.repo.Complete(id)
}
