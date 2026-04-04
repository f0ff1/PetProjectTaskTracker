package app

import (
	"context"
	"fmt"
	"log/slog"

	"TaskTracker/config"
	myerrors "TaskTracker/errors"
	"TaskTracker/factory"
	"TaskTracker/internal/handler"
)

type Application struct {
	handler *handler.CLIHandler
}

func NewApp() (*Application, error) {
	handlers := []struct {
		name    string
		tryFunc func() (*handler.CLIHandler, error)
	}{
		{"PostgreSQL", tryPostgres},
		{"JSON", tryJSON},
		{"In-memory", tryInMemory},
	}

	for _, handler := range handlers {
		cliHandler, err := handler.tryFunc()
		if err != nil {
			slog.Error("Провалена ошибка создания handler", "handler", handler.name, "err", err)
			continue
		}

		if cliHandler != nil {
			return &Application{handler: cliHandler}, nil
		}
	}

	return nil, myerrors.ErrWrongTypeRepo
}

func (a *Application) Run(ctx context.Context) error {
	if a == nil || a.handler == nil {
		slog.Error("Ошибка: nil handler, неверный тип репозитория", "err", myerrors.ErrWrongTypeRepo)
		return myerrors.ErrWrongTypeRepo
	}
	a.handler.Run(ctx)
	return nil
}

func tryInMemory() (*handler.CLIHandler, error) {
	cliHandler, err := factory.CreateCLIHandler(factory.InMemory, "", "")
	return cliHandler, err
}

func tryJSON() (*handler.CLIHandler, error) {
	jsonPath := "C:/GoLand/GoCourse/TaskTracker/data/data.json"
	cliHandler, err := factory.CreateCLIHandler(factory.JSON, "", jsonPath)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания JSON application %w", err)
	}

	return cliHandler, nil
}

func tryPostgres() (*handler.CLIHandler, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Не удалось загрузить конфиг PostgreSQL: %w", err)
	}

	dsn := cfg.GetDSN()

	cliHandler, err := factory.CreateCLIHandler(factory.Postgres, dsn, "")

	if err != nil {

		return nil, fmt.Errorf("Ошибка подключения к PostgreSQL: %w", err)
	}

	return cliHandler, nil
}
