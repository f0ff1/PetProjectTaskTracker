package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"

	myErrors "TaskTracker/errors"
)

func runMigrations(connStr string) error {

	migrationsPath := getMigrationsPath()

	log.Printf("Путь к миграциям: %s", migrationsPath)

	migration, err := migrate.New(
		"file://"+migrationsPath,
		connStr,
	)

	if err != nil {
		return fmt.Errorf("%v : %w", myErrors.ErrCreateMigration, err)
	}
	defer migration.Close()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%v : %w", myErrors.ErrCantUseMigration, err)
	}

	log.Println("Миграции применены")
	return nil
}

func getMigrationsPath() string {

	possiblePaths := []string{
		"migrations",       // /app/migrations (Docker/Railway)
		"./migrations",     // ./migrations
		"../migrations",    // ../migrations (локально при запуске из cmd/)
		"../../migrations", // на два уровня выше
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("✅ Найдены миграции: %s", path)
			return path
		}
	}

	// Если не нашли, возвращаем стандартный путь
	log.Println("⚠️ Миграции не найдены, использую './migrations'")
	return "migrations"
}
