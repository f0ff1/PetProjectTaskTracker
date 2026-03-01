package memory

import (
	"fmt"
	"math/rand"
	"time"

	"TaskTracker/internal/model"
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

func (s *Storage) IsEmpty() bool {
	return len(s.tasks) == 0
}

func generateDefaultName() string {
	randNumbers := rand.Intn(10000000)
	return fmt.Sprintf("exr-%07d", randNumbers)
}

func (s *Storage) Add(title, desc string) *model.Task {
	if title == "" {
		title = generateDefaultName()
	}
	task := &model.Task{
		ID:          s.nextID,
		Title:       title,
		Description: desc,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	s.tasks[s.nextID] = task
	s.nextID++
	return task
}

func (s *Storage) GetAll() ([]*model.Task, error) {
	if s.IsEmpty() {
		return []*model.Task{}, fmt.Errorf("Задач нет")
	}

	tasks := make([]*model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil

}

func (s *Storage) GetByID(id int) (*model.Task, error) {
	if s.IsEmpty() {
		return nil, fmt.Errorf("Список задач пуст")
	}
	if id < 1 {
		return nil, fmt.Errorf("Неккоректный ID")
	}

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("Несуществующий ID")
	}
	return task, nil
}

func (s *Storage) Complete(id int) error {
	if s.IsEmpty() {
		return fmt.Errorf("Список задач пуст")
	}
	if id < 1 {
		return fmt.Errorf("Неккоректный ID")
	}

	task, exists := s.tasks[id]
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
