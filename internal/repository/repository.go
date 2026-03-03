package repository

import "TaskTracker/internal/model"

type Repository interface {
	Add(title, description string, tags []string) (*model.Task, error)
	GetAll() ([]*model.Task, error)
	GetByID(id int) (*model.Task, error)
	Complete(id int) error
	GetByTag(tag string) ([]*model.Task, error)
}
