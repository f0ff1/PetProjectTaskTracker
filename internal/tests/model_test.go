package tests

import (
	"testing"
	"time"

	"TaskTracker/internal/model"
)

// TestTask_Creation проверяет создание структуры Task
func TestTask_Creation(t *testing.T) {
	t.Parallel()

	now := time.Now()
	task := model.Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   now,
		CompletedAt: nil,
		Tags:        []string{"test"},
	}

	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected Title = 'Test Task', got %q", task.Title)
	}

	if task.Description != "Test Description" {
		t.Errorf("Expected Description = 'Test Description', got %q", task.Description)
	}

	if task.Completed {
		t.Error("Expected Completed = false")
	}

	if !task.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt = %v, got %v", now, task.CreatedAt)
	}

	if task.CompletedAt != nil {
		t.Error("Expected CompletedAt = nil")
	}

	if len(task.Tags) != 1 || task.Tags[0] != "test" {
		t.Errorf("Expected Tags = ['test'], got %v", task.Tags)
	}
}

// TestTask_Completed проверяет завершенную задачу
func TestTask_Completed(t *testing.T) {
	t.Parallel()

	now := time.Now()
	completedTime := time.Now().Add(time.Hour)

	task := model.Task{
		ID:          2,
		Title:       "Completed Task",
		Description: "This task is done",
		Completed:   true,
		CreatedAt:   now,
		CompletedAt: &completedTime,
		Tags:        []string{"done"},
	}

	if !task.Completed {
		t.Error("Expected Completed = true")
	}

	if task.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}

	if !task.CompletedAt.Equal(completedTime) {
		t.Errorf("Expected CompletedAt = %v, got %v", completedTime, *task.CompletedAt)
	}
}

// TestTask_EmptyTags проверяет задачу без тегов
func TestTask_EmptyTags(t *testing.T) {
	t.Parallel()

	task := model.Task{
		ID:          3,
		Title:       "No Tags Task",
		Description: "Task without tags",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        []string{},
	}

	if len(task.Tags) != 0 {
		t.Errorf("Expected empty tags, got %v", task.Tags)
	}
}

// TestTask_NilTags проверяет задачу с nil тегами
func TestTask_NilTags(t *testing.T) {
	t.Parallel()

	task := model.Task{
		ID:          4,
		Title:       "Nil Tags Task",
		Description: "Task with nil tags",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        nil,
	}

	if task.Tags != nil && len(task.Tags) != 0 {
		t.Errorf("Expected nil or empty tags, got %v", task.Tags)
	}
}

// TestTask_MultipleTags проверяет задачу с несколькими тегами
func TestTask_MultipleTags(t *testing.T) {
	t.Parallel()

	tags := []string{"urgent", "work", "important", "todo"}

	task := model.Task{
		ID:          5,
		Title:       "Multi-tag Task",
		Description: "Task with multiple tags",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        tags,
	}

	if len(task.Tags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(task.Tags))
	}

	for i, tag := range tags {
		if task.Tags[i] != tag {
			t.Errorf("Expected tag %q at position %d, got %q", tag, i, task.Tags[i])
		}
	}
}

// TestTask_Fields проверяет все поля Task
func TestTask_Fields(t *testing.T) {
	t.Parallel()

	task := &model.Task{
		ID:          100,
		Title:       "Test",
		Description: "Desc",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        []string{"test"},
	}

	// Проверяем что поля доступны для изменения
	task.Completed = true
	if !task.Completed {
		t.Error("Should be able to modify Completed field")
	}

	task.Title = "New Title"
	if task.Title != "New Title" {
		t.Error("Should be able to modify Title field")
	}

	task.Description = "New Description"
	if task.Description != "New Description" {
		t.Error("Should be able to modify Description field")
	}

	now := time.Now()
	task.CompletedAt = &now
	if task.CompletedAt == nil {
		t.Error("Should be able to modify CompletedAt field")
	}
}

// TestTask_TimeComparison проверяет сравнение времени в Task
func TestTask_TimeComparison(t *testing.T) {
	t.Parallel()

	createdTime := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	completedTime := time.Date(2026, 2, 1, 14, 0, 0, 0, time.UTC)

	task := model.Task{
		ID:          6,
		Title:       "Time Task",
		Description: "Testing time fields",
		Completed:   true,
		CreatedAt:   createdTime,
		CompletedAt: &completedTime,
		Tags:        []string{},
	}

	if task.CreatedAt.After(*task.CompletedAt) {
		t.Error("CompletedAt should be after CreatedAt")
	}

	if !task.CreatedAt.Before(*task.CompletedAt) {
		t.Error("CreatedAt should be before CompletedAt")
	}
}

// TestTask_Pointer проверяет работу с указателями на Task
func TestTask_Pointer(t *testing.T) {
	t.Parallel()

	task1 := model.Task{
		ID:          7,
		Title:       "Task 1",
		Description: "Description 1",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        []string{"tag1"},
	}

	task2 := &task1

	// Изменяем через указатель
	task2.Title = "Modified Title"

	if task1.Title != "Modified Title" {
		t.Error("Modifying through pointer should modify original")
	}
}

// TestTask_TagModification проверяет изменение тегов
func TestTask_TagModification(t *testing.T) {
	t.Parallel()

	task := model.Task{
		ID:          8,
		Title:       "Tag Mod Task",
		Description: "Testing tag modification",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        []string{"tag1", "tag2"},
	}

	// Добавляем новый тег
	task.Tags = append(task.Tags, "tag3")

	if len(task.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(task.Tags))
	}

	if task.Tags[2] != "tag3" {
		t.Errorf("Expected 'tag3', got %q", task.Tags[2])
	}
}

// TestTask_ZeroValues проверяет Task с нулевыми значениями
func TestTask_ZeroValues(t *testing.T) {
	t.Parallel()

	var task model.Task

	if task.ID != 0 {
		t.Errorf("Expected ID = 0, got %d", task.ID)
	}

	if task.Title != "" {
		t.Errorf("Expected empty Title, got %q", task.Title)
	}

	if task.Description != "" {
		t.Errorf("Expected empty Description, got %q", task.Description)
	}

	if task.Completed {
		t.Error("Expected Completed = false")
	}

	if !task.CreatedAt.IsZero() {
		t.Errorf("Expected zero CreatedAt, got %v", task.CreatedAt)
	}

	if task.CompletedAt != nil {
		t.Errorf("Expected nil CompletedAt, got %v", task.CompletedAt)
	}

	if task.Tags != nil {
		t.Errorf("Expected nil Tags, got %v", task.Tags)
	}
}

// TestTask_TagOrder проверяет порядок тегов
func TestTask_TagOrder(t *testing.T) {
	t.Parallel()

	tags := []string{"first", "second", "third", "fourth", "fifth"}

	task := model.Task{
		ID:          9,
		Title:       "Tag Order Task",
		Description: "Testing tag order",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        tags,
	}

	for i, expectedTag := range tags {
		if task.Tags[i] != expectedTag {
			t.Errorf("Tag at position %d: expected %q, got %q", i, expectedTag, task.Tags[i])
		}
	}
}

// TestTask_Copy проверяет копирование Task
func TestTask_Copy(t *testing.T) {
	t.Parallel()

	original := model.Task{
		ID:          10,
		Title:       "Original Task",
		Description: "Original Description",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        []string{"tag1", "tag2"},
	}

	// Создаем копию
	copy := original

	// Изменяем копию
	copy.Title = "Modified Title"

	// Проверяем что оригинал не изменился
	if original.Title != "Original Task" {
		t.Error("Modifying copy should not affect original")
	}

	if copy.Title != "Modified Title" {
		t.Error("Copy modification failed")
	}
}

// TestTask_LongDescription проверяет задачу с длинным описанием
func TestTask_LongDescription(t *testing.T) {
	t.Parallel()

	longDesc := "This is a very long description. " // Повторяем много раз
	for i := 0; i < 100; i++ {
		longDesc += "This is a very long description. "
	}

	task := model.Task{
		ID:          11,
		Title:       "Long Desc Task",
		Description: longDesc,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        []string{},
	}

	if len(task.Description) != len(longDesc) {
		t.Error("Long description not stored correctly")
	}

	if task.Description != longDesc {
		t.Error("Long description content mismatch")
	}
}

// TestTask_SpecialCharacters проверяет Task со специальными символами
func TestTask_SpecialCharacters(t *testing.T) {
	t.Parallel()

	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"

	task := model.Task{
		ID:          12,
		Title:       specialChars,
		Description: specialChars,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        []string{specialChars},
	}

	if task.Title != specialChars {
		t.Error("Special characters in Title not preserved")
	}

	if task.Description != specialChars {
		t.Error("Special characters in Description not preserved")
	}

	if task.Tags[0] != specialChars {
		t.Error("Special characters in Tags not preserved")
	}
}

// TestTask_CompletedAtNil проверяет что CompletedAt правильно nil
func TestTask_CompletedAtNil(t *testing.T) {
	t.Parallel()

	task := model.Task{
		ID:          13,
		Title:       "Not Completed",
		Description: "Should have nil CompletedAt",
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Tags:        []string{},
	}

	if task.CompletedAt != nil {
		t.Error("CompletedAt should be nil for incomplete task")
	}

	// Попытка разыменования должна привести к панике
	// но мы не проверяем это в тесте для безопасности
}

// TestTask_CompletedAtDeref проверяет разыменование CompletedAt
func TestTask_CompletedAtDeref(t *testing.T) {
	t.Parallel()

	completedTime := time.Now()
	task := model.Task{
		ID:          14,
		Title:       "Completed",
		Description: "Should have CompletedAt",
		Completed:   true,
		CreatedAt:   time.Now(),
		CompletedAt: &completedTime,
		Tags:        []string{},
	}

	if task.CompletedAt == nil {
		t.Error("CompletedAt should not be nil for completed task")
	}

	if !(*task.CompletedAt).Equal(completedTime) {
		t.Errorf("CompletedAt dereferenced value mismatch")
	}
}
