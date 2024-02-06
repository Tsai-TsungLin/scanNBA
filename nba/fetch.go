package nba

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

// GameData 結構體用於表示比賽數據
type GameData struct {
	Date    string  `json:"date"`
	Matches []Match `json:"matches"`
}

// Match 結構體用於表示一場比賽的詳細信息
type Match struct {
	AwayTeam         string              `json:"awayTeam"`
	HomeTeam         string              `json:"homeTeam"`
	Time             string              `json:"time"`
	InitialOdds      string              `json:"initialOdds"`
	CurrentOdds      string              `json:"currentOdds"`
	InitialOverUnder string              `json:"initialOverUnder"`
	CurrentOverUnder string              `json:"currentOverUnder"`
	AwayOverUnder    string              `json:"awayOverUnder"`
	HomeOverUnder    string              `json:"homeOverUnder"`
	AwayInjuries     []map[string]string `json:"awayInjuries"`
	AwayDish         []string            `json:"awayDish"`
	HomeInjuries     []map[string]string `json:"homeInjuries"`
	HomeDish         []string            `json:"homeDish"`
}

func GetNBAData() []byte {
	startTime := time.Now()

	today := time.Now()

	// 檢查時間是否已過中午12點
	var newTime string
	if today.Hour() >= 12 {
		newTime = today.Format("20060102")
	} else {
		newTime = today.AddDate(0, 0, -1).Format("20060102") // 如果沒有過中午12點，則減去一天
	}

	// 獲取數據 API URL
	url := "https://site.api.espn.com/apis/v2/scoreboard/header?sport=basketball&league=nba&dates=" + newTime + "&tz=America%2FNew_York&showAirings=buy%2Clive%2Creplay&showZipLookup=true&buyWindow=1m&lang=en&region=us&contentorigin=espn"

	log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var nbaTeam NBATeam2

	if unmarshalErr := json.Unmarshal(body, &nbaTeam); unmarshalErr != nil {
		panic(unmarshalErr)
	}

	if len(nbaTeam.Sports[0].Leagues[0].Events) == 0 {
		fmt.Println("今天", newTime, "沒有比賽")
		fmt.Println("Spend Time:", time.Since(startTime))
		return nil
	}

	chineseTeam := TeamInit()

	// var games []string
	var dateOfGames time.Time
	matchCount := 0

	var wg sync.WaitGroup

	var games GameData

	for _, sport := range nbaTeam.Sports {
		for _, league := range sport.Leagues {
			for _, event := range league.Events {
				wg.Add(1)
				go func(event Event) {
					defer wg.Done()
					var homeTeam, awayTeam string
					// 創建一個新的Match實例
					match := Match{
						AwayTeam:     chineseTeam[homeTeam],
						HomeTeam:     chineseTeam[awayTeam],
						Time:         startTime.Format("15:04"),
						InitialOdds:  event.Odds.PointSpread.Away.Open.Line,  // 假設這是初盤
						CurrentOdds:  event.Odds.PointSpread.Away.Close.Line, // 假設這是現在盤口
						AwayInjuries: nil,
						AwayDish:     nil,
						HomeInjuries: nil,
						HomeDish:     nil,
					}

					for _, competitor := range event.Competitors {
						if competitor.HomeAway == "home" {
							homeTeam = competitor.DisplayName
						} else if competitor.HomeAway == "away" {
							awayTeam = competitor.DisplayName
						}
					}

					match.AwayTeam = chineseTeam[awayTeam]
					match.HomeTeam = chineseTeam[homeTeam]

					if dateOfGames.IsZero() {
						// 設置日期為第一個事件的日期
						dateOfGames = event.Date
					}
					games.Date = dateOfGames.Format("2006-01-02")
					// 转换开赛时间为您所需的时区，这里假设为 UTC+8
					startTime := event.Date.In(time.FixedZone("UTC+8", 8*60*60))

					match.Time = startTime.Format("15:04")

					// 检查是否有赔率信息，并确定是哪个队伍让分

					if event.Odds.PointSpread.Away.Open.Line != "" {
						match.CurrentOverUnder = event.Odds.Total.Over.Close.Line[1:]
						match.InitialOverUnder = event.Odds.Total.Over.Open.Line[1:]
						if event.Odds.AwayTeamOdds.Favorite {

							match.InitialOdds = event.Odds.PointSpread.Away.Open.Line
							match.CurrentOdds = event.Odds.PointSpread.Away.Close.Line

							// 調用函數計算主客場得分
							homeScore, awayScore, err := calculateTeamScores(match.CurrentOverUnder, event.Odds.PointSpread.Away.Close.Line)
							if err != nil {
								// 處理錯誤
								log.Printf("計算得分錯誤: %v", err)
							} else {
								match.HomeOverUnder = homeScore
								match.AwayOverUnder = awayScore
							}

						} else if event.Odds.HomeTeamOdds.Favorite {

							match.InitialOdds = string([]byte(event.Odds.PointSpread.Home.Open.Line)[1:])
							match.CurrentOdds = string([]byte(event.Odds.PointSpread.Home.Close.Line)[1:])

							// 調用函數計算主客場得分
							homeScore, awayScore, err := calculateTeamScores(match.CurrentOverUnder, event.Odds.PointSpread.Home.Close.Line[1:])
							if err != nil {
								// 處理錯誤
								log.Printf("計算得分錯誤: %v", err)
							} else {
								match.HomeOverUnder = homeScore
								match.AwayOverUnder = awayScore
							}
						} else {
							match.CurrentOverUnder = "赔率信息不明确"
							match.InitialOverUnder = "赔率信息不明确"
							match.InitialOdds = "赔率信息不明确"
							match.CurrentOdds = "赔率信息不明确"
						}
					} else {
						match.CurrentOverUnder = "尚未開盤"
						match.InitialOverUnder = "尚未開盤"
						match.InitialOdds = "尚未開盤"
						match.CurrentOdds = "尚未開盤"
					}

					innerWg := &sync.WaitGroup{} // 內部WaitGroup用於等待隊伍傷兵信息

					// 為客隊和主隊的傷兵信息開啟goroutine
					innerWg.Add(2)
					go func() {
						defer innerWg.Done()
						awayTeamInjuries := getInjuryFetch(awayTeam)
						awayDish := getDishFetch(awayTeam)
						match.AwayInjuries = awayTeamInjuries
						match.AwayDish = awayDish
					}()

					go func() {
						defer innerWg.Done()
						homeTeamInjuries := getInjuryFetch(homeTeam)
						homeDish := getDishFetch(homeTeam)
						match.HomeInjuries = homeTeamInjuries
						match.HomeDish = homeDish
					}()

					innerWg.Wait() // 等待內部的兩個goroutine完成

					matchCount++

					games.Matches = append(games.Matches, match)

				}(event)
			}
		}
	}

	wg.Wait()

	sort.Slice(games.Matches, func(i, j int) bool {
		if games.Matches[i].Time == games.Matches[j].Time {
			// 如果時間相同，則根據隊伍名稱排序
			if games.Matches[i].AwayTeam == games.Matches[j].AwayTeam {
				return games.Matches[i].HomeTeam < games.Matches[j].HomeTeam
			}
			return games.Matches[i].AwayTeam < games.Matches[j].AwayTeam
		}
		// 首先根據時間排序
		return games.Matches[i].Time < games.Matches[j].Time
	})

	jsonData, err := json.Marshal(games)
	if err != nil {
		log.Println("json Marshal error:", err)
	}

	return jsonData
}
