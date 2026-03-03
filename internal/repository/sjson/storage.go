package sjson

import (
	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
	"encoding/json"
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
			return nil, myErrors.ErrCantCreateJsonFile
		}
		return storage, nil
	}

	if err := storage.loadFromFile(); err != nil {
		return nil, myErrors.ErrCantReadJsonData
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
		return myErrors.ErrWrongPath
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
		return myErrors.ErrCantReadJsonData
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

func (s *JSONStorage) Add(title, description string, tags []string) (*model.Task, error) {
	s.mu.Lock()
	task := &model.Task{
		ID:          s.nextID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        tags,
	}

	s.tasks[s.nextID] = task
	s.nextID++
	s.mu.Unlock()

	if err := s.saveToFile(); err != nil {
		return nil, myErrors.ErrCantSaveTaskToJson
	}

	return task, nil

}

func (s *JSONStorage) GetAll() ([]*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *JSONStorage) GetByID(id int) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, myErrors.ErrIdNotExists
	}
	return task, nil

}

func (s *JSONStorage) GetByTag(tag string) ([]*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

func (s *JSONStorage) Complete(id int) (*model.Task, error) {

	task, err := s.GetByID(id)
	if err != nil {
		return nil, myErrors.ErrIdNotExists
	}

	if task.Completed {
		return nil, myErrors.ErrTaskAlredyComplete
	}

	s.mu.Lock()
	completeTime := time.Now()
	task.Completed = true
	task.CompletedAt = &completeTime
	s.mu.Unlock()

	if err := s.saveToFile(); err != nil {
		return nil, myErrors.ErrCantSaveTaskToJson
	}
	return task, nil

}
