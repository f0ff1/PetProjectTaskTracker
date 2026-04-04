package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"TaskTracker/internal/model"
)

func (s *DataBaseRepo) GetTasksStats(ctx context.Context, userID int) (*model.TaskStats, error) {
	stats := &model.TaskStats{}

	if err := s.fillBasicStats(ctx, userID, stats); err != nil {
		return nil, fmt.Errorf("базовая статистика: %w", err)
	}

	if err := s.fillTopTags(ctx, userID, stats); err != nil {
		return nil, fmt.Errorf("топ-теги: %w", err)
	}

	if err := s.fillBestDay(ctx, userID, stats); err != nil {
		return nil, fmt.Errorf("лучший день: %w", err)
	}

	if err := s.fillLastDays(ctx, userID, stats); err != nil {
		return nil, fmt.Errorf("последние 7 дней: %w", err)
	}

	return stats, nil
}

func (s *DataBaseRepo) fillBasicStats(ctx context.Context, userID int, stats *model.TaskStats) error {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN completed THEN 1 END) as completed
		FROM tasks
		WHERE user_id = $1
	`

	err := s.dbPool.QueryRow(ctx, query, userID).Scan(&stats.Total, &stats.Completed)
	if err != nil {
		return err
	}

	stats.Pending = stats.Total - stats.Completed
	if stats.Total > 0 {
		stats.Rate = float64(stats.Completed) / float64(stats.Total) * 100
	}
	return nil
}

func (s *DataBaseRepo) fillTopTags(ctx context.Context, userID int, stats *model.TaskStats) error {
	query := `
		SELECT tag, COUNT(*) as count
		FROM (
			SELECT unnest(tags) as tag
			FROM tasks
			WHERE user_id = $1 AND tags IS NOT NULL AND array_length(tags, 1) > 0
		) t
		WHERE tag != ''
		GROUP BY tag
		ORDER BY count DESC
		LIMIT 3
	`

	rows, err := s.dbPool.Query(ctx, query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	stats.TopTags = []model.TagStat{}
	for rows.Next() {
		var tag string
		var count int
		if err := rows.Scan(&tag, &count); err != nil {
			return err
		}
		stats.TopTags = append(stats.TopTags, model.TagStat{
			Name:  tag,
			Count: count,
		})
	}
	return rows.Err()
}

func (s *DataBaseRepo) fillBestDay(ctx context.Context, userID int, stats *model.TaskStats) error {
	query := `
		SELECT TO_CHAR(created_at, 'Day') as weekday, COUNT(*) as count
		FROM tasks
		WHERE user_id = $1
		GROUP BY weekday
		ORDER BY count DESC
		LIMIT 1
	`

	var day string
	var count int
	err := s.dbPool.QueryRow(ctx, query, userID).Scan(&day, &count)
	if err != nil {
		stats.BestDay = ""
		return nil
	}

	stats.BestDay = strings.TrimSpace(day)
	return nil
}

func (s *DataBaseRepo) fillLastDays(ctx context.Context, userID int, stats *model.TaskStats) error {
	query := `
		SELECT DATE(created_at) as day, COUNT(*) as count
		FROM tasks
		WHERE user_id = $1 AND created_at >= CURRENT_DATE - INTERVAL '7 days'
		GROUP BY DATE(created_at)
		ORDER BY day
	`

	rows, err := s.dbPool.Query(ctx, query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	stats.Last7Days = make(map[string]int)
	for rows.Next() {
		var day time.Time
		var cnt int
		if err := rows.Scan(&day, &cnt); err != nil {
			return err
		}
		stats.Last7Days[day.Format("2006-01-02")] = cnt
	}
	return nil
}
