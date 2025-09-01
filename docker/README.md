# Docker 開發環境

這個目錄包含資安情報平台的 Docker 開發環境配置。

## 服務概覽

| 服務                | 描述       | 端口                    | 管理介面               |
| ------------------- | ---------- | ----------------------- | ---------------------- |
| **PostgreSQL**      | 主要資料庫 | 5432                    | PgAdmin (5050)         |
| **Redis**           | 快取和排程 | 6379                    | Redis Commander (8081) |
| **PgAdmin**         | 資料庫管理 | 5050                    | http://localhost:5050  |
| **Redis Commander** | Redis 管理 | 8081                    | http://localhost:8081  |
| **MailHog**         | 郵件測試   | 1025 (SMTP), 8025 (Web) | http://localhost:8025  |

## 快速開始

### 1. 啟動所有服務

```bash
# 在專案根目錄
docker-compose -f docker/docker-compose.yml up -d

# 或者使用 make 指令 (如果已設定)
make docker-up
```

### 2. 檢查服務狀態

```bash
docker-compose -f docker/docker-compose.yml ps
```

### 3. 檢視日誌

```bash
# 檢視所有服務日誌
docker-compose -f docker/docker-compose.yml logs -f

# 檢視特定服務日誌
docker-compose -f docker/docker-compose.yml logs -f postgres
```

### 4. 停止服務

```bash
docker-compose -f docker/docker-compose.yml down

# 停止並刪除 volumes (注意：會刪除所有資料)
docker-compose -f docker/docker-compose.yml down -v
```

## 服務詳細資訊

### PostgreSQL 資料庫

- **映像**: `postgres:15-alpine`
- **預設資料庫**: `security_intel`
- **使用者**: `postgres`
- **密碼**: `postgres`
- **連線字串**: `postgres://postgres:postgres@localhost:5432/security_intel?sslmode=disable`

#### 預設測試帳戶

- **管理員**: admin@security-intel.com / admin123
- **測試使用者**: test@security-intel.com / test123

### Redis

- **映像**: `redis:7-alpine`
- **配置**: 啟用 AOF 持久化
- **連線**: `redis://localhost:6379`

### PgAdmin

- **映像**: `dpage/pgadmin4:latest`
- **URL**: http://localhost:5050
- **登入**: admin@security-intel.com / admin
- **伺服器設定**:
  - 主機: postgres
  - 端口: 5432
  - 使用者: postgres
  - 密碼: postgres

### Redis Commander

- **映像**: `rediscommander/redis-commander:latest`
- **URL**: http://localhost:8081
- **自動連線**: 已配置連線到 Redis

### MailHog

- **映像**: `mailhog/mailhog:latest`
- **SMTP**: localhost:1025
- **Web UI**: http://localhost:8025
- **用途**: 攔截開發環境的郵件

## 資料持久化

以下 volumes 用於資料持久化：

- `postgres_data`: PostgreSQL 資料
- `redis_data`: Redis 資料

## 網路配置

所有服務都在 `security-intel-network` 網路中，可以互相通信。

## 健康檢查

PostgreSQL 和 Redis 都配置了健康檢查：

```bash
# 檢查 PostgreSQL 健康狀態
docker-compose -f docker/docker-compose.yml exec postgres pg_isready -U postgres

# 檢查 Redis 健康狀態
docker-compose -f docker/docker-compose.yml exec redis redis-cli ping
```

## 疑難排解

### 常見問題

1. **端口被佔用**

   ```bash
   # 檢查端口使用情況
   lsof -i :5432  # PostgreSQL
   lsof -i :6379  # Redis
   lsof -i :5050  # PgAdmin
   ```

2. **資料庫連線失敗**

   ```bash
   # 檢查容器是否正在運行
   docker-compose -f docker/docker-compose.yml ps

   # 檢查 PostgreSQL 日誌
   docker-compose -f docker/docker-compose.yml logs postgres
   ```

3. **清理和重建**

   ```bash
   # 停止並刪除所有資料
   docker-compose -f docker/docker-compose.yml down -v

   # 重新啟動
   docker-compose -f docker/docker-compose.yml up -d
   ```

### 重設資料庫

如果需要重設資料庫：

```bash
# 停止服務
docker-compose -f docker/docker-compose.yml stop postgres

# 刪除 PostgreSQL volume
docker volume rm docker_postgres_data

# 重新啟動 PostgreSQL
docker-compose -f docker/docker-compose.yml up -d postgres
```

## 開發提示

1. **連線字串**: 在後端 `.env` 檔案中使用：

   ```
   DATABASE_DSN=postgres://postgres:postgres@localhost:5432/security_intel?sslmode=disable
   ```

2. **Redis 連線**: 在程式中使用：

   ```
   REDIS_URL=redis://localhost:6379
   ```

3. **郵件測試**: 設定 SMTP 為 `localhost:1025`，然後在 http://localhost:8025 查看郵件

## 監控

### 資源使用情況

```bash
# 檢視容器資源使用情況
docker stats

# 檢視特定容器
docker stats security-intel-db security-intel-redis
```

### 日誌監控

```bash
# 即時監控所有服務日誌
docker-compose -f docker/docker-compose.yml logs -f --tail=100
```

## 備份和還原

### 備份資料庫

```bash
docker-compose -f docker/docker-compose.yml exec postgres pg_dump -U postgres security_intel > backup.sql
```

### 還原資料庫

```bash
docker-compose -f docker/docker-compose.yml exec -T postgres psql -U postgres security_intel < backup.sql
```
