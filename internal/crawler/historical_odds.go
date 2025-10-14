package crawler

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// FetchHistoricalSpread 取得歷史比賽的盤口資料
// 目前從 NBA odds API 取得（只有當天或最近的比賽）
// TODO: 整合 titan007 或其他資料源以取得更早期的歷史盤口
func FetchHistoricalSpread(gameID string) (homeSpread float64, hasSpread bool) {
	// 嘗試從 NBA odds API 取得
	odds, err := FetchOdds()
	if err != nil {
		log.Printf("無法取得賠率資料: %v", err)
		return 0, false
	}

	// 建立盤口查找表
	oddsMap := BuildOddsMap(odds)

	// 查找該場比賽的盤口
	if spreadInfo, ok := oddsMap[gameID]; ok && spreadInfo.Found {
		// 使用開盤盤口（opening spread）來判斷過盤
		return spreadInfo.HomeOpeningSpread, true
	}

	return 0, false
}

// CalculateSpreadResult 計算是否過盤
// actualDiff: 實際分差（主隊得分 - 客隊得分 或 客隊得分 - 主隊得分）
// spread: 盤口讓分（主隊讓分 或 客隊讓分）
// 返回: true=過盤(W), false=沒過盤(L)
func CalculateSpreadResult(actualDiff int, spread float64) bool {
	// 過盤判斷：實際分差 + 盤口 > 0 為過盤
	return float64(actualDiff)+spread > 0
}

// FormatSpreadDisplay 格式化盤口顯示
// spread: 盤口值（正數表示受讓，負數表示讓分）
// isHome: 是否為主隊
// 返回格式: "主讓5.5" 或 "客讓3.5" 或 "主受3.5" 或 "客受2.5"
func FormatSpreadDisplay(spread float64, isHome bool) string {
	teamType := "主"
	if !isHome {
		teamType = "客"
	}

	absSpread := spread
	if absSpread < 0 {
		absSpread = -absSpread
	}

	// 負數表示讓分，正數表示受讓
	if spread < 0 {
		return fmt.Sprintf("%s讓%.1f", teamType, absSpread)
	} else if spread > 0 {
		return fmt.Sprintf("%s受%.1f", teamType, absSpread)
	} else {
		return "無讓分"
	}
}

// ParseGameIDToScheduleID 將 NBA GameID 轉換為可能的其他系統 ID
// 例如: "0012500019" -> 可能對應到 titan007 或其他系統的 ID
// TODO: 實作 GameID 映射邏輯
func ParseGameIDToScheduleID(gameID string) (int, error) {
	// NBA GameID 格式: 001SSGGGG
	// SS = 賽季年份後兩位
	// GGGG = 比賽編號

	if len(gameID) < 10 {
		return 0, fmt.Errorf("invalid gameID format: %s", gameID)
	}

	// 提取比賽編號
	gameNumStr := gameID[5:]
	gameNum, err := strconv.Atoi(strings.TrimLeft(gameNumStr, "0"))
	if err != nil {
		return 0, fmt.Errorf("failed to parse game number: %w", err)
	}

	return gameNum, nil
}

// FetchTitan007HistoricalSpread 從 titan007 取得歷史盤口
// TODO: 實作從 titan007 API 取得歷史比賽盤口
func FetchTitan007HistoricalSpread(gameDate string, homeTeam string, awayTeam string) (homeSpread float64, hasSpread bool) {
	// 這裡需要：
	// 1. 根據日期和球隊名稱查詢 titan007
	// 2. 解析 JS 資料取得盤口
	// 3. 返回盤口值

	// 暫時返回無盤口
	return 0, false
}
