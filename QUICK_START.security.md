# 🚀 Semgrep 安全問題自動修復 - 快速指南

## 立即開始

### 1. 安裝工具
```bash
# 安裝 Semgrep 和修復工具
pip install semgrep pyyaml requests
```

### 2. 執行自動修復
```bash
# 使用 Makefile（推薦）
make -f Makefile.security security-full

# 或直接使用 Python 腳本
python scripts/security_auto_fix.py
```

### 3. 查看結果
```bash
# 查看修復報告
cat security-fix-report.md

# 查看問題摘要
make -f Makefile.security security-summary
```

## 🔧 針對您的 79 個問題

根據您的 Semgrep 掃描結果，這個工具可以自動修復以下問題：

### Nginx 配置問題 (15個)
- ✅ `add_header` 缺少 `always` 標誌
- ✅ H2C 走私漏洞

### Dockerfile 問題 (3個)  
- ✅ 缺少非 root 用戶

### Terraform 問題 (61個)
- ✅ ELB 缺少日誌記錄
- ✅ 不安全的 TLS 版本
- ✅ CloudWatch 日誌未加密
- ✅ KMS 缺少輪換
- ✅ 子網路公共 IP 地址問題

## 📊 預期修復效果

執行後，您應該看到：
- **修復前**: 79 個問題
- **修復後**: 約 16 個問題（剩餘的需要手動處理）
- **自動修復**: 約 63 個問題

## ⚡ 快速命令

```bash
# 查看當前問題
make -f Makefile.security security-summary

# 只執行掃描
make -f Makefile.security security-scan

# 只執行修復
make -f Makefile.security security-fix

# 驗證修復效果
make -f Makefile.security security-verify

# 清理備份檔案
make -f Makefile.security security-clean
```

## 🔄 GitHub Actions 自動化

1. 進入 GitHub 專案頁面
2. 點擊 Actions 標籤
3. 選擇 "安全問題自動修復"
4. 點擊 "Run workflow"
5. 等待自動修復完成並建立 PR

## ⚠️ 重要提醒

1. **備份保護**: 所有檔案都會自動備份
2. **測試環境**: 建議先在測試分支執行
3. **手動審查**: 修復後請審查所有變更
4. **資源配置**: 某些 Terraform 修復可能需要額外 AWS 資源

## 🆘 遇到問題？

1. 檢查 `security-fix-report.md`
2. 查看 `backups/` 目錄中的原始檔案
3. 重新執行掃描驗證效果
4. 提交 Issue 尋求協助

---

**立即開始修復您的 79 個安全問題！** 🚀
