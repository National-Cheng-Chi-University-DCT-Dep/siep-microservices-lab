# =============================================================================
# Ultimate Security Intelligence Platform - RDS Module Variables
# RDS 模組的變數定義
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

variable "private_subnet_ids" {
  description = "私有子網路 ID 列表"
  type        = list(string)
  validation {
    condition     = length(var.private_subnet_ids) >= 2
    error_message = "至少需要 2 個私有子網路以確保高可用性。"
  }
}

variable "allowed_security_groups" {
  description = "允許存取資料庫的安全群組 ID 列表"
  type        = list(string)
  default     = []
}

variable "allowed_cidr_blocks" {
  description = "允許存取資料庫的 CIDR 區塊列表"
  type        = list(string)
  default     = []
}

# -----------------------------------------------------------------------------
# 資料庫實例設定變數
# -----------------------------------------------------------------------------

variable "db_instance_class" {
  description = "RDS 實例類型"
  type        = string
  default     = "db.t3.micro"
  validation {
    condition     = can(regex("^db\\.", var.db_instance_class))
    error_message = "實例類型必須以 'db.' 開頭。"
  }
}

variable "db_allocated_storage" {
  description = "RDS 初始儲存空間 (GB)"
  type        = number
  default     = 20
  validation {
    condition     = var.db_allocated_storage >= 20 && var.db_allocated_storage <= 65536
    error_message = "儲存空間必須在 20GB 到 65536GB 之間。"
  }
}

variable "db_max_allocated_storage" {
  description = "RDS 最大儲存空間 (GB)"
  type        = number
  default     = 100
  validation {
    condition     = var.db_max_allocated_storage >= 20
    error_message = "最大儲存空間必須至少 20GB。"
  }
}

variable "db_engine_version" {
  description = "PostgreSQL 版本"
  type        = string
  default     = "15.4"
}

# -----------------------------------------------------------------------------
# 資料庫認證設定變數
# -----------------------------------------------------------------------------

variable "db_username" {
  description = "資料庫主使用者名稱"
  type        = string
  default     = "security_user"
  validation {
    condition     = can(regex("^[a-zA-Z][a-zA-Z0-9_]*$", var.db_username))
    error_message = "使用者名稱必須以字母開頭，只能包含字母、數字和底線。"
  }
}

variable "db_database_name" {
  description = "預設資料庫名稱"
  type        = string
  default     = "security_intel"
  validation {
    condition     = can(regex("^[a-zA-Z][a-zA-Z0-9_]*$", var.db_database_name))
    error_message = "資料庫名稱必須以字母開頭，只能包含字母、數字和底線。"
  }
}

# -----------------------------------------------------------------------------
# 備份和維護設定變數
# -----------------------------------------------------------------------------

variable "db_backup_retention_period" {
  description = "資料庫備份保留天數"
  type        = number
  default     = 7
  validation {
    condition     = var.db_backup_retention_period >= 0 && var.db_backup_retention_period <= 35
    error_message = "備份保留期間必須在 0 到 35 天之間。"
  }
}

variable "db_skip_final_snapshot" {
  description = "刪除實例時跳過最終快照"
  type        = bool
  default     = true
}

variable "enable_deletion_protection" {
  description = "啟用刪除保護"
  type        = bool
  default     = false
}

# -----------------------------------------------------------------------------
# 高可用性和效能設定變數
# -----------------------------------------------------------------------------

variable "create_read_replica" {
  description = "建立讀取副本"
  type        = bool
  default     = false
}

variable "read_replica_instance_class" {
  description = "讀取副本實例類型"
  type        = string
  default     = "db.t3.micro"
}

variable "create_manual_snapshot" {
  description = "建立手動快照"
  type        = bool
  default     = false
}

# -----------------------------------------------------------------------------
# 監控設定變數
# -----------------------------------------------------------------------------

variable "log_retention_days" {
  description = "CloudWatch 日誌保留天數"
  type        = number
  default     = 14
  validation {
    condition     = contains([1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, 3653], var.log_retention_days)
    error_message = "日誌保留天數必須是 CloudWatch Logs 支援的值。"
  }
}

variable "sns_topic_arn" {
  description = "CloudWatch 警報通知的 SNS Topic ARN"
  type        = string
  default     = null
}

# -----------------------------------------------------------------------------
# 安全性設定變數
# -----------------------------------------------------------------------------

variable "enable_performance_insights" {
  description = "啟用效能洞察"
  type        = bool
  default     = true
}

variable "performance_insights_retention_period" {
  description = "效能洞察資料保留天數"
  type        = number
  default     = 7
  validation {
    condition     = contains([7, 731], var.performance_insights_retention_period)
    error_message = "效能洞察保留期間必須是 7 或 731 天。"
  }
}

# =============================================================================
# 變數說明
# =============================================================================
# 1. private_subnet_ids: RDS 將部署在私有子網路中以提高安全性
# 2. db_password: 將由 random_password 資源自動生成
# 3. 備份設定: 生產環境建議設定較長的保留期間
# 4. 刪除保護: 生產環境強烈建議啟用
# 5. 監控: 效能洞察和 CloudWatch 日誌有助於故障排除
# ============================================================================= 