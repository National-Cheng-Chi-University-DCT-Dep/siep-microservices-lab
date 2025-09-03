#!/bin/bash

# =============================================================================
# Docker 建置測試腳本
# 測試修復後的 Docker 建置是否成功
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

# 主函數
main() {
    log_info "開始測試 Docker 建置..."
    
    # 設定建置參數
    VERSION=${VERSION:-"v1.0.0"}
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    COMMIT_SHA=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
    IMAGE_NAME="security-intel-backend"
    TAG="latest"
    
    log_info "建置參數："
    log_info "  - VERSION: $VERSION"
    log_info "  - BUILD_TIME: $BUILD_TIME"
    log_info "  - COMMIT_SHA: $COMMIT_SHA"
    log_info "  - IMAGE_NAME: $IMAGE_NAME"
    log_info "  - TAG: $TAG"
    
    # 檢查 Docker 是否可用
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安裝或不在 PATH 中"
        exit 1
    fi
    
    # 檢查 Dockerfile 是否存在
    if [ ! -f "docker/Dockerfile" ]; then
        log_error "找不到 docker/Dockerfile"
        exit 1
    fi
    
    # 執行 Docker 建置
    log_info "開始執行 Docker 建置..."
    
    docker build \
        -t "${IMAGE_NAME}:${TAG}" \
        --build-arg VERSION="${VERSION}" \
        --build-arg BUILD_TIME="${BUILD_TIME}" \
        --build-arg COMMIT_SHA="${COMMIT_SHA}" \
        -f docker/Dockerfile \
        . 2>&1 | tee docker-build.log
    
    # 檢查建置結果
    if [ ${PIPESTATUS[0]} -eq 0 ]; then
        log_success "Docker 建置成功！"
        log_info "映像檔: ${IMAGE_NAME}:${TAG}"
        
        # 顯示映像檔資訊
        log_info "映像檔資訊："
        docker images "${IMAGE_NAME}:${TAG}" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
        
        # 清理建置快取（可選）
        if [ "${CLEANUP:-false}" = "true" ]; then
            log_info "清理 Docker 建置快取..."
            docker builder prune -f
        fi
        
    else
        log_error "Docker 建置失敗！"
        log_info "請檢查 docker-build.log 檔案以獲取詳細錯誤資訊"
        exit 1
    fi
}

# 顯示使用說明
show_usage() {
    echo "使用方法: $0 [選項]"
    echo ""
    echo "選項："
    echo "  -v, --version VERSION    設定版本號 (預設: v1.0.0)"
    echo "  -t, --tag TAG            設定標籤 (預設: latest)"
    echo "  -c, --cleanup            建置完成後清理快取"
    echo "  -h, --help               顯示此說明"
    echo ""
    echo "範例："
    echo "  $0                        # 使用預設參數建置"
    echo "  $0 -v v1.1.0 -t stable    # 建置 v1.1.0 版本並標籤為 stable"
    echo "  $0 -c                     # 建置並清理快取"
}

# 解析命令列參數
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        -c|--cleanup)
            CLEANUP="true"
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            log_error "未知選項: $1"
            show_usage
            exit 1
            ;;
    esac
done

# 執行主函數
main "$@"
