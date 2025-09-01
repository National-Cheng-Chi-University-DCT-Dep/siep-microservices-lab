-- Migration: 20241202000000_add_supabase_auth.up.sql
-- 將用戶系統與 Supabase Auth 整合

-- 新增 supabase_id 欄位到 users 表
ALTER TABLE users 
ADD COLUMN supabase_id UUID;

-- 新增索引以提高查詢效能
CREATE INDEX idx_users_supabase_id ON users(supabase_id);

-- 新增 auth_provider 欄位，儲存認證提供者來源
ALTER TABLE users
ADD COLUMN auth_provider VARCHAR(20) DEFAULT 'email';

-- 修改 email_verified 欄位的預設值，因為 Supabase 會處理郵件驗證
ALTER TABLE users
ALTER COLUMN email_verified SET DEFAULT false;

-- 新增用於 OAuth 的欄位
ALTER TABLE users
ADD COLUMN avatar_url TEXT,
ADD COLUMN full_name VARCHAR(100);

-- 新增 Supabase Auth 相關設定表
CREATE TABLE IF NOT EXISTS auth_settings (
    id SERIAL PRIMARY KEY,
    provider_name VARCHAR(50) NOT NULL,
    is_enabled BOOLEAN DEFAULT false,
    client_id VARCHAR(255),
    client_secret VARCHAR(255),
    redirect_uri TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 新增預設的認證提供者
INSERT INTO auth_settings (provider_name, is_enabled) VALUES
    ('google', false),
    ('github', false),
    ('facebook', false),
    ('twitter', false),
    ('email', true);

-- 將現有用戶預設為使用 Email 認證
COMMENT ON COLUMN users.auth_provider IS '認證提供者: email, google, github, facebook, twitter 等';

-- 更新 users 表的註解
COMMENT ON TABLE users IS '使用者資料表，整合 Supabase Auth';
