# NBA Scanner

NBA 賽事即時分析工具，提供今日賽程、傷兵報告、讓分走勢、分節比分、球員數據等完整資訊。

## 功能特色

- 📅 **今日賽程**：自動抓取 NBA 官方賽程資料
- 🏥 **傷兵報告**：即時更新球隊傷兵狀態（ESPN 資料源）
- 📊 **讓分走勢**：顯示開盤讓分與開賽前盤口
- 🏀 **戰績追蹤**：顯示主客隊近五場戰績（勝敗/讓分結果）
- 🎯 **分節比分**：顯示 Q1-Q4 各節得分（進行中/已結束比賽）
- 💯 **球員數據**：即時顯示球員得分、籃板、助攻（進行中/已結束比賽）
- 📈 **讓分分析**：上半場/全場讓分與大小分結果分析
- 🌐 **Web 介面**：美觀的網頁介面，支援行動裝置
- 🔄 **自動更新**：每 5 分鐘自動刷新資料

## 快速開始

### 方法一：Docker（推薦）

```bash
# 建置 Docker 映像檔
docker build -t nba-scanner .

# 運行容器
docker run -p 8081:8081 nba-scanner

# 或使用 docker-compose
docker-compose up -d
```

### 方法二：本地編譯

```bash
# 安裝依賴
go mod download

# 編譯
go build -o nba-scanner

# 啟動 Web 服務
./nba-scanner --server --port 8081
```

開啟瀏覽器訪問：http://localhost:8081

## 專案架構

```
scanNBA/
├── cmd/                    # CLI 指令
│   └── root.go
├── internal/
│   ├── crawler/           # 資料爬取
│   │   ├── schedule.go    # 賽程資料
│   │   ├── odds.go        # 賠率資料
│   │   ├── injury.go      # 傷兵資料
│   │   ├── history.go     # 戰績資料
│   │   └── history_cache.go # 戰績快取（1小時）
│   ├── logic/             # 業務邏輯
│   │   └── api.go         # API 處理
│   ├── models/            # 資料模型
│   │   ├── nba_schedule.go
│   │   ├── nba_odds.go
│   │   ├── team_history.go
│   │   ├── api_response.go
│   │   └── teammap.go     # 中英文隊名對照
│   └── server/            # Web 服務
│       ├── server.go
│       └── static/
│           └── index.html
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## API 資料來源

| 資料類型 | API 端點 | 說明 |
|---------|---------|------|
| 賽程 | `https://cdn.nba.com/static/json/liveData/scoreboard/todaysScoreboard_00.json` | 比賽時間、比分、各節得分 |
| 賠率 | `https://cdn.nba.com/static/json/liveData/odds/odds_todaysGames.json` | 開盤讓分（僅開賽前） |
| 球員數據 | `https://cdn.nba.com/static/json/liveData/boxscore/boxscore_{gameId}.json` | 球員得分、籃板、助攻 |
| 傷兵 | `https://www.espn.com/nba/injuries` | 球隊傷兵清單 |
| 戰績 | `https://cdn.nba.com/static/json/staticData/scheduleLeagueV2_9.json` | 近五場戰績 |

**重要提醒**：NBA 官方賠率 API 不提供即時讓分更新，僅顯示開賽前的盤口資料。

## 技術特點

### 效能優化
- ⚡ 平行抓取多個資料源（使用 goroutines）
- 💾 戰績資料快取（1 小時 TTL）
- 🔍 O(1) 時間複雜度的資料查找（使用 map）

### 資料處理
- 🕐 自動轉換為台北時區（UTC+8）
- 🏷️ 隊名中英文對照（30 支球隊完整對應）
- 🎨 傷兵狀態色彩標記（Day-To-Day/Out/Questionable）
- 📈 讓分分析（上半場/全場實際分差計算）

### 讓分分析說明
系統會自動計算：
1. **實際分差**：主隊得分 - 客隊得分
   - 顯示格式：`主-24`（主隊贏 24 分）或 `客-15`（客隊贏 15 分）
2. **讓分結果**：實際分差 + 讓分值
   - 🟢 綠色：讓分贏盤（結果 > 0）
   - 🔴 紅色：讓分輸盤（結果 < 0）
3. **上半場分析**：Q1+Q2 得分計算
4. **全場分析**：總得分計算

**範例**：
- 盤口：主隊 -6.5
- 上半場：主 70 - 客 46 = **主-24**（實際領先 24 分）
- 讓分結果：24 + (-6.5) = 17.5 ✅ **贏盤**（綠色顯示）

### Web 功能
- 📱 響應式設計（手機/平板/電腦）
- 🎭 收合/展開介面（球員數據/戰績）
- 🔄 自動刷新（5 分鐘）
- 🎯 懸停顯示詳細資訊
- 🏀 場上球員標記（進行中比賽用籃球圖示標示）

## GCP 部署指南

詳見 [DEPLOY_GCP.md](DEPLOY_GCP.md) - 完整的 GCP VM + Docker 部署流程，使用最便宜的配置避免產生費用。

## 開發

### 本地開發環境

```bash
# 安裝依賴
go mod download

# 執行（開發模式）
go run main.go --server --port 8081

# 編譯
go build -o nba-scanner

# 測試
go test ./...

# 檢查端口佔用
lsof -i :8081

# 停止所有 Go 進程
pkill -f "go run main.go"
```

### 更新靜態文件

由於使用 `//go:embed` 嵌入靜態文件，修改 HTML/CSS/JS 後需要重新編譯：

```bash
go build -o nba-scanner
```

## 環境變數

| 變數 | 說明 | 預設值 |
|-----|------|--------|
| `TZ` | 時區設定 | `Asia/Taipei` |
| `PORT` | 服務端口 | `8081` |

## 常見問題

**Q: 為什麼比賽進行中盤口沒有更新？**
A: NBA 官方 API 不提供即時賠率更新，只有開賽前的盤口資料。顯示的「開賽前盤口」是比賽開始前的最後盤口。

**Q: 為什麼沒有顯示賠率資料？**
A: 確認當天有比賽，且賠率 API 已更新資料（通常在比賽日前一天晚上更新）。

**Q: 戰績資料多久更新一次？**
A: 使用 1 小時快取機制，每小時自動更新。

**Q: 如何清除快取？**
A: 重啟服務即可清除所有快取。

**Q: 瀏覽器看不到最新變更？**
A: 使用強制刷新（Cmd+Shift+R 或 Ctrl+F5）清除瀏覽器快取。

**Q: 為什麼有些球員數據是空的？**
A: 球員數據僅在比賽進行中或結束後才會顯示。未開始的比賽不會有球員數據。

**Q: 讓分計算是如何運作的？**
A:
- 實際分差 = 主隊得分 - 客隊得分
- 讓分結果 = 實際分差 + 讓分值
- 例：主隊 -6.5，實際贏 10 分，則結果為 10 + (-6.5) = 3.5（贏盤）

## 授權

MIT License

## 更新日誌

### v2.1.0 (2025-10-13)
- ✨ 新增分節比分顯示（Q1-Q4）
- 💯 新增球員即時數據（得分/籃板/助攻）
- 📈 新增讓分分析（上半場/全場讓分大小分結果）
- 🏀 新增場上球員標記（籃球圖示）
- 🎨 優化版面配置（三欄式佈局）
- 🔧 修正讓分計算邏輯
- 📝 標註賠率 API 限制（僅開賽前數據）

### v2.0.0 (2025-10-13)
- ✨ 重構為模組化架構
- 🌐 新增 Web 介面
- 📊 新增近五場戰績功能
- 🐳 支援 Docker 部署
- ⚡ 效能優化（平行抓取、快取機制）
- 🕐 時區自動轉換（台北時間）
- 🎨 UI/UX 改善（收合、懸停提示）

### v1.0.0
- 🎉 初始版本
- 基本賽程、傷兵、讓分功能
