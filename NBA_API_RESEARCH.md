# NBA 賽程 API 歷史戰績調查報告

## 調查日期
2025-10-09

## 調查目標
研究 NBA 官方 API 是否能查詢球隊的歷史戰績和對戰紀錄，包括：
- 球隊最近 5-10 場比賽的輸贏紀錄
- 每場比賽的對手是誰
- 比賽的比分
- 比賽日期

## 調查結果

### ✅ 可用的 NBA 官方 API

#### 1. **完整賽程 API (推薦使用)**
```
https://cdn.nba.com/static/json/staticData/scheduleLeagueV2_9.json
```

**功能特點：**
- ✅ 包含整個賽季的所有比賽資料
- ✅ 包含比賽狀態（未開始/進行中/已結束）
- ✅ 包含比分（已結束的比賽）
- ✅ 包含球隊戰績（wins/losses）
- ✅ 可以透過 teamId 過濾特定球隊的比賽
- ✅ 包含比賽日期和時間

**資料結構範例：**
```json
{
  "leagueSchedule": {
    "seasonYear": "2025-26",
    "gameDates": [
      {
        "gameDate": "10/05/2025 00:00:00",
        "games": [
          {
            "gameId": "0012500030",
            "gameStatus": 3,
            "gameStatusText": "Final",
            "homeTeam": {
              "teamId": 1610612744,
              "teamName": "Warriors",
              "teamCity": "Golden State",
              "teamTricode": "GSW",
              "wins": 1,
              "losses": 0,
              "score": 111
            },
            "awayTeam": {
              "teamId": 1610612747,
              "teamName": "Lakers",
              "teamCity": "Los Angeles",
              "teamTricode": "LAL",
              "wins": 0,
              "losses": 2,
              "score": 103
            }
          }
        ]
      }
    ]
  }
}
```

**使用方式：**
1. 下載完整賽程 JSON
2. 根據 teamId 過濾比賽
3. 篩選 gameStatus=3 的已結束比賽
4. 根據分數判斷輸贏

#### 2. **今日賽程 API**
```
https://cdn.nba.com/static/json/liveData/scoreboard/todaysScoreboard_00.json
或
https://nba-prod-us-east-1-mediaops-stats.s3.amazonaws.com/NBA/liveData/scoreboard/todaysScoreboard_00.json
```

**功能特點：**
- ✅ 即時更新今日比賽
- ✅ 包含即時比分
- ✅ 包含球隊戰績
- ✅ 包含比賽狀態

**限制：**
- ❌ 只有當天的比賽
- ❌ 無法查詢歷史比賽

#### 3. **比賽詳細資料 API (Boxscore)**
```
https://cdn.nba.com/static/json/liveData/boxscore/boxscore_{gameId}.json
```

**功能特點：**
- ✅ 詳細的比賽統計資料
- ✅ 球員個人數據
- ✅ 球隊數據
- ✅ 每節比分

**範例：**
```
https://cdn.nba.com/static/json/liveData/boxscore/boxscore_0012500001.json
```

### ❌ 無法使用的 API

#### 1. **歷史日期 Scoreboard**
```
https://cdn.nba.com/static/json/liveData/scoreboard/20251007/scoreboard.json
```
**問題：** Access Denied - 無法訪問歷史日期的 scoreboard

#### 2. **stats.nba.com API**
```
https://stats.nba.com/stats/leaguegamefinder?...
```
**問題：**
- 需要複雜的認證和 headers
- 經常超時
- 不穩定

### 🔧 實作方案

#### 方案一：使用完整賽程 API（推薦）

**優點：**
- ✅ 穩定可靠
- ✅ 一次請求獲取所有資料
- ✅ 包含完整的戰績資訊
- ✅ 無需認證

**實作步驟：**
```go
// 1. 獲取完整賽程
resp, _ := http.Get("https://cdn.nba.com/static/json/staticData/scheduleLeagueV2_9.json")

// 2. 解析 JSON
var schedule ScheduleLeague
json.Unmarshal(body, &schedule)

// 3. 過濾特定球隊的比賽
for _, dateInfo := range schedule.LeagueSchedule.GameDates {
    for _, game := range dateInfo.Games {
        if game.HomeTeam.TeamID == targetTeamID || game.AwayTeam.TeamID == targetTeamID {
            // 處理比賽資料
        }
    }
}

// 4. 只取已結束的比賽 (gameStatus == 3)
// 5. 根據分數判斷輸贏
// 6. 限制返回數量（例如最近 10 場）
```

**測試結果：**
```
Golden State Warriors 最近戰績:
1. [W] 10/05/2025 主場 vs Los Angeles Lakers (0-2) - 111-103
2. [W] 10/08/2025 主場 vs Portland Trail Blazers (0-1) - 129-123
戰績: 2 勝 0 敗 (勝率 100.0%)

Los Angeles Lakers 最近戰績:
1. [L] 10/03/2025 主場 vs Phoenix Suns (1-0) - 81-103
2. [L] 10/05/2025 客場 vs Golden State Warriors (1-0) - 103-111
戰績: 0 勝 2 敗 (勝率 0.0%)
```

#### 方案二：結合多個 API

**適用場景：** 需要更詳細的比賽數據

```go
// 1. 從完整賽程 API 獲取 gameId
scheduleAPI := "https://cdn.nba.com/static/json/staticData/scheduleLeagueV2_9.json"

// 2. 根據 gameId 獲取詳細比賽數據
boxscoreAPI := fmt.Sprintf(
    "https://cdn.nba.com/static/json/liveData/boxscore/boxscore_%s.json",
    gameId
)
```

### 📊 可獲取的資料

使用完整賽程 API 可以獲取：

| 資料項目 | 可獲取 | 說明 |
|---------|--------|------|
| 比賽日期 | ✅ | gameDate |
| 對手球隊 | ✅ | homeTeam/awayTeam |
| 比分 | ✅ | score (已結束比賽) |
| 輸贏結果 | ✅ | 根據分數判斷 |
| 對手戰績 | ✅ | wins/losses |
| 主客場 | ✅ | 根據 teamId 判斷 |
| 比賽狀態 | ✅ | gameStatus/gameStatusText |
| 賽季類型 | ✅ | gameLabel (季前賽/例行賽/季後賽) |

### 🚀 建議實作

**1. 創建歷史戰績查詢功能：**
```go
// internal/crawler/history.go
func FetchTeamHistory(teamID int, limit int) ([]TeamGame, error) {
    // 實作邏輯見 test_history_api.go
}
```

**2. 整合到現有 API：**
```go
// internal/logic/api.go
func GetTeamRecentGames(teamID int) (*TeamHistoryResponse, error) {
    games, err := crawler.FetchTeamHistory(teamID, 10)
    // 轉換為 API 回應格式
}
```

**3. 添加 Web API 端點：**
```go
// GET /api/team/{teamId}/history?limit=10
router.GET("/api/team/:teamId/history", server.GetTeamHistory)
```

### 📝 球隊 ID 對照表

常用球隊 ID：
- Golden State Warriors: 1610612744
- Los Angeles Lakers: 1610612747
- Boston Celtics: 1610612738
- Brooklyn Nets: 1610612751
- Phoenix Suns: 1610612756
- Miami Heat: 1610612748
- (完整列表請參考 internal/models/teammap.go)

### ⚠️ 注意事項

1. **API 版本：** scheduleLeagueV2_9.json 中的版本號可能會變更
2. **快取策略：** 建議快取賽程資料，避免頻繁請求
3. **錯誤處理：** API 可能返回空資料或錯誤狀態
4. **賽季更新：** 新賽季開始時 API 內容會更新

### 🔗 相關 API 文檔

- nba_api (Python): https://github.com/swar/nba_api
- NBA Stats 非官方文檔: https://github.com/swar/nba_api/tree/master/docs
- SportsDataIO NBA API: https://sportsdata.io/developers/api-documentation/nba

### 📈 替代方案（如果官方 API 不可用）

1. **第三方 API：**
   - RapidAPI - API-NBA: https://rapidapi.com/api-sports/api/api-nba
   - SportsDataIO: https://sportsdata.io/
   - BallDontLie: https://www.balldontlie.io/

2. **網頁爬蟲：**
   - ESPN NBA Scores
   - NBA.com 官網
   - Basketball Reference

## 結論

✅ **NBA 官方 API 可以查詢球隊歷史戰績**

**最佳方案：** 使用 `scheduleLeagueV2_9.json` API
- 包含完整賽季資料
- 穩定可靠
- 無需認證
- 一次請求獲取所有所需資料

**實作已驗證：**
- 測試代碼：`test_history_api.go`
- 成功獲取勇士、湖人、賽爾提克的最近戰績
- 包含日期、對手、比分、勝負等完整資訊

**下一步建議：**
1. 將功能整合到現有專案
2. 添加快取機制
3. 創建 Web API 端點
4. 添加前端顯示介面
