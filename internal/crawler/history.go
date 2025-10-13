package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nba-scanner/internal/models"
	"sort"
	"time"
)

// convertGameDateTime 將 ISO 時間轉換為台北時間 "2025/10/02 08:00:00"
func convertGameDateTime(isoTime string) string {
	// 解析 ISO 格式 "2025-10-02T12:00:00Z"
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		// 嘗試解析沒有 Z 的格式
		t, err = time.Parse("2006-01-02T15:04:05", isoTime)
		if err != nil {
			return isoTime // 解析失敗，返回原始字串
		}
	}

	// 轉換為台北時區 (UTC+8)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.FixedZone("CST", 8*60*60) // 如果找不到時區，使用固定偏移
	}
	tTaipei := t.In(loc)

	// 格式化為 "2025/10/02 08:00:00"
	return tTaipei.Format("2006/01/02 15:04:05")
}

// FetchTeamHistory 抓取球隊近期戰績
func FetchTeamHistory(teamID int, limit int) (*models.TeamHistory, error) {
	// 抓取完整賽季賽程
	url := "https://cdn.nba.com/static/json/staticData/scheduleLeagueV2_9.json"

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
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var schedule models.FullSchedule
	if err := json.Unmarshal(body, &schedule); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 收集該球隊的所有比賽
	var games []models.GameResult
	for _, gameDate := range schedule.LeagueSchedule.GameDates {
		for _, game := range gameDate.Games {
			// 只處理已結束的比賽
			if game.GameStatus != 3 {
				continue
			}

			// 檢查是否為該球隊的比賽
			var isHome bool
			var opponent string
			var score string
			var result string

			if game.HomeTeam.TeamID == teamID {
				// 主場比賽
				isHome = true
				opponent = game.AwayTeam.TeamCity + " " + game.AwayTeam.TeamName
				score = fmt.Sprintf("%d-%d", game.HomeTeam.Score, game.AwayTeam.Score)
				if game.HomeTeam.Score > game.AwayTeam.Score {
					result = "W"
				} else {
					result = "L"
				}

				// 轉換為中文隊名
				opponentCN := models.TeamMap[opponent]
				if opponentCN == "" {
					opponentCN = opponent // 如果找不到對應，保持原文
				}

				games = append(games, models.GameResult{
					Date:     convertGameDateTime(game.GameDateTimeEst),
					Opponent: opponentCN,
					IsHome:   isHome,
					Score:    score,
					Result:   result,
				})
			} else if game.AwayTeam.TeamID == teamID {
				// 客場比賽
				isHome = false
				opponent = game.HomeTeam.TeamCity + " " + game.HomeTeam.TeamName
				score = fmt.Sprintf("%d-%d", game.AwayTeam.Score, game.HomeTeam.Score)
				if game.AwayTeam.Score > game.HomeTeam.Score {
					result = "W"
				} else {
					result = "L"
				}

				// 轉換為中文隊名
				opponentCN := models.TeamMap[opponent]
				if opponentCN == "" {
					opponentCN = opponent // 如果找不到對應，保持原文
				}

				games = append(games, models.GameResult{
					Date:     convertGameDateTime(game.GameDateTimeEst),
					Opponent: opponentCN,
					IsHome:   isHome,
					Score:    score,
					Result:   result,
				})
			}
		}
	}

	// 按日期排序（最新的在前）
	sort.Slice(games, func(i, j int) bool {
		dateI, _ := time.Parse("2006/01/02 15:04:05", games[i].Date)
		dateJ, _ := time.Parse("2006/01/02 15:04:05", games[j].Date)
		return dateI.After(dateJ)
	})

	// 只取最近 N 場
	if len(games) > limit {
		games = games[:limit]
	}

	// 計算勝負場數
	wins := 0
	losses := 0
	for _, game := range games {
		if game.Result == "W" {
			wins++
		} else {
			losses++
		}
	}

	return &models.TeamHistory{
		TeamID:      teamID,
		RecentGames: games,
		WinCount:    wins,
		LossCount:   losses,
	}, nil
}

// GetTeamIDFromName 從球隊名稱獲取 TeamID（簡化版，實際應該從 API 獲取）
func GetTeamIDFromGame(game *models.Game) int {
	return game.HomeTeam.TeamID
}
