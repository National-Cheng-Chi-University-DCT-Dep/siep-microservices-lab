# =============================================================================
# Ultimate Security Intelligence Platform - Security Module Variables
# 安全性模組的變數定義
# =============================================================================

# -----------------------------------------------------------------------------
# 基本設定變數
# -----------------------------------------------------------------------------

variable "project_name" {
  description = "專案名稱"
  type        = string
}

variable "environment" {
  description = "環境名稱"
  type        = string
}

variable "aws_region" {
  description = "AWS 區域"
  type        = string
}

variable "common_tags" {
  description = "通用標籤"
  type        = map(string)
  default     = {}
}

# -----------------------------------------------------------------------------
# 網路設定變數
# -----------------------------------------------------------------------------

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR 區塊"
  type        = string
}

# -----------------------------------------------------------------------------
# IAM 角色設定變數
# -----------------------------------------------------------------------------

variable "create_ecs_task_execution_role" {
  description = "建立 ECS 任務執行角色"
  type        = bool
  default     = true
}

variable "create_ecs_task_role" {
  description = "建立 ECS 任務角色"
  type        = bool
  default     = true
}

variable "ecs_task_role_arn" {
  description = "現有的 ECS 任務角色 ARN（如果不建立新的）"
  type        = string
  default     = null
}

variable "s3_bucket_arns" {
  description = "允許存取的 S3 儲存桶 ARN 列表"
  type        = list(string)
  default     = []
}

# -----------------------------------------------------------------------------
# Secrets Manager 設定變數
# -----------------------------------------------------------------------------

variable "jwt_secret" {
  description = "JWT 簽名密鑰"
  type        = string
  sensitive   = true
}

variable "encryption_key" {
  description = "資料加密金鑰"
  type        = string
  sensitive   = true
}

variable "api_key" {
  description = "API 金鑰"
  type        = string
  sensitive   = true
}

variable "third_party_tokens" {
  description = "第三方服務 Token"
  type        = map(string)
  default     = {}
  sensitive   = true
}

# -----------------------------------------------------------------------------
# KMS 設定變數
# -----------------------------------------------------------------------------

variable "kms_deletion_window_in_days" {
  description = "KMS 金鑰刪除等待期間（天）"
  type        = number
  default     = 30
  validation {
    condition     = var.kms_deletion_window_in_days >= 7 && var.kms_deletion_window_in_days <= 30
    error_message = "KMS 金鑰刪除等待期間必須在 7 到 30 天之間。"
  }
}

# -----------------------------------------------------------------------------
# WAF 設定變數
# -----------------------------------------------------------------------------

variable "enable_waf" {
  description = "啟用 WAF (Web Application Firewall)"
  type        = bool
  default     = true
}

variable "waf_rate_limit" {
  description = "WAF 速率限制（每 5 分鐘的請求數）"
  type        = number
  default     = 2000
  validation {
    condition     = var.waf_rate_limit >= 100 && var.waf_rate_limit <= 20000000
    error_message = "WAF 速率限制必須在 100 到 20,000,000 之間。"
  }
}

# -----------------------------------------------------------------------------
# SSL/TLS 設定變數
# -----------------------------------------------------------------------------

variable "enable_ssl" {
  description = "啟用 SSL/TLS 憑證"
  type        = bool
  default     = false
}

variable "domain_name" {
  description = "主要域名（用於 SSL 憑證）"
  type        = string
  default     = null
}

variable "subject_alternative_names" {
  description = "SSL 憑證的替代域名列表"
  type        = list(string)
  default     = []
}

# -----------------------------------------------------------------------------
# 安全性警報設定變數
# -----------------------------------------------------------------------------

variable "enable_security_alerts" {
  description = "啟用安全性警報"
  type        = bool
  default     = true
}

variable "security_alert_email" {
  description = "安全性警報通知郵件地址"
  type        = string
  default     = null
}

variable "enable_security_monitoring" {
  description = "啟用安全性監控"
  type        = bool
  default     = true
}

# -----------------------------------------------------------------------------
# VPC Flow Logs 設定變數
# -----------------------------------------------------------------------------

variable "enable_vpc_flow_logs" {
  description = "啟用 VPC Flow Logs"
  type        = bool
  default     = true
}

# -----------------------------------------------------------------------------
# 監控設定變數
# -----------------------------------------------------------------------------

variable "log_retention_days" {
  description = "CloudWatch 日誌保留天數"
  type        = number
  default     = 30
  validation {
    condition     = contains([1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, 3653], var.log_retention_days)
    error_message = "日誌保留天數必須是 CloudWatch Logs 支援的值。"
  }
}

# -----------------------------------------------------------------------------
# 安全性設定變數
# -----------------------------------------------------------------------------

variable "allowed_cidr_blocks" {
  description = "允許存取的 CIDR 區塊列表"
  type        = list(string)
  default     = []
}

variable "enable_database_encryption" {
  description = "啟用資料庫加密"
  type        = bool
  default     = true
}

variable "enable_secrets_encryption" {
  description = "啟用 Secrets Manager 加密"
  type        = bool
  default     = true
}

variable "password_policy" {
  description = "密碼政策設定"
  type = object({
    minimum_password_length        = number
    require_lowercase_characters   = bool
    require_uppercase_characters   = bool
    require_numbers                = bool
    require_symbols                = bool
    allow_users_to_change_password = bool
    max_password_age               = number
    password_reuse_prevention      = number
  })
  default = {
    minimum_password_length        = 12
    require_lowercase_characters   = true
    require_uppercase_characters   = true
    require_numbers                = true
    require_symbols                = true
    allow_users_to_change_password = true
    max_password_age               = 90
    password_reuse_prevention      = 24
  }
}

# -----------------------------------------------------------------------------
# 網路安全設定變數
# -----------------------------------------------------------------------------

variable "enable_network_acls" {
  description = "啟用網路 ACL"
  type        = bool
  default     = false
}

variable "enable_security_groups_logging" {
  description = "啟用安全群組日誌記錄"
  type        = bool
  default     = true
}

# -----------------------------------------------------------------------------
# 合規性設定變數
# -----------------------------------------------------------------------------

variable "enable_config_rules" {
  description = "啟用 AWS Config 規則"
  type        = bool
  default     = false
}

variable "enable_cloudtrail" {
  description = "啟用 CloudTrail"
  type        = bool
  default     = true
}

variable "enable_guardduty" {
  description = "啟用 GuardDuty"
  type        = bool
  default     = false
}

# =============================================================================
# 變數說明
# =============================================================================
# 1. jwt_secret, encryption_key, api_key: 敏感資訊，請使用環境變數或 tfvars 檔案設定
# 2. domain_name: 設定後將自動建立 SSL 憑證
# 3. security_alert_email: 設定後將接收安全性警報通知
# 4. WAF 設定: 提供基本的 Web 應用程式防火牆保護
# 5. VPC Flow Logs: 記錄網路流量，用於安全性分析
# 6. 合規性功能: 可根據需求啟用進階合規性監控
# ============================================================================= 