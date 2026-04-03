package database

import (
	"context"
	"log"
	"log/slog"
	"time"

	myErrors "TaskTracker/errors"
	"TaskTracker/internal/model"
)

func (s *DataBaseRepo) StartStatsUpdater(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		s.UpdateStats(ctx)

		for {
			select {
			case <-ticker.C:
				s.UpdateStats(ctx)
			case <-ctx.Done():
				slog.Info("Остановка обновления статистики")
				return
			}
		}

	}()
}

func (s *DataBaseRepo) UpdateStats(ctx context.Context) {
	if s.cache.getIsUpdating() {
		return
	}

	s.cache.setUpdating(true)

	defer s.cache.setUpdating(false)

	stats, err := s.GetTasksStats(ctx)
	if err != nil {
		slog.Error("Ошибка обновления статистики", "err", err)
		return
	}
	slog.Info("Статистика обновлена")
	s.cache.set(stats)
}

func (s *DataBaseRepo) GetStatsFromCache(ctx context.Context) (*model.TaskStats, error) {

	if !s.cache.hasStats() {
		return nil, myErrors.ErrStatsDoesntWritten
	}

	stats, _, _ := s.cache.getAllCacheData()
	return stats, nil

}

func (s *DataBaseRepo) GetStatsWithInfo(ctx context.Context) (*model.TaskStats, time.Time, bool, error) {
	stats, lastUpdate, isUpdating := s.cache.getAllCacheData()
	return stats, lastUpdate, isUpdating, nil
}

func (s *DataBaseRepo) GetStatsWithRefresh(ctx context.Context, forceRefresh bool) (*model.TaskStats, error) {
	s.cache.mu.RLock()
	cacheAge := time.Since(s.cache.lastUpdate)
	hasStats := s.cache.stats != nil
	s.cache.mu.RUnlock()

	// Если кэш свежий (менее 1 минуты) и не нужен форс-обновление
	if !forceRefresh && hasStats && cacheAge < 1*time.Minute {
		return s.GetStatsFromCache(ctx)
	}

	// Иначе обновляем
	stats, err := s.GetTasksStats(ctx)
	if err != nil {
		// Если ошибка, но есть кэш — вернем его
		if hasStats {
			log.Printf("⚠️ Ошибка обновления статистики, возвращаю кэш: %v", err)
			return s.GetStatsFromCache(ctx)
		}
		return nil, err
	}

	// Сохраняем в кэш
	s.cache.mu.Lock()
	s.cache.stats = stats
	s.cache.lastUpdate = time.Now()
	s.cache.isUpdating = false
	s.cache.mu.Unlock()

	return stats, nil
}
