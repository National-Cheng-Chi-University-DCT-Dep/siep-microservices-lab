# Robot Framework 測試套件

## 概述

這是 Ultimate Security Intelligence Platform 的 Robot Framework 自動化測試套件，用於測試後端 API 的各項功能。

## 功能特色

- 🚀 **完整的 API 測試覆蓋**：涵蓋認證、威脅情報、收集器等所有 API 端點
- 🔧 **自動化環境設定**：自動啟動所需的服務和資料庫
- 📊 **詳細的測試報告**：生成 HTML 和 XML 格式的測試報告
- 🏷️ **靈活的測試標籤**：支援按標籤執行特定類型的測試
- 🔄 **並行測試執行**：支援多執行緒並行測試以提升效率
- 🐳 **Docker 整合**：自動管理 Docker 服務的啟停

## 測試套件結構

```
tests/
├── api/                           # API測試套件
│   ├── auth_tests.robot          # 認證API測試
│   ├── threat_intelligence_tests.robot  # 威脅情報API測試
│   └── collector_tests.robot     # 收集器API測試
├── config/                        # 測試配置
│   └── test_config.robot         # 通用測試配置和關鍵字
├── results/                       # 測試結果（自動生成）
├── requirements.txt               # Python依賴
├── run_tests.sh                  # 測試執行腳本
└── README.md                     # 本文檔
```

## 安裝和設定

### 系統需求

- Python 3.8+
- Docker 和 Docker Compose
- Go 1.23+（用於後端服務）
- 網路連接（用於外部 API 測試）

### 快速開始

1. **執行所有測試**：

   ```bash
   cd tests
   chmod +x run_tests.sh
   ./run_tests.sh
   ```

2. **執行特定測試套件**：

   ```bash
   ./run_tests.sh --suite auth          # 只執行認證測試
   ./run_tests.sh --suite threat        # 只執行威脅情報測試
   ./run_tests.sh --suite collector     # 只執行收集器測試
   ```

3. **按標籤執行測試**：
   ```bash
   ./run_tests.sh --tags positive       # 執行正向測試
   ./run_tests.sh --tags negative       # 執行負向測試
   ./run_tests.sh --tags security       # 執行安全測試
   ```

## 詳細使用說明

### 測試執行腳本選項

```bash
./run_tests.sh [選項]

選項:
  -s, --suite SUITE     執行特定測試套件 (auth, threat, collector, all)
  -t, --tags TAGS       執行帶有特定標籤的測試
  -v, --variables FILE  載入變數檔案
  -p, --parallel        並行執行測試
  -V, --verbose         顯示詳細輸出
  -c, --clean           清理之前的測試結果
      --setup-only      只設定環境，不執行測試
  -h, --help            顯示幫助訊息
```

### 常用執行範例

```bash
# 執行所有測試，顯示詳細輸出
./run_tests.sh --verbose

# 並行執行威脅情報的正向測試
./run_tests.sh --suite threat --tags positive --parallel

# 只設定環境，手動執行測試
./run_tests.sh --setup-only

# 清理舊結果並執行新測試
./run_tests.sh --clean --suite auth

# 執行特定標籤組合
./run_tests.sh --tags "create AND positive"
```

### 手動執行 Robot Framework

如果需要更細緻的控制，可以直接使用 Robot Framework：

```bash
# 啟動Python虛擬環境
source venv/bin/activate

# 執行特定測試檔案
robot --outputdir results/manual api/auth_tests.robot

# 執行帶標籤的測試
robot --include positive --outputdir results/manual api/

# 執行特定測試案例
robot --test "Test User Login Success" --outputdir results/manual api/auth_tests.robot
```

## 測試標籤說明

### 功能標籤

- `auth` - 認證相關測試
- `threat` - 威脅情報相關測試
- `collector` - 收集器相關測試

### 測試類型標籤

- `positive` - 正向測試（預期成功的操作）
- `negative` - 負向測試（預期失敗的操作）
- `security` - 安全性測試
- `performance` - 效能測試

### 操作標籤

- `create` - 建立操作測試
- `read` / `get` - 讀取操作測試
- `update` - 更新操作測試
- `delete` - 刪除操作測試
- `search` - 搜尋操作測試
- `batch` - 批量操作測試

### 特殊標籤

- `smoke` - 冒煙測試（基本功能驗證）
- `regression` - 回歸測試
- `integration` - 整合測試

## 測試配置

### 環境變數

測試套件支援通過環境變數覆蓋預設配置：

```bash
# API設定
export TEST_BASE_URL=http://localhost:8080
export TEST_API_VERSION=v1

# 資料庫設定
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_NAME=security_intelligence_test
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=password

# 測試用戶設定
export TEST_USER_USERNAME=testuser
export TEST_USER_EMAIL=testuser@example.com
export TEST_USER_PASSWORD=testpassword123
```

### 自定義變數檔案

創建自定義變數檔案 `custom_vars.py`：

```python
# custom_vars.py
BASE_URL = "https://your-api-server.com"
TEST_USER_USERNAME = "your_test_user"
TEST_USER_PASSWORD = "your_test_password"
```

使用自定義變數：

```bash
./run_tests.sh --variables custom_vars.py
```

## 測試案例說明

### 認證測試 (auth_tests.robot)

- **使用者註冊測試**

  - 成功註冊新使用者
  - 重複使用者名稱/郵箱錯誤處理
  - 無效輸入驗證

- **使用者登入測試**

  - 使用者名稱/郵箱登入
  - 無效認證錯誤處理
  - 令牌生成驗證

- **令牌管理測試**

  - 令牌刷新功能
  - 無效令牌處理
  - 令牌過期處理

- **個人檔案管理測試**
  - 取得使用者資訊
  - 更新個人檔案
  - 密碼修改

### 威脅情報測試 (threat_intelligence_tests.robot)

- **CRUD 操作測試**

  - 建立威脅情報記錄
  - 讀取威脅情報詳情
  - 更新威脅情報資訊
  - 刪除威脅情報記錄

- **查詢和搜尋測試**

  - 列表查詢功能
  - 條件篩選功能
  - 分頁功能測試
  - 全文搜尋功能
  - IP/域名特定搜尋

- **批量操作測試**

  - 批量建立威脅情報
  - 批量更新操作
  - 批量刪除操作

- **統計功能測試**
  - 威脅情報統計資料
  - 趨勢分析資料

### 收集器測試 (collector_tests.robot)

- **單一 IP 收集測試**

  - 惡意 IP 資訊收集
  - 清潔 IP 處理
  - 無效 IP 格式處理

- **批量 IP 收集測試**

  - 多 IP 批量收集
  - 混合結果處理
  - 錯誤處理機制

- **收集器限制測試**

  - 速率限制驗證
  - 資料量限制測試
  - 超時處理測試

- **資料品質測試**
  - 信心分數驗證
  - 威脅分類驗證
  - 元資料完整性檢查

## 測試結果和報告

### 報告檔案

測試執行後會在 `results/latest/` 目錄下生成：

- `report.html` - 主要測試報告（建議檢視）
- `log.html` - 詳細執行日誌
- `output.xml` - 機器可讀的 XML 格式結果
- `summary.txt` - 測試摘要文字檔案

### 持續整合支援

測試套件支援 CI/CD 整合，返回標準退出碼：

- `0` - 所有測試通過
- `1-250` - 失敗測試數量
- `251` - 意外錯誤
- `252` - 無效命令列參數

#### Jenkins 整合範例

```groovy
pipeline {
    agent any
    stages {
        stage('API Tests') {
            steps {
                sh 'cd tests && ./run_tests.sh --parallel --clean'
            }
            post {
                always {
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: 'tests/results/latest',
                        reportFiles: 'report.html',
                        reportName: 'Robot Framework Report'
                    ])
                }
            }
        }
    }
}
```

#### GitHub Actions 整合範例

```yaml
name: API Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run API Tests
        run: |
          cd tests
          chmod +x run_tests.sh
          ./run_tests.sh --parallel --clean
      - name: Upload Test Results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: test-results
          path: tests/results/latest/
```

## 故障排除

### 常見問題

1. **後端服務啟動失敗**

   ```bash
   # 檢查連接埠是否被佔用
   lsof -i :8080

   # 手動啟動後端服務
   cd ../backend
   make run
   ```

2. **資料庫連接錯誤**

   ```bash
   # 檢查Docker服務狀態
   docker-compose -f ../docker/docker-compose.yml ps

   # 重新啟動資料庫
   docker-compose -f ../docker/docker-compose.yml restart postgres
   ```

3. **Python 依賴問題**

   ```bash
   # 重新建立虛擬環境
   rm -rf venv
   python3 -m venv venv
   source venv/bin/activate
   pip install -r requirements.txt
   ```

4. **測試資料衝突**
   ```bash
   # 清理測試資料庫
   cd ../backend
   make migrate-down
   make migrate-up
   ```

### 除錯技巧

1. **使用詳細模式執行**：

   ```bash
   ./run_tests.sh --verbose --suite auth
   ```

2. **只設定環境進行手動測試**：

   ```bash
   ./run_tests.sh --setup-only
   # 在另一個終端機執行特定測試
   source venv/bin/activate
   robot --loglevel DEBUG api/auth_tests.robot
   ```

3. **檢查測試日誌**：

   ```bash
   # 檢視最新測試日誌
   open results/latest/log.html

   # 或使用文字檢視器
   grep -i error results/latest/log.html
   ```

4. **逐步執行測試案例**：
   ```bash
   robot --test "Test User Login Success" --loglevel DEBUG api/auth_tests.robot
   ```

## 貢獻指南

### 添加新測試案例

1. 選擇適當的測試檔案或建立新檔案
2. 使用描述性的測試案例名稱
3. 添加適當的標籤
4. 遵循現有的測試結構和命名規範
5. 包含正向和負向測試案例
6. 添加適當的文檔和註解

### 測試案例範本

```robot
Test Case Name
    [Documentation]    測試案例的詳細說明
    [Tags]    feature_name    test_type    operation

    # 準備測試資料
    ${test_data}=    Prepare Test Data

    # 執行測試操作
    ${response}=    Perform API Operation    ${test_data}

    # 驗證結果
    Verify Response Success    ${response}
    Verify Expected Data    ${response}    ${test_data}

    # 驗證副作用（如資料庫變更）
    Verify Database State    ${expected_state}
```

### 程式碼審查檢查清單

- [ ] 測試案例名稱清晰且具描述性
- [ ] 包含適當的文檔和標籤
- [ ] 測試資料準備和清理正確
- [ ] 錯誤處理和邊界條件測試
- [ ] 回應驗證完整
- [ ] 無硬編碼值（使用變數）
- [ ] 遵循專案編碼規範

## 授權

本測試套件遵循 MIT 授權條款，詳見專案根目錄的 LICENSE 檔案。
