package tests

import (
	"regexp"
	"testing"
	"time"

	"TaskTracker/internal/model"
	"TaskTracker/internal/repository/memory"
)

func TestNewStorage(t *testing.T) {
	t.Parallel() // Можно запускать параллельно, т.к. тесты независимы

	storage := memory.NewStorage()

	if storage == nil {
		t.Fatal("NewStorage() returned nil")
	}

	// Проверяем через вызов метода IsEmpty (косвенная проверка инициализации)
	if !storage.IsEmpty() {
		t.Error("New storage should be empty")
	}
}

func TestStorage_IsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		prepare  func(*memory.Storage)
		expected bool
	}{
		{
			name:     "new storage should be empty",
			prepare:  func(s *memory.Storage) {},
			expected: true,
		},
		{
			name: "storage with one task should not be empty",
			prepare: func(s *memory.Storage) {
				s.Add("Task 1", "Description 1")
			},
			expected: false,
		},
		{
			name: "storage with multiple tasks should not be empty",
			prepare: func(s *memory.Storage) {
				s.Add("Task 1", "Description 1")
				s.Add("Task 2", "Description 2")
				s.Add("Task 3", "Description 3")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := memory.NewStorage()
			tt.prepare(storage)

			result := storage.IsEmpty()

			if result != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStorage_Add_WithCustomTitle(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	title := "Купить молоко"
	description := "Обязательно свежее"

	task := storage.Add(title, description)

	// Проверяем ID
	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}

	// Проверяем название
	if task.Title != title {
		t.Errorf("Expected Title = %q, got %q", title, task.Title)
	}

	// Проверяем описание
	if task.Description != description {
		t.Errorf("Expected Description = %q, got %q", description, task.Description)
	}

	// Проверяем статус выполнения
	if task.Completed != false {
		t.Error("Expected Completed = false")
	}

	// Проверяем время создания (должно быть близко к текущему)
	if time.Since(task.CreatedAt) > time.Second {
		t.Error("CreatedAt is too far in the past")
	}

	// Проверяем что CompletedAt = nil
	if task.CompletedAt != nil {
		t.Error("Expected CompletedAt = nil")
	}

	// Проверяем что задача сохранилась в хранилище
	if storage.IsEmpty() {
		t.Error("Storage should not be empty after adding task")
	}

	// Проверяем что nextID увеличился
	nextTask := storage.Add("Another task", "")
	if nextTask.ID != 2 {
		t.Errorf("Expected next task ID = 2, got %d", nextTask.ID)
	}
}

func TestStorage_Add_WithEmptyTitle(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	description := "Задача без названия"

	task := storage.Add("", description)

	// Проверяем ID
	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}

	// Проверяем формат сгенерированного названия
	pattern := `^exr-\d{7}$`
	matched, err := regexp.MatchString(pattern, task.Title)
	if err != nil {
		t.Fatalf("Regex error: %v", err)
	}
	if !matched {
		t.Errorf("Generated title %q doesn't match pattern %q", task.Title, pattern)
	}

	// Проверяем длину (exr- + 7 цифр = 11 символов)
	if len(task.Title) != 11 {
		t.Errorf("Expected title length 11, got %d for title %q", len(task.Title), task.Title)
	}

	// Проверяем описание
	if task.Description != description {
		t.Errorf("Expected Description = %q, got %q", description, task.Description)
	}

	// Проверяем что задача сохранилась
	retrieved, err := storage.GetByID(1)
	if err != nil {
		t.Fatalf("Failed to get task by ID: %v", err)
	}
	if retrieved.Title != task.Title {
		t.Errorf("Retrieved task title = %q, want %q", retrieved.Title, task.Title)
	}
}

func TestStorage_Add_MultipleTasks(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	// Добавляем несколько задач
	task1 := storage.Add("Task 1", "Description 1")
	task2 := storage.Add("", "Description 2") // С автогенерацией
	task3 := storage.Add("Task 3", "")

	// Проверяем ID
	if task1.ID != 1 {
		t.Errorf("Task1 ID = %d, want 1", task1.ID)
	}
	if task2.ID != 2 {
		t.Errorf("Task2 ID = %d, want 2", task2.ID)
	}
	if task3.ID != 3 {
		t.Errorf("Task3 ID = %d, want 3", task3.ID)
	}

	// Проверяем количество задач в хранилище
	allTasks, err := storage.GetAll()
	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks in storage, got %d", len(allTasks))
	}

	// Проверяем что все задачи разные
	if task1 == task2 || task1 == task3 || task2 == task3 {
		t.Error("Tasks should be different objects")
	}

	// Проверяем уникальность названий (у task2 должно быть сгенерированное)
	if task1.Title == task2.Title {
		t.Error("Task1 and Task2 should have different titles")
	}
	if task2.Title == "Task 2" {
		t.Error("Task2 should have auto-generated title")
	}
}

func TestStorage_GetAll_Empty(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	tasks, err := storage.GetAll()

	// Должна быть ошибка
	if err == nil {
		t.Error("Expected error for empty storage, got nil")
	}

	// Проверяем сообщение об ошибке
	expectedErr := "Задач нет"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}

	// Должен вернуться пустой срез, не nil
	if tasks == nil {
		t.Error("Expected non-nil slice, got nil")
	}
	if len(tasks) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(tasks))
	}
}

func TestStorage_GetAll_WithTasks(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	// Добавляем задачи
	expectedTasks := []*model.Task{
		storage.Add("Задача 1", "Описание 1"),
		storage.Add("Задача 2", "Описание 2"),
		storage.Add("", "Описание 3"),
	}

	tasks, err := storage.GetAll()

	// Не должно быть ошибки
	if err != nil {
		t.Fatalf("GetAll() returned error: %v", err)
	}

	// Проверяем количество
	if len(tasks) != len(expectedTasks) {
		t.Fatalf("Expected %d tasks, got %d", len(expectedTasks), len(tasks))
	}

	// Создаем карту для проверки наличия всех задач
	taskMap := make(map[int]*model.Task)
	for _, task := range tasks {
		taskMap[task.ID] = task
	}

	// Проверяем что все ожидаемые задачи присутствуют
	for _, expected := range expectedTasks {
		found, exists := taskMap[expected.ID]
		if !exists {
			t.Errorf("Task with ID %d not found in result", expected.ID)
			continue
		}
		if found.Title != expected.Title {
			t.Errorf("Task %d: Expected title %q, got %q", expected.ID, expected.Title, found.Title)
		}
	}
}

func TestStorage_GetByID_Success(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	// Добавляем задачи
	added1 := storage.Add("Task 1", "Description 1")
	added2 := storage.Add("Task 2", "Description 2")

	tests := []struct {
		name     string
		id       int
		expected *model.Task
	}{
		{
			name:     "get first task",
			id:       1,
			expected: added1,
		},
		{
			name:     "get second task",
			id:       2,
			expected: added2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := storage.GetByID(tt.id)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if task != tt.expected {
				t.Error("Returned task doesn't match expected")
			}
			if task.ID != tt.id {
				t.Errorf("Expected ID %d, got %d", tt.id, task.ID)
			}
		})
	}
}

func TestStorage_GetByID_Errors(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	storage.Add("Task", "Description") // ID = 1

	tests := []struct {
		name        string
		id          int
		expectedErr string
	}{
		{
			name:        "id = 0",
			id:          0,
			expectedErr: "Неккоректный ID",
		},
		{
			name:        "negative id",
			id:          -5,
			expectedErr: "Неккоректный ID",
		},
		{
			name:        "non-existent id",
			id:          999,
			expectedErr: "Несуществующий ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := storage.GetByID(tt.id)

			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error %q, got %q", tt.expectedErr, err.Error())
			}
			if task != nil {
				t.Errorf("Expected nil task, got %v", task)
			}
		})
	}
}

func TestStorage_GetByID_EmptyStorage(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	task, err := storage.GetByID(1)

	if err == nil {
		t.Error("Expected error for empty storage, got nil")
	}
	expectedErr := "Список задач пуст"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
	if task != nil {
		t.Errorf("Expected nil task, got %v", task)
	}
}

func TestStorage_Complete_Success(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	task := storage.Add("Test Task", "Description")

	// Проверяем начальное состояние
	if task.Completed {
		t.Error("New task should not be completed")
	}
	if task.CompletedAt != nil {
		t.Error("New task should have nil CompletedAt")
	}

	// Выполняем задачу
	err := storage.Complete(task.ID)

	if err != nil {
		t.Fatalf("Complete() failed: %v", err)
	}

	// Проверяем состояние после выполнения
	if !task.Completed {
		t.Error("Task should be completed after Complete()")
	}
	if task.CompletedAt == nil {
		t.Error("CompletedAt should be set after Complete()")
	}

	// Проверяем что время завершения близко к текущему
	if time.Since(*task.CompletedAt) > time.Second {
		t.Error("CompletedAt is too far in the past")
	}

	// Проверяем что задача обновилась в хранилище
	retrieved, _ := storage.GetByID(task.ID)
	if !retrieved.Completed {
		t.Error("Task in storage should be completed")
	}
}

func TestStorage_Complete_AlreadyCompleted(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	task := storage.Add("Test Task", "Description")

	// Первое выполнение
	err1 := storage.Complete(task.ID)
	if err1 != nil {
		t.Fatalf("First Complete() failed: %v", err1)
	}

	// Сохраняем время первого завершения
	firstCompleteTime := *task.CompletedAt

	// Даем небольшую задержку
	time.Sleep(10 * time.Millisecond)

	// Второе выполнение (должно быть ошибка)
	err2 := storage.Complete(task.ID)

	if err2 == nil {
		t.Error("Expected error for already completed task, got nil")
	}
	expectedErr := "Задача уже выполнена"
	if err2.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err2.Error())
	}

	// Проверяем что статус и время не изменились
	if !task.Completed {
		t.Error("Task should remain completed")
	}
	if *task.CompletedAt != firstCompleteTime {
		t.Error("CompletedAt should not change")
	}
}

func TestStorage_Complete_Errors(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()
	storage.Add("Task", "Description") // ID = 1

	tests := []struct {
		name        string
		id          int
		expectedErr string
	}{
		{
			name:        "id = 0",
			id:          0,
			expectedErr: "Неккоректный ID",
		},
		{
			name:        "negative id",
			id:          -1,
			expectedErr: "Неккоректный ID",
		},
		{
			name:        "non-existent id",
			id:          999,
			expectedErr: "Несуществующий ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Complete(tt.id)

			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestStorage_Complete_EmptyStorage(t *testing.T) {
	t.Parallel()

	storage := memory.NewStorage()

	err := storage.Complete(1)

	if err == nil {
		t.Error("Expected error for empty storage, got nil")
	}
	expectedErr := "Список задач пуст"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestGenerateDefaultName_Format(t *testing.T) {
	t.Parallel()

	// Проверяем несколько раз для разных случайных чисел
	for i := 0; i < 100; i++ {
		// Вызываем через добавление задачи с пустым названием
		storage := memory.NewStorage()
		task := storage.Add("", "")

		// Проверяем формат
		name := task.Title

		// Длина должна быть 11 (exr- + 7 цифр)
		if len(name) != 11 {
			t.Errorf("Iteration %d: Expected length 11, got %d for name %q", i, len(name), name)
		}

		// Проверяем префикс
		if name[:4] != "exr-" {
			t.Errorf("Iteration %d: Expected prefix 'exr-', got %q", i, name[:4])
		}

		// Проверяем что остальные символы - цифры
		for j := 4; j < len(name); j++ {
			if name[j] < '0' || name[j] > '9' {
				t.Errorf("Iteration %d: Character %c at position %d is not a digit", i, name[j], j)
			}
		}
	}
}

// Тест на параллельное выполнение (если понадобится)
func TestStorage_Concurrent(t *testing.T) {
	storage := memory.NewStorage()

	// Запускаем несколько горутин для параллельного добавления задач
	const numTasks = 10
	done := make(chan bool)

	for i := 0; i < numTasks; i++ {
		go func(index int) {
			title := ""
			if index%2 == 0 {
				title = "Custom Task"
			}
			task := storage.Add(title, "Description")

			// Базовая проверка
			if task.ID < 1 || task.ID > numTasks {
				t.Errorf("Invalid task ID: %d", task.ID)
			}
			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < numTasks; i++ {
		<-done
	}

	// Проверяем что все задачи добавились
	if storage.IsEmpty() {
		t.Error("Storage should not be empty after concurrent adds")
	}

	tasks, err := storage.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(tasks) != numTasks {
		t.Errorf("Expected %d tasks, got %d", numTasks, len(tasks))
	}
}
