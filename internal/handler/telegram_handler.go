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
)

// Состояния диалога
const (
	StateIdle                = ""
	StateAwaitingTitle       = "awaiting_title"
	StateAwaitingDescription = "awaiting_description"
	StateAwaitingTags        = "awaiting_tags"
)

type PendingTask struct {
	Title       string
	Description string
	Tags        []string
}

type TelegramHandler struct {
	bot         *tgbotapi.BotAPI
	taskService *service.TaskService
	extService  service.ExtendedTaskService

	userStates   map[int64]string
	pendingTasks map[int64]*PendingTask
	mu           sync.RWMutex
}

func NewTelegramHandler(bot *tgbotapi.BotAPI, taskService *service.TaskService, extService service.ExtendedTaskService) *TelegramHandler {
	return &TelegramHandler{
		bot:          bot,
		taskService:  taskService,
		extService:   extService,
		userStates:   make(map[int64]string),
		pendingTasks: make(map[int64]*PendingTask),
	}
}

// ========== УПРАВЛЕНИЕ СОСТОЯНИЯМИ ==========

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

// ========== ЗАПУСК И GRACEFUL SHUTDOWN ==========

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

// ========== ОБРАБОТКА СООБЩЕНИЙ ==========

func (h *TelegramHandler) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	text := strings.TrimSpace(msg.Text)
	chatID := msg.Chat.ID

	currentState := h.getState(chatID)

	if currentState != StateIdle {
		h.handleDialogState(ctx, chatID, text, currentState)
		return
	}

	switch {
	case text == "/start" || text == "/help":
		h.sendHelpMessage(chatID)

	case text == "/add":
		h.startAddTaskDialog(chatID)

	case strings.HasPrefix(text, "/add "):
		h.handleFastAdd(ctx, chatID, text)

	case text == "/list":
		h.handleListCommand(ctx, chatID)

	case strings.HasPrefix(text, "/complete"):
		h.handleCompleteCommand(ctx, chatID, text)

	case strings.HasPrefix(text, "/delete"):
		h.handleDeleteCommand(ctx, chatID, text)

	case strings.HasPrefix(text, "/find"):
		h.handleFindCommand(ctx, chatID, text)

	case text == "/stats":
		h.handleStatsCommand(ctx, chatID)

	case strings.HasPrefix(text, "/tag"):
		h.handleTagCommand(ctx, chatID, text)

	default:
		h.sendMessage(chatID, "❌ Неизвестная команда. Введите /help", false)
	}
}

// ========== ОБРАБОТКА INLINE-КНОПОК ==========

func (h *TelegramHandler) handleCallback(ctx context.Context, query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID
	data := query.Data

	h.bot.Send(tgbotapi.NewCallback(query.ID, ""))

	switch {
	case strings.HasPrefix(data, "complete_"):
		idStr := strings.TrimPrefix(data, "complete_")
		id, _ := strconv.Atoi(idStr)

		task, err := h.taskService.CompleteTask(ctx, id)
		if err != nil {
			h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err), false)
			return
		}
		h.updateTaskMessage(chatID, messageID, task)

	case strings.HasPrefix(data, "delete_"):
		if h.extService == nil {
			h.sendMessage(chatID, "❌ Удаление недоступно", false)
			return
		}

		idStr := strings.TrimPrefix(data, "delete_")
		id, _ := strconv.Atoi(idStr)

		if err := h.extService.DeleteTask(ctx, id); err != nil {
			h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err), false)
			return
		}

		deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
		h.bot.Send(deleteMsg)
		h.sendMessage(chatID, fmt.Sprintf("🗑️ Задача #%d удалена", id), false)
	}
}

// ========== ДИАЛОГ ДОБАВЛЕНИЯ ЗАДАЧИ ==========

func (h *TelegramHandler) startAddTaskDialog(chatID int64) {
	h.setState(chatID, StateAwaitingTitle)
	h.sendMessage(chatID, "📝 Введите название задачи:", false)
}

func (h *TelegramHandler) handleDialogState(ctx context.Context, chatID int64, input string, state string) {
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

		task, err := h.taskService.AddTask(ctx, pending.Title, pending.Description, pending.Tags)
		h.clearState(chatID)

		if err != nil {
			h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка при создании задачи: %v", err), false)
			return
		}

		h.sendTaskCard(chatID, task)
		h.sendMessage(chatID, "✅ Задача успешно создана! Используйте /list для просмотра.", false)
	}
}

// ========== ОТПРАВКА СООБЩЕНИЙ ==========

func (h *TelegramHandler) sendMessage(chatID int64, text string, markdown bool) {
	msg := tgbotapi.NewMessage(chatID, text)
	if markdown {
		msg.ParseMode = "Markdown"
	}

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка отправки сообщения: %v", err)
		// Попробуем отправить без Markdown
		if markdown {
			plainMsg := tgbotapi.NewMessage(chatID, text)
			h.bot.Send(plainMsg)
		}
	}
}

func (h *TelegramHandler) sendMarkdown(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
}

// ========== КОМАНДЫ ==========

func (h *TelegramHandler) sendHelpMessage(chatID int64) {
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
	h.sendMarkdown(chatID, helpText)
}

func (h *TelegramHandler) handleFastAdd(ctx context.Context, chatID int64, text string) {
	rawTitle := strings.TrimPrefix(text, "/add ")
	rawTitle = strings.TrimSpace(rawTitle)

	title, tags := h.parseTitleAndTags(rawTitle)

	task, err := h.taskService.AddTask(ctx, title, "", tags)
	if err != nil {
		h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err), false)
		return
	}

	h.sendTaskCard(chatID, task)
}

func (h *TelegramHandler) handleListCommand(ctx context.Context, chatID int64) {
	tasks, err := h.taskService.GetAllTasks(ctx)
	if err != nil {
		h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err), false)
		return
	}

	if len(tasks) == 0 {
		h.sendMessage(chatID, "📭 У вас нет задач.", false)
		return
	}

	// Отправляем счетчик
	h.sendMessage(chatID, fmt.Sprintf("📊 Найдено задач: %d", len(tasks)), false)

	// Отправляем каждую задачу, НО не прерываем цикл при ошибке
	successCount := 0
	for i, task := range tasks {
		log.Printf("🔄 Отправка задачи %d/%d (ID=%d)", i+1, len(tasks), task.ID)

		err := h.sendTaskCardSafe(chatID, task)
		if err != nil {
			log.Printf("⚠️ Не удалось отправить задачу #%d: %v", task.ID, err)
			// Отправляем упрощенную версию
			fallbackMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Задача #%d не отображается: %s", task.ID, task.Title))
			h.bot.Send(fallbackMsg)
		} else {
			successCount++
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Отправляем итог
	completed := 0
	for _, t := range tasks {
		if t.Completed {
			completed++
		}
	}
	summary := fmt.Sprintf("📊 *Итого:* %d задач | ✅ Выполнено: %d | ⏳ В работе: %d\n✅ Успешно отправлено: %d",
		len(tasks), completed, len(tasks)-completed, successCount)
	h.sendMarkdown(chatID, summary)
}

func (h *TelegramHandler) sendTaskCardSafe(chatID int64, task *model.Task) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🔥 PANIC в sendTaskCard для задачи #%d: %v", task.ID, r)
		}
	}()

	if task == nil {
		return fmt.Errorf("task is nil")
	}

	return h.sendTaskCard(chatID, task)
}

func (h *TelegramHandler) handleCompleteCommand(ctx context.Context, chatID int64, text string) {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 {
		h.sendMessage(chatID, "❌ Использование: /complete <ID>", false)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendMessage(chatID, "❌ ID должен быть числом", false)
		return
	}

	task, err := h.taskService.CompleteTask(ctx, id)
	if err != nil {
		h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err), false)
		return
	}

	h.sendMessage(chatID, fmt.Sprintf("✅ Задача \"%s\" выполнена!", task.Title), false)
}

func (h *TelegramHandler) handleDeleteCommand(ctx context.Context, chatID int64, text string) {
	if h.extService == nil {
		h.sendMessage(chatID, "❌ Удаление недоступно", false)
		return
	}

	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		h.sendMessage(chatID, "❌ Использование: /delete <ID>", false)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendMessage(chatID, "❌ ID должен быть числом", false)
		return
	}

	if err := h.extService.DeleteTask(ctx, id); err != nil {
		h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err), false)
		return
	}

	h.sendMessage(chatID, fmt.Sprintf("🗑️ Задача #%d удалена", id), false)
}

func (h *TelegramHandler) handleFindCommand(ctx context.Context, chatID int64, text string) {
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		h.sendMessage(chatID, "❌ Использование: /find <ID>", false)
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		h.sendMessage(chatID, "❌ ID должен быть числом", false)
		return
	}

	task, err := h.taskService.GetTaskById(ctx, id)
	if err != nil {
		h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err), false)
		return
	}

	h.sendTaskCard(chatID, task)
}

func (h *TelegramHandler) handleTagCommand(ctx context.Context, chatID int64, text string) {
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		h.sendMessage(chatID, "❌ Использование: /tag <тег>", false)
		return
	}

	tag := parts[1]
	tasks, err := h.taskService.GetTasksByTag(ctx, tag)
	if err != nil {
		h.sendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err), false)
		return
	}

	if len(tasks) == 0 {
		h.sendMessage(chatID, fmt.Sprintf("📭 Нет задач с тегом \"%s\"", tag), false)
		return
	}

	for _, task := range tasks {
		h.sendTaskCard(chatID, task)
		time.Sleep(50 * time.Millisecond)
	}

	summary := fmt.Sprintf("📊 *Найдено задач с тегом \"%s\":* %d", tag, len(tasks))
	h.sendMarkdown(chatID, summary)
}

func (h *TelegramHandler) handleStatsCommand(ctx context.Context, chatID int64) {
	if h.extService == nil {
		h.sendMessage(chatID, "❌ Статистика доступна только при использовании PostgreSQL", false)
		return
	}

	// Сначала показываем, что обновляем
	msg := h.sendMessageTemp(chatID, "📊 Обновляю статистику...")

	// Принудительно обновляем статистику
	stats, err := h.extService.GetStatsForce(ctx)
	if err != nil {
		h.editMessage(chatID, msg.MessageID, fmt.Sprintf("❌ Ошибка: %v", err))
		return
	}

	if stats == nil || stats.Total == 0 {
		h.editMessage(chatID, msg.MessageID, "📭 Нет данных для статистики")
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

	h.editMessage(chatID, msg.MessageID, statsText)
}

// sendMessageTemp отправляет временное сообщение и возвращает его
func (h *TelegramHandler) sendMessageTemp(chatID int64, text string) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatID, text)
	sentMsg, _ := h.bot.Send(msg)
	return sentMsg
}

// editMessage редактирует существующее сообщение
func (h *TelegramHandler) editMessage(chatID int64, messageID int, text string) {
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	h.bot.Send(editMsg)
}

// ========== ОТОБРАЖЕНИЕ ЗАДАЧ ==========

func (h *TelegramHandler) sendTaskCard(chatID int64, task *model.Task) error {
	log.Printf("📤 Начинаю отправку задачи #%d", task.ID)

	statusEmoji := "❌"
	statusText := "НЕ ВЫПОЛНЕНА"
	if task.Completed {
		statusEmoji = "✅"
		statusText = "ВЫПОЛНЕНА"
	}

	var sb strings.Builder

	// Потенциально опасное место: экранирование
	title := escapeMarkdown(task.Title)
	log.Printf("  Название: %s", title)

	if len(title) > 40 {
		title = title[:37] + "..."
	}

	sb.WriteString("┌────────────────────────────────────────┐\n")
	sb.WriteString(fmt.Sprintf("│ 📋 *Задача #%d*\n", task.ID))
	sb.WriteString(fmt.Sprintf("│ %s *%s*\n", statusEmoji, statusText))
	sb.WriteString("├────────────────────────────────────────┤\n")
	sb.WriteString(fmt.Sprintf("│ 📝 *Название:* %s\n", title))

	// Опасное место: описание
	if task.Description != "" {
		desc := escapeMarkdown(task.Description)
		log.Printf("  Описание: %s", desc)
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ 📄 *Описание:* %s\n", desc))
	}

	// Опасное место: теги
	if len(task.Tags) > 0 {
		tagsStr := strings.Join(task.Tags, " ")
		log.Printf("  Теги: %s", tagsStr)
		if len(tagsStr) > 35 {
			tagsStr = tagsStr[:32] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ 🏷️ *Теги:* `%s`\n", tagsStr))
	}

	sb.WriteString(fmt.Sprintf("│ ⏰ *Создана:* %s\n", task.CreatedAt.Format("02.01.2006 15:04")))

	if task.Completed && task.CompletedAt != nil {
		sb.WriteString(fmt.Sprintf("│ ✅ *Завершена:* %s\n", task.CompletedAt.Format("02.01.2006 15:04")))
	}

	sb.WriteString("└────────────────────────────────────────┘")

	var rows [][]tgbotapi.InlineKeyboardButton
	if !task.Completed {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Выполнить", fmt.Sprintf("complete_%d", task.ID)),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_%d", task.ID)),
		))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_%d", task.ID)),
		))
	}

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, sb.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = inlineKeyboard

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка отправки: %v", err)
		return err
	}

	log.Printf("✅ Задача #%d отправлена", task.ID)
	return nil
}
func (h *TelegramHandler) updateTaskMessage(chatID int64, messageID int, task *model.Task) {
	var sb strings.Builder

	sb.WriteString("┌────────────────────────────────────────┐\n")
	sb.WriteString(fmt.Sprintf("│ 📋 *Задача #%d*\n", task.ID))
	sb.WriteString("│ ✅ *ВЫПОЛНЕНА*\n")
	sb.WriteString("├────────────────────────────────────────┤\n")

	title := escapeMarkdown(task.Title)
	if len(title) > 40 {
		title = title[:37] + "..."
	}
	sb.WriteString(fmt.Sprintf("│ 📝 *Название:* %s\n", title))

	if task.Description != "" {
		desc := escapeMarkdown(task.Description)
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ 📄 *Описание:* %s\n", desc))
	}

	if len(task.Tags) > 0 {
		tagsStr := strings.Join(task.Tags, " ")
		if len(tagsStr) > 35 {
			tagsStr = tagsStr[:32] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ 🏷️ *Теги:* `%s`\n", tagsStr))
	}

	sb.WriteString(fmt.Sprintf("│ ⏰ *Создана:* %s\n", task.CreatedAt.Format("02.01.2006 15:04")))

	if task.CompletedAt != nil {
		sb.WriteString(fmt.Sprintf("│ ✅ *Завершена:* %s\n", task.CompletedAt.Format("02.01.2006 15:04")))
	}

	sb.WriteString("└────────────────────────────────────────┘")

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_%d", task.ID)),
		),
	)

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, sb.String())
	editMsg.ParseMode = "Markdown"
	editMsg.ReplyMarkup = &inlineKeyboard
	h.bot.Send(editMsg)
}

// ========== ПАРСИНГ И ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ==========

func (h *TelegramHandler) parseTitleAndTags(raw string) (string, []string) {
	words := strings.Fields(raw)
	titleWords := []string{}
	tags := []string{}

	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			if tag != "" {
				tags = append(tags, tag)
			}
		} else {
			titleWords = append(titleWords, word)
		}
	}

	title := strings.Join(titleWords, " ")
	return title, tags
}

func (h *TelegramHandler) parseTags(input string) []string {
	words := strings.Fields(input)
	tags := []string{}
	seen := make(map[string]bool)

	for _, word := range words {
		tag := strings.TrimPrefix(word, "#")
		tag = strings.TrimSpace(tag)
		if tag != "" && !seen[tag] {
			seen[tag] = true
			tags = append(tags, tag)
		}
	}
	return tags
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}
