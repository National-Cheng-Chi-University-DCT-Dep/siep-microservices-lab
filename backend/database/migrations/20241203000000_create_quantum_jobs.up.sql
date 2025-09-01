-- Migration: 20241203000000_create_quantum_jobs.up.sql
-- 創建量子任務資料表，用於儲存量子運算任務資訊

-- 首先創建一個新的 ENUM 類型用於任務狀態
CREATE TYPE quantum_job_status AS ENUM (
    'pending',    -- 等待執行
    'running',    -- 正在執行中
    'completed',  -- 已完成
    'failed'      -- 執行失敗
);

-- 創建量子任務資料表
CREATE TABLE quantum_jobs (
    -- 基本資訊欄位
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- 任務狀態資訊
    status quantum_job_status NOT NULL DEFAULT 'pending',
    priority INT NOT NULL DEFAULT 1,    -- 優先順序: 1 (最低) 到 10 (最高)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    
    -- 任務參數與結果
    input_params JSONB NOT NULL,        -- 輸入參數
    results JSONB,                      -- 任務結果
    error_message TEXT,                 -- 錯誤訊息（如果有）
    
    -- 執行資訊
    execution_time_seconds INT,         -- 執行時間（秒）
    quantum_backend VARCHAR(100),       -- 使用的量子後端
    is_simulation BOOLEAN DEFAULT true, -- 是否是模擬器
    shots INT DEFAULT 1024,             -- 量子實驗的樣本數
    
    -- 統計與分析資訊
    confidence_score DECIMAL(5,2),      -- 信心分數 (0-100)
    is_malicious BOOLEAN,               -- 是否是惡意威脅
    
    -- 額外資訊
    tags TEXT[],                        -- 標籤
    notes TEXT,                         -- 附加筆記
    source VARCHAR(50)                  -- 來源標識 (如 "api", "scheduled", "manual")
);

-- 創建索引以提高查詢效能
CREATE INDEX idx_quantum_jobs_user_id ON quantum_jobs(user_id);
CREATE INDEX idx_quantum_jobs_status ON quantum_jobs(status);
CREATE INDEX idx_quantum_jobs_created_at ON quantum_jobs(created_at);
CREATE INDEX idx_quantum_jobs_is_malicious ON quantum_jobs(is_malicious);

-- 創建觸發器函數，用於自動更新 updated_at 欄位
CREATE OR REPLACE FUNCTION update_quantum_job_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 創建觸發器
CREATE TRIGGER trigger_update_quantum_job_updated_at
BEFORE UPDATE ON quantum_jobs
FOR EACH ROW
EXECUTE FUNCTION update_quantum_job_updated_at();

-- 創建觸發器，在任務狀態變更時自動更新相關時間戳
CREATE OR REPLACE FUNCTION update_quantum_job_timestamps()
RETURNS TRIGGER AS $$
BEGIN
    -- 當狀態變更為 'running' 時更新 started_at
    IF NEW.status = 'running' AND OLD.status != 'running' THEN
        NEW.started_at = CURRENT_TIMESTAMP;
    END IF;
    
    -- 當狀態變更為 'completed' 或 'failed' 時更新 completed_at
    IF (NEW.status = 'completed' OR NEW.status = 'failed') AND 
       (OLD.status != 'completed' AND OLD.status != 'failed') THEN
        NEW.completed_at = CURRENT_TIMESTAMP;
        
        -- 計算執行時間（如果有 started_at）
        IF NEW.started_at IS NOT NULL THEN
            NEW.execution_time_seconds = EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - NEW.started_at))::INT;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 創建觸發器
CREATE TRIGGER trigger_update_quantum_job_timestamps
BEFORE UPDATE ON quantum_jobs
FOR EACH ROW
WHEN (OLD.status IS DISTINCT FROM NEW.status)
EXECUTE FUNCTION update_quantum_job_timestamps();

-- 創建量子任務日誌表，用於記錄任務狀態變更
CREATE TABLE quantum_job_logs (
    id SERIAL PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES quantum_jobs(id) ON DELETE CASCADE,
    old_status quantum_job_status,
    new_status quantum_job_status NOT NULL,
    message TEXT,
    created_by VARCHAR(50) NOT NULL, -- 誰或什麼系統創建了此日誌
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 創建索引
CREATE INDEX idx_quantum_job_logs_job_id ON quantum_job_logs(job_id);

-- 創建觸發器函數，用於自動記錄任務狀態變更
CREATE OR REPLACE FUNCTION log_quantum_job_status_change()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO quantum_job_logs (job_id, old_status, new_status, message, created_by)
    VALUES (
        NEW.id,
        OLD.status,
        NEW.status,
        CASE 
            WHEN NEW.status = 'running' THEN '任務開始執行'
            WHEN NEW.status = 'completed' THEN '任務已完成'
            WHEN NEW.status = 'failed' THEN '任務執行失敗: ' || COALESCE(NEW.error_message, '未知錯誤')
            ELSE '任務狀態已更新'
        END,
        CURRENT_USER
    );
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 創建觸發器
CREATE TRIGGER trigger_log_quantum_job_status_change
AFTER UPDATE ON quantum_jobs
FOR EACH ROW
WHEN (OLD.status IS DISTINCT FROM NEW.status)
EXECUTE FUNCTION log_quantum_job_status_change();

-- 創建函數，用於獲取下一個待處理的任務
CREATE OR REPLACE FUNCTION get_next_pending_quantum_job()
RETURNS TABLE (
    job_id UUID,
    job_params JSONB
) AS $$
BEGIN
    RETURN QUERY
    WITH next_job AS (
        SELECT id, input_params
        FROM quantum_jobs
        WHERE status = 'pending'
        ORDER BY priority DESC, created_at ASC
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    )
    UPDATE quantum_jobs qj
    SET status = 'running',
        started_at = CURRENT_TIMESTAMP
    FROM next_job
    WHERE qj.id = next_job.id
    RETURNING qj.id, qj.input_params;
END;
$$ LANGUAGE plpgsql;

-- 添加注釋
COMMENT ON TABLE quantum_jobs IS '量子任務資料表，用於儲存量子運算任務資訊';
COMMENT ON TABLE quantum_job_logs IS '量子任務日誌表，用於記錄任務狀態變更';
COMMENT ON FUNCTION get_next_pending_quantum_job() IS '獲取下一個待處理的量子任務，並更新其狀態為執行中';
COMMENT ON COLUMN quantum_jobs.input_params IS 'JSON 格式的輸入參數，包含威脅數據和分析設定';
COMMENT ON COLUMN quantum_jobs.results IS 'JSON 格式的任務結果，包含預測結果和分析數據';
