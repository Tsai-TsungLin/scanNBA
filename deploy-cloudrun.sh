#!/bin/bash

# NBA Scanner GCP 部署腳本
# 使用方法: ./deploy.sh

set -e  # 遇到錯誤立即退出

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 設定變數
PROJECT_ID=$(gcloud config get-value project)
REGION="asia-east1"
SERVICE_NAME="nba-scanner"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo -e "${GREEN}🚀 NBA Scanner 部署腳本${NC}"
echo "================================"
echo "專案 ID: $PROJECT_ID"
echo "區域: $REGION"
echo "服務名稱: $SERVICE_NAME"
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

# 確認是否繼續
read -p "是否繼續部署? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}⏹️  部署已取消${NC}"
    exit 0
fi

# 步驟 1: 建置 Docker 映像
echo -e "\n${GREEN}📦 步驟 1/3: 建置 Docker 映像...${NC}"
gcloud builds submit --tag $IMAGE_NAME

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 映像建置成功${NC}"
else
    echo -e "${RED}❌ 映像建置失敗${NC}"
    exit 1
fi

# 步驟 2: 部署到 Cloud Run
echo -e "\n${GREEN}☁️  步驟 2/3: 部署到 Cloud Run...${NC}"
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8081 \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 5 \
  --min-instances 0 \
  --timeout 60s \
  --set-env-vars TZ=Asia/Taipei

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 部署成功${NC}"
else
    echo -e "${RED}❌ 部署失敗${NC}"
    exit 1
fi

# 步驟 3: 取得服務 URL
echo -e "\n${GREEN}📍 步驟 3/3: 取得服務資訊...${NC}"
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region $REGION --format="value(status.url)")

echo ""
echo "================================"
echo -e "${GREEN}✅ 部署完成！${NC}"
echo "================================"
echo -e "服務 URL: ${YELLOW}$SERVICE_URL${NC}"
echo ""
echo "📊 查看日誌:"
echo "  gcloud run services logs tail $SERVICE_NAME --region $REGION"
echo ""
echo "🔄 更新服務:"
echo "  ./deploy.sh"
echo ""
echo "🗑️  刪除服務:"
echo "  gcloud run services delete $SERVICE_NAME --region $REGION"
echo "================================"
