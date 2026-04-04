package model

type TagStat struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type TaskStats struct {
	Total     int     `json:"total"`
	Completed int     `json:"completed"`
	Pending   int     `json:"pending"`
	Rate      float64 `json:"rate"`

	TopTags []TagStat `json:"top_tags"`
	BestDay string    `json:"best_day"`

	Last7Days map[string]int `json:"last_7_days,omitempty"` // {"2024-01-01": 5, ...}
}
