-- Migration: 20241203000000_create_quantum_jobs.down.sql
-- 移除量子任務相關的資料表和類型

-- 移除觸發器
DROP TRIGGER IF EXISTS trigger_log_quantum_job_status_change ON quantum_jobs;
DROP TRIGGER IF EXISTS trigger_update_quantum_job_timestamps ON quantum_jobs;
DROP TRIGGER IF EXISTS trigger_update_quantum_job_updated_at ON quantum_jobs;

-- 移除函數
DROP FUNCTION IF EXISTS log_quantum_job_status_change();
DROP FUNCTION IF EXISTS update_quantum_job_timestamps();
DROP FUNCTION IF EXISTS update_quantum_job_updated_at();
DROP FUNCTION IF EXISTS get_next_pending_quantum_job();

-- 移除日誌表
DROP TABLE IF EXISTS quantum_job_logs;

-- 移除主資料表
DROP TABLE IF EXISTS quantum_jobs;

-- 移除自定義類型
DROP TYPE IF EXISTS quantum_job_status;
