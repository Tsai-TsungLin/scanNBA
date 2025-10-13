package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"nba-scanner/internal/logic"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

// Start 啟動 HTTP Server
func Start(port int) error {
	// API endpoint
	http.HandleFunc("/api/games", handleGamesAPI)

	// 靜態檔案（HTML, CSS, JS）
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("無法載入靜態檔案: %w", err)
	}
	http.Handle("/", http.FileServer(http.FS(staticFS)))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("🏀 NBA Scanner 啟動於 http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

// handleGamesAPI 處理 API 請求
func handleGamesAPI(w http.ResponseWriter, r *http.Request) {
	// 設定 CORS 和 JSON header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 取得比賽資料
	games, err := logic.GetTodayGames()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 回傳 JSON
	json.NewEncoder(w).Encode(games)
}
