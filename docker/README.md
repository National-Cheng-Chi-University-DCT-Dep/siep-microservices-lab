# Ultimate Security Intelligence Platform - Docker Setup

## 概述

這個 Dockerfile 和相關腳本用於初始化 Ultimate Security Intelligence Platform 的所有後端服務。

## 服務架構

### 核心服務

- **Backend API** (Go): 主要的後端 API 服務，提供 RESTful API 和 gRPC 介面
- **PostgreSQL**: 主要資料庫，儲存威脅情報、安全事件、資產等資料
- **Redis**: 快取服務，用於會話管理和資料快取
- **MQTT (Mosquitto)**: 訊息佇列服務，用於即時通訊和事件通知

### 監控服務

- **Prometheus**: 監控和指標收集服務
- **Grafana**: 監控視覺化儀表板

## 目錄結構

```
docker/
├── Dockerfile                    # 主要的多階段 Dockerfile
├── README.md                     # 本檔案
├── scripts/                      # 服務初始化腳本
│   ├── init-services.sh         # 主要服務初始化協調腳本
│   ├── wait-for-it.sh           # 等待服務啟動的通用腳本
│   ├── init-database.sh         # 資料庫初始化腳本
│   ├── init-redis.sh            # Redis 初始化腳本
│   ├── init-mqtt.sh             # MQTT 初始化腳本
│   ├── init-monitoring.sh       # 監控服務初始化腳本
│   └── start-services.sh        # 服務啟動主腳本
├── supervisor/                  # Supervisor 配置
│   ├── supervisord.conf         # 主要 Supervisor 配置
│   └── conf.d/                  # 個別服務配置
├── database/                    # 資料庫相關檔案
│   └── init/                    # 資料庫初始化腳本
│       └── 01-init-database.sh # PostgreSQL 初始化腳本
├── redis/                       # Redis 配置
│   └── redis-sit.conf          # Redis 配置檔案
├── mosquitto/                   # MQTT 配置
│   └── config-sit/             # Mosquitto 配置目錄
│       └── mosquitto.conf      # Mosquitto 配置檔案
├── prometheus/                  # Prometheus 配置
│   └── config-sit/             # Prometheus 配置目錄
│       ├── prometheus.yml      # 主要配置檔案
│       └── alert_rules.yml     # 告警規則
└── grafana/                     # Grafana 配置
    ├── provisioning-sit/        # 自動配置目錄
    │   ├── datasources/         # 資料來源配置
    │   │   └── prometheus.yml  # Prometheus 資料來源
    │   └── dashboards/          # 儀表板配置
    │       └── dashboards.yml  # 儀表板自動載入配置
    └── dashboards-sit/          # 儀表板檔案
        └── security-intelligence-overview.json # 基本儀表板
```

## 建置和使用

### 1. 建置 Docker 鏡像

```bash
# 建置多服務鏡像
docker build -t security-intel-backend:latest \
  --build-arg VERSION=v1.0.0 \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg COMMIT_SHA=$(git rev-parse HEAD) \
  -f docker/Dockerfile .
```

### 2. 設定環境變數

建立 `.env` 檔案：

```bash
# 資料庫設定
POSTGRES_DB=security_intel_sit
POSTGRES_USER=sit_user
POSTGRES_PASSWORD=your_secure_password

# Redis 設定
REDIS_PASSWORD=your_redis_password

# MQTT 設定
MQTT_USERNAME=mqtt_user
MQTT_PASSWORD=your_mqtt_password

# JWT 設定
JWT_SECRET=your_jwt_secret_key

# 應用程式設定
APP_ENV=staging
APP_DEBUG=false
```

### 3. 運行容器

```bash
# 運行多服務容器
docker run -d \
  --name security-intel-backend \
  --env-file .env \
  -p 8080:8080 \
  -p 5432:5432 \
  -p 6379:6379 \
  -p 1883:1883 \
  -p 9090:9090 \
  -p 3000:3000 \
  security-intel-backend:latest
```

### 4. 檢查服務狀態

```bash
# 檢查容器狀態
docker ps

# 查看服務日誌
docker logs security-intel-backend

# 進入容器
docker exec -it security-intel-backend /bin/bash
```

## 服務端口

| 服務           | 端口 | 說明                     |
| -------------- | ---- | ------------------------ |
| Backend API    | 8080 | RESTful API 和 gRPC 服務 |
| PostgreSQL     | 5432 | 資料庫服務               |
| Redis          | 6379 | 快取服務                 |
| MQTT           | 1883 | 訊息佇列服務             |
| MQTT WebSocket | 9001 | WebSocket 連接           |
| Prometheus     | 9090 | 監控服務                 |
| Grafana        | 3000 | 監控儀表板               |

## 健康檢查

所有服務都包含健康檢查機制：

- **Backend API**: `http://localhost:8080/api/v1/health`
- **PostgreSQL**: `pg_isready` 命令
- **Redis**: `redis-cli ping` 命令
- **MQTT**: 端口連接檢查
- **Prometheus**: `http://localhost:9090/-/healthy`
- **Grafana**: `http://localhost:3000/api/health`

## 監控和日誌

### 監控

- **Prometheus**: 收集所有服務的指標
- **Grafana**: 提供視覺化儀表板
- **告警**: 設定自動告警規則

### 日誌

- 所有服務的日誌都輸出到標準輸出
- 可以使用 Docker 日誌功能查看：`docker logs security-intel-backend`
- Supervisor 管理所有服務的日誌

## 故障排除

### 常見問題

1. **服務啟動失敗**

   - 檢查環境變數是否正確設定
   - 查看容器日誌：`docker logs security-intel-backend`
   - 確認端口沒有被佔用

2. **資料庫連接失敗**

   - 確認 PostgreSQL 服務已啟動
   - 檢查資料庫憑證是否正確
   - 確認資料庫初始化腳本已執行

3. **監控服務無法訪問**
   - 確認 Prometheus 和 Grafana 服務已啟動
   - 檢查防火牆設定
   - 確認端口映射正確

### 調試模式

```bash
# 以調試模式運行
docker run -it \
  --name security-intel-backend-debug \
  --env-file .env \
  security-intel-backend:latest \
  /bin/bash

# 手動啟動服務
/usr/local/bin/init-services.sh
```

## 安全考量

1. **密碼管理**

   - 使用強密碼
   - 定期更換密碼
   - 不要在程式碼中硬編碼密碼

2. **網路安全**

   - 限制容器網路訪問
   - 使用防火牆規則
   - 考慮使用 VPN 或私有網路

3. **資料安全**
   - 加密敏感資料
   - 定期備份資料
   - 實施存取控制

## 開發和測試

### 開發環境

```bash
# 建置開發版本
docker build -t security-intel-backend:dev \
  --build-arg VERSION=dev \
  --target backend-service \
  -f docker/Dockerfile .
```

### 測試

```bash
# 運行測試
docker run --rm \
  --env-file .env \
  security-intel-backend:latest \
  /scripts/run-tests.sh
```

## 貢獻

歡迎貢獻程式碼！請遵循以下步驟：

1. Fork 專案
2. 建立功能分支
3. 提交變更
4. 建立 Pull Request

## 授權

本專案採用 MIT 授權條款。詳見 [LICENSE](../LICENSE) 檔案。
