-- Migration: 20241202000000_add_supabase_auth.down.sql
-- 回滾 Supabase Auth 整合

-- 刪除 auth_settings 表
DROP TABLE IF EXISTS auth_settings;

-- 移除 users 表中的 Supabase 相關欄位
ALTER TABLE users
DROP COLUMN IF EXISTS supabase_id,
DROP COLUMN IF EXISTS auth_provider,
DROP COLUMN IF EXISTS avatar_url,
DROP COLUMN IF EXISTS full_name;

-- 移除索引
DROP INDEX IF EXISTS idx_users_supabase_id;

-- 更新 users 表的註解
COMMENT ON TABLE users IS '使用者資料表';
