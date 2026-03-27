package handler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"TaskTracker/internal/model"
	"TaskTracker/internal/service"
)

type CLIHandler struct {
	baseService *service.TaskService
	pgService   *service.PostgresTaskService
	reader      *bufio.Reader
}

// Конструктор для базового сервиса (in-memory, JSON)
func NewCLIHandler(service *service.TaskService) *CLIHandler {
	return &CLIHandler{
		baseService: service,
		pgService:   nil,
		reader:      bufio.NewReader(os.Stdin),
	}
}

// Конструктор для PostgreSQL сервиса
func NewPostgresCLIHandler(service *service.PostgresTaskService) *CLIHandler {
	return &CLIHandler{
		baseService: service.TaskService,
		pgService:   service,
		reader:      bufio.NewReader(os.Stdin),
	}
}

func (h *CLIHandler) Run() {
	for {
		h.printMenu()
		choice := h.readInput()

		if !h.handleChoice(choice) {
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
	if h.pgService != nil {
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

func (h *CLIHandler) handleChoice(choice string) bool {
	switch choice {
	case "1":
		h.handleAdd()
	case "2":
		h.handleList()
	case "3":
		h.handleFindById()
	case "4":
		h.handleFineByTag()
	case "5":
		h.handleComplete()
	case "6":
		// Удаление доступно только для PostgreSQL
		if h.pgService != nil {
			h.handleDelete()
		} else {
			fmt.Println("❌ Удаление задач недоступно для этого типа хранилища")
		}
	case "7":
		// Статистика доступна только для PostgreSQL
		if h.pgService != nil {
			h.handleGetStats()
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

func (h *CLIHandler) handleAdd() {
	fmt.Println("=======================")
	fmt.Print("Введите название задачи: ")
	title := h.readInput()

	fmt.Print("Введите описание вашей задачи: ")
	description := h.readInput()

	fmt.Print("Введите теги (через запятую/пробел): ")
	tags := h.parseTags(h.readInput())

	task, err := h.baseService.AddTask(title, description, tags)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}
	fmt.Printf("\n✅ Задача '%s' добавлена с ID %d\n", task.Title, task.ID)
}

func (h *CLIHandler) handleDelete() {
	id := h.readID()
	err := h.pgService.DeleteTask(id)
	if err != nil {
		fmt.Println("❌ Ошибка: ", err)
		return
	}
	fmt.Printf("✅ Задача с ID %d удалена\n", id)
}

func (h *CLIHandler) handleList() {
	tasks, err := h.baseService.GetAllTasks()
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

func (h *CLIHandler) handleFindById() {
	id := h.readID()

	task, err := h.baseService.GetTaskById(id)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	h.printTask(task)
}

func (h *CLIHandler) handleFineByTag() {
	fmt.Print("Введите тег для поиска: ")
	tag := strings.TrimSpace(h.readInput())

	tasks, err := h.baseService.GetTasksByTag(tag)
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

func (h *CLIHandler) handleComplete() {
	id := h.readID()

	task, err := h.baseService.CompleteTask(id)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}
	fmt.Printf("✅ Задача '%s' отмечена как выполненная\n", task.Title)
}

func (h *CLIHandler) handleGetStats() {
	stats, err := h.pgService.GetStats()
	if err != nil {
		fmt.Printf("❌ Ошибка при получении статистики: %v\n", err)
		return
	}

	if len(stats) == 0 {
		fmt.Println("📊 Статистика недоступна или нет данных")
		return
	}

	fmt.Println("\n📊 ТОП-3 ПОПУЛЯРНЫХ ТЕГА:")
	for i, stat := range stats {
		fmt.Printf("  %d. %s\n", i+1, stat)
	}
	fmt.Println()
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
