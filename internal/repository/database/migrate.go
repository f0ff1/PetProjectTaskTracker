package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	myErrors "TaskTracker/errors"
)

func runMigrations(connStr string) error {
	log.Println("🚀 ЗАПУСК МИГРАЦИЙ")

	migrationsPath := getMigrationsPath()
	log.Printf("📁 Путь к миграциям: %s", migrationsPath)

	// Проверяем, что файлы реально читаются
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Printf("❌ Не могу прочитать папку %s: %v", migrationsPath, err)
	} else {
		log.Printf("📄 Найдено файлов: %d", len(files))
		for _, f := range files {
			log.Printf("   - %s", f.Name())
		}
	}

	log.Printf("🔗 Подключение к БД: %s", maskDSN(connStr))

	migration, err := migrate.New(
		"file://"+migrationsPath,
		connStr,
	)

	if err != nil {
		log.Printf("❌ Ошибка создания мигратора: %v", err)
		return fmt.Errorf("%v : %w", myErrors.ErrCreateMigration, err)
	}
	defer migration.Close()

	// Получаем текущую версию
	version, dirty, err := migration.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Printf("⚠️ Ошибка получения версии: %v", err)
	} else {
		log.Printf("📌 Текущая версия: %d, dirty: %v", version, dirty)
	}

	log.Println("🔄 Применяем миграции...")

	if err := migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("✅ Нет новых миграций (уже применены)")
			return nil
		}
		log.Printf("❌ Ошибка применения миграций: %v", err)
		return fmt.Errorf("%v : %w", myErrors.ErrCantUseMigration, err)
	}

	log.Println("✅ Миграции успешно применены")
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

	log.Println("⚠️ Миграции не найдены, использую './migrations'")
	return "migrations"
}

func maskDSN(dsn string) string {
	if len(dsn) < 30 {
		return "***"
	}
	return dsn[:15] + "..." + dsn[len(dsn)-15:]
}
