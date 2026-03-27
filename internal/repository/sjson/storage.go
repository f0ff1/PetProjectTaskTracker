package sjson

import (
	"encoding/json"
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
			return nil, myErrors.ErrCantCreateJsonFile
		}
		return storage, nil
	}

	if err := storage.loadFromFile(); err != nil {
		return nil, myErrors.ErrCantReadJsonData
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
		return myErrors.ErrWrongPath
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
		return myErrors.ErrCantReadJsonData
	}
	defer file.Close()

	var data JSONData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	jR.nextID = data.NextID
	jR.tasks = data.Tasks
	return nil

}

func (jR *JSONRepo) Add(title, description string, tags []string) (*model.Task, error) {
	jR.mu.Lock()
	task := &model.Task{
		ID:          jR.nextID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        tags,
	}

	jR.tasks[jR.nextID] = task
	jR.nextID++
	jR.mu.Unlock()

	if err := jR.saveToFile(); err != nil {
		return nil, myErrors.ErrCantSaveTaskToJson
	}

	return task, nil

}

func (jR *JSONRepo) GetAll() ([]*model.Task, error) {
	jR.mu.RLock()
	defer jR.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(jR.tasks))
	for _, task := range jR.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (jR *JSONRepo) GetByID(id int) (*model.Task, error) {
	jR.mu.RLock()
	defer jR.mu.RUnlock()

	task, exists := jR.tasks[id]
	if !exists {
		return nil, myErrors.ErrIdNotExists
	}
	return task, nil

}

func (jR *JSONRepo) GetByTag(tag string) ([]*model.Task, error) {
	jR.mu.RLock()
	defer jR.mu.RUnlock()

	taggetTasks := make([]*model.Task, 0)
	tasks, _ := jR.GetAll()
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

func (jR *JSONRepo) Complete(id int) (*model.Task, error) {

	task, err := jR.GetByID(id)
	if err != nil {
		return nil, myErrors.ErrIdNotExists
	}

	if task.Completed {
		return nil, myErrors.ErrTaskAlredyComplete
	}

	jR.mu.Lock()
	completeTime := time.Now()
	task.Completed = true
	task.CompletedAt = &completeTime
	jR.mu.Unlock()

	if err := jR.saveToFile(); err != nil {
		return nil, myErrors.ErrCantSaveTaskToJson
	}
	return task, nil

}

func (jR *JSONRepo) DeleteByID(id int) error {
	_, err := jR.GetByID(id)
	if err != nil {
		return myErrors.ErrIdNotExists
	}
	jR.mu.Lock()
	delete(jR.tasks, id)
	jR.mu.Unlock()

	if err := jR.saveToFile(); err != nil {
		return myErrors.ErrCantSaveTaskToJson
	}
	return nil

}
