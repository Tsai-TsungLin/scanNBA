package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Titan007SpreadCache 快取 titan007 過盤資料
var (
	titan007SpreadCache      map[string][]string // map[teamName] = ["贏", "輸", "贏", "贏", "輸"]
	titan007SpreadCacheMutex sync.RWMutex
	titan007SpreadCacheTime  time.Time
)

// FetchTitan007Spreads 抓取 titan007 整季的過盤資料
func FetchTitan007Spreads() (map[string][]string, error) {
	// 檢查快取（1小時有效）
	titan007SpreadCacheMutex.RLock()
	if titan007SpreadCache != nil && time.Since(titan007SpreadCacheTime) < 1*time.Hour {
		defer titan007SpreadCacheMutex.RUnlock()
		return titan007SpreadCache, nil
	}
	titan007SpreadCacheMutex.RUnlock()

	// 計算當前賽季
	// NBA 賽季從10月開始，但 titan007 的資料可能會延遲
	// 在新賽季剛開始時（10月初），先使用上一個賽季的資料
	now := time.Now()
	var season string
	month := int(now.Month())
	day := now.Day()

	if month < 7 {
		// 1-6月：使用上一年的賽季
		season = fmt.Sprintf("%02d-%02d", (now.Year()-1)%100, now.Year()%100)
	} else if month >= 10 && day < 20 {
		// 10月前20天：還是使用上一個賽季（季前賽/賽季初）
		season = fmt.Sprintf("%02d-%02d", (now.Year()-1)%100, now.Year()%100)
	} else {
		// 其他時間：使用當前賽季
		season = fmt.Sprintf("%02d-%02d", now.Year()%100, (now.Year()+1)%100)
	}

	version := now.Format("2006010215")
	url := fmt.Sprintf("https://nba.titan007.com/jsData/letGoal/%s/l1.js?version=%s", season, version)

	log.Printf("抓取 titan007 過盤資料: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("titan007 請求失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("titan007 返回狀態碼: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取 titan007 回應失敗: %w", err)
	}

	// 解析 titan007 的 JS 資料
	spreadMap, err := parseTitan007Spreads(string(body))
	if err != nil {
		return nil, fmt.Errorf("解析 titan007 資料失敗: %w", err)
	}

	// 更新快取
	titan007SpreadCacheMutex.Lock()
	titan007SpreadCache = spreadMap
	titan007SpreadCacheTime = time.Now()
	titan007SpreadCacheMutex.Unlock()

	log.Printf("成功快取 titan007 過盤資料，共 %d 支球隊", len(spreadMap))

	return spreadMap, nil
}

// parseTitan007Spreads 解析 titan007 的 JS 資料
func parseTitan007Spreads(jsData string) (map[string][]string, error) {
	result := make(map[string][]string)

	// 找出分號位置（用於定位資料區塊）
	r1 := regexp.MustCompile(";")
	r2 := r1.FindAllStringSubmatchIndex(jsData, -1)

	if len(r2) < 3 {
		return nil, fmt.Errorf("資料格式不符合預期（分號數量不足）")
	}

	// 提取隊伍資料區塊
	Team := string([]byte(jsData[r2[0][0]+17 : r2[1][0]-1]))
	frontside := regexp.MustCompile(`\[(.*?)]`)
	frontsideTeam := frontside.FindAllStringSubmatchIndex(Team, -1)

	// 建立 TeamNumber -> TeamName 的映射
	matchMap := make(map[int64]string)

	for i := range frontsideTeam {
		TeamData := Team[frontsideTeam[i][0]:frontsideTeam[i][1]]
		// 用字串分割（更可靠）
		parts := regexp.MustCompile(",").Split(TeamData[1:len(TeamData)-1], -1)

		if len(parts) < 4 {
			continue
		}

		TeamNumber, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}

		// 球隊名稱在第3個位置（index 3）
		teamName := parts[3]
		if len(teamName) > 0 && teamName[0] == '\'' && teamName[len(teamName)-1] == '\'' {
			teamName = teamName[1 : len(teamName)-1]
		}
		matchMap[TeamNumber] = teamName
	}

	// 提取過盤資料區塊
	data := string([]byte(jsData[r2[1][0]+20 : r2[2][0]-1]))
	frontsidedata := frontside.FindAllStringSubmatchIndex(data, -1)

	// 解析每支球隊的過盤資料
	for i := range frontsidedata {
		winPercentData := data[frontsidedata[i][0]:frontsidedata[i][1]]

		// 用字串分割（更可靠）
		parts := regexp.MustCompile(",").Split(winPercentData[1:len(winPercentData)-1], -1)

		if len(parts) < 17 {
			continue
		}

		TeamNumber, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}

		// 檢查是否為目標球隊
		if teamName, ok := matchMap[TeamNumber]; ok {
			// 近5場過盤結果在 index 12-16
			// 資料格式: [TeamNum, ..., value12, value13, value14, value15, value16]
			// "0" = 沒過盤（輸）, "2" = 過盤（贏）
			spreads := make([]string, 5)

			for j := 0; j < 5; j++ {
				spreadValue := parts[12+j]
				if spreadValue == "2" {
					spreads[j] = "W" // 過盤
				} else {
					spreads[j] = "L" // 沒過盤（0 或其他值）
				}
			}

			result[teamName] = spreads
		}
	}

	return result, nil
}

// GetTeamSpreads 獲取指定球隊的近5場過盤結果
func GetTeamSpreads(teamName string) ([]string, bool) {
	spreadMap, err := FetchTitan007Spreads()
	if err != nil {
		log.Printf("抓取 titan007 過盤資料失敗: %v", err)
		return nil, false
	}

	// 處理 LA Clippers 特殊情況
	if teamName == "Los Angeles Clippers" {
		teamName = "LA Clippers"
	}

	if spreads, ok := spreadMap[teamName]; ok {
		return spreads, true
	}

	return nil, false
}
