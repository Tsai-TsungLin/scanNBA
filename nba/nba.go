package nba

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type NBATeam2 struct {
	Sports []struct {
		Leagues []struct {
			Events []struct {
				Date        time.Time `json:"date"`
				Competitors []struct {
					DisplayName string `json:"displayName"`
					HomeAway    string `json:"homeAway"`
				} `json:"competitors"`
				Odds struct {
					OverUnder   float64 `json:"overUnder"`
					PointSpread struct {
						Away struct {
							Open struct {
								Line string `json:"line"`
							} `json:"open"`
							Close struct {
								Line string `json:"line"`
							} `json:"close"`
						} `json:"away"`
						Home struct {
							Open struct {
								Line string `json:"line"`
							} `json:"open"`
							Close struct {
								Line string `json:"line"`
							} `json:"close"`
						} `json:"home"`
					} `json:"pointSpread"`
					AwayTeamOdds struct {
						Favorite bool `json:"favorite"`
					} `json:"awayTeamOdds"`
					HomeTeamOdds struct {
						Favorite bool `json:"favorite"`
					} `json:"homeTeamOdds"`
				} `json:"odds"`
			} `json:"events"`
		} `json:"leagues"`
	} `json:"sports"`
}

func GetNBATeams() {
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

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
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
		return
	}

	chineseTeam := TeamInit()
	var games []string
	var dateOfGames time.Time
	matchCount := 0
	for _, sport := range nbaTeam.Sports {
		for _, league := range sport.Leagues {
			for _, event := range league.Events {
				var homeTeam, awayTeam string
				for _, competitor := range event.Competitors {
					if competitor.HomeAway == "home" {
						homeTeam = competitor.DisplayName
					} else if competitor.HomeAway == "away" {
						awayTeam = competitor.DisplayName
					}
				}

				if dateOfGames.IsZero() {
					// 設置日期為第一個事件的日期
					dateOfGames = event.Date
				}

				// 转换开赛时间为您所需的时区，这里假设为 UTC+8
				startTime := event.Date.In(time.FixedZone("UTC+8", 8*60*60))

				// 检查是否有赔率信息，并确定是哪个队伍让分
				var oddsInfo string
				if event.Odds.PointSpread.Away.Open.Line != "" {
					if event.Odds.AwayTeamOdds.Favorite {
						// 如果客队是让分的队伍
						oddsInfo = fmt.Sprintf("初盤 %s %s, 大小分: %.1f ", chineseTeam[awayTeam], event.Odds.PointSpread.Away.Open.Line, event.Odds.OverUnder)
					} else if event.Odds.HomeTeamOdds.Favorite {
						// 如果主队是让分的队伍
						oddsInfo = fmt.Sprintf("初盤 %s %s, 大小分: %.1f ", chineseTeam[homeTeam], event.Odds.PointSpread.Home.Open.Line, event.Odds.OverUnder)
					} else {
						oddsInfo = "赔率信息不明确"
					}
				} else {
					oddsInfo = "尚未開盤"
				}

				// 获取客队和主队的伤兵信息
				awayTeamInjuries := getInjury(awayTeam)
				homeTeamInjuries := getInjury(homeTeam)
				matchCount++
				// 将伤兵信息添加到比赛信息中
				gameInfo := fmt.Sprintf("%d. %s %s %s (主)\n %s\n%s 傷兵:\n%s\n---------------------------------\n%s 傷兵:\n%s\n", matchCount, chineseTeam[awayTeam], startTime.Format("15:04"), chineseTeam[homeTeam], oddsInfo, chineseTeam[awayTeam], formatInjury(awayTeamInjuries), chineseTeam[homeTeam], formatInjury(homeTeamInjuries))
				games = append(games, gameInfo)
			}
		}
	}

	fmt.Printf("今天 %s 有 %d 場比賽\n", dateOfGames.Format("2006-01-02"), len(games))
	for i, game := range games {
		fmt.Printf("%d. %s\n", i+1, game)
	}
}
