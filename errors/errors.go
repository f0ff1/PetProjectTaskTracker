package errors

import "errors"

var (
	ErrCantCreateJsonFile   = errors.New("Не удалось создать файл JSON")
	ErrCantReadJsonData     = errors.New("Не удалось прочесть данные из файла JSON")
	ErrWrongPath            = errors.New("Некорректный путь к созданию / перезаписи файла JSON")
	ErrCantSaveTaskToJson   = errors.New("Не удалось сохранить task в JSON")
	ErrTasksNotFound        = errors.New("Нет задач для отображения")
	ErrIdNotExists          = errors.New("Несуществующий ID")
	ErrWrongTypeID          = errors.New("Неккоректный ID")
	ErrTasksWithTagNotFound = errors.New("Задач с таким тегом не существует")
	ErrWrongTag             = errors.New("Неккоректный тег")
	ErrTaskAlredyComplete   = errors.New("Задача уже выполнена")
	ErrCantConnectToDB      = errors.New("Невозможно подключиться к Базе Данных.")
	ErrPingEx               = errors.New("Провалена проверка подключения к БД.")
	ErrCantCreateDBTable    = errors.New("Невозможно создать такую таблицу")
	ErrCantSaveTaskToDB     = errors.New("Невозможно сохрнаить задачу в БД")
	ErrCantReadTable        = errors.New("Невозможно прочитать данные из таблицы")
	ErrTableIsEmpty         = errors.New("Строки не найдены")
	ErrCantDeleteTask       = errors.New("Невозможно удалить задачу. Возможно, она не существует")
	ErrWrongTypeRepo        = errors.New("Неверный тип репозитория")
)
