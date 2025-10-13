package models

// TeamHistory 球隊歷史戰績
type TeamHistory struct {
	TeamID       int            `json:"teamId"`
	TeamName     string         `json:"teamName"`
	RecentGames  []GameResult   `json:"recentGames"`
	WinCount     int            `json:"winCount"`
	LossCount    int            `json:"lossCount"`
}

// GameResult 單場比賽結果
type GameResult struct {
	Date       string `json:"date"`       // 比賽日期
	Opponent   string `json:"opponent"`   // 對手
	IsHome     bool   `json:"isHome"`     // 是否主場
	Score      string `json:"score"`      // 比分 (如: "111-103")
	Result     string `json:"result"`     // W 或 L
}

// FullSchedule 完整賽季賽程
type FullSchedule struct {
	LeagueSchedule struct {
		GameDates []GameDate `json:"gameDates"`
	} `json:"leagueSchedule"`
}

// GameDate 某日的比賽
type GameDate struct {
	GameDate string           `json:"gameDate"`
	Games    []ScheduledGame  `json:"games"`
}

// ScheduledGame 賽程中的比賽
type ScheduledGame struct {
	GameID          string `json:"gameId"`
	GameCode        string `json:"gameCode"`
	GameStatus      int    `json:"gameStatus"`   // 1=未開始 2=進行中 3=已結束
	GameStatusText  string `json:"gameStatusText"`
	GameDateTimeEst string `json:"gameDateTimeEst"` // 東岸時間，格式: 2025-10-02T12:00:00Z
	HomeTeam        ScheduledTeam `json:"homeTeam"`
	AwayTeam        ScheduledTeam `json:"awayTeam"`
}

// ScheduledTeam 賽程中的球隊資訊
type ScheduledTeam struct {
	TeamID   int    `json:"teamId"`
	TeamName string `json:"teamName"`
	TeamCity string `json:"teamCity"`
	Score    int    `json:"score"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
}
