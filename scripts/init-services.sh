#!/bin/bash

# =============================================================================
# Ultimate Security Intelligence Platform - Service Initialization Script
# 初始化所有後端服務的協調腳本
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

# 等待服務啟動函數
wait_for_service() {
    local service_name=$1
    local host=$2
    local port=$3
    local max_attempts=${4:-30}
    local attempt=1

    log_info "等待服務 $service_name 啟動 ($host:$port)..."

    while [ $attempt -le $max_attempts ]; do
        if nc -z $host $port 2>/dev/null; then
            log_success "$service_name 服務已啟動"
            return 0
        fi
        
        log_info "嘗試 $attempt/$max_attempts - $service_name 尚未就緒，等待 2 秒..."
        sleep 2
        attempt=$((attempt + 1))
    done

    log_error "$service_name 服務啟動超時"
    return 1
}

# 檢查環境變數
check_environment() {
    log_info "檢查環境變數..."

    # 必要的環境變數
    required_vars=(
        "POSTGRES_DB"
        "POSTGRES_USER"
        "POSTGRES_PASSWORD"
        "REDIS_PASSWORD"
        "JWT_SECRET"
    )

    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            log_error "缺少必要的環境變數: $var"
            exit 1
        fi
    done

    log_success "環境變數檢查完成"
}

# 初始化資料庫
init_database() {
    log_info "初始化 PostgreSQL 資料庫..."

    # 啟動 PostgreSQL
    if [ ! -f /var/lib/postgresql/data/postgresql.conf ]; then
        log_info "首次啟動 PostgreSQL，初始化資料庫..."
        su postgres -c "initdb -D /var/lib/postgresql/data --encoding=UTF-8 --locale=C"
    fi

    # 啟動 PostgreSQL 服務
    su postgres -c "pg_ctl -D /var/lib/postgresql/data -l /var/lib/postgresql/data/postgresql.log start"

    # 等待 PostgreSQL 啟動
    wait_for_service "PostgreSQL" "localhost" "5432" 30

    # 執行初始化腳本
    if [ -d "/docker-entrypoint-initdb.d" ]; then
        log_info "執行資料庫初始化腳本..."
        for script in /docker-entrypoint-initdb.d/*.sh; do
            if [ -f "$script" ]; then
                log_info "執行腳本: $script"
                bash "$script"
            fi
        done
    fi

    log_success "PostgreSQL 資料庫初始化完成"
}

# 初始化 Redis
init_redis() {
    log_info "初始化 Redis 快取服務..."

    # 啟動 Redis 服務
    redis-server /etc/redis/redis.conf --daemonize yes

    # 等待 Redis 啟動
    wait_for_service "Redis" "localhost" "6379" 15

    # 測試 Redis 連接
    if redis-cli ping >/dev/null 2>&1; then
        log_success "Redis 服務初始化完成"
    else
        log_error "Redis 服務啟動失敗"
        exit 1
    fi
}

# 初始化 MQTT
init_mqtt() {
    log_info "初始化 MQTT 訊息佇列服務..."

    # 啟動 Mosquitto 服務
    /usr/sbin/mosquitto -c /etc/mosquitto/mosquitto.conf -d

    # 等待 MQTT 啟動
    wait_for_service "MQTT" "localhost" "1883" 15

    log_success "MQTT 服務初始化完成"
}

# 初始化監控服務
init_monitoring() {
    log_info "初始化監控服務..."

    # 啟動 Prometheus
    su nobody -s /bin/bash -c "prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/var/lib/prometheus/data --storage.tsdb.retention.time=72h --web.enable-lifecycle" &
    PROMETHEUS_PID=$!

    # 等待 Prometheus 啟動
    wait_for_service "Prometheus" "localhost" "9090" 20

    # 啟動 Grafana
    su grafana -s /bin/bash -c "grafana-server --config=/etc/grafana/grafana.ini --homepath=/usr/share/grafana" &
    GRAFANA_PID=$!

    # 等待 Grafana 啟動
    wait_for_service "Grafana" "localhost" "3000" 20

    log_success "監控服務初始化完成"
}

# 初始化後端 API
init_backend() {
    log_info "初始化後端 API 服務..."

    # 切換到應用程式目錄
    cd /app

    # 執行資料庫遷移
    if [ -f "./server" ]; then
        log_info "執行資料庫遷移..."
        # 這裡可以添加資料庫遷移命令
        # ./server migrate
    fi

    # 啟動後端服務
    su appuser -s /bin/bash -c "./server" &
    BACKEND_PID=$!

    # 等待後端服務啟動
    wait_for_service "Backend API" "localhost" "8080" 30

    log_success "後端 API 服務初始化完成"
}

# 健康檢查
health_check() {
    log_info "執行健康檢查..."

    local services=(
        "PostgreSQL:5432"
        "Redis:6379"
        "MQTT:1883"
        "Prometheus:9090"
        "Grafana:3000"
        "Backend API:8080"
    )

    local all_healthy=true

    for service in "${services[@]}"; do
        local name=$(echo $service | cut -d: -f1)
        local port=$(echo $service | cut -d: -f2)

        if nc -z localhost $port 2>/dev/null; then
            log_success "$name 健康檢查通過"
        else
            log_error "$name 健康檢查失敗"
            all_healthy=false
        fi
    done

    if [ "$all_healthy" = true ]; then
        log_success "所有服務健康檢查通過"
        return 0
    else
        log_error "部分服務健康檢查失敗"
        return 1
    fi
}

# 清理函數
cleanup() {
    log_info "收到停止信號，正在清理服務..."
    
    # 停止所有後台進程
    jobs -p | xargs -r kill
    
    # 等待進程結束
    wait
    
    log_success "服務清理完成"
    exit 0
}

# 主函數
main() {
    log_info "開始初始化 Ultimate Security Intelligence Platform 後端服務"

    # 設定信號處理
    trap cleanup SIGTERM SIGINT

    # 檢查環境變數
    check_environment

    # 初始化各項服務
    init_database
    init_redis
    init_mqtt
    init_monitoring
    init_backend

    # 執行健康檢查
    health_check

    log_success "所有後端服務初始化完成！"
    log_info "服務端口："
    log_info "  - Backend API: http://localhost:8080"
    log_info "  - PostgreSQL: localhost:5432"
    log_info "  - Redis: localhost:6379"
    log_info "  - MQTT: localhost:1883"
    log_info "  - Prometheus: http://localhost:9090"
    log_info "  - Grafana: http://localhost:3000"

    # 保持腳本運行
    while true; do
        sleep 30
        health_check
    done
}

# 執行主函數
main "$@"
