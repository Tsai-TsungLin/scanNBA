package logic

import (
	"fmt"
	"nba-scanner/internal/crawler"
	"nba-scanner/internal/models"
	"sync"
	"time"
)

// PKTeam 主要功能：抓取並顯示今日所有比賽資訊
func PKTeam() {
	start := time.Now()

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

	// 1. 抓取賽程
	go func() {
		defer wg.Done()
		sb, err := crawler.FetchSchedule()
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
		fmt.Println("抓取資料時發生錯誤：")
		for _, err := range errors {
			fmt.Println("  -", err)
		}
		return
	}

	// 建立賠率查找表
	oddsMap := crawler.BuildOddsMap(odds)

	// 顯示比賽資訊
	fmt.Printf("今天 %s 有 %d 場比賽\n\n", scoreboard.Scoreboard.GameDate, len(scoreboard.Scoreboard.Games))

	for i, game := range scoreboard.Scoreboard.Games {
		displayGame(i+1, &game, oddsMap, injuryMap)
	}

	fmt.Printf("\n執行時間：%v\n", time.Since(start))
}

// PKTeamOnStartTime 根據開賽時間篩選比賽
func PKTeamOnStartTime(st string) {
	if _, err := time.Parse("15:04", st); err != nil {
		fmt.Println("時間格式錯誤，請用 15:04 格式")
		return
	}

	start := time.Now()

	// 平行抓取三個資料源（同上）
	var (
		scoreboard *models.NBAScoreboard
		odds       *models.NBAOdds
		injuryMap  map[string][]string
		wg         sync.WaitGroup
		mu         sync.Mutex
		errors     []error
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		sb, err := crawler.FetchSchedule()
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

	go func() {
		defer wg.Done()
		im := crawler.FetchInjuryMap()
		mu.Lock()
		injuryMap = im
		mu.Unlock()
	}()

	wg.Wait()

	if len(errors) > 0 {
		fmt.Println("抓取資料時發生錯誤：")
		for _, err := range errors {
			fmt.Println("  -", err)
		}
		return
	}

	oddsMap := crawler.BuildOddsMap(odds)

	// 篩選指定時間的比賽
	var matchedGames []models.Game
	for _, game := range scoreboard.Scoreboard.Games {
		gameTime, err := crawler.ConvertUTCToLocal(game.GameTimeUTC)
		if err != nil {
			continue
		}
		if gameTime == st {
			matchedGames = append(matchedGames, game)
		}
	}

	if len(matchedGames) == 0 {
		fmt.Printf("今天 %s 沒有比賽\n", st)
		return
	}

	fmt.Printf("今天 %s 有 %d 場比賽\n\n", st, len(matchedGames))

	for i, game := range matchedGames {
		displayGame(i+1, &game, oddsMap, injuryMap)
	}

	fmt.Printf("\n執行時間：%v\n", time.Since(start))
}

// displayGame 顯示單場比賽資訊
func displayGame(index int, game *models.Game, oddsMap map[string]models.SpreadInfo, injuryMap map[string][]string) {
	homeTeam := game.HomeTeam.GetFullTeamName()
	awayTeam := game.AwayTeam.GetFullTeamName()

	gameTime, _ := crawler.ConvertUTCToLocal(game.GameTimeUTC)

	// 取得中文隊名，如果找不到就顯示原名
	homeTeamCN := models.TeamMap[homeTeam]
	if homeTeamCN == "" {
		homeTeamCN = homeTeam
	}
	awayTeamCN := models.TeamMap[awayTeam]
	if awayTeamCN == "" {
		awayTeamCN = awayTeam
	}

	// 顯示比賽基本資訊（含戰績）
	awayRecord := ""
	homeRecord := ""
	if game.AwayTeam.Wins > 0 || game.AwayTeam.Losses > 0 {
		awayRecord = fmt.Sprintf(" (%d-%d)", game.AwayTeam.Wins, game.AwayTeam.Losses)
	}
	if game.HomeTeam.Wins > 0 || game.HomeTeam.Losses > 0 {
		homeRecord = fmt.Sprintf(" (%d-%d)", game.HomeTeam.Wins, game.HomeTeam.Losses)
	}

	// 顯示比賽狀態和比分
	statusText := getGameStatusText(game)
	fmt.Printf("%d. %s%s  %s  %s%s(主)  %s\n",
		index, awayTeamCN, awayRecord, gameTime, homeTeamCN, homeRecord, statusText)

	// 如果比賽進行中或已結束，顯示比分
	if game.GameStatus >= 2 { // 2=進行中, 3=結束
		fmt.Printf("   比分：%s %d - %d %s", awayTeamCN, game.AwayTeam.Score, game.HomeTeam.Score, homeTeamCN)
		if game.GameStatus == 2 { // 進行中
			fmt.Printf("  [Q%d %s]", game.Period, game.GameClock)
		}
		fmt.Println()
	}

	// 顯示賠率
	if spread, ok := oddsMap[game.GameID]; ok && spread.Found {
		fmt.Printf("   讓分盤：開盤 主隊%.1f/客隊%.1f → 目前 主隊%s/客隊%s\n",
			spread.HomeOpeningSpread, spread.AwayOpeningSpread,
			spread.HomeSpread, spread.AwaySpread)
	} else {
		fmt.Println("   讓分盤：無資料")
	}

	// 顯示傷兵
	fmt.Println("\n  ---------------------------------")
	printInjuries(awayTeam, awayTeamCN, injuryMap)
	fmt.Println("  ---------------------------------")
	printInjuries(homeTeam, homeTeamCN, injuryMap)
	fmt.Println()
}

// getGameStatusText 取得比賽狀態文字
func getGameStatusText(game *models.Game) string {
	switch game.GameStatus {
	case 1: // 未開始
		return ""
	case 2: // 進行中
		return fmt.Sprintf("[進行中 Q%d]", game.Period)
	case 3: // 已結束
		return "[已結束]"
	default:
		return ""
	}
}

// printInjuries 顯示球隊傷兵名單
func printInjuries(teamName string, teamCN string, injuryMap map[string][]string) {
	fmt.Printf("  %s 傷兵名單\n", teamCN)

	// 處理 LA Clippers 的特殊情況
	if teamName == "Los Angeles Clippers" {
		teamName = "LA Clippers"
	}

	if injuries, ok := injuryMap[teamName]; ok && len(injuries) > 0 {
		for _, injury := range injuries {
			fmt.Printf("  %s\n", injury)
		}
	} else {
		fmt.Println("  沒有傷兵")
	}
}
