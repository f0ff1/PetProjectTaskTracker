package memory

import (
	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
	"time"
)

type Storage struct {
	tasks  map[int]*model.Task
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		tasks:  make(map[int]*model.Task),
		nextID: 1,
	}
}

func (s *Storage) Add(title, desc string, tags []string) (*model.Task, error) {

	task := &model.Task{
		ID:          s.nextID,
		Title:       title,
		Description: desc,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        tags,
	}

	s.tasks[s.nextID] = task
	s.nextID++
	return task, nil
}

func (s *Storage) GetAll() ([]*model.Task, error) {

	tasks := make([]*model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil

}

func (s *Storage) GetByID(id int) (*model.Task, error) {

	task, exists := s.tasks[id]
	if !exists {
		return nil, myErrors.ErrIdNotExists
	}
	return task, nil
}

func (s *Storage) GetByTag(tag string) ([]*model.Task, error) {

	taggetTasks := make([]*model.Task, 0)
	tasks, _ := s.GetAll()
	for _, task := range tasks {
		taskTags := task.Tags
		tagsMap := make(map[string]bool)
		for _, task := range taskTags {
			tagsMap[task] = true
		}
		if tagsMap[tag] {
			taggetTasks = append(taggetTasks, task)
		}

	}

	return taggetTasks, nil

}

func (s *Storage) Complete(id int) (*model.Task, error) {
	task, err := s.GetByID(id)
	if err != nil {
		return nil, myErrors.ErrIdNotExists
	}

	if task.Completed {
		return nil, myErrors.ErrTaskAlredyComplete
	}

	completeTime := time.Now()
	task.Completed = true
	task.CompletedAt = &completeTime
	return task, nil

}
