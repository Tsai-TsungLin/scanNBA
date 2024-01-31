package main

import (
	"example/scanNBA/nba"
	"log"
	"net/http"
)

func main() {
	// 設置靜態文件服務
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// 設置API端點
	http.HandleFunc("/api/games", gamesHandler)

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func gamesHandler(w http.ResponseWriter, r *http.Request) {
	// 創建一些假數據來返回
	games := nba.GetNBAData()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(games)
}
