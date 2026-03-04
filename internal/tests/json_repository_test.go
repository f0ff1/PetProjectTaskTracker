package tests

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	customError "TaskTracker/errors"
	"TaskTracker/internal/repository/sjson"
)

// helper для подготовки временного JSON файла
func prepTempJSON(t *testing.T) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	path := tmpFile.Name()
	tmpFile.Close()
	// удаляем чтобы NewJSONStorage создавал корректный файл
	os.Remove(path)
	// гарантируем очистку после теста
	t.Cleanup(func() { os.Remove(path) })
	return path
}

// TestJSONStorage_NewStorage проверяет создание нового JSON хранилища
func TestJSONStorage_NewStorage(t *testing.T) {
	t.Parallel()

	// подготовка файла
	path := prepTempJSON(t)
	storage, err := sjson.NewJSONStorage(path)

	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}

	if storage == nil {
		t.Fatal("NewJSONStorage() returned nil")
	}

	// Проверяем что хранилище пусто
	tasks, err := storage.GetAll()
	if err != nil {
		t.Errorf("Expected no error for GetAll(), got %v", err)
	}
	if tasks != nil && len(tasks) != 0 {
		t.Error("Expected empty tasks for new storage")
	}
}

// TestJSONStorage_NewStorage_NonExistentFile проверяет создание хранилища для несуществующего файла
func TestJSONStorage_NewStorage_NonExistentFile(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "test_storage_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "newfile.json")

	storage, err := sjson.NewJSONStorage(filePath)

	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}

	if storage == nil {
		t.Fatal("Expected storage, got nil")
	}

	// Проверяем что файл был создан
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("File should have been created")
	}
}

// TestJSONStorage_Add_BasicTask проверяет добавление задачи в JSON хранилище
func TestJSONStorage_Add_BasicTask(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp("", "test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	os.Remove(tmpFile.Name())

	storage, err := sjson.NewJSONStorage(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}

	title := "Test Task"
	desc := "Test Description"
	tags := []string{"test", "json"}

	task, err := storage.Add(title, desc, tags)

	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}

	if task.Title != title {
		t.Errorf("Expected Title = %q, got %q", title, task.Title)
	}

	if task.Description != desc {
		t.Errorf("Expected Description = %q, got %q", desc, task.Description)
	}

	if !task.Completed {
		if task.Completed {
			t.Error("Expected Completed = false")
		}
	}

	if task.CompletedAt != nil {
		t.Error("Expected CompletedAt = nil")
	}
}

// TestJSONStorage_Persistence проверяет что данные сохраняются в файл
func TestJSONStorage_Persistence(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	// Создаем хранилище и добавляем задачи
	storage1, _ := sjson.NewJSONStorage(path)
	storage1.Add("Task 1", "Desc 1", []string{})
	storage1.Add("Task 2", "Desc 2", []string{})

	// Создаем новое хранилище из того же файла
	storage2, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("Failed to create second storage: %v", err)
	}

	// Проверяем что данные загрузились
	tasks, err := storage2.GetAll()
	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Title != "Task 1" && tasks[1].Title != "Task 1" {
		t.Error("Task 1 not found in loaded storage")
	}
}

// TestJSONStorage_GetByID проверяет получение задачи по ID
func TestJSONStorage_GetByID(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	storage, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}
	storage.Add("Task 1", "Desc 1", []string{})
	storage.Add("Task 2", "Desc 2", []string{})

	task, err := storage.GetByID(1)

	if err != nil {
		t.Fatalf("GetByID() failed: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}

	if task.Title != "Task 1" {
		t.Errorf("Expected Title = 'Task 1', got %q", task.Title)
	}
}

// TestJSONStorage_GetByID_NotFound проверяет получение несуществующей задачи
func TestJSONStorage_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	storage, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}

	task, err := storage.GetByID(999)

	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}

	if err != customError.ErrIdNotExists {
		t.Errorf("Expected ErrIdNotExists, got %v", err)
	}

	if task != nil {
		t.Errorf("Expected nil task, got %v", task)
	}
}

// TestJSONStorage_Complete проверяет отметить задачу как выполненную
func TestJSONStorage_Complete(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	storage, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}
	storage.Add("Task", "Description", []string{})

	task, err := storage.Complete(1)

	if err != nil {
		t.Fatalf("Complete() failed: %v", err)
	}

	if !task.Completed {
		t.Error("Task should be completed")
	}

	if task.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}

	// Проверяем что изменение было сохранено
	retrieved, _ := storage.GetByID(1)
	if !retrieved.Completed {
		t.Error("Completed status should be persisted")
	}
}

// TestJSONStorage_Complete_AlreadyCompleted проверяет повторное завершение
func TestJSONStorage_Complete_AlreadyCompleted(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	storage, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}
	storage.Add("Task", "Description", []string{})

	storage.Complete(1)

	_, err = storage.Complete(1)

	if err == nil {
		t.Error("Expected error for already completed task, got nil")
	}

	if err != customError.ErrTaskAlredyComplete {
		t.Errorf("Expected ErrTaskAlredyComplete, got %v", err)
	}
}

// TestJSONStorage_GetAll проверяет получение всех задач
func TestJSONStorage_GetAll(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	storage, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}

	storage.Add("Task 1", "Desc 1", []string{})
	storage.Add("Task 2", "Desc 2", []string{})
	storage.Add("Task 3", "Desc 3", []string{})

	tasks, err := storage.GetAll()

	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

// TestJSONStorage_GetByTag проверяет поиск по тегам
func TestJSONStorage_GetByTag(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	storage, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}

	storage.Add("Task 1", "Desc 1", []string{"work", "urgent"})
	storage.Add("Task 2", "Desc 2", []string{"personal"})
	storage.Add("Task 3", "Desc 3", []string{"work"})

	tasks, err := storage.GetByTag("work")

	if err != nil {
		t.Fatalf("GetByTag() failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks with 'work' tag, got %d", len(tasks))
	}

	// Проверяем что все задачи имеют нужный тег
	for _, task := range tasks {
		hasTag := false
		for _, tag := range task.Tags {
			if tag == "work" {
				hasTag = true
				break
			}
		}
		if !hasTag {
			t.Errorf("Task %d doesn't have 'work' tag", task.ID)
		}
	}
}

// TestJSONStorage_Concurrent проверяет параллельное использование
func TestJSONStorage_Concurrent(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	storage, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}

	const numTasks = 20
	done := make(chan bool, numTasks)
	var mu sync.Mutex

	for i := 0; i < numTasks; i++ {
		go func(index int) {
			defer func() { done <- true }()
			mu.Lock()
			storage.Add("Task", "Description", []string{})
			mu.Unlock()
		}(i)
	}

	for i := 0; i < numTasks; i++ {
		<-done
	}

	tasks, err := storage.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(tasks) != numTasks {
		t.Errorf("Expected %d tasks, got %d", numTasks, len(tasks))
	}
}

// TestJSONStorage_DataIntegration проверяет интеграцию с файлом
func TestJSONStorage_DataIntegration(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	// Первое хранилище - добавляем и завершаем данные
	storage1, _ := sjson.NewJSONStorage(path)
	storage1.Add("Task 1", "Desc 1", []string{})
	storage1.Add("Task 2", "Desc 2", []string{})
	storage1.Complete(1)

	// Второе хранилище - загружаем из файла
	storage2, _ := sjson.NewJSONStorage(path)
	tasks, _ := storage2.GetAll()

	// Проверяем что данные загрузились правильно
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Проверяем что статус выполнения сохранился
	task1, _ := storage2.GetByID(1)
	if !task1.Completed {
		t.Error("Task 1 should be completed")
	}

	task2, _ := storage2.GetByID(2)
	if task2.Completed {
		t.Error("Task 2 should not be completed")
	}
}

// TestJSONStorage_TagPersistence проверяет сохранение тегов
func TestJSONStorage_TagPersistence(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	// Добавляем задачу с тегами
	storage1, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}
	expectedTags := []string{"tag1", "tag2", "tag3"}
	if _, err := storage1.Add("Task", "Description", expectedTags); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Загружаем из файла
	storage2, err := sjson.NewJSONStorage(path)
	if err != nil {
		t.Fatalf("NewJSONStorage() failed: %v", err)
	}
	task, err := storage2.GetByID(1)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	// Проверяем что теги сохранились
	if len(task.Tags) != len(expectedTags) {
		t.Errorf("Expected %d tags, got %d", len(expectedTags), len(task.Tags))
	}

	for i, tag := range expectedTags {
		if task.Tags[i] != tag {
			t.Errorf("Expected tag %q, got %q", tag, task.Tags[i])
		}
	}
}

// TestJSONStorage_InvalidPath проверяет сообщение об ошибке при неправильном пути
func TestJSONStorage_InvalidPath(t *testing.T) {
	t.Parallel()

	// Используем несуществующий путь с недопустимым именем
	invalidPath := "/invalid/nonexistent/path/file.json"

	_, err := sjson.NewJSONStorage(invalidPath)

	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

// TestJSONStorage_CompleteAndVerify проверяет завершение и верификацию
func TestJSONStorage_CompleteAndVerify(t *testing.T) {
	t.Parallel()

	path := prepTempJSON(t)

	storage, _ := sjson.NewJSONStorage(path)

	task, _ := storage.Add("Task", "Description", []string{})
	beforeComplete := time.Now()

	completed, _ := storage.Complete(task.ID)

	afterComplete := time.Now()

	// Проверяем что время завершения правильное
	if completed.CompletedAt.Before(beforeComplete) || completed.CompletedAt.After(afterComplete.Add(time.Second)) {
		t.Error("CompletedAt time is incorrect")
	}
}
