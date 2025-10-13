package models

// APIResponse API 回應格式
type APIResponse struct {
	Date  string     `json:"date"`
	Games []GameInfo `json:"games"`
}

// GameInfo 單場比賽資訊（用於前端）
type GameInfo struct {
	GameID         string          `json:"gameId"`         // 比賽 ID
	GameTime       string          `json:"gameTime"`       // 開賽時間或比賽狀態 (Q4 04:00)
	GameStatus     int             `json:"gameStatus"`     // 1=未開始, 2=進行中, 3=已結束
	GameStatusText string          `json:"gameStatusText"` // 狀態文字
	ScoreDisplay   string          `json:"scoreDisplay"`   // 比分顯示 "[98-128]" 或空字串
	HomeTeam       TeamInfo        `json:"homeTeam"`
	AwayTeam       TeamInfo        `json:"awayTeam"`
	HomeScore      int             `json:"homeScore"`      // 主隊比分
	AwayScore      int             `json:"awayScore"`      // 客隊比分
	Spread         SpreadDisplay   `json:"spread"`
	HomeInjuries   []string        `json:"homeInjuries"`
	AwayInjuries   []string        `json:"awayInjuries"`
	HomeHistory    *TeamHistory    `json:"homeHistory,omitempty"`
	AwayHistory    *TeamHistory    `json:"awayHistory,omitempty"`
	HomePlayers    []PlayerDisplay `json:"homePlayers,omitempty"`    // 主隊球員數據（進行中時才有）
	AwayPlayers    []PlayerDisplay `json:"awayPlayers,omitempty"`    // 客隊球員數據（進行中時才有）
	PeriodScores   *PeriodScores   `json:"periodScores,omitempty"`   // 各節比分（進行中或已結束才有）
}

// PeriodScores 各節比分顯示
type PeriodScores struct {
	HomePeriods []int `json:"homePeriods"` // 主隊各節得分 [Q1, Q2, Q3, Q4]
	AwayPeriods []int `json:"awayPeriods"` // 客隊各節得分 [Q1, Q2, Q3, Q4]
}

// TeamInfo 球隊資訊
type TeamInfo struct {
	NameEN string `json:"nameEN"`
	NameCN string `json:"nameCN"`
	Wins   int    `json:"wins"`   // 勝場
	Losses int    `json:"losses"` // 敗場
}

// SpreadDisplay 即時盤口顯示（只顯示主隊）
type SpreadDisplay struct {
	Opening string `json:"opening"` // 開盤讓分
	Current string `json:"current"` // 即時讓分
	HasData bool   `json:"hasData"` // 是否有賠率資料
}
