#!/bin/bash

# =============================================================================
# wait-for-it.sh - 等待服務啟動的通用腳本
# 用於確保依賴服務在啟動前已經就緒
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

# 顯示使用說明
usage() {
    echo "用法: $0 <host>:<port> [--timeout=<timeout>] [--strict] [--quiet] [--command=<command>]"
    echo ""
    echo "選項:"
    echo "  --timeout=<timeout>    等待超時時間（秒），預設為 15"
    echo "  --strict              嚴格模式，如果服務未就緒則退出"
    echo "  --quiet               安靜模式，不輸出日誌"
    echo "  --command=<command>   服務就緒後執行的命令"
    echo ""
    echo "範例:"
    echo "  $0 localhost:5432"
    echo "  $0 postgres:5432 --timeout=30"
    echo "  $0 redis:6379 --command='echo \"Redis is ready\"'"
    exit 1
}

# 解析參數
parse_args() {
    TIMEOUT=15
    STRICT=false
    QUIET=false
    COMMAND=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --timeout=*)
                TIMEOUT="${1#*=}"
                shift
                ;;
            --strict)
                STRICT=true
                shift
                ;;
            --quiet)
                QUIET=true
                shift
                ;;
            --command=*)
                COMMAND="${1#*=}"
                shift
                ;;
            -h|--help)
                usage
                ;;
            *)
                if [[ -z "$HOST_PORT" ]]; then
                    HOST_PORT="$1"
                else
                    log_error "未知參數: $1"
                    usage
                fi
                shift
                ;;
        esac
    done
    
    if [[ -z "$HOST_PORT" ]]; then
        log_error "缺少主機和端口參數"
        usage
    fi
    
    # 解析主機和端口
    if [[ "$HOST_PORT" =~ ^([^:]+):([0-9]+)$ ]]; then
        HOST="${BASH_REMATCH[1]}"
        PORT="${BASH_REMATCH[2]}"
    else
        log_error "無效的主機端口格式: $HOST_PORT"
        usage
    fi
}

# 檢查服務是否就緒
check_service() {
    local host=$1
    local port=$2
    
    # 使用 nc 檢查端口
    if command -v nc >/dev/null 2>&1; then
        nc -z "$host" "$port" 2>/dev/null
        return $?
    fi
    
    # 使用 /dev/tcp 作為備選方案
    if timeout 1 bash -c "</dev/tcp/$host/$port" 2>/dev/null; then
        return 0
    fi
    
    return 1
}

# 等待服務啟動
wait_for_service() {
    local host=$1
    local port=$2
    local timeout=$3
    local attempt=1
    
    if [[ "$QUIET" != "true" ]]; then
        log_info "等待服務啟動: $host:$port (超時: ${timeout}s)"
    fi
    
    while [[ $attempt -le $timeout ]]; do
        if check_service "$host" "$port"; then
            if [[ "$QUIET" != "true" ]]; then
                log_success "服務已就緒: $host:$port"
            fi
            return 0
        fi
        
        if [[ "$QUIET" != "true" ]]; then
            log_info "嘗試 $attempt/$timeout - 服務尚未就緒，等待 1 秒..."
        fi
        
        sleep 1
        attempt=$((attempt + 1))
    done
    
    if [[ "$QUIET" != "true" ]]; then
        log_error "服務啟動超時: $host:$port"
    fi
    
    return 1
}

# 主函數
main() {
    # 解析命令行參數
    parse_args "$@"
    
    # 等待服務啟動
    if wait_for_service "$HOST" "$PORT" "$TIMEOUT"; then
        # 服務已就緒，執行命令（如果指定）
        if [[ -n "$COMMAND" ]]; then
            if [[ "$QUIET" != "true" ]]; then
                log_info "執行命令: $COMMAND"
            fi
            eval "$COMMAND"
        fi
        exit 0
    else
        # 服務未就緒
        if [[ "$STRICT" == "true" ]]; then
            log_error "嚴格模式：服務未就緒，退出"
            exit 1
        else
            log_warning "服務未就緒，但非嚴格模式，繼續執行"
            exit 0
        fi
    fi
}

# 執行主函數
main "$@"
