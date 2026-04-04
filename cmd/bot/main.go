package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"TaskTracker/config"
	"TaskTracker/internal/handler"
	"TaskTracker/internal/repository/database"
	"TaskTracker/internal/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Ошибка загрузки конфига: %v", err)
	}

	repo, err := database.NewPostgresRepo(cfg.GetDSN())
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer repo.Close()

	//repo.StartStatsUpdater(ctx, 5*time.Minute)

	taskSvc := service.NewTaskService(repo)
	extendedSvc := service.NewPostgresTaskService(repo)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("❌ Ошибка создания бота: %v", err)
	}

	log.Printf("✅ Бот запущен: @%s", bot.Self.UserName)

	tgHandler := handler.NewTelegramHandler(bot, taskSvc, extendedSvc)

	errCh := make(chan error, 1)
	go func() {
		if err := tgHandler.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("\n🛑 Получен сигнал завершения, останавливаюсь...")
		time.Sleep(2 * time.Second)
		log.Println("✅ Приложение завершено")
	case err := <-errCh:
		log.Fatalf("❌ Ошибка работы бота: %v", err)
	}
}
