package memory

import (
	"sync"
	"time"

	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
)

type InMemoryRepo struct {
	mu     sync.RWMutex
	tasks  map[int]*model.Task
	nextID int
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		mu:     sync.RWMutex{},
		tasks:  make(map[int]*model.Task),
		nextID: 1,
	}
}

func (iM *InMemoryRepo) Add(title, desc string, tags []string) (*model.Task, error) {
	iM.mu.Lock()
	defer iM.mu.Unlock()

	task := &model.Task{
		ID:          iM.nextID,
		Title:       title,
		Description: desc,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        tags,
	}

	iM.tasks[iM.nextID] = task
	iM.nextID++
	return task, nil
}

func (iM *InMemoryRepo) GetAll() ([]*model.Task, error) {
	iM.mu.RLock()
	defer iM.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(iM.tasks))
	for _, task := range iM.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil

}

func (iM *InMemoryRepo) GetByID(id int) (*model.Task, error) {
	iM.mu.RLock()
	defer iM.mu.RUnlock()

	task, exists := iM.tasks[id]
	if !exists {
		return nil, myErrors.ErrIdNotExists
	}
	return task, nil
}

func (iM *InMemoryRepo) GetByTag(tag string) ([]*model.Task, error) {
	iM.mu.RLock()
	defer iM.mu.RUnlock()

	taggetTasks := make([]*model.Task, 0)
	for _, task := range iM.tasks {
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

func (iM *InMemoryRepo) Complete(id int) (*model.Task, error) {
	iM.mu.Lock()
	defer iM.mu.Unlock()

	task, exists := iM.tasks[id]
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

func (iM *InMemoryRepo) DeleteByID(id int) error {
	iM.mu.Lock()
	defer iM.mu.Unlock()

	_, exists := iM.tasks[id]
	if !exists {
		return myErrors.ErrIdNotExists
	}

	delete(iM.tasks, id)
	return nil
}
