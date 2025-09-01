# 貢獻指南 (Contributing Guide)

我們非常歡迎您為「資安情報平台」做出貢獻！本指南將協助您了解如何參與專案開發，並確保我們能夠維持高品質的程式碼標準。

## 📋 目錄

- [開始之前](#開始之前)
- [嚴謹標準流程](#嚴謹標準流程)
- [開發環境設定](#開發環境設定)
- [提交程式碼](#提交程式碼)
- [Pull Request 流程](#pull-request-流程)
- [程式碼規範](#程式碼規範)
- [測試要求](#測試要求)

## 開始之前

### 必備條件

- 閱讀並同意遵守我們的 [行為準則](CODE_OF_CONDUCT.md)
- 熟悉 Git 和 GitHub 的基本操作
- 具備 Go、TypeScript、PostgreSQL 的基礎知識

### 設定開發環境

1. **Fork 此專案**到您的 GitHub 帳戶
2. **Clone 您 fork 的儲存庫**：

   ```bash
   git clone https://github.com/YOUR-USERNAME/security-intel-platform.git
   cd security-intel-platform
   ```

3. **設定 upstream remote**：
   ```bash
   git remote add upstream https://github.com/ORIGINAL-OWNER/security-intel-platform.git
   ```

## 嚴謹標準流程

我們採用嚴格的開發流程，確保程式碼品質和系統穩定性。每個貢獻者都必須遵循以下步驟：

### 1. ERD 更新

> 🎯 **目標**：確保資料庫結構設計合理且一致

**何時需要**：當您的變更涉及資料庫結構時

**步驟**：

1. 使用 [dbdiagram.io](https://dbdiagram.io/) 或 [draw.io](https://draw.io/) 更新 ERD
2. 將 ERD 檔案放置在 `docs/erd/` 目錄下
3. 在 PR 描述中說明資料庫結構變更的原因

### 2. 撰寫/更新 Migration 檔案

> 🎯 **目標**：確保資料庫變更可追蹤且可重複執行

**重要規則**：

- 📛 **絕不可修改現有的 migration 檔案**
- ✅ **只能在 `database/migrations/` 目錄新增新的 migration 檔案**
- 📄 **migration 檔案必須經過 code review**

**檔案命名格式**：

```
YYYYMMDDHHMMSS_descriptive_name.sql
```

**範例**：

```sql
-- 20241201120000_add_threat_intelligence_table.sql
CREATE TABLE threat_intelligence (
    id SERIAL PRIMARY KEY,
    ip_address INET NOT NULL,
    threat_type VARCHAR(50) NOT NULL,
    severity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3. 更新 GORM Model

> 🎯 **目標**：保持 Go 結構與資料庫 schema 同步

**規則**：

- 📍 **只在 `internal/model/` 目錄維護 GORM struct**
- 🚫 **不可混用 json tag 或 DTO 屬性**
- ✅ **Model 必須與資料庫 schema 完全對應**

**範例**：

```go
// internal/model/threat_intelligence.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type ThreatIntelligence struct {
    ID          uint           `gorm:"primaryKey"`
    IPAddress   string         `gorm:"type:inet;not null"`
    ThreatType  string         `gorm:"size:50;not null"`
    Severity    int            `gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

### 4. 同步 DTO/VO 結構

> 🎯 **目標**：分離 API 介面與資料模型

**目錄結構**：

- `internal/dto/`：API 輸入結構 (Request)
- `internal/vo/`：API 輸出結構 (Response)

**重要規則**：

- 🚫 **嚴禁 handler 或 service 直接使用 model struct**
- 📝 **必須包含 binding/json tag**
- 🔄 **API 版本變動時，建立新 VO 結構，保留舊版相容**

**範例**：

```go
// internal/dto/threat_intelligence.go
package dto

type CreateThreatIntelligenceRequest struct {
    IPAddress  string `json:"ip_address" binding:"required,ip"`
    ThreatType string `json:"threat_type" binding:"required,max=50"`
    Severity   int    `json:"severity" binding:"required,min=1,max=10"`
}

// internal/vo/threat_intelligence.go
package vo

type ThreatIntelligenceResponse struct {
    ID         uint   `json:"id"`
    IPAddress  string `json:"ip_address"`
    ThreatType string `json:"threat_type"`
    Severity   int    `json:"severity"`
    CreatedAt  string `json:"created_at"`
}
```

### 5. Handler 層資料轉換

> 🎯 **目標**：確保資料轉換的安全性和一致性

**規則**：

- 📦 **使用 `copier` 或 `mapstructure` 進行自動 mapping**
- 🔄 **Handler/service 層只處理 DTO/VO**
- 🚫 **不可直接操作 model**

**範例**：

```go
// internal/handler/threat_intelligence.go
package handler

import (
    "github.com/jinzhu/copier"
    "your-project/internal/dto"
    "your-project/internal/vo"
    "your-project/internal/service"
)

func (h *ThreatIntelligenceHandler) CreateThreatIntelligence(c *gin.Context) {
    var req dto.CreateThreatIntelligenceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 使用 service 層處理業務邏輯
    result, err := h.service.CreateThreatIntelligence(req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // 轉換為 VO
    var response vo.ThreatIntelligenceResponse
    copier.Copy(&response, result)

    c.JSON(201, response)
}
```

### 6. Swagger 註解與自動生成

> 🎯 **目標**：確保 API 文件與實作同步

**規則**：

- 📝 **只在 VO 結構加上 Swagger 註解**
- 🔄 **執行 `make swagger` 生成 swagger.json**
- ✅ **確保 API 文件完整且正確**

**範例**：

```go
// internal/vo/threat_intelligence.go
package vo

// ThreatIntelligenceResponse 威脅情報回應
type ThreatIntelligenceResponse struct {
    ID         uint   `json:"id" example:"1"`                                    // 威脅情報 ID
    IPAddress  string `json:"ip_address" example:"192.168.1.1"`                  // IP 位址
    ThreatType string `json:"threat_type" example:"malware"`                     // 威脅類型
    Severity   int    `json:"severity" example:"8"`                              // 嚴重程度 (1-10)
    CreatedAt  string `json:"created_at" example:"2024-01-01T12:00:00Z"`        // 建立時間
}
```

## 提交程式碼

### Commit Message 格式

我們使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type 類型**：

- `feat`: 新功能
- `fix`: 錯誤修正
- `docs`: 文件更新
- `style`: 程式碼格式變更
- `refactor`: 重構
- `test`: 測試相關
- `chore`: 維護任務

**範例**：

```
feat(api): add threat intelligence CRUD endpoints

- 新增威脅情報的建立、讀取、更新、刪除 API
- 實作對應的 DTO 和 VO 結構
- 加入 Swagger 註解

Closes #123
```

### Branch 命名規則

```
<type>/<issue-number>-<short-description>
```

**範例**：

- `feat/123-add-threat-intel-api`
- `fix/456-database-connection-timeout`
- `docs/789-update-contributing-guide`

## Pull Request 流程

### 1. 建立 Pull Request

- 📝 **使用 PR 模板**（系統會自動載入）
- 📄 **提供清晰的變更描述**
- 🔗 **連結相關的 Issue**
- 📋 **勾選 checklist 確認所有步驟完成**

### 2. CI 檢查

您的 PR 必須通過以下檢查：

- ✅ **GitHub Actions** - 程式碼品質檢查
- ✅ **所有測試** - 單元測試和整合測試
- ✅ **Swagger 生成** - API 文件同步
- ✅ **Migration 驗證** - 資料庫變更檢查

### 3. Code Review

- 👥 **至少需要一位維護者審查**
- 🔄 **根據回饋進行修改**
- ✅ **所有對話標記為 resolved**

### 4. 合併

- 🎯 **使用 "Squash and merge"**
- 🗑️ **刪除 feature branch**

## 程式碼規範

### Go 程式碼

- 📋 **遵循 `gofmt` 格式**
- 🔍 **通過 `golangci-lint` 檢查**
- 📝 **為公開函數和結構添加註解**
- 🧪 **為新功能編寫測試**

### TypeScript 程式碼

- 📋 **遵循 ESLint 規則**
- 📝 **使用 TypeScript 嚴格模式**
- 🎨 **使用 Prettier 進行格式化**
- 🧪 **為元件編寫測試**

## 測試要求

### 後端測試

```bash
# 運行所有測試
make test

# 運行特定套件測試
go test ./internal/service/...

# 生成覆蓋率報告
make test-coverage
```

### 前端測試

```bash
# 運行所有測試
npm test

# 運行特定測試
npm test -- --testNamePattern="ThreatIntelligence"

# 生成覆蓋率報告
npm run test:coverage
```

### 測試覆蓋率要求

- 📊 **整體覆蓋率 ≥ 80%**
- 🎯 **新功能覆蓋率 ≥ 90%**
- 🚫 **不可降低現有覆蓋率**

## 問題回報

如果您在開發過程中遇到問題：

1. 📋 **檢查現有的 Issues**
2. 🆕 **建立新的 Issue**（如果不存在）
3. 🏷️ **使用適當的標籤**
4. 📧 **可以聯繫維護者**

## 技術支援

- 📚 **文件**: [docs/](docs/)
- 💬 **Discussions**: 使用 GitHub Discussions
- 📧 **電子郵件**: [INSERT EMAIL]

感謝您的貢獻！🎉
