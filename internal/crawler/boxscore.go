package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nba-scanner/internal/models"
)

// FetchBoxscore 抓取比賽的 boxscore 數據
func FetchBoxscore(gameID string) (*models.BoxscoreResponse, error) {
	url := fmt.Sprintf("https://cdn.nba.com/static/json/liveData/boxscore/boxscore_%s.json", gameID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch boxscore: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("boxscore API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read boxscore response: %w", err)
	}

	var boxscore models.BoxscoreResponse
	if err := json.Unmarshal(body, &boxscore); err != nil {
		return nil, fmt.Errorf("failed to parse boxscore JSON: %w", err)
	}

	return &boxscore, nil
}

// BuildPlayerDisplayList 建立顯示用的球員列表
func BuildPlayerDisplayList(players []models.BoxscorePlayer) []models.PlayerDisplay {
	var result []models.PlayerDisplay

	for _, p := range players {
		// 只顯示有上場的球員
		if p.Played != "1" {
			continue
		}

		// 格式化上場時間
		minutes := formatMinutes(p.Statistics.Minutes)

		result = append(result, models.PlayerDisplay{
			Name:      p.Name,
			NameI:     p.NameI,
			JerseyNum: p.JerseyNum,
			Position:  p.Position,
			IsStarter: p.Starter == "1",
			OnCourt:   p.OnCourt == "1",
			Points:    p.Statistics.Points,
			Rebounds:  p.Statistics.ReboundsTotal,
			Assists:   p.Statistics.Assists,
			Minutes:   minutes,
		})
	}

	return result
}

// formatMinutes 格式化分鐘數，從 "PT21M42.00S" 轉為 "21:42"
func formatMinutes(minutes string) string {
	if minutes == "" {
		return "0:00"
	}

	// NBA API 格式: PT21M42.00S
	var min, sec int
	fmt.Sscanf(minutes, "PT%dM%d", &min, &sec)

	return fmt.Sprintf("%d:%02d", min, sec)
}

// SortPlayersByStarterAndPoints 按先發/替補和得分排序
func SortPlayersByStarterAndPoints(players []models.PlayerDisplay) []models.PlayerDisplay {
	// 分成先發和替補
	var starters, benchers []models.PlayerDisplay

	for _, p := range players {
		if p.IsStarter {
			starters = append(starters, p)
		} else {
			benchers = append(benchers, p)
		}
	}

	// 排序函數：按得分降序
	sortByPoints := func(list []models.PlayerDisplay) {
		for i := 0; i < len(list)-1; i++ {
			for j := i + 1; j < len(list); j++ {
				if list[i].Points < list[j].Points {
					list[i], list[j] = list[j], list[i]
				}
			}
		}
	}

	sortByPoints(starters)
	sortByPoints(benchers)

	// 合併：先發在前，替補在後
	return append(starters, benchers...)
}
