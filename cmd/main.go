package main

import (
	"context"
	"fmt"

	"TaskTracker/config"
	myerrors "TaskTracker/errors"
	"TaskTracker/factory"

)

func main() {
	ctx := context.Background()
	readyHandler := getReadyHandler()
	if readyHandler == nil {
		panic(myerrors.ErrWrongTypeRepo)
	} else {
		readyHandler.Run(ctx)
	}
}

func getReadyHandler() interface{ Run(ctx context.Context) } {
	if handler := tryPostgres(); handler != nil {
		return handler
	}

	if handler := tryJSON(); handler != nil {
		return handler
	}

	return tryInMemory()
}

func tryInMemory() interface{ Run(ctx context.Context) } {
	fmt.Println("💾 Использую In-Memory хранилище")
	cliHandler, err := factory.CreateCLIHandler(factory.InMemory, "", "")
	if err != nil {
		fmt.Printf("❌ Ошибка создания In-Memory: %v\n", err)
		return nil
	}
	return cliHandler
}

func tryJSON() interface{ Run(ctx context.Context) } {
	jsonPath := "C:/GoLand/GoCourse/TaskTracker/data/data.json"
	fmt.Printf("📁 Пробую использовать JSON файл: %s\n", jsonPath)
	cliHandler, err := factory.CreateCLIHandler(factory.JSON, "", jsonPath)
	if err != nil {
		fmt.Printf("⚠️ Ошибка загрузки JSON: %v\n", err)
		return nil
	}
	fmt.Println("✅ Использую JSON хранилище")

	return cliHandler
}

func tryPostgres() interface{ Run(ctx context.Context) } {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("⚠️ Не удалось загрузить конфиг PostgreSQL: %v\n", err)
		return nil
	}
	
	dsn := cfg.GetDSN()
	fmt.Printf("🔌 Подключение к PostgreSQL: %s\n", dsn)
	cliHandler, err := factory.CreateCLIHandler(factory.Postgres, dsn, "")
	if err != nil {
		fmt.Printf("⚠️ Ошибка подключения к PostgreSQL: %v\n", err)
		return nil
	}

	fmt.Println("✅ Использую PostgreSQL хранилище")
	return cliHandler
}
