package task

import (
	"fmt"
	"time"

	"TaskTracker/model"

)

type Manager struct {
	tasks  map[int]*model.Task
	nextID int
}

func NewManager() *Manager {
	return &Manager{
		tasks:  make(map[int]*model.Task),
		nextID: 1,
	}
}

func (m *Manager) IsEmpty() bool {
	return len(m.tasks) == 0
}

func (m *Manager) Add(title, desc string) *model.Task {
	task := &model.Task{
		ID:          m.nextID,
		Title:       title,
		Description: desc,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	m.tasks[m.nextID] = task
	m.nextID++
	return task
}

func (m *Manager) GetAll() ([]*model.Task, error) {
	if m.IsEmpty() {
		return []*model.Task{}, fmt.Errorf("Задач нет")
	}

	tasks := make([]*model.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil

}

func (m *Manager) GetByID(id int) (*model.Task, error) {
	if m.IsEmpty() {
		return nil, fmt.Errorf("Список задач пуст")
	}
	if id < 1 {
		return nil, fmt.Errorf("Неккоректный ID")
	}

	task, exists := m.tasks[id]
	if !exists {
		return nil, fmt.Errorf("Несуществующий ID")
	}
	return task, nil
}

func (m *Manager) Complete(id int) error {
	if m.IsEmpty() {
		return fmt.Errorf("Список задач пуст")
	}
	if id < 1 {
		return fmt.Errorf("Неккоректный ID")
	}

	task, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("Несуществующий ID")
	}

	if task.Completed {
		return fmt.Errorf("Задача уже выполнена")
	}

	completeTime := time.Now()
	task.Completed = true
	task.CompletedAt = &completeTime
	return nil

}
