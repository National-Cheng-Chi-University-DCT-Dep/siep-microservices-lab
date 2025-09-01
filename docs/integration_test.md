# 整合測試指南

## 前置準備

### 1. 資料庫設定

確保 PostgreSQL 已啟動並創建資料庫：

```bash
# 連線到 PostgreSQL
psql -U postgres

# 建立資料庫
CREATE DATABASE security_intel;

# 連線到資料庫
\c security_intel;
```

### 2. 執行 Migration

```bash
cd backend
# 如果有 golang-migrate 工具
migrate -path database/migrations -database "postgres://postgres:password@localhost:5432/security_intel?sslmode=disable" up

# 或手動執行 SQL 檔案
psql -U postgres -d security_intel -f database/migrations/20241201120000_initial_setup.up.sql
```

### 3. 插入測試資料

```bash
psql -U postgres -d security_intel -f scripts/test_data.sql
```

## 後端測試

### 1. 啟動後端服務

```bash
cd backend

# 複製環境變數範例
cp test.env .env

# 啟動服務
go run cmd/server/main.go
```

### 2. 測試 API Endpoints

#### 健康檢查

```bash
curl http://localhost:8080/api/v1/health
```

#### 取得威脅情報列表

```bash
curl http://localhost:8080/api/v1/threats
```

#### 搜尋威脅情報

```bash
curl "http://localhost:8080/api/v1/threats?threat_type=malware&severity=high"
```

#### 建立威脅情報

```bash
curl -X POST http://localhost:8080/api/v1/threats \
  -H "Content-Type: application/json" \
  -d '{
    "ip_address": "1.2.3.4",
    "threat_type": "malware",
    "severity": "high",
    "confidence_score": 80,
    "source": "Manual Test",
    "description": "手動測試建立的威脅情報"
  }'
```

#### IP 查詢

```bash
curl "http://localhost:8080/api/v1/threats/lookup/ip?ip_address=192.168.1.100"
```

#### 統計資訊

```bash
curl http://localhost:8080/api/v1/threats/stats
```

## 前端測試

### 1. 安裝依賴並啟動

```bash
cd frontend
npm install
npm run dev
```

### 2. 瀏覽器測試

開啟 http://localhost:3000

#### 功能檢查清單

- [ ] 頁面正常載入，顯示標題和導航
- [ ] 統計卡片顯示正確數據（總數、高危、關鍵等）
- [ ] 威脅情報列表顯示測試資料
- [ ] 搜尋功能正常運作
  - [ ] IP 地址搜尋
  - [ ] 威脅類型篩選
  - [ ] 嚴重程度篩選
- [ ] 威脅情報卡片顯示詳細資訊
- [ ] 收集器功能（如果有 AbuseIPDB API Key）

## 收集器測試

### 1. 設定 AbuseIPDB API Key

```bash
# 在 .env 檔案中設定
ABUSEIPDB_API_KEY=your_actual_api_key
```

### 2. 測試收集功能

#### 單一 IP 收集

```bash
curl -X POST http://localhost:8080/api/v1/collector/ip \
  -H "Content-Type: application/json" \
  -d '{"ip_address": "8.8.8.8"}'
```

#### 批量 IP 收集

```bash
curl -X POST http://localhost:8080/api/v1/collector/bulk-ip \
  -H "Content-Type: application/json" \
  -d '{"ip_addresses": ["8.8.8.8", "1.1.1.1"]}'
```

## 常見問題排除

### 1. 資料庫連線失敗

檢查：

- PostgreSQL 服務是否啟動
- 資料庫名稱、使用者名稱、密碼是否正確
- 防火牆設定

### 2. CORS 錯誤

確認 `.env` 檔案中的 `CORS_ALLOWED_ORIGINS` 包含前端地址。

### 3. API 代理失敗

檢查 `frontend/next.config.ts` 中的 proxy 設定是否正確。

### 4. 前端組件錯誤

檢查：

- 所有組件檔案都已建立
- TypeScript 編譯無錯誤
- API 回應格式與前端期望一致

## 測試結果驗證

### 成功標準

1. **後端服務正常啟動**，所有 API endpoints 回應正確
2. **前端應用正常載入**，所有組件顯示正確
3. **前後端通信正常**，API 呼叫成功
4. **資料庫操作正常**，CRUD 功能運作
5. **搜尋功能正常**，篩選條件生效
6. **收集器功能正常**（如果有 API Key）

### 測試報告

記錄以下資訊：

- 測試環境（作業系統、瀏覽器版本等）
- 測試時間
- 測試結果（通過/失敗）
- 遇到的問題和解決方案
- 效能觀察（載入時間、回應時間等）

## 後續步驟

整合測試成功後，可以進行：

1. 部署到開發環境
2. 加入更多測試資料
3. 進行效能測試
4. 開始 Phase 2 開發
