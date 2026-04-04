package service

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"TaskTracker/internal/model"
	"TaskTracker/internal/utils"
)

type ReminderService struct {
	taskService *TaskService
	extService  ExtendedTaskService
	bot         *tgbotapi.BotAPI
}

func NewRemindereService(taskService *TaskService, extService ExtendedTaskService, bot *tgbotapi.BotAPI) *ReminderService {
	return &ReminderService{taskService: taskService, extService: extService, bot: bot}
}

func (s *ReminderService) StartReminderChecker(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Println("Сервис напоминаний запущен")

		for {
			select {
			case <-ticker.C:
				s.checkAndSendRimenders(ctx)
			case <-ctx.Done():
				log.Println("Сервис напоминаний остановлен")
				return

			}
		}
	}()
}

func (s *ReminderService) checkAndSendRimenders(ctx context.Context) {
	log.Printf("[REMINDER CHECK] Starting check at %s", time.Now().Format("02.01.2006 15:04:05"))

	tasks, err := s.extService.GetTasksForReminder(ctx)
	if err != nil {
		log.Printf("[REMINDER CHECK] Error getting tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		log.Printf("[REMINDER CHECK] No tasks to remind")
		return
	}

	log.Printf("[REMINDER CHECK] Found %d tasks for reminder", len(tasks))

	for _, task := range tasks {
		s.sendReminder(ctx, task)
	}
}

func (s *ReminderService) sendReminder(ctx context.Context, task *model.Task) {
	log.Printf("[REMINDER SEND] Sending reminder for task %d (ID: %d), Title: %s", task.UserTaskID, task.ID, task.Title)

	user, err := s.extService.GetUserByID(ctx, task.UserID)

	if err != nil {
		log.Printf("[REMINDER SEND] User not found for task #%d: %v", task.ID, err)
		return
	}

	if task.DueDate == nil {
		log.Printf("[REMINDER SEND] No DueDate for task %d", task.ID)
		return
	}

	timeLeft := time.Until(*task.DueDate)

	var timeLeftText string

	switch {
	case timeLeft <= 0:
		timeLeftText = "СРОК ИСТЁК !"
	case timeLeft < 1*time.Hour:
		timeLeftText = fmt.Sprintf("⏰ Осталось %d минут!", int(timeLeft.Minutes()))
	case timeLeft <= 24*time.Hour:
		timeLeftText = fmt.Sprintf("⏰ Осталось %d часов!", int(timeLeft.Hours()))
	default:
		timeLeftText = fmt.Sprintf("📅 Осталось %d дней", int(timeLeft.Hours()/24))
	}

	// Формируем текст напоминания
	reminderText := fmt.Sprintf(
		"🔔 *Напоминание о задаче*\n\n"+
			"📌 *%s*\n"+
			"%s\n\n"+
			"⏰ Срок: %s\n"+
			"%s",
		utils.EscapeMarkdown(task.Title),
		utils.EscapeMarkdown(task.Description),
		task.DueDate.Format("02.01.2006 15:04"),
		timeLeftText,
	)

	// Отправляем сообщение в Telegram
	msg := tgbotapi.NewMessage(user.TelegramID, reminderText)
	msg.ParseMode = "Markdown"

	_, err = s.bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки напоминания пользователю %d: %v", user.ID, err)
		return
	}

	// Отмечаем, что напоминание отправлено
	err = s.extService.MarkReminderSent(ctx, task.ID)
	if err != nil {
		log.Printf("Ошибка при обновлении статуса напоминания для задачи #%d: %v", task.ID, err)
		return
	}

	log.Printf("Напоминание отправлено пользователю %s (ID: %d) для задачи #%d", user.Username, user.ID, task.ID)
}
