package nba

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetNBATeams() {
	// startTime := time.Now()

	today := time.Now()

	// 檢查時間是否已過中午12點
	var newTime string
	if today.Hour() >= 12 {
		newTime = today.Format("20060102")
	} else {
		newTime = today.AddDate(0, 0, -1).Format("20060102") // 如果沒有過中午12點，則減去一天
	}

	// 獲取數據 API URL
	url := "https://site.api.espn.com/apis/v2/scoreboard/header?sport=basketball&league=nba&dates=" + newTime + "&tz=America%2FNew_York&showAirings=buy%2Clive%2Creplay&showZipLookup=true&buyWindow=1m&lang=en&region=us&contentorigin=espn"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result NBATeam

	if unmarshalErr := json.Unmarshal(body, &result); unmarshalErr != nil {
		panic(unmarshalErr)
	}

	log.Println(result.Sports[0].Leagues[0].Events[0].Name)
}
