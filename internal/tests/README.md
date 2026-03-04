# TaskTracker Tests Package

Полный пакет тестов для проекта TaskTracker, обеспечивающий покрытие всех компонентов приложения.

## Структура тестов

### 1. `memory_repository_test.go` - Тесты Repository (Memory)
Тесты для в памяти хранилища задач (`internal/repository/memory/storage.go`).

**Протестированные функции:**
- `NewStorage()` - создание нового хранилища
- `Add()` - добавление задачи с валидацией
- `GetAll()` - получение всех задач
- `GetByID()` - получение задачи по ID с обработкой ошибок
- `Complete()` - отметить задачу как выполненную
- `GetByTag()` - поиск задач по тегу

**Все тесты:**
- `TestMemory_NewStorage` - инициализация хранилища
- `TestMemory_Add_BasicTask` - добавление простой задачи
- `TestMemory_Add_EmptyTitle` - обработка пустого названия
- `TestMemory_Add_WithoutTags` - добавление без тегов
- `TestMemory_Add_MultipleTasks` - добавление нескольких задач
- `TestMemory_GetAll_Empty` - получение из пустого хранилища
- `TestMemory_GetAll_WithTasks` - получение всех задач
- `TestMemory_GetByID_Success` - успешный поиск по ID
- `TestMemory_GetByID_NotFound` - ошибка при поиске несуществующей задачи
- `TestMemory_GetByTag_Empty` - поиск в пустом хранилище
- `TestMemory_GetByTag_Success` - успешный поиск по тегу
- `TestMemory_GetByTag_NotFound` - поиск несуществующего тега
- `TestMemory_Complete_Success` - завершение задачи
- `TestMemory_Complete_AlreadyCompleted` - ошибка при повторном завершении
- `TestMemory_Complete_NotFound` - ошибка при завершении несуществующей задачи
- `TestMemory_Concurrent` - параллельное добавление задач
- `TestMemory_GetByTag_MultipleTags` - поиск с несколькими тегами
- `TestMemory_TaskIntegrity` - целостность данных при операциях
- `TestMemory_LongTitle` - работа с длинными названиями
- `TestMemory_EmptyTagList` - работа с пустым списком тегов

**Количество тестов: 20**

### 2. `json_repository_test.go` - Тесты Repository (JSON)
Тесты для JSON хранилища задач (`internal/repository/sjson/storage.go`).

**Протестированные функции:**
- `NewJSONStorage()` - создание и инициализация хранилища
- `Add()` - добавление с сохранением в файл
- `GetAll()` - получение с загрузкой из файла
- `GetByID()` - получение по ID
- `Complete()` - завершение с сохранением
- `GetByTag()` - поиск по тегам с хранением

**Все тесты:**
- `TestJSONStorage_NewStorage` - создание нового хранилища
- `TestJSONStorage_NewStorage_NonExistentFile` - создание для несуществующего файла
- `TestJSONStorage_Add_BasicTask` - добавление задачи
- `TestJSONStorage_Persistence` - сохранение и загрузка данных
- `TestJSONStorage_GetByID` - получение по ID
- `TestJSONStorage_GetByID_NotFound` - ошибка при поиске
- `TestJSONStorage_Complete` - завершение с сохранением
- `TestJSONStorage_Complete_AlreadyCompleted` - ошибка при повторном завершении
- `TestJSONStorage_GetAll` - получение всех задач
- `TestJSONStorage_GetByTag` - поиск по тегам
- `TestJSONStorage_Concurrent` - параллельные операции
- `TestJSONStorage_DataIntegration` - интеграция с файлом
- `TestJSONStorage_TagPersistence` - сохранение тегов
- `TestJSONStorage_InvalidPath` - обработка неправильного пути
- `TestJSONStorage_CompleteAndVerify` - завершение и верификация времени

**Количество тестов: 15**

### 3. `service_test.go` - Тесты Service
Тесты для бизнес-логики (`internal/service/service.go`).

**Протестированные функции:**
- `NewNewTaskService()` - создание сервиса
- `AddTask()` - добавление с автогенерацией названия
- `GetAllTasks()` - получение всех
- `GetTaskById()` - получение по ID с валидацией
- `GetTasksByTag()` - получение по тегу с проверкой
- `CompleteTask()` - завершение с валидацией

**Все тесты:**
- `TestTaskService_NewService` - создание сервиса
- `TestTaskService_AddTask_WithTitle` - добавление с названием
- `TestTaskService_AddTask_WithoutTitle` - автогенерация названия
- `TestTaskService_AddTask_MultipleTasks` - добавление нескольких
- `TestTaskService_GetAllTasks` - получение всех
- `TestTaskService_GetAllTasks_Empty` - получение из пустого
- `TestTaskService_GetTaskById_Success` - успешное получение
- `TestTaskService_GetTaskById_InvalidID` - ошибка при неправильном ID
- `TestTaskService_GetTaskById_NotFound` - ошибка при поиске
- `TestTaskService_GetTasksByTag` - поиск по тегу
- `TestTaskService_GetTasksByTag_InvalidTag` - ошибка при пустом теге
- `TestTaskService_GetTasksByTag_NotFound` - поиск несуществующего тега
- `TestTaskService_CompleteTask_Success` - успешное завершение
- `TestTaskService_CompleteTask_InvalidID` - ошибка при неправильном ID
- `TestTaskService_CompleteTask_AlreadyCompleted` - ошибка при повторном завершении
- `TestTaskService_CompleteTask_NotFound` - ошибка при поиске
- `TestTaskService_WorkflowScenario` - типичный workflow
- `TestTaskService_AddTask_WithManyTags` - добавление с множеством тегов
- `TestTaskService_AddTask_EmptyDescription` - пустое описание
- `TestTaskService_GetAllTasks_Order` - порядок задач

**Количество тестов: 20**

### 4. `model_test.go` - Тесты Model
Тесты для структур данных (`internal/model/task.go`).

**Протестированные поля:**
- Все поля структуры Task (ID, Title, Description, Completed, CreatedAt, CompletedAt, Tags)

**Все тесты:**
- `TestTask_Creation` - создание структуры
- `TestTask_Completed` - завершенная задача
- `TestTask_EmptyTags` - пустой список тегов
- `TestTask_NilTags` - nil теги
- `TestTask_MultipleTags` - несколько тегов
- `TestTask_Fields` - изменение полей
- `TestTask_TimeComparison` - сравнение времени
- `TestTask_Pointer` - работа с указателями
- `TestTask_TagModification` - изменение тегов
- `TestTask_ZeroValues` - нулевые значения
- `TestTask_TagOrder` - порядок тегов
- `TestTask_Copy` - копирование структуры
- `TestTask_LongDescription` - длинное описание
- `TestTask_SpecialCharacters` - специальные символы
- `TestTask_CompletedAtNil` - nil CompletedAt
- `TestTask_CompletedAtDeref` - разыменование CompletedAt

**Количество тестов: 16**

### 5. `errors_test.go` - Тесты Errors Package
Тесты для пакета ошибок (`errors/errors.go`).

**Протестированные ошибки:**
- `ErrCantCreateJsonFile` - ошибка создания JSON
- `ErrCantReadJsonData` - ошибка чтения JSON
- `ErrWrongPath` - неправильный путь
- `ErrCantSaveTaskToJson` - ошибка сохранения
- `ErrTasksNotFound` - задачи не найдены
- `ErrIdNotExists` - ID не существует
- `ErrWrongTypeID` - неправильный тип ID
- `ErrTasksWithTagNotFound` - задач с тегом не найдено
- `ErrWrongTag` - неправильный тег
- `ErrTaskAlredyComplete` - задача уже завершена

**Все тесты:**
- `TestErrorVariables` - определение всех ошибок
- `TestErrorComparison` - сравнение ошибок
- `TestErrCantCreateJsonFile` - тест ErrCantCreateJsonFile
- `TestErrCantReadJsonData` - тест ErrCantReadJsonData
- `TestErrWrongPath` - тест ErrWrongPath
- `TestErrCantSaveTaskToJson` - тест ErrCantSaveTaskToJson
- `TestErrTasksNotFound` - тест ErrTasksNotFound
- `TestErrIdNotExists` - тест ErrIdNotExists
- `TestErrWrongTypeID` - тест ErrWrongTypeID
- `TestErrTasksWithTagNotFound` - тест ErrTasksWithTagNotFound
- `TestErrWrongTag` - тест ErrWrongTag
- `TestErrTaskAlreadyComplete` - тест ErrTaskAlredyComplete
- `TestErrorMessages` - сообщения об ошибках
- `TestErrorTypeAssertion` - типизация ошибок
- `TestErrorInComparison` - использование в сравнении
- `TestErrorsAreStable` - стабильность ошибок
- `TestErrorStrings` - строковое представление
- `TestErrorUsagePattern` - типичный паттерн использования
- `TestAllErrorsDefined` - определение всех ошибок

**Количество тестов: 19**

### 6. `handler_test.go` - Тесты Handler
Тесты для CLI handler (`internal/handler/cli_handler.go`).

**Протестированные функции:**
- `NewCLIHandler()` - создание handler
- Парсинг тегов
- Чтение ввода
- Интеграция с сервисом

**Все тесты:**
- `TestCLIHandler_Creation` - создание handler
- `TestCLIHandler_ParseTags_SingleTag` - парсинг одного тега
- `TestCLIHandler_ParseTags_CommaSeparated` - запятые
- `TestCLIHandler_ParseTags_SpaceSeparated` - пробелы
- `TestCLIHandler_ParseTags_SemicolonSeparated` - точка с запятой
- `TestCLIHandler_ParseTags_PipeSeparated` - pipe символ
- `TestCLIHandler_ParseTags_MixedSeparators` - смешанные разделители
- `TestCLIHandler_ParseTags_Empty` - пустая строка
- `TestCLIHandler_ParseTags_Whitespace` - только пробелы
- `TestCLIHandler_ReadInput_SimpleInput` - чтение ввода
- `TestCLIHandler_Service_Integration` - интеграция с сервисом
- `TestCLIHandler_ErrorHandling` - обработка ошибок
- `TestCLIHandler_MultipleOperations` - несколько операций
- `TestCLIHandler_StringValidation` - валидация строк
- `TestCLIHandler_ConcurrentOperations` - параллельные операции
- `TestCLIHandler_TagProcessing` - обработка тегов
- `TestCLIHandler_TaskCompletion` - завершение задач
- `TestCLIHandler_WorkflowWithHandler` - полный workflow
- `TestCLIHandler_RepositoryIntegration` - интеграция с repository
- `TestCLIHandler_Stability` - стабильность

**Количество тестов: 20**

### 7. `integration_test.go` - Интеграционные тесты
Тесты для интеграции всех компонентов вместе.

**Все тесты:**
- `TestIntegration_MemoryRepository_Service_Handler` - полная интеграция с Memory
- `TestIntegration_JSONRepository_Service_Handler` - полная интеграция с JSON
- `TestIntegration_MultipleServices` - несколько сервисов с одним хранилищем
- `TestIntegration_ComplexWorkflow` - сложный workflow
- `TestIntegration_JSONPersistence` - сохранение данных в JSON
- `TestIntegration_ConcurrentOperations` - параллельные операции через интеграцию
- `TestIntegration_ServiceWithDifferentStorages` - сервис с разными хранилищами
- `TestIntegration_FullApplicationSimulation` - полное приложение
- `TestIntegration_LargeDataSet` - работа с большим количеством задач

**Количество тестов: 9**

## Статистика

| Компонент       | Количество тестов |
|-----------------|-------------------|
| Memory Repo     | 20                |
| JSON Repo       | 15                |
| Service         | 20                |
| Model           | 16                |
| Errors          | 19                |
| Handler         | 20                |
| Integration     | 9                 |
| **ИТОГО**       | **119**           |

## Запуск тестов

### Все тесты
```bash
go test ./internal/tests -v
```

### Конкретный файл
```bash
go test ./internal/tests -v -run TestMemory
```

### Конкретный тест
```bash
go test ./internal/tests -v -run TestMemory_Add_BasicTask
```

### С покрытием
```bash
go test ./internal/tests -cover
go test ./internal/tests -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Параллельный запуск
```bash
go test ./internal/tests -v -race
```

### Бенчмарки (если добавить)
```bash
go test ./internal/tests -bench=.
```

## Особенности тестов

### Параллелизм
- Большинство тестов запускаются параллельно с флагом `t.Parallel()`
- Безопасны для параллельного выполнения

### Таблично-ориентированные тесты
- Используются для тестирования вариантов
- Облегчают добавление новых тестовых случаев

### Изоляция
- Каждый тест использует своё хранилище
- JSON тесты используют временные файлы
- Нет зависимостей между тестами

### Функциональное покрытие
- Все публичные функции протестированы
- Граничные случаи (пустые данные, несуществующие элементы)
- Ошибки и исключительные ситуации
- Параллельные операции

## Как добавить новый тест

1. Выберите подходящий файл или создайте новый
2. Используйте паттерн: `TestPackage_Function_Case`
3. Добавьте `t.Parallel()` если возможно
4. Используйте descriptive ошибки
5. Проверьте что тест работает: `go test ./internal/tests -v`

## Примеры использования

### Добавление задачи
```go
repo := memory.NewStorage()
svc := service.NewNewTaskService(repo)
task, err := svc.AddTask("Title", "Description", []string{"tag"})
```

### Получение всех задач
```go
tasks, err := svc.GetAllTasks()
```

### Поиск по ID
```go
task, err := svc.GetTaskById(1)
```

### Поиск по тегу
```go
tasks, err := svc.GetTasksByTag("work")
```

### Завершение задачи
```go
completed, err := svc.CompleteTask(1)
```

## Для разработчиков

- Используйте `testing.T` для отчёта об ошибках
- Используйте `t.Fatalf()` для критических ошибок
- Используйте `t.Errorf()` для тестовых ошибок
- Используйте `t.Run()` для подтестов
- Добавляйте `t.Parallel()` для независимых тестов

## Заметки

Все тесты разработаны чтобы быть:
- **Независимыми** - не зависят друг от друга
- **Воспроизводимыми** - дают одинаковые результаты
- **Быстрыми** - выполняются за приемлемое время
- **Читаемыми** - ясно выражают назначение
- **Полными** - покрывают основной функционал
