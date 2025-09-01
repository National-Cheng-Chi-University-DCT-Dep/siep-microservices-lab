# Ultimate Security Intelligence Platform - 部署指南與操作手冊

本文檔提供完整的部署指南和日常操作手冊，包含開發、測試、生產環境的部署流程。

## 📋 目錄

- [系統需求](#系統需求)
- [環境準備](#環境準備)
- [部署流程](#部署流程)
- [環境管理](#環境管理)
- [監控與日誌](#監控與日誌)
- [備份與恢復](#備份與恢復)
- [故障排除](#故障排除)
- [維運指南](#維運指南)
- [安全性設定](#安全性設定)

## 🖥️ 系統需求

### 硬體需求

| 環境 | CPU     | 記憶體 | 存儲  | 網路    |
| ---- | ------- | ------ | ----- | ------- |
| 開發 | 4 核心  | 8GB    | 50GB  | 100Mbps |
| 測試 | 8 核心  | 16GB   | 100GB | 1Gbps   |
| 生產 | 16 核心 | 32GB   | 500GB | 1Gbps   |

### 軟體需求

- **作業系統**: Ubuntu 20.04 LTS 或 CentOS 8+
- **Docker**: 20.10.0+
- **Docker Compose**: 2.0.0+
- **Git**: 2.25.0+
- **OpenSSL**: 1.1.1+

## 🔧 環境準備

### 1. 系統更新

```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# CentOS/RHEL
sudo yum update -y
```

### 2. 安裝 Docker

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# CentOS/RHEL
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker
```

### 3. 安裝 Docker Compose

```bash
# 下載並安裝
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version
```

### 4. 系統優化

```bash
# 增加檔案描述符限制
echo "* soft nofile 65535" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65535" | sudo tee -a /etc/security/limits.conf

# 設定 vm.max_map_count (for Elasticsearch)
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

## 🚀 部署流程

### 1. 專案克隆

```bash
git clone https://github.com/your-org/Ultimate-Security-Intelligence-Platform.git
cd Ultimate-Security-Intelligence-Platform
```

### 2. 環境設定

```bash
# 複製環境變數範例
cp env.example .env

# 編輯環境變數
nano .env
```

### 3. SSL 憑證設定

```bash
# 建立 SSL 目錄
mkdir -p nginx/ssl

# 生成自簽名憑證 (開發環境)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/default.key \
  -out nginx/ssl/default.crt \
  -subj "/C=TW/ST=Taiwan/L=Taipei/O=Organization/CN=localhost"

# 生產環境請使用 Let's Encrypt 或商業憑證
# certbot certonly --nginx -d your-domain.com
```

### 4. 資料庫初始化

```bash
# 建立資料庫初始化目錄
mkdir -p database/init

# 如果有 migration 檔案，請放置在此目錄
cp backend/database/migrations/*.sql database/init/
```

## 🌍 環境管理

### 開發環境部署

```bash
# 啟動開發環境
docker-compose -f docker-compose.dev.yml up -d

# 檢查服務狀態
docker-compose -f docker-compose.dev.yml ps

# 查看日誌
docker-compose -f docker-compose.dev.yml logs -f

# 停止服務
docker-compose -f docker-compose.dev.yml down
```

### 測試環境部署 (SIT)

```bash
# 設定測試環境變數
cp env.example .env.sit
nano .env.sit

# 啟動測試環境
docker-compose -f docker-compose.sit.yml --env-file .env.sit up -d

# 執行自動化測試
docker-compose -f docker-compose.sit.yml run --rm test-runner

# 停止服務
docker-compose -f docker-compose.sit.yml down
```

### 生產環境部署

```bash
# 設定生產環境變數
cp env.example .env.prod
nano .env.prod

# 啟動生產環境
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d

# 檢查服務狀態
docker-compose -f docker-compose.prod.yml ps

# 停止服務
docker-compose -f docker-compose.prod.yml down
```

## 📊 監控與日誌

### 服務監控

- **Grafana**: http://localhost:3001 (生產環境)
- **Prometheus**: http://localhost:9090 (生產環境)
- **Elasticsearch**: http://localhost:9200 (生產環境)
- **Kibana**: http://localhost:5601 (生產環境)

### 日誌查看

```bash
# 查看所有服務日誌
docker-compose -f docker-compose.prod.yml logs -f

# 查看特定服務日誌
docker-compose -f docker-compose.prod.yml logs -f backend
docker-compose -f docker-compose.prod.yml logs -f frontend
docker-compose -f docker-compose.prod.yml logs -f nginx

# 查看實時日誌
docker-compose -f docker-compose.prod.yml logs -f --tail=100
```

### 健康檢查

```bash
# 檢查服務健康狀態
curl -f http://localhost/health
curl -f http://localhost:8080/health
curl -f http://localhost:3000/api/health

# 檢查 Nginx 狀態
curl -f http://localhost:8080/nginx_status
```

## 💾 備份與恢復

### 資料庫備份

```bash
# 手動備份
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U security_user security_intel > backup_$(date +%Y%m%d_%H%M%S).sql

# 自動備份腳本
cat > backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/security_intel_$DATE.sql"

mkdir -p $BACKUP_DIR
docker-compose -f docker-compose.prod.yml exec -T postgres pg_dump -U security_user security_intel > $BACKUP_FILE

# 上傳到 S3
aws s3 cp $BACKUP_FILE s3://security-intel-backups/

# 刪除本地超過 7 天的備份
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
EOF

chmod +x backup.sh
```

### 資料庫恢復

```bash
# 停止服務
docker-compose -f docker-compose.prod.yml stop backend

# 恢復資料庫
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U security_user security_intel < backup_file.sql

# 重新啟動服務
docker-compose -f docker-compose.prod.yml start backend
```

### 卷備份

```bash
# 備份 Docker 卷
docker run --rm -v security-intel-postgres-data:/data -v $(pwd)/backups:/backup alpine tar czf /backup/postgres-data-$(date +%Y%m%d).tar.gz /data

# 恢復 Docker 卷
docker run --rm -v security-intel-postgres-data:/data -v $(pwd)/backups:/backup alpine tar xzf /backup/postgres-data-20231201.tar.gz
```

## 🔧 故障排除

### 常見問題

#### 1. 服務啟動失敗

```bash
# 檢查日誌
docker-compose -f docker-compose.prod.yml logs service_name

# 檢查資源使用情況
docker stats

# 檢查磁碟空間
df -h
```

#### 2. 資料庫連接失敗

```bash
# 檢查 PostgreSQL 狀態
docker-compose -f docker-compose.prod.yml exec postgres pg_isready

# 檢查網路連接
docker-compose -f docker-compose.prod.yml exec backend ping postgres

# 檢查環境變數
docker-compose -f docker-compose.prod.yml exec backend env | grep DB_
```

#### 3. Nginx 無法啟動

```bash
# 檢查配置語法
docker-compose -f docker-compose.prod.yml exec nginx nginx -t

# 檢查埠口佔用
sudo netstat -tlnp | grep :80
sudo netstat -tlnp | grep :443
```

#### 4. 記憶體不足

```bash
# 檢查記憶體使用
free -h
docker stats

# 清理無用的容器和鏡像
docker system prune -a
```

### 效能調優

#### 1. 資料庫優化

```bash
# 進入 PostgreSQL 容器
docker-compose -f docker-compose.prod.yml exec postgres psql -U security_user security_intel

# 檢查慢查詢
SELECT query, calls, total_time, mean_time FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;

# 檢查索引使用情況
SELECT tablename, indexname, num_scans, tuples_read, tuples_fetched
FROM pg_stat_user_indexes ORDER BY num_scans DESC;
```

#### 2. 快取優化

```bash
# 檢查 Redis 狀態
docker-compose -f docker-compose.prod.yml exec redis redis-cli info

# 檢查快取命中率
docker-compose -f docker-compose.prod.yml exec redis redis-cli info stats | grep keyspace
```

## 🛡️ 安全性設定

### 1. 防火牆設定

```bash
# Ubuntu UFW
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# CentOS Firewalld
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --reload
```

### 2. SSL 憑證更新

```bash
# Let's Encrypt 自動更新
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -

# 手動更新
sudo certbot renew
docker-compose -f docker-compose.prod.yml restart nginx
```

### 3. 密碼輪換

```bash
# 生成強密碼
openssl rand -base64 32

# 更新環境變數
nano .env.prod

# 重新啟動受影響的服務
docker-compose -f docker-compose.prod.yml restart backend redis
```

## 📋 維運指南

### 日常檢查清單

- [ ] 檢查服務狀態
- [ ] 查看系統資源使用情況
- [ ] 檢查日誌是否有異常
- [ ] 確認備份完成
- [ ] 檢查監控指標
- [ ] 更新安全補丁

### 定期維護

```bash
# 週間維護腳本
cat > weekly_maintenance.sh << 'EOF'
#!/bin/bash
echo "開始週間維護 $(date)"

# 清理 Docker 資源
docker system prune -f

# 更新系統
sudo apt update && sudo apt upgrade -y

# 檢查磁碟空間
df -h

# 輪換日誌
sudo logrotate /etc/logrotate.conf

# 備份資料庫
./backup.sh

echo "週間維護完成 $(date)"
EOF

chmod +x weekly_maintenance.sh
```

### 緊急程序

```bash
# 快速重啟所有服務
docker-compose -f docker-compose.prod.yml restart

# 緊急停止
docker-compose -f docker-compose.prod.yml stop

# 查看最近的錯誤
docker-compose -f docker-compose.prod.yml logs --since 1h | grep ERROR

# 資源使用情況
docker stats --no-stream
```

## 🔄 CI/CD 整合

### Drone CI 設定

1. 在 Drone 中設定必要的 Secrets
2. 推送程式碼觸發自動建置
3. 查看建置狀態和日誌

### 手動部署

```bash
# 拉取最新鏡像
docker-compose -f docker-compose.prod.yml pull

# 重新啟動服務
docker-compose -f docker-compose.prod.yml up -d

# 檢查更新狀態
docker-compose -f docker-compose.prod.yml ps
```

## 📞 支援與聯繫

- **技術支援**: tech-support@your-domain.com
- **安全問題**: security@your-domain.com
- **緊急聯繫**: +886-xxx-xxx-xxx

---

**注意**: 請定期更新此文檔，確保所有資訊都是最新的。在生產環境中進行任何更改前，請務必在測試環境中驗證。
