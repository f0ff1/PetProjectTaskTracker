package main

import (
	"fmt"

	"TaskTracker/config"
	"TaskTracker/internal/handler"
	"TaskTracker/internal/repository"
	"TaskTracker/internal/repository/postgres"
	"TaskTracker/internal/service"

)

func main() {
	// var repo repository.Repository
	// jsonRepo, err := sjson.NewJSONStorage("data/data.json")
	// if err != nil {
	// 	fmt.Printf("⚠️ JSON storage error: %v\n", err)
	// 	fmt.Println("🔄 Использую in-memory хранилище")
	// 	repo = memory.NewStorage()
	// } else {
	// 	repo = jsonRepo
	// }

	var repo repository.Repository
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	dsn := cfg.GetDSN()
	fmt.Println(dsn)
	postgreSQLRepo, err := postgres.NewPostgresStorage(dsn)
	if err != nil {
		panic(err)
	}
	repo = postgreSQLRepo

	taskService := service.NewNewTaskService(repo)

	cliHandler := handler.NewCLIHandler(taskService)

	fmt.Println("🚀 TaskTracker запущен!")
	cliHandler.Run()

}
