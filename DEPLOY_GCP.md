# GCP 部署指南 - VM + Docker 最便宜方案

本指南提供 **最便宜的 GCP 部署方案**，適合使用免費試用帳號，並每半年輪換以避免產生費用。

## 📋 方案概覽

使用 **Compute Engine (VM) + Docker**，優勢：
- ✅ 使用 e2-micro 實例（免費額度內）
- ✅ 完全控制部署環境
- ✅ 適合長時間運行的服務
- ✅ 輪換帳號策略避免收費

**預估成本**：每月 $0（在免費額度內）

---

## 🎯 部署策略

### 免費額度說明
每個新的 GCP 帳號提供：
1. **$300 美金試用額度**（90 天內有效）
2. **永久免費額度**：
   - 1 個 e2-micro VM 實例（美國三個區域之一）
   - 30 GB 標準硬碟
   - 1 GB 網路流量（每月）

### 半年輪換策略
1. 使用新 Gmail 帳號註冊 GCP
2. 綁定信用卡（不會自動扣款，需手動升級）
3. 使用 6 個月後準備新帳號
4. 匯出 Docker 映像，在新帳號重新部署

---

## 🚀 完整部署流程

### 步驟 1：準備 GCP 帳號

1. **註冊新的 Gmail 帳號**（如果需要輪換）
   - 前往 [Gmail](https://mail.google.com)
   - 註冊新帳號

2. **註冊 GCP 免費試用**
   - 前往 [GCP Console](https://console.cloud.google.com)
   - 選擇「免費試用」
   - 填寫信用卡資訊（僅驗證用，不會自動扣款）

3. **確認不會自動收費**
   - GCP 預設不會在試用期結束後自動收費
   - 需手動升級為付費帳戶才會開始計費
   - 設定預算警報（下文說明）

---

### 步驟 2：安裝 Google Cloud SDK（本地電腦）

#### macOS
```bash
brew install --cask google-cloud-sdk
```

#### Linux
```bash
curl https://sdk.cloud.google.com | bash
exec -l $SHELL
```

#### Windows
下載安裝器：https://cloud.google.com/sdk/docs/install

#### 登入 GCP
```bash
# 登入
gcloud auth login

# 設定專案 ID（自訂，例如：nba-scanner-2025）
gcloud config set project nba-scanner-2025

# 創建專案
gcloud projects create nba-scanner-2025 --name="NBA Scanner"

# 列出所有專案
gcloud projects list
```

---

### 步驟 3：創建最便宜的 VM 實例

```bash
# 創建 e2-micro 實例（最便宜，符合免費額度）
gcloud compute instances create nba-scanner-vm \
  --zone=us-central1-a \
  --machine-type=e2-micro \
  --image-family=ubuntu-2004-lts \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=10GB \
  --boot-disk-type=pd-standard \
  --tags=http-server,nba-scanner
```

**重要參數說明**：
- `--zone=us-central1-a`：美國中部（免費額度適用區域）
  - 其他免費區域：`us-west1`、`us-east1`
  - ⚠️ 避免使用 asia-east1（台灣），不在免費額度內
- `--machine-type=e2-micro`：最小實例（免費）
  - 2 個共享 vCPU
  - 1 GB 記憶體
- `--boot-disk-size=10GB`：最小硬碟（免費額度 30GB 內）
- `--boot-disk-type=pd-standard`：標準硬碟（最便宜）

---

### 步驟 4：設定防火牆規則（開放 8081 端口）

```bash
# 允許 8081 端口的 HTTP 流量
gcloud compute firewall-rules create allow-nba-scanner \
  --allow tcp:8081 \
  --target-tags nba-scanner \
  --description="Allow NBA Scanner web traffic on port 8081"

# 查看防火牆規則
gcloud compute firewall-rules list
```

---

### 步驟 5：連線到 VM

```bash
# SSH 連線到 VM
gcloud compute ssh nba-scanner-vm --zone=us-central1-a
```

成功連線後，會進入 VM 的終端機。

---

### 步驟 6：在 VM 上安裝 Docker

```bash
# 更新套件清單
sudo apt-get update

# 安裝 Docker
sudo apt-get install -y docker.io

# 啟動 Docker 服務
sudo systemctl start docker
sudo systemctl enable docker

# 將當前使用者加入 docker 群組（避免每次都要 sudo）
sudo usermod -aG docker $USER

# 驗證 Docker 安裝
docker --version
```

**重要**：執行 `usermod` 後需要重新登入 SSH 才會生效：
```bash
# 登出
exit

# 重新登入
gcloud compute ssh nba-scanner-vm --zone=us-central1-a
```

---

### 步驟 7：上傳專案到 VM

有兩種方式：

#### 方式 A：使用 gcloud scp（推薦）

**在本地電腦執行**：
```bash
# 進入專案目錄
cd /path/to/scanNBA

# 壓縮專案（排除不必要的檔案）
tar -czf nba-scanner.tar.gz \
  --exclude='node_modules' \
  --exclude='.git' \
  --exclude='*.log' \
  --exclude='nba-scan' \
  --exclude='nba-scanner' \
  .

# 上傳到 VM
gcloud compute scp nba-scanner.tar.gz nba-scanner-vm:~ --zone=us-central1-a
```

**在 VM 上執行**：
```bash
# 解壓縮
mkdir -p ~/scanNBA
tar -xzf nba-scanner.tar.gz -C ~/scanNBA

# 進入目錄
cd ~/scanNBA
```

#### 方式 B：使用 Git（如果專案在 GitHub）

**在 VM 上執行**：
```bash
# 安裝 Git
sudo apt-get install -y git

# Clone 專案
git clone https://github.com/YOUR_USERNAME/scanNBA.git
cd scanNBA
```

---

### 步驟 8：在 VM 上建置 Docker 映像

**在 VM 上執行**：
```bash
# 確認在專案目錄
cd ~/scanNBA

# 檢查 Dockerfile 是否存在
ls -l Dockerfile

# 建置 Docker 映像
docker build -t nba-scanner:latest .

# 查看映像
docker images
```

建置過程需要 5-10 分鐘（視網路速度）。

---

### 步驟 9：運行 Docker 容器

```bash
# 運行容器（背景執行）
docker run -d \
  --name nba-scanner \
  --restart unless-stopped \
  -p 8081:8081 \
  -e TZ=Asia/Taipei \
  nba-scanner:latest

# 查看容器狀態
docker ps

# 查看日誌
docker logs -f nba-scanner

# 停止查看日誌：按 Ctrl+C
```

**參數說明**：
- `-d`：背景執行
- `--name nba-scanner`：容器名稱
- `--restart unless-stopped`：自動重啟（除非手動停止）
- `-p 8081:8081`：端口映射（主機:容器）
- `-e TZ=Asia/Taipei`：設定時區

---

### 步驟 10：取得外部 IP 並測試

**在本地電腦執行**：
```bash
# 取得 VM 的外部 IP
gcloud compute instances list

# 輸出範例：
# NAME              ZONE           MACHINE_TYPE  EXTERNAL_IP
# nba-scanner-vm    us-central1-a  e2-micro      34.123.45.67
```

**測試訪問**：
```bash
# 使用 curl 測試
curl http://34.123.45.67:8081

# 或在瀏覽器開啟
# http://34.123.45.67:8081
```

🎉 **部署完成！** 現在可以使用外部 IP 訪問服務了。

---

## 🔧 常用管理指令

### Docker 容器管理

```bash
# 查看運行中的容器
docker ps

# 查看所有容器（包含停止的）
docker ps -a

# 停止容器
docker stop nba-scanner

# 啟動容器
docker start nba-scanner

# 重啟容器
docker restart nba-scanner

# 刪除容器
docker rm -f nba-scanner

# 查看容器日誌
docker logs nba-scanner

# 即時查看日誌
docker logs -f nba-scanner

# 進入容器內部
docker exec -it nba-scanner /bin/sh
```

### VM 管理

```bash
# 停止 VM（不會刪除，但停止計費）
gcloud compute instances stop nba-scanner-vm --zone=us-central1-a

# 啟動 VM
gcloud compute instances start nba-scanner-vm --zone=us-central1-a

# 重啟 VM
gcloud compute instances reset nba-scanner-vm --zone=us-central1-a

# 刪除 VM（⚠️ 會永久刪除）
gcloud compute instances delete nba-scanner-vm --zone=us-central1-a

# 查看 VM 狀態
gcloud compute instances describe nba-scanner-vm --zone=us-central1-a
```

---

## 🔄 更新應用程式

當程式碼有變更時：

### 方法 1：完整重建（推薦）

```bash
# 在本地電腦：重新打包並上傳
cd /path/to/scanNBA
tar -czf nba-scanner.tar.gz \
  --exclude='node_modules' \
  --exclude='.git' \
  --exclude='*.log' \
  --exclude='nba-scan' \
  --exclude='nba-scanner' \
  .
gcloud compute scp nba-scanner.tar.gz nba-scanner-vm:~ --zone=us-central1-a

# 在 VM 上：更新並重建
gcloud compute ssh nba-scanner-vm --zone=us-central1-a

# 解壓縮新版本
cd ~/scanNBA
rm -rf *
tar -xzf ~/nba-scanner.tar.gz

# 停止舊容器
docker stop nba-scanner
docker rm nba-scanner

# 重建映像
docker build -t nba-scanner:latest .

# 啟動新容器
docker run -d \
  --name nba-scanner \
  --restart unless-stopped \
  -p 8081:8081 \
  -e TZ=Asia/Taipei \
  nba-scanner:latest

# 查看日誌確認運行
docker logs -f nba-scanner
```

### 方法 2：使用 Git（如果專案在 GitHub）

```bash
# 在 VM 上執行
cd ~/scanNBA
git pull origin main

# 重建並重啟
docker stop nba-scanner
docker rm nba-scanner
docker build -t nba-scanner:latest .
docker run -d \
  --name nba-scanner \
  --restart unless-stopped \
  -p 8081:8081 \
  -e TZ=Asia/Taipei \
  nba-scanner:latest
```

---

## 💰 成本監控與預算警報

### 設定預算警報（避免意外收費）

1. **前往 GCP Console**
   - https://console.cloud.google.com/billing

2. **建立預算警報**
   ```bash
   # 或使用 CLI 建立
   gcloud billing budgets create \
     --billing-account=YOUR_BILLING_ACCOUNT_ID \
     --display-name="NBA Scanner Budget Alert" \
     --budget-amount=5USD \
     --threshold-rule=percent=50 \
     --threshold-rule=percent=90 \
     --threshold-rule=percent=100
   ```

3. **設定 Email 通知**
   - 在 Console 中設定當達到預算 50%/90%/100% 時發送 Email

### 查看當前費用

```bash
# 查看當前費用
gcloud billing accounts list

# 在 Console 查看詳細費用
# https://console.cloud.google.com/billing/reports
```

### 免費額度使用狀況

- **查看免費額度**：https://console.cloud.google.com/billing/freetrial
- **查看 VM 運行時間**：
  ```bash
  gcloud compute instances list --format="table(name,status,creationTimestamp)"
  ```

---

## 🔄 半年輪換策略

當接近 6 個月或試用額度用完時：

### 步驟 1：匯出 Docker 映像

**在舊 VM 上執行**：
```bash
# 匯出 Docker 映像為檔案
docker save nba-scanner:latest | gzip > nba-scanner-image.tar.gz

# 下載到本地
# 在本地電腦執行
gcloud compute scp nba-scanner-vm:~/nba-scanner-image.tar.gz . --zone=us-central1-a
```

### 步驟 2：註冊新的 GCP 帳號

1. 使用新的 Gmail 帳號
2. 註冊 GCP 免費試用
3. 綁定信用卡

### 步驟 3：在新帳號部署

1. 按照「完整部署流程」重新執行步驟 2-10
2. 或上傳已匯出的 Docker 映像：
   ```bash
   # 上傳映像檔到新 VM
   gcloud compute scp nba-scanner-image.tar.gz nba-scanner-vm:~ --zone=us-central1-a

   # 在新 VM 上載入映像
   docker load < nba-scanner-image.tar.gz

   # 運行容器
   docker run -d \
     --name nba-scanner \
     --restart unless-stopped \
     -p 8081:8081 \
     -e TZ=Asia/Taipei \
     nba-scanner:latest
   ```

### 步驟 4：刪除舊帳號資源（可選）

**⚠️ 確認新環境運行正常後再刪除**：
```bash
# 在舊帳號刪除 VM
gcloud compute instances delete nba-scanner-vm --zone=us-central1-a

# 刪除防火牆規則
gcloud compute firewall-rules delete allow-nba-scanner

# 關閉計費（在 Console 操作）
# https://console.cloud.google.com/billing
```

---

## 🔒 安全性建議

### 1. 限制 SSH 訪問

```bash
# 只允許特定 IP 訪問 SSH
gcloud compute firewall-rules create allow-ssh-from-my-ip \
  --allow tcp:22 \
  --source-ranges YOUR_IP_ADDRESS/32 \
  --target-tags nba-scanner

# 刪除預設的 SSH 規則（可選）
gcloud compute firewall-rules delete default-allow-ssh
```

### 2. 啟用自動更新

```bash
# 在 VM 上啟用自動安全更新
sudo apt-get install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

### 3. 定期更新 Docker 映像

```bash
# 更新基礎映像（在 Dockerfile 中使用固定版本）
# 定期重建映像以獲取安全更新
docker build --no-cache -t nba-scanner:latest .
```

### 4. 設定 HTTPS（可選，使用 Caddy）

如果需要 HTTPS，推薦使用 Caddy（自動申請 Let's Encrypt 證書）：

```bash
# 在 VM 上安裝 Caddy
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install -y caddy

# 設定 Caddyfile
sudo tee /etc/caddy/Caddyfile <<EOF
your-domain.com {
    reverse_proxy localhost:8081
}
EOF

# 重啟 Caddy
sudo systemctl restart caddy

# 開放 443 端口
gcloud compute firewall-rules create allow-https \
  --allow tcp:443 \
  --target-tags nba-scanner
```

---

## 🐛 故障排除

### 問題 1：容器無法啟動

```bash
# 查看詳細日誌
docker logs nba-scanner

# 檢查端口佔用
sudo lsof -i :8081

# 手動測試程式
cd ~/scanNBA
go run main.go --server --port 8081
```

### 問題 2：無法從外部訪問

```bash
# 檢查防火牆規則
gcloud compute firewall-rules list

# 檢查容器是否運行
docker ps

# 檢查端口映射
docker port nba-scanner

# 測試本地訪問
curl http://localhost:8081
```

### 問題 3：VM 記憶體不足

```bash
# 查看記憶體使用
free -h

# 查看 Docker 容器資源使用
docker stats nba-scanner

# 重啟容器
docker restart nba-scanner
```

### 問題 4：超過免費額度

- **確認區域**：必須在 us-central1/us-west1/us-east1
- **確認實例類型**：必須是 e2-micro
- **確認只有一個實例**：`gcloud compute instances list`
- **確認硬碟大小**：≤ 30GB

### 問題 5：Docker 建置失敗

```bash
# 清理 Docker 快取
docker system prune -a

# 確認硬碟空間
df -h

# 重新建置
docker build --no-cache -t nba-scanner:latest .
```

---

## 📊 效能優化

### 1. 啟用 Docker 日誌輪替

```bash
# 設定 Docker daemon
sudo tee /etc/docker/daemon.json <<EOF
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF

# 重啟 Docker
sudo systemctl restart docker

# 重啟容器
docker stop nba-scanner && docker start nba-scanner
```

### 2. 設定 VM 自動重啟（預防當機）

```bash
# 啟用自動重啟
gcloud compute instances update nba-scanner-vm \
  --zone=us-central1-a \
  --restart-on-failure
```

### 3. 監控服務健康狀態

建立簡單的監控腳本：

```bash
# 在 VM 上建立監控腳本
cat > ~/monitor.sh <<'EOF'
#!/bin/bash
if ! docker ps | grep -q nba-scanner; then
    echo "Container down, restarting..."
    docker start nba-scanner
fi
EOF

chmod +x ~/monitor.sh

# 設定 cron job（每 5 分鐘檢查一次）
(crontab -l 2>/dev/null; echo "*/5 * * * * ~/monitor.sh") | crontab -
```

---

## 📚 參考資源

- [GCP 免費額度說明](https://cloud.google.com/free)
- [Compute Engine 定價](https://cloud.google.com/compute/pricing)
- [Docker 官方文件](https://docs.docker.com/)
- [GCP SDK 文件](https://cloud.google.com/sdk/docs)
- [防火牆規則設定](https://cloud.google.com/vpc/docs/firewalls)

---

## ✅ 部署檢查清單

- [ ] 註冊 GCP 免費試用帳號
- [ ] 設定預算警報（$5 USD）
- [ ] 創建 e2-micro VM（us-central1-a）
- [ ] 設定防火牆規則（8081 端口）
- [ ] 安裝 Docker
- [ ] 上傳專案並建置映像
- [ ] 運行 Docker 容器
- [ ] 測試外部訪問
- [ ] 設定自動重啟
- [ ] 設定監控腳本
- [ ] 記錄外部 IP 地址
- [ ] 標記 6 個月輪換日期

---

**祝部署順利！** 🚀

如遇到問題，請參考「故障排除」章節或檢查 GCP Console 的日誌。
