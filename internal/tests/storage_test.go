package tests

import (
	"strings"
	"testing"
	"time"

	"TaskTracker/internal/repository/memory"
)

// TestNewStorage проверяет создание нового хранилища
func TestNewStorage(t *testing.T) {
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

// TestStorage_IsEmpty проверяет методы для работы

// TestStorage_Add_WithCustomTitle проверяет добавление с кастомным названием
func TestStorage_Add_WithCustomTitle(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	title := "Купить молоко"
	description := "Обязательно свежее"
	tags := []string{"покупки"}

	task, err := storage.Add(title, description, tags)

	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}

	if task.Title != title {
		t.Errorf("Expected Title = %q, got %q", title, task.Title)
	}

	if task.Description != description {
		t.Errorf("Expected Description = %q, got %q", description, task.Description)
	}

	if task.Completed {
		t.Error("Expected Completed = false")
	}

	if time.Since(task.CreatedAt) > time.Second {
		t.Error("CreatedAt is too far in the past")
	}

	if task.CompletedAt != nil {
		t.Error("Expected CompletedAt = nil")
	}

	if len(task.Tags) != 1 || task.Tags[0] != "покупки" {
		t.Errorf("Expected Tags = ['покупки'], got %v", task.Tags)
	}
}

// TestStorage_Add_WithEmptyTitle проверяет добавление с пустым названием
func TestStorage_Add_WithEmptyTitle(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	description := "Задача без названия"

	task, err := storage.Add("", description, []string{})

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

	if task.Description != description {
		t.Errorf("Expected Description = %q, got %q", description, task.Description)
	}
}

// TestStorage_Add_MultipleTasks проверяет добавление нескольких задач
func TestStorage_Add_MultipleTasks(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task1, err1 := storage.Add("Task 1", "Description 1", []string{})
	task2, err2 := storage.Add("", "Description 2", []string{})
	task3, err3 := storage.Add("Task 3", "", []string{})

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("Add() failed")
	}

	if task1.ID != 1 {
		t.Errorf("Task1 ID = %d, want 1", task1.ID)
	}
	if task2.ID != 2 {
		t.Errorf("Task2 ID = %d, want 2", task2.ID)
	}
	if task3.ID != 3 {
		t.Errorf("Task3 ID = %d, want 3", task3.ID)
	}

	allTasks, err := storage.GetAll()
	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks in storage, got %d", len(allTasks))
	}

	if task1.Title == task2.Title {
		t.Error("Task1 and Task2 should have different titles")
	}
}

// TestStorage_GetAll_Empty проверяет получение из пустого хранилища
func TestStorage_GetAll_Empty(t *testing.T) {
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

// TestStorage_GetAll_WithTasks проверяет получение всех задач
func TestStorage_GetAll_WithTasks(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	storage.Add("Задача 1", "Описание 1", []string{})
	storage.Add("Задача 2", "Описание 2", []string{})
	storage.Add("", "Описание 3", []string{})

	tasks, err := storage.GetAll()

	if err != nil {
		t.Fatalf("GetAll() returned error: %v", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(tasks))
	}

	if tasks[0].Title != "Задача 1" {
		t.Errorf("Expected first task title 'Задача 1', got %q", tasks[0].Title)
	}
}

// TestStorage_GetByID_Success проверяет получение по ID
func TestStorage_GetByID_Success(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	added1, _ := storage.Add("Task 1", "Description 1", []string{})
	added2, _ := storage.Add("Task 2", "Description 2", []string{})

	task, err := storage.GetByID(1)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if task == nil {
		t.Fatal("Expected task, got nil")
	}
	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
	if task.Title != "Task 1" {
		t.Errorf("Expected title 'Task 1', got %q", task.Title)
	}

	task2, err := storage.GetByID(2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if task2.Title != "Task 2" {
		t.Errorf("Expected title 'Task 2', got %q", task2.Title)
	}

	_ = added1
	_ = added2
}

// TestStorage_GetByID_NotFound проверяет получение несуществующего ID
func TestStorage_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	storage.Add("Task", "Description", []string{})

	task, err := storage.GetByID(999)

	if err == nil {
		t.Fatal("Expected error for non-existent ID, got nil")
	}
	if task != nil {
		t.Errorf("Expected nil task, got %v", task)
	}
}

// TestStorage_Complete_Success проверяет завершение задачи
func TestStorage_Complete_Success(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	task, _ := storage.Add("Test Task", "Description", []string{})

	if task.Completed {
		t.Error("New task should not be completed")
	}
	if task.CompletedAt != nil {
		t.Error("New task should have nil CompletedAt")
	}

	completed, err := storage.Complete(task.ID)

	if err != nil {
		t.Fatalf("Complete() failed: %v", err)
	}

	if !completed.Completed {
		t.Error("Task should be completed after Complete()")
	}
	if completed.CompletedAt == nil {
		t.Error("CompletedAt should be set after Complete()")
	}

	if time.Since(*completed.CompletedAt) > time.Second {
		t.Error("CompletedAt is too far in the past")
	}

	retrieved, _ := storage.GetByID(task.ID)
	if !retrieved.Completed {
		t.Error("Task in storage should be completed")
	}
}

// TestStorage_Complete_AlreadyCompleted проверяет повторное завершение
func TestStorage_Complete_AlreadyCompleted(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	task, _ := storage.Add("Test Task", "Description", []string{})

	_, err1 := storage.Complete(task.ID)
	if err1 != nil {
		t.Fatalf("First Complete() failed: %v", err1)
	}

	firstCompleteTime := *task.CompletedAt

	time.Sleep(10 * time.Millisecond)

	completed2, err2 := storage.Complete(task.ID)

	if err2 == nil {
		t.Error("Expected error for already completed task, got nil")
	}
	if completed2 != nil {
		t.Error("Expected nil task for error case")
	}

	if *task.CompletedAt != firstCompleteTime {
		t.Error("CompletedAt should not change")
	}
}

// TestStorage_Complete_NotFound проверяет завершение несуществующей задачи
func TestStorage_Complete_NotFound(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	_, err := storage.Complete(999)

	if err == nil {
		t.Fatal("Expected error for non-existent task, got nil")
	}
}

// TestStorage_GetByTag_Empty проверяет поиск в пустом хранилище
func TestStorage_GetByTag_Empty(t *testing.T) {
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

// TestStorage_GetByTag_Success проверяет поиск по тегу
func TestStorage_GetByTag_Success(t *testing.T) {
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

// TestStorage_GetByTag_MultipleTags проверяет поиск с множественными тегами
func TestStorage_GetByTag_MultipleTags(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	storage.Add("Task 1", "Desc 1", []string{"a", "b", "c"})
	storage.Add("Task 2", "Desc 2", []string{"b", "c"})
	storage.Add("Task 3", "Desc 3", []string{"c"})
	storage.Add("Task 4", "Desc 4", []string{"d"})

	tests := []struct {
		tag           string
		expectedCount int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 1},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
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

// TestStorage_TaskIntegrity проверяет целостность данных
func TestStorage_TaskIntegrity(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task1, _ := storage.Add("Task 1", "Description 1", []string{"tag1"})
	task2, _ := storage.Add("Task 2", "Description 2", []string{"tag2"})

	storage.Complete(task1.ID)

	retrieved, _ := storage.GetByID(task2.ID)

	if retrieved.Completed {
		t.Error("Task 2 should not be completed")
	}

	if retrieved.Title != "Task 2" {
		t.Errorf("Task 2 title changed: got %q", retrieved.Title)
	}

	retrieved1, _ := storage.GetByID(task1.ID)

	if !retrieved1.Completed {
		t.Error("Task 1 should be completed")
	}
}

// TestStorage_LongTitle проверяет работу с длинными названиями
func TestStorage_LongTitle(t *testing.T) {
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

// TestStorage_EmptyTagList проверяет работу с пустым списком тегов
func TestStorage_EmptyTagList(t *testing.T) {
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

// TestStorage_Concurrent проверяет параллельное добавление
func TestStorage_Concurrent(t *testing.T) {
	storage := memory.NewStorage()

	const numTasks = 10
	done := make(chan bool, numTasks)

	for i := 0; i < numTasks; i++ {
		go func(index int) {
			title := "Task"
			if index%2 == 0 {
				title = ""
			}
			storage.Add(title, "Description", []string{})
			done <- true
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
