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
	service *service.TaskService
	reader  *bufio.Reader
}

func NewCLIHandler(service *service.TaskService) *CLIHandler {
	return &CLIHandler{
		service: service,
		reader:  bufio.NewReader(os.Stdin),
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
	// Просто используем фиксированную ширину, без центрирования
	// Но с красивым оформлением

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
		menuItemStyle.Render("6. 🚪 Выход"),
	}

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

	// Просто печатаем красивое меню без центрирования
	fmt.Println(borderStyle.Render(content))
	fmt.Print(" ")
}

func (h *CLIHandler) handleChoice(choice string) bool {
	switch choice {
	case "1":
		h.HandleAdd()
	case "2":
		h.handleList()
	case "3":
		h.handleFindById()
	case "4":
		h.handleFineByTag()
	case "5":
		h.handleComplete()
	case "6":
		fmt.Println("👋 До свидания!")
		return false
	default:
		fmt.Println("❌ Неверный пункт. Введите 1-6")
	}
	return true

}

func (h *CLIHandler) HandleAdd() {
	fmt.Println("=======================")
	fmt.Print("Введите название задачи: ")
	title := h.readInput()

	fmt.Print("Введите описание вашей задачи: ")
	description := h.readInput()

	fmt.Print("Введите теги (через запятую/пробел): ")
	tags := h.parseTags(h.readInput())

	task, err := h.service.AddTask(title, description, tags)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}
	fmt.Printf("\n✅ Задача '%s' добавлена с ID %d\n", task.Title, task.ID)
}

func (h *CLIHandler) handleList() {
	tasks, err := h.service.GetAllTasks()
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

	task, err := h.service.GetTaskById(id)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	h.printTask(task)
}

func (h *CLIHandler) handleFineByTag() {
	fmt.Print("Введите тег для поиска: ")
	tag := strings.TrimSpace(h.readInput())

	tasks, err := h.service.GetTasksByTag(tag)
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

	task, err := h.service.CompleteTask(id)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}
	fmt.Printf("✅ Задача '%s' отмечена как выполненная\n", task.Title)

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
			Padding(0, 2). // Уменьшил паддинг до 0 сверху/снизу
			Width(60)      // Уменьшил ширину

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1) // Добавил отступ после заголовка

	fieldStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Width(12) // Фиксированная ширина для поля

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))
)

func (h *CLIHandler) printTask(task *model.Task) {
	status := "❌ - Не выполнена"
	if task.Completed {
		status = "✅ - Выполнена"
	}
	// fmt.Println("``````````````````````````````````````````````````````````````````````````````````")
	// fmt.Printf("📂 ID: %d | ✏️ Название: %s | 🔄 Статус: %s\n", task.ID, task.Title, status)
	// fmt.Printf("📝 Описание: %s\n🏷️ Тэги: %s\n", task.Description, task.Tags)
	// fmt.Printf("⏰ Время создания: %s\n", task.CreatedAt.Format("02.01.2006 15:04:05"))
	// if task.Completed {
	// 	fmt.Printf("⏰ Время завершения: %s\n", task.CompletedAt.Format("02.01.2006 15:04:05"))
	// }
	// fmt.Println("``````````````````````````````````````````````````````````````````````````````````")

	// Собираем содержимое построчно
	content := strings.Builder{}

	// Заголовок
	content.WriteString(headerStyle.Render("📋 ДЕТАЛИ ЗАДАЧИ") + "\n\n")

	// Строка ID - эмодзи отдельно, значение с fieldStyle
	content.WriteString("📂 " + fieldStyle.Render("ID:") + " " +
		valueStyle.Render(fmt.Sprintf("%d", task.ID)) + "\n")

	// Строка Названия
	content.WriteString("✏️ " + fieldStyle.Render("Название:") + " " +
		valueStyle.Render(task.Title) + "\n")

	// Строка Статуса
	content.WriteString("🔄 " + fieldStyle.Render("Статус:") + " " +
		valueStyle.Render(status) + "\n")

	// Строка Описания
	content.WriteString("📝 " + fieldStyle.Render("Описание:") + " " +
		valueStyle.Render(task.Description) + "\n")

	// Строка Тэгов
	content.WriteString("🏷️ " + fieldStyle.Render("Тэги:") + " " +
		valueStyle.Render(fmt.Sprintf("%v", task.Tags)) + "\n")

	// Строка Создано
	content.WriteString("⏰ " + fieldStyle.Render("Создано:") + " " +
		valueStyle.Render(task.CreatedAt.Format("02.01.2006 15:04:05")) + "\n")

	// Строка Завершено - только если задача выполнена
	if task.Completed {
		content.WriteString("⏰ " + fieldStyle.Render("Завершено:") + " " +
			valueStyle.Render(task.CompletedAt.Format("02.01.2006 15:04:05")))
	}

	// Оборачиваем всё в рамку
	boxedContent := borderStyle.Render(content.String())
	fmt.Println(boxedContent)
}
