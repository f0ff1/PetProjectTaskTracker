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

func (s *Storage) Add(title, desc string, tags []string) *model.Task {
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
		Tags:        tags,
	}

	s.tasks[s.nextID] = task
	fmt.Printf("Задаче: %s добавили тэги: %s", title, tags)
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

func (s *Storage) GetByTag(tag string) ([]*model.Task, error) {
	if s.IsEmpty() {
		return []*model.Task{}, fmt.Errorf("Задач не обнаружено")
	}
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

	if len(taggetTasks) < 1 {
		return taggetTasks, fmt.Errorf("Задач с таким тегом не существует")
	}
	return taggetTasks, nil

}

func (s *Storage) Complete(id int) error {
	task, err := s.GetByID(id)
	if err != nil {
		return fmt.Errorf("Ошибка: %w", err)
	}

	if task.Completed {
		return fmt.Errorf("Задача уже выполнена")
	}

	completeTime := time.Now()
	task.Completed = true
	task.CompletedAt = &completeTime
	return nil

}
