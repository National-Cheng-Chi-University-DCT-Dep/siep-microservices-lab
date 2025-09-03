#!/bin/bash

# =============================================================================
# 服務啟動腳本
# 用於啟動所有後端服務的主協調腳本
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

# 設定信號處理
trap cleanup SIGTERM SIGINT

# 主函數
main() {
    log_info "開始啟動 Ultimate Security Intelligence Platform 後端服務"
    
    # 執行服務初始化腳本
    /usr/local/bin/init-services.sh
    
    log_success "所有後端服務啟動完成！"
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
    done
}

# 執行主函數
main "$@"
