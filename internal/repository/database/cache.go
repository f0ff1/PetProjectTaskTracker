package database

import (
	"sync"
	"time"

	"TaskTracker/internal/model"
)

type userStatsCache struct {
	mu         sync.RWMutex
	stats      *model.TaskStats
	lastUpdate time.Time
	isUpdating bool
}

type statsCache struct {
	mu    sync.RWMutex
	cache map[int]*userStatsCache // key = userID
}

func newStatsCache() *statsCache {
	return &statsCache{
		cache: make(map[int]*userStatsCache),
	}
}

func (c *statsCache) getOrCreate(userID int) *userStatsCache {
	if _, exists := c.cache[userID]; !exists {
		c.cache[userID] = &userStatsCache{
			stats:      nil,
			lastUpdate: time.Time{},
			isUpdating: false,
		}
	}
	return c.cache[userID]
}

func (c *statsCache) getAllCacheData(userID int) (*model.TaskStats, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userCache := c.getOrCreate(userID)
	return userCache.stats, userCache.lastUpdate, userCache.isUpdating
}

func (c *statsCache) set(userID int, stats *model.TaskStats) {
	c.mu.Lock()
	defer c.mu.Unlock()

	userCache := c.getOrCreate(userID)
	userCache.stats = stats
	userCache.lastUpdate = time.Now()
	userCache.isUpdating = false
}

func (c *statsCache) setUpdating(userID int, updating bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	userCache := c.getOrCreate(userID)
	userCache.isUpdating = updating
}

func (c *statsCache) getIsUpdating(userID int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userCache := c.getOrCreate(userID)
	return userCache.isUpdating
}

func (c *statsCache) hasStats(userID int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userCache := c.getOrCreate(userID)
	return userCache.stats != nil
}

func (c *statsCache) getLastUpdate(userID int) time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userCache := c.getOrCreate(userID)
	return userCache.lastUpdate
}
