package errors

import "errors"

var (
	// JSON file errors
	ErrCantCreateJsonFile = errors.New("Не удалось создать файл JSON")
	ErrCantReadJsonData   = errors.New("Не удалось прочесть данные из файла JSON")
	ErrWrongPath          = errors.New("Некорректный путь к созданию / перезаписи файла JSON")
	ErrCantSaveTaskToJson = errors.New("Не удалось сохранить task в JSON")

	// Task-related errors
	ErrTasksNotFound        = errors.New("Нет задач для отображения")
	ErrIdNotExists          = errors.New("Несуществующий ID")
	ErrWrongTypeID          = errors.New("Неккоректный ID")
	ErrTasksWithTagNotFound = errors.New("Задач с таким тегом не существует")
	ErrWrongTag             = errors.New("Неккоректный тег")
	ErrTaskAlredyComplete   = errors.New("Задача уже выполнена")
	ErrTaskNotFound         = errors.New("Задача не найдена")
	ErrCantDeleteTask       = errors.New("Невозможно удалить задачу. Возможно, она не существует")

	// Database errors
	ErrCantConnectToDB   = errors.New("Невозможно подключиться к Базе Данных.")
	ErrPingEx            = errors.New("Провалена проверка подключения к БД.")
	ErrCantCreateDBTable = errors.New("Невозможно создать такую таблицу")
	ErrCantSaveTaskToDB  = errors.New("Невозможно сохранить задачу в БД")
	ErrCantReadTable     = errors.New("Невозможно прочитать данные из таблицы")
	ErrTableIsEmpty      = errors.New("Строки не найдены")

	// Repository errors
	ErrWrongTypeRepo = errors.New("Неверный тип репозитория")

	// Stats errors
	ErrStatsDoesntWritten = errors.New("Статистика еще не записана")

	// Migration errors
	ErrCreateMigration  = errors.New("Ошибка создания мигратора")
	ErrCantUseMigration = errors.New("Ошибка применения миграций")

	// Auth errors
	ErrUserNotFound   = errors.New("Пользователь не найден")
	ErrCantCreateUser = errors.New("Невозможно создать пользователя")
)

// GetUserFriendlyMessage returns a user-friendly error message for display in Telegram
func GetUserFriendlyMessage(err error) string {
	if err == nil {
		return ""
	}

	// Check for specific errors and return friendly messages
	switch {
	case errors.Is(err, ErrTasksNotFound), errors.Is(err, ErrTableIsEmpty):
		return "📭 Задач не найдено"
	case errors.Is(err, ErrIdNotExists), errors.Is(err, ErrWrongTypeID):
		return "❌ Неверный ID задачи"
	case errors.Is(err, ErrTasksWithTagNotFound), errors.Is(err, ErrWrongTag):
		return "📭 Задач с таким тегом не найдено"
	case errors.Is(err, ErrTaskAlredyComplete):
		return "✅ Задача уже выполнена"
	case errors.Is(err, ErrTaskNotFound):
		return "❌ Задача не найдена или уже удалена"
	case errors.Is(err, ErrCantDeleteTask):
		return "❌ Не удалось удалить задачу"
	case errors.Is(err, ErrCantConnectToDB), errors.Is(err, ErrPingEx):
		return "⚠️ Ой, что-то пошло не так... Попробуйте позже"
	case errors.Is(err, ErrCantSaveTaskToDB), errors.Is(err, ErrCantCreateDBTable):
		return "⚠️ Не удалось сохранить задачу. Попробуйте позже"
	case errors.Is(err, ErrCantReadTable):
		return "⚠️ Не удалось получить данные. Попробуйте позже"
	case errors.Is(err, ErrUserNotFound):
		return "❌ Пользователь не найден"
	case errors.Is(err, ErrCantCreateUser):
		return "⚠️ Не удалось создать профиль. Попробуйте позже"
	case errors.Is(err, ErrStatsDoesntWritten):
		return "📭 Статистика еще не доступна"
	case errors.Is(err, ErrCantCreateJsonFile), errors.Is(err, ErrCantReadJsonData), errors.Is(err, ErrWrongPath):
		return "⚠️ Ошибка работы с файлами. Попробуйте позже"
	case errors.Is(err, ErrCantSaveTaskToJson):
		return "⚠️ Не удалось сохранить задачу. Попробуйте позже"
	}

	// If it's a wrapped error, try to find the underlying error
	if cause := errors.Unwrap(err); cause != nil {
		return GetUserFriendlyMessage(cause)
	}

	// Generic fallback message
	return "⚠️ Ой, что-то пошло не так... Попробуйте позже"
}
