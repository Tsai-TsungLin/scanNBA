package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nba-scanner/internal/models"
)

// FetchOdds 抓取今日賠率
func FetchOdds() (*models.NBAOdds, error) {
	url := "https://cdn.nba.com/static/json/liveData/odds/odds_todaysGames.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch odds: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("odds API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read odds response: %w", err)
	}

	var odds models.NBAOdds
	if err := json.Unmarshal(body, &odds); err != nil {
		return nil, fmt.Errorf("failed to parse odds JSON: %w", err)
	}

	return &odds, nil
}

// BuildOddsMap 建立 gameId -> SpreadInfo 的 map
func BuildOddsMap(odds *models.NBAOdds) map[string]models.SpreadInfo {
	oddsMap := make(map[string]models.SpreadInfo)

	for _, game := range odds.Games {
		spreadInfo := game.GetSupermatchSpread()
		oddsMap[game.GameID] = spreadInfo
	}

	return oddsMap
}
