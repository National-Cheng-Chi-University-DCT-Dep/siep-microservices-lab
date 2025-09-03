#!/bin/bash

# =============================================================================
# MQTT 初始化腳本
# 用於初始化 Mosquitto MQTT 訊息佇列服務
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
    log_info "檢查 MQTT 環境變數..."
    
    if [[ -z "$MQTT_USERNAME" ]]; then
        log_warning "未設定 MQTT_USERNAME，使用預設配置"
    fi
    
    if [[ -z "$MQTT_PASSWORD" ]]; then
        log_warning "未設定 MQTT_PASSWORD，使用預設配置"
    fi
    
    log_success "環境變數檢查完成"
}

# 建立 Mosquitto 配置
create_config() {
    log_info "建立 Mosquitto 配置..."
    
    # 建立配置目錄
    mkdir -p /etc/mosquitto/conf.d
    
    # 主配置檔案
    cat > /etc/mosquitto/mosquitto.conf <<EOF
# Mosquitto MQTT Broker 配置
pid_file /var/run/mosquitto.pid
persistence true
persistence_location /mosquitto/data/
log_dest file /mosquitto/log/mosquitto.log
log_dest stdout
log_type all
log_timestamp true
connection_messages true
log_timestamp_format %Y-%m-%dT%H:%M:%S

# 監聽設定
listener 1883
protocol mqtt

# WebSocket 支援
listener 9001
protocol websockets

# 允許匿名連接（開發環境）
allow_anonymous true

# 最大連接數
max_connections 1000

# 連接超時
connection_messages true
log_connections true
log_disconnections true

# 保留訊息設定
max_queued_messages 100
max_inflight_messages 20
EOF

    # 如果設定了認證，添加認證配置
    if [[ -n "$MQTT_USERNAME" && -n "$MQTT_PASSWORD" ]]; then
        cat > /etc/mosquitto/conf.d/auth.conf <<EOF
# 認證配置
allow_anonymous false
password_file /etc/mosquitto/passwd
EOF
        
        # 建立密碼檔案
        mosquitto_passwd -c /etc/mosquitto/passwd "$MQTT_USERNAME" <<< "$MQTT_PASSWORD"
    fi
    
    log_success "Mosquitto 配置建立完成"
}

# 建立預設主題
create_topics() {
    log_info "建立預設 MQTT 主題..."
    
    # 等待 Mosquitto 啟動
    sleep 3
    
    # 建立預設主題結構
    mosquitto_pub -h localhost -t "security_intel/status" -m "online" -r
    mosquitto_pub -h localhost -t "security_intel/version" -m "1.0.0" -r
    mosquitto_pub -h localhost -t "security_intel/environment" -m "sit" -r
    
    log_success "預設主題建立完成"
}

# 主函數
main() {
    log_info "開始初始化 MQTT 服務"
    
    # 檢查環境變數
    check_env
    
    # 建立配置
    create_config
    
    # 啟動 Mosquitto
    /usr/sbin/mosquitto -c /etc/mosquitto/mosquitto.conf &
    MOSQUITTO_PID=$!
    
    # 等待 Mosquitto 啟動
    sleep 3
    
    # 建立預設主題
    create_topics
    
    log_success "MQTT 服務初始化完成"
    
    # 等待 Mosquitto 進程
    wait $MOSQUITTO_PID
}

# 執行主函數
main "$@"
