package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"TaskTracker/internal/model"
	"TaskTracker/internal/repository"
	"TaskTracker/internal/repository/memory"
)

func readID(reader *bufio.Reader) (int, error) {
	fmt.Print("\nВведите ID задачи: ")
	strId, _ := reader.ReadString('\n')
	strId = strings.TrimSpace(strId)
	id, err := strconv.Atoi(strings.TrimSpace(strId))
	if err != nil || id < 1 {
		return 0, fmt.Errorf("некорректный ID")
	}
	return id, nil
}

func printTaskData(task *model.Task) {
	status := "❌ - Не выполнена"
	if task.Completed {
		status = "✅ - Выполнена"
	}
	fmt.Printf("📂 ID: %d | ✏️ Название: %s | 🔄 Статус: %s\n", task.ID, task.Title, status)
	fmt.Printf("📝 Описание: %s\n🏷️ Тэги: %s\n", task.Description, task.Tags)
	fmt.Printf("⏰ Время создания: %s\n", task.CreatedAt.Format("02.01.2006 15:04:05"))
	if task.Completed {
		fmt.Printf("⏰Время завершения: %s\n", task.CompletedAt.Format("02.01.2006 15:04:05"))
	}
}

func addTask(reader *bufio.Reader, repo repository.Repository) {
	fmt.Println("=======================")
	fmt.Print("Введите название задачи: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Введите описание вашей задачи: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Введите тег/и: ")
	strTags, _ := reader.ReadString('\n')
	replacer := strings.NewReplacer(",", " ", ";", " ", "|", " ")
	normalized := replacer.Replace(strTags)
	sliceTags := strings.Fields(normalized)

	task := repo.Add(title, description, sliceTags)
	fmt.Printf("\n✅ Задача: '%s' добавлена успешно.\n", task.Title)
}

func listTasks(repo repository.Repository) {
	if repo.IsEmpty() {
		fmt.Println("Задач нет")
		return
	}

	tasks, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, task := range tasks {
		printTaskData(task)
	}

}

func findTaskById(reader *bufio.Reader, repo repository.Repository) {
	var id, err = readID(reader)
	if err != nil {
		fmt.Println("❌ Ошибка ввода:", err)
		return
	}

	task, err := repo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	printTaskData(task)

}

func findTaskByTag(reader *bufio.Reader, repo repository.Repository) {
	fmt.Print("\nВведите Тэг задачи: ")
	inputTag, _ := reader.ReadString('\n')
	inputTag = strings.TrimSpace(inputTag)
	words := strings.Split(inputTag, " ")
	if len(words) > 1 {
		inputTag = words[0]
	}

	tasks, err := repo.GetByTag(inputTag)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, task := range tasks {
		printTaskData(task)
	}
}

func completeTask(reader *bufio.Reader, repo repository.Repository) {
	var id, err = readID(reader)
	if err != nil {
		fmt.Println("❌ Ошибка ввода:", err)
		return
	}

	isCompleteErr := repo.Complete(id)

	if isCompleteErr != nil {
		fmt.Println(isCompleteErr)
		return
	}

	task, err := repo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("✅ Задача '%s' отмечена как выполненная в %s\n",
		task.Title, task.CompletedAt.Format("02.01.2006 15:04:05"))

}

func main() {
	var repo repository.Repository
	repo = memory.NewStorage()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n=== Менеджер задач ===")
		fmt.Println("1. ➕ Добавить задачу")
		fmt.Println("2. 📋 Список всех задач")
		fmt.Println("3. 📋 Найти задачу по ID")
		fmt.Println("4. 🏷️ Найти задачу по Тэгу")
		fmt.Println("5. ✏️ Отметить задачу как выполненную")
		fmt.Println("6. 🚪 Выход")
		fmt.Print("\nВыберите действие: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			addTask(reader, repo)
		case "2":
			listTasks(repo)
		case "3":
			findTaskById(reader, repo)
		case "4":
			findTaskByTag(reader, repo)
		case "5":
			completeTask(reader, repo)
		case "6":
			return
		default:
			fmt.Println("Нет такого пункта, ебанат")
		}
	}

}
