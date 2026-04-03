package main

import (
	"context"
	"fmt"
	"os"

	"TaskTracker/internal/app"
)

func main() {
	ctx := context.Background()

	app, err := app.NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации приложения: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка выполнения приложения: %v\n", err)
		os.Exit(1)
	}
}
