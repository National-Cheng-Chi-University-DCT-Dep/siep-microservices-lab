# =============================================================================
# 資安情報平台後端 Dockerfile
# 多階段建置，最小化最終鏡像體積並提高安全性
# =============================================================================

# -----------------------------------------------------------------------------
# 第一階段：建置階段
# -----------------------------------------------------------------------------
FROM golang:1.23-alpine AS builder

# 設定建置參數
ARG VERSION=latest
ARG BUILD_TIME
ARG COMMIT_SHA

# 設定工作目錄
WORKDIR /app

# 安裝必要的建置工具
RUN apk add --no-cache git ca-certificates tzdata

# 複製 go.mod 和 go.sum 檔案
COPY go.mod go.sum ./

# 下載依賴套件
RUN go mod download && go mod verify

# 複製原始碼
COPY . .

# 建置應用程式
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commitSHA=${COMMIT_SHA}" \
    -a -installsuffix cgo \
    -o server cmd/server/main.go

# 建置健康檢查工具
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o healthcheck scripts/healthcheck.go 2>/dev/null || echo "No healthcheck script found"

# -----------------------------------------------------------------------------
# 第二階段：最終執行階段
# -----------------------------------------------------------------------------
FROM alpine:3.18

# 設定標籤
LABEL maintainer="Ultimate Security Intelligence Platform Team"
LABEL version="${VERSION}"
LABEL description="Security Intelligence Platform Backend Service"

# 設定環境變數
ENV GIN_MODE=release
ENV TZ=Asia/Taipei

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

# 從建置階段複製二進位檔案
COPY --from=builder /app/server .
COPY --from=builder /app/healthcheck . 2>/dev/null || echo "No healthcheck binary"

# 複製必要的配置檔案
COPY --from=builder /app/database/migrations ./database/migrations

# 建立必要的目錄
RUN mkdir -p /app/logs /app/tmp && \
    chown -R appuser:appuser /app

# 切換到非特權使用者
USER appuser

# 暴露端口
EXPOSE 8080 9090

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

# 設定容器啟動命令
CMD ["./server"]

# =============================================================================
# 建置範例：
# docker build -t security-intel-backend:latest \
#   --build-arg VERSION=v1.0.0 \
#   --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
#   --build-arg COMMIT_SHA=$(git rev-parse HEAD) \
#   .
# =============================================================================

