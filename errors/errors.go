package errors

import "errors"

var (
	ErrCantCreateJsonFile   = errors.New("Не удалось создать файл JSON")
	ErrCantReadJsonData     = errors.New("Не удалось прочесть данные из файла JSON")
	ErrWrongPath            = errors.New("Некорректный путь к созданию / перезаписи файла JSON")
	ErrCantSaveTaskToJson   = errors.New("Не удалось сохранить ебаную таску в JSON")
	ErrTasksNotFound        = errors.New("Нема задач вообще")
	ErrIdNotExists          = errors.New("Несуществующий ID")
	ErrWrongTypeID          = errors.New("Неккоректный ID")
	ErrTasksWithTagNotFound = errors.New("Задач с таким тегом не существует")
	ErrWrongTag             = errors.New("Неккоректный тег")
	ErrTaskAlredyComplete   = errors.New("Задача уже выполнена")
)
