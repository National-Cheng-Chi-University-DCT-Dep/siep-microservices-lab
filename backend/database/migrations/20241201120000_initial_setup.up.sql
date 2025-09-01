-- 建立基本的 enum 類型
CREATE TYPE user_role AS ENUM ('admin', 'premium', 'basic');
CREATE TYPE threat_type AS ENUM ('malware', 'phishing', 'spam', 'botnet', 'scanner', 'ddos', 'bruteforce', 'other');
CREATE TYPE severity_level AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE collection_status AS ENUM ('pending', 'in_progress', 'completed', 'failed');

-- 啟用必要的 PostgreSQL 擴展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 建立使用者表
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role DEFAULT 'basic',
    is_active BOOLEAN DEFAULT true,
    subscription_expires_at TIMESTAMP,
    api_quota INTEGER DEFAULT 1000,
    api_usage INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- 建立威脅情報表
CREATE TABLE threat_intelligence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address INET NOT NULL,
    domain VARCHAR(253),
    threat_type threat_type NOT NULL,
    severity severity_level NOT NULL,
    confidence_score INTEGER CHECK (confidence_score >= 0 AND confidence_score <= 100),
    description TEXT,
    source VARCHAR(100) NOT NULL,
    external_id VARCHAR(100),
    country_code VARCHAR(2),
    asn INTEGER,
    isp VARCHAR(200),
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    tags TEXT[],
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 建立情報來源表
CREATE TABLE intelligence_sources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    url VARCHAR(500),
    api_key_required BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    collection_interval INTEGER DEFAULT 3600, -- 秒
    last_collection TIMESTAMP,
    total_collected INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 建立收集任務表
CREATE TABLE collection_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_id UUID NOT NULL REFERENCES intelligence_sources(id) ON DELETE CASCADE,
    status collection_status DEFAULT 'pending',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    records_collected INTEGER DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 建立 API 金鑰表
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP,
    quota INTEGER DEFAULT 1000,
    usage INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP
);

-- 建立 API 請求日誌表
CREATE TABLE api_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    api_key_id UUID REFERENCES api_keys(id) ON DELETE SET NULL,
    endpoint VARCHAR(200) NOT NULL,
    method VARCHAR(10) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    response_code INTEGER,
    response_time INTEGER, -- 毫秒
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引以提高查詢效能
-- 威脅情報表索引
CREATE INDEX idx_threat_intelligence_ip ON threat_intelligence(ip_address);
CREATE INDEX idx_threat_intelligence_domain ON threat_intelligence(domain);
CREATE INDEX idx_threat_intelligence_type ON threat_intelligence(threat_type);
CREATE INDEX idx_threat_intelligence_severity ON threat_intelligence(severity);
CREATE INDEX idx_threat_intelligence_source ON threat_intelligence(source);
CREATE INDEX idx_threat_intelligence_created_at ON threat_intelligence(created_at);
CREATE INDEX idx_threat_intelligence_last_seen ON threat_intelligence(last_seen);
CREATE INDEX idx_threat_intelligence_country ON threat_intelligence(country_code);
CREATE INDEX idx_threat_intelligence_tags ON threat_intelligence USING GIN(tags);
CREATE INDEX idx_threat_intelligence_metadata ON threat_intelligence USING GIN(metadata);

-- 使用者表索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(is_active);

-- API 相關索引
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_active ON api_keys(is_active);
CREATE INDEX idx_api_requests_user_id ON api_requests(user_id);
CREATE INDEX idx_api_requests_created_at ON api_requests(created_at);
CREATE INDEX idx_api_requests_endpoint ON api_requests(endpoint);

-- 收集任務索引
CREATE INDEX idx_collection_jobs_source_id ON collection_jobs(source_id);
CREATE INDEX idx_collection_jobs_status ON collection_jobs(status);
CREATE INDEX idx_collection_jobs_created_at ON collection_jobs(created_at);

-- 建立觸發器以自動更新 updated_at 欄位
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 套用觸發器到各表
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_threat_intelligence_updated_at BEFORE UPDATE ON threat_intelligence
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_intelligence_sources_updated_at BEFORE UPDATE ON intelligence_sources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_collection_jobs_updated_at BEFORE UPDATE ON collection_jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 插入初始情報來源
INSERT INTO intelligence_sources (name, url, api_key_required, collection_interval) VALUES 
('AbuseIPDB', 'https://api.abuseipdb.com/api/v2', true, 3600),
('Manual Entry', NULL, false, 0),
('Internal Detection', NULL, false, 0)
ON CONFLICT (name) DO NOTHING;

-- 插入初始管理員使用者（僅開發環境）
INSERT INTO users (username, email, password_hash, role) VALUES 
('admin', 'admin@security-intel.com', crypt('admin123', gen_salt('bf')), 'admin'),
('demo', 'demo@security-intel.com', crypt('demo123', gen_salt('bf')), 'basic')
ON CONFLICT (username) DO NOTHING;

-- 插入一些測試威脅情報資料
INSERT INTO threat_intelligence (ip_address, threat_type, severity, confidence_score, description, source, country_code) VALUES 
('192.168.1.100', 'malware', 'high', 85, '測試惡意軟體 IP', 'Manual Entry', 'TW'),
('10.0.0.1', 'scanner', 'medium', 70, '測試掃描器 IP', 'Manual Entry', 'US'),
('172.16.0.1', 'phishing', 'critical', 95, '測試釣魚網站 IP', 'Manual Entry', 'CN'),
('203.0.113.1', 'ddos', 'high', 88, '測試 DDoS 攻擊來源', 'Manual Entry', 'RU'),
('198.51.100.1', 'bruteforce', 'medium', 75, '測試暴力破解攻擊', 'Manual Entry', 'BR')
ON CONFLICT DO NOTHING; 