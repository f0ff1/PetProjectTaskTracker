package sjson

import (
	"TaskTracker/internal/model"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type JSONData struct {
	NextID int                 `json:"next_id"`
	Tasks  map[int]*model.Task `json:"tasks"`
}

type JSONStorage struct {
	filePath string
	tasks    map[int]*model.Task
	nextID   int
	mu       sync.RWMutex
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	storage := &JSONStorage{
		filePath: path,
		tasks:    make(map[int]*model.Task),
		nextID:   1,
		mu:       sync.RWMutex{},
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		//Если файла нема нихуя
		if err := storage.saveToFile(); err != nil {
			return nil, fmt.Errorf("Не получилось создать твой ебаный файл йоу: %w", err)
		}
		return storage, nil
	}

	if err := storage.loadFromFile(); err != nil {
		return nil, fmt.Errorf("Не смог прочитать ебаные данные: %w", err)
	}
	return storage, nil
}

func (s *JSONStorage) saveToFile() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := JSONData{
		NextID: s.nextID,
		Tasks:  s.tasks,
	}

	file, err := os.Create(s.filePath)
	if err != nil {
		return fmt.Errorf("Не удалось создать/переписать файл по такому ебаному пути: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(data)
}

func (s *JSONStorage) loadFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		return fmt.Errorf("Не смог прочитать ебаные данные: %w", err)
	}
	defer file.Close()

	var data JSONData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	s.nextID = data.NextID
	s.tasks = data.Tasks
	return nil

}

func (s *JSONStorage) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tasks) == 0
}

func generateDefaultName() string {
	randNumbers := rand.Intn(10000000)
	return fmt.Sprintf("json-exr-%07d", randNumbers)
}

func (s *JSONStorage) Add(title, description string, tags []string) *model.Task {
	s.mu.Lock()

	if title == "" {
		title = generateDefaultName()
	}

	task := &model.Task{
		ID:          s.nextID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	s.tasks[s.nextID] = task
	s.nextID++
	s.mu.Unlock()

	if err := s.saveToFile(); err != nil {
		fmt.Printf("Не удалось сохранить ебаную таску в JSON %s\n", err)
	}

	return task

}

func (s *JSONStorage) GetAll() ([]*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.IsEmpty() {
		return []*model.Task{}, fmt.Errorf("Задач нема нихуя")
	}

	tasks := make([]*model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *JSONStorage) GetByID(id int) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.IsEmpty() {
		return nil, fmt.Errorf("Задач нема нихуя")
	}

	if id < 1 {
		return nil, fmt.Errorf("Не может быть ID < 1, тупица ебаная")
	}

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("Несуществующий ID")
	}
	return task, nil

}

func (s *JSONStorage) GetByTag(tag string) ([]*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

func (s *JSONStorage) Complete(id int) error {

	if s.IsEmpty() {
		return fmt.Errorf("Список задач пуст")
	}

	s.mu.Lock()
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
	s.mu.Unlock()

	if err := s.saveToFile(); err != nil {
		return fmt.Errorf("ошибка сохранения: %w", err)
	}
	return nil

}
