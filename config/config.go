package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Для облака (приоритет)
	DatabaseURL string

	// Для локальной разработки
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Telegram
	TelegramBotToken string
}

func LoadConfig() (*Config, error) {
	// Пытаемся загрузить .env файл (только для локальной разработки)
	loadEnvFile()

	cfg := &Config{
		// Локальные переменные (дефолты)
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "tasktracker"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Облачные переменные
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
	}

	// Проверка: токен обязателен везде
	if cfg.TelegramBotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN не задан")
	}

	// Если есть DATABASE_URL - используем облачный режим (Railway)
	if cfg.DatabaseURL != "" {
		fmt.Println("☁️ Режим: Облако (Railway)")
		fmt.Printf("🔗 DATABASE_URL: %s\n", maskString(cfg.DatabaseURL))
		return cfg, nil
	}

	// Иначе используем локальный режим
	fmt.Println("💻 Режим: Локальная разработка")

	// Проверяем локальные параметры
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD не задан (для локальной разработки)")
	}

	fmt.Printf("🔗 Подключение: %s:%d/%s\n", cfg.DBHost, cfg.DBPort, cfg.DBName)

	return cfg, nil
}

// GetDSN возвращает строку подключения к БД
func (c *Config) GetDSN() string {
	// Приоритет: DATABASE_URL (облако)
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	// Локальный режим
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// loadEnvFile ищет и загружает .env файл (только для локальной разработки)
func loadEnvFile() {
	// Возможные имена файлов
	possibleNames := []string{".env", "dbconfig.env"}

	// Возможные пути (от текущей директории)
	paths := []string{"."}

	// Если запуск из cmd/, добавляем родительскую директорию
	if wd, err := os.Getwd(); err == nil {
		if filepath.Base(wd) == "cmd" {
			paths = append(paths, "..")
		}
	}

	for _, name := range possibleNames {
		for _, dir := range paths {
			fullPath := filepath.Join(dir, name)
			if _, err := os.Stat(fullPath); err == nil {
				if err := godotenv.Load(fullPath); err == nil {
					fmt.Printf("✅ Загружен конфиг: %s\n", fullPath)
					return
				}
			}
		}
	}

	// В Railway просто игнорируем
	fmt.Println("ℹ️ Файл .env не найден, используем переменные окружения")
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
	if len(s) < 20 {
		return "***"
	}
	return s[:10] + "..." + s[len(s)-10:]
}
