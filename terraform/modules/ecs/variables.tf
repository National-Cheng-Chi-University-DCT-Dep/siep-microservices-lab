# =============================================================================
# Ultimate Security Intelligence Platform - ECS Module Variables
# ECS 模組的變數定義
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

variable "public_subnet_ids" {
  description = "公開子網路 ID 列表（用於 ALB）"
  type        = list(string)
  validation {
    condition     = length(var.public_subnet_ids) >= 2
    error_message = "至少需要 2 個公開子網路以確保 ALB 的高可用性。"
  }
}

variable "private_subnet_ids" {
  description = "私有子網路 ID 列表（用於 ECS 任務）"
  type        = list(string)
  validation {
    condition     = length(var.private_subnet_ids) >= 2
    error_message = "至少需要 2 個私有子網路以確保 ECS 的高可用性。"
  }
}

variable "allowed_cidr_blocks" {
  description = "允許存取 ALB 的 CIDR 區塊列表"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

# -----------------------------------------------------------------------------
# 應用程式設定變數
# -----------------------------------------------------------------------------

variable "app_image" {
  description = "應用程式 Docker 映像"
  type        = string
}

variable "container_port" {
  description = "容器端口"
  type        = number
  default     = 8080
  validation {
    condition     = var.container_port > 0 && var.container_port < 65536
    error_message = "容器端口必須在 1 到 65535 之間。"
  }
}

variable "task_cpu" {
  description = "ECS 任務 CPU 單位"
  type        = number
  default     = 256
  validation {
    condition     = contains([256, 512, 1024, 2048, 4096], var.task_cpu)
    error_message = "CPU 必須是 256, 512, 1024, 2048, 或 4096。"
  }
}

variable "task_memory" {
  description = "ECS 任務記憶體 (MB)"
  type        = number
  default     = 512
  validation {
    condition     = var.task_memory >= 512 && var.task_memory <= 30720
    error_message = "記憶體必須在 512MB 到 30720MB 之間。"
  }
}

variable "desired_count" {
  description = "ECS 服務所需的任務數量"
  type        = number
  default     = 1
  validation {
    condition     = var.desired_count >= 1
    error_message = "所需任務數量必須至少為 1。"
  }
}

# -----------------------------------------------------------------------------
# 負載平衡器設定變數
# -----------------------------------------------------------------------------

variable "alb_health_check_path" {
  description = "ALB 健康檢查路徑"
  type        = string
  default     = "/health"
}

variable "alb_health_check_interval" {
  description = "ALB 健康檢查間隔 (秒)"
  type        = number
  default     = 30
  validation {
    condition     = var.alb_health_check_interval >= 5 && var.alb_health_check_interval <= 300
    error_message = "健康檢查間隔必須在 5 到 300 秒之間。"
  }
}

variable "alb_health_check_timeout" {
  description = "ALB 健康檢查超時 (秒)"
  type        = number
  default     = 5
  validation {
    condition     = var.alb_health_check_timeout >= 2 && var.alb_health_check_timeout <= 120
    error_message = "健康檢查超時必須在 2 到 120 秒之間。"
  }
}

variable "alb_healthy_threshold" {
  description = "ALB 健康檢查成功次數閾值"
  type        = number
  default     = 2
  validation {
    condition     = var.alb_healthy_threshold >= 2 && var.alb_healthy_threshold <= 10
    error_message = "健康檢查成功次數閾值必須在 2 到 10 之間。"
  }
}

variable "alb_unhealthy_threshold" {
  description = "ALB 健康檢查失敗次數閾值"
  type        = number
  default     = 2
  validation {
    condition     = var.alb_unhealthy_threshold >= 2 && var.alb_unhealthy_threshold <= 10
    error_message = "健康檢查失敗次數閾值必須在 2 到 10 之間。"
  }
}

# -----------------------------------------------------------------------------
# 環境變數和機密設定變數
# -----------------------------------------------------------------------------

variable "environment_variables" {
  description = "ECS 任務環境變數"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "secrets_from_secrets_manager" {
  description = "從 AWS Secrets Manager 取得的機密變數"
  type = list(object({
    name      = string
    valueFrom = string
  }))
  default = []
}

variable "secrets_manager_arns" {
  description = "允許存取的 Secrets Manager ARN 列表"
  type        = list(string)
  default     = []
}

# -----------------------------------------------------------------------------
# 自動擴展設定變數
# -----------------------------------------------------------------------------

variable "enable_auto_scaling" {
  description = "啟用自動擴展"
  type        = bool
  default     = true
}

variable "min_capacity" {
  description = "最小任務數量"
  type        = number
  default     = 1
}

variable "max_capacity" {
  description = "最大任務數量"
  type        = number
  default     = 10
}

variable "cpu_target_value" {
  description = "CPU 使用率目標值 (%)"
  type        = number
  default     = 70
  validation {
    condition     = var.cpu_target_value > 0 && var.cpu_target_value <= 100
    error_message = "CPU 目標值必須在 1% 到 100% 之間。"
  }
}

variable "memory_target_value" {
  description = "記憶體使用率目標值 (%)"
  type        = number
  default     = 80
  validation {
    condition     = var.memory_target_value > 0 && var.memory_target_value <= 100
    error_message = "記憶體目標值必須在 1% 到 100% 之間。"
  }
}

# -----------------------------------------------------------------------------
# 服務發現設定變數
# -----------------------------------------------------------------------------

variable "enable_service_discovery" {
  description = "啟用服務發現"
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

variable "enable_deletion_protection" {
  description = "啟用 ALB 刪除保護"
  type        = bool
  default     = false
}

# =============================================================================
# 變數說明
# =============================================================================
# 1. app_image: 應該是完整的 Docker 映像 URI，例如：
#    123456789012.dkr.ecr.us-east-1.amazonaws.com/my-app:latest
# 2. task_cpu 和 task_memory: 必須是 Fargate 支援的組合
# 3. environment_variables: 明碼環境變數
# 4. secrets_from_secrets_manager: 從 AWS Secrets Manager 取得的機密
# 5. 自動擴展: 基於 CPU 和記憶體使用率進行擴展
# 6. 服務發現: 啟用後可通過 DNS 名稱訪問服務
# ============================================================================= 