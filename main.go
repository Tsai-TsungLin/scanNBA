package main

import (
	"embed"
	"example/scanNBA/nba"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// 创建一个新的HTTP多路复用器
	mux := http.NewServeMux()

	// 设置静态文件服务
	// 使用embed文件系统代替http.Dir
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// 设置API端点
	mux.HandleFunc("/api/games", gamesHandler)

	log.Println("Server starting on :8080...")
	err = http.ListenAndServe(":8080", mux)
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
