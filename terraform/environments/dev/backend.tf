# =============================================================================
# Ultimate Security Intelligence Platform - Terraform Backend Configuration
# 設定遠端狀態後端，使用 S3 儲存 tfstate 檔案
# =============================================================================

terraform {
  # 最低 Terraform 版本要求
  required_version = ">= 1.0"
  
  # 必要的 Provider 版本
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
  }
  
  # 遠端狀態後端設定
  backend "s3" {
    # 請將以下設定替換為您的實際值
    bucket         = "security-intel-tfstate-2024"    # 您的 S3 儲存桶名稱
    key            = "dev/terraform.tfstate"          # 狀態檔案路徑
    region         = "ap-northeast-1"                 # AWS 區域
    encrypt        = true                             # 啟用加密
    dynamodb_table = "security-intel-tfstate-lock"   # DynamoDB 表格用於狀態鎖定
    
    # 可選：使用 KMS 加密
    # kms_key_id = "alias/terraform-state-key"
  }
}

# =============================================================================
# Provider 設定
# =============================================================================

# AWS Provider 設定
provider "aws" {
  region = var.aws_region
  
  # 預設標籤，會自動套用到所有支援的資源
  default_tags {
    tags = {
      Environment = var.environment
      Project     = var.project_name
      ManagedBy   = "terraform"
      Owner       = var.owner
    }
  }
}

# =============================================================================
# 設定說明
# =============================================================================
# 1. 請先在 AWS 中手動建立 S3 儲存桶和 DynamoDB 表格
# 2. 確保您的 AWS 憑證有適當的權限
# 3. 建議啟用 S3 儲存桶的版本控制和加密
# 4. DynamoDB 表格用於防止多人同時修改狀態檔案
# ============================================================================= 