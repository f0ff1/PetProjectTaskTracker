package tests

import (
	"testing"

	customError "TaskTracker/errors"
	"TaskTracker/internal/repository/memory"
	"TaskTracker/internal/service"
)

// TestTaskService_NewService проверяет создание нового сервиса
func TestTaskService_NewService(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	if svc == nil {
		t.Fatal("NewNewTaskService() returned nil")
	}
}

// TestTaskService_AddTask_WithTitle проверяет добавление задачи с названием
func TestTaskService_AddTask_WithTitle(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	title := "Test Task"
	desc := "Test Description"
	tags := []string{"test"}

	task, err := svc.AddTask(title, desc, tags)

	if err != nil {
		t.Fatalf("AddTask() failed: %v", err)
	}

	if task.Title != title {
		t.Errorf("Expected Title = %q, got %q", title, task.Title)
	}

	if task.Description != desc {
		t.Errorf("Expected Description = %q, got %q", desc, task.Description)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", task.ID)
	}
}

// TestTaskService_AddTask_WithoutTitle проверяет добавление задачи без названия
func TestTaskService_AddTask_WithoutTitle(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	task, err := svc.AddTask("", "Description", []string{})

	if err != nil {
		t.Fatalf("AddTask() failed: %v", err)
	}

	// Должно быть сгенерировано название
	if task.Title == "" {
		t.Error("Expected auto-generated title, got empty string")
	}

	// Проверяем что название начинается с "def-name-exr-"
	if len(task.Title) < 13 {
		t.Errorf("Generated title is too short: %q", task.Title)
	}
}

// TestTaskService_AddTask_MultipleTasks проверяет добавление нескольких задач
func TestTaskService_AddTask_MultipleTasks(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	task1, _ := svc.AddTask("Task 1", "Desc 1", []string{})
	task2, _ := svc.AddTask("Task 2", "Desc 2", []string{})
	task3, _ := svc.AddTask("", "Desc 3", []string{}) // С автогенерацией

	if task1.ID != 1 {
		t.Errorf("Task1 ID = %d, want 1", task1.ID)
	}
	if task2.ID != 2 {
		t.Errorf("Task2 ID = %d, want 2", task2.ID)
	}
	if task3.ID != 3 {
		t.Errorf("Task3 ID = %d, want 3", task3.ID)
	}
}

// TestTaskService_GetAllTasks проверяет получение всех задач
func TestTaskService_GetAllTasks(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	svc.AddTask("Task 1", "Desc 1", []string{})
	svc.AddTask("Task 2", "Desc 2", []string{})

	tasks, err := svc.GetAllTasks()

	if err != nil {
		t.Fatalf("GetAllTasks() failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

// TestTaskService_GetAllTasks_Empty проверяет получение всех задач из пустого хранилища
func TestTaskService_GetAllTasks_Empty(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	tasks, err := svc.GetAllTasks()

	if err != nil {
		t.Errorf("Expected no error for empty storage, got %v", err)
	}

	if tasks != nil && len(tasks) != 0 {
		t.Error("Expected empty tasks for empty storage")
	}
}

// TestTaskService_GetTaskById_Success проверяет получение задачи по ID
func TestTaskService_GetTaskById_Success(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	added, _ := svc.AddTask("Test Task", "Description", []string{})

	task, err := svc.GetTaskById(added.ID)

	if err != nil {
		t.Fatalf("GetTaskById() failed: %v", err)
	}

	if task.ID != added.ID {
		t.Errorf("Expected ID = %d, got %d", added.ID, task.ID)
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected Title = 'Test Task', got %q", task.Title)
	}
}

// TestTaskService_GetTaskById_InvalidID проверяет получение задачи с неправильным ID
func TestTaskService_GetTaskById_InvalidID(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	tests := []struct {
		name      string
		id        int
		expectErr error
	}{
		{"zero id", 0, customError.ErrIdNotExists},
		{"negative id", -1, customError.ErrIdNotExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := svc.GetTaskById(tt.id)

			if err == nil {
				t.Error("Expected error, got nil")
			}

			if err != tt.expectErr {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if task != nil {
				t.Errorf("Expected nil task, got %v", task)
			}
		})
	}
}

// TestTaskService_GetTaskById_NotFound проверяет получение несуществующей задачи
func TestTaskService_GetTaskById_NotFound(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	task, err := svc.GetTaskById(999)

	if err == nil {
		t.Error("Expected error for non-existent ID, got nil")
	}

	if task != nil {
		t.Errorf("Expected nil task, got %v", task)
	}
}

// TestTaskService_GetTasksByTag проверяет получение задач по тегу
func TestTaskService_GetTasksByTag(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	svc.AddTask("Task 1", "Desc 1", []string{"work", "urgent"})
	svc.AddTask("Task 2", "Desc 2", []string{"personal"})
	svc.AddTask("Task 3", "Desc 3", []string{"work"})

	tasks, err := svc.GetTasksByTag("work")

	if err != nil {
		t.Fatalf("GetTasksByTag() failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks with 'work' tag, got %d", len(tasks))
	}
}

// TestTaskService_GetTasksByTag_InvalidTag проверяет получение задач с пустым тегом
func TestTaskService_GetTasksByTag_InvalidTag(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	svc.AddTask("Task", "Description", []string{"tag"})

	tasks, err := svc.GetTasksByTag("")

	if err == nil {
		t.Error("Expected error for empty tag, got nil")
	}

	if err != customError.ErrWrongTag {
		t.Errorf("Expected ErrWrongTag, got %v", err)
	}

	if tasks != nil {
		t.Errorf("Expected nil tasks, got %v", tasks)
	}
}

// TestTaskService_GetTasksByTag_NotFound проверяет поиск несуществующего тега
func TestTaskService_GetTasksByTag_NotFound(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	svc.AddTask("Task", "Description", []string{"tag"})

	tasks, err := svc.GetTasksByTag("nonexistent")

	if err != nil {
		t.Fatalf("GetTasksByTag() failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

// TestTaskService_CompleteTask_Success проверяет отметить задачу как выполненную
func TestTaskService_CompleteTask_Success(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	added, _ := svc.AddTask("Task", "Description", []string{})

	completed, err := svc.CompleteTask(added.ID)

	if err != nil {
		t.Fatalf("CompleteTask() failed: %v", err)
	}

	if !completed.Completed {
		t.Error("Task should be completed")
	}

	if completed.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}
}

// TestTaskService_CompleteTask_InvalidID проверяет завершение с неправильным ID
func TestTaskService_CompleteTask_InvalidID(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	tests := []struct {
		name      string
		id        int
		expectErr error
	}{
		{"zero id", 0, customError.ErrWrongTypeID},
		{"negative id", -1, customError.ErrWrongTypeID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := svc.CompleteTask(tt.id)

			if err == nil {
				t.Error("Expected error, got nil")
			}

			if err != tt.expectErr {
				t.Errorf("Expected error %v, got %v", tt.expectErr, err)
			}

			if task != nil {
				t.Errorf("Expected nil task, got %v", task)
			}
		})
	}
}

// TestTaskService_CompleteTask_AlreadyCompleted проверяет повторное завершение
func TestTaskService_CompleteTask_AlreadyCompleted(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	added, _ := svc.AddTask("Task", "Description", []string{})

	svc.CompleteTask(added.ID)

	_, err := svc.CompleteTask(added.ID)

	if err == nil {
		t.Error("Expected error for already completed task, got nil")
	}

	if err != customError.ErrTaskAlredyComplete {
		t.Errorf("Expected ErrTaskAlredyComplete, got %v", err)
	}
}

// TestTaskService_CompleteTask_NotFound проверяет завершение несуществующей задачи
func TestTaskService_CompleteTask_NotFound(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	task, err := svc.CompleteTask(999)

	if err == nil {
		t.Error("Expected error for non-existent task, got nil")
	}

	if task != nil {
		t.Errorf("Expected nil task, got %v", task)
	}
}

// TestTaskService_WorkflowScenario проверяет типичный workflow
func TestTaskService_WorkflowScenario(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	// Step 1: Добавляем несколько задач
	task1, _ := svc.AddTask("Buy milk", "Fresh milk", []string{"shopping", "urgent"})
	task2, _ := svc.AddTask("Finish project", "Complete the feature", []string{"work"})
	task3, _ := svc.AddTask("", "Go for a run", []string{"health"}) // Auto-generated title

	if task1.ID != 1 || task2.ID != 2 || task3.ID != 3 {
		t.Error("Task IDs are not sequential")
	}

	// Step 2: Получаем все задачи
	allTasks, _ := svc.GetAllTasks()
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(allTasks))
	}

	// Step 3: Ищем задачи по тегу
	workTasks, _ := svc.GetTasksByTag("work")
	if len(workTasks) != 1 {
		t.Errorf("Expected 1 work task, got %d", len(workTasks))
	}

	// Step 4: Извлекаем задачу по ID
	retrieved, _ := svc.GetTaskById(1)
	if retrieved.Title != "Buy milk" {
		t.Errorf("Expected 'Buy milk', got %q", retrieved.Title)
	}

	// Step 5: Завершаем задачу
	completed, _ := svc.CompleteTask(1)
	if !completed.Completed {
		t.Error("Task should be completed")
	}

	// Step 6: Проверяем что завершенная задача не появляется как незавершенная
	retrieved2, _ := svc.GetTaskById(1)
	if !retrieved2.Completed {
		t.Error("Retrieved task should be completed")
	}
}

// TestTaskService_AddTask_WithManyTags проверяет добавление задачи с большим количеством тегов
func TestTaskService_AddTask_WithManyTags(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	tags := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}

	task, err := svc.AddTask("Task", "Description", tags)

	if err != nil {
		t.Fatalf("AddTask() failed: %v", err)
	}

	if len(task.Tags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(task.Tags))
	}

	// Проверяем что можем найти по каждому тегу
	for _, tag := range tags {
		foundTasks, _ := svc.GetTasksByTag(tag)
		if len(foundTasks) != 1 {
			t.Errorf("Expected 1 task with tag %q, got %d", tag, len(foundTasks))
		}
	}
}

// TestTaskService_AddTask_EmptyDescription проверяет добавление задачи с пустым описанием
func TestTaskService_AddTask_EmptyDescription(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	task, err := svc.AddTask("Task", "", []string{})

	if err != nil {
		t.Fatalf("AddTask() failed: %v", err)
	}

	if task.Description != "" {
		t.Error("Description should be empty")
	}
}

// TestTaskService_GetAllTasks_Order проверяет порядок задач
func TestTaskService_GetAllTasks_Order(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	task1, _ := svc.AddTask("First", "First", []string{})
	task2, _ := svc.AddTask("Second", "Second", []string{})
	task3, _ := svc.AddTask("Third", "Third", []string{})

	tasks, _ := svc.GetAllTasks()

	// Проверяем что все ID присутствуют
	ids := make(map[int]bool)
	for _, task := range tasks {
		ids[task.ID] = true
	}

	if !ids[task1.ID] || !ids[task2.ID] || !ids[task3.ID] {
		t.Error("Not all tasks are present in the list")
	}
}
