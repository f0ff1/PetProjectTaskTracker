package memory

import (
	"context"
	"fmt"
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

func (iM *InMemoryRepo) Add(ctx context.Context, userID int, task *model.Task) (*model.Task, error) {
	iM.mu.Lock()
	defer iM.mu.Unlock()

	task.ID = iM.nextID
	task.Completed = false
	task.CreatedAt = time.Now()
	task.CompletedAt = nil

	iM.tasks[iM.nextID] = task
	iM.nextID++
	return task, nil
}

func (iM *InMemoryRepo) GetAllTasksByUser(ctx context.Context, userID int) ([]*model.Task, error) {
	iM.mu.RLock()
	defer iM.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(iM.tasks))
	for _, task := range iM.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil

}

func (iM *InMemoryRepo) GetByID(ctx context.Context, userID int, taskID int) (*model.Task, error) {
	iM.mu.RLock()
	defer iM.mu.RUnlock()

	task, exists := iM.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("Задача с ID %d не найдена: %w", taskID, myErrors.ErrIdNotExists)
	}
	return task, nil
}

func (iM *InMemoryRepo) GetByTag(ctx context.Context, userID int, tag string) ([]*model.Task, error) {
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

func (iM *InMemoryRepo) Complete(ctx context.Context, userID int, taskID int) (*model.Task, error) {
	iM.mu.Lock()
	defer iM.mu.Unlock()

	task, exists := iM.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("Задача с ID %d не найдена: %w", taskID, myErrors.ErrIdNotExists)
	}

	if task.Completed {
		return nil, fmt.Errorf("Задача с ID %d уже завершена: %w", taskID, myErrors.ErrTaskAlredyComplete)
	}

	completeTime := time.Now()
	task.Completed = true
	task.CompletedAt = &completeTime
	return task, nil

}

func (iM *InMemoryRepo) DeleteByID(ctx context.Context, userID int, id int) error {
	iM.mu.Lock()
	defer iM.mu.Unlock()

	_, exists := iM.tasks[id]
	if !exists {
		return fmt.Errorf("Задача с ID %d не найдена: %w", id, myErrors.ErrIdNotExists)
	}

	delete(iM.tasks, id)
	return nil
}
