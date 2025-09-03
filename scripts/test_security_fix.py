#!/usr/bin/env python3
"""
測試安全自動修復工具
"""

import json
import tempfile
import os
from pathlib import Path
from scripts.security_auto_fix import SecurityAutoFixer


def create_test_nginx_config():
    """建立測試用的 Nginx 配置"""
    config = '''
server {
    listen 80;
    server_name example.com;
    
    location / {
        add_header Content-Type text/plain;
        proxy_pass http://backend;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
'''
    return config


def create_test_dockerfile():
    """建立測試用的 Dockerfile"""
    dockerfile = '''
FROM python:3.9-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt

COPY . .
CMD ["python", "app.py"]
'''
    return dockerfile


def create_test_terraform():
    """建立測試用的 Terraform 配置"""
    terraform = '''
resource "aws_lb" "main" {
  name               = "test-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets           = var.public_subnet_ids

  enable_deletion_protection = false
}

resource "aws_kms_key" "rds" {
  description = "KMS key for RDS encryption"
  deletion_window_in_days = 7
}

resource "aws_cloudwatch_log_group" "postgresql" {
  name              = "/aws/rds/instance/test-postgres/postgresql"
  retention_in_days = var.log_retention_days
}

resource "aws_subnet" "public" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-west-2a"
  map_public_ip_on_launch = true
}
'''
    return terraform


def create_test_semgrep_results():
    """建立測試用的 Semgrep 結果"""
    results = {
        "results": [
            {
                "check_id": "generic.nginx.security.missing-always-flag",
                "path": "nginx/test.conf",
                "start": {"line": 7, "col": 9},
                "end": {"line": 7, "col": 35},
                "extra": {
                    "message": "Missing 'always' flag in add_header directive"
                }
            },
            {
                "check_id": "generic.nginx.security.possible-h2c-smuggling",
                "path": "nginx/test.conf",
                "start": {"line": 9, "col": 9},
                "end": {"line": 9, "col": 35},
                "extra": {
                    "message": "Possible H2C smuggling vulnerability"
                }
            },
            {
                "check_id": "dockerfile.security.missing-user",
                "path": "Dockerfile",
                "start": {"line": 7, "col": 1},
                "end": {"line": 7, "col": 25},
                "extra": {
                    "message": "Missing non-root user"
                }
            },
            {
                "check_id": "terraform.aws.security.aws-elb-access-logs-not-enabled",
                "path": "main.tf",
                "start": {"line": 1, "col": 1},
                "end": {"line": 10, "col": 1},
                "extra": {
                    "message": "ELB has no logging"
                }
            },
            {
                "check_id": "terraform.aws.security.aws-kms-no-rotation",
                "path": "main.tf",
                "start": {"line": 12, "col": 1},
                "end": {"line": 15, "col": 1},
                "extra": {
                    "message": "KMS has no rotation"
                }
            },
            {
                "check_id": "terraform.aws.security.aws-cloudwatch-log-group-unencrypted",
                "path": "main.tf",
                "start": {"line": 17, "col": 1},
                "end": {"line": 20, "col": 1},
                "extra": {
                    "message": "CloudWatch log group is unencrypted"
                }
            },
            {
                "check_id": "terraform.aws.security.aws-subnet-has-public-ip-address",
                "path": "main.tf",
                "start": {"line": 22, "col": 1},
                "end": {"line": 27, "col": 1},
                "extra": {
                    "message": "Subnet has public IP address"
                }
            }
        ]
    }
    return results


def test_security_auto_fix():
    """測試安全自動修復功能"""
    print("🧪 開始測試安全自動修復工具...")
    
    # 建立臨時目錄
    with tempfile.TemporaryDirectory() as temp_dir:
        os.chdir(temp_dir)
        
        # 建立測試檔案
        print("📝 建立測試檔案...")
        
        # Nginx 配置
        nginx_dir = Path("nginx")
        nginx_dir.mkdir()
        with open(nginx_dir / "test.conf", "w") as f:
            f.write(create_test_nginx_config())
        
        # Dockerfile
        with open("Dockerfile", "w") as f:
            f.write(create_test_dockerfile())
        
        # Terraform
        with open("main.tf", "w") as f:
            f.write(create_test_terraform())
        
        # Semgrep 結果
        with open("semgrep-results.json", "w") as f:
            json.dump(create_test_semgrep_results(), f)
        
        print(f"✅ 測試檔案已建立: {temp_dir}")
        
        # 執行修復
        print("🔧 執行自動修復...")
        fixer = SecurityAutoFixer("semgrep-results.json")
        fixes_applied = fixer.run_fixes()
        
        # 驗證結果
        print("🔍 驗證修復結果...")
        
        # 檢查 Nginx 修復
        with open(nginx_dir / "test.conf", "r") as f:
            nginx_content = f.read()
            if "always;" in nginx_content:
                print("✅ Nginx 配置已修復")
            else:
                print("❌ Nginx 配置修復失敗")
        
        # 檢查 Dockerfile 修復
        with open("Dockerfile", "r") as f:
            dockerfile_content = f.read()
            if "USER appuser" in dockerfile_content:
                print("✅ Dockerfile 已修復")
            else:
                print("❌ Dockerfile 修復失敗")
        
        # 檢查 Terraform 修復
        with open("main.tf", "r") as f:
            terraform_content = f.read()
            terraform_fixes = 0
            if "access_logs" in terraform_content:
                terraform_fixes += 1
            if "enable_key_rotation = true" in terraform_content:
                terraform_fixes += 1
            if "kms_key_id" in terraform_content:
                terraform_fixes += 1
            if "map_public_ip_on_launch = false" in terraform_content:
                terraform_fixes += 1
            
            print(f"✅ Terraform 修復了 {terraform_fixes}/4 個問題")
        
        # 檢查備份
        if Path("backups").exists():
            backup_files = list(Path("backups").glob("*.backup"))
            print(f"✅ 建立了 {len(backup_files)} 個備份檔案")
        else:
            print("❌ 沒有建立備份檔案")
        
        # 檢查報告
        if Path("security-fix-report.md").exists():
            print("✅ 修復報告已生成")
        else:
            print("❌ 修復報告生成失敗")
        
        print(f"\n🎉 測試完成！總共修復了 {fixes_applied} 個問題")


if __name__ == "__main__":
    test_security_auto_fix()
