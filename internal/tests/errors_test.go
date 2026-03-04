package tests

import (
	"errors"
	"testing"

	customError "TaskTracker/errors"
)

// TestErrorVariables проверяет определение переменных ошибок
func TestErrorVariables(t *testing.T) {
	t.Parallel()

	// Проверяем что ошибки определены
	if customError.ErrCantCreateJsonFile == nil {
		t.Error("ErrCantCreateJsonFile should not be nil")
	}

	if customError.ErrCantReadJsonData == nil {
		t.Error("ErrCantReadJsonData should not be nil")
	}

	if customError.ErrWrongPath == nil {
		t.Error("ErrWrongPath should not be nil")
	}

	if customError.ErrCantSaveTaskToJson == nil {
		t.Error("ErrCantSaveTaskToJson should not be nil")
	}

	if customError.ErrTasksNotFound == nil {
		t.Error("ErrTasksNotFound should not be nil")
	}

	if customError.ErrIdNotExists == nil {
		t.Error("ErrIdNotExists should not be nil")
	}

	if customError.ErrWrongTypeID == nil {
		t.Error("ErrWrongTypeID should not be nil")
	}

	if customError.ErrTasksWithTagNotFound == nil {
		t.Error("ErrTasksWithTagNotFound should not be nil")
	}

	if customError.ErrWrongTag == nil {
		t.Error("ErrWrongTag should not be nil")
	}

	if customError.ErrTaskAlredyComplete == nil {
		t.Error("ErrTaskAlredyComplete should not be nil")
	}
}

// TestErrorComparison проверяет сравнение ошибок
func TestErrorComparison(t *testing.T) {
	t.Parallel()

	// Проверяем что ошибки уникальны
	if customError.ErrIdNotExists == customError.ErrWrongTypeID {
		t.Error("Different errors should not be equal")
	}

	if customError.ErrTasksNotFound == customError.ErrIdNotExists {
		t.Error("Different errors should not be equal")
	}

	if customError.ErrCantCreateJsonFile == customError.ErrCantReadJsonData {
		t.Error("Different errors should not be equal")
	}
}

// TestErrCantCreateJsonFile проверяет ошибку создания JSON файла
func TestErrCantCreateJsonFile(t *testing.T) {
	t.Parallel()

	err := customError.ErrCantCreateJsonFile

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	// Проверяем что это обычная ошибка
	if !errors.Is(err, customError.ErrCantCreateJsonFile) {
		t.Error("Error comparison failed")
	}
}

// TestErrCantReadJsonData проверяет ошибку чтения JSON файла
func TestErrCantReadJsonData(t *testing.T) {
	t.Parallel()

	err := customError.ErrCantReadJsonData

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrCantReadJsonData) {
		t.Error("Error comparison failed")
	}
}

// TestErrWrongPath проверяет ошибку неправильного пути
func TestErrWrongPath(t *testing.T) {
	t.Parallel()

	err := customError.ErrWrongPath

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrWrongPath) {
		t.Error("Error comparison failed")
	}
}

// TestErrCantSaveTaskToJson проверяет ошибку сохранения в JSON
func TestErrCantSaveTaskToJson(t *testing.T) {
	t.Parallel()

	err := customError.ErrCantSaveTaskToJson

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrCantSaveTaskToJson) {
		t.Error("Error comparison failed")
	}
}

// TestErrTasksNotFound проверяет ошибку когда задач не найдено
func TestErrTasksNotFound(t *testing.T) {
	t.Parallel()

	err := customError.ErrTasksNotFound

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrTasksNotFound) {
		t.Error("Error comparison failed")
	}
}

// TestErrIdNotExists проверяет ошибку несуществующего ID
func TestErrIdNotExists(t *testing.T) {
	t.Parallel()

	err := customError.ErrIdNotExists

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrIdNotExists) {
		t.Error("Error comparison failed")
	}
}

// TestErrWrongTypeID проверяет ошибку неправильного типа ID
func TestErrWrongTypeID(t *testing.T) {
	t.Parallel()

	err := customError.ErrWrongTypeID

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrWrongTypeID) {
		t.Error("Error comparison failed")
	}
}

// TestErrTasksWithTagNotFound проверяет ошибку когда задач с тегом не найдено
func TestErrTasksWithTagNotFound(t *testing.T) {
	t.Parallel()

	err := customError.ErrTasksWithTagNotFound

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrTasksWithTagNotFound) {
		t.Error("Error comparison failed")
	}
}

// TestErrWrongTag проверяет ошибку неправильного тега
func TestErrWrongTag(t *testing.T) {
	t.Parallel()

	err := customError.ErrWrongTag

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrWrongTag) {
		t.Error("Error comparison failed")
	}
}

// TestErrTaskAlreadyComplete проверяет ошибку уже завершенной задачи
func TestErrTaskAlreadyComplete(t *testing.T) {
	t.Parallel()

	err := customError.ErrTaskAlredyComplete

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if !errors.Is(err, customError.ErrTaskAlredyComplete) {
		t.Error("Error comparison failed")
	}
}

// TestErrorMessages проверяет сообщения об ошибках
func TestErrorMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		err   error
		empty bool
	}{
		{"ErrCantCreateJsonFile", customError.ErrCantCreateJsonFile, false},
		{"ErrCantReadJsonData", customError.ErrCantReadJsonData, false},
		{"ErrWrongPath", customError.ErrWrongPath, false},
		{"ErrCantSaveTaskToJson", customError.ErrCantSaveTaskToJson, false},
		{"ErrTasksNotFound", customError.ErrTasksNotFound, false},
		{"ErrIdNotExists", customError.ErrIdNotExists, false},
		{"ErrWrongTypeID", customError.ErrWrongTypeID, false},
		{"ErrTasksWithTagNotFound", customError.ErrTasksWithTagNotFound, false},
		{"ErrWrongTag", customError.ErrWrongTag, false},
		{"ErrTaskAlredyComplete", customError.ErrTaskAlredyComplete, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.empty {
				if tt.err.Error() != "" {
					t.Errorf("Expected empty message, got %q", tt.err.Error())
				}
			} else {
				if tt.err.Error() == "" {
					t.Errorf("Expected non-empty message")
				}
			}
		})
	}
}

// TestErrorTypeAssertion проверяет типизацию ошибок
func TestErrorTypeAssertion(t *testing.T) {
	t.Parallel()

	var err error

	// Все ошибки должны реализовывать интерфейс error
	err = customError.ErrIdNotExists
	if err == nil {
		t.Error("Error should not be nil after assignment")
	}

	err = customError.ErrWrongTypeID
	if err == nil {
		t.Error("Error should not be nil after assignment")
	}

	err = customError.ErrTaskAlredyComplete
	if err == nil {
		t.Error("Error should not be nil after assignment")
	}
}

// TestErrorInComparison проверяет использование ошибок в сравнении
func TestErrorInComparison(t *testing.T) {
	t.Parallel()

	var err1 error = customError.ErrIdNotExists
	var err2 error = customError.ErrIdNotExists

	if err1 != err2 {
		t.Error("Same errors should be equal")
	}

	err1 = customError.ErrIdNotExists
	err2 = customError.ErrWrongTypeID

	if err1 == err2 {
		t.Error("Different errors should not be equal")
	}
}

// TestErrorsAreConstant проверяет что ошибки примерно константны
func TestErrorsAreStable(t *testing.T) {
	t.Parallel()

	err1 := customError.ErrIdNotExists
	err2 := customError.ErrIdNotExists

	if err1 != err2 {
		t.Error("Same error variable should always be equal")
	}

	if err1.Error() != err2.Error() {
		t.Error("Error messages should be consistent")
	}
}

// TestErrorStrings проверяет строковое представление ошибок
func TestErrorStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err error
	}{
		{customError.ErrCantCreateJsonFile},
		{customError.ErrCantReadJsonData},
		{customError.ErrWrongPath},
		{customError.ErrCantSaveTaskToJson},
		{customError.ErrTasksNotFound},
		{customError.ErrIdNotExists},
		{customError.ErrWrongTypeID},
		{customError.ErrTasksWithTagNotFound},
		{customError.ErrWrongTag},
		{customError.ErrTaskAlredyComplete},
	}

	for _, tt := range tests {
		// Проверяем что String() работает
		str := tt.err.Error()
		if str == "" {
			t.Error("Error.Error() should return non-empty string")
		}

		// Проверяем что возвращается русский текст (закомментировано, т.к. зависит от реализации)
		// Просто проверяем что функция работает
		if len(str) == 0 {
			t.Error("Error string is empty")
		}
	}
}

// TestErrorUsagePattern проверяет типичный паттерн использования ошибок
func TestErrorUsagePattern(t *testing.T) {
	t.Parallel()

	// Типичный паттерн: возвращаем ошибку и проверяем ее
	testFunc := func() error {
		return customError.ErrIdNotExists
	}

	err := testFunc()

	if err != customError.ErrIdNotExists {
		t.Error("Error not properly returned or compared")
	}

	if errors.Is(err, customError.ErrIdNotExists) {
		// Success
	} else {
		t.Error("Error.Is() should match")
	}
}

// TestAllErrorsDefined проверяет что все ошибки определены
func TestAllErrorsDefined(t *testing.T) {
	t.Parallel()

	errors := []error{
		customError.ErrCantCreateJsonFile,
		customError.ErrCantReadJsonData,
		customError.ErrWrongPath,
		customError.ErrCantSaveTaskToJson,
		customError.ErrTasksNotFound,
		customError.ErrIdNotExists,
		customError.ErrWrongTypeID,
		customError.ErrTasksWithTagNotFound,
		customError.ErrWrongTag,
		customError.ErrTaskAlredyComplete,
	}

	for i, err := range errors {
		if err == nil {
			t.Errorf("Error at index %d is nil", i)
		}
		if err.Error() == "" {
			t.Errorf("Error at index %d has empty message", i)
		}
	}
}
