package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	DatabaseURL string

	TelegramBotToken string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load("../dbconfig.env"); err != nil {
		slog.Warn("Предупреждение: файл .env не найден, используем переменные окружения", "err", err)
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "myapp_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		DatabaseURL: getEnv("DATABASE_URL", ""),

		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
	}

	if cfg.TelegramBotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN не задан")
	}

	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD не задан")
	}
	return cfg, nil
}

func getEnv(key, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defValue
}

// PGX
func (c *Config) GetDSN() string {

	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	fmt.Printf("USER=%s PASS=%s\n", c.DBUser, c.DBPassword)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode)
}

// PQ
func (c *Config) GetConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}
