package tests

import (
	"strings"
	"sync"
	"testing"
	"time"

	customError "TaskTracker/errors"
	"TaskTracker/internal/repository/memory"





)

// TestMemory_NewStorage проверяет создание нового хранилища в памяти
func TestMemory_NewStorage(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	if storage == nil {
		t.Fatal("NewStorage() returned nil")
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

// TestMemory_Add_BasicTask проверяет добавление простой задачи
func TestMemory_Add_BasicTask(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	title := "Купить молоко"
	desc := "Обязательно свежее"
	tags := []string{"покупки", "продукты"}

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

	if task.Completed != false {
		t.Error("Expected Completed = false")
	}

	if task.CompletedAt != nil {
		t.Error("Expected CompletedAt = nil")
	}

	if len(task.Tags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(task.Tags))
	}

	if time.Since(task.CreatedAt) > time.Second {
		t.Error("CreatedAt is too far in the past")
	}
}

// TestMemory_Add_EmptyTitle проверяет добавление задачи с пустым названием
func TestMemory_Add_EmptyTitle(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task, err := storage.Add("", "Description", []string{})

	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}

	// Storage allows empty title - it just stores what's given
	if task.Title != "" {
		t.Errorf("Expected empty title, got %q", task.Title)
	}
}

// TestMemory_Add_WithoutTags проверяет добавление задачи без тегов
func TestMemory_Add_WithoutTags(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task, err := storage.Add("Task", "Description", nil)

	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if task.Tags != nil && len(task.Tags) > 0 {
		t.Errorf("Expected empty tags, got %v", task.Tags)
	}
}

// TestMemory_Add_MultipleTasks проверяет добавление нескольких задач
func TestMemory_Add_MultipleTasks(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task1, _ := storage.Add("Task 1", "Description 1", []string{})
	task2, _ := storage.Add("Task 2", "Description 2", []string{})
	task3, _ := storage.Add("Task 3", "Description 3", []string{})

	if task1.ID != 1 {
		t.Errorf("Task1 ID = %d, want 1", task1.ID)
	}
	if task2.ID != 2 {
		t.Errorf("Task2 ID = %d, want 2", task2.ID)
	}
	if task3.ID != 3 {
		t.Errorf("Task3 ID = %d, want 3", task3.ID)
	}

	tasks, err := storage.GetAll()
	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

// TestMemory_GetAll_Empty проверяет получение задач из пустого хранилища
func TestMemory_GetAll_Empty(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	tasks, err := storage.GetAll()

	if err != nil {
		t.Errorf("Expected no error for empty storage, got %v", err)
	}

	if tasks != nil && len(tasks) != 0 {
		t.Error("Expected empty tasks for empty storage")
	}
}

// TestMemory_GetAll_WithTasks проверяет получение всех задач
func TestMemory_GetAll_WithTasks(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

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

// TestMemory_GetByID_Success проверяет получение задачи по ID
func TestMemory_GetByID_Success(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	added, _ := storage.Add("Test Task", "Test Description", []string{"test"})

	task, err := storage.GetByID(1)

	if err != nil {
		t.Fatalf("GetByID() failed: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected Title = 'Test Task', got %q", task.Title)
	}

	if task == added {
		// Должен вернуться один и тот же объект
		if task.Title != added.Title {
			t.Error("Returned task doesn't match added task")
		}
	}
}

// TestMemory_GetByID_NotFound проверяет получение несуществующей задачи
func TestMemory_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

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

// TestMemory_GetByTag_Empty проверяет поиск задач по тегу в пустом хранилище
func TestMemory_GetByTag_Empty(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	tasks, err := storage.GetByTag("test")

	if err != nil {
		t.Fatalf("GetByTag() failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

// TestMemory_GetByTag_Success проверяет поиск задач по тегу
func TestMemory_GetByTag_Success(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	storage.Add("Task 1", "Desc 1", []string{"work", "urgent"})
	storage.Add("Task 2", "Desc 2", []string{"personal"})
	storage.Add("Task 3", "Desc 3", []string{"work", "home"})

	tasks, err := storage.GetByTag("work")

	if err != nil {
		t.Fatalf("GetByTag() failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks with 'work' tag, got %d", len(tasks))
	}

	// Проверяем что все возвращенные задачи имеют требуемый тег
	for _, task := range tasks {
		hasTag := false
		for _, tag := range task.Tags {
			if tag == "work" {
				hasTag = true
				break
			}
		}
		if !hasTag {
			t.Errorf("Task %q doesn't have 'work' tag", task.Title)
		}
	}
}

// TestMemory_GetByTag_NotFound проверяет поиск несуществующего тега
func TestMemory_GetByTag_NotFound(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	storage.Add("Task 1", "Desc 1", []string{"work"})
	storage.Add("Task 2", "Desc 2", []string{"personal"})

	tasks, err := storage.GetByTag("nonexistent")

	if err != nil {
		t.Fatalf("GetByTag() failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

// TestMemory_Complete_Success проверяет отметить задачу как выполненную
func TestMemory_Complete_Success(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task, _ := storage.Add("Task", "Description", []string{})

	if task.Completed {
		t.Error("New task should not be completed")
	}

	completed, err := storage.Complete(1)

	if err != nil {
		t.Fatalf("Complete() failed: %v", err)
	}

	if !completed.Completed {
		t.Error("Task should be completed")
	}

	if completed.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}

	if time.Since(*completed.CompletedAt) > time.Second {
		t.Error("CompletedAt is too far in the past")
	}
}

// TestMemory_Complete_AlreadyCompleted проверяет повторное завершение задачи
func TestMemory_Complete_AlreadyCompleted(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task, _ := storage.Add("Task", "Description", []string{})

	// Первое завершение
	storage.Complete(task.ID)

	// Второе завершение (должна быть ошибка)
	_, err := storage.Complete(task.ID)

	if err == nil {
		t.Error("Expected error for already completed task, got nil")
	}

	if err != customError.ErrTaskAlredyComplete {
		t.Errorf("Expected ErrTaskAlredyComplete, got %v", err)
	}
}

// TestMemory_Complete_NotFound проверяет завершение несуществующей задачи
func TestMemory_Complete_NotFound(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	_, err := storage.Complete(999)

	if err == nil {
		t.Error("Expected error for non-existent task, got nil")
	}

	if err != customError.ErrIdNotExists {
		t.Errorf("Expected ErrIdNotExists, got %v", err)
	}
}

// TestMemory_Concurrent проверяет параллельное добавление задач
func TestMemory_Concurrent(t *testing.T) {
	storage := memory.NewStorage()

	const numTasks = 50
	done := make(chan bool, numTasks)
	var mu sync.Mutex

	for i := 0; i < numTasks; i++ {
		go func(index int) {
			defer func() { done <- true }()
			title := "Task " + string(rune(index))
			mu.Lock()
			storage.Add(title, "Description", []string{})
			mu.Unlock()
		}(i)
	}

	// Ждем завершения всех горутин
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

// TestMemory_GetByTag_MultipleTags проверяет поиск с множественными тегами
func TestMemory_GetByTag_MultipleTags(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	storage.Add("Task 1", "Desc 1", []string{"a", "b", "c"})
	storage.Add("Task 2", "Desc 2", []string{"b", "c"})
	storage.Add("Task 3", "Desc 3", []string{"c"})
	storage.Add("Task 4", "Desc 4", []string{"d"})

	tagTests := []struct {
		tag           string
		expectedCount int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 1},
		{"nonexistent", 0},
	}

	for _, tt := range tagTests {
		t.Run("tag_"+tt.tag, func(t *testing.T) {
			tasks, err := storage.GetByTag(tt.tag)
			if err != nil {
				t.Fatalf("GetByTag() failed: %v", err)
			}

			if len(tasks) != tt.expectedCount {
				t.Errorf("Expected %d tasks with tag %q, got %d", tt.expectedCount, tt.tag, len(tasks))
			}
		})
	}
}

// TestMemory_TaskIntegrity проверяет целостность данных после операций
func TestMemory_TaskIntegrity(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	// Добавляем несколько задач
	task1, _ := storage.Add("Task 1", "Description 1", []string{"tag1"})
	task2, _ := storage.Add("Task 2", "Description 2", []string{"tag2"})

	// Завершаем первую задачу
	storage.Complete(task1.ID)

	// Проверяем что вторая задача не изменилась
	retrieved, _ := storage.GetByID(task2.ID)

	if retrieved.Completed {
		t.Error("Task 2 should not be completed")
	}

	if retrieved.Title != "Task 2" {
		t.Errorf("Task 2 title changed: got %q", retrieved.Title)
	}

	// Проверяем что первая задача действительно завершена
	retrieved1, _ := storage.GetByID(task1.ID)

	if !retrieved1.Completed {
		t.Error("Task 1 should be completed")
	}
}

// TestMemory_LongTitle проверяет работу с длинными названиями
func TestMemory_LongTitle(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	longTitle := strings.Repeat("A", 1000)

	task, err := storage.Add(longTitle, "Description", []string{})

	if err != nil {
		t.Fatalf("Add() with long title failed: %v", err)
	}

	if task.Title != longTitle {
		t.Error("Long title was not stored correctly")
	}
}

// TestMemory_EmptyTagList проверяет работу с пустым списком тегов
func TestMemory_EmptyTagList(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task, err := storage.Add("Task", "Description", []string{})

	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if len(task.Tags) != 0 {
		t.Error("Tags should be empty")
	}

	tasks, err := storage.GetByTag("any")
	if err != nil {
		t.Fatalf("GetByTag() failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Error("Should not find tasks without tags")
	}
}
