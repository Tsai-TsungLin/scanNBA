package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nba-scanner/internal/models"
	"time"
)

// FetchSchedule 抓取今日賽程
func FetchSchedule() (*models.NBAScoreboard, error) {
	url := "https://nba-prod-us-east-1-mediaops-stats.s3.amazonaws.com/NBA/liveData/scoreboard/todaysScoreboard_00.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("schedule API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read schedule response: %w", err)
	}

	var scoreboard models.NBAScoreboard
	if err := json.Unmarshal(body, &scoreboard); err != nil {
		return nil, fmt.Errorf("failed to parse schedule JSON: %w", err)
	}

	return &scoreboard, nil
}

// ConvertUTCToLocal 將 UTC 時間轉為本地顯示時間 (+8 時區)
func ConvertUTCToLocal(utcTimeStr string) (string, error) {
	t, err := time.Parse(time.RFC3339, utcTimeStr)
	if err != nil {
		return "", err
	}

	// 轉換到 UTC+8
	localTime := t.Add(8 * time.Hour)
	return localTime.Format("15:04"), nil
}
