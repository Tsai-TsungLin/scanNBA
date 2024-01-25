package nba

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
