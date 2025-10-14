#!/bin/bash

# NBA Scanner GCP VM 部署腳本
# 使用方法: ./deploy.sh [VM_NAME] [ZONE]
# 範例: ./deploy.sh nba-scanner-vm us-central1-a

set -e  # 遇到錯誤立即退出

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 設定變數（可從參數覆蓋）
VM_NAME="${1:-nba-scanner-vm}"
ZONE="${2:-us-central1-a}"
PROJECT_ID=$(gcloud config get-value project 2>/dev/null)

echo -e "${GREEN}🚀 NBA Scanner VM 部署腳本${NC}"
echo "================================"
echo "專案 ID: $PROJECT_ID"
echo "VM 名稱: $VM_NAME"
echo "區域: $ZONE"
echo "================================"
echo ""

# 檢查是否已登入
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" > /dev/null 2>&1; then
    echo -e "${RED}❌ 請先登入 GCP: gcloud auth login${NC}"
    exit 1
fi

# 檢查專案 ID
if [ -z "$PROJECT_ID" ]; then
    echo -e "${RED}❌ 請先設定 GCP 專案: gcloud config set project YOUR_PROJECT_ID${NC}"
    exit 1
fi

# 檢查 VM 是否存在
echo -e "${BLUE}🔍 檢查 VM 狀態...${NC}"
if ! gcloud compute instances describe $VM_NAME --zone=$ZONE > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  VM 不存在，是否要創建新 VM？${NC}"
    echo "這將創建一個 e2-micro 實例（免費額度內）"
    read -p "創建新 VM? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${GREEN}📦 創建 VM 實例...${NC}"
        gcloud compute instances create $VM_NAME \
          --zone=$ZONE \
          --machine-type=e2-micro \
          --image-family=ubuntu-2004-lts \
          --image-project=ubuntu-os-cloud \
          --boot-disk-size=10GB \
          --boot-disk-type=pd-standard \
          --tags=http-server,nba-scanner

        echo -e "${GREEN}🔥 設定防火牆規則...${NC}"
        if ! gcloud compute firewall-rules describe allow-nba-scanner > /dev/null 2>&1; then
            gcloud compute firewall-rules create allow-nba-scanner \
              --allow tcp:8081 \
              --target-tags nba-scanner \
              --description="Allow NBA Scanner web traffic on port 8081"
        fi

        echo -e "${GREEN}⏳ 等待 VM 啟動...${NC}"
        sleep 30
    else
        echo -e "${RED}❌ 部署已取消${NC}"
        exit 1
    fi
fi

# 確認是否繼續部署
echo ""
read -p "是否繼續部署到 VM? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}⏹️  部署已取消${NC}"
    exit 0
fi

# 步驟 1: 打包專案
echo -e "\n${GREEN}📦 步驟 1/5: 打包專案檔案...${NC}"
tar -czf nba-scanner.tar.gz \
  --exclude='node_modules' \
  --exclude='.git' \
  --exclude='*.log' \
  --exclude='nba-scan' \
  --exclude='nba-scanner' \
  --exclude='nba-scanner.tar.gz' \
  .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 專案打包成功${NC}"
else
    echo -e "${RED}❌ 專案打包失敗${NC}"
    exit 1
fi

# 步驟 2: 上傳到 VM
echo -e "\n${GREEN}📤 步驟 2/5: 上傳專案到 VM...${NC}"
gcloud compute scp nba-scanner.tar.gz $VM_NAME:~ --zone=$ZONE

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 上傳成功${NC}"
    rm nba-scanner.tar.gz
else
    echo -e "${RED}❌ 上傳失敗${NC}"
    rm nba-scanner.tar.gz
    exit 1
fi

# 步驟 3: 在 VM 上安裝 Docker（如果尚未安裝）
echo -e "\n${GREEN}🐳 步驟 3/5: 檢查並安裝 Docker...${NC}"
gcloud compute ssh $VM_NAME --zone=$ZONE --command="
    if ! command -v docker &> /dev/null; then
        echo '安裝 Docker...'
        sudo apt-get update
        sudo apt-get install -y docker.io
        sudo systemctl start docker
        sudo systemctl enable docker
        sudo usermod -aG docker \$USER
        echo 'Docker 安裝完成'
    else
        echo 'Docker 已安裝'
    fi
"

# 步驟 4: 解壓並建置 Docker 映像
echo -e "\n${GREEN}🏗️  步驟 4/5: 建置 Docker 映像...${NC}"
gcloud compute ssh $VM_NAME --zone=$ZONE --command="
    # 解壓縮
    mkdir -p ~/scanNBA
    tar -xzf ~/nba-scanner.tar.gz -C ~/scanNBA

    # 建置 Docker 映像
    cd ~/scanNBA
    sudo docker build -t nba-scanner:latest .

    echo 'Docker 映像建置完成'
"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Docker 映像建置成功${NC}"
else
    echo -e "${RED}❌ Docker 映像建置失敗${NC}"
    exit 1
fi

# 步驟 5: 停止舊容器並啟動新容器
echo -e "\n${GREEN}🚀 步驟 5/5: 部署容器...${NC}"
gcloud compute ssh $VM_NAME --zone=$ZONE --command="
    # 停止並刪除舊容器（如果存在）
    if sudo docker ps -a | grep -q nba-scanner; then
        echo '停止舊容器...'
        sudo docker stop nba-scanner 2>/dev/null || true
        sudo docker rm nba-scanner 2>/dev/null || true
    fi

    # 啟動新容器
    echo '啟動新容器...'
    sudo docker run -d \
      --name nba-scanner \
      --restart unless-stopped \
      -p 8081:8081 \
      -e TZ=Asia/Taipei \
      nba-scanner:latest

    # 等待容器啟動
    sleep 3

    # 檢查容器狀態
    if sudo docker ps | grep -q nba-scanner; then
        echo '✅ 容器運行中'
        sudo docker ps | grep nba-scanner
    else
        echo '❌ 容器啟動失敗，查看日誌:'
        sudo docker logs nba-scanner
        exit 1
    fi
"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 容器部署成功${NC}"
else
    echo -e "${RED}❌ 容器部署失敗${NC}"
    exit 1
fi

# 取得外部 IP
echo -e "\n${GREEN}📍 取得服務資訊...${NC}"
EXTERNAL_IP=$(gcloud compute instances describe $VM_NAME --zone=$ZONE --format="get(networkInterfaces[0].accessConfigs[0].natIP)")

echo ""
echo "================================"
echo -e "${GREEN}✅ 部署完成！${NC}"
echo "================================"
echo -e "服務 URL: ${YELLOW}http://$EXTERNAL_IP:8081${NC}"
echo ""
echo "📊 查看日誌:"
echo "  gcloud compute ssh $VM_NAME --zone=$ZONE --command='sudo docker logs -f nba-scanner'"
echo ""
echo "🔄 重啟服務:"
echo "  gcloud compute ssh $VM_NAME --zone=$ZONE --command='sudo docker restart nba-scanner'"
echo ""
echo "🗑️  停止服務:"
echo "  gcloud compute ssh $VM_NAME --zone=$ZONE --command='sudo docker stop nba-scanner'"
echo ""
echo "💻 SSH 連線:"
echo "  gcloud compute ssh $VM_NAME --zone=$ZONE"
echo ""
echo "🔥 刪除 VM:"
echo "  gcloud compute instances delete $VM_NAME --zone=$ZONE"
echo "================================"
