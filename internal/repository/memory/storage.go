package memory

import (
	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
	"sync"
	"time"
)

type Storage struct {
	mu     sync.RWMutex
	tasks  map[int]*model.Task
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		tasks:  make(map[int]*model.Task),
		nextID: 1,
	}
}

func (s *Storage) Add(title, desc string, tags []string) (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil

}

func (s *Storage) GetByID(id int) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, myErrors.ErrIdNotExists
	}
	return task, nil
}

func (s *Storage) GetByTag(tag string) ([]*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	taggetTasks := make([]*model.Task, 0)
	for _, task := range s.tasks {
		taskTags := task.Tags
		tagsMap := make(map[string]bool)
		for _, t := range taskTags {
			tagsMap[t] = true
		}
		if tagsMap[tag] {
			taggetTasks = append(taggetTasks, task)
		}
	}

	return taggetTasks, nil

}

func (s *Storage) Complete(id int) (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
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
