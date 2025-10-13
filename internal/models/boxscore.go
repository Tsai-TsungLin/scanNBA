package models

// BoxscoreResponse Boxscore API 回應
type BoxscoreResponse struct {
	Game BoxscoreGame `json:"game"`
}

// BoxscoreGame 比賽數據
type BoxscoreGame struct {
	GameID     string         `json:"gameId"`
	GameStatus int            `json:"gameStatus"`
	HomeTeam   BoxscoreTeam   `json:"homeTeam"`
	AwayTeam   BoxscoreTeam   `json:"awayTeam"`
}

// BoxscoreTeam 球隊數據
type BoxscoreTeam struct {
	TeamID   int              `json:"teamId"`
	TeamName string           `json:"teamName"`
	Players  []BoxscorePlayer `json:"players"`
}

// BoxscorePlayer 球員數據
type BoxscorePlayer struct {
	PersonID   int              `json:"personId"`
	Name       string           `json:"name"`
	NameI      string           `json:"nameI"`
	JerseyNum  string           `json:"jerseyNum"`
	Position   string           `json:"position"`
	Starter    string           `json:"starter"`    // "1" = 先發, "0" = 替補
	OnCourt    string           `json:"oncourt"`    // "1" = 在場上, "0" = 不在場上
	Played     string           `json:"played"`     // "1" = 有上場, "0" = 未上場
	Statistics PlayerStatistics `json:"statistics"`
}

// PlayerStatistics 球員統計
type PlayerStatistics struct {
	Points              int     `json:"points"`
	ReboundsTotal       int     `json:"reboundsTotal"`
	Assists             int     `json:"assists"`
	Steals              int     `json:"steals"`
	Blocks              int     `json:"blocks"`
	FieldGoalsMade      int     `json:"fieldGoalsMade"`
	FieldGoalsAttempted int     `json:"fieldGoalsAttempted"`
	ThreePointersMade   int     `json:"threePointersMade"`
	Minutes             string  `json:"minutes"`
	PlusMinusPoints     float64 `json:"plusMinusPoints"`
}

// PlayerDisplay 用於前端顯示的球員資料
type PlayerDisplay struct {
	Name      string `json:"name"`
	NameI     string `json:"nameI"`
	JerseyNum string `json:"jerseyNum"`
	Position  string `json:"position"`
	IsStarter bool   `json:"isStarter"`
	OnCourt   bool   `json:"onCourt"`
	Points    int    `json:"points"`
	Rebounds  int    `json:"rebounds"`
	Assists   int    `json:"assists"`
	Minutes   string `json:"minutes"`
}
