# GCP 部署指南

本指南將引導您在 Google Cloud Platform (GCP) 上部署 NBA Scanner。

## 部署方案概覽

推薦使用 **Cloud Run**，原因：
- ✅ 全託管服務，無需管理伺服器
- ✅ 自動擴展（流量大時自動增加實例）
- ✅ 按使用量計費（無流量時接近零成本）
- ✅ 支援容器部署
- ✅ 內建 HTTPS

替代方案：Compute Engine (VM)，適合需要更多控制權的情況。

---

## 方案一：Cloud Run 部署（推薦）

### 前置準備

1. **安裝 Google Cloud SDK**
   ```bash
   # macOS
   brew install --cask google-cloud-sdk

   # 其他系統請參考：https://cloud.google.com/sdk/docs/install
   ```

2. **登入 GCP**
   ```bash
   gcloud auth login
   ```

3. **創建專案（如果還沒有）**
   ```bash
   # 創建新專案
   gcloud projects create nba-scanner-project --name="NBA Scanner"

   # 設定為當前專案
   gcloud config set project nba-scanner-project

   # 啟用計費（需要在 GCP Console 中操作）
   # https://console.cloud.google.com/billing
   ```

4. **啟用必要的 API**
   ```bash
   gcloud services enable cloudbuild.googleapis.com
   gcloud services enable run.googleapis.com
   gcloud services enable artifactregistry.googleapis.com
   ```

### 部署步驟

#### 步驟 1：建置並推送 Docker 映像

```bash
# 進入專案目錄
cd /path/to/scanNBA

# 設定變數
PROJECT_ID=$(gcloud config get-value project)
REGION=asia-east1  # 台灣機房

# 建置映像並推送到 Artifact Registry
gcloud builds submit --tag gcr.io/$PROJECT_ID/nba-scanner
```

#### 步驟 2：部署到 Cloud Run

```bash
gcloud run deploy nba-scanner \
  --image gcr.io/$PROJECT_ID/nba-scanner \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 5 \
  --timeout 60s \
  --set-env-vars TZ=Asia/Taipei
```

#### 步驟 3：取得服務 URL

部署完成後會顯示服務 URL，格式為：
```
https://nba-scanner-xxxxx-xx.a.run.app
```

訪問此 URL 即可使用服務！

### 設定自訂域名（可選）

```bash
# 1. 在 Cloud Run 中對應域名
gcloud run domain-mappings create \
  --service nba-scanner \
  --domain nba.yourdomain.com \
  --region $REGION

# 2. 根據提示在 DNS 提供商處設定記錄
```

### 更新服務

當程式碼有變更時：

```bash
# 重新建置並部署
gcloud builds submit --tag gcr.io/$PROJECT_ID/nba-scanner
gcloud run deploy nba-scanner \
  --image gcr.io/$PROJECT_ID/nba-scanner \
  --region $REGION
```

### 監控與日誌

```bash
# 查看日誌
gcloud run services logs read nba-scanner --region $REGION

# 即時查看日誌
gcloud run services logs tail nba-scanner --region $REGION

# 在 GCP Console 查看
# https://console.cloud.google.com/run
```

### 成本估算

Cloud Run 免費額度（每月）：
- 200 萬次請求
- 36 萬 GB-秒
- 18 萬 vCPU-秒

預估成本（低流量）：**$0 - $5/月**

---

## 方案二：Compute Engine (VM) 部署

### 步驟 1：創建 VM 實例

```bash
# 創建 e2-micro 實例（免費額度）
gcloud compute instances create nba-scanner-vm \
  --zone=asia-east1-b \
  --machine-type=e2-micro \
  --image-family=cos-stable \
  --image-project=cos-cloud \
  --boot-disk-size=10GB \
  --tags=http-server,https-server

# 允許 HTTP/HTTPS 流量
gcloud compute firewall-rules create allow-http-8080 \
  --allow tcp:8080 \
  --target-tags http-server
```

### 步驟 2：連線到 VM

```bash
gcloud compute ssh nba-scanner-vm --zone=asia-east1-b
```

### 步驟 3：在 VM 上安裝 Docker

```bash
# 安裝 Docker
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# 重新登入以套用權限
exit
gcloud compute ssh nba-scanner-vm --zone=asia-east1-b
```

### 步驟 4：部署應用

```bash
# 拉取映像（從本地推送或使用 gcr.io）
docker pull gcr.io/YOUR_PROJECT_ID/nba-scanner

# 運行容器
docker run -d \
  --name nba-scanner \
  --restart unless-stopped \
  -p 8080:8080 \
  -e TZ=Asia/Taipei \
  gcr.io/YOUR_PROJECT_ID/nba-scanner

# 查看日誌
docker logs -f nba-scanner
```

### 步驟 5：取得外部 IP

```bash
gcloud compute instances list
```

訪問 `http://EXTERNAL_IP:8080`

### 設定 HTTPS（使用 Nginx + Let's Encrypt）

```bash
# 安裝 Nginx
sudo yum install -y nginx certbot python3-certbot-nginx

# 設定 Nginx 反向代理
sudo tee /etc/nginx/conf.d/nba-scanner.conf <<EOF
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF

# 啟動 Nginx
sudo systemctl start nginx
sudo systemctl enable nginx

# 申請 SSL 證書
sudo certbot --nginx -d your-domain.com
```

### 成本估算

e2-micro 實例（免費額度）：
- 每月 730 小時免費（1 個實例）
- 30 GB 標準硬碟

預估成本：**$0/月**（在免費額度內）

---

## 快速部署腳本

創建 `deploy.sh` 自動化部署：

```bash
#!/bin/bash

# 設定變數
PROJECT_ID=$(gcloud config get-value project)
REGION="asia-east1"
SERVICE_NAME="nba-scanner"

echo "🚀 開始部署 NBA Scanner..."

# 建置並推送映像
echo "📦 建置 Docker 映像..."
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME

# 部署到 Cloud Run
echo "☁️  部署到 Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 5 \
  --set-env-vars TZ=Asia/Taipei

echo "✅ 部署完成！"
echo "📍 服務 URL："
gcloud run services describe $SERVICE_NAME --region $REGION --format="value(status.url)"
```

使用方法：
```bash
chmod +x deploy.sh
./deploy.sh
```

---

## 故障排除

### 問題 1：部署失敗 - 權限不足

```bash
# 確認計費已啟用
gcloud beta billing projects describe PROJECT_ID

# 授予必要權限
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="user:YOUR_EMAIL" \
  --role="roles/run.admin"
```

### 問題 2：容器啟動失敗

```bash
# 查看詳細日誌
gcloud run services logs read nba-scanner --region asia-east1 --limit 50

# 本地測試 Docker 映像
docker run -p 8080:8080 gcr.io/PROJECT_ID/nba-scanner
```

### 問題 3：無法訪問服務

```bash
# 確認服務允許未驗證訪問
gcloud run services add-iam-policy-binding nba-scanner \
  --region asia-east1 \
  --member="allUsers" \
  --role="roles/run.invoker"
```

### 問題 4：超過免費額度

- Cloud Run: 監控請求數，考慮設定 `--min-instances 0`
- VM: 確保只運行 1 個 e2-micro 實例

---

## 安全建議

1. **啟用 Cloud Armor**（防 DDoS）
   ```bash
   gcloud compute security-policies create nba-scanner-policy
   ```

2. **設定 CORS**（如果有前端分離）
   - 在 `server.go` 中加入 CORS 中介軟體

3. **監控異常流量**
   - 在 GCP Console 設定警報

4. **定期更新依賴**
   ```bash
   go get -u ./...
   go mod tidy
   ```

---

## 持續部署（CI/CD）

使用 Cloud Build 自動部署：

創建 `cloudbuild.yaml`：

```yaml
steps:
  # 建置 Docker 映像
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/nba-scanner', '.']

  # 推送到 Container Registry
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/nba-scanner']

  # 部署到 Cloud Run
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    entrypoint: gcloud
    args:
      - 'run'
      - 'deploy'
      - 'nba-scanner'
      - '--image'
      - 'gcr.io/$PROJECT_ID/nba-scanner'
      - '--region'
      - 'asia-east1'
      - '--platform'
      - 'managed'
      - '--allow-unauthenticated'

images:
  - 'gcr.io/$PROJECT_ID/nba-scanner'
```

設定 GitHub 自動部署：
```bash
gcloud builds triggers create github \
  --repo-name=YOUR_REPO \
  --repo-owner=YOUR_GITHUB_USERNAME \
  --branch-pattern="^main$" \
  --build-config=cloudbuild.yaml
```

---

## 參考資源

- [Cloud Run 官方文件](https://cloud.google.com/run/docs)
- [Compute Engine 定價](https://cloud.google.com/compute/pricing)
- [GCP 免費額度](https://cloud.google.com/free)
- [Container Registry 文件](https://cloud.google.com/container-registry/docs)
