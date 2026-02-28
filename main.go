package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time // указатель, nil если не выполнена
}

var tasks = make(map[int]*Task)
var nextID = 1

func addTask(reader *bufio.Reader) {
	fmt.Print("Введите название задачи: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Введите описание вашей задачи: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	task := &Task{
		ID:          nextID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	tasks[nextID] = task
	fmt.Printf("✅ Задача: '%s' добавлена успешно.\n", task.Title)
	nextID++
}

func readID(reader *bufio.Reader) (int, error) {
	strId, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(strId))
	if err != nil || id < 1 {
		return 0, fmt.Errorf("некорректный ID")
	}
	return id, nil
}

func printTaskData(task *Task) {
	status := "❌ - Не выполнена"
	if task.Completed {
		status = "✅ - Выполнена"
	}
	fmt.Printf("📂 ID: %d | ✏️ Название: %s | 🔄 Статус: %s\n", task.ID, task.Title, status)
	fmt.Printf("📝 Описание: %s\n", task.Description)
	fmt.Printf("⏰ Время создания: %s\n", task.CreatedAt.Format("02.01.2006 15:04:05"))
	if task.Completed {
		fmt.Printf("⏰Время завершения: %s\n", task.CompletedAt.Format("02.01.2006 15:04:05"))
	}
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("Задач нет.")
		return
	}

	fmt.Println("\n=== Список задач ===")
	for _, task := range tasks {
		printTaskData(task)
		fmt.Println()
	}

}

func findTaskById(reader *bufio.Reader) {
	if len(tasks) == 0 {
		fmt.Println("Задач нет.")
		return
	}

	var id, err = readID(reader)
	if err != nil {
		fmt.Println("❌ Ошибка ввода:", err)
		return
	}

	task, exists := tasks[id]
	if !exists {
		fmt.Println("Указанной задачи не существует")
		return
	}
	printTaskData(task)
	fmt.Println()

}

func completeTask(reader *bufio.Reader) {
	if len(tasks) == 0 {
		fmt.Println("Задач нет.")
		return
	}

	var id, err = readID(reader)
	if err != nil {
		fmt.Println("❌ Ошибка ввода:", err)
		return
	}

	task, exists := tasks[id]
	if !exists {
		fmt.Println("Указанной задачи не существует")
		return
	}

	if task.Completed {
		fmt.Println("Задача уже выполнена")
		return
	}
	timeNow := time.Now()
	task.Completed = true
	task.CompletedAt = &timeNow

	fmt.Printf("✅ Задача '%s' отмечена как выполненная в %s\n",
		task.Title, task.CompletedAt.Format("02.01.2006 15:04:05"))

}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n=== Менеджер задач ===")
		fmt.Println("1. ➕ Добавить задачу")
		fmt.Println("2. 📋 Список всех задач")
		fmt.Println("3. 📋 Найти задачу по ID")
		fmt.Println("4. ✏️ Отметить задачу как выполненную")
		fmt.Println("5. 🚪 Выход")
		fmt.Print("Выберите действие: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			addTask(reader)
		case "2":
			listTasks()
		case "3":
			findTaskById(reader)
		case "4":
			completeTask(reader)
		case "5":
			return
		default:
			fmt.Println("Нет такого пункта, ебанат")
		}
	}

}
