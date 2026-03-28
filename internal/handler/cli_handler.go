package handler

import (
	"bufio"
	"context"
	"fmt"
	"os"
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
			fmt.Println("❌ Удаление задач недоступно для этого типа хранилища")
		}
	case "7":
		// Статистика доступна только для PostgreSQL
		if h.extendedService != nil {

			h.handleStats(defCtx)
		} else {
			fmt.Println("❌ Статистика недоступна для этого типа хранилища")
		}
	case "8":
		fmt.Println("👋 До свидания!")
		return false
	default:
		fmt.Println("❌ Неверный пункт. Введите 1-8")
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
		fmt.Printf("❌ Ошибка: %v\n", err)
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
		fmt.Println("❌ Ошибка: ", err)
		return
	}
	fmt.Printf("✅ Задача с ID %d удалена\n", id)
}

func (h *CLIHandler) handleList(defCtx context.Context) {
	ctx, cancel := context.WithTimeout(defCtx, 5*time.Second)
	defer cancel()
	tasks, err := h.baseService.GetAllTasks(ctx)
	if err != nil {
		fmt.Println("❌ Ошибка: ", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("📭 Список задач пуст")
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
		fmt.Printf("❌ Ошибка: %v\n", err)
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
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Printf("📭 Задач с тегом '%s' не найдено\n", tag)
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
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}
	fmt.Printf("✅ Задача '%s' отмечена как выполненная\n", task.Title)
}

func (h *CLIHandler) handleStats(defCtx context.Context) {
	fmt.Println("\n📊 Загрузка статистики...")
	stats, lastUpdated, isUpdating, err := h.extendedService.GetStatsWithInfo(defCtx)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	if isUpdating {
		fmt.Println("🔄 Статистика обновляется в фоне...")
	}

	if time.Since(lastUpdated) > 5*time.Minute {
		fmt.Println("⚠️ Данные могут быть устаревшими")
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
	fmt.Printf("   ├─ Всего задач:      %d\n", stats.TotalTasks)
	fmt.Printf("   ├─ ✅ Выполнено:      %d (%.1f%%)\n",
		stats.CompletedTasks, stats.CompletionRate)
	fmt.Printf("   ├─ ⏳ Ожидает:        %d\n", stats.PendingTasks)
	fmt.Printf("   └─ 🏷️ Уникальных тегов: %d\n", stats.TotalUniqueTags)

	// 2. Временная статистика
	if stats.AvgCompletionTime > 0 {
		hours := int(stats.AvgCompletionTime)
		minutes := int((stats.AvgCompletionTime - float64(hours)) * 60)
		fmt.Printf("\n⏱️ ВРЕМЯ ВЫПОЛНЕНИЯ:\n")
		fmt.Printf("   └─ Среднее: %d ч %d мин\n", hours, minutes)
	}

	// 3. Топ-теги
	if len(stats.TopTags) > 0 {
		fmt.Println("\n🏷️ ПОПУЛЯРНЫЕ ТЕГИ:")
		for i, tag := range stats.TopTags {
			bar := h.createBar(tag.UsageCount, stats.TotalTasks, 20)
			fmt.Printf("   %d. %-12s %s (%d задач)\n",
				i+1, tag.Tag, bar, tag.UsageCount)
		}
	}

	// 4. Активность по дням (график)
	if len(stats.TasksByDay) > 0 {
		fmt.Println("\n📅 АКТИВНОСТЬ ПО ДНЯМ (последние 7 дней):")
		maxCount := h.getMaxCount(stats.TasksByDay)

		// Показываем только последние 7 дней
		startIdx := len(stats.TasksByDay) - 7
		if startIdx < 0 {
			startIdx = 0
		}

		for i := startIdx; i < len(stats.TasksByDay); i++ {
			day := stats.TasksByDay[i]
			bar := h.createBar(day.Count, maxCount, 30)
			fmt.Printf("   %s %s (%d)\n", day.Date, bar, day.Count)
		}
	}

	// 5. Самый продуктивный день
	if stats.MostProductiveDay != "" {
		fmt.Printf("\n🔥 САМЫЙ ПРОДУКТИВНЫЙ ДЕНЬ: %s\n", stats.MostProductiveDay)
	}

	// 6. Распределение по часам
	if len(stats.TasksByHour) > 0 {
		fmt.Println("\n⏰ АКТИВНОСТЬ ПО ЧАСАМ:")
		maxCount := h.getMaxHourCount(stats.TasksByHour)

		// Создаем график для 24 часов
		hourMap := make(map[int]int)
		for _, h := range stats.TasksByHour {
			hourMap[h.Hour] = h.Count
		}

		for hour := 0; hour < 24; hour++ {
			count := hourMap[hour]
			if count > 0 {
				bar := h.createBar(count, maxCount, 30)
				fmt.Printf("   %02d:00 %s (%d задач)\n", hour, bar, count)
			}
		}
	}

	// 7. Распределение по дням недели
	if len(stats.TasksByWeekday) > 0 {
		fmt.Println("\n📆 РАСПРЕДЕЛЕНИЕ ПО ДНЯМ НЕДЕЛИ:")
		maxCount := h.getMaxWeekdayCount(stats.TasksByWeekday)

		for _, wd := range stats.TasksByWeekday {
			bar := h.createBar(wd.Count, maxCount, 25)
			fmt.Printf("   %-10s %s (%d задач)\n", wd.Weekday, bar, wd.Count)
		}
	}

	// 8. Статистика выполнения по дням
	if len(stats.CompletionByDay) > 0 {
		fmt.Println("\n✅ ВЫПОЛНЕННЫЕ ЗАДАЧИ ПО ДНЯМ (последние 7 дней):")
		startIdx := len(stats.CompletionByDay) - 7
		if startIdx < 0 {
			startIdx = 0
		}

		for i := startIdx; i < len(stats.CompletionByDay); i++ {
			day := stats.CompletionByDay[i]
			fmt.Printf("   %s: %d задач\n", day.Date, day.Count)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
}

// Всякая доп хуйня
func (h *CLIHandler) readInput() string {
	input, _ := h.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (h *CLIHandler) readID() int {
	fmt.Print("Введите ID задачи: ")
	idStr := h.readInput()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
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

func (h *CLIHandler) getMaxCount(dailyStats []model.DailyStat) int {
	max := 0
	for _, ds := range dailyStats {
		if ds.Count > max {
			max = ds.Count
		}
	}
	return max
}

func (h *CLIHandler) getMaxHourCount(hourlyStats []model.HourlyStat) int {
	max := 0
	for _, hs := range hourlyStats {
		if hs.Count > max {
			max = hs.Count
		}
	}
	return max
}

func (h *CLIHandler) getMaxWeekdayCount(weekdayStats []model.WeekdayStat) int {
	max := 0
	for _, ws := range weekdayStats {
		if ws.Count > max {
			max = ws.Count
		}
	}
	return max
}
