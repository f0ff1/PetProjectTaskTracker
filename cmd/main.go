package main

import (
	"fmt"

	"TaskTracker/internal/handler"
	"TaskTracker/internal/repository"
	"TaskTracker/internal/repository/memory"
	"TaskTracker/internal/repository/sjson"
	"TaskTracker/internal/service"
)

func main() {
	var repo repository.Repository
	jsonRepo, err := sjson.NewJSONStorage("data/data.json")
	if err != nil {
		fmt.Printf("⚠️ JSON storage error: %v\n", err)
		fmt.Println("🔄 Использую in-memory хранилище")
		repo = memory.NewStorage()
	} else {
		repo = jsonRepo
	}

	taskService := service.NewNewTaskService(repo)

	cliHandler := handler.NewCLIHandler(taskService)

	fmt.Println("🚀 TaskTracker запущен!")
	cliHandler.Run()

}
