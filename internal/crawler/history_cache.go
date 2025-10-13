package crawler

import (
	"nba-scanner/internal/models"
	"sync"
	"time"
)

// HistoryCache 戰績快取
type HistoryCache struct {
	data      map[int]*models.TeamHistory // teamID -> 戰績
	mu        sync.RWMutex
	lastUpdate time.Time
	cacheDuration time.Duration
}

var (
	historyCache *HistoryCache
	once         sync.Once
)

// GetHistoryCache 取得快取單例
func GetHistoryCache() *HistoryCache {
	once.Do(func() {
		historyCache = &HistoryCache{
			data:          make(map[int]*models.TeamHistory),
			cacheDuration: 1 * time.Hour, // 快取 1 小時
		}
	})
	return historyCache
}

// Get 從快取取得球隊戰績
func (c *HistoryCache) Get(teamID int) (*models.TeamHistory, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 檢查快取是否過期
	if time.Since(c.lastUpdate) > c.cacheDuration {
		return nil, false
	}

	history, exists := c.data[teamID]
	return history, exists
}

// Set 設定球隊戰績到快取
func (c *HistoryCache) Set(teamID int, history *models.TeamHistory) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[teamID] = history
	c.lastUpdate = time.Now()
}

// Clear 清除所有快取
func (c *HistoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[int]*models.TeamHistory)
}

// IsExpired 檢查快取是否過期
func (c *HistoryCache) IsExpired() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return time.Since(c.lastUpdate) > c.cacheDuration
}

// FetchTeamHistoryWithCache 使用快取的版本
func FetchTeamHistoryWithCache(teamID int, limit int) (*models.TeamHistory, error) {
	cache := GetHistoryCache()

	// 先嘗試從快取取得
	if history, exists := cache.Get(teamID); exists {
		return history, nil
	}

	// 快取未命中，抓取新資料
	history, err := FetchTeamHistory(teamID, limit)
	if err != nil {
		return nil, err
	}

	// 存入快取
	cache.Set(teamID, history)

	return history, nil
}
