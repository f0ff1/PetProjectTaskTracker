package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Строка подключения к БД (уже готова)
	DatabaseURL string

	// Telegram
	TelegramBotToken string
}

func LoadConfig() (*Config, error) {
	// Загружаем подходящий .env файл
	if err := loadEnvFile(); err != nil {
		fmt.Printf("⚠️ %v\n", err)
	}

	// Получаем переменные
	dbURL := os.Getenv("DATABASE_URL")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	// Если нет DATABASE_URL, собираем из локальных переменных
	if dbURL == "" {
		dbURL = buildDSNFromLocal()
	}

	// Проверяем обязательные параметры
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN не задан")
	}

	if dbURL == "" {
		return nil, fmt.Errorf("не удалось собрать DATABASE_URL")
	}

	// Логируем режим работы
	if os.Getenv("DB_HOST") != "" || os.Getenv("DATABASE_URL") != "" {
		fmt.Println("☁️ Режим: Облако (Railway)")
	} else {
		fmt.Println("💻 Режим: Локальная разработка")
	}

	fmt.Printf("🔗 БД: %s\n", maskString(dbURL))

	return &Config{
		DatabaseURL:      dbURL,
		TelegramBotToken: botToken,
	}, nil
}

func (c *Config) GetDSN() string {
	if !strings.Contains(c.DatabaseURL, "timezone") {
		return c.DatabaseURL + "&timezone=Europe/Moscow"
	}
	return c.DatabaseURL
}

// loadEnvFile ищет и загружает .env файл
func loadEnvFile() error {
	// Возможные имена файлов (в порядке приоритета)
	possibleNames := []string{".env", "dbconfig.env"}

	// Возможные пути
	paths := []string{"."}

	// Если запуск из cmd/, добавляем родительскую папку
	if wd, err := os.Getwd(); err == nil {
		if filepath.Base(wd) == "cmd" {
			paths = append(paths, "..")
		}
	}

	// Ищем файл
	for _, name := range possibleNames {
		for _, dir := range paths {
			fullPath := filepath.Join(dir, name)
			if _, err := os.Stat(fullPath); err == nil {
				if err := godotenv.Load(fullPath); err == nil {
					fmt.Printf("✅ Загружен конфиг: %s\n", fullPath)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("файл .env не найден, используем переменные окружения")
}

// buildDSNFromLocal собирает DSN из локальных переменных
func buildDSNFromLocal() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnvAsInt("DB_PORT", 5432)
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "tasktracker")
	sslmode := getEnv("DB_SSLMODE", "disable")

	if password == "" {
		return ""
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)
}

// getEnv получает переменную окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// maskString скрывает часть строки для логов
func maskString(s string) string {
	if s == "" {
		return "<empty>"
	}

	// Скрываем пароль в строке postgres://user:pass@host/db
	if strings.HasPrefix(s, "postgres://") {
		// Находим позицию после пароля
		atIndex := strings.Index(s, "@")
		colonIndex := strings.Index(s[8:], ":")

		if colonIndex > 0 && atIndex > 0 {
			user := s[8 : 8+colonIndex]
			rest := s[atIndex:]
			return fmt.Sprintf("postgres://%s:***%s", user, rest)
		}
	}

	if len(s) < 20 {
		return "***"
	}
	return s[:10] + "..." + s[len(s)-10:]
}
