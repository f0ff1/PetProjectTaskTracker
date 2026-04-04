package factory

import (
	"context"
	"fmt"
	"time"

	myerrors "TaskTracker/errors"
	"TaskTracker/internal/handler"
	"TaskTracker/internal/repository/database"

	// "TaskTracker/internal/repository/memory"
	// "TaskTracker/internal/repository/sjson"
	"TaskTracker/internal/service"
)

type StorageType string

const (
	InMemory      StorageType = "memory"
	JSON          StorageType = "json"
	Postgres      StorageType = "postgres"
	defaultUserID             = 1
)

func CreateCLIHandler(storageType StorageType, connString string, jsonPath string) (*handler.CLIHandler, error) {

	switch storageType {
	// case InMemory:
	// 	repo := memory.NewInMemoryRepo()
	// 	svc := service.NewTaskService(repo)
	// 	return handler.NewCLIHandler(svc, nil, defaultUserID), nil
	// case JSON:
	// 	repo, err := sjson.NewJSONRepo(jsonPath)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("Ошибка при создании JSON репозитория: %w", err)
	// 	}
	// 	svc := service.NewTaskService(repo)
	// 	return handler.NewCLIHandler(svc, nil), nil
	case Postgres:
		repo, err := database.NewPostgresRepo(connString)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при создании Postgres репозитория: %w", err)
		}
		svc := service.NewPostgresTaskService(repo)

		repo.StartStatsUpdater(context.Background(), defaultUserID, 5*time.Minute)
		return handler.NewCLIHandler(svc.TaskService, svc, defaultUserID), nil
	default:
		return nil, fmt.Errorf("Ошибка при создании сервиса: %w", myerrors.ErrWrongTypeRepo)
	}
}
