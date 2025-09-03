# =============================================================================
# Ultimate Security Intelligence Platform - Backend Services Dockerfile
# 多階段建置，初始化所有後端服務
# =============================================================================

# -----------------------------------------------------------------------------
# 第一階段：後端 API 建置階段
# -----------------------------------------------------------------------------
FROM golang:1.23-alpine AS backend-builder

# 設定建置參數
ARG VERSION=latest
ARG BUILD_TIME
ARG COMMIT_SHA

# 設定工作目錄
WORKDIR /app

# 安裝必要的建置工具
RUN apk add --no-cache git ca-certificates tzdata

# 複製 go.mod 和 go.sum 檔案
COPY ../backend/go.mod ../backend/go.sum ./

# 下載依賴套件
RUN go mod download && go mod verify

# 複製後端原始碼
COPY ../backend/ .

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
# 第二階段：資料庫初始化腳本建置
# -----------------------------------------------------------------------------
FROM postgres:15-alpine AS db-init

# 安裝必要的工具
RUN apk add --no-cache bash curl

# 複製資料庫初始化腳本
COPY ../database/init/ /docker-entrypoint-initdb.d/

# 設定權限
RUN chmod +x /docker-entrypoint-initdb.d/*.sh

# -----------------------------------------------------------------------------
# 第三階段：Redis 配置準備
# -----------------------------------------------------------------------------
FROM redis:7-alpine AS redis-config

# 複製 Redis 配置檔案
COPY ../redis/redis-sit.conf /usr/local/etc/redis/redis.conf

# -----------------------------------------------------------------------------
# 第四階段：MQTT 配置準備
# -----------------------------------------------------------------------------
FROM eclipse-mosquitto:2.0 AS mqtt-config

# 複製 MQTT 配置檔案
COPY ../mosquitto/config-sit/ /mosquitto/config/

# -----------------------------------------------------------------------------
# 第五階段：Prometheus 配置準備
# -----------------------------------------------------------------------------
FROM prom/prometheus:latest AS prometheus-config

# 複製 Prometheus 配置檔案
COPY ../prometheus/config-sit/ /etc/prometheus/

# -----------------------------------------------------------------------------
# 第六階段：Grafana 配置準備
# -----------------------------------------------------------------------------
FROM grafana/grafana:latest AS grafana-config

# 複製 Grafana 配置檔案
COPY ../grafana/provisioning-sit/ /etc/grafana/provisioning/
COPY ../grafana/dashboards-sit/ /var/lib/grafana/dashboards/

# -----------------------------------------------------------------------------
# 第七階段：最終後端服務鏡像
# -----------------------------------------------------------------------------
FROM alpine:3.18 AS backend-service

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
    postgresql-client \
    redis \
    && rm -rf /var/cache/apk/*

# 設定工作目錄
WORKDIR /app

# 從建置階段複製二進位檔案
COPY --from=backend-builder /app/server .
COPY --from=backend-builder /app/healthcheck . 2>/dev/null || echo "No healthcheck binary"

# 複製必要的配置檔案
COPY --from=backend-builder /app/database/migrations ./database/migrations

# 建立必要的目錄
RUN mkdir -p /app/logs /app/tmp /app/uploads && \
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

# -----------------------------------------------------------------------------
# 第八階段：資料庫服務鏡像
# -----------------------------------------------------------------------------
FROM postgres:15-alpine AS database-service

# 設定標籤
LABEL maintainer="Ultimate Security Intelligence Platform Team"
LABEL description="PostgreSQL Database Service"

# 設定環境變數
ENV TZ=Asia/Taipei

# 複製初始化腳本
COPY --from=db-init /docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/

# 建立資料目錄
RUN mkdir -p /var/lib/postgresql/data /backups && \
    chown -R postgres:postgres /var/lib/postgresql/data /backups

# 暴露端口
EXPOSE 5432

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD pg_isready -U postgres || exit 1

# -----------------------------------------------------------------------------
# 第九階段：快取服務鏡像
# -----------------------------------------------------------------------------
FROM redis:7-alpine AS cache-service

# 設定標籤
LABEL maintainer="Ultimate Security Intelligence Platform Team"
LABEL description="Redis Cache Service"

# 設定環境變數
ENV TZ=Asia/Taipei

# 複製配置檔案
COPY --from=redis-config /usr/local/etc/redis/redis.conf /usr/local/etc/redis/redis.conf

# 建立資料目錄
RUN mkdir -p /data && chown -R redis:redis /data

# 暴露端口
EXPOSE 6379

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD redis-cli --raw incr ping || exit 1

# 設定啟動命令
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]

# -----------------------------------------------------------------------------
# 第十階段：訊息佇列服務鏡像
# -----------------------------------------------------------------------------
FROM eclipse-mosquitto:2.0 AS mqtt-service

# 設定標籤
LABEL maintainer="Ultimate Security Intelligence Platform Team"
LABEL description="MQTT Message Broker Service"

# 設定環境變數
ENV TZ=Asia/Taipei

# 複製配置檔案
COPY --from=mqtt-config /mosquitto/config/ /mosquitto/config/

# 建立必要目錄
RUN mkdir -p /mosquitto/data /mosquitto/log && \
    chown -R mosquitto:mosquitto /mosquitto

# 暴露端口
EXPOSE 1883 9001

# 設定啟動命令
CMD ["/usr/sbin/mosquitto", "-c", "/mosquitto/config/mosquitto.conf"]

# -----------------------------------------------------------------------------
# 第十一階段：監控服務鏡像
# -----------------------------------------------------------------------------
FROM prom/prometheus:latest AS monitoring-service

# 設定標籤
LABEL maintainer="Ultimate Security Intelligence Platform Team"
LABEL description="Prometheus Monitoring Service"

# 設定環境變數
ENV TZ=Asia/Taipei

# 複製配置檔案
COPY --from=prometheus-config /etc/prometheus/ /etc/prometheus/

# 建立資料目錄
RUN mkdir -p /prometheus && chown -R nobody:nobody /prometheus

# 暴露端口
EXPOSE 9090

# 設定啟動命令
CMD ["--config.file=/etc/prometheus/prometheus.yml", \
     "--storage.tsdb.path=/prometheus", \
     "--storage.tsdb.retention.time=72h", \
     "--web.enable-lifecycle"]

# -----------------------------------------------------------------------------
# 第十二階段：視覺化服務鏡像
# -----------------------------------------------------------------------------
FROM grafana/grafana:latest AS visualization-service

# 設定標籤
LABEL maintainer="Ultimate Security Intelligence Platform Team"
LABEL description="Grafana Visualization Service"

# 設定環境變數
ENV TZ=Asia/Taipei

# 複製配置檔案
COPY --from=grafana-config /etc/grafana/provisioning/ /etc/grafana/provisioning/
COPY --from=grafana-config /var/lib/grafana/dashboards/ /var/lib/grafana/dashboards/

# 建立資料目錄
RUN mkdir -p /var/lib/grafana && chown -R grafana:grafana /var/lib/grafana

# 暴露端口
EXPOSE 3000

# 設定啟動命令
CMD ["grafana-server", "--config=/etc/grafana/grafana.ini"]

# -----------------------------------------------------------------------------
# 第十三階段：初始化腳本建置
# -----------------------------------------------------------------------------
FROM alpine:3.18 AS init-scripts

# 安裝必要的工具
RUN apk add --no-cache bash curl postgresql-client redis

# 複製初始化腳本
COPY ../scripts/init-services.sh /usr/local/bin/
COPY ../scripts/wait-for-it.sh /usr/local/bin/

# 設定權限
RUN chmod +x /usr/local/bin/*.sh

# 建立工作目錄
WORKDIR /scripts

# 複製服務初始化腳本
COPY ../scripts/init-database.sh ./
COPY ../scripts/init-redis.sh ./
COPY ../scripts/init-mqtt.sh ./
COPY ../scripts/init-monitoring.sh ./

# 設定權限
RUN chmod +x *.sh

# -----------------------------------------------------------------------------
# 第十四階段：最終多服務鏡像
# -----------------------------------------------------------------------------
FROM alpine:3.18 AS multi-service

# 設定標籤
LABEL maintainer="Ultimate Security Intelligence Platform Team"
LABEL description="Multi-Service Backend Platform"

# 安裝必要的工具和服務
RUN apk add --no-cache \
    bash \
    curl \
    postgresql-client \
    redis \
    mosquitto \
    prometheus \
    grafana \
    supervisor \
    tzdata \
    && rm -rf /var/cache/apk/*

# 設定時區
ENV TZ=Asia/Taipei

# 建立必要的目錄
RUN mkdir -p \
    /app/logs \
    /app/uploads \
    /var/lib/postgresql/data \
    /var/lib/redis/data \
    /var/lib/mosquitto/data \
    /var/lib/mosquitto/log \
    /var/lib/prometheus/data \
    /var/lib/grafana/data \
    /etc/supervisor/conf.d

# 從各階段複製檔案
COPY --from=backend-builder /app/server /app/
COPY --from=db-init /docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/
COPY --from=redis-config /usr/local/etc/redis/redis.conf /etc/redis/redis.conf
COPY --from=mqtt-config /mosquitto/config/ /etc/mosquitto/
COPY --from=prometheus-config /etc/prometheus/ /etc/prometheus/
COPY --from=grafana-config /etc/grafana/provisioning/ /etc/grafana/provisioning/
COPY --from=grafana-config /var/lib/grafana/dashboards/ /var/lib/grafana/dashboards/
COPY --from=init-scripts /usr/local/bin/ /usr/local/bin/
COPY --from=init-scripts /scripts/ /scripts/

# 複製 Supervisor 配置
COPY supervisor/supervisord.conf /etc/supervisor/supervisord.conf
COPY supervisor/conf.d/ /etc/supervisor/conf.d/

# 設定權限
RUN chmod +x /usr/local/bin/*.sh /scripts/*.sh && \
    chown -R postgres:postgres /var/lib/postgresql/data && \
    chown -R redis:redis /var/lib/redis/data && \
    chown -R mosquitto:mosquitto /var/lib/mosquitto && \
    chown -R nobody:nobody /var/lib/prometheus/data && \
    chown -R grafana:grafana /var/lib/grafana

# 暴露所有服務端口
EXPOSE 8080 9090 5432 6379 1883 9001 9090 3000

# 複製啟動腳本
COPY scripts/start-services.sh /start-services.sh
RUN chmod +x /start-services.sh

# 設定啟動命令
CMD ["/start-services.sh"]

# =============================================================================
# 建置範例：
# docker build -t security-intel-backend:latest \
#   --build-arg VERSION=v1.0.0 \
#   --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
#   --build-arg COMMIT_SHA=$(git rev-parse HEAD) \
#   .
# =============================================================================
