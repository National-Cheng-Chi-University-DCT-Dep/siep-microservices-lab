#!/bin/bash

# =============================================================================
# 資料庫初始化腳本
# 用於初始化 PostgreSQL 資料庫和建立必要的表格
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
    log_info "檢查資料庫環境變數..."
    
    if [[ -z "$POSTGRES_DB" ]]; then
        log_error "缺少 POSTGRES_DB 環境變數"
        exit 1
    fi
    
    if [[ -z "$POSTGRES_USER" ]]; then
        log_error "缺少 POSTGRES_USER 環境變數"
        exit 1
    fi
    
    if [[ -z "$POSTGRES_PASSWORD" ]]; then
        log_error "缺少 POSTGRES_PASSWORD 環境變數"
        exit 1
    fi
    
    log_success "環境變數檢查完成"
}

# 建立資料庫使用者
create_user() {
    log_info "建立資料庫使用者: $POSTGRES_USER"
    
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        CREATE USER IF NOT EXISTS $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';
        ALTER USER $POSTGRES_USER CREATEDB;
        GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO $POSTGRES_USER;
EOSQL
    
    log_success "資料庫使用者建立完成"
}

# 建立資料庫
create_database() {
    log_info "建立資料庫: $POSTGRES_DB"
    
    createdb --username "$POSTGRES_USER" --owner "$POSTGRES_USER" "$POSTGRES_DB" 2>/dev/null || log_warning "資料庫已存在"
    
    log_success "資料庫建立完成"
}

# 建立基本表格
create_tables() {
    log_info "建立基本表格..."
    
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        -- 使用者表格
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            role VARCHAR(20) DEFAULT 'user',
            is_active BOOLEAN DEFAULT true,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        -- 威脅情報表格
        CREATE TABLE IF NOT EXISTS threat_intelligence (
            id SERIAL PRIMARY KEY,
            indicator VARCHAR(255) NOT NULL,
            indicator_type VARCHAR(50) NOT NULL,
            threat_type VARCHAR(100),
            confidence_level INTEGER CHECK (confidence_level >= 1 AND confidence_level <= 10),
            source VARCHAR(255),
            description TEXT,
            tags TEXT[],
            first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        -- 安全事件表格
        CREATE TABLE IF NOT EXISTS security_events (
            id SERIAL PRIMARY KEY,
            event_type VARCHAR(100) NOT NULL,
            severity VARCHAR(20) NOT NULL,
            source_ip INET,
            destination_ip INET,
            source_port INTEGER,
            destination_port INTEGER,
            protocol VARCHAR(10),
            description TEXT,
            raw_data JSONB,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        -- 資產管理表格
        CREATE TABLE IF NOT EXISTS assets (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            asset_type VARCHAR(50) NOT NULL,
            ip_address INET,
            mac_address MACADDR,
            hostname VARCHAR(255),
            os_info TEXT,
            status VARCHAR(20) DEFAULT 'active',
            tags TEXT[],
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        -- 漏洞管理表格
        CREATE TABLE IF NOT EXISTS vulnerabilities (
            id SERIAL PRIMARY KEY,
            cve_id VARCHAR(20),
            title VARCHAR(255) NOT NULL,
            description TEXT,
            severity VARCHAR(20) NOT NULL,
            cvss_score DECIMAL(3,1),
            affected_assets INTEGER[],
            status VARCHAR(20) DEFAULT 'open',
            discovered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            resolved_at TIMESTAMP,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        -- 日誌表格
        CREATE TABLE IF NOT EXISTS audit_logs (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id),
            action VARCHAR(100) NOT NULL,
            resource_type VARCHAR(50),
            resource_id INTEGER,
            details JSONB,
            ip_address INET,
            user_agent TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
EOSQL
    
    log_success "基本表格建立完成"
}

# 建立索引
create_indexes() {
    log_info "建立資料庫索引..."
    
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        -- 威脅情報索引
        CREATE INDEX IF NOT EXISTS idx_threat_intelligence_indicator ON threat_intelligence(indicator);
        CREATE INDEX IF NOT EXISTS idx_threat_intelligence_type ON threat_intelligence(indicator_type);
        CREATE INDEX IF NOT EXISTS idx_threat_intelligence_created_at ON threat_intelligence(created_at);
        
        -- 安全事件索引
        CREATE INDEX IF NOT EXISTS idx_security_events_type ON security_events(event_type);
        CREATE INDEX IF NOT EXISTS idx_security_events_severity ON security_events(severity);
        CREATE INDEX IF NOT EXISTS idx_security_events_created_at ON security_events(created_at);
        CREATE INDEX IF NOT EXISTS idx_security_events_source_ip ON security_events(source_ip);
        
        -- 資產索引
        CREATE INDEX IF NOT EXISTS idx_assets_name ON assets(name);
        CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(asset_type);
        CREATE INDEX IF NOT EXISTS idx_assets_ip ON assets(ip_address);
        
        -- 漏洞索引
        CREATE INDEX IF NOT EXISTS idx_vulnerabilities_cve ON vulnerabilities(cve_id);
        CREATE INDEX IF NOT EXISTS idx_vulnerabilities_severity ON vulnerabilities(severity);
        CREATE INDEX IF NOT EXISTS idx_vulnerabilities_status ON vulnerabilities(status);
        
        -- 審計日誌索引
        CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
        CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
        CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
EOSQL
    
    log_success "資料庫索引建立完成"
}

# 插入初始資料
insert_initial_data() {
    log_info "插入初始資料..."
    
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        -- 建立預設管理員使用者
        INSERT INTO users (username, email, password_hash, role) 
        VALUES ('admin', 'admin@security-intel.com', '\$2a\$10\$default_hash_placeholder', 'admin')
        ON CONFLICT (username) DO NOTHING;
        
        -- 插入範例威脅情報
        INSERT INTO threat_intelligence (indicator, indicator_type, threat_type, confidence_level, source, description, tags) 
        VALUES 
            ('192.168.1.100', 'ip', 'malware', 8, 'internal_scan', '可疑的內部 IP 地址', ARRAY['malware', 'internal']),
            ('example.com', 'domain', 'phishing', 7, 'external_feed', '釣魚網站域名', ARRAY['phishing', 'external'])
        ON CONFLICT DO NOTHING;
        
        -- 插入範例資產
        INSERT INTO assets (name, asset_type, ip_address, hostname, os_info, status) 
        VALUES 
            ('Web Server 01', 'server', '192.168.1.10', 'web01.local', 'Ubuntu 20.04 LTS', 'active'),
            ('Database Server', 'server', '192.168.1.20', 'db01.local', 'CentOS 8', 'active')
        ON CONFLICT DO NOTHING;
EOSQL
    
    log_success "初始資料插入完成"
}

# 主函數
main() {
    log_info "開始初始化 PostgreSQL 資料庫"
    
    # 檢查環境變數
    check_env
    
    # 建立資料庫使用者
    create_user
    
    # 建立資料庫
    create_database
    
    # 建立基本表格
    create_tables
    
    # 建立索引
    create_indexes
    
    # 插入初始資料
    insert_initial_data
    
    log_success "PostgreSQL 資料庫初始化完成"
}

# 執行主函數
main "$@"
