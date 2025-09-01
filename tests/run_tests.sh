#!/bin/bash

# Robot Framework 測試執行腳本
# 作者: Ultimate Security Intelligence Platform Team
# 版本: 1.0

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"
TESTS_DIR="$SCRIPT_DIR"
RESULTS_DIR="$TESTS_DIR/results"
VENV_DIR="$TESTS_DIR/venv"

# 預設值
TEST_SUITE=""
TAGS=""
VARIABLES=""
PARALLEL=false
VERBOSE=false
CLEAN=false
SETUP_ONLY=false

# 顯示幫助
show_help() {
    echo -e "${BLUE}Robot Framework 測試執行器${NC}"
    echo ""
    echo "用法: $0 [選項]"
    echo ""
    echo "選項:"
    echo "  -s, --suite SUITE     執行特定測試套件 (auth, threat, collector, all)"
    echo "  -t, --tags TAGS       執行帶有特定標籤的測試 (positive, negative, etc.)"
    echo "  -v, --variables FILE  載入變數檔案"
    echo "  -p, --parallel        並行執行測試"
    echo "  -V, --verbose         顯示詳細輸出"
    echo "  -c, --clean           清理之前的測試結果"
    echo "      --setup-only      只設定環境，不執行測試"
    echo "  -h, --help            顯示此幫助訊息"
    echo ""
    echo "範例:"
    echo "  $0 -s auth                    # 執行認證測試"
    echo "  $0 -t positive               # 執行所有正向測試"
    echo "  $0 -s threat -t create       # 執行威脅情報建立測試"
    echo "  $0 -p -V                     # 並行執行所有測試，顯示詳細輸出"
    echo ""
}

# 日誌函數
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    if [ "$VERBOSE" = true ]; then
        echo -e "${BLUE}[DEBUG]${NC} $1"
    fi
}

# 檢查依賴
check_dependencies() {
    log_info "檢查系統依賴..."
    
    # 檢查Python
    if ! command -v python3 &> /dev/null; then
        log_error "Python 3 未安裝"
        exit 1
    fi
    
    # 檢查pip
    if ! command -v pip3 &> /dev/null; then
        log_error "pip3 未安裝"
        exit 1
    fi
    
    # 檢查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安裝"
        exit 1
    fi
    
    # 檢查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安裝"
        exit 1
    fi
    
    # 檢查後端伺服器
    if [ ! -f "$BACKEND_DIR/cmd/server/main.go" ]; then
        log_error "後端專案未找到: $BACKEND_DIR"
        exit 1
    fi
    
    log_info "所有依賴檢查通過"
}

# 設定Python虛擬環境
setup_virtualenv() {
    log_info "設定Python虛擬環境..."
    
    if [ ! -d "$VENV_DIR" ]; then
        log_debug "建立虛擬環境: $VENV_DIR"
        python3 -m venv "$VENV_DIR"
    fi
    
    log_debug "啟動虛擬環境"
    source "$VENV_DIR/bin/activate"
    
    log_debug "升級pip"
    pip install --upgrade pip
    
    log_debug "安裝測試依賴"
    pip install -r "$TESTS_DIR/requirements.txt"
    
    log_info "Python環境設定完成"
}

# 啟動服務
start_services() {
    log_info "啟動測試服務..."
    
    # 啟動Docker服務
    log_debug "啟動Docker服務"
    cd "$PROJECT_ROOT"
    docker-compose -f docker/docker-compose.yml up -d
    
    # 等待資料庫啟動
    log_debug "等待資料庫啟動"
    sleep 10
    
    # 執行資料庫遷移
    log_debug "執行資料庫遷移"
    cd "$BACKEND_DIR"
    make migrate-up || true
    
    # 啟動後端服務
    log_debug "啟動後端服務"
    go run cmd/server/main.go &
    BACKEND_PID=$!
    echo $BACKEND_PID > /tmp/backend_test.pid
    
    # 等待後端服務啟動
    log_debug "等待後端服務啟動"
    for i in {1..30}; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            log_info "後端服務已啟動"
            break
        fi
        sleep 2
        if [ $i -eq 30 ]; then
            log_error "後端服務啟動超時"
            exit 1
        fi
    done
}

# 停止服務
stop_services() {
    log_info "停止測試服務..."
    
    # 停止後端服務
    if [ -f /tmp/backend_test.pid ]; then
        BACKEND_PID=$(cat /tmp/backend_test.pid)
        if ps -p $BACKEND_PID > /dev/null 2>&1; then
            log_debug "停止後端服務 (PID: $BACKEND_PID)"
            kill $BACKEND_PID
            sleep 2
            if ps -p $BACKEND_PID > /dev/null 2>&1; then
                kill -9 $BACKEND_PID
            fi
        fi
        rm -f /tmp/backend_test.pid
    fi
    
    # 停止Docker服務
    log_debug "停止Docker服務"
    cd "$PROJECT_ROOT"
    docker-compose -f docker/docker-compose.yml down
}

# 清理測試結果
clean_results() {
    if [ "$CLEAN" = true ]; then
        log_info "清理之前的測試結果..."
        rm -rf "$RESULTS_DIR"
    fi
}

# 準備結果目錄
prepare_results_dir() {
    log_debug "準備結果目錄: $RESULTS_DIR"
    mkdir -p "$RESULTS_DIR"
    
    # 建立時間戳目錄
    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    CURRENT_RUN_DIR="$RESULTS_DIR/run_$TIMESTAMP"
    mkdir -p "$CURRENT_RUN_DIR"
    
    # 建立符號連結到最新結果
    ln -sfn "$CURRENT_RUN_DIR" "$RESULTS_DIR/latest"
    
    echo "$CURRENT_RUN_DIR"
}

# 執行測試
run_tests() {
    local results_dir=$1
    
    log_info "執行Robot Framework測試..."
    
    # 啟動虛擬環境
    source "$VENV_DIR/bin/activate"
    
    # 建構robot命令
    local robot_cmd="robot"
    local robot_args=""
    
    # 設定輸出目錄
    robot_args="$robot_args --outputdir $results_dir"
    
    # 設定日誌級別
    if [ "$VERBOSE" = true ]; then
        robot_args="$robot_args --loglevel DEBUG"
    else
        robot_args="$robot_args --loglevel INFO"
    fi
    
    # 添加標籤過濾
    if [ -n "$TAGS" ]; then
        robot_args="$robot_args --include $TAGS"
    fi
    
    # 添加變數檔案
    if [ -n "$VARIABLES" ]; then
        robot_args="$robot_args --variablefile $VARIABLES"
    fi
    
    # 設定並行執行
    if [ "$PARALLEL" = true ]; then
        robot_args="$robot_args --processes 4"
    fi
    
    # 選擇測試套件
    local test_files=""
    case "$TEST_SUITE" in
        "auth")
            test_files="$TESTS_DIR/api/auth_tests.robot"
            ;;
        "threat")
            test_files="$TESTS_DIR/api/threat_intelligence_tests.robot"
            ;;
        "collector")
            test_files="$TESTS_DIR/api/collector_tests.robot"
            ;;
        "all"|"")
            test_files="$TESTS_DIR/api/"
            ;;
        *)
            log_error "未知的測試套件: $TEST_SUITE"
            exit 1
            ;;
    esac
    
    # 執行測試
    log_debug "執行命令: $robot_cmd $robot_args $test_files"
    
    set +e
    $robot_cmd $robot_args $test_files
    local exit_code=$?
    set -e
    
    # 生成測試報告摘要
    generate_summary "$results_dir" "$exit_code"
    
    return $exit_code
}

# 生成測試摘要
generate_summary() {
    local results_dir=$1
    local exit_code=$2
    
    log_info "生成測試摘要..."
    
    local output_file="$results_dir/output.xml"
    local report_file="$results_dir/report.html"
    local log_file="$results_dir/log.html"
    
    if [ -f "$output_file" ]; then
        # 使用rebot生成更好的報告
        source "$VENV_DIR/bin/activate"
        rebot --outputdir "$results_dir" --name "Security Intelligence Platform Tests" "$output_file"
        
        # 提取測試統計
        local total_tests=$(grep -o 'stat[^>]*>.*</stat>' "$output_file" | head -1 | sed 's/.*>\(.*\)<.*/\1/' || echo "Unknown")
        local passed_tests=$(grep -o 'stat[^>]*pass="[0-9]*"' "$output_file" | head -1 | sed 's/.*pass="\([0-9]*\)".*/\1/' || echo "0")
        local failed_tests=$(grep -o 'stat[^>]*fail="[0-9]*"' "$output_file" | head -1 | sed 's/.*fail="\([0-9]*\)".*/\1/' || echo "0")
        
        # 建立摘要檔案
        cat > "$results_dir/summary.txt" << EOF
測試執行摘要
============
執行時間: $(date)
測試套件: ${TEST_SUITE:-all}
標籤過濾: ${TAGS:-none}
並行執行: $PARALLEL

測試結果:
- 總計: $total_tests
- 通過: $passed_tests
- 失敗: $failed_tests
- 退出碼: $exit_code

報告檔案:
- HTML報告: $report_file
- 詳細日誌: $log_file
- XML輸出: $output_file
EOF
        
        log_info "測試完成 - 通過: $passed_tests, 失敗: $failed_tests"
        log_info "報告位置: $report_file"
    else
        log_warn "未找到測試輸出檔案: $output_file"
    fi
}

# 清理函數
cleanup() {
    log_info "執行清理..."
    stop_services
    
    # 停用虛擬環境
    if [ -n "$VIRTUAL_ENV" ]; then
        deactivate
    fi
}

# 主要執行函數
main() {
    # 解析命令列參數
    while [[ $# -gt 0 ]]; do
        case $1 in
            -s|--suite)
                TEST_SUITE="$2"
                shift 2
                ;;
            -t|--tags)
                TAGS="$2"
                shift 2
                ;;
            -v|--variables)
                VARIABLES="$2"
                shift 2
                ;;
            -p|--parallel)
                PARALLEL=true
                shift
                ;;
            -V|--verbose)
                VERBOSE=true
                shift
                ;;
            -c|--clean)
                CLEAN=true
                shift
                ;;
            --setup-only)
                SETUP_ONLY=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "未知選項: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 設定清理處理器
    trap cleanup EXIT
    
    log_info "開始Robot Framework測試執行"
    
    # 檢查依賴
    check_dependencies
    
    # 清理舊結果
    clean_results
    
    # 設定環境
    setup_virtualenv
    
    # 準備結果目錄
    local results_dir=$(prepare_results_dir)
    
    # 啟動服務
    start_services
    
    if [ "$SETUP_ONLY" = true ]; then
        log_info "環境設定完成，使用 --setup-only 選項，跳過測試執行"
        log_info "後端服務運行在: http://localhost:8080"
        log_info "使用 Ctrl+C 停止服務"
        
        # 等待用戶中斷
        while true; do
            sleep 1
        done
    else
        # 執行測試
        run_tests "$results_dir"
        local test_exit_code=$?
        
        log_info "測試執行完成，退出碼: $test_exit_code"
        exit $test_exit_code
    fi
}

# 執行主函數
main "$@" 