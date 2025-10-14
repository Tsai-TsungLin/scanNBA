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

// ConvertUTCToLocal 將 EST 時間轉為台北時間
// NBA API 的 GameTimeUTC 欄位名稱有誤導性，實際是 EST 時間（UTC-4）
func ConvertUTCToLocal(estTimeStr string) (string, error) {
	// 解析時間（移除 Z，因為這不是真正的 UTC 時間）
	estTimeNoZ := estTimeStr
	if len(estTimeStr) > 0 && estTimeStr[len(estTimeStr)-1] == 'Z' {
		estTimeNoZ = estTimeStr[:len(estTimeStr)-1]
	}

	// 解析為無時區的時間
	t, err := time.Parse("2006-01-02T15:04:05", estTimeNoZ)
	if err != nil {
		return "", err
	}

	// 將這個時間視為 EST 時區（UTC-4 夏令時）
	estLocation := time.FixedZone("EST", -4*60*60)
	tEST := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, estLocation)

	// 轉換為台北時區 (UTC+8)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.FixedZone("CST", 8*60*60)
	}
	tTaipei := tEST.In(loc)

	return tTaipei.Format("15:04"), nil
}
