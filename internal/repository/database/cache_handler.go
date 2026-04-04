package database

import (
	"context"
	"log"
	"log/slog"
	"time"

	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
)

func (s *DataBaseRepo) StartStatsUpdater(ctx context.Context, userID int, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		s.UpdateStats(ctx, userID)

		for {
			select {
			case <-ticker.C:
				s.UpdateStats(ctx, userID)
			case <-ctx.Done():
				slog.Info("Остановка обновления статистики")
				return
			}
		}

	}()
}

func (s *DataBaseRepo) UpdateStats(ctx context.Context, userID int) {
	if s.cache.getIsUpdating(userID) {
		return
	}

	s.cache.setUpdating(userID, true)

	defer s.cache.setUpdating(userID, false)

	stats, err := s.GetTasksStats(ctx, userID)
	if err != nil {
		slog.Error("Ошибка обновления статистики", "user_id", userID, "err", err)
		return
	}
	slog.Info("Статистика обновлена", "user_id", userID)
	s.cache.set(userID, stats)
}

func (s *DataBaseRepo) GetStatsFromCache(ctx context.Context, userID int) (*model.TaskStats, error) {

	if !s.cache.hasStats(userID) {
		return nil, myErrors.ErrStatsDoesntWritten
	}

	stats, _, _ := s.cache.getAllCacheData(userID)
	return stats, nil

}

func (s *DataBaseRepo) GetStatsWithInfo(ctx context.Context, userID int) (*model.TaskStats, time.Time, bool, error) {
	stats, lastUpdate, isUpdating := s.cache.getAllCacheData(userID)
	return stats, lastUpdate, isUpdating, nil
}

func (s *DataBaseRepo) GetStatsWithRefresh(ctx context.Context, userID int, forceRefresh bool) (*model.TaskStats, error) {
	s.cache.mu.RLock()
	cacheAge := time.Since(s.cache.getLastUpdate(userID))
	hasStats := s.cache.hasStats(userID)
	s.cache.mu.RUnlock()

	// Если кэш свежий (менее 1 минуты) и не нужен форс-обновление
	if !forceRefresh && hasStats && cacheAge < 1*time.Minute {
		return s.GetStatsFromCache(ctx, userID)
	}

	// Иначе обновляем
	stats, err := s.GetTasksStats(ctx, userID)
	if err != nil {
		// Если ошибка, но есть кэш — вернем его
		if hasStats {
			log.Printf("⚠️ Ошибка обновления статистики, возвращаю кэш: %v", err)
			return s.GetStatsFromCache(ctx, userID)
		}
		return nil, err
	}

	// Сохраняем в кэш
	s.cache.set(userID, stats)

	return stats, nil
}
