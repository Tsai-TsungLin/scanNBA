package models

// NBAScoreboard 新的賽程 API 回應結構
type NBAScoreboard struct {
	Scoreboard struct {
		GameDate string `json:"gameDate"`
		Games    []Game `json:"games"`
	} `json:"scoreboard"`
}

// Game 比賽資訊
type Game struct {
	GameID         string `json:"gameId"`
	GameCode       string `json:"gameCode"`
	GameStatus     int    `json:"gameStatus"`
	GameStatusText string `json:"gameStatusText"`
	Period         int    `json:"period"`
	GameClock      string `json:"gameClock"`
	GameTimeUTC    string `json:"gameTimeUTC"`
	HomeTeam       Team   `json:"homeTeam"`
	AwayTeam       Team   `json:"awayTeam"`
}

// Team 球隊資訊
type Team struct {
	TeamID      int            `json:"teamId"`
	TeamName    string         `json:"teamName"`
	TeamCity    string         `json:"teamCity"`
	TeamTricode string         `json:"teamTricode"`
	Score       int            `json:"score"`
	Wins        int            `json:"wins"`
	Losses      int            `json:"losses"`
	Periods     []PeriodScore  `json:"periods"`
}

// PeriodScore 單節比分
type PeriodScore struct {
	Period     int    `json:"period"`
	PeriodType string `json:"periodType"`
	Score      int    `json:"score"`
}

// GetFullTeamName 取得完整隊名 (City + Name)
func (t *Team) GetFullTeamName() string {
	return t.TeamCity + " " + t.TeamName
}
