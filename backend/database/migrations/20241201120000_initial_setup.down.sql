-- 刪除觸發器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_threat_intelligence_updated_at ON threat_intelligence;
DROP TRIGGER IF EXISTS update_intelligence_sources_updated_at ON intelligence_sources;
DROP TRIGGER IF EXISTS update_collection_jobs_updated_at ON collection_jobs;

-- 刪除函數
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 刪除表格（注意外鍵約束的順序）
DROP TABLE IF EXISTS api_requests;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS collection_jobs;
DROP TABLE IF EXISTS threat_intelligence;
DROP TABLE IF EXISTS intelligence_sources;
DROP TABLE IF EXISTS users;

-- 刪除 enum 類型
DROP TYPE IF EXISTS collection_status;
DROP TYPE IF EXISTS severity_level;
DROP TYPE IF EXISTS threat_type;
DROP TYPE IF EXISTS user_role;

-- 刪除擴展（如果沒有其他依賴的話）
-- DROP EXTENSION IF EXISTS "pgcrypto";
-- DROP EXTENSION IF EXISTS "uuid-ossp"; 