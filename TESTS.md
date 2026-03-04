# TaskTracker - Testing Documentation

## Overview
This document describes all the tests created for the TaskTracker project. The test suite includes unit tests for individual components, integration tests for workflows, and edge case handling.

## Test Files Structure

### 1. **internal/model/task_test.go** - Task Model Tests
Tests for the `Task` struct and its properties.

#### Test Cases:
- **TestTaskCreation**: Tests creating Task objects with various properties
- **TestTaskWithTags**: Tests Task with multiple tags
- **TestTaskCompletedAtTime**: Tests completed_at timestamp handling
- **TestTaskNotCompleted**: Tests incomplete task creation

**Coverage**: Task struct initialization, field assignments, tag management

---

### 2. **internal/repository/memory/storage_test.go** - Memory Storage Tests
Tests for the in-memory task storage implementation.

#### Test Cases:
- **TestMemoryStorageAdd**: Tests adding tasks to memory storage
  - Simple task addition
  - Tasks without tags
  - Multiple tags per task

- **TestMemoryStorageGetAll**: Tests retrieving all tasks
  - Empty storage handling
  - Multiple task retrieval

- **TestMemoryStorageGetByID**: Tests finding tasks by ID
  - Valid ID lookup
  - Non-existing ID handling
  - Invalid ID ranges

- **TestMemoryStorageGetByTag**: Tests filtering tasks by tags
  - Single tag filtering
  - Multi-tag scenarios
  - Non-existing tags

- **TestMemoryStorageComplete**: Tests task completion
  - Marking task as complete
  - Preventing double completion
  - Non-existing task completion attempts

- **TestMemoryStorageIDIncrement**: Tests ID auto-increment
  - Sequential ID assignment

- **TestMemoryStorageConcurrency**: Tests concurrent task additions
  - Multi-goroutine safety

**Coverage**: Full CRUD operations, error handling, ID management

---

### 3. **internal/repository/sjson/storage_test.go** - JSON Storage Tests
Tests for the persistent JSON file storage implementation.

#### Test Cases:
- **TestJSONStorageCreate**: Tests file creation and initialization
- **TestJSONStorageAddAndLoad**: Tests data persistence across sessions
- **TestJSONStorageGetAll**: Tests retrieving all persisted tasks
- **TestJSONStorageGetByID**: Tests finding persisted tasks by ID
- **TestJSONStorageGetByTag**: Tests tag-based filtering with persistence
- **TestJSONStorageComplete**: Tests completing and persisting completed tasks
- **TestJSONStoragePersistence**: Tests full read-write-read cycle
- **TestJSONStorageIDPersistence**: Tests ID counter persistence
- **TestJSONStorageIDIncrement**: Tests sequential ID generation
- **TestJSONStorageCompleteDateTime**: Tests timestamps on completion
- **TestJSONStorageEmptyTags**: Tests handling of tasks without tags
- **TestJSONStorageFileNotFound**: Tests error handling for missing paths
- **TestJSONStorageConcurrentWrites**: Tests concurrent file operations

**Coverage**: Persistence layer, file I/O, data serialization, concurrency

---

### 4. **internal/service/service_test.go** - Service Layer Tests
Tests for the business logic service layer.

#### Test Cases:
- **TestAddTaskWithTitle**: Tests task creation via service
  - With explicit title
  - With default title generation
  - Multiple tags

- **TestGetAllTasks**: Tests retrieving all tasks via service
  - Empty list handling
  - Multiple task retrieval

- **TestGetTaskById**: Tests task lookup by ID
  - Valid IDs
  - Invalid/negative IDs
  - Non-existing tasks

- **TestGetTasksByTag**: Tests tag-based filtering
  - Valid tag searches
  - Empty tag validation
  - Non-existing tags

- **TestCompleteTask**: Tests task completion service
  - Single completion
  - Double completion prevention
  - Invalid ID handling

- **TestServiceWithMultipleTasks**: Tests service with multiple tasks
  - Complex workflows
  - State management

- **TestServiceWithMockRepository**: Tests service with mock repository
  - Dependency injection
  - Method invocation verification

**Coverage**: Service methods, validation rules, business logic

---

### 5. **internal/handler/cli_handler_test.go** - CLI Handler Tests
Tests for the command-line interface handler.

#### Test Cases:
- **TestParseTagsWithComma**: Tests tag parsing with various delimiters
  - Comma-separated tags
  - Space-separated tags
  - Semicolon-separated tags
  - Pipe-separated tags
  - Mixed delimiters

- **TestReadInput**: Tests user input reading and trimming
  - Simple input
  - Input with leading/trailing spaces
  - Empty input

- **TestReadID**: Tests ID input parsing
  - Valid positive IDs
  - Large IDs
  - Zero and negative values

- **TestHandleChoice**: Tests menu choice handling
  - All menu options (1-6)
  - Exit condition
  - Invalid choices

- **TestCLIHandlerCreation**: Tests handler initialization
- **TestParseTagsEdgeCases**: Tests edge cases in tag parsing
  - Only separators
  - Single character tags
  - Tags with numbers

- **TestPrintTask**: Tests task display formatting
- **TestPrintMenu**: Tests menu display rendering
- **TestTaskPrintingWithNilCompletedAt**: Tests display of incomplete tasks
- **TestTaskPrintingWithCompletedAt**: Tests display of completed tasks

**Coverage**: User input handling, CLI logic, display formatting

---

### 6. **internal/tests/integration_test.go** - Integration Tests
Tests for complete workflows combining multiple components.

#### Test Cases:
- **TestFullWorkflowWithMemoryStorage**: Complete workflow test
  - Add multiple tasks
  - Query tasks
  - Filter by tags
  - Mark tasks complete

- **TestFullWorkflowWithJSONStorage**: Complete workflow with persistence
  - Add tasks
  - Save to file
  - Load from file
  - Verify data integrity

- **TestMultipleTasksWithDifferentStates**: Tests multiple task states
  - Various task states
  - State transitions
  - Tag associations

- **TestErrorHandling**: Integration error scenarios
  - Invalid IDs
  - Empty tags
  - Non-existing items

- **TestDefaultTitleGeneration**: Tests automatic title generation
  - Generated vs explicit titles
  - Title uniqueness

- **TestTaskTagManagement**: Tests comprehensive tag operations
  - Multi-tag assignments
  - Tag-based queries

- **TestIDValidation**: Tests ID validation across operations
- **TestTaskCompletionTime**: Tests timestamp accuracy
- **TestTaskPersistenceWithIDIncrement**: Tests ID continuity across sessions
- **TestEmptyOperations**: Tests operations on empty storage

**Coverage**: End-to-end workflows, data persistence, state management

---

## Running the Tests

### Run all tests:
```bash
go test ./...
```

### Run specific test file:
```bash
go test ./internal/service -v
```

### Run specific test:
```bash
go test -run TestAddTaskWithTitle ./internal/service -v
```

### Run tests with coverage:
```bash
go test ./... -cover
```

### Generate coverage report:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Test Coverage Summary

| Component | Coverage | Notes |
|-----------|----------|-------|
| Model | Comprehensive | Task struct and properties |
| Memory Storage | Comprehensive | All CRUD operations |
| JSON Storage | Comprehensive | Persistence and file I/O |
| Service | Comprehensive | Business logic and validation |
| Handler | Comprehensive | CLI operations and formatting |
| Integration | Comprehensive | Full workflows |

---

## Key Testing Scenarios

### 1. **Data Persistence**
- Tasks saved to JSON are correctly loaded
- ID counters persist across sessions
- Completed state persists

### 2. **Error Handling**
- Invalid IDs are rejected
- Empty tags are validated
- Non-existing tasks return errors
- Already completed tasks cannot be completed again

### 3. **Concurrency**
- Multiple concurrent writes are handled safely
- Lock mechanisms work correctly
- Race conditions are prevented

### 4. **Data Integrity**
- Tags are correctly associated with tasks
- Task timestamps are accurate
- Completed times are set correctly
- All fields persist correctly

### 5. **State Management**
- Tasks transition between states correctly
- Completed state is immutable
- IDs are unique and sequential

---

## Testing Best Practices Used

1. **Table-driven tests**: Using test case tables for multiple scenarios
2. **Cleanup**: Removing temporary files after JSON storage tests
3. **Isolation**: Each test is independent and doesn't affect others
4. **Mock objects**: Using mock repositories for service tests
5. **Edge cases**: Testing boundary conditions and error paths
6. **Integration testing**: Testing components working together
7. **Concurrent testing**: Verifying thread-safety

---

## Error Scenarios Covered

- Invalid task IDs (negative, zero, non-existing)
- Empty or invalid tags
- Attempting to complete already completed tasks
- File I/O errors
- Concurrent access scenarios
- Empty storage operations
- Invalid user input

---

## Notes for Developers

- All tests use temporary files for JSON storage to avoid side effects
- Memory storage tests are fast and can be run frequently
- Integration tests verify the system works as a whole
- Tests follow Go testing conventions (Test prefix, descriptive names)
- Each test is self-contained and can run independently
- Error messages are clear and helpful for debugging

---

## Future Test Enhancements

1. Add stress tests with large datasets
2. Add performance benchmarks
3. Add fuzz testing for input validation
4. Add mutation testing to verify test quality
5. Add property-based testing with frameworks like go-quickcheck
