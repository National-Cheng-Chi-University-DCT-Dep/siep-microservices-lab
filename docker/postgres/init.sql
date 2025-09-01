-- 資安情報平台資料庫初始化腳本
-- 此腳本會在 PostgreSQL 容器首次啟動時自動執行

-- 建立資料庫 (如果不存在)
-- 由於 POSTGRES_DB 環境變數已經建立了 security_intel 資料庫，這裡不需要重複建立

-- 啟用必要的 PostgreSQL 擴展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 建立基本的 enum 類型
CREATE TYPE user_role AS ENUM ('admin', 'premium', 'basic');
CREATE TYPE threat_type AS ENUM ('malware', 'phishing', 'spam', 'botnet', 'scanner', 'other');
CREATE TYPE severity_level AS ENUM ('low', 'medium', 'high', 'critical');

-- 建立使用者表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role DEFAULT 'basic',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- 建立威脅情報表
CREATE TABLE IF NOT EXISTS threat_intelligence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address INET NOT NULL,
    threat_type threat_type NOT NULL,
    severity severity_level NOT NULL,
    confidence_score INTEGER CHECK (confidence_score >= 0 AND confidence_score <= 100),
    description TEXT,
    source VARCHAR(100),
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 建立 API 金鑰表
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP
);

-- 建立索引以提高查詢效能
CREATE INDEX IF NOT EXISTS idx_threat_intelligence_ip ON threat_intelligence(ip_address);
CREATE INDEX IF NOT EXISTS idx_threat_intelligence_type ON threat_intelligence(threat_type);
CREATE INDEX IF NOT EXISTS idx_threat_intelligence_severity ON threat_intelligence(severity);
CREATE INDEX IF NOT EXISTS idx_threat_intelligence_created_at ON threat_intelligence(created_at);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);

-- 建立觸發器以自動更新 updated_at 欄位
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_threat_intelligence_updated_at BEFORE UPDATE ON threat_intelligence
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 插入一些初始測試資料
INSERT INTO users (username, email, password_hash, role) VALUES 
('admin', 'admin@security-intel.com', crypt('admin123', gen_salt('bf')), 'admin'),
('testuser', 'test@security-intel.com', crypt('test123', gen_salt('bf')), 'basic')
ON CONFLICT (username) DO NOTHING;

INSERT INTO threat_intelligence (ip_address, threat_type, severity, confidence_score, description, source) VALUES 
('192.168.1.100', 'malware', 'high', 85, '測試惡意軟體 IP', 'manual'),
('10.0.0.1', 'scanner', 'medium', 70, '測試掃描器 IP', 'manual'),
('172.16.0.1', 'phishing', 'critical', 95, '測試釣魚網站 IP', 'manual')
ON CONFLICT DO NOTHING;

-- 顯示初始化完成訊息
DO $$
BEGIN
    RAISE NOTICE '資安情報平台資料庫初始化完成！';
    RAISE NOTICE '預設管理員帳戶: admin@security-intel.com / admin123';
    RAISE NOTICE '預設測試帳戶: test@security-intel.com / test123';
END $$; 