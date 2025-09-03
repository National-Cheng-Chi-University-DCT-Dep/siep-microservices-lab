# 🔒 Semgrep 安全問題自動修復工具

這個工具集可以自動修復 Semgrep 掃描發現的常見安全問題，包括 Nginx 配置、Dockerfile 和 Terraform 相關的安全漏洞。

## 📋 功能特色

- 🔍 **自動掃描**: 使用 Semgrep 自動檢測安全問題
- 🔧 **智能修復**: 自動修復常見的安全配置問題
- 📄 **詳細報告**: 生成完整的修復報告
- 🔄 **備份保護**: 自動備份原始檔案
- ✅ **驗證機制**: 修復後自動驗證效果

## 🚀 快速開始

### 1. 安裝依賴

```bash
# 安裝 Semgrep 和相關工具
make -f Makefile.security security-install

# 或者手動安裝
pip install semgrep pyyaml requests
```

### 2. 執行完整修復流程

```bash
# 執行完整的安全修復流程
make -f Makefile.security security-full
```

### 3. 分步驟執行

```bash
# 1. 執行安全掃描
make -f Makefile.security security-scan

# 2. 查看問題摘要
make -f Makefile.security security-summary

# 3. 自動修復問題
make -f Makefile.security security-fix

# 4. 驗證修復效果
make -f Makefile.security security-verify

# 5. 查看修復報告
make -f Makefile.security security-report
```

## 🔧 支援修復的問題類型

### Nginx 配置問題
- ✅ **缺少 `always` 標誌**: 修復 `add_header` 指令缺少 `always` 標誌
- ✅ **H2C 走私漏洞**: 限制 Upgrade 標頭以防止 HTTP/2 over cleartext 走私

### Dockerfile 安全問題
- ✅ **缺少非 root 用戶**: 自動添加非 root 用戶以提升容器安全性

### Terraform 安全問題
- ✅ **ELB 缺少日誌**: 啟用 Application Load Balancer 存取日誌
- ✅ **不安全的 TLS 版本**: 升級到安全的 TLS 1.3 版本
- ✅ **CloudWatch 日誌未加密**: 添加 KMS 加密配置
- ✅ **KMS 缺少輪換**: 啟用金鑰自動輪換
- ✅ **子網路公共 IP 地址**: 修復子網路自動分配公共 IP 的問題

## 📁 檔案結構

```
.
├── .github/workflows/
│   └── security-auto-fix.yml    # GitHub Actions 自動修復工作流程
├── scripts/
│   └── security_auto_fix.py    # Python 自動修復腳本
├── Makefile.security            # Makefile 命令集
├── semgrep-results.json         # Semgrep 掃描結果
├── security-fix-report.md       # 修復報告
└── backups/                     # 備份檔案目錄
```

## 🔄 GitHub Actions 自動化

### 手動觸發
在 GitHub 專案頁面，進入 Actions 標籤，選擇 "安全問題自動修復" 工作流程，點擊 "Run workflow" 按鈕。

### 自動排程
工作流程會每週一凌晨 2 點自動執行，並在發現問題時建立 Pull Request。

### 輸入參數
- `fix_nginx`: 是否修復 Nginx 配置問題 (預設: true)
- `fix_dockerfile`: 是否修復 Dockerfile 問題 (預設: true)
- `fix_terraform`: 是否修復 Terraform 問題 (預設: true)
- `create_pr`: 是否建立 Pull Request (預設: true)

## 📊 使用範例

### 查看當前安全問題
```bash
make -f Makefile.security security-summary
```

輸出範例：
```
📊 安全問題摘要：
總問題數量: 79
Nginx 問題: 15
Dockerfile 問題: 3
Terraform 問題: 61
```

### 執行自動修復
```bash
make -f Makefile.security security-fix
```

輸出範例：
```
🔍 開始分析 Semgrep 結果...
發現 79 個問題
✅ 修復了 15 個 Nginx 配置問題
✅ 修復了 3 個 Dockerfile 問題
✅ 修復了 45 個 Terraform 問題

🎉 總共修復了 63 個安全問題
📄 修復報告已生成: security-fix-report.md
```

### 驗證修復效果
```bash
make -f Makefile.security security-verify
```

輸出範例：
```
🔍 驗證修復效果...
修復前問題數量: 79
修復後問題數量: 16
✅ 成功修復了 63 個安全問題
```

## ⚠️ 注意事項

### 重要提醒
1. **備份保護**: 所有原始檔案都會自動備份到 `backups/` 目錄
2. **測試環境**: 建議先在測試環境中驗證修復效果
3. **手動審查**: 自動修復後請手動審查所有變更
4. **資源配置**: 某些 Terraform 修復可能需要額外的 AWS 資源

### 限制
- 只能修復常見的、可自動化的安全問題
- 複雜的業務邏輯問題需要手動處理
- 某些修復可能需要額外的配置或資源

## 🔧 進階使用

### 自定義修復規則
可以修改 `scripts/security_auto_fix.py` 來添加自定義的修復規則：

```python
def fix_custom_issues(self) -> int:
    """修復自定義問題"""
    fixes = 0
    # 添加自定義修復邏輯
    return fixes
```

### 整合到 CI/CD
在 `.github/workflows/ci.yml` 中添加安全檢查：

```yaml
- name: 安全檢查
  run: |
    make -f Makefile.security security-scan
    make -f Makefile.security security-summary
```

### 定期執行
設定 cron 任務定期執行安全檢查：

```bash
# 每週執行一次安全檢查
0 2 * * 1 cd /path/to/project && make -f Makefile.security security-full
```

## 📞 支援

如果遇到問題或有建議，請：

1. 檢查 `security-fix-report.md` 中的詳細報告
2. 查看 `backups/` 目錄中的原始檔案
3. 重新執行掃描驗證修復效果
4. 提交 Issue 或 Pull Request

## 📄 授權

本工具遵循 MIT 授權條款。
