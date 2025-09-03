#!/bin/bash

# =============================================================================
# Docker 建置問題修復腳本
# 修復所有缺少的目錄和檔案問題
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
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 檢查並建立目錄
check_and_create_dir() {
    local dir="$1"
    if [ ! -d "$dir" ]; then
        log_info "建立目錄: $dir"
        mkdir -p "$dir"
        log_success "目錄已建立: $dir"
    else
        log_info "目錄已存在: $dir"
    fi
}

# 檢查並建立檔案
check_and_create_file() {
    local file="$1"
    local content="$2"
    if [ ! -f "$file" ]; then
        log_info "建立檔案: $file"
        echo "$content" > "$file"
        log_success "檔案已建立: $file"
    else
        log_info "檔案已存在: $file"
    fi
}

# 主函數
main() {
    log_info "開始修復 Docker 建置問題..."
    
    # 1. 建立缺少的目錄
    log_info "檢查並建立缺少的目錄..."
    
    check_and_create_dir "supervisor/conf.d"
    check_and_create_dir "database/init"
    check_and_create_dir "scripts"
    check_and_create_dir "grafana/provisioning-sit"
    check_and_create_dir "grafana/dashboards-sit"
    check_and_create_dir "redis"
    check_and_create_dir "mosquitto/config-sit"
    check_and_create_dir "prometheus/config-sit"
    check_and_create_dir "backend"
    
    # 2. 檢查並建立必要的檔案
    log_info "檢查並建立必要的檔案..."
    
    # 檢查 scripts/start-services.sh
    if [ ! -f "scripts/start-services.sh" ]; then
        log_warning "缺少 scripts/start-services.sh，請確保此檔案存在"
    fi
    
    # 檢查 supervisor/conf.d/ 下的配置檔案
    if [ ! -f "supervisor/conf.d/postgresql.conf" ]; then
        log_warning "缺少 supervisor/conf.d/postgresql.conf"
    fi
    
    if [ ! -f "supervisor/conf.d/redis.conf" ]; then
        log_warning "缺少 supervisor/conf.d/redis.conf"
    fi
    
    if [ ! -f "supervisor/conf.d/mosquitto.conf" ]; then
        log_warning "缺少 supervisor/conf.d/mosquitto.conf"
    fi
    
    if [ ! -f "supervisor/conf.d/prometheus.conf" ]; then
        log_warning "缺少 supervisor/conf.d/prometheus.conf"
    fi
    
    if [ ! -f "supervisor/conf.d/grafana.conf" ]; then
        log_warning "缺少 supervisor/conf.d/grafana.conf"
    fi
    
    if [ ! -f "supervisor/conf.d/backend.conf" ]; then
        log_warning "缺少 supervisor/conf.d/backend.conf"
    fi
    
    # 3. 檢查配置檔案
    log_info "檢查配置檔案..."
    
    if [ ! -f "redis/redis-sit.conf" ]; then
        log_warning "缺少 redis/redis-sit.conf"
    fi
    
    if [ ! -f "mosquitto/config-sit/mosquitto.conf" ]; then
        log_warning "缺少 mosquitto/config-sit/mosquitto.conf"
    fi
    
    if [ ! -f "prometheus/config-sit/prometheus.yml" ]; then
        log_warning "缺少 prometheus/config-sit/prometheus.yml"
    fi
    
    if [ ! -f "grafana/provisioning-sit/datasources/prometheus.yml" ]; then
        log_warning "缺少 grafana/provisioning-sit/datasources/prometheus.yml"
    fi
    
    if [ ! -f "grafana/provisioning-sit/dashboards/dashboards.yml" ]; then
        log_warning "缺少 grafana/provisioning-sit/dashboards/dashboards.yml"
    fi
    
    # 4. 檢查後端檔案
    log_info "檢查後端檔案..."
    
    if [ ! -f "backend/go.mod" ]; then
        log_warning "缺少 backend/go.mod"
    fi
    
    if [ ! -f "backend/go.sum" ]; then
        log_warning "缺少 backend/go.sum"
    fi
    
    # 5. 檢查資料庫初始化檔案
    log_info "檢查資料庫初始化檔案..."
    
    if [ ! -f "database/init/01-init-database.sh" ]; then
        log_warning "缺少 database/init/01-init-database.sh"
    fi
    
    # 6. 設定檔案權限
    log_info "設定檔案權限..."
    
    if [ -f "scripts/start-services.sh" ]; then
        chmod +x scripts/start-services.sh
        log_success "已設定 scripts/start-services.sh 執行權限"
    fi
    
    if [ -d "scripts" ]; then
        chmod +x scripts/*.sh 2>/dev/null || true
        log_success "已設定 scripts 目錄下所有 .sh 檔案執行權限"
    fi
    
    # 7. 驗證修復結果
    log_info "驗證修復結果..."
    
    local missing_files=0
    
    # 檢查關鍵檔案
    critical_files=(
        "scripts/start-services.sh"
        "supervisor/conf.d/postgresql.conf"
        "supervisor/conf.d/redis.conf"
        "supervisor/conf.d/mosquitto.conf"
        "supervisor/conf.d/prometheus.conf"
        "supervisor/conf.d/grafana.conf"
        "supervisor/conf.d/backend.conf"
        "redis/redis-sit.conf"
        "mosquitto/config-sit/mosquitto.conf"
        "prometheus/config-sit/prometheus.yml"
        "backend/go.mod"
        "backend/go.sum"
        "database/init/01-init-database.sh"
    )
    
    for file in "${critical_files[@]}"; do
        if [ ! -f "$file" ]; then
            log_error "缺少關鍵檔案: $file"
            ((missing_files++))
        fi
    done
    
    if [ $missing_files -eq 0 ]; then
        log_success "所有關鍵檔案都存在！"
    else
        log_warning "仍有 $missing_files 個關鍵檔案缺少"
    fi
    
    # 8. 顯示建置建議
    log_info "Docker 建置建議："
    echo "1. 確保所有必要的檔案都存在"
    echo "2. 使用以下命令建置："
    echo "   docker build -t security-intel-backend:latest \\"
    echo "     --build-arg VERSION=v1.0.0 \\"
    echo "     --build-arg BUILD_TIME=\$(date -u +\"%Y-%m-%dT%H:%M:%SZ\") \\"
    echo "     --build-arg COMMIT_SHA=\$(git rev-parse HEAD) \\"
    echo "     -f docker/Dockerfile ."
    echo ""
    echo "3. 如果仍有問題，請檢查："
    echo "   - 檔案路徑是否正確"
    echo "   - 檔案權限是否設定"
    echo "   - Docker 建置上下文是否正確"
    
    log_success "Docker 建置問題修復完成！"
}

# 執行主函數
main "$@"
