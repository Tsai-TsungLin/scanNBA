package nba

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Get the injuriers of the nba team
func getInjuryFetch(searchTeam string) []map[string]string {

	var results []map[string]string // 切片，存儲每個球員的信息

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

				link, exists := g.Find(".AnchorLink").Attr("href") // 提取 href 属性
				if !exists {
					link = "" // 如果 href 属性不存在，設置為空字符串
				}

				if status == "Day-To-Day" {
					status = "GTD"
				}

				if status != "" {

					name = fmt.Sprintf("%-25s", name)
					status = fmt.Sprintf("%-15s", status)

					player := map[string]string{
						"name":   strings.TrimSpace(name),   // 移除多餘的空白
						"status": strings.TrimSpace(status), // 移除多餘的空白
						"link":   strings.TrimSpace(link),
					}

					// 將這個map添加到結果切片中
					results = append(results, player)

				}
			}))
		}
	})

	return results

}

// 取得近五場輸贏盤口資訊
func getDishFetch(searchTeam string) []string {
	var (
		season string
	)

	if searchTeam == "LA Clippers" {
		searchTeam = "Los Angeles Clippers"
	}

	// if searchTeam == "Portland Trail Blazers" {
	// 	searchTeam = "Portland Trail_blazers"
	// }

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
	result := []string{}

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

			result = append(result, one, two, three, four, five)
		}
	}

	return result

}

// calculateTeamScores 根据总大小分和让分计算主客场的大小分
func calculateTeamScores(totalOverUnder, pointSpread string) (string, string, error) {

	total, err := strconv.ParseFloat(totalOverUnder, 64)
	if err != nil {
		return "", "", fmt.Errorf("解析总大小分时出错: %v", err)
	}

	spread, err := strconv.ParseFloat(pointSpread, 64)
	if err != nil {
		return "", "", fmt.Errorf("解析让分时出错: %v", err)
	}

	// 计算主客场的预期得分
	homeScore := (total / 2) - (spread / 2)
	awayScore := (total / 2) + (spread / 2)

	if spread < 0 {
		homeScore = homeScore - 0.5
		awayScore = awayScore - 0.5
	}
	// 确保主队和客队的总得分不超过整场的大小分
	if homeScore+awayScore > total {
		return "", "", fmt.Errorf("主队和客队的总得分超过了整场的大小分")
	}

	// 四舍五入到最接近的整数
	homeScoreRounded := math.Round(homeScore)
	awayScoreRounded := math.Round(awayScore)

	// 将得分转换为字符串
	awayScoreStr := fmt.Sprintf("%.0f", homeScoreRounded)
	homeScoreStr := fmt.Sprintf("%.0f", awayScoreRounded)

	return homeScoreStr, awayScoreStr, nil
}
