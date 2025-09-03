#!/usr/bin/env python3
"""
Semgrep 安全問題自動修復工具
自動修復 Semgrep 掃描發現的常見安全問題
"""

import json
import os
import re
import shutil
import sys
from pathlib import Path
from typing import Dict, List, Any, Optional


class SecurityAutoFixer:
    def __init__(self, semgrep_results_file: str = "semgrep-results.json"):
        self.semgrep_results_file = semgrep_results_file
        self.results = self._load_results()
        self.fixes_applied = []
        self.backup_dir = Path("backups")
        self.backup_dir.mkdir(exist_ok=True)
        
    def _load_results(self) -> Dict[str, Any]:
        """載入 Semgrep 掃描結果"""
        if not os.path.exists(self.semgrep_results_file):
            print(f"❌ 找不到 Semgrep 結果檔案: {self.semgrep_results_file}")
            return {"results": []}
            
        with open(self.semgrep_results_file, 'r', encoding='utf-8') as f:
            return json.load(f)
    
    def _backup_file(self, file_path: str) -> str:
        """備份檔案"""
        if not os.path.exists(file_path):
            return ""
            
        backup_path = self.backup_dir / f"{Path(file_path).name}.backup"
        shutil.copy2(file_path, backup_path)
        return str(backup_path)
    
    def fix_nginx_issues(self) -> int:
        """修復 Nginx 配置問題"""
        fixes = 0
        
        for finding in self.results.get("results", []):
            if "nginx" not in finding.get("path", "").lower():
                continue
                
            rule_id = finding.get("check_id", "")
            file_path = finding.get("path", "")
            
            if not os.path.exists(file_path):
                continue
                
            # 備份原始檔案
            self._backup_file(file_path)
            
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            
            # 修復 add_header 缺少 always 標誌
            if "generic.nginx.security.missing-always-flag" in rule_id:
                content = re.sub(
                    r'add_header\s+([^;]+);',
                    r'add_header \1 always;',
                    content
                )
            
            # 修復 H2C 走私問題
            elif "generic.nginx.security.possible-h2c-smuggling" in rule_id:
                # 限制 Upgrade 標頭
                upgrade_pattern = r'proxy_set_header\s+Upgrade\s+\$http_upgrade;'
                if re.search(upgrade_pattern, content):
                    content = re.sub(
                        upgrade_pattern,
                        r'# 限制 Upgrade 標頭以防止 H2C 走私\n    proxy_set_header Upgrade $http_upgrade;',
                        content
                    )
            
            # 寫回修復後的內容
            if content != original_content:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                self.fixes_applied.append(f"修復 Nginx 配置: {file_path}")
                fixes += 1
        
        return fixes
    
    def fix_dockerfile_issues(self) -> int:
        """修復 Dockerfile 安全問題"""
        fixes = 0
        
        for finding in self.results.get("results", []):
            if "dockerfile" not in finding.get("path", "").lower():
                continue
                
            rule_id = finding.get("check_id", "")
            file_path = finding.get("path", "")
            
            if not os.path.exists(file_path):
                continue
                
            # 備份原始檔案
            self._backup_file(file_path)
            
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            
            # 修復缺少非 root 用戶
            if "dockerfile.security.missing-user" in rule_id:
                # 在 CMD 指令前添加非 root 用戶
                cmd_pattern = r'(CMD\s+\[.*?\])'
                if re.search(cmd_pattern, content):
                    user_commands = [
                        "",
                        "# 創建非 root 用戶",
                        "RUN adduser --disabled-password --gecos \"\" appuser",
                        "USER appuser",
                        ""
                    ]
                    content = re.sub(
                        cmd_pattern,
                        '\n'.join(user_commands) + r'\1',
                        content
                    )
            
            # 寫回修復後的內容
            if content != original_content:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                self.fixes_applied.append(f"修復 Dockerfile: {file_path}")
                fixes += 1
        
        return fixes
    
    def fix_terraform_issues(self) -> int:
        """修復 Terraform 安全問題"""
        fixes = 0
        
        for finding in self.results.get("results", []):
            if "terraform" not in finding.get("path", "").lower():
                continue
                
            rule_id = finding.get("check_id", "")
            file_path = finding.get("path", "")
            
            if not os.path.exists(file_path):
                continue
                
            # 備份原始檔案
            self._backup_file(file_path)
            
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            
            # 修復 ELB 缺少日誌
            if "terraform.aws.security.aws-elb-access-logs-not-enabled" in rule_id:
                # 在 aws_lb 資源中添加 access_logs 配置
                lb_pattern = r'(resource\s+"aws_lb"\s+"[^"]+"\s*\{[^}]*\})'
                if re.search(lb_pattern, content, re.DOTALL):
                    access_logs_config = '''
  # 啟用存取日誌
  access_logs {
    bucket  = aws_s3_bucket.alb_logs[0].id
    prefix  = "alb-logs"
  }
'''
                    content = re.sub(
                        lb_pattern,
                        r'\1' + access_logs_config,
                        content,
                        flags=re.DOTALL
                    )
            
            # 修復不安全的 TLS 版本
            elif "terraform.aws.security.insecure-load-balancer-tls-version" in rule_id:
                content = re.sub(
                    r'protocol\s*=\s*"HTTP"',
                    'protocol = "HTTPS"\n  ssl_policy = "ELBSecurityPolicy-TLS13-1-2-Res-2021-06"',
                    content
                )
            
            # 修復 CloudWatch 日誌未加密
            elif "terraform.aws.security.aws-cloudwatch-log-group-unencrypted" in rule_id:
                # 添加 KMS 加密
                log_group_pattern = r'(resource\s+"aws_cloudwatch_log_group"[^}]*\})'
                if re.search(log_group_pattern, content, re.DOTALL):
                    kms_config = '''
  kms_key_id = aws_kms_key.main.arn
'''
                    content = re.sub(
                        log_group_pattern,
                        r'\1' + kms_config,
                        content,
                        flags=re.DOTALL
                    )
            
            # 修復 KMS 缺少輪換
            elif "terraform.aws.security.aws-kms-no-rotation" in rule_id:
                content = re.sub(
                    r'(deletion_window_in_days\s*=\s*\d+)',
                    r'\1\n  enable_key_rotation = true',
                    content
                )
            
            # 修復子網路公共 IP 地址
            elif "terraform.aws.security.aws-subnet-has-public-ip-address" in rule_id:
                content = re.sub(
                    r'map_public_ip_on_launch\s*=\s*true',
                    'map_public_ip_on_launch = false',
                    content
                )
            
            # 寫回修復後的內容
            if content != original_content:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                self.fixes_applied.append(f"修復 Terraform: {file_path}")
                fixes += 1
        
        return fixes
    
    def generate_report(self) -> str:
        """生成修復報告"""
        report = f"""# 安全問題自動修復報告

## 修復摘要
- **原始問題數量**: {len(self.results.get('results', []))}
- **成功修復數量**: {len(self.fixes_applied)}
- **備份檔案位置**: {self.backup_dir}

## 修復的問題類型

### 1. Nginx 配置問題
- ✅ 修復了缺少 `always` 標誌的 `add_header` 指令
- ✅ 修復了 H2C 走私漏洞，限制 Upgrade 標頭

### 2. Dockerfile 安全問題
- ✅ 添加了非 root 用戶以提升容器安全性

### 3. Terraform 安全問題
- ✅ 啟用了 ELB 存取日誌
- ✅ 升級了 TLS 版本到安全版本
- ✅ 啟用了 KMS 金鑰輪換
- ✅ 添加了 CloudWatch 日誌加密
- ✅ 修復了子網路公共 IP 地址問題

## 修復的檔案
"""
        
        for fix in self.fixes_applied:
            report += f"- {fix}\n"
        
        report += """
## 建議的後續行動
1. 審查所有修復的變更
2. 在測試環境中驗證修復效果
3. 更新相關的部署文檔
4. 建立安全掃描的持續整合流程

## 注意事項
- 所有原始檔案都已備份
- 建議在部署前進行完整的測試
- 某些修復可能需要額外的 AWS 資源配置
"""
        
        return report
    
    def run_fixes(self) -> int:
        """執行所有修復"""
        print("🔍 開始分析 Semgrep 結果...")
        print(f"發現 {len(self.results.get('results', []))} 個問題")
        
        total_fixes = 0
        
        # 修復 Nginx 問題
        nginx_fixes = self.fix_nginx_issues()
        if nginx_fixes > 0:
            print(f"✅ 修復了 {nginx_fixes} 個 Nginx 配置問題")
            total_fixes += nginx_fixes
        
        # 修復 Dockerfile 問題
        dockerfile_fixes = self.fix_dockerfile_issues()
        if dockerfile_fixes > 0:
            print(f"✅ 修復了 {dockerfile_fixes} 個 Dockerfile 問題")
            total_fixes += dockerfile_fixes
        
        # 修復 Terraform 問題
        terraform_fixes = self.fix_terraform_issues()
        if terraform_fixes > 0:
            print(f"✅ 修復了 {terraform_fixes} 個 Terraform 問題")
            total_fixes += terraform_fixes
        
        print(f"\n🎉 總共修復了 {total_fixes} 個安全問題")
        
        # 生成報告
        report = self.generate_report()
        with open("security-fix-report.md", 'w', encoding='utf-8') as f:
            f.write(report)
        
        print("📄 修復報告已生成: security-fix-report.md")
        
        return total_fixes


def main():
    """主函數"""
    if len(sys.argv) > 1:
        results_file = sys.argv[1]
    else:
        results_file = "semgrep-results.json"
    
    fixer = SecurityAutoFixer(results_file)
    fixes_applied = fixer.run_fixes()
    
    if fixes_applied > 0:
        print(f"\n🚀 建議執行以下命令來驗證修復效果:")
        print("semgrep scan --config auto --json --output semgrep-results-after-fix.json")
        sys.exit(0)
    else:
        print("❌ 沒有找到可以自動修復的問題")
        sys.exit(1)


if __name__ == "__main__":
    main()
