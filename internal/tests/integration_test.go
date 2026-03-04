package tests

import (
	"os"
	"path/filepath"
	"testing"

	"TaskTracker/internal/handler"
	"TaskTracker/internal/repository/memory"
	"TaskTracker/internal/repository/sjson"
	"TaskTracker/internal/service"
)

// TestIntegration_MemoryRepository_Service_Handler проверяет интеграцию всех компонентов с Memory хранилищем
func TestIntegration_MemoryRepository_Service_Handler(t *testing.T) {
	t.Parallel()

	// Инициализируем компоненты
	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Fatal("Handler creation failed")
	}

	// Тестируем добавление
	task1, err := svc.AddTask("Project Setup", "Initialize the project", []string{"work", "setup"})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// Тестируем получение
	retrieved, err := svc.GetTaskById(task1.ID)
	if err != nil {
		t.Fatalf("GetTaskById failed: %v", err)
	}

	if retrieved.Title != "Project Setup" {
		t.Errorf("Retrieved task title mismatch: got %q", retrieved.Title)
	}

	// Тестируем поиск по тегу
	tasks, err := svc.GetTasksByTag("work")
	if err != nil {
		t.Fatalf("GetTasksByTag failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 work task, got %d", len(tasks))
	}

	// Тестируем завершение
	completed, err := svc.CompleteTask(task1.ID)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	if !completed.Completed {
		t.Error("Task should be completed")
	}
}

// TestIntegration_JSONRepository_Service_Handler проверяет интеграцию с JSON хранилищем
func TestIntegration_JSONRepository_Service_Handler(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp("", "test_integration_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	os.Remove(tmpFile.Name())

	// Инициализируем компоненты с JSON хранилищем
	repo, err := sjson.NewJSONStorage(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewJSONStorage failed: %v", err)
	}

	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Fatal("Handler creation failed")
	}

	// Добавляем задачу
	task, _ := svc.AddTask("Database Setup", "Configure database connection", []string{"db", "setup"})

	// Закрываем первый handler
	_ = h

	// Создаем новый repo из файла (загружаем сохраненные данные)
	repo2, err := sjson.NewJSONStorage(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewJSONStorage failed: %v", err)
	}

	svc2 := service.NewNewTaskService(repo2)

	// Проверяем что данные загрузились
	retrieved, err := svc2.GetTaskById(task.ID)
	if err != nil {
		t.Fatalf("GetTaskById failed: %v", err)
	}

	if retrieved.Title != "Database Setup" {
		t.Errorf("Expected 'Database Setup', got %q", retrieved.Title)
	}
}

// TestIntegration_MultipleServices проверяет несколько сервисов с одним хранилищем
func TestIntegration_MultipleServices(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()

	// Создаем два сервиса с одним хранилищем
	svc1 := service.NewNewTaskService(repo)
	svc2 := service.NewNewTaskService(repo)

	// Добавляем через первый сервис
	svc1.AddTask("Task 1", "Description 1", []string{})

	// Получаем через второй сервис
	tasks, err := svc2.GetAllTasks()

	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	// Добавляем через второй сервис
	svc2.AddTask("Task 2", "Description 2", []string{})

	// Получаем через первый сервис
	tasks, err = svc1.GetAllTasks()

	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

// TestIntegration_ComplexWorkflow проверяет сложный workflow
func TestIntegration_ComplexWorkflow(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)
	h := handler.NewCLIHandler(svc)

	if h == nil {
		t.Fatal("Handler creation failed")
	}

	// Этап 1: Добавляем несколько задач
	taskData := []struct {
		title string
		desc  string
		tags  []string
	}{
		{"Buy groceries", "Milk, bread, eggs", []string{"shopping", "urgent"}},
		{"Write report", "Quarterly report", []string{"work"}},
		{"Exercise", "Morning run", []string{"health", "morning"}},
		{"Call doctor", "Schedule appointment", []string{"health", "important"}},
		{"Team meeting", "Project discussion", []string{"work", "meeting"}},
	}

	for _, td := range taskData {
		svc.AddTask(td.title, td.desc, td.tags)
	}

	// Этап 2: Получаем все задачи
	allTasks, err := svc.GetAllTasks()
	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}

	if len(allTasks) != 5 {
		t.Errorf("Expected 5 tasks, got %d", len(allTasks))
	}

	// Этап 3: Ищем по тегам
	workTasks, _ := svc.GetTasksByTag("work")
	if len(workTasks) != 2 {
		t.Errorf("Expected 2 work tasks, got %d", len(workTasks))
	}

	healthTasks, _ := svc.GetTasksByTag("health")
	if len(healthTasks) != 2 {
		t.Errorf("Expected 2 health tasks, got %d", len(healthTasks))
	}

	shoppingTasks, _ := svc.GetTasksByTag("shopping")
	if len(shoppingTasks) != 1 {
		t.Errorf("Expected 1 shopping task, got %d", len(shoppingTasks))
	}

	// Этап 4: Завершаем несколько задач
	for _, task := range allTasks {
		if task.Tags != nil {
			for _, tag := range task.Tags {
				if tag == "urgent" {
					svc.CompleteTask(task.ID)
					break
				}
			}
		}
	}

	// Этап 5: Проверяем что правильные задачи завершены
	for _, task := range allTasks {
		retrieved, _ := svc.GetTaskById(task.ID)

		isUrgent := false
		if retrieved.Tags != nil {
			for _, tag := range retrieved.Tags {
				if tag == "urgent" {
					isUrgent = true
					break
				}
			}
		}

		if isUrgent && !retrieved.Completed {
			t.Errorf("Urgent task %d should be completed", task.ID)
		} else if !isUrgent && retrieved.Completed {
			t.Errorf("Non-urgent task %d should not be completed", task.ID)
		}
	}
}

// TestIntegration_JSONPersistence проверяет сохранение данных в JSON
func TestIntegration_JSONPersistence(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "test_persistence_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "tasks.json")

	// Фаза 1: Создание и добавление данных
	{
		repo, _ := sjson.NewJSONStorage(filePath)
		svc := service.NewNewTaskService(repo)

		svc.AddTask("Task 1", "Description 1", []string{"tag1"})
		svc.AddTask("Task 2", "Description 2", []string{"tag2"})
		svc.AddTask("Task 3", "Description 3", []string{"tag1", "tag2"})

		completed, _ := svc.AddTask("Task 4", "Description 4", []string{})
		svc.CompleteTask(completed.ID)
	}

	// Проверяем что файл существует
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("JSON file should have been created")
	}

	// Фаза 2: Загрузка данных из файла
	{
		repo, _ := sjson.NewJSONStorage(filePath)
		svc := service.NewNewTaskService(repo)

		tasks, _ := svc.GetAllTasks()
		if len(tasks) != 4 {
			t.Errorf("Expected 4 tasks after reload, got %d", len(tasks))
		}

		tag1Tasks, _ := svc.GetTasksByTag("tag1")
		if len(tag1Tasks) != 2 {
			t.Errorf("Expected 2 tag1 tasks, got %d", len(tag1Tasks))
		}

		completed, _ := svc.GetTaskById(4)
		if !completed.Completed {
			t.Error("Task 4 should be completed after reload")
		}
	}
}

// TestIntegration_ConcurrentOperations проверяет параллельные операции через интеграцию
func TestIntegration_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	const numGoroutines = 50
	numTasksPerGoroutine := 10

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer func() { done <- true }()
			for j := 0; j < numTasksPerGoroutine; j++ {
				svc.AddTask("Task", "Description", []string{})
			}
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Проверяем что все задачи добавились
	tasks, err := svc.GetAllTasks()
	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}

	expectedCount := numGoroutines * numTasksPerGoroutine
	if len(tasks) != expectedCount {
		t.Errorf("Expected %d tasks, got %d", expectedCount, len(tasks))
	}
}

// TestIntegration_ServiceWithDifferentStorages проверяет сервис с разными хранилищами
func TestIntegration_ServiceWithDifferentStorages(t *testing.T) {
	t.Parallel()

	memRepo := memory.NewStorage()
	memSvc := service.NewNewTaskService(memRepo)

	// Добавляем в Memory
	memSvc.AddTask("Memory Task", "Description", []string{})

	// Проверяем в Memory
	memTasks, _ := memSvc.GetAllTasks()
	if len(memTasks) != 1 {
		t.Errorf("Expected 1 memory task, got %d", len(memTasks))
	}

	// Создаем JSON хранилище
	tmpFile, _ := os.CreateTemp("", "test_*.json")
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	os.Remove(tmpFile.Name())

	jsonRepo, err := sjson.NewJSONStorage(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewJSONStorage failed: %v", err)
	}
	jsonSvc := service.NewNewTaskService(jsonRepo)

	// Добавляем в JSON
	if _, err := jsonSvc.AddTask("JSON Task", "Description", []string{}); err != nil {
		t.Fatalf("AddTask to JSON storage failed: %v", err)
	}

	// Проверяем в JSON
	jsonTasks, _ := jsonSvc.GetAllTasks()
	if len(jsonTasks) != 1 {
		t.Errorf("Expected 1 JSON task, got %d", len(jsonTasks))
	}

	// Проверяем что данные разделены (Memory != JSON)
	if len(memTasks) == len(jsonTasks) {
		memTitle := memTasks[0].Title
		jsonTitle := jsonTasks[0].Title

		if memTitle == jsonTitle {
			t.Error("Memory and JSON storages should have different data")
		}
	}
}

// TestIntegration_FullApplicationSimulation проверяет полное приложение
func TestIntegration_FullApplicationSimulation(t *testing.T) {
	t.Parallel()

	tmpDir, _ := os.MkdirTemp("", "test_app_*")
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "app_tasks.json")

	// Запускаем приложение (первый раз)
	{
		repo, _ := sjson.NewJSONStorage(filePath)
		svc := service.NewNewTaskService(repo)
		h := handler.NewCLIHandler(svc)

		if h == nil {
			t.Fatal("Handler should be created")
		}

		// Пользователь добавляет задачи
		svc.AddTask("Learn Go", "Complete Go tutorials", []string{"learning", "programming"})
		svc.AddTask("Build CLI", "Create a CLI application", []string{"programming", "project"})
		svc.AddTask("Write tests", "Add unit tests", []string{"testing", "programming"})
	}

	// Перезапускаем приложение (данные должны загрузиться)
	{
		repo, _ := sjson.NewJSONStorage(filePath)
		svc := service.NewNewTaskService(repo)
		h := handler.NewCLIHandler(svc)

		if h == nil {
			t.Fatal("Handler should be created on restart")
		}

		// Проверяем что данные загрузились
		tasks, _ := svc.GetAllTasks()
		if len(tasks) != 3 {
			t.Errorf("Expected 3 tasks on restart, got %d", len(tasks))
		}

		// Пользователь находит задачи по тегу
		progTasks, _ := svc.GetTasksByTag("programming")
		if len(progTasks) != 3 {
			t.Errorf("Expected 3 programming tasks, got %d", len(progTasks))
		}

		// Пользователь завершает задачу
		svc.CompleteTask(1)

		// Проверяем что задача завершена
		task1, _ := svc.GetTaskById(1)
		if !task1.Completed {
			t.Error("Task 1 should be completed")
		}
	}

	// Третий запуск - проверяем что статус сохранился
	{
		repo, _ := sjson.NewJSONStorage(filePath)
		svc := service.NewNewTaskService(repo)

		task1, _ := svc.GetTaskById(1)
		if !task1.Completed {
			t.Error("Task 1 should remain completed after restart")
		}

		task2, _ := svc.GetTaskById(2)
		if task2.Completed {
			t.Error("Task 2 should not be completed")
		}
	}
}

// TestIntegration_LargeDataSet проверяет работу с большим количеством данных
func TestIntegration_LargeDataSet(t *testing.T) {
	repo := memory.NewStorage()
	svc := service.NewNewTaskService(repo)

	const numberOfTasks = 1000

	// Добавляем много задач
	for i := 0; i < numberOfTasks; i++ {
		svc.AddTask("Task"+string(rune(i%1000)), "Description", []string{"bulk"})
	}

	// Получаем все задачи
	tasks, err := svc.GetAllTasks()
	if err != nil {
		t.Fatalf("GetAllTasks failed: %v", err)
	}

	if len(tasks) != numberOfTasks {
		t.Errorf("Expected %d tasks, got %d", numberOfTasks, len(tasks))
	}

	// Находим по тегу
	bulkTasks, _ := svc.GetTasksByTag("bulk")
	if len(bulkTasks) != numberOfTasks {
		t.Errorf("Expected %d bulk tasks, got %d", numberOfTasks, len(bulkTasks))
	}

	// Завершаем половину
	for i := 1; i <= numberOfTasks/2; i++ {
		svc.CompleteTask(i)
	}

	// Проверяем что завершилось правильно
	task := tasks[0]
	if task.ID <= numberOfTasks/2 {
		retrieved, _ := svc.GetTaskById(task.ID)
		if !retrieved.Completed {
			t.Error("Task should be completed in large dataset")
		}
	}
}
