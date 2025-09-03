#!/bin/bash

# =============================================================================
# Redis 初始化腳本
# 用於初始化 Redis 快取服務和設定基本配置
# =============================================================================

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日誌函數
log_info() {
    echo -e "${BLUE}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# 檢查環境變數
check_env() {
    log_info "檢查 Redis 環境變數..."
    
    if [[ -z "$REDIS_PASSWORD" ]]; then
        log_warning "未設定 REDIS_PASSWORD，使用預設配置"
    fi
    
    log_success "環境變數檢查完成"
}

# 建立 Redis 配置
create_config() {
    log_info "建立 Redis 配置..."
    
    cat > /etc/redis/redis.conf <<EOF
# Redis 配置檔案
bind 0.0.0.0
port 6379
timeout 0
tcp-keepalive 300
daemonize no
supervised no
pidfile /var/run/redis_6379.pid
loglevel notice
logfile ""
databases 16
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data
maxmemory 256mb
maxmemory-policy allkeys-lru
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
EOF

    # 如果設定了密碼，添加認證配置
    if [[ -n "$REDIS_PASSWORD" ]]; then
        echo "requirepass $REDIS_PASSWORD" >> /etc/redis/redis.conf
    fi
    
    log_success "Redis 配置建立完成"
}

# 初始化 Redis 資料
init_data() {
    log_info "初始化 Redis 資料..."
    
    # 等待 Redis 啟動
    sleep 2
    
    # 設定基本快取資料
    redis-cli --raw <<EOF
SET "security_intel:version" "1.0.0"
SET "security_intel:startup_time" "$(date -u +%s)"
SET "security_intel:environment" "sit"
EXPIRE "security_intel:version" 86400
EXPIRE "security_intel:startup_time" 86400
EXPIRE "security_intel:environment" 86400
EOF
    
    log_success "Redis 資料初始化完成"
}

# 主函數
main() {
    log_info "開始初始化 Redis 服務"
    
    # 檢查環境變數
    check_env
    
    # 建立配置
    create_config
    
    # 啟動 Redis
    redis-server /etc/redis/redis.conf &
    REDIS_PID=$!
    
    # 等待 Redis 啟動
    sleep 3
    
    # 初始化資料
    init_data
    
    log_success "Redis 服務初始化完成"
    
    # 等待 Redis 進程
    wait $REDIS_PID
}

# 執行主函數
main "$@"
