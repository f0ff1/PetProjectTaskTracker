package factory

import (
	"fmt"

	myerrors "TaskTracker/errors"
	"TaskTracker/internal/handler"
	"TaskTracker/internal/repository/memory"
	"TaskTracker/internal/repository/postgres"
	"TaskTracker/internal/repository/sjson"
	"TaskTracker/internal/service"
)

type StorageType string

const (
	InMemory StorageType = "memory"
	JSON     StorageType = "json"
	Postgres StorageType = "postgres"
)

func CreateCLIHandler(storageType StorageType, connString string, jsonPath string) (*handler.CLIHandler, error) {

	switch storageType {
	case InMemory:
		repo := memory.NewInMemoryRepo()
		svc := service.NewTaskService(repo)
		return handler.NewCLIHandler(svc, nil), nil
	case JSON:
		repo, err := sjson.NewJSONRepo(jsonPath)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при создании JSON репозитория: %w", err)
		}
		svc := service.NewTaskService(repo)
		return handler.NewCLIHandler(svc, nil), nil
	case Postgres:
		repo, err := postgres.NewPostgresRepo(connString)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при создании Postgres репозитория: %w", err)
		}
		svc := service.NewPostgresTaskService(repo)
		return handler.NewCLIHandler(svc.TaskService, svc), nil
	default:
		return nil, fmt.Errorf("Ошибка при создании сервиса: %w", myerrors.ErrWrongTypeRepo)
	}
}
