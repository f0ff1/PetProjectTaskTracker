package main

import (
	"fmt"

	"TaskTracker/config"
	myerrors "TaskTracker/errors"
	"TaskTracker/factory"
	"TaskTracker/internal/handler"
	"TaskTracker/internal/service"

)

func main() {
	readyHandler := getReadyHandler()
	if readyHandler == nil {
		panic(myerrors.ErrWrongTypeRepo)
	} else {
		readyHandler.Run()
	}

}

func getReadyHandler() interface{ Run() } {
	if handler := tryPostgres(); handler != nil {
		return handler
	}

	if handler := tryJSON(); handler != nil {
		return handler
	}

	return tryInMemory()
}

func tryInMemory() interface{ Run() } {
	fmt.Println("💾 Использую In-Memory хранилище")
	svc, err := factory.CreateTaskService(factory.InMemory, "", "")
	if err != nil {
		fmt.Printf("❌ Ошибка создания In-Memory: %v\n", err)
		return nil
	}
	if memSvc, ok := svc.(*service.TaskService); ok {
		return handler.NewCLIHandler(memSvc)
	}
	return nil
}

func tryJSON() interface{ Run() } {
	jsonPath := "C:/GoLand/GoCourse/TaskTracker/data/data.json"
	fmt.Printf("📁 Пробую использовать JSON файл: %s\n", jsonPath)
	svc, err := factory.CreateTaskService(factory.JSON, "", jsonPath)
	if err != nil {
		fmt.Printf("⚠️ Ошибка загрузки JSON: %v\n", err)
		return nil
	}
	if jsonSvc, ok := svc.(*service.TaskService); ok {
		fmt.Println("✅ Использую JSON хранилище")
		return handler.NewCLIHandler(jsonSvc)
	}

	return nil
}

func tryPostgres() interface{ Run() } {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("⚠️ Не удалось загрузить конфиг PostgreSQL: %v\n", err)
		return nil
	}
	dsn := cfg.GetDSN()
	fmt.Printf("🔌 Подключение к PostgreSQL: %s\n", dsn)
	svc, err := factory.CreateTaskService(factory.Postgres, dsn, "")
	if err != nil {
		fmt.Printf("⚠️ Ошибка подключения к PostgreSQL: %v\n", err)
		return nil
	}

	if pgSvc, ok := svc.(*service.PostgresTaskService); ok {
		fmt.Println("✅ Использую PostgreSQL хранилище")
		return handler.NewPostgresCLIHandler(pgSvc)
	}
	return nil
}
