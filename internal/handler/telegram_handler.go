package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"TaskTracker/internal/model"
	"TaskTracker/internal/service"
	"TaskTracker/internal/utils"
)

type TelegramHandler struct {
	bot         *tgbotapi.BotAPI
	taskService *service.TaskService
	extService  service.ExtendedTaskService

	userStates   map[int64]string
	pendingTasks map[int64]*PendingTask
	users        map[int64]*model.User
	mu           sync.RWMutex
}

func NewTelegramHandler(bot *tgbotapi.BotAPI, taskService *service.TaskService, extService service.ExtendedTaskService) *TelegramHandler {
	return &TelegramHandler{
		bot:          bot,
		taskService:  taskService,
		extService:   extService,
		userStates:   make(map[int64]string),
		pendingTasks: make(map[int64]*PendingTask),
		users:        make(map[int64]*model.User),
	}
}

func (h *TelegramHandler) getUser(ctx context.Context, chatID int64, msg *tgbotapi.Message) (*model.User, error) {

	h.mu.RLock()
	user, exists := h.users[chatID]
	h.mu.RUnlock()

	if exists && user != nil {
		return user, nil
	}

	user, err := h.extService.GetOrCreateUser(ctx, msg.From.ID, msg.From.UserName, msg.From.FirstName, msg.From.LastName)
	if err != nil {
		return nil, err
	}

	h.mu.Lock()
	h.users[chatID] = user
	h.mu.Unlock()

	log.Printf("👤 Пользователь: ID=%d, Telegram=%d, Admin=%v", user.ID, user.TelegramID, user.IsAdmin)
	return user, nil
}

func (h *TelegramHandler) setState(chatID int64, state string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.userStates[chatID] = state
}

func (h *TelegramHandler) getState(chatID int64) string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.userStates[chatID]
}

func (h *TelegramHandler) clearState(chatID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.userStates, chatID)
	delete(h.pendingTasks, chatID)
}

func (h *TelegramHandler) setPendingTask(chatID int64, task *PendingTask) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.pendingTasks[chatID] = task
}

func (h *TelegramHandler) getPendingTask(chatID int64) *PendingTask {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.pendingTasks[chatID]
}

func (h *TelegramHandler) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	log.Println("✅ Бот запущен, ожидание сообщений...")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Контекст отменен, завершаю работу бота...")
			return nil

		case update, ok := <-updates:
			if !ok {
				log.Println("❌ Канал обновлений закрыт")
				return nil
			}

			if update.Message != nil {
				go h.handleMessage(ctx, update.Message)
			}

			if update.CallbackQuery != nil {
				go h.handleCallback(ctx, update.CallbackQuery)
			}
		}
	}
}

func (h *TelegramHandler) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ID

	// Получаем или создаем пользователя
	user, err := h.getUser(ctx, chatID, msg)
	if err != nil {
		h.handleAndLogError(chatID, "getUser", err)
		return
	}

	currentState := h.getState(chatID)

	if currentState != StateIdle {
		h.handleDialogState(ctx, chatID, user.ID, text, currentState)
		return
	}

	// Обработчики простых команд
	commandHandlers := map[string]func(){
		"/start": func() { h.sendHelpMessage(chatID, user.IsAdmin) },
		"/help":  func() { h.sendHelpMessage(chatID, user.IsAdmin) },
		"/add":   func() { h.startAddTaskDialog(chatID) },
		"/list":  func() { h.handleListCommand(ctx, chatID, user.ID) },
		"/stats": func() { h.handleStatsCommand(ctx, chatID, user.ID) },
	}

	if handler, ok := commandHandlers[text]; ok {
		handler()
		return
	}

	// Обработка команд с параметрами
	switch {
	case strings.HasPrefix(text, "/add "):
		h.handleFastAdd(ctx, chatID, user.ID, text)

	case strings.HasPrefix(text, "/complete"):
		h.handleCompleteCommand(ctx, chatID, user.ID, text)

	case strings.HasPrefix(text, "/delete"):
		h.handleDeleteCommand(ctx, chatID, user.ID, text)

	case strings.HasPrefix(text, "/find"):
		h.handleFindCommand(ctx, chatID, user.ID, text)

	case strings.HasPrefix(text, "/tag"):
		h.handleTagCommand(ctx, chatID, user.ID, text)

	// Админские команды
	case user.IsAdmin && strings.HasPrefix(text, "/admintask"):
		h.handleAdminTasks(ctx, chatID)
	case user.IsAdmin && strings.HasPrefix(text, "/adminuser"):
		h.handleAdminUsers(ctx, chatID)

	default:
		h.sendMessage(chatID, "❌ Неизвестная команда. Введите /help", false)
	}
}

func (h *TelegramHandler) handleCallback(ctx context.Context, query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID
	data := query.Data

	h.bot.Send(tgbotapi.NewCallback(query.ID, ""))

	// Получаем пользователя из кэша
	h.mu.RLock()
	user := h.users[chatID]
	h.mu.RUnlock()

	if user == nil {
		h.sendMessage(chatID, "❌ Ошибка: пользователь не найден", false)
		return
	}

	switch {
	case strings.HasPrefix(data, "complete_"):
		idStr := strings.TrimPrefix(data, "complete_")
		id, _ := strconv.Atoi(idStr)

		task, err := h.taskService.CompleteTask(ctx, user.ID, id)
		if err != nil {
			h.handleAndLogErrorWithContext(chatID, "CompleteTask", err, user.ID, fmt.Sprintf("TaskID=%d", id))
			return
		}
		h.updateTaskMessage(chatID, messageID, task)

	case strings.HasPrefix(data, "delete_"):
		idStr := strings.TrimPrefix(data, "delete_")
		taskID, _ := strconv.Atoi(idStr)

		if err := h.extService.DeleteTask(ctx, user.ID, taskID); err != nil {
			h.handleAndLogErrorWithContext(chatID, "DeleteTask", err, user.ID, fmt.Sprintf("TaskID=%d", taskID))
			return
		}

		deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
		h.bot.Send(deleteMsg)
		h.sendMessage(chatID, fmt.Sprintf("🗑️ Задача #%d удалена", taskID), false)
	}
}

func (h *TelegramHandler) startAddTaskDialog(chatID int64) {
	h.setState(chatID, StateAwaitingTitle)
	h.sendMessage(chatID, "📝 Введите название задачи:", false)
}

func (h *TelegramHandler) handleDialogState(ctx context.Context, chatID int64, userID int, input string, state string) {
	switch state {
	case StateAwaitingTitle:
		if input == "" {
			h.sendMessage(chatID, "❌ Название не может быть пустым. Попробуйте еще раз:", false)
			return
		}

		h.setPendingTask(chatID, &PendingTask{
			Title:       input,
			Description: "",
			Tags:        []string{},
		})
		h.setState(chatID, StateAwaitingDescription)
		h.sendMessage(chatID, "📄 Введите описание задачи (или /skip чтобы пропустить):", false)

	case StateAwaitingDescription:
		pending := h.getPendingTask(chatID)
		if pending == nil {
			h.clearState(chatID)
			h.sendMessage(chatID, "❌ Ошибка: начните заново с /add", false)
			return
		}

		if input != "/skip" && input != "" {
			pending.Description = input
		}
		h.setPendingTask(chatID, pending)
		h.setState(chatID, StateAwaitingTags)
		h.sendMessage(chatID, "🏷️ Введите теги через пробел (или /skip чтобы пропустить):\n\nПример: работа важное дом", false)

	case StateAwaitingTags:
		pending := h.getPendingTask(chatID)
		if pending == nil {
			h.clearState(chatID)
			h.sendMessage(chatID, "❌ Ошибка: начните заново с /add", false)
			return
		}

		if input != "/skip" && input != "" {
			pending.Tags = h.parseTags(input)
		}
		h.setPendingTask(chatID, pending)
		h.setState(chatID, StateAwaitingDueDate)
		h.sendMessage(chatID, "📅 Введите срок выполнения (или /skip чтобы пропустить):\n\nФормат: ДД.МММ.ГГГГ ЧЧ:МИ\nПример: 15.04.2026 18:30", false)

	case StateAwaitingDueDate:
		pending := h.getPendingTask(chatID)
		if pending == nil {
			h.clearState(chatID)
			h.sendMessage(chatID, "❌ Ошибка: начните заново с /add", false)
			return
		}

		if input != "/skip" && input != "" {
			// Проверяем формат даты
			if _, err := time.Parse(dateFormat, input); err != nil {
				h.sendMessage(chatID, "❌ Неправильный формат даты. Используйте: ДД.МММ.ГГГГ ЧЧ:МИ\nПример: 15.04.2026 18:30", false)
				return
			}
			// Сохраняем дату как строку
			pending.DueDate = &input
		}
		h.setPendingTask(chatID, pending)
		h.setState(chatID, StateAwaitingReminderOffset)
		h.sendMessage(chatID, "⏰ Введите время напоминания перед сроком (или /skip):\n\nПримеры: 30m, 1h, 2h, 1d", false)

	case StateAwaitingReminderOffset:
		pending := h.getPendingTask(chatID)
		if pending == nil {
			h.clearState(chatID)
			h.sendMessage(chatID, "❌ Ошибка: начните заново с /add", false)
			return
		}

		if input != "/skip" && input != "" {
			// Проверяем формат offset (минимальная проверка)
			if !isValidDuration(input) {
				h.sendMessage(chatID, "❌ Неправильный формат. Используйте: 30m, 1h, 2h, 1d\n(m=минуты, h=часы, d=дни)", false)
				return
			}
			pending.ReminderOffset = &input
		}

		// Преобразуем ReminderOffset в строку (поле в модели тип string)
		reminderOffset := ""
		if pending.ReminderOffset != nil {
			reminderOffset = *pending.ReminderOffset
		}

		// Создаем задачу с данными
		task, err := h.taskService.AddTaskWithReminder(ctx, userID, pending.Title, pending.Description, pending.Tags, pending.DueDate, reminderOffset)
		h.clearState(chatID)

		if err != nil {
			h.handleAndLogErrorWithContext(chatID, "AddTaskWithReminder", err, userID, fmt.Sprintf("Title=%s Tags=%v DueDate=%v", pending.Title, pending.Tags, pending.DueDate))
			return
		}

		h.sendTaskCard(chatID, task)
		h.sendMessage(chatID, "✅ Задача успешно создана! Используйте /list для просмотра.", false)
	}
}

func (h *TelegramHandler) sendHelpMessage(chatID int64, isAdmin bool) {
	helpText := `📋 *TaskTracker Bot - Команды*

/add - Добавить задачу (пошагово)
/add <название> #тег - Быстрое добавление
/list - Список всех задач
/complete <ID> - Отметить задачу как выполненную
/delete <ID> - Удалить задачу
/find <ID> - Найти задачу по ID
/tag <тег> - Найти задачи по тегу
/stats - Показать статистику
/help - Показать эту справку`

	if isAdmin {
		helpText += `

*👑 Админ-команды:*
/admintasks - Все задачи всех пользователей
/adminusers - Список всех пользователей`
	}

	h.sendMarkdown(chatID, helpText)
}

func (h *TelegramHandler) handleFastAdd(ctx context.Context, chatID int64, userID int, text string) {
	rawTitle := strings.TrimPrefix(text, "/add ")
	rawTitle = strings.TrimSpace(rawTitle)

	title, tags := h.parseTitleAndTags(rawTitle)

	task, err := h.taskService.AddTask(ctx, userID, title, "", tags)
	if err != nil {
		h.handleAndLogErrorWithContext(chatID, "AddTask_FastAdd", err, userID, fmt.Sprintf("Title=%s Tags=%v", title, tags))
		return
	}

	h.sendTaskCard(chatID, task)
}

func (h *TelegramHandler) handleListCommand(ctx context.Context, chatID int64, userID int) {
	tasks, err := h.taskService.GetAllTasks(ctx, userID)
	if err != nil {
		h.handleAndLogErrorWithContext(chatID, "GetAllTasks", err, userID, "ListCommand")
		return
	}

	if len(tasks) == 0 {
		h.sendMessage(chatID, "📭 У вас нет задач.", false)
		return
	}

	completed := 0
	for _, task := range tasks {
		h.sendTaskCard(chatID, task)
		if task.Completed {
			completed++
		}
		time.Sleep(100 * time.Millisecond)
	}

	summary := fmt.Sprintf("📊 *Итого:* %d задач | ✅ Выполнено: %d | ⏳ В работе: %d",
		len(tasks), completed, len(tasks)-completed)
	h.sendMarkdown(chatID, summary)
}

func (h *TelegramHandler) handleCompleteCommand(ctx context.Context, chatID int64, userID int, text string) {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		h.sendMessage(chatID, "❌ Использование: /complete <ID>", false)
		return
	}

	taskID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendMessage(chatID, "❌ ID должен быть числом", false)
		return
	}

	task, err := h.taskService.CompleteTask(ctx, userID, taskID)
	if err != nil {
		h.handleAndLogErrorWithContext(chatID, "CompleteTask", err, userID, fmt.Sprintf("TaskID=%d", taskID))
		return
	}

	h.sendMessage(chatID, fmt.Sprintf("✅ Задача \"%s\" выполнена!", task.Title), false)
}

func (h *TelegramHandler) handleDeleteCommand(ctx context.Context, chatID int64, userID int, text string) {
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		h.sendMessage(chatID, "❌ Использование: /delete <ID>", false)
		return
	}

	taskID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendMessage(chatID, "❌ ID должен быть числом", false)
		return
	}

	if err := h.extService.DeleteTask(ctx, userID, taskID); err != nil {
		h.handleAndLogErrorWithContext(chatID, "DeleteTask", err, userID, fmt.Sprintf("TaskID=%d", taskID))
		return
	}

	h.sendMessage(chatID, fmt.Sprintf("🗑️ Задача #%d удалена", taskID), false)
}

func (h *TelegramHandler) handleFindCommand(ctx context.Context, chatID int64, userID int, text string) {
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		h.sendMessage(chatID, "❌ Использование: /find <ID>", false)
		return
	}

	taskID, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendMessage(chatID, "❌ ID должен быть числом", false)
		return
	}

	task, err := h.taskService.GetTaskById(ctx, userID, taskID)
	if err != nil {
		h.handleAndLogErrorWithContext(chatID, "GetTaskById", err, userID, fmt.Sprintf("TaskID=%d", taskID))
		return
	}

	h.sendTaskCard(chatID, task)
}

func (h *TelegramHandler) handleTagCommand(ctx context.Context, chatID int64, userID int, text string) {
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		h.sendMessage(chatID, "❌ Использование: /tag <тег>", false)
		return
	}

	tag := parts[1]
	tasks, err := h.taskService.GetTasksByTag(ctx, userID, tag)
	if err != nil {
		h.handleAndLogErrorWithContext(chatID, "GetTasksByTag", err, userID, fmt.Sprintf("Tag=%s", tag))
		return
	}

	if len(tasks) == 0 {
		h.sendMessage(chatID, fmt.Sprintf("📭 Нет задач с тегом \"%s\"", tag), false)
		return
	}

	for _, task := range tasks {
		h.sendTaskCard(chatID, task)
		time.Sleep(100 * time.Millisecond)
	}

	summary := fmt.Sprintf("📊 *Найдено задач с тегом \"%s\":* %d", tag, len(tasks))
	h.sendMarkdown(chatID, summary)
}

func (h *TelegramHandler) handleStatsCommand(ctx context.Context, chatID int64, userID int) {
	if h.extService == nil {
		h.sendMessage(chatID, "❌ Статистика доступна только при использовании PostgreSQL", false)
		return
	}

	h.sendMessage(chatID, "📊 Обновляю статистику...", false)

	stats, err := h.extService.GetStatsForce(ctx, userID)
	if err != nil {
		h.handleAndLogErrorWithContext(chatID, "GetStatsForce", err, userID, "StatsCommand")
		return
	}

	if stats == nil || stats.Total == 0 {
		h.sendMessage(chatID, "📭 Нет данных для статистики", false)
		return
	}

	statsText := fmt.Sprintf(
		"📊 *Статистика*\n\n"+
			"📈 Всего задач: %d\n"+
			"✅ Выполнено: %d (%.1f%%)\n"+
			"⏳ Ожидает: %d\n\n"+
			"🏷️ Популярные теги:\n",
		stats.Total, stats.Completed, stats.Rate, stats.Pending,
	)

	for i, tag := range stats.TopTags {
		statsText += fmt.Sprintf("   %d. %s (%d задач)\n", i+1, tag.Name, tag.Count)
	}

	if stats.BestDay != "" {
		statsText += fmt.Sprintf("\n🔥 Самый продуктивный день: %s\n", stats.BestDay)
	}

	h.sendMarkdown(chatID, statsText)
}

func (h *TelegramHandler) buildTaskMessage(task *model.Task, showCompleteButton bool) (string, *tgbotapi.InlineKeyboardMarkup) {
	if task == nil {
		return "", nil
	}

	statusEmoji := "❌"
	statusText := "НЕ ВЫПОЛНЕНА"
	if task.Completed {
		statusEmoji = "✅"
		statusText = "ВЫПОЛНЕНА"
	}

	var sb strings.Builder

	sb.WriteString("┌────────────────────────────────────────┐\n")
	sb.WriteString(fmt.Sprintf("│ 📋 *Задача #%d*\n", task.UserTaskID))
	sb.WriteString(fmt.Sprintf("│ %s *%s*\n", statusEmoji, statusText))
	sb.WriteString("├────────────────────────────────────────┤\n")

	title := utils.EscapeMarkdown(task.Title)

	title = utils.Truncate(title, 40)
	sb.WriteString(fmt.Sprintf("│ 📝 *Название:* %s\n", title))

	if task.Description != "" {
		desc := utils.EscapeMarkdown(task.Description)

		desc = utils.Truncate(desc, 40)
		sb.WriteString(fmt.Sprintf("│ 📄 *Описание:* %s\n", desc))
	}

	if len(task.Tags) > 0 {
		tagsStr := strings.Join(task.Tags, " ")

		tagsStr = utils.Truncate(tagsStr, 35)
		sb.WriteString(fmt.Sprintf("│ 🏷️ *Теги:* `%s`\n", tagsStr))
	}

	sb.WriteString(fmt.Sprintf("│ ⏰ *Создана:* %s\n", task.CreatedAt.Format(dateFormat)))

	if task.Completed && task.CompletedAt != nil {
		sb.WriteString(fmt.Sprintf("│ ✅ *Завершена:* %s\n", task.CompletedAt.Format(dateFormat)))
	}

	if task.DueDate != nil {
		sb.WriteString(fmt.Sprintf("│ 📅 *Срок:* %s\n", task.DueDate.Format(dateFormat)))
	}

	if task.ReminderOffset != nil && *task.ReminderOffset != "" {
		sb.WriteString(fmt.Sprintf("│ 🔔 *Напоминание:* за %s до срока\n", *task.ReminderOffset))
	}

	sb.WriteString("└────────────────────────────────────────┘")

	var rows [][]tgbotapi.InlineKeyboardButton

	if showCompleteButton && !task.Completed {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Выполнить", fmt.Sprintf("complete_%d", task.UserTaskID)),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_%d", task.UserTaskID)),
		))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_%d", task.UserTaskID)),
		))
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return sb.String(), &inlineKeyboard
}

func (h *TelegramHandler) sendTaskCard(chatID int64, task *model.Task) {
	text, keyboard := h.buildTaskMessage(task, !task.Completed)
	if text == "" {
		return
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}

func (h *TelegramHandler) updateTaskMessage(chatID int64, messageID int, task *model.Task) {
	text, keyboard := h.buildTaskMessage(task, false)
	if text == "" {
		return
	}

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	editMsg.ReplyMarkup = keyboard
	h.bot.Send(editMsg)
}

func (h *TelegramHandler) SendReminderMessage(chatID int64, task *model.Task, timeLeftText string) {
	message := fmt.Sprintf(
		"🔔 *Напоминание о задаче!*\n\n"+
			"📋 *Задача #%d:* %s\n"+
			"%s\n\n"+
			"✅ Выполнить: `/complete %d`\n"+
			"🗑️ Удалить: `/delete %d`",
		task.UserTaskID,
		utils.EscapeMarkdown(task.Title),
		timeLeftText,
		task.UserTaskID,
		task.UserTaskID,
	)

	if task.Description != "" {
		message += fmt.Sprintf("\n📝 %s", utils.EscapeMarkdown(utils.Truncate(task.Description, 100)))
	}

	if task.DueDate != nil {
		message += fmt.Sprintf("\n\n📅 *Дедлайн:* %s", task.DueDate.Format("02.01.2006 15:04"))
	}

	h.sendMarkdown(chatID, message)
}
