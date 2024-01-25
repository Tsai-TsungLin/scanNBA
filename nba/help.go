package nba

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

	log.Println("URL: ", url)

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

// formatDish 格式化過盤狀況信息
func formatDish(dish map[string]string, team string) string {
	if result, ok := dish[team]; ok {
		return result
	}
	return "無可用數據"
}
