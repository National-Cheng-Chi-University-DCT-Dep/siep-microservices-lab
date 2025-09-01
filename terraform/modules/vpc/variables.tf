# =============================================================================
# Ultimate Security Intelligence Platform - VPC Module Variables
# VPC 模組的變數定義
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

variable "vpc_cidr" {
  description = "VPC CIDR 區塊"
  type        = string
  validation {
    condition     = can(cidrhost(var.vpc_cidr, 0))
    error_message = "VPC CIDR 必須是有效的 CIDR 區塊。"
  }
}

variable "availability_zones" {
  description = "可用性區域列表"
  type        = list(string)
  validation {
    condition     = length(var.availability_zones) >= 2
    error_message = "至少需要 2 個可用性區域以確保高可用性。"
  }
}

variable "public_subnet_cidrs" {
  description = "公開子網路 CIDR 區塊列表"
  type        = list(string)
  validation {
    condition     = length(var.public_subnet_cidrs) >= 2
    error_message = "至少需要 2 個公開子網路以確保高可用性。"
  }
}

variable "private_subnet_cidrs" {
  description = "私有子網路 CIDR 區塊列表"
  type        = list(string)
  validation {
    condition     = length(var.private_subnet_cidrs) >= 2
    error_message = "至少需要 2 個私有子網路以確保高可用性。"
  }
}

# -----------------------------------------------------------------------------
# 可選功能變數
# -----------------------------------------------------------------------------

variable "enable_vpc_flow_logs" {
  description = "啟用 VPC Flow Logs"
  type        = bool
  default     = true
}

variable "enable_s3_endpoint" {
  description = "啟用 S3 VPC Endpoint"
  type        = bool
  default     = true
}

variable "enable_network_acls" {
  description = "啟用自定義網路 ACL"
  type        = bool
  default     = false
}

# -----------------------------------------------------------------------------
# 安全性設定變數
# -----------------------------------------------------------------------------

variable "ssh_allowed_cidr" {
  description = "允許 SSH 連接的 CIDR 區塊"
  type        = string
  default     = "10.0.0.0/8"
}

# -----------------------------------------------------------------------------
# 監控設定變數
# -----------------------------------------------------------------------------

variable "log_retention_days" {
  description = "日誌保留天數"
  type        = number
  default     = 14
  validation {
    condition     = contains([1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, 3653], var.log_retention_days)
    error_message = "日誌保留天數必須是 CloudWatch Logs 支援的值。"
  }
}

# =============================================================================
# 變數說明
# =============================================================================
# 1. vpc_cidr: 建議使用私有 IP 範圍，如 10.0.0.0/16
# 2. availability_zones: 需要與 subnet_cidrs 數量匹配
# 3. 子網路 CIDR: 需要在 VPC CIDR 範圍內
# 4. 安全性設定: 生產環境請限制 SSH 存取來源
# ============================================================================= 