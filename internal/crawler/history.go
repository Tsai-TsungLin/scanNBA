package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nba-scanner/internal/models"
	"sort"
	"time"
)

// convertGameDateTime 將 ISO 時間轉換為台北時間 "2025/10/02 08:00:00"
func convertGameDateTime(isoTime string) string {
	// NBA API 的 gameDateTimeEst 名為 EST，帶 Z 結尾但實際是 EST 時區的時間值
	// 例如: "2025-10-07T21:30:00Z" 表示 EST 21:30（不是 UTC 21:30）
	// 需要先解析，然後加上 EST->UTC 的偏移（+4 夏令時或 +5 標準時間），再轉台北

	// 移除 Z，當成無時區的時間解析
	isoTimeNoZ := isoTime
	if len(isoTime) > 0 && isoTime[len(isoTime)-1] == 'Z' {
		isoTimeNoZ = isoTime[:len(isoTime)-1]
	}

	// 解析為無時區的時間
	t, err := time.Parse("2006-01-02T15:04:05", isoTimeNoZ)
	if err != nil {
		return isoTime // 解析失敗，返回原始字串
	}

	// 將這個時間視為 EST 時區（UTC-5 標準時間，UTC-4 夏令時）
	// 10月通常還在夏令時（EDT），使用 UTC-4
	estLocation := time.FixedZone("EST", -4*60*60)
	tEST := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, estLocation)

	// 轉換為台北時區 (UTC+8)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.FixedZone("CST", 8*60*60)
	}
	tTaipei := tEST.In(loc)

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

	// 先找出球隊英文名稱（用於查詢 titan007）
	var teamNameEN string
	for _, gameDate := range schedule.LeagueSchedule.GameDates {
		for _, game := range gameDate.Games {
			if game.HomeTeam.TeamID == teamID {
				teamNameEN = game.HomeTeam.TeamCity + " " + game.HomeTeam.TeamName
				break
			} else if game.AwayTeam.TeamID == teamID {
				teamNameEN = game.AwayTeam.TeamCity + " " + game.AwayTeam.TeamName
				break
			}
		}
		if teamNameEN != "" {
			break
		}
	}

	// 從 titan007 獲取近5場過盤結果
	titan007Spreads, hasTitan007 := GetTeamSpreads(teamNameEN)
	if !hasTitan007 {
		log.Printf("警告：未從 titan007 獲取到 %s 的過盤資料", teamNameEN)
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

				// 計算實際勝負
				var gameResult string
				if game.HomeTeam.Score > game.AwayTeam.Score {
					gameResult = "W"
				} else {
					gameResult = "L"
				}

				// 使用 titan007 的過盤結果（按時間順序，最新的在前）
				var spreadResult string
				var spread string
				var hasSpread bool

				// 當前收集的比賽數量對應 titan007Spreads 的 index
				currentIndex := len(games)
				if hasTitan007 && currentIndex < len(titan007Spreads) {
					spreadResult = titan007Spreads[currentIndex]
					spread = "有盤口"
					hasSpread = true
					result = spreadResult
				} else {
					spread = "無盤口"
					spreadResult = gameResult // 無盤口時，過盤結果等於實際勝負
					result = gameResult
					hasSpread = false
				}

				// 轉換為中文隊名
				opponentCN := models.TeamMap[opponent]
				if opponentCN == "" {
					opponentCN = opponent
				}

				// 分離日期和時間
				gameDateTime := convertGameDateTime(game.GameDateTimeEst)
				dateTimeParts := splitDateTime(gameDateTime)

				// 提取 NBA 原始日期（用於跳轉查詢）
				nbaGameDate := extractNBAGameDate(game.GameDateTimeEst)

				games = append(games, models.GameResult{
					GameID:       game.GameID,
					GameDate:     nbaGameDate,
					Date:         dateTimeParts[0],
					Time:         dateTimeParts[1],
					Opponent:     opponentCN,
					VsIndicator:  "vs",
					IsHome:       isHome,
					Score:        score,
					GameResult:   gameResult,
					SpreadResult: spreadResult,
					Result:       result,
					Spread:       spread,
					HasSpread:    hasSpread,
				})
			} else if game.AwayTeam.TeamID == teamID {
				// 客場比賽
				isHome = false
				opponent = game.HomeTeam.TeamCity + " " + game.HomeTeam.TeamName
				score = fmt.Sprintf("%d-%d", game.AwayTeam.Score, game.HomeTeam.Score)

				// 計算實際勝負
				var gameResult string
				if game.AwayTeam.Score > game.HomeTeam.Score {
					gameResult = "W"
				} else {
					gameResult = "L"
				}

				// 使用 titan007 的過盤結果（按時間順序，最新的在前）
				var spreadResult string
				var spread string
				var hasSpread bool

				// 當前收集的比賽數量對應 titan007Spreads 的 index
				currentIndex := len(games)
				if hasTitan007 && currentIndex < len(titan007Spreads) {
					spreadResult = titan007Spreads[currentIndex]
					spread = "有盤口"
					hasSpread = true
					result = spreadResult
				} else {
					spread = "無盤口"
					spreadResult = gameResult // 無盤口時，過盤結果等於實際勝負
					result = gameResult
					hasSpread = false
				}

				// 轉換為中文隊名
				opponentCN := models.TeamMap[opponent]
				if opponentCN == "" {
					opponentCN = opponent
				}

				// 分離日期和時間
				gameDateTime := convertGameDateTime(game.GameDateTimeEst)
				dateTimeParts := splitDateTime(gameDateTime)

				// 提取 NBA 原始日期（用於跳轉查詢）
				nbaGameDate := extractNBAGameDate(game.GameDateTimeEst)

				games = append(games, models.GameResult{
					GameID:       game.GameID,
					GameDate:     nbaGameDate,
					Date:         dateTimeParts[0],
					Time:         dateTimeParts[1],
					Opponent:     opponentCN,
					VsIndicator:  "@",
					IsHome:       isHome,
					Score:        score,
					GameResult:   gameResult,
					SpreadResult: spreadResult,
					Result:       result,
					Spread:       spread,
					HasSpread:    hasSpread,
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

// splitDateTime 分離日期和時間
// 輸入: "2025/10/02 08:00:00"
// 輸出: ["2025/10/02", "08:00"]
func splitDateTime(datetime string) []string {
	parts := []string{"", ""}

	// 使用空格分割日期和時間
	if len(datetime) > 10 {
		parts[0] = datetime[:10] // 日期部分
		if len(datetime) > 11 {
			timePart := datetime[11:]
			// 只取 HH:MM 部分
			if len(timePart) >= 5 {
				parts[1] = timePart[:5]
			}
		}
	} else {
		parts[0] = datetime
	}

	return parts
}

// GetTeamIDFromName 從球隊名稱獲取 TeamID（簡化版，實際應該從 API 獲取）
func GetTeamIDFromGame(game *models.Game) int {
	return game.HomeTeam.TeamID
}

// extractNBAGameDate 從 GameDateTimeEst 提取 NBA 原始比賽日期
// 輸入: "2025-10-07T21:30:00Z"（實際為 EST 時區）
// 輸出: "2025-10-07"（NBA API 的原始日期，用於跳轉查詢）
func extractNBAGameDate(isoTime string) string {
	// GameDateTimeEst 格式: "2025-10-07T21:30:00Z"
	// 雖然帶 Z，但實際是 EST 時區的日期
	// 直接取前 10 個字符作為日期
	if len(isoTime) >= 10 {
		return isoTime[:10] // "2025-10-07"
	}
	return isoTime
}
