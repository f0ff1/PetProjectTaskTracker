package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"TaskTracker/internal/model"
)

func (h *TelegramHandler) handleAdminTasks(ctx context.Context, chatID int64) {
	if h.extService == nil {
		h.sendMessage(chatID, "❌ Админ-панель доступна только при использовании PostgreSQL", false)
		return
	}

	tasks, err := h.extService.GetAllTasksForAdmin(ctx)
	if err != nil {
		log.Printf("❌ ERROR [GetAllTasksForAdmin]: %v\n", err)
		h.sendMessage(chatID, "⚠️ Не удалось получить все задачи. Попробуйте позже", false)
		return
	}

	if len(tasks) == 0 {
		h.sendMessage(chatID, "📭 Нет задач ни у одного пользователя.", false)
		return
	}

	h.sendMessage(chatID, fmt.Sprintf("👑 *Всего задач в системе:* %d\n\n", len(tasks)), true)

	// Группируем по пользователям
	userTasks := make(map[int][]*model.Task)
	for _, task := range tasks {
		userTasks[task.UserID] = append(userTasks[task.UserID], task)
	}

	for userID, userTaskList := range userTasks {
		user, _ := h.extService.GetUserByID(ctx, userID)
		username := "unknown"
		if user != nil {
			username = user.Username
			if username == "" {
				username = user.FirstName
			}
		}

		h.sendMessage(chatID, fmt.Sprintf("👤 *Пользователь: %s* (ID: %d) — %d задач", username, userID, len(userTaskList)), true)

		for _, task := range userTaskList {
			status := "❌"
			if task.Completed {
				status = "✅"
			}
			h.sendMessage(chatID, fmt.Sprintf("  %s #%d: %s", status, task.ID, task.Title), false)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func (h *TelegramHandler) handleAdminUsers(ctx context.Context, chatID int64) {
	if h.extService == nil {
		h.sendMessage(chatID, "❌ Админ-панель доступна только при использовании PostgreSQL", false)
		return
	}

	users, err := h.extService.GetAllUsers(ctx)
	if err != nil {
		log.Printf("❌ ERROR [GetAllUsers]: %v\n", err)
		h.sendMessage(chatID, "⚠️ Не удалось получить список пользователей. Попробуйте позже", false)
		return
	}

	if len(users) == 0 {
		h.sendMessage(chatID, "📭 Нет пользователей.", false)
		return
	}

	var sb strings.Builder
	sb.WriteString("👑 *Список пользователей*\n\n")

	for _, u := range users {
		adminBadge := ""
		if u.IsAdmin {
			adminBadge = " 👑"
		}
		name := u.Username
		if name == "" {
			name = u.FirstName
		}
		if name == "" {
			name = strconv.FormatInt(u.TelegramID, 10)
		}
		sb.WriteString(fmt.Sprintf("• *%s* (ID: %d)%s\n", name, u.ID, adminBadge))
		sb.WriteString(fmt.Sprintf("  Telegram: `%d`\n", u.TelegramID))
		sb.WriteString(fmt.Sprintf("  Активен: %s\n", u.LastActive.Format("02.01.2006 15:04")))
		sb.WriteString("\n")
	}

	h.sendMarkdown(chatID, sb.String())
}
