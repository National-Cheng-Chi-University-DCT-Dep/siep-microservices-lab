# =============================================================================
# Ultimate Security Intelligence Platform - Development Environment Variables
# 開發環境的變數定義
# =============================================================================

# -----------------------------------------------------------------------------
# 基本設定變數
# -----------------------------------------------------------------------------

variable "aws_region" {
  description = "AWS 區域"
  type        = string
  default     = "ap-northeast-1"
}

variable "environment" {
  description = "環境名稱"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "專案名稱"
  type        = string
  default     = "security-intel"
}

variable "owner" {
  description = "專案負責人"
  type        = string
  default     = "security-team"
}

# -----------------------------------------------------------------------------
# 網路設定變數
# -----------------------------------------------------------------------------

variable "vpc_cidr" {
  description = "VPC CIDR 區塊"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "可用性區域列表"
  type        = list(string)
  default     = ["ap-northeast-1a", "ap-northeast-1c"]
}

variable "public_subnet_cidrs" {
  description = "公開子網路 CIDR 區塊"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "私有子網路 CIDR 區塊"
  type        = list(string)
  default     = ["10.0.11.0/24", "10.0.12.0/24"]
}

# -----------------------------------------------------------------------------
# 資料庫設定變數
# -----------------------------------------------------------------------------

variable "db_instance_class" {
  description = "RDS 實例類型"
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage" {
  description = "RDS 儲存空間 (GB)"
  type        = number
  default     = 20
}

variable "db_max_allocated_storage" {
  description = "RDS 最大儲存空間 (GB)"
  type        = number
  default     = 100
}

variable "db_engine_version" {
  description = "PostgreSQL 版本"
  type        = string
  default     = "15.4"
}

variable "db_username" {
  description = "資料庫使用者名稱"
  type        = string
  default     = "security_user"
}

variable "db_database_name" {
  description = "資料庫名稱"
  type        = string
  default     = "security_intel"
}

variable "db_backup_retention_period" {
  description = "資料庫備份保留天數"
  type        = number
  default     = 7
}

variable "db_skip_final_snapshot" {
  description = "刪除時跳過最終快照"
  type        = bool
  default     = true
}

# -----------------------------------------------------------------------------
# ECS 設定變數
# -----------------------------------------------------------------------------

variable "ecs_task_cpu" {
  description = "ECS Task CPU 配置"
  type        = number
  default     = 256
}

variable "ecs_task_memory" {
  description = "ECS Task 記憶體配置 (MB)"
  type        = number
  default     = 512
}

variable "ecs_desired_count" {
  description = "ECS 服務所需的任務數量"
  type        = number
  default     = 1
}

variable "ecs_container_port" {
  description = "容器端口"
  type        = number
  default     = 8080
}

variable "backend_image_tag" {
  description = "後端容器映像標籤"
  type        = string
  default     = "latest"
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
}

variable "alb_health_check_timeout" {
  description = "ALB 健康檢查超時 (秒)"
  type        = number
  default     = 5
}

variable "alb_healthy_threshold" {
  description = "ALB 健康檢查成功次數閾值"
  type        = number
  default     = 2
}

variable "alb_unhealthy_threshold" {
  description = "ALB 健康檢查失敗次數閾值"
  type        = number
  default     = 2
}

# -----------------------------------------------------------------------------
# 安全性設定變數
# -----------------------------------------------------------------------------

variable "allowed_cidr_blocks" {
  description = "允許存取的 CIDR 區塊"
  type        = list(string)
  default     = ["0.0.0.0/0"] # 開發環境允許所有來源，生產環境應限制
}

variable "enable_deletion_protection" {
  description = "啟用刪除保護"
  type        = bool
  default     = false # 開發環境設為 false，生產環境應設為 true
}

# -----------------------------------------------------------------------------
# 監控設定變數
# -----------------------------------------------------------------------------

variable "enable_cloudwatch_logs" {
  description = "啟用 CloudWatch 日誌"
  type        = bool
  default     = true
}

variable "log_retention_in_days" {
  description = "日誌保留天數"
  type        = number
  default     = 14
}

# -----------------------------------------------------------------------------
# 標籤設定變數
# -----------------------------------------------------------------------------

variable "common_tags" {
  description = "通用標籤"
  type        = map(string)
  default = {
    Project     = "security-intel"
    Environment = "dev"
    ManagedBy   = "terraform"
    Team        = "security-team"
  }
}

# =============================================================================
# 說明
# =============================================================================
# 1. 這些變數可以透過 terraform.tfvars 檔案覆蓋
# 2. 敏感資訊請使用 AWS Secrets Manager 或環境變數
# 3. 生產環境建議使用更嚴格的預設值
# ============================================================================= 