package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"TaskTracker/config"
	"TaskTracker/internal/handler"
	"TaskTracker/internal/repository"
	"TaskTracker/internal/repository/postgres"
	"TaskTracker/internal/service"

)

func main() {
	app := &cli.App{
		Name:    "Task Tracker",
		Usage:   "Пет-проект по добавлению задач с разной реализацией хранилища",
		Version: "4.0",
		Authors: []*cli.Author{
			{
				Name: "Андрей Кантур",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "storage",
				Aliases: []string{"st"},
				Value:   "postgres",
				Usage:   "Вид хранилища (memory, json, postgres)",
			},
			&cli.StringFlag{
				Name:    "json-path",
				Aliases: []string{"j"},
				Value:   "C:/GoLand/GoCourse/TaskTracker/data/data.json",
				Usage:   "Путь к JSON файлу (для json хранилища)",
			},
		},
		Action: func(c *cli.Context) error {
			requestedStorageType := c.String("storage")
			// jsonPath := c.String("json-path")
			var repo repository.Repository
			var usedStorage string

			switch requestedStorageType {
			case "postgres":
				cfg, configErr := config.LoadConfig()
				if configErr != nil {
					fmt.Printf("❌ Ошибка загрузки конфига PostgreSQL: %v\n", configErr)
				}
				dsn := cfg.GetDSN()
				fmt.Printf("🔌 Подключение к PostgreSQL: %s\n", dsn)
				pgRepo, pgErr := postgres.NewPostgresStorage(dsn)
				if pgErr != nil {
					fmt.Printf("❌ Ошибка подключения к PostgreSQL: %v\n", pgErr)
					goto tryJSON
				}
				repo = pgRepo
				usedStorage = "PostgreSQL"
				fmt.Println("✅ PostgreSQL хранилище инициализировано")
				break

			tryJSON:
				fallthrough

			case "json":
				// fmt.Printf("📁 Попытка использовать JSON файл: %s\n", jsonPath)
				// jsonRepo, jsonErr := sjson.NewJSONStorage(jsonPath)
				// if jsonErr != nil {
				// 	fmt.Printf("❌ Ошибка JSON хранилища: %v\n", jsonErr)
				// 	goto tryMemory
				// }
				// repo = jsonRepo
				// usedStorage = "JSON"
				fmt.Println("✅ JSON хранилище инициализировано")
				break

			// tryMemory:
			// 	fallthrough
			case "memory":
				// repo = memory.NewStorage()
				// usedStorage = "in-memory"
				fmt.Println("✅ In-memory хранилище инициализировано")
			default:
				return fmt.Errorf("неизвестный тип хранилища: %s (доступные: postgres, json, memory)", requestedStorageType)
			}
			if requestedStorageType != usedStorage {
				fmt.Printf("⚠️ ВНИМАНИЕ: Запрошено хранилище '%s', но используется '%s'\n", requestedStorageType, usedStorage)
			}
			taskService := service.NewTaskService(repo)
			cliHandler := handler.NewCLIHandler(taskService)

			fmt.Printf("🚀 TaskTracker запущен с хранилищем: %s\n", usedStorage)
			cliHandler.Run()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
