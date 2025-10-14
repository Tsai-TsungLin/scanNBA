package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nba-scanner/internal/models"
	"time"
)

// FetchScheduleForDate 從完整賽季 API 取得指定日期的比賽
func FetchScheduleForDate(targetDate time.Time) (*models.NBAScoreboard, error) {
	url := "https://cdn.nba.com/static/json/staticData/scheduleLeagueV2_9.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch full schedule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("full schedule API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var schedule models.FullSchedule
	if err := json.Unmarshal(body, &schedule); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 格式化目標日期為 MM/DD/YYYY 00:00:00（符合 API 格式）
	targetDateStr := targetDate.Format("01/02/2006 00:00:00")

	// 查找該日期的比賽
	var games []models.Game
	for _, gameDate := range schedule.LeagueSchedule.GameDates {
		if gameDate.GameDate == targetDateStr {
			// 轉換為 NBAScoreboard 格式
			for _, g := range gameDate.Games {
				// 將 GameDateTimeEst 轉換為 UTC 格式
				gameTimeUTC := g.GameDateTimeEst

				game := models.Game{
					GameID:         g.GameID,
					GameCode:       g.GameCode,
					GameStatus:     g.GameStatus,
					GameStatusText: g.GameStatusText,
					Period:         0, // 預設值
					GameClock:      "",
					GameTimeUTC:    gameTimeUTC,
					HomeTeam: models.Team{
						TeamID:      g.HomeTeam.TeamID,
						TeamName:    g.HomeTeam.TeamName,
						TeamCity:    g.HomeTeam.TeamCity,
						TeamTricode: "",  // ScheduledTeam 沒有 TeamTricode
						Score:       g.HomeTeam.Score,
						Wins:        g.HomeTeam.Wins,
						Losses:      g.HomeTeam.Losses,
						Periods:     []models.PeriodScore{},
					},
					AwayTeam: models.Team{
						TeamID:      g.AwayTeam.TeamID,
						TeamName:    g.AwayTeam.TeamName,
						TeamCity:    g.AwayTeam.TeamCity,
						TeamTricode: "",  // ScheduledTeam 沒有 TeamTricode
						Score:       g.AwayTeam.Score,
						Wins:        g.AwayTeam.Wins,
						Losses:      g.AwayTeam.Losses,
						Periods:     []models.PeriodScore{},
					},
				}

				// 如果比賽進行中或已結束，從 boxscore 取得詳細數據
				if g.GameStatus == 2 || g.GameStatus == 3 {
					if boxscore, err := FetchBoxscore(g.GameID); err == nil {
						// 更新比賽狀態和時鐘
						game.Period = boxscore.Game.Period
						game.GameClock = boxscore.Game.GameClock
						game.GameStatusText = boxscore.Game.GameStatusText

						// 更新比分
						game.HomeTeam.Score = boxscore.Game.HomeTeam.Score
						game.AwayTeam.Score = boxscore.Game.AwayTeam.Score

						// 更新各節得分
						game.HomeTeam.Periods = boxscore.Game.HomeTeam.Periods
						game.AwayTeam.Periods = boxscore.Game.AwayTeam.Periods
					}
				}

				games = append(games, game)
			}
			break
		}
	}

	// 構建 NBAScoreboard 回應
	scoreboard := &models.NBAScoreboard{
		Scoreboard: struct {
			GameDate string        `json:"gameDate"`
			Games    []models.Game `json:"games"`
		}{
			GameDate: targetDateStr,
			Games:    games,
		},
	}

	return scoreboard, nil
}

// ShouldShowTomorrow 判斷是否應該顯示明天的比賽
// 台北時間 12:00 或 18:00 之後顯示明天
func ShouldShowTomorrow() bool {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		loc = time.FixedZone("CST", 8*60*60)
	}

	now := time.Now().In(loc)
	hour := now.Hour()

	// 12:00 (中午) 或 18:00 (晚上6點) 之後顯示明天
	return hour >= 12 || hour >= 18
}
