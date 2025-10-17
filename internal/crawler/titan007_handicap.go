package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Titan007TeamIDMap NBA 球隊 ID 映射表 (Titan007 TeamID -> NBA 英文名稱)
var Titan007TeamIDMap = map[int]string{
	1:  "Los Angeles Lakers",
	2:  "Boston Celtics",
	3:  "Miami Heat",
	4:  "Brooklyn Nets",
	5:  "New York Knicks",
	6:  "Orlando Magic",
	7:  "Philadelphia 76ers",
	8:  "Washington Wizards",
	9:  "Detroit Pistons",
	10: "Indiana Pacers",
	11: "New Orleans Pelicans",
	12: "Milwaukee Bucks",
	13: "Atlanta Hawks",
	14: "Chicago Bulls",
	15: "Toronto Raptors",
	16: "Cleveland Cavaliers",
	17: "Dallas Mavericks",
	18: "San Antonio Spurs",
	19: "Minnesota Timberwolves",
	20: "Utah Jazz",
	21: "Houston Rockets",
	22: "Memphis Grizzlies",
	23: "Denver Nuggets",
	24: "Sacramento Kings",
	25: "Portland Trail Blazers",
	26: "Phoenix Suns",
	27: "Golden State Warriors",
	28: "Oklahoma City Thunder",
	29: "Los Angeles Clippers",
	30: "Charlotte Hornets",
}

// Titan007HandicapCache 快取 titan007 盤口戰績資料
var (
	titan007HandicapCache      map[string][]string // map[teamName] = ["W", "L", "W", "W", "L"]
	titan007HandicapCacheMutex sync.RWMutex
	titan007HandicapCacheTime  time.Time
)

// HandicapGame 盤口戰績單場比賽
type HandicapGame struct {
	GameID       int
	GameType     int     // 1=常規賽, 2=季後賽, 3=季前賽
	GameTime     string  // 2025/10/05 08:00
	HomeTeamID   int
	AwayTeamID   int
	HomeScore    int
	AwayScore    int
	Spread       float64 // 盤口
	SpreadResult int     // 1=贏盤, 2=走盤, 3=輸盤
}

// FetchTitan007TeamHandicap 從 HandicapDetail 頁面抓取球隊盤口戰績
func FetchTitan007TeamHandicap(teamNameEN string, limit int) ([]string, error) {
	// 先從 letGoal API 取得 TeamID 映射
	teamIDMap, err := fetchTeamIDMap()
	if err != nil {
		return nil, fmt.Errorf("無法取得球隊 ID 映射: %w", err)
	}

	// 查找球隊 ID
	teamID := -1
	for id, name := range teamIDMap {
		if name == teamNameEN {
			teamID = id
			break
		}
	}

	if teamID == -1 {
		return nil, fmt.Errorf("找不到球隊: %s", teamNameEN)
	}

	// 抓取該球隊的盤口戰績頁面
	games, err := fetchHandicapDetailPage(teamID)
	if err != nil {
		return nil, err
	}

	// 只取最近的 N 場比賽
	startIdx := len(games) - limit
	if startIdx < 0 {
		startIdx = 0
	}

	// 先收集結果（此時順序是由舊到新）
	temp := make([]string, 0, limit)
	for i := startIdx; i < len(games); i++ {
		game := games[i]
		if game.SpreadResult == 1 {
			temp = append(temp, "W") // 贏盤
		} else if game.SpreadResult == 3 {
			temp = append(temp, "L") // 輸盤
		} else {
			temp = append(temp, "W") // 走盤當作平手，這裡簡化為 W
		}
	}

	// 反轉順序，使其由新到舊（配合 history.go 的 games 排序）
	results := make([]string, len(temp))
	for i := range temp {
		results[i] = temp[len(temp)-1-i]
	}

	return results, nil
}

// fetchTeamIDMap 取得 TeamID 映射（直接使用寫死的映射表）
func fetchTeamIDMap() (map[int]string, error) {
	return Titan007TeamIDMap, nil
}

// fetchHandicapDetailPage 抓取 HandicapDetail 頁面
func fetchHandicapDetailPage(teamID int) ([]HandicapGame, error) {
	// 構建 URL
	season := "2025-2026" // 當前賽季，後續可以動態計算
	url := fmt.Sprintf("https://nba.titan007.com/cn/Team/HandicapDetail.aspx?sclassid=1&teamid=%d&matchseason=%s&halfOrAll=0", teamID, season)

	log.Printf("抓取 titan007 盤口戰績: %s", url)

	// 建立 HTTP 請求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("建立請求失敗: %w", err)
	}

	// 設定 headers 模擬瀏覽器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://nba.titan007.com/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("請求失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP 狀態碼: %d", resp.StatusCode)
	}

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取失敗: %w", err)
	}

	html := string(body)

	// 解析 handicapDetail 陣列
	games, err := parseHandicapDetail(html)
	if err != nil {
		return nil, fmt.Errorf("解析失敗: %w", err)
	}

	log.Printf("成功抓取 %d 場盤口戰績", len(games))

	return games, nil
}

// parseHandicapDetail 解析 handicapDetail JavaScript 陣列
func parseHandicapDetail(html string) ([]HandicapGame, error) {
	// 提取 handicapDetail 陣列
	// 格式: var handicapDetail = [[670554,1,3,'2025/10/05 08:00',3,6,118,126,61,50,-4.5,1],...];
	r := regexp.MustCompile(`var handicapDetail = (\[\[.*?\]\]);`)
	matches := r.FindStringSubmatch(html)

	if len(matches) < 2 {
		return nil, fmt.Errorf("找不到 handicapDetail 資料")
	}

	// 清理資料（去掉單引號換成雙引號以符合 JSON 格式）
	dataStr := strings.ReplaceAll(matches[1], "'", "\"")

	// 解析 JSON
	var rawData [][]interface{}
	if err := json.Unmarshal([]byte(dataStr), &rawData); err != nil {
		return nil, fmt.Errorf("JSON 解析失敗: %w", err)
	}

	// 轉換為 HandicapGame 結構
	games := make([]HandicapGame, 0, len(rawData))
	for _, game := range rawData {
		if len(game) < 12 {
			continue
		}

		games = append(games, HandicapGame{
			GameID:       int(game[0].(float64)),
			GameType:     int(game[2].(float64)),
			GameTime:     game[3].(string),
			HomeTeamID:   int(game[4].(float64)),
			AwayTeamID:   int(game[5].(float64)),
			HomeScore:    int(game[6].(float64)),
			AwayScore:    int(game[7].(float64)),
			Spread:       game[10].(float64),
			SpreadResult: int(game[11].(float64)),
		})
	}

	return games, nil
}

// GetTeamHandicapSpreads 獲取指定球隊的近N場盤口結果（替換舊的 GetTeamSpreads）
func GetTeamHandicapSpreads(teamName string, limit int) ([]string, bool) {
	// 處理 LA Clippers 特殊情況
	if teamName == "Los Angeles Clippers" {
		teamName = "LA Clippers"
	}

	spreads, err := FetchTitan007TeamHandicap(teamName, limit)
	if err != nil {
		log.Printf("抓取 titan007 盤口戰績失敗 (%s): %v", teamName, err)
		return nil, false
	}

	return spreads, true
}
