package old

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Schedule struct {
	Context struct {
		User struct {
			CountryCode  string `json:"countryCode"`
			CountryName  string `json:"countryName"`
			Locale       string `json:"locale"`
			TimeZone     string `json:"timeZone"`
			TimeZoneCity string `json:"timeZoneCity"`
		} `json:"user"`
		Device struct {
			Clazz interface{} `json:"clazz"`
		} `json:"device"`
	} `json:"context"`
	Error struct {
		Detail  interface{} `json:"detail"`
		IsError string      `json:"isError"`
		Message interface{} `json:"message"`
	} `json:"error"`
	Payload struct {
		League struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"league"`
		Season struct {
			IsCurrent               string `json:"isCurrent"`
			RosterSeasonType        int    `json:"rosterSeasonType"`
			RosterSeasonYear        string `json:"rosterSeasonYear"`
			RosterSeasonYearDisplay string `json:"rosterSeasonYearDisplay"`
			ScheduleSeasonType      int    `json:"scheduleSeasonType"`
			ScheduleSeasonYear      string `json:"scheduleSeasonYear"`
			ScheduleYearDisplay     string `json:"scheduleYearDisplay"`
			StatsSeasonType         int    `json:"statsSeasonType"`
			StatsSeasonYear         string `json:"statsSeasonYear"`
			StatsSeasonYearDisplay  string `json:"statsSeasonYearDisplay"`
			Year                    string `json:"year"`
			YearDisplay             string `json:"yearDisplay"`
		} `json:"season"`
		Date struct {
			Games []struct {
				Profile struct {
					ArenaLocation string      `json:"arenaLocation"`
					ArenaName     string      `json:"arenaName"`
					AwayTeamID    string      `json:"awayTeamId"`
					DateTimeEt    string      `json:"dateTimeEt"`
					GameID        string      `json:"gameId"`
					HomeTeamID    string      `json:"homeTeamId"`
					Number        string      `json:"number"`
					ScheduleCode  interface{} `json:"scheduleCode"`
					SeasonType    string      `json:"seasonType"`
					Sequence      string      `json:"sequence"`
					UtcMillis     string      `json:"utcMillis"`
				} `json:"profile"`
				Boxscore struct {
					Attendance            string      `json:"attendance"`
					AwayScore             int         `json:"awayScore"`
					GameLength            interface{} `json:"gameLength"`
					HomeScore             int         `json:"homeScore"`
					LeadChanges           interface{} `json:"leadChanges"`
					OfficialsDisplayName1 interface{} `json:"officialsDisplayName1"`
					OfficialsDisplayName2 interface{} `json:"officialsDisplayName2"`
					OfficialsDisplayName3 interface{} `json:"officialsDisplayName3"`
					Period                string      `json:"period"`
					PeriodClock           interface{} `json:"periodClock"`
					Status                string      `json:"status"`
					StatusDesc            interface{} `json:"statusDesc"`
					Ties                  interface{} `json:"ties"`
				} `json:"boxscore"`
				Urls         []interface{} `json:"urls"`
				Broadcasters []interface{} `json:"broadcasters"`
				HomeTeam     struct {
					Profile struct {
						Abbr              string `json:"abbr"`
						City              string `json:"city"`
						CityEn            string `json:"cityEn"`
						Code              string `json:"code"`
						Conference        string `json:"conference"`
						DisplayAbbr       string `json:"displayAbbr"`
						DisplayConference string `json:"displayConference"`
						Division          string `json:"division"`
						ID                string `json:"id"`
						IsAllStarTeam     bool   `json:"isAllStarTeam"`
						IsLeagueTeam      bool   `json:"isLeagueTeam"`
						LeagueID          string `json:"leagueId"`
						Name              string `json:"name"`
						NameEn            string `json:"nameEn"`
					} `json:"profile"`
					Matchup struct {
						ConfRank   string      `json:"confRank"`
						DivRank    string      `json:"divRank"`
						Losses     string      `json:"losses"`
						SeriesText interface{} `json:"seriesText"`
						Wins       string      `json:"wins"`
					} `json:"matchup"`
					Score struct {
						Assists                int     `json:"assists"`
						BiggestLead            int     `json:"biggestLead"`
						Blocks                 int     `json:"blocks"`
						BlocksAgainst          int     `json:"blocksAgainst"`
						DefRebs                int     `json:"defRebs"`
						Disqualifications      int     `json:"disqualifications"`
						Ejections              int     `json:"ejections"`
						FastBreakPoints        int     `json:"fastBreakPoints"`
						Fga                    int     `json:"fga"`
						Fgm                    int     `json:"fgm"`
						Fgpct                  float64 `json:"fgpct"`
						FlagrantFouls          int     `json:"flagrantFouls"`
						Fouls                  int     `json:"fouls"`
						Fta                    int     `json:"fta"`
						Ftm                    int     `json:"ftm"`
						Ftpct                  float64 `json:"ftpct"`
						FullTimeoutsRemaining  int     `json:"fullTimeoutsRemaining"`
						Mins                   int     `json:"mins"`
						OffRebs                int     `json:"offRebs"`
						Ot10Score              int     `json:"ot10Score"`
						Ot1Score               int     `json:"ot1Score"`
						Ot2Score               int     `json:"ot2Score"`
						Ot3Score               int     `json:"ot3Score"`
						Ot4Score               int     `json:"ot4Score"`
						Ot5Score               int     `json:"ot5Score"`
						Ot6Score               int     `json:"ot6Score"`
						Ot7Score               int     `json:"ot7Score"`
						Ot8Score               int     `json:"ot8Score"`
						Ot9Score               int     `json:"ot9Score"`
						PointsInPaint          int     `json:"pointsInPaint"`
						PointsOffTurnovers     int     `json:"pointsOffTurnovers"`
						Q1Score                int     `json:"q1Score"`
						Q2Score                int     `json:"q2Score"`
						Q3Score                int     `json:"q3Score"`
						Q4Score                int     `json:"q4Score"`
						Rebs                   int     `json:"rebs"`
						Score                  int     `json:"score"`
						Seconds                int     `json:"seconds"`
						ShortTimeoutsRemaining int     `json:"shortTimeoutsRemaining"`
						Steals                 int     `json:"steals"`
						TechnicalFouls         int     `json:"technicalFouls"`
						Tpa                    int     `json:"tpa"`
						Tpm                    int     `json:"tpm"`
						Tppct                  float64 `json:"tppct"`
						Turnovers              int     `json:"turnovers"`
					} `json:"score"`
					PointGameLeader   interface{} `json:"pointGameLeader"`
					AssistGameLeader  interface{} `json:"assistGameLeader"`
					ReboundGameLeader interface{} `json:"reboundGameLeader"`
				} `json:"homeTeam"`
				AwayTeam struct {
					Profile struct {
						Abbr              string `json:"abbr"`
						City              string `json:"city"`
						CityEn            string `json:"cityEn"`
						Code              string `json:"code"`
						Conference        string `json:"conference"`
						DisplayAbbr       string `json:"displayAbbr"`
						DisplayConference string `json:"displayConference"`
						Division          string `json:"division"`
						ID                string `json:"id"`
						IsAllStarTeam     bool   `json:"isAllStarTeam"`
						IsLeagueTeam      bool   `json:"isLeagueTeam"`
						LeagueID          string `json:"leagueId"`
						Name              string `json:"name"`
						NameEn            string `json:"nameEn"`
					} `json:"profile"`
					Matchup struct {
						ConfRank   string      `json:"confRank"`
						DivRank    string      `json:"divRank"`
						Losses     string      `json:"losses"`
						SeriesText interface{} `json:"seriesText"`
						Wins       string      `json:"wins"`
					} `json:"matchup"`
					Score struct {
						Assists                int     `json:"assists"`
						BiggestLead            int     `json:"biggestLead"`
						Blocks                 int     `json:"blocks"`
						BlocksAgainst          int     `json:"blocksAgainst"`
						DefRebs                int     `json:"defRebs"`
						Disqualifications      int     `json:"disqualifications"`
						Ejections              int     `json:"ejections"`
						FastBreakPoints        int     `json:"fastBreakPoints"`
						Fga                    int     `json:"fga"`
						Fgm                    int     `json:"fgm"`
						Fgpct                  float64 `json:"fgpct"`
						FlagrantFouls          int     `json:"flagrantFouls"`
						Fouls                  int     `json:"fouls"`
						Fta                    int     `json:"fta"`
						Ftm                    int     `json:"ftm"`
						Ftpct                  float64 `json:"ftpct"`
						FullTimeoutsRemaining  int     `json:"fullTimeoutsRemaining"`
						Mins                   int     `json:"mins"`
						OffRebs                int     `json:"offRebs"`
						Ot10Score              int     `json:"ot10Score"`
						Ot1Score               int     `json:"ot1Score"`
						Ot2Score               int     `json:"ot2Score"`
						Ot3Score               int     `json:"ot3Score"`
						Ot4Score               int     `json:"ot4Score"`
						Ot5Score               int     `json:"ot5Score"`
						Ot6Score               int     `json:"ot6Score"`
						Ot7Score               int     `json:"ot7Score"`
						Ot8Score               int     `json:"ot8Score"`
						Ot9Score               int     `json:"ot9Score"`
						PointsInPaint          int     `json:"pointsInPaint"`
						PointsOffTurnovers     int     `json:"pointsOffTurnovers"`
						Q1Score                int     `json:"q1Score"`
						Q2Score                int     `json:"q2Score"`
						Q3Score                int     `json:"q3Score"`
						Q4Score                int     `json:"q4Score"`
						Rebs                   int     `json:"rebs"`
						Score                  int     `json:"score"`
						Seconds                int     `json:"seconds"`
						ShortTimeoutsRemaining int     `json:"shortTimeoutsRemaining"`
						Steals                 int     `json:"steals"`
						TechnicalFouls         int     `json:"technicalFouls"`
						Tpa                    int     `json:"tpa"`
						Tpm                    int     `json:"tpm"`
						Tppct                  float64 `json:"tppct"`
						Turnovers              int     `json:"turnovers"`
					} `json:"score"`
					PointGameLeader   interface{} `json:"pointGameLeader"`
					AssistGameLeader  interface{} `json:"assistGameLeader"`
					ReboundGameLeader interface{} `json:"reboundGameLeader"`
				} `json:"awayTeam"`
				IfNecessary bool        `json:"ifNecessary"`
				SeriesText  interface{} `json:"seriesText"`
			} `json:"games"`
			DateMillis string `json:"dateMillis"`
			GameCount  string `json:"gameCount"`
		} `json:"date"`
		NextAvailableDateMillis string `json:"nextAvailableDateMillis"`
		UtcMillis               string `json:"utcMillis"`
	} `json:"payload"`
	Timestamp string `json:"timestamp"`
}

<<<<<<< HEAD:nba/old/old_nba.go
type Game struct {
	Profile struct {
		ArenaLocation string      `json:"arenaLocation"`
		ArenaName     string      `json:"arenaName"`
		AwayTeamID    string      `json:"awayTeamId"`
		DateTimeEt    string      `json:"dateTimeEt"`
		GameID        string      `json:"gameId"`
		HomeTeamID    string      `json:"homeTeamId"`
		Number        string      `json:"number"`
		ScheduleCode  interface{} `json:"scheduleCode"`
		SeasonType    string      `json:"seasonType"`
		Sequence      string      `json:"sequence"`
		UtcMillis     string      `json:"utcMillis"`
	} `json:"profile"`
	Boxscore struct {
		Attendance            string      `json:"attendance"`
		AwayScore             int         `json:"awayScore"`
		GameLength            interface{} `json:"gameLength"`
		HomeScore             int         `json:"homeScore"`
		LeadChanges           interface{} `json:"leadChanges"`
		OfficialsDisplayName1 interface{} `json:"officialsDisplayName1"`
		OfficialsDisplayName2 interface{} `json:"officialsDisplayName2"`
		OfficialsDisplayName3 interface{} `json:"officialsDisplayName3"`
		Period                string      `json:"period"`
		PeriodClock           interface{} `json:"periodClock"`
		Status                string      `json:"status"`
		StatusDesc            interface{} `json:"statusDesc"`
		Ties                  interface{} `json:"ties"`
	} `json:"boxscore"`
	Urls         []interface{} `json:"urls"`
	Broadcasters []interface{} `json:"broadcasters"`
	HomeTeam     struct {
		Profile struct {
			Abbr              string `json:"abbr"`
			City              string `json:"city"`
			CityEn            string `json:"cityEn"`
			Code              string `json:"code"`
			Conference        string `json:"conference"`
			DisplayAbbr       string `json:"displayAbbr"`
			DisplayConference string `json:"displayConference"`
			Division          string `json:"division"`
			ID                string `json:"id"`
			IsAllStarTeam     bool   `json:"isAllStarTeam"`
			IsLeagueTeam      bool   `json:"isLeagueTeam"`
			LeagueID          string `json:"leagueId"`
			Name              string `json:"name"`
			NameEn            string `json:"nameEn"`
		} `json:"profile"`
		Matchup struct {
			ConfRank   string      `json:"confRank"`
			DivRank    string      `json:"divRank"`
			Losses     string      `json:"losses"`
			SeriesText interface{} `json:"seriesText"`
			Wins       string      `json:"wins"`
		} `json:"matchup"`
		Score struct {
			Assists                int     `json:"assists"`
			BiggestLead            int     `json:"biggestLead"`
			Blocks                 int     `json:"blocks"`
			BlocksAgainst          int     `json:"blocksAgainst"`
			DefRebs                int     `json:"defRebs"`
			Disqualifications      int     `json:"disqualifications"`
			Ejections              int     `json:"ejections"`
			FastBreakPoints        int     `json:"fastBreakPoints"`
			Fga                    int     `json:"fga"`
			Fgm                    int     `json:"fgm"`
			Fgpct                  float64 `json:"fgpct"`
			FlagrantFouls          int     `json:"flagrantFouls"`
			Fouls                  int     `json:"fouls"`
			Fta                    int     `json:"fta"`
			Ftm                    int     `json:"ftm"`
			Ftpct                  float64 `json:"ftpct"`
			FullTimeoutsRemaining  int     `json:"fullTimeoutsRemaining"`
			Mins                   int     `json:"mins"`
			OffRebs                int     `json:"offRebs"`
			Ot10Score              int     `json:"ot10Score"`
			Ot1Score               int     `json:"ot1Score"`
			Ot2Score               int     `json:"ot2Score"`
			Ot3Score               int     `json:"ot3Score"`
			Ot4Score               int     `json:"ot4Score"`
			Ot5Score               int     `json:"ot5Score"`
			Ot6Score               int     `json:"ot6Score"`
			Ot7Score               int     `json:"ot7Score"`
			Ot8Score               int     `json:"ot8Score"`
			Ot9Score               int     `json:"ot9Score"`
			PointsInPaint          int     `json:"pointsInPaint"`
			PointsOffTurnovers     int     `json:"pointsOffTurnovers"`
			Q1Score                int     `json:"q1Score"`
			Q2Score                int     `json:"q2Score"`
			Q3Score                int     `json:"q3Score"`
			Q4Score                int     `json:"q4Score"`
			Rebs                   int     `json:"rebs"`
			Score                  int     `json:"score"`
			Seconds                int     `json:"seconds"`
			ShortTimeoutsRemaining int     `json:"shortTimeoutsRemaining"`
			Steals                 int     `json:"steals"`
			TechnicalFouls         int     `json:"technicalFouls"`
			Tpa                    int     `json:"tpa"`
			Tpm                    int     `json:"tpm"`
			Tppct                  float64 `json:"tppct"`
			Turnovers              int     `json:"turnovers"`
		} `json:"score"`
		PointGameLeader   interface{} `json:"pointGameLeader"`
		AssistGameLeader  interface{} `json:"assistGameLeader"`
		ReboundGameLeader interface{} `json:"reboundGameLeader"`
	} `json:"homeTeam"`
	AwayTeam struct {
		Profile struct {
			Abbr              string `json:"abbr"`
			City              string `json:"city"`
			CityEn            string `json:"cityEn"`
			Code              string `json:"code"`
			Conference        string `json:"conference"`
			DisplayAbbr       string `json:"displayAbbr"`
			DisplayConference string `json:"displayConference"`
			Division          string `json:"division"`
			ID                string `json:"id"`
			IsAllStarTeam     bool   `json:"isAllStarTeam"`
			IsLeagueTeam      bool   `json:"isLeagueTeam"`
			LeagueID          string `json:"leagueId"`
			Name              string `json:"name"`
			NameEn            string `json:"nameEn"`
		} `json:"profile"`
		Matchup struct {
			ConfRank   string      `json:"confRank"`
			DivRank    string      `json:"divRank"`
			Losses     string      `json:"losses"`
			SeriesText interface{} `json:"seriesText"`
			Wins       string      `json:"wins"`
		} `json:"matchup"`
		Score struct {
			Assists                int     `json:"assists"`
			BiggestLead            int     `json:"biggestLead"`
			Blocks                 int     `json:"blocks"`
			BlocksAgainst          int     `json:"blocksAgainst"`
			DefRebs                int     `json:"defRebs"`
			Disqualifications      int     `json:"disqualifications"`
			Ejections              int     `json:"ejections"`
			FastBreakPoints        int     `json:"fastBreakPoints"`
			Fga                    int     `json:"fga"`
			Fgm                    int     `json:"fgm"`
			Fgpct                  float64 `json:"fgpct"`
			FlagrantFouls          int     `json:"flagrantFouls"`
			Fouls                  int     `json:"fouls"`
			Fta                    int     `json:"fta"`
			Ftm                    int     `json:"ftm"`
			Ftpct                  float64 `json:"ftpct"`
			FullTimeoutsRemaining  int     `json:"fullTimeoutsRemaining"`
			Mins                   int     `json:"mins"`
			OffRebs                int     `json:"offRebs"`
			Ot10Score              int     `json:"ot10Score"`
			Ot1Score               int     `json:"ot1Score"`
			Ot2Score               int     `json:"ot2Score"`
			Ot3Score               int     `json:"ot3Score"`
			Ot4Score               int     `json:"ot4Score"`
			Ot5Score               int     `json:"ot5Score"`
			Ot6Score               int     `json:"ot6Score"`
			Ot7Score               int     `json:"ot7Score"`
			Ot8Score               int     `json:"ot8Score"`
			Ot9Score               int     `json:"ot9Score"`
			PointsInPaint          int     `json:"pointsInPaint"`
			PointsOffTurnovers     int     `json:"pointsOffTurnovers"`
			Q1Score                int     `json:"q1Score"`
			Q2Score                int     `json:"q2Score"`
			Q3Score                int     `json:"q3Score"`
			Q4Score                int     `json:"q4Score"`
			Rebs                   int     `json:"rebs"`
			Score                  int     `json:"score"`
			Seconds                int     `json:"seconds"`
			ShortTimeoutsRemaining int     `json:"shortTimeoutsRemaining"`
			Steals                 int     `json:"steals"`
			TechnicalFouls         int     `json:"technicalFouls"`
			Tpa                    int     `json:"tpa"`
			Tpm                    int     `json:"tpm"`
			Tppct                  float64 `json:"tppct"`
			Turnovers              int     `json:"turnovers"`
		} `json:"score"`
		PointGameLeader   interface{} `json:"pointGameLeader"`
		AssistGameLeader  interface{} `json:"assistGameLeader"`
		ReboundGameLeader interface{} `json:"reboundGameLeader"`
	} `json:"awayTeam"`
	IfNecessary bool        `json:"ifNecessary"`
	SeriesText  interface{} `json:"seriesText"`
}

// GameResult 用于存储和排序比赛信息// GameInfo 用于存储和排序比赛信息
type GameInfo struct {
	Description string
	GameTime    time.Time
}

=======
>>>>>>> main:nba.go
// 把隊伍存進map
func TeamInit() map[string]string {

	teamMap := make(map[string]string)
	teamMap["Atlanta Hawks"] = "老鷹"
	teamMap["Boston Celtics"] = "提克"
	teamMap["Brooklyn Nets"] = "籃網"
	teamMap["Cleveland Cavaliers"] = "騎士"
	teamMap["Charlotte Hornets"] = "黃蜂"
	teamMap["Chicago Bulls"] = "公牛"
	teamMap["Dallas Mavericks"] = "小牛"
	teamMap["Denver Nuggets"] = "金塊"
	teamMap["Detroit Pistons"] = "活塞"
	teamMap["Golden State Warriors"] = "勇士"
	teamMap["Houston Rockets"] = "火箭"
	teamMap["Indiana Pacers"] = "溜馬"
	teamMap["Los Angeles Lakers"] = "湖人"
	teamMap["LA Clippers"] = "快艇"
	teamMap["Memphis Grizzlies"] = "灰熊"
	teamMap["Miami Heat"] = "熱火"
	teamMap["Milwaukee Bucks"] = "公鹿"
	teamMap["Minnesota Timberwolves"] = "灰狼"
	teamMap["New Orleans Pelicans"] = "鵜鶘"
	teamMap["New York Knicks"] = "尼克"
	teamMap["Oklahoma City Thunder"] = "雷霆"
	teamMap["Orlando Magic"] = "魔術"
	teamMap["Philadelphia 76ers"] = "76人"
	teamMap["Phoenix Suns"] = "太陽"
	teamMap["Portland Trail Blazers"] = "拓荒"
	teamMap["Sacramento Kings"] = "國王"
	teamMap["San Antonio Spurs"] = "馬刺"
	teamMap["Toronto Raptors"] = "暴龍"
	teamMap["Utah Jazz"] = "爵士"
	teamMap["Washington Wizards"] = "巫師"

	return teamMap
}

// Get the injuriers of the nba team
func getInjury(searchTeam string) (result []string) {

	if searchTeam == "Los Angeles Clippers" {
		searchTeam = "LA Clippers"
	}

	if searchTeam == "Los Angeles Clippers" {
		searchTeam = "LA Clippers"
	}

	if searchTeam == "Portland Trail_blazers" {
		searchTeam = "Portland Trail Blazers"
	}

	res, err := http.Get("https://www.espn.com/nba/injuries")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".Table__league-injuries").Each(func(i int, s *goquery.Selection) {
		team := s.Find(".injuries__teamName").Text()

		if strings.EqualFold(team, searchTeam) {

			s.Find(".Table__even").Each((func(i int, g *goquery.Selection) {
				name := g.Find(".AnchorLink").Text()
				status := g.Find(".col-stat").Text()
				comment := g.Find(".col-desc").Text()

				if status != "" {
					name = fmt.Sprintf("%-25s", name)
					status = fmt.Sprintf("%-15s", status)
					comment = fmt.Sprintf("%-5s", comment)
					injury := name + status + comment
					result = append(result, injury)
				}
			}))
		}
	})

	return result

}

// 取得Injury的comment
func GetInjuryComment(searchTeam string) (result string) {

	foundInjuries := false

	if searchTeam == "Los Angeles Clippers" {
		searchTeam = "LA Clippers"
	}

	if searchTeam == "Portland Trail_blazers" {
		searchTeam = "Portland Trail Blazers"
	}

	res, err := http.Get("https://www.espn.com/nba/injuries")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	teamMap := TeamInit()

	doc.Find(".Table__league-injuries").Each(func(i int, s *goquery.Selection) {
		team := s.Find(".injuries__teamName").Text()

		if strings.EqualFold(team, searchTeam) {

			result = result + teamMap[searchTeam] + "--"

			s.Find(".Table__even").Each((func(ii int, g *goquery.Selection) {
				name := g.Find(".AnchorLink").Text()
				status := g.Find(".col-stat").Text()
				comment := g.Find(".col-desc").Text()

				if status != "" {
					foundInjuries = true
					commentResult := sortComment(name, comment)
					result = result + commentResult

				}

			}))
		}
	})

	if !foundInjuries {
		result = result + teamMap[searchTeam] + "-全陣容 /"
	}

	return result
}

// comment 分類
func sortComment(name, comment string) (result string) {
	switch {
	case strings.Contains(comment, "out"):
		result = result + name + "不上/"
	case strings.Contains(comment, "miss"):
		result = result + name + "不上/"
	case strings.Contains(comment, "will not play"):
		result = result + name + "不上/"
	case strings.Contains(comment, "won't"):
		result = result + name + "不上/"
	case strings.Contains(comment, "question"):
		result = result + name + "可能不上/"
	case strings.Contains(comment, "doubtful"):
		result = result + name + "可能不上/"
	case strings.Contains(comment, "day-to-day"):
		result = result + name + "不確定會不會上/"
	case strings.Contains(comment, "probable"):
		result = result + name + "可能會上/"
	default:
		result = result + name + "不確定會不會上/"
	}

	return result
}

<<<<<<< HEAD:nba/old/old_nba.go
=======
// Get the nba game of the day
>>>>>>> main:nba.go
func PKTeam() {
	startTime := time.Now()

	// 轉換時區到 UTC-4
	zone := time.FixedZone("", -4*60*60)
	today := time.Now()
	newTime := today.In(zone).Format("2006-01-02")

	// 獲取數據 API URL
	url := "https://in.global.nba.com/stats2/scores/daily.json?gameDate=" + newTime + "&locale=en&tz=%2B8&countryCode=TW#"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var result Schedule
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if unmarshalErr := json.Unmarshal(body, &result); unmarshalErr != nil {
		panic(unmarshalErr)
	}

	if result.Payload.Date.GameCount == "" {
		fmt.Println("今天", newTime, "沒有比賽")
		fmt.Println("Spend Time:", time.Since(startTime))
		return
	}

	var wg sync.WaitGroup
	gameInfoChan := make(chan GameInfo, len(result.Payload.Date.Games))
	commentChan := make(chan string, len(result.Payload.Date.Games))

	for _, game := range result.Payload.Date.Games {
		wg.Add(1)
		go func(g Game) {
			defer wg.Done()
			gameInfo := processGame(g)
			gameInfoChan <- gameInfo

			// 獲取傷病評論
			comment := GetInjuryComment(g.HomeTeam.Profile.City + " " + g.HomeTeam.Profile.Name)
			comment += " / " + GetInjuryComment(g.AwayTeam.Profile.City+" "+g.AwayTeam.Profile.Name)
			commentChan <- comment
		}(game)
	}

	wg.Wait()
	close(gameInfoChan)
	close(commentChan)

	var gameInfos []GameInfo
	for gi := range gameInfoChan {
		gameInfos = append(gameInfos, gi)
	}

	// 按比賽時間排序
	sort.Slice(gameInfos, func(i, j int) bool {
		return gameInfos[i].GameTime.Before(gameInfos[j].GameTime)
	})

	// 打印排序後的比賽信息
	fmt.Printf("今天 %s 有 %s 場比賽 \n", newTime, result.Payload.Date.GameCount)
	for _, gi := range gameInfos {
		fmt.Println(gi.Description)
	}

	var finalComments string
	for comment := range commentChan {
		finalComments += comment + "\n\n"
	}
	fmt.Println(finalComments)

	fmt.Println("Spend Time:", time.Since(startTime))
}

<<<<<<< HEAD:nba/old/old_nba.go
func processGame(game Game) GameInfo {
	// 從 game 中提取比賽時間和球隊信息
	layout := "2006-01-02T15:04"
	gameTime, _ := time.Parse(layout, game.Profile.DateTimeEt)

	gameTime = gameTime.Add(time.Hour * 13)

	// 假設的 getInjury 和 getDish 函數
	AwayTeamInjury := getInjury(game.AwayTeam.Profile.City + " " + game.AwayTeam.Profile.Name)
	HomeTeamInjury := getInjury(game.HomeTeam.Profile.City + " " + game.HomeTeam.Profile.Name)
	AwayTeamDish := getDish(game.AwayTeam.Profile.City + " " + game.AwayTeam.Profile.Name)
	HomeTeamDish := getDish(game.HomeTeam.Profile.City + " " + game.HomeTeam.Profile.Name)

	description := fmt.Sprintf(
		"%s %s vs %s\n\n  ---------------------------------\n  %s injury 名單 \n%s\n\n  ---------------------------------\n  %s injury 名單 \n%s\n\n  %s 近期過盤狀況: %s\n  %s 近期過盤狀況: %s\n",
		gameTime.Format("2006-01-02 15:04"),
		game.AwayTeam.Profile.City+" "+game.AwayTeam.Profile.Name,
		game.HomeTeam.Profile.City+" "+game.HomeTeam.Profile.Name,
		game.AwayTeam.Profile.City+" "+game.AwayTeam.Profile.Name,
		formatInjury(AwayTeamInjury),
		game.HomeTeam.Profile.City+" "+game.HomeTeam.Profile.Name,
		formatInjury(HomeTeamInjury),
		game.AwayTeam.Profile.City+" "+game.AwayTeam.Profile.Name,
		formatDish(AwayTeamDish, game.AwayTeam.Profile.City+" "+game.AwayTeam.Profile.Name),
		game.HomeTeam.Profile.City+" "+game.HomeTeam.Profile.Name,
		formatDish(HomeTeamDish, game.HomeTeam.Profile.City+" "+game.HomeTeam.Profile.Name))

	return GameInfo{
		Description: description,
		GameTime:    gameTime,
	}
}

// formatInjury 格式化傷病信息
func formatInjury(injuries []string) string {
	if len(injuries) == 0 {
		return "  沒有傷兵\n"
	}
	formatted := "  "
	for _, injury := range injuries {
		formatted += injury + "\n  "
	}
	return formatted
}

=======
>>>>>>> main:nba.go
// Get the nba game of the day , search by "StartTime"
func PKTeamOnStartTime(st string) {
	if len(st) != 5 {
		fmt.Println("你輸入的時間" + st + "格式錯誤，eg: '11:00'(請用半形) ")
		return
	}

	_, errCheckTime := time.Parse("15:04", st)
	if errCheckTime != nil {
		fmt.Println("時間格式錯誤: ", errCheckTime)
		return
	}

	startTime := time.Now()
	//轉成UTC-4
	zone := time.FixedZone("", -4*60*60)
	today := time.Now()
	newTime := today.In(zone).Format("2006-01-02")

	//get data api url
	url := "https://in.global.nba.com/stats2/scores/daily.json?gameDate=" + newTime + "&locale=en&tz=%2B8&countryCode=TW#"
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	var result Schedule

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if unmarshalErr := json.Unmarshal(body, &result); unmarshalErr != nil {
		panic(unmarshalErr)
	}

	var msg string
	var layout string = "2006-01-02T15:04"
	var count int = 0 //計算是否有隊伍

	for _, v := range result.Payload.Date.Games {

		t, _ := time.Parse(layout, v.Profile.DateTimeEt)

		if v.AwayTeam.Profile.Name == "Clippers" {
			v.AwayTeam.Profile.City = "Los Angeles"
		}

		if v.HomeTeam.Profile.Name == "Clippers" {
			v.HomeTeam.Profile.City = "Los Angeles"
		}

		//若st == 開打時間
		if st == t.Add(time.Hour*13).Format("15:04") {
			count++

			AwayTeam := v.AwayTeam.Profile.City + " " + v.AwayTeam.Profile.Name
			HomeTeam := v.HomeTeam.Profile.City + " " + v.HomeTeam.Profile.Name
			AwayTeamInjury := getInjury(AwayTeam)
			HomeTeamInjury := getInjury(HomeTeam)
			msg = msg + fmt.Sprint(count) + ". " + AwayTeam + "  " + t.Add(time.Hour*13).Format("15:04") + "  " + HomeTeam + "(主)  " + "\n"

			msg = msg + "\n  ---------------------------------\n"
			msg = msg + "  " + AwayTeam + " injury 名單 \n"

			if len(AwayTeamInjury) == 0 {
				msg = msg + "  沒有傷兵\n\n"
			} else {
				for _, v := range AwayTeamInjury {
					msg = msg + "  " + v + "\n"
				}

			}

			msg = msg + "\n  ---------------------------------\n"

			msg = msg + "  " + HomeTeam + " injury 名單 \n"
			if len(HomeTeamInjury) == 0 {
				msg = msg + "  沒有傷兵\n\n"
			} else {
				for _, v := range HomeTeamInjury {
					msg = msg + "  " + v + "\n"
				}
				msg = msg + "\n"
			}
			AwayTeamDish := getDish(AwayTeam)
			HomeTeamDish := getDish(HomeTeam)
			msg = msg + "  " + AwayTeam + " 近期過盤狀況: " + AwayTeamDish[AwayTeam] + "\n"
			msg = msg + "  " + HomeTeam + " 近期過盤狀況: " + HomeTeamDish[HomeTeam] + "\n\n"
		}

	}

	if count == 0 {
		msg = msg + "今天 " + st + " 沒有比賽\n"
	} else {
		msg = "今天 " + st + " 有 " + fmt.Sprint(count) + " 場比賽 \n" + msg
	}

	fmt.Println(msg)
	fmt.Println("Spend Time:", time.Since(startTime))
}

// 取得近五場輸贏盤口資訊
func getDish(searchTeam string) map[string]string {
	var (
		season string
	)

	year := time.Now()
	month := time.Now().Format("01")
	monthR, _ := strconv.ParseInt(month, 10, 64)
	if monthR < 7 {
		season = year.AddDate(-1, 0, 0).Format("06") + "-" + year.Format("06")
	} else {
		season = year.Format("06") + "-" + year.AddDate(1, 0, 0).Format("06")
	}
	yearR := year.Format("2006010215")

	// http://nba.titan007.com/cn/LetGoal.aspx?SclassID=1&matchSeason=2022-2023
	url := "https://nba.titan007.com/jsData/letGoal/" + season + "/l1.js?version=" + yearR
	log.Println(url)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var Team string

	r1, _ := regexp.Compile(";")
	r2 := r1.FindAllStringSubmatchIndex(string(sitemap), -1)

	Team = (string([]byte(sitemap[r2[0][0]+17 : r2[1][0]-1])))
	frontside, _ := regexp.Compile(`\[(.*?)]`)

	frontsideTeam := frontside.FindAllStringSubmatchIndex(string(Team), -1)

	matchMap := make(map[int64]string)
	result := make(map[string]string)

	for i := range frontsideTeam {
		//把隊伍資料放進map，稍後用在拿勝率
		TeamData := (Team[frontsideTeam[i][0]:frontsideTeam[i][1]])
		match, _ := regexp.Compile(",")
		matchR := match.FindAllStringSubmatchIndex(TeamData, -1)
		TeamNumber, _ := strconv.ParseInt(TeamData[1:matchR[0][0]], 10, 64) //TeamData[1:TeamData[matchR[0][0]-1]]

		if searchTeam == TeamData[matchR[2][1]+1:matchR[3][0]-1] {
			matchMap[TeamNumber] = TeamData[matchR[2][1]+1 : matchR[3][0]-1]
		}
	}

	data := (string([]byte(sitemap[r2[1][0]+20 : r2[2][0]-1])))
	frontsidedata := frontside.FindAllStringSubmatchIndex(string(data), -1)

	//拿取勝率等資料
	for i := range frontsidedata {
		winPercentData := (data[frontsidedata[i][0]:frontsidedata[i][1]])
		match, _ := regexp.Compile(",")
		matchR := match.FindAllStringSubmatchIndex(winPercentData, -1)
		TeamNumber, _ := strconv.ParseInt(winPercentData[matchR[0][0]+1:matchR[1][0]], 10, 64)

		//比對Map內是否有這個TeamNumber，有的話加入result map內回傳
		if _, ok := matchMap[TeamNumber]; ok {

			one := changeWinLose(winPercentData[matchR[11][0]+1 : matchR[12][1]-1])
			two := changeWinLose(winPercentData[matchR[12][0]+1 : matchR[13][1]-1])
			three := changeWinLose(winPercentData[matchR[13][0]+1 : matchR[14][1]-1])
			four := changeWinLose(winPercentData[matchR[14][0]+1 : matchR[15][1]-1])
			five := changeWinLose(winPercentData[matchR[15][0]+1 : matchR[15][1]+1])

			result[matchMap[TeamNumber]] = one + "," + two + "," + three + "," + four + "," + five
		}
	}

	return result

}

func changeWinLose(r string) (result string) {
	if r == "0" {
		result = "贏"
	} else {
		result = "輸"
	}
	return
}

<<<<<<< HEAD:nba/old/old_nba.go
// formatDish 格式化過盤狀況信息
func formatDish(dish map[string]string, team string) string {
	if result, ok := dish[team]; ok {
		return result
	}
	return "無可用數據"
}
=======
// func main() {
// 	PKTeam()
// }
>>>>>>> main:nba.go
