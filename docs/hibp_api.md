# Have I Been Pwned (HIBP) API 整合文檔

## 概述

本平台完整整合了 Have I Been Pwned API v3，提供全面的數據泄露監控功能。所有端點都需要 JWT 認證（除了密碼檢查端點）。

## 基礎 URL

```
https://your-domain.com/api/v1/hibp
```

## 認證

除了密碼檢查端點外，所有端點都需要在請求頭中包含 JWT 令牌：

```
Authorization: Bearer YOUR_JWT_TOKEN
```

## 端點列表

### 1. 帳戶泄露檢查

#### 檢查帳戶泄露事件

```http
GET /api/v1/hibp/account/{account}/breaches
```

**參數：**

- `account` (path): 電子郵件地址或用戶名
- `include_unverified` (query): 是否包含未驗證的泄露事件 (true/false)

**響應範例：**

```json
{
  "success": true,
  "message": "Account breaches retrieved successfully",
  "data": {
    "account": "user@example.com",
    "breaches": [
      {
        "Name": "Adobe",
        "Title": "Adobe",
        "Domain": "adobe.com",
        "BreachDate": "2013-10-04",
        "PwnCount": 152445165,
        "Description": "In October 2013, 153 million Adobe accounts were breached...",
        "DataClasses": ["Email addresses", "Passwords", "Usernames"],
        "IsVerified": true,
        "IsSensitive": false
      }
    ],
    "count": 1
  }
}
```

#### 檢查帳戶 Paste 記錄

```http
GET /api/v1/hibp/account/{account}/pastes
```

**響應範例：**

```json
{
  "success": true,
  "message": "Account pastes retrieved successfully",
  "data": {
    "account": "user@example.com",
    "pastes": [
      {
        "Source": "Pastebin",
        "Id": "abc123",
        "Title": "Data dump",
        "Date": "2023-01-01T00:00:00Z",
        "EmailCount": 1000
      }
    ],
    "count": 1
  }
}
```

#### 處理帳戶泄露並創建威脅情報

```http
POST /api/v1/hibp/account/{account}/process
```

此端點會檢查帳戶泄露並自動創建威脅情報記錄。

### 2. 域名泄露監控

#### 檢查域名泄露

```http
GET /api/v1/hibp/domain/{domain}/breaches
```

**響應範例：**

```json
{
  "success": true,
  "message": "Domain breaches retrieved successfully",
  "data": {
    "domain": "example.com",
    "breaches": {
      "user1@example.com": ["Adobe", "LinkedIn"],
      "user2@example.com": ["Adobe"]
    },
    "count": 2
  }
}
```

#### 獲取已訂閱域名

```http
GET /api/v1/hibp/domains/subscribed
```

### 3. 泄露事件查詢

#### 獲取所有泄露事件

```http
GET /api/v1/hibp/breaches?domain={domain}
```

**查詢參數：**

- `domain` (optional): 按域名過濾

#### 獲取特定泄露事件

```http
GET /api/v1/hibp/breach/{name}
```

#### 獲取最新泄露事件

```http
GET /api/v1/hibp/breach/latest
```

#### 獲取數據類別

```http
GET /api/v1/hibp/dataclasses
```

### 4. 密碼安全檢查

#### 檢查密碼是否已被泄露

```http
GET /api/v1/hibp/password/check?password={password}
```

**注意：** 此端點無需認證

**響應範例：**

```json
{
  "success": true,
  "message": "Password hash check completed",
  "data": {
    "pwned_count": 12345,
    "is_pwned": true
  }
}
```

### 5. Stealer Logs (需要高級訂閱)

#### 按郵箱查詢竊取器日誌

```http
GET /api/v1/hibp/stealer/email/{email}
```

#### 按網站域名查詢竊取器日誌

```http
GET /api/v1/hibp/stealer/website/{domain}
```

#### 按郵箱域名查詢竊取器日誌

```http
GET /api/v1/hibp/stealer/emaildomain/{domain}
```

### 6. 系統狀態

#### 獲取訂閱狀態

```http
GET /api/v1/hibp/subscription/status
```

**響應範例：**

```json
{
  "success": true,
  "message": "Subscription status retrieved successfully",
  "data": {
    "Name": "Pwned 1",
    "Description": "Basic subscription",
    "Price": 100,
    "PwnCount": 1000000,
    "BreachCount": 100,
    "PasteCount": 50,
    "DomainSearchEnabled": true,
    "DomainSearchCapacity": 1000,
    "DomainSearchUsed": 500,
    "StealerLogEnabled": false
  }
}
```

#### 獲取 HIBP 統計信息

```http
GET /api/v1/hibp/stats
```

#### 健康檢查

```http
GET /api/v1/hibp/health
```

## 錯誤處理

### 常見錯誤響應

```json
{
  "success": false,
  "message": "Error message",
  "error": {
    "code": "ERROR_CODE",
    "message": "Detailed error message",
    "details": "Additional error details"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 錯誤代碼

| 代碼                    | 描述              |
| ----------------------- | ----------------- |
| `INVALID_ACCOUNT`       | 無效的帳戶參數    |
| `INVALID_DOMAIN`        | 無效的域名參數    |
| `INVALID_PASSWORD`      | 無效的密碼參數    |
| `HIBP_API_ERROR`        | HIBP API 調用失敗 |
| `HIBP_PROCESSING_ERROR` | HIBP 數據處理失敗 |
| `HIBP_UNHEALTHY`        | HIBP 服務不健康   |
| `UNAUTHORIZED`          | 未授權訪問        |
| `RATE_LIMIT_EXCEEDED`   | 超過速率限制      |

## 速率限制

HIBP API 有速率限制，具體限制取決於您的訂閱等級：

- **Pwned 1**: 每分鐘 10 個請求
- **Pwned 2**: 每分鐘 50 個請求
- **Pwned 3**: 每分鐘 100 個請求
- **Pwned 4**: 每分鐘 200 個請求

當超過速率限制時，API 會返回 429 狀態碼和 `retry-after` 頭。

## 使用範例

### cURL 範例

```bash
# 檢查帳戶泄露
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "https://your-domain.com/api/v1/hibp/account/user@example.com/breaches"

# 檢查密碼安全性
curl "https://your-domain.com/api/v1/hibp/password/check?password=yourpassword"

# 檢查域名泄露
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "https://your-domain.com/api/v1/hibp/domain/example.com/breaches"
```

### JavaScript 範例

```javascript
// 檢查帳戶泄露
const checkAccountBreaches = async (email) => {
  const response = await fetch(
    `/api/v1/hibp/account/${encodeURIComponent(email)}/breaches`,
    {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("access_token")}`,
      },
    }
  );

  const data = await response.json();
  return data;
};

// 檢查密碼安全性
const checkPassword = async (password) => {
  const response = await fetch(
    `/api/v1/hibp/password/check?password=${encodeURIComponent(password)}`
  );
  const data = await response.json();
  return data;
};
```

## 配置

### 環境變數

在 `.env` 文件中設置 HIBP API 金鑰：

```bash
HIBP_API_KEY=your_hibp_api_key_here
```

### 獲取 API 金鑰

1. 訪問 [Have I Been Pwned API 頁面](https://haveibeenpwned.com/API/Key)
2. 選擇適合的訂閱等級
3. 獲取 API 金鑰
4. 將金鑰添加到環境變數中

## 安全注意事項

1. **密碼檢查**：密碼檢查使用 k-anonymity 技術，不會將完整密碼發送到服務器
2. **API 金鑰保護**：請妥善保護您的 HIBP API 金鑰
3. **速率限制**：請遵守 HIBP 的速率限制政策
4. **數據隱私**：所有查詢都會記錄在 HIBP 的日誌中

## 支援

如有問題，請參考：

- [Have I Been Pwned API 文檔](https://haveibeenpwned.com/API/v3)
- [HIBP API 使用條款](https://haveibeenpwned.com/API/v3#AcceptableUse)
- [HIBP 常見問題](https://haveibeenpwned.com/FAQs)
