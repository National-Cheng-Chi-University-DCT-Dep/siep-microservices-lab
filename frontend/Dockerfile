# =============================================================================
# 資安情報平台前端 Dockerfile
# 多階段建置，最小化最終鏡像體積並提高安全性
# =============================================================================

# -----------------------------------------------------------------------------
# 第一階段：依賴安裝階段
# -----------------------------------------------------------------------------
FROM node:18-alpine AS dependencies

# 設定工作目錄
WORKDIR /app

# 複製套件定義檔案
COPY package*.json ./

# 安裝套件
RUN npm ci --only=production && npm cache clean --force

# -----------------------------------------------------------------------------
# 第二階段：建置階段
# -----------------------------------------------------------------------------
FROM node:18-alpine AS builder

# 設定建置參數
ARG VERSION=latest
ARG BUILD_TIME
ARG COMMIT_SHA
ARG NEXT_PUBLIC_API_URL=http://localhost:8080

# 設定工作目錄
WORKDIR /app

# 複製套件定義檔案
COPY package*.json ./

# 安裝所有依賴（包括開發依賴）
RUN npm ci

# 複製原始碼
COPY . .

# 設定環境變數
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
ENV NEXT_PUBLIC_VERSION=${VERSION}
ENV NEXT_PUBLIC_BUILD_TIME=${BUILD_TIME}
ENV NEXT_PUBLIC_COMMIT_SHA=${COMMIT_SHA}

# 建置 Next.js 應用程式
RUN npm run build

# -----------------------------------------------------------------------------
# 第三階段：最終執行階段
# -----------------------------------------------------------------------------
FROM node:18-alpine AS runtime

# 設定標籤
LABEL maintainer="Ultimate Security Intelligence Platform Team"
LABEL version="${VERSION}"
LABEL description="Security Intelligence Platform Frontend Service"

# 設定環境變數
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000

# 建立非特權使用者
RUN addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -s /bin/sh -D appuser

# 安裝必要的運行時依賴
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# 設定工作目錄
WORKDIR /app

# 從依賴階段複製 node_modules
COPY --from=dependencies /app/node_modules ./node_modules

# 從建置階段複製建置結果
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/next.config.* ./

# 建立必要的目錄並設定權限
RUN mkdir -p /app/.next/cache && \
    chown -R appuser:appuser /app

# 切換到非特權使用者
USER appuser

# 暴露端口
EXPOSE 3000

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/api/health || exit 1

# 設定容器啟動命令
CMD ["npm", "start"]

# =============================================================================
# 建置範例：
# docker build -t security-intel-frontend:latest \
#   --build-arg VERSION=v1.0.0 \
#   --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
#   --build-arg COMMIT_SHA=$(git rev-parse HEAD) \
#   --build-arg NEXT_PUBLIC_API_URL=https://api.yourdomain.com \
#   .
# ============================================================================= 