#!/bin/bash

# =============================================================================
# 監控服務初始化腳本
# 用於初始化 Prometheus 和 Grafana 監控服務
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

# 建立 Prometheus 配置
create_prometheus_config() {
    log_info "建立 Prometheus 配置..."
    
    cat > /etc/prometheus/prometheus.yml <<EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "alert_rules.yml"

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'security-intel-backend'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s

  - job_name: 'postgres'
    static_configs:
      - targets: ['localhost:5432']
    metrics_path: '/metrics'

  - job_name: 'redis'
    static_configs:
      - targets: ['localhost:6379']
    metrics_path: '/metrics'

  - job_name: 'mosquitto'
    static_configs:
      - targets: ['localhost:1883']
    metrics_path: '/metrics'
EOF

    # 建立告警規則
    cat > /etc/prometheus/alert_rules.yml <<EOF
groups:
  - name: security_intel_alerts
    rules:
      - alert: BackendServiceDown
        expr: up{job="security-intel-backend"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Backend service is down"
          description: "The security intelligence backend service has been down for more than 1 minute"

      - alert: DatabaseServiceDown
        expr: up{job="postgres"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database service is down"
          description: "The PostgreSQL database service has been down for more than 1 minute"

      - alert: CacheServiceDown
        expr: up{job="redis"} == 0
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Cache service is down"
          description: "The Redis cache service has been down for more than 1 minute"
EOF

    log_success "Prometheus 配置建立完成"
}

# 建立 Grafana 配置
create_grafana_config() {
    log_info "建立 Grafana 配置..."
    
    # 建立資料來源配置
    mkdir -p /etc/grafana/provisioning/datasources
    cat > /etc/grafana/provisioning/datasources/prometheus.yml <<EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://localhost:9090
    isDefault: true
    editable: true
EOF

    # 建立儀表板配置
    mkdir -p /etc/grafana/provisioning/dashboards
    cat > /etc/grafana/provisioning/dashboards/dashboards.yml <<EOF
apiVersion: 1

providers:
  - name: 'Security Intelligence'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /var/lib/grafana/dashboards
EOF

    # 建立基本儀表板
    mkdir -p /var/lib/grafana/dashboards
    cat > /var/lib/grafana/dashboards/security-intelligence-overview.json <<EOF
{
  "dashboard": {
    "id": null,
    "title": "Security Intelligence Overview",
    "tags": ["security", "intelligence"],
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Service Status",
        "type": "stat",
        "targets": [
          {
            "expr": "up",
            "legendFormat": "{{job}}"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 0, "y": 0}
      },
      {
        "id": 2,
        "title": "HTTP Requests",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{handler}}"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 12, "y": 0}
      }
    ],
    "time": {"from": "now-1h", "to": "now"},
    "refresh": "30s"
  }
}
EOF

    log_success "Grafana 配置建立完成"
}

# 啟動監控服務
start_monitoring_services() {
    log_info "啟動監控服務..."
    
    # 啟動 Prometheus
    su nobody -s /bin/bash -c "prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/var/lib/prometheus/data --storage.tsdb.retention.time=72h --web.enable-lifecycle" &
    PROMETHEUS_PID=$!
    
    # 等待 Prometheus 啟動
    sleep 5
    
    # 啟動 Grafana
    su grafana -s /bin/bash -c "grafana-server --config=/etc/grafana/grafana.ini --homepath=/usr/share/grafana" &
    GRAFANA_PID=$!
    
    # 等待 Grafana 啟動
    sleep 5
    
    log_success "監控服務啟動完成"
}

# 主函數
main() {
    log_info "開始初始化監控服務"
    
    # 建立 Prometheus 配置
    create_prometheus_config
    
    # 建立 Grafana 配置
    create_grafana_config
    
    # 啟動監控服務
    start_monitoring_services
    
    log_success "監控服務初始化完成"
    
    # 等待進程
    wait $PROMETHEUS_PID $GRAFANA_PID
}

# 執行主函數
main "$@"
