package database

import (
	"sync"
	"time"

	"TaskTracker/internal/model"
)

type statsCache struct {
	mu         sync.RWMutex
	stats      *model.TaskStats
	lastUpdate time.Time
	isUpdating bool
}

func newStatsCache() *statsCache {
	return &statsCache{
		stats:      nil,
		lastUpdate: time.Time{},
		isUpdating: false,
	}
}

func (c *statsCache) getAllCacheData() (*model.TaskStats, time.Time, bool) {

	c.mu.RLock()

	defer c.mu.RUnlock()

	return c.stats, c.lastUpdate, c.isUpdating
}

func (c *statsCache) set(stats *model.TaskStats) {

	c.mu.Lock()

	defer c.mu.Unlock()

	c.stats = stats
	c.lastUpdate = time.Now()
	c.isUpdating = false
}

func (c *statsCache) setUpdating(updating bool) {

	c.mu.Lock()

	c.isUpdating = updating

	c.mu.Unlock()
}

func (c *statsCache) getIsUpdating() bool {

	c.mu.Lock()

	defer c.mu.Unlock()

	return c.isUpdating
}

func (c *statsCache) hasStats() bool {

	c.mu.RLock()

	defer c.mu.RUnlock()

	return c.stats != nil
}
