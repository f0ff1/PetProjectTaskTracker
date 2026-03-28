package sjson

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
)

type JSONData struct {
	NextID int                 `json:"next_id"`
	Tasks  map[int]*model.Task `json:"tasks"`
}

type JSONRepo struct {
	filePath string
	tasks    map[int]*model.Task
	nextID   int
	mu       sync.RWMutex
}

func NewJSONRepo(path string) (*JSONRepo, error) {
	storage := &JSONRepo{
		filePath: path,
		tasks:    make(map[int]*model.Task),
		nextID:   1,
		mu:       sync.RWMutex{},
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		//Если файла нема нихуя
		if err := storage.saveToFile(); err != nil {
			return nil, fmt.Errorf("Ошибка при создании файла JSON: %w | %w", myErrors.ErrCantCreateJsonFile, err)
		}
		return storage, nil
	}

	if err := storage.loadFromFile(); err != nil {
		return nil, fmt.Errorf("Ошибка при чтении файла JSON: %w | %w", myErrors.ErrCantReadJsonData, err)
	}
	return storage, nil
}

func (jR *JSONRepo) saveToFile() error {
	jR.mu.RLock()
	defer jR.mu.RUnlock()

	data := JSONData{
		NextID: jR.nextID,
		Tasks:  jR.tasks,
	}

	file, err := os.Create(jR.filePath)
	if err != nil {
		return fmt.Errorf("Ошибка при создании файла JSON: %w | %w", myErrors.ErrWrongPath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(data)
}

func (jR *JSONRepo) loadFromFile() error {
	jR.mu.Lock()
	defer jR.mu.Unlock()

	file, err := os.Open(jR.filePath)
	if err != nil {
		return fmt.Errorf("Ошибка при чтении файла JSON: %w | %w", myErrors.ErrCantReadJsonData, err)
	}
	defer file.Close()

	var data JSONData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("Ошибка при декодировании данных из файла JSON: %w | %w", myErrors.ErrCantReadJsonData, err)
	}

	jR.nextID = data.NextID
	jR.tasks = data.Tasks
	return nil

}

func (jR *JSONRepo) Add(ctx context.Context, task *model.Task) (*model.Task, error) {
	jR.mu.Lock()
	task.ID = jR.nextID
	task.Completed = false
	task.CreatedAt = time.Now()
	task.CompletedAt = nil

	jR.tasks[jR.nextID] = task
	jR.nextID++
	jR.mu.Unlock()

	if err := jR.saveToFile(); err != nil {
		return nil, fmt.Errorf("Ошибка добавления задачи: %w | %w", myErrors.ErrCantSaveTaskToJson, err)
	}

	return task, nil

}

func (jR *JSONRepo) GetAll(ctx context.Context) ([]*model.Task, error) {
	jR.mu.RLock()
	defer jR.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(jR.tasks))
	for _, task := range jR.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (jR *JSONRepo) GetByID(ctx context.Context, id int) (*model.Task, error) {
	jR.mu.RLock()
	defer jR.mu.RUnlock()

	task, exists := jR.tasks[id]
	if !exists {
		return nil, fmt.Errorf("Ошибка при чтении задачи: %w", myErrors.ErrIdNotExists)
	}
	return task, nil

}

func (jR *JSONRepo) GetByTag(ctx context.Context, tag string) ([]*model.Task, error) {
	jR.mu.RLock()
	defer jR.mu.RUnlock()

	taggetTasks := make([]*model.Task, 0)
	tasks, _ := jR.GetAll(ctx)
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

func (jR *JSONRepo) Complete(ctx context.Context, id int) (*model.Task, error) {

	task, err := jR.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении задачи: %w | %w", myErrors.ErrIdNotExists, err)
	}

	if task.Completed {
		return nil, fmt.Errorf("Ошибка при завершении задачи: %w | %w", myErrors.ErrTaskAlredyComplete, err)
	}

	jR.mu.Lock()
	completeTime := time.Now()
	task.Completed = true
	task.CompletedAt = &completeTime
	jR.mu.Unlock()

	if err := jR.saveToFile(); err != nil {
		return nil, fmt.Errorf("Ошибка при сохранении задачи: %w | %w", myErrors.ErrCantSaveTaskToJson, err)
	}
	return task, nil

}

func (jR *JSONRepo) DeleteByID(ctx context.Context, id int) error {
	_, err := jR.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении задачи: %w | %w", myErrors.ErrIdNotExists, err)
	}
	jR.mu.Lock()
	delete(jR.tasks, id)
	jR.mu.Unlock()

	if err := jR.saveToFile(); err != nil {
		return fmt.Errorf("Ошибка при сохранении задачи: %w | %w", myErrors.ErrCantSaveTaskToJson, err)
	}
	return nil

}
