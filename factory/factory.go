package factory

import (
	myerrors "TaskTracker/errors"
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

func CreateTaskService(storageType StorageType, connString string, jsonPath string) (interface{}, error) {

	switch storageType {
	case InMemory:
		repo := memory.NewInMemoryRepo()
		return service.NewTaskService(repo), nil
	case JSON:
		repo, err := sjson.NewJSONRepo(jsonPath)
		if err != nil {
			return nil, err
		}
		return service.NewTaskService(repo), nil
	case Postgres:
		repo, err := postgres.NewPostgresRepo(connString)
		if err != nil {
			return nil, err
		}
		return service.NewPostgresTaskService(repo), nil
	default:
		return nil, myerrors.ErrWrongTypeRepo
	}
}
