package model

// Основная структура статистики
type TaskStats struct {
	// Основная статистика
	TotalTasks     int     `json:"total_tasks"`
	CompletedTasks int     `json:"completed_tasks"`
	PendingTasks   int     `json:"pending_tasks"`
	CompletionRate float64 `json:"completion_rate"` // процент выполнения

	// Статистика по тегам
	TopTags         []TagStat `json:"top_tags"`
	TotalUniqueTags int       `json:"total_unique_tags"`

	// Временная статистика
	TasksByDay        []DailyStat `json:"tasks_by_day"`              // задачи по дням
	AvgCompletionTime float64     `json:"avg_completion_time_hours"` // среднее время выполнения

	// Активность
	MostProductiveDay string `json:"most_productive_day"` // самый продуктивный день

	// Дополнительная статистика
	TasksByHour     []HourlyStat  `json:"tasks_by_hour,omitempty"`     // задачи по часам
	CompletionByDay []DailyStat   `json:"completion_by_day,omitempty"` // выполненные по дням
	TasksByWeekday  []WeekdayStat `json:"tasks_by_weekday,omitempty"`  // задачи по дням недели
}

// Статистика по тегу
type TagStat struct {
	Tag        string `json:"tag"`
	UsageCount int    `json:"usage_count"`
}

// Статистика по дням
type DailyStat struct {
	Date  string `json:"date"` // YYYY-MM-DD
	Count int    `json:"count"`
}

// Статистика по часам
type HourlyStat struct {
	Hour  int `json:"hour"` // 0-23
	Count int `json:"count"`
}

// Статистика по дням недели
type WeekdayStat struct {
	Weekday string `json:"weekday"` // Monday, Tuesday, etc.
	Count   int    `json:"count"`
}
