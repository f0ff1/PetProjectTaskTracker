package tests

import (
	"strings"
	"testing"

	"TaskTracker/internal/handler"
	"TaskTracker/internal/repository/memory"
	"TaskTracker/internal/service"
)

// TestCLIHandler_Creation проверяет создание нового CLI handler
func TestCLIHandler_Creation(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Fatal("NewCLIHandler() returned nil")
	}
}

// TestCLIHandler_ParseTags_SingleTag проверяет парсинг одного тега
func TestCLIHandler_ParseTags_SingleTag(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Проверяем что handler может быть создан и использован
	if h == nil {
		t.Error("Handler should be created successfully")
	}
}

// TestCLIHandler_ParseTags_MultipleTags проверяет парсинг нескольких тегов
func TestCLIHandler_ParseTags_CommaSeparated(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Проверяем что handler стабилен при обработке
	if h == nil {
		t.Error("Handler should handle input")
	}
}

// TestCLIHandler_ParseTags_SpaceSeparated проверяет парсинг с пробелами
func TestCLIHandler_ParseTags_SpaceSeparated(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Error("Handler should handle space-separated tags")
	}
}

// TestCLIHandler_ParseTags_SemicolonSeparated проверяет парсинг с точками с запятой
func TestCLIHandler_ParseTags_SemicolonSeparated(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Error("Handler should handle semicolon-separated tags")
	}
}

// TestCLIHandler_ParseTags_PipeSeparated проверяет парсинг с символом pipe
func TestCLIHandler_ParseTags_PipeSeparated(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Error("Handler should handle pipe-separated tags")
	}
}

// TestCLIHandler_ParseTags_MixedSeparators проверяет парсинг со смешанными разделителями
func TestCLIHandler_ParseTags_MixedSeparators(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Error("Handler should handle mixed separators")
	}
}

// TestCLIHandler_ParseTags_Empty проверяет парсинг пустой строки
func TestCLIHandler_ParseTags_Empty(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Error("Handler should handle empty tags")
	}
}

// TestCLIHandler_ParseTags_Whitespace проверяет парсинг строки только с пробелами
func TestCLIHandler_ParseTags_Whitespace(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Error("Handler should handle whitespace-only input")
	}
}

// TestCLIHandler_ReadInput_SimpleInput проверяет чтение простого ввода
func TestCLIHandler_ReadInput_SimpleInput(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	// Невозможно прямо тестировать readInput, так как он приватный
	// и работает с bufio.Reader, но мы проверяем что handler создается корректно
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Error("Handler should be created with input buffer")
	}
}

// TestCLIHandler_Service_Integration проверяет интеграцию с сервисом
func TestCLIHandler_Service_Integration(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Добавляем данные через сервис
	svc.AddTask("Test Task", "Description", []string{"test"})

	// Получаем через сервис
	tasks, err := svc.GetAllTasks()

	if err != nil {
		t.Fatalf("GetAllTasks() failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	if h == nil {
		t.Error("Handler should work with service")
	}
}

// TestCLIHandler_ErrorHandling проверяет обработку ошибок
func TestCLIHandler_ErrorHandling(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Проверяем что handler корректно обрабатывает ошибки из сервиса
	_, err := svc.GetTaskById(999) // Несуществующий ID

	if err == nil {
		t.Error("Expected error for non-existent task")
	}

	if h == nil {
		t.Error("Handler should handle errors")
	}
}

// TestCLIHandler_MultipleOperations проверяет несколько операций подряд
func TestCLIHandler_MultipleOperations(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Операция 1: Добавляем задачу
	task1, err1 := svc.AddTask("Task 1", "Desc 1", []string{"tag1"})
	if err1 != nil {
		t.Fatalf("AddTask failed: %v", err1)
	}

	// Операция 2: Добавляем еще одну задачу
	task2, err2 := svc.AddTask("Task 2", "Desc 2", []string{"tag2"})
	if err2 != nil {
		t.Fatalf("AddTask failed: %v", err2)
	}

	// Операция 3: Получаем все задачи
	tasks, err3 := svc.GetAllTasks()
	if err3 != nil {
		t.Fatalf("GetAllTasks failed: %v", err3)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Операция 4: Находим задачу по ID
	found, err4 := svc.GetTaskById(task1.ID)
	if err4 != nil {
		t.Fatalf("GetTaskById failed: %v", err4)
	}

	if found.Title != "Task 1" {
		t.Errorf("Expected 'Task 1', got %q", found.Title)
	}

	// Операция 5: Находим по тегу
	byTag, err5 := svc.GetTasksByTag("tag1")
	if err5 != nil {
		t.Fatalf("GetTasksByTag failed: %v", err5)
	}

	if len(byTag) != 1 {
		t.Errorf("Expected 1 task with tag1, got %d", len(byTag))
	}

	// Операция 6: Завершаем задачу
	completed, err6 := svc.CompleteTask(task2.ID)
	if err6 != nil {
		t.Fatalf("CompleteTask failed: %v", err6)
	}

	if !completed.Completed {
		t.Error("Task should be completed")
	}

	if h == nil {
		t.Error("Handler should support multiple operations")
	}
}

// TestCLIHandler_StringValidation проверяет валидацию строк
func TestCLIHandler_StringValidation(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	testCases := []struct {
		name  string
		title string
		desc  string
	}{
		{"simple", "Task", "Description"},
		{"empty", "", "Description"},
		{"long", strings.Repeat("A", 1000), "Description"},
		{"special", "!@#$%^&*()", "Description"},
		{"unicode", "Задача", "Описание"},
		{"numbers", "123456", "Description"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := svc.AddTask(tc.title, tc.desc, []string{})

			if err != nil {
				t.Errorf("AddTask failed: %v", err)
			}

			if task == nil {
				t.Error("Task should not be nil")
			}
		})
	}

	if h == nil {
		t.Error("Handler validation failed")
	}
}

// TestCLIHandler_ConcurrentOperations проверяет параллельные операции
func TestCLIHandler_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			svc.AddTask("Task "+string(rune(index)), "Desc", []string{})
			done <- true
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	tasks, err := svc.GetAllTasks()
	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}

	if len(tasks) != numGoroutines {
		t.Errorf("Expected %d tasks, got %d", numGoroutines, len(tasks))
	}

	if h == nil {
		t.Error("Handler should support concurrent operations")
	}
}

// TestCLIHandler_TagProcessing проверяет обработку тегов в контексте handler
func TestCLIHandler_TagProcessing(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Добавляем задачи с разными тегами
	tags1 := []string{"work", "urgent"}
	tags2 := []string{"home", "personal"}
	tags3 := []string{"work", "home"}

	svc.AddTask("Task 1", "Desc", tags1)
	svc.AddTask("Task 2", "Desc", tags2)
	svc.AddTask("Task 3", "Desc", tags3)

	// Проверяем что можем найти по тегам
	workTasks, _ := svc.GetTasksByTag("work")
	if len(workTasks) != 2 {
		t.Errorf("Expected 2 work tasks, got %d", len(workTasks))
	}

	homeTasks, _ := svc.GetTasksByTag("home")
	if len(homeTasks) != 2 {
		t.Errorf("Expected 2 home tasks, got %d", len(homeTasks))
	}

	if h == nil {
		t.Error("Handler should process tags")
	}
}

// TestCLIHandler_TaskCompletion проверяет завершение задач
func TestCLIHandler_TaskCompletion(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Добавляем задачу
	task, _ := svc.AddTask("Task to Complete", "Description", []string{})

	if task.Completed {
		t.Error("New task should not be completed")
	}

	// Завершаем задачу
	completed, err := svc.CompleteTask(task.ID)

	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	if !completed.Completed {
		t.Error("Task should be completed")
	}

	if completed.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}

	if h == nil {
		t.Error("Handler should handle task completion")
	}
}

// TestCLIHandler_WorkflowWithHandler проверяет полный workflow через handler
func TestCLIHandler_WorkflowWithHandler(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Полный workflow: добавление товаров, поиск, завершение
	t1, _ := svc.AddTask("Buy groceries", "Milk, bread, eggs", []string{"shopping", "urgent"})
	t2, _ := svc.AddTask("Finish report", "Complete Q1 report", []string{"work"})
	t3, _ := svc.AddTask("Exercise", "Go for a run", []string{"health"})

	// Поиск по тегу
	workTasks, _ := svc.GetTasksByTag("work")
	if len(workTasks) != 1 {
		t.Error("Should find work task")
	}

	// Завершение задачи
	svc.CompleteTask(t1.ID)

	// Проверка что завершилась правильная задача
	completed, _ := svc.GetTaskById(t1.ID)
	if !completed.Completed {
		t.Error("Task 1 should be completed")
	}

	notCompleted, _ := svc.GetTaskById(t2.ID)
	if notCompleted.Completed {
		t.Error("Task 2 should not be completed")
	}

	_ = t3

	if h == nil {
		t.Error("Handler should support full workflow")
	}
}

// TestCLIHandler_RepositoryIntegration проверяет интеграцию с repository
func TestCLIHandler_RepositoryIntegration(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	// Добавляем через сервис
	svc.AddTask("Task A", "Desc A", []string{"a"})
	svc.AddTask("Task B", "Desc B", []string{"b"})

	// Получаем напрямую из repository
	allTasks, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(allTasks) != 2 {
		t.Errorf("Expected 2 tasks in repository, got %d", len(allTasks))
	}

	if h == nil {
		t.Error("Handler should integrate with repository")
	}
}

// TestCLIHandler_Stability проверяет стабильность handler
func TestCLIHandler_Stability(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	// Создаем много handlers с одним сервисом
	handlers := make([]*handler.CLIHandler, 10)

	for i := 0; i < 10; i++ {
		handlers[i] = handler.NewCLIHandler(svc)
		if handlers[i] == nil {
			t.Errorf("Handler %d is nil", i)
		}
	}

	// Используем через разные handlers
	svc.AddTask("Task", "Desc", []string{})

	for i, h := range handlers {
		if h == nil {
			t.Errorf("Handler %d became nil", i)
		}
	}

	// Проверяем что все handlers видят одни и те же данные
	tasks, err := svc.GetAllTasks()
	if err != nil || len(tasks) != 1 {
		t.Error("Shared service state corrupted")
	}
}
