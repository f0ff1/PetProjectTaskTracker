package handler

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"TaskTracker/internal/model"
	"TaskTracker/internal/service"
)

type CLIHandler struct {
	baseService     *service.TaskService
	extendedService service.ExtendedTaskService
	reader          *bufio.Reader
}

func NewCLIHandler(service *service.TaskService, extendedSvc service.ExtendedTaskService) *CLIHandler {
	return &CLIHandler{
		baseService:     service,
		extendedService: extendedSvc,
		reader:          bufio.NewReader(os.Stdin),
	}
}

func (h *CLIHandler) Run(ctx context.Context) {
	for {
		h.printMenu()
		choice := h.readInput()
		if !h.handleChoice(ctx, choice) {
			break
		}
	}
}

func (h *CLIHandler) printMenu() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("236")).
		Padding(0, 2).
		MarginTop(1).
		MarginBottom(1).
		Align(lipgloss.Center).
		Width(40)

	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 3).
		Width(46)

	menuItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		PaddingLeft(2)

	menuItems := []string{
		menuItemStyle.Render("1. ➕ Добавить задачу"),
		menuItemStyle.Render("2. 📋 Список всех задач"),
		menuItemStyle.Render("3. 📋 Найти задачу по ID"),
		menuItemStyle.Render("4. 🏷️ Найти задачу по Тэгу"),
		menuItemStyle.Render("5. ✏️ Отметить задачу как выполненную"),
	}

	// Добавляем пункты меню, которые доступны только для PostgreSQL
	if h.extendedService != nil {
		menuItems = append(menuItems,
			menuItemStyle.Render("6. 🗑️ Удалить задачу по ID"),
			menuItemStyle.Render("7. 📊 Статистика"),
		)
	}

	menuItems = append(menuItems, menuItemStyle.Render("8. 🚪 Выход"))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("📋 МЕНЕДЖЕР ЗАДАЧ"),
		"",
		lipgloss.JoinVertical(lipgloss.Left, menuItems...),
		"",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Render("⚡ Выберите действие:"),
	)

	fmt.Println(borderStyle.Render(content))
	fmt.Print(" ")
}

func (h *CLIHandler) handleChoice(defCtx context.Context, choice string) bool {
	switch choice {
	case "1":

		h.handleAdd(defCtx)
	case "2":

		h.handleList(defCtx)
	case "3":

		h.handleFindById(defCtx)
	case "4":

		h.handleFineByTag(defCtx)
	case "5":

		h.handleComplete(defCtx)
	case "6":
		// Удаление доступно только для PostgreSQL
		if h.extendedService != nil {
			h.handleDelete(defCtx)
		} else {
			slog.Warn("Удаление задач недоступно для этого типа хранилища")
		}
	case "7":
		// Статистика доступна только для PostgreSQL
		if h.extendedService != nil {
			h.handleStats(defCtx)
		} else {
			slog.Warn("Статистика недоступна для этого типа хранилища")
		}
	case "8":
		slog.Info("До свидания")
		return false
	default:
		slog.Warn("Неверный пункт. Введите 1-8")
	}
	return true
}

func (h *CLIHandler) handleAdd(defCtx context.Context) {
	fmt.Println("=======================")
	fmt.Print("Введите название задачи: ")
	title := h.readInput()

	fmt.Print("Введите описание вашей задачи: ")
	description := h.readInput()

	fmt.Print("Введите теги (через запятую/пробел): ")
	tags := h.parseTags(h.readInput())

	ctx, cancel := context.WithTimeout(defCtx, 1*time.Second)
	defer cancel()

	task, err := h.baseService.AddTask(ctx, title, description, tags)
	if err != nil {
		slog.Error("Ошибка добавления задачи", "err", err)
		return
	}
	fmt.Printf("\n✅ Задача '%s' добавлена с ID %d\n", task.Title, task.ID)
}

func (h *CLIHandler) handleDelete(defCtx context.Context) {
	id := h.readID()
	ctx, cancel := context.WithTimeout(defCtx, 1*time.Second)
	defer cancel()
	err := h.extendedService.DeleteTask(ctx, id)
	if err != nil {
		slog.Error("Ошибка удаления задачи", "err", err, "id", id)
		return
	}
	fmt.Printf("✅ Задача с ID %d удалена\n", id)
}

func (h *CLIHandler) handleList(defCtx context.Context) {
	ctx, cancel := context.WithTimeout(defCtx, 5*time.Second)
	defer cancel()
	tasks, err := h.baseService.GetAllTasks(ctx)
	if err != nil {
		slog.Error("Ошибка получения списка задач", "err", err)
		return
	}

	if len(tasks) == 0 {
		slog.Info("Список задач пуст")
		return
	}

	for _, task := range tasks {
		h.printTask(task)
	}
}

func (h *CLIHandler) handleFindById(defCtx context.Context) {
	id := h.readID()
	ctx, cancel := context.WithTimeout(defCtx, 2*time.Second)
	defer cancel()
	task, err := h.baseService.GetTaskById(ctx, id)
	if err != nil {
		slog.Error("Ошибка при поиске задачи по ID", "err", err, "id", id)
		return
	}

	h.printTask(task)
}

func (h *CLIHandler) handleFineByTag(defCtx context.Context) {
	fmt.Print("Введите тег для поиска: ")
	tag := strings.TrimSpace(h.readInput())
	ctx, cancel := context.WithTimeout(defCtx, 1*time.Second)
	defer cancel()
	tasks, err := h.baseService.GetTasksByTag(ctx, tag)
	if err != nil {
		slog.Error("Ошибка при поиске задач по тегу", "err", err, "tag", tag)
		return
	}

	if len(tasks) == 0 {
		slog.Info("Задач с тегом не найдено", "tag", tag)
		return
	}

	for _, task := range tasks {
		h.printTask(task)
	}
}

func (h *CLIHandler) handleComplete(defCtx context.Context) {
	id := h.readID()

	ctx, cancel := context.WithTimeout(defCtx, 1*time.Second)
	defer cancel()

	task, err := h.baseService.CompleteTask(ctx, id)
	if err != nil {
		slog.Error("Ошибка при отмечании задачи как выполненной", "err", err, "id", id)
		return
	}
	fmt.Printf("✅ Задача '%s' отмечена как выполненная\n", task.Title)
}

func (h *CLIHandler) handleStats(defCtx context.Context) {
	slog.Info("Загрузка статистики")
	stats, lastUpdated, isUpdating, err := h.extendedService.GetStatsWithInfo(defCtx)
	if err != nil {
		slog.Error("Ошибка получения статистики", "err", err)
		return
	}

	if isUpdating {
		slog.Info("Статистика обновляется в фоне")
	}

	if time.Since(lastUpdated) > 5*time.Minute {
		slog.Warn("Данные статистики могут быть устаревшими", "lastUpdated", lastUpdated)
	}

	fmt.Printf("📅 Последнее обновление: %s\n", lastUpdated.Format("15:04:05"))

	h.printStats(stats)
}

func (h *CLIHandler) printStats(stats *model.TaskStats) {

	// Заголовок
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊  ДЕТАЛЬНАЯ СТАТИСТИКА ЗАДАЧ  📊")
	fmt.Println(strings.Repeat("=", 60))

	// 1. Основная статистика
	fmt.Println("\n📈 ОСНОВНАЯ СТАТИСТИКА:")
	fmt.Printf("   ├─ Всего задач:      %d\n", stats.Total)
	fmt.Printf("   ├─ ✅ Выполнено:      %d (%.1f%%)\n", stats.Completed, stats.Rate)
	fmt.Printf("   ├─ ⏳ Ожидает:        %d\n", stats.Pending)

	// 2. Топ-теги
	if len(stats.TopTags) > 0 {
		fmt.Println("\n🏷️ ПОПУЛЯРНЫЕ ТЕГИ:")
		for i, tag := range stats.TopTags {
			bar := h.createBar(tag.Count, stats.Total, 20)
			fmt.Printf("   %d. %-12s %s (%d задач)\n", i+1, tag.Name, bar, tag.Count)
		}
	} else {
		fmt.Println("\n🏷️ ПОПУЛЯРНЫЕ ТЕГИ:")
		fmt.Println("   └─ Нет тегов")
	}

	// 4. Самый продуктивный день
	if stats.BestDay != "" {
		fmt.Printf("\n🔥 САМЫЙ ПРОДУКТИВНЫЙ ДЕНЬ: %s\n", stats.BestDay)
	}

	// 5. Активность за последние 7 дней
	if len(stats.Last7Days) > 0 {
		fmt.Println("\n📅 АКТИВНОСТЬ ПО ДНЯМ (последние 7 дней):")

		// Находим максимальное значение для масштабирования графика
		maxCount := 0
		for _, count := range stats.Last7Days {
			if count > maxCount {
				maxCount = count
			}
		}

		// Сортируем дни (от старых к новым или наоборот)
		days := make([]string, 0, len(stats.Last7Days))
		for day := range stats.Last7Days {
			days = append(days, day)
		}
		sort.Strings(days) // сортируем по дате (формат YYYY-MM-DD)

		// Показываем последние 7 дней
		startIdx := 0
		if len(days) > 7 {
			startIdx = len(days) - 7
		}

		for i := startIdx; i < len(days); i++ {
			day := days[i]
			count := stats.Last7Days[day]
			bar := h.createBar(count, maxCount, 30)

			// Форматируем дату для красивого вывода
			formattedDay := day
			if t, err := time.Parse("2006-01-02", day); err == nil {
				formattedDay = t.Format("02 Jan")
			}
			fmt.Printf("   %-10s %s (%d)\n", formattedDay, bar, count)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
}

// Всякая доп хуйня
func (h *CLIHandler) readInput() string {
	input, err := h.reader.ReadString('\n')
	if err != nil {
		slog.Error("Ошибка чтения ввода", "err", err)
		return ""
	}
	return strings.TrimSpace(input)
}

func (h *CLIHandler) readID() int {
	fmt.Print("Введите ID задачи: ")
	idStr := h.readInput()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("Ошибка преобразования ID", "err", err, "input", idStr)
		return 0
	}

	return id
}

func (h *CLIHandler) parseTags(input string) []string {
	replacer := strings.NewReplacer(",", " ", ";", " ", "|", " ")
	normalized := replacer.Replace(input)
	return strings.Fields(normalized)
}

var (
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 2).
			Width(60)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	fieldStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Width(12)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))
)

func (h *CLIHandler) printTask(task *model.Task) {
	status := "❌ - Не выполнена"
	if task.Completed {
		status = "✅ - Выполнена"
	}

	content := strings.Builder{}

	content.WriteString(headerStyle.Render("📋 ДЕТАЛИ ЗАДАЧИ") + "\n\n")

	content.WriteString("📂 " + fieldStyle.Render("ID:") + " " +
		valueStyle.Render(fmt.Sprintf("%d", task.ID)) + "\n")

	content.WriteString("✏️ " + fieldStyle.Render("Название:") + " " +
		valueStyle.Render(task.Title) + "\n")

	content.WriteString("🔄 " + fieldStyle.Render("Статус:") + " " +
		valueStyle.Render(status) + "\n")

	content.WriteString("📝 " + fieldStyle.Render("Описание:") + " " +
		valueStyle.Render(task.Description) + "\n")

	content.WriteString("🏷️ " + fieldStyle.Render("Тэги:") + " " +
		valueStyle.Render(fmt.Sprintf("%v", task.Tags)) + "\n")

	content.WriteString("⏰ " + fieldStyle.Render("Создано:") + " " +
		valueStyle.Render(task.CreatedAt.Format("02.01.2006 15:04:05")) + "\n")

	if task.CompletedAt != nil {
		if task.Completed {
			content.WriteString("⏰ " + fieldStyle.Render("Завершено:") + " " +
				valueStyle.Render(task.CompletedAt.Format("02.01.2006 15:04:05")))
		}
	} else {
		content.WriteString("⏰ " + fieldStyle.Render("Не завершено"))
	}

	boxedContent := borderStyle.Render(content.String())
	fmt.Println(boxedContent)
}

func (h *CLIHandler) createBar(value, max int, width int) string {
	if max == 0 {
		return ""
	}

	filled := int(float64(value) / float64(max) * float64(width))
	if filled < 1 && value > 0 {
		filled = 1
	}

	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
	return bar
}
