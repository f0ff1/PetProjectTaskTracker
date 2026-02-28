package task

import "TaskTracker/model"

type Repository interface {
	Add(title, description string) *model.Task
	GetAll() ([]*model.Task, error)
	GetByID(id int) (*model.Task, error)
	Complete(id int) error
	IsEmpty() bool
}
