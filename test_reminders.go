package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"TaskTracker/config"
	"TaskTracker/internal/repository/database"
	"TaskTracker/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	repo, err := database.NewPostgresRepo(cfg.GetDSN())
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer repo.Close()

	taskSvc := service.NewTaskService(repo)
	extSvc := service.NewPostgresTaskService(repo)

	ctx := context.Background()

	// Test 1: Add task without reminder
	fmt.Println("🧪 Test 1: Add task without reminder")
	task1, err := taskSvc.AddTask(ctx, 1, "Test task 1", "Description", []string{"test"})
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success: Task #%d created\n", task1.UserTaskID)
	}

	// Test 2: Add task with reminder
	fmt.Println("\n🧪 Test 2: Add task with reminder")
	dueDateStr := time.Now().Add(1 * time.Hour).Format("02.01.2006 15:04")
	task2, err := extSvc.AddTaskWithReminder(ctx, 1, "Test task 2", "Description", []string{"test"}, &dueDateStr, "30m")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success: Task #%d created with reminder\n", task2.UserTaskID)
	}

	// Test 3: Get all tasks
	fmt.Println("\n🧪 Test 3: Get all tasks")
	tasks, err := taskSvc.GetAllTasks(ctx, 1)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success: Found %d tasks\n", len(tasks))
		for _, t := range tasks {
			fmt.Printf("  - #%d: %s\n", t.UserTaskID, t.Title)
		}
	}

	// Test 4: Get tasks for reminder
	fmt.Println("\n🧪 Test 4: Get tasks for reminder")
	reminderTasks, err := extSvc.GetTasksForReminder(ctx)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Success: Found %d tasks for reminder\n", len(reminderTasks))
	}

	fmt.Println("\n🎉 All tests completed!")
}
