# =============================================================================
# Ultimate Security Intelligence Platform - Development Environment
# 開發環境的主要配置，整合所有模組
# =============================================================================

# -----------------------------------------------------------------------------
# VPC 網路模組
# -----------------------------------------------------------------------------

module "vpc" {
  source = "../../modules/vpc"

  project_name         = var.project_name
  environment          = var.environment
  aws_region           = var.aws_region
  common_tags          = var.common_tags
  vpc_cidr             = var.vpc_cidr
  availability_zones   = var.availability_zones
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  enable_vpc_flow_logs = true
  enable_s3_endpoint   = true
  enable_network_acls  = false
  log_retention_days   = var.log_retention_in_days
}

# -----------------------------------------------------------------------------
# 安全性模組
# -----------------------------------------------------------------------------

module "security" {
  source = "../../modules/security"

  project_name                   = var.project_name
  environment                    = var.environment
  aws_region                     = var.aws_region
  common_tags                    = var.common_tags
  vpc_id                         = module.vpc.vpc_id
  vpc_cidr                       = module.vpc.vpc_cidr_block
  create_ecs_task_execution_role = true
  create_ecs_task_role           = true
  jwt_secret                     = var.jwt_secret
  encryption_key                 = var.encryption_key
  api_key                        = var.api_key
  third_party_tokens             = var.third_party_tokens
  enable_waf                     = var.enable_waf
  enable_ssl                     = var.enable_ssl
  domain_name                    = var.domain_name
  subject_alternative_names      = var.subject_alternative_names
  enable_security_alerts         = var.enable_security_alerts
  security_alert_email           = var.security_alert_email
  enable_vpc_flow_logs           = var.enable_vpc_flow_logs
  log_retention_days             = var.log_retention_in_days
}

# -----------------------------------------------------------------------------
# RDS 資料庫模組
# -----------------------------------------------------------------------------

module "rds" {
  source = "../../modules/rds"

  project_name               = var.project_name
  environment                = var.environment
  common_tags                = var.common_tags
  vpc_id                     = module.vpc.vpc_id
  private_subnet_ids         = module.vpc.private_subnet_ids
  allowed_security_groups    = [module.security.database_client_security_group_id]
  db_instance_class          = var.db_instance_class
  db_allocated_storage       = var.db_allocated_storage
  db_max_allocated_storage   = var.db_max_allocated_storage
  db_engine_version          = var.db_engine_version
  db_username                = var.db_username
  db_database_name           = var.db_database_name
  db_backup_retention_period = var.db_backup_retention_period
  db_skip_final_snapshot     = var.db_skip_final_snapshot
  enable_deletion_protection = var.enable_deletion_protection
  create_read_replica        = false # 開發環境不需要讀取副本
  log_retention_days         = var.log_retention_in_days
  sns_topic_arn              = module.security.security_alerts_topic_arn
}

# -----------------------------------------------------------------------------
# ECS 運算模組
# -----------------------------------------------------------------------------

module "ecs" {
  source = "../../modules/ecs"

  project_name               = var.project_name
  environment                = var.environment
  aws_region                 = var.aws_region
  common_tags                = var.common_tags
  vpc_id                     = module.vpc.vpc_id
  public_subnet_ids          = module.vpc.public_subnet_ids
  private_subnet_ids         = module.vpc.private_subnet_ids
  allowed_cidr_blocks        = var.allowed_cidr_blocks
  app_image                  = var.backend_image_tag
  container_port             = var.ecs_container_port
  task_cpu                   = var.ecs_task_cpu
  task_memory                = var.ecs_task_memory
  desired_count              = var.ecs_desired_count
  alb_health_check_path      = var.alb_health_check_path
  alb_health_check_interval  = var.alb_health_check_interval
  alb_health_check_timeout   = var.alb_health_check_timeout
  alb_healthy_threshold      = var.alb_healthy_threshold
  alb_unhealthy_threshold    = var.alb_unhealthy_threshold
  enable_auto_scaling        = true
  min_capacity               = 1
  max_capacity               = 3 # 開發環境限制較小的擴展
  cpu_target_value           = 70
  memory_target_value        = 80
  enable_service_discovery   = false # 開發環境暫時不啟用
  log_retention_days         = var.log_retention_in_days
  sns_topic_arn              = module.security.security_alerts_topic_arn
  enable_deletion_protection = var.enable_deletion_protection
  secrets_manager_arns       = [module.security.app_secrets_arn]

  # 環境變數配置
  environment_variables = [
    {
      name  = "ENVIRONMENT"
      value = var.environment
    },
    {
      name  = "DATABASE_HOST"
      value = module.rds.db_instance_endpoint
    },
    {
      name  = "DATABASE_PORT"
      value = tostring(module.rds.db_instance_port)
    },
    {
      name  = "DATABASE_NAME"
      value = module.rds.db_instance_database_name
    },
    {
      name  = "LOG_LEVEL"
      value = "debug"
    },
    {
      name  = "REDIS_HOST"
      value = "localhost" # 暫時使用本地 Redis
    },
    {
      name  = "REDIS_PORT"
      value = "6379"
    }
  ]

  # 從 Secrets Manager 獲取的機密
  secrets_from_secrets_manager = module.security.ecs_secrets_configuration
}

# -----------------------------------------------------------------------------
# CloudWatch Dashboard (監控儀表板)
# -----------------------------------------------------------------------------

resource "aws_cloudwatch_dashboard" "main" {
  dashboard_name = "${var.project_name}-${var.environment}-dashboard"

  dashboard_body = jsonencode({
    widgets = [
      {
        type   = "metric"
        x      = 0
        y      = 0
        width  = 12
        height = 6

        properties = {
          metrics = [
            ["AWS/ECS", "CPUUtilization", "ServiceName", module.ecs.ecs_service_name, "ClusterName", module.ecs.ecs_cluster_name],
            [".", "MemoryUtilization", ".", ".", ".", "."]
          ]
          period = 300
          stat   = "Average"
          region = var.aws_region
          title  = "ECS Service Metrics"
        }
      },
      {
        type   = "metric"
        x      = 0
        y      = 6
        width  = 12
        height = 6

        properties = {
          metrics = [
            ["AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", module.rds.db_instance_id],
            [".", "DatabaseConnections", ".", "."],
            [".", "FreeStorageSpace", ".", "."]
          ]
          period = 300
          stat   = "Average"
          region = var.aws_region
          title  = "RDS Metrics"
        }
      },
      {
        type   = "metric"
        x      = 0
        y      = 12
        width  = 12
        height = 6

        properties = {
          metrics = [
            ["AWS/ApplicationELB", "RequestCount", "LoadBalancer", module.ecs.alb_arn],
            [".", "TargetResponseTime", ".", "."],
            [".", "HTTPCode_Target_2XX_Count", ".", "."],
            [".", "HTTPCode_Target_4XX_Count", ".", "."],
            [".", "HTTPCode_Target_5XX_Count", ".", "."]
          ]
          period = 300
          stat   = "Sum"
          region = var.aws_region
          title  = "ALB Metrics"
        }
      }
    ]
  })

  tags = var.common_tags
}

# -----------------------------------------------------------------------------
# Route 53 DNS 記錄 (如果有域名)
# -----------------------------------------------------------------------------

data "aws_route53_zone" "main" {
  count = var.domain_name != null ? 1 : 0
  name  = var.domain_name
}

resource "aws_route53_record" "main" {
  count = var.domain_name != null ? 1 : 0

  zone_id = data.aws_route53_zone.main[0].zone_id
  name    = var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}"
  type    = "A"

  alias {
    name                   = module.ecs.alb_dns_name
    zone_id                = module.ecs.alb_zone_id
    evaluate_target_health = true
  }
}

# -----------------------------------------------------------------------------
# SSM Parameters (系統參數)
# -----------------------------------------------------------------------------

resource "aws_ssm_parameter" "database_endpoint" {
  name  = "/${var.project_name}/${var.environment}/database/endpoint"
  type  = "String"
  value = module.rds.db_instance_endpoint

  tags = var.common_tags
}

resource "aws_ssm_parameter" "database_port" {
  name  = "/${var.project_name}/${var.environment}/database/port"
  type  = "String"
  value = tostring(module.rds.db_instance_port)

  tags = var.common_tags
}

resource "aws_ssm_parameter" "alb_dns_name" {
  name  = "/${var.project_name}/${var.environment}/alb/dns_name"
  type  = "String"
  value = module.ecs.alb_dns_name

  tags = var.common_tags
}

resource "aws_ssm_parameter" "application_url" {
  name  = "/${var.project_name}/${var.environment}/application/url"
  type  = "String"
  value = var.domain_name != null ? "https://${var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}"}" : module.ecs.application_url

  tags = var.common_tags
}

# -----------------------------------------------------------------------------
# 本地值
# -----------------------------------------------------------------------------

locals {
  # 環境特定的標籤
  environment_tags = {
    Environment = var.environment
    Terraform   = "true"
    Repository  = "security-intelligence-platform"
  }

  # 合併所有標籤
  all_tags = merge(var.common_tags, local.environment_tags)
}

# =============================================================================
# 輸出說明
# =============================================================================
# 1. 此配置建立了完整的開發環境，包含 VPC、安全性、RDS、ECS
# 2. 所有模組都經過適當的配置，確保安全性和可維護性
# 3. CloudWatch Dashboard 提供了監控視圖
# 4. SSM Parameters 儲存了重要的系統參數
# 5. 如果有域名，會自動建立 DNS 記錄
# ============================================================================= 