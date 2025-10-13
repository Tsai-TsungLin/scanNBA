// ============ internal/models/schedule.go ============
package models

type Schedule struct {
	Payload struct {
		Date struct {
			Games []struct {
				Profile struct {
					DateTimeEt string `json:"dateTimeEt"`
					HomeTeamID string `json:"homeTeamId"`
					AwayTeamID string `json:"awayTeamId"`
				} `json:"profile"`
				HomeTeam struct {
					Profile struct {
						City string `json:"city"`
						Name string `json:"name"`
					} `json:"profile"`
				} `json:"homeTeam"`
				AwayTeam struct {
					Profile struct {
						City string `json:"city"`
						Name string `json:"name"`
					} `json:"profile"`
				} `json:"awayTeam"`
			} `json:"games"`
			GameCount string `json:"gameCount"`
		} `json:"date"`
	} `json:"payload"`
}
