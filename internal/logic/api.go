package logic

import (
	"fmt"
	"log"
	"nba-scanner/internal/crawler"
	"nba-scanner/internal/models"
	"sync"
	"time"
)

// GetTodayGames 取得今日比賽資料（供 API 使用）
// 台北時間 12:00 之後會顯示明天的比賽
func GetTodayGames() (*models.APIResponse, error) {
	// 平行抓取三個資料源
	var (
		scoreboard *models.NBAScoreboard
		odds       *models.NBAOdds
		injuryMap  map[string][]string
		wg         sync.WaitGroup
		mu         sync.Mutex
		errors     []error
	)

	wg.Add(3)

	// 1. 抓取賽程（根據時間決定顯示今天還是昨天/明天）
	go func() {
		defer wg.Done()

		// 判斷要顯示哪一天的比賽
		loc, _ := time.LoadLocation("Asia/Taipei")
		now := time.Now().In(loc)

		var targetDate time.Time
		// 台北時間 14:00 之前顯示昨天，14:00 之後顯示今天
		// 這樣可以確保早上還能看到昨晚/今早的比賽
		if now.Hour() < 14 {
			targetDate = now.AddDate(0, 0, -1) // 昨天
		} else {
			targetDate = now // 今天
		}

		// 從完整賽季 API 取得比賽
		sb, err := crawler.FetchScheduleForDate(targetDate)
		if err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("賽程錯誤: %w", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		scoreboard = sb
		mu.Unlock()
	}()

	// 2. 抓取賠率
	go func() {
		defer wg.Done()
		od, err := crawler.FetchOdds()
		if err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("賠率錯誤: %w", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		odds = od
		mu.Unlock()
	}()

	// 3. 抓取傷兵
	go func() {
		defer wg.Done()
		im := crawler.FetchInjuryMap()
		mu.Lock()
		injuryMap = im
		mu.Unlock()
	}()

	wg.Wait()

	// 檢查錯誤
	if len(errors) > 0 {
		return nil, fmt.Errorf("抓取資料失敗: %v", errors)
	}

	// 建立賠率查找表
	oddsMap := crawler.BuildOddsMap(odds)

	// 轉換為 API 回應格式
	// 使用台北時間的今天日期
	taipeiLocation, _ := time.LoadLocation("Asia/Taipei")
	todayDate := time.Now().In(taipeiLocation).Format("2006-01-02")

	response := &models.APIResponse{
		Date:  todayDate,
		Games: make([]models.GameInfo, 0, len(scoreboard.Scoreboard.Games)),
	}

	for _, game := range scoreboard.Scoreboard.Games {
		gameInfo := buildGameInfo(&game, oddsMap, injuryMap)
		response.Games = append(response.Games, gameInfo)
	}

	return response, nil
}

// buildGameInfo 建立單場比賽資訊
func buildGameInfo(game *models.Game, oddsMap map[string]models.SpreadInfo, injuryMap map[string][]string) models.GameInfo {
	homeTeam := game.HomeTeam.GetFullTeamName()
	awayTeam := game.AwayTeam.GetFullTeamName()

	// 根據比賽狀態決定顯示時間或狀態
	var gameTimeDisplay, scoreDisplay string
	if game.GameStatus == 2 { // 進行中
		// 顯示 "Q4 04:00" 格式
		gameTimeDisplay = fmt.Sprintf("Q%d %s", game.Period, formatGameClock(game.GameClock))
		// 顯示即時比分 "[98-128]"
		scoreDisplay = fmt.Sprintf("[%d-%d]", game.AwayTeam.Score, game.HomeTeam.Score)
	} else if game.GameStatus == 3 { // 已結束
		gameTimeDisplay = "已結束"
		// 顯示最終比分 "[111-109]"
		scoreDisplay = fmt.Sprintf("[%d-%d]", game.AwayTeam.Score, game.HomeTeam.Score)
	} else { // 未開始
		gameTime, err := crawler.ConvertUTCToLocal(game.GameTimeUTC)
		if err != nil {
			log.Printf("時間轉換錯誤: %v", err)
			gameTimeDisplay = "00:00"
		} else {
			gameTimeDisplay = gameTime
		}
		scoreDisplay = "" // 未開始不顯示比分
	}

	// 取得中文隊名
	homeTeamCN := models.TeamMap[homeTeam]
	if homeTeamCN == "" {
		homeTeamCN = homeTeam
	}
	awayTeamCN := models.TeamMap[awayTeam]
	if awayTeamCN == "" {
		awayTeamCN = awayTeam
	}

	// 取得盤口資訊（根據比賽狀態決定顯示哪個盤口）
	spreadDisplay := models.SpreadDisplay{HasData: false}
	if spread, ok := oddsMap[game.GameID]; ok && spread.Found {
		spreadDisplay.HasData = true

		if game.GameStatus == 1 { // 未開始：顯示開盤和當前盤口
			spreadDisplay.Opening = fmt.Sprintf("%.1f", spread.HomeOpeningSpread)
			spreadDisplay.Current = spread.HomeSpread // 當前盤口
		} else { // 進行中或已結束：顯示開盤和最後盤口
			// 已結束時，current 就是開打前的最後盤口
			spreadDisplay.Opening = fmt.Sprintf("%.1f", spread.HomeOpeningSpread)
			spreadDisplay.Current = spread.HomeSpread
		}
	}

	// 取得傷兵資訊
	homeInjuries := getInjuriesForTeam(homeTeam, injuryMap)
	awayInjuries := getInjuriesForTeam(awayTeam, injuryMap)

	// 取得主隊近五場戰績（使用快取）
	var homeHistory *models.TeamHistory
	if history, err := crawler.FetchTeamHistoryWithCache(game.HomeTeam.TeamID, 5); err == nil {
		homeHistory = history
	} else {
		log.Printf("取得主隊戰績失敗 (TeamID: %d): %v", game.HomeTeam.TeamID, err)
	}

	// 取得客隊近五場戰績（使用快取）
	var awayHistory *models.TeamHistory
	if history, err := crawler.FetchTeamHistoryWithCache(game.AwayTeam.TeamID, 5); err == nil {
		awayHistory = history
	} else {
		log.Printf("取得客隊戰績失敗 (TeamID: %d): %v", game.AwayTeam.TeamID, err)
	}

	// 如果比賽進行中或已結束，取得球員數據和各節比分
	var homePlayers, awayPlayers []models.PlayerDisplay
	var periodScores *models.PeriodScores
	if game.GameStatus == 2 || game.GameStatus == 3 { // 進行中或已結束
		if boxscore, err := crawler.FetchBoxscore(game.GameID); err == nil {
			// 處理主隊球員
			homePlayerList := crawler.BuildPlayerDisplayList(boxscore.Game.HomeTeam.Players)
			homePlayers = crawler.SortPlayersByStarterAndPoints(homePlayerList)

			// 處理客隊球員
			awayPlayerList := crawler.BuildPlayerDisplayList(boxscore.Game.AwayTeam.Players)
			awayPlayers = crawler.SortPlayersByStarterAndPoints(awayPlayerList)
		} else {
			log.Printf("取得球員數據失敗 (GameID: %s): %v", game.GameID, err)
		}

		// 處理各節比分
		periodScores = buildPeriodScores(game)
	}

	return models.GameInfo{
		GameID: game.GameID,
		GameTime:       gameTimeDisplay,
		GameStatus:     game.GameStatus,
		GameStatusText: game.GameStatusText,
		ScoreDisplay:   scoreDisplay,
		HomeTeam: models.TeamInfo{
			NameEN: homeTeam,
			NameCN: homeTeamCN,
			Wins:   game.HomeTeam.Wins,
			Losses: game.HomeTeam.Losses,
		},
		AwayTeam: models.TeamInfo{
			NameEN: awayTeam,
			NameCN: awayTeamCN,
			Wins:   game.AwayTeam.Wins,
			Losses: game.AwayTeam.Losses,
		},
		HomeScore:    game.HomeTeam.Score,
		AwayScore:    game.AwayTeam.Score,
		Spread:       spreadDisplay,
		HomeInjuries: homeInjuries,
		AwayInjuries: awayInjuries,
		HomeHistory:  homeHistory,
		AwayHistory:  awayHistory,
		HomePlayers:  homePlayers,
		AwayPlayers:  awayPlayers,
		PeriodScores: periodScores,
	}
}

// buildPeriodScores 建立各節比分資訊
func buildPeriodScores(game *models.Game) *models.PeriodScores {
	// 確保有四節的資料（如果比賽還在進行中，未完成的節顯示 0）
	homePeriods := make([]int, 4)
	awayPeriods := make([]int, 4)

	// 填入主隊各節得分
	for _, period := range game.HomeTeam.Periods {
		if period.Period >= 1 && period.Period <= 4 {
			homePeriods[period.Period-1] = period.Score
		}
	}

	// 填入客隊各節得分
	for _, period := range game.AwayTeam.Periods {
		if period.Period >= 1 && period.Period <= 4 {
			awayPeriods[period.Period-1] = period.Score
		}
	}

	return &models.PeriodScores{
		HomePeriods: homePeriods,
		AwayPeriods: awayPeriods,
	}
}

// formatGameClock 格式化比賽時鐘，從 "PT04M03.00S" 轉為 "04:03"
func formatGameClock(clock string) string {
	if clock == "" {
		return "00:00"
	}

	// NBA API 格式: PT04M03.00S
	// 需要解析出分鐘和秒數
	var minutes, seconds int
	fmt.Sscanf(clock, "PT%dM%d", &minutes, &seconds)

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// getInjuriesForTeam 取得球隊傷兵清單
func getInjuriesForTeam(teamName string, injuryMap map[string][]string) []string {
	// 處理 LA Clippers 特殊情況
	if teamName == "Los Angeles Clippers" {
		teamName = "LA Clippers"
	}

	if injuries, ok := injuryMap[teamName]; ok && len(injuries) > 0 {
		return injuries
	}
	return []string{}
}
