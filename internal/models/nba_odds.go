package models

// NBAOdds 賠率 API 回應結構
type NBAOdds struct {
	Games []OddsGame `json:"games"`
}

// OddsGame 單場比賽的賠率資訊
type OddsGame struct {
	GameID     string       `json:"gameId"`
	HomeTeamID string       `json:"homeTeamId"`
	AwayTeamID string       `json:"awayTeamId"`
	Markets    []OddsMarket `json:"markets"`
}

// OddsMarket 賠率市場
type OddsMarket struct {
	Name       string      `json:"name"`
	OddsTypeID int         `json:"odds_type_id"`
	Books      []Bookmaker `json:"books"`
}

// Bookmaker 莊家資訊
type Bookmaker struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Outcomes []Outcome `json:"outcomes"`
}

// Outcome 賠率結果
type Outcome struct {
	Type          string  `json:"type"`
	Odds          string  `json:"odds"`
	Spread        string  `json:"spread,omitempty"`
	OpeningSpread float64 `json:"opening_spread,omitempty"`
}

// SpreadInfo 整理後的讓分盤資訊
type SpreadInfo struct {
	HomeSpread        string
	AwaySpread        string
	HomeOpeningSpread float64
	AwayOpeningSpread float64
	Found             bool
}

// GetSupermatchSpread 取得 Supermatch 的讓分盤資訊
func (og *OddsGame) GetSupermatchSpread() SpreadInfo {
	result := SpreadInfo{Found: false}

	// 找到 spread market (name = "spread")
	for _, market := range og.Markets {
		if market.Name == "spread" {
			// 找 Supermatch bookmaker (id: "sr:book:818")
			for _, bookmaker := range market.Books {
				if bookmaker.ID == "sr:book:818" {
					result.Found = true
					// 解析 home 和 away 的賠率
					for _, outcome := range bookmaker.Outcomes {
						if outcome.Type == "home" {
							result.HomeSpread = outcome.Spread
							result.HomeOpeningSpread = outcome.OpeningSpread
						} else if outcome.Type == "away" {
							result.AwaySpread = outcome.Spread
							result.AwayOpeningSpread = outcome.OpeningSpread
						}
					}
					return result
				}
			}
		}
	}

	return result
}
