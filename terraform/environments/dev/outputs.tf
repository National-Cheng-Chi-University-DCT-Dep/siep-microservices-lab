# =============================================================================
# Ultimate Security Intelligence Platform - Development Environment Outputs
# 開發環境的輸出值
# =============================================================================

# -----------------------------------------------------------------------------
# 網路資源輸出
# -----------------------------------------------------------------------------

output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "vpc_cidr_block" {
  description = "VPC CIDR 區塊"
  value       = module.vpc.vpc_cidr_block
}

output "public_subnet_ids" {
  description = "公開子網路 ID 列表"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "私有子網路 ID 列表"
  value       = module.vpc.private_subnet_ids
}

output "internet_gateway_id" {
  description = "Internet Gateway ID"
  value       = module.vpc.internet_gateway_id
}

output "nat_gateway_ids" {
  description = "NAT Gateway ID 列表"
  value       = module.vpc.nat_gateway_ids
}

# -----------------------------------------------------------------------------
# 應用程式存取資訊
# -----------------------------------------------------------------------------

output "application_url" {
  description = "應用程式存取 URL"
  value       = var.domain_name != null ? "https://${var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}"}" : module.ecs.application_url
}

output "application_health_check_url" {
  description = "應用程式健康檢查 URL"
  value       = "${module.ecs.application_url}${var.alb_health_check_path}"
}

output "alb_dns_name" {
  description = "Application Load Balancer DNS 名稱"
  value       = module.ecs.alb_dns_name
}

output "alb_zone_id" {
  description = "Application Load Balancer Zone ID"
  value       = module.ecs.alb_zone_id
}

# -----------------------------------------------------------------------------
# 資料庫資訊
# -----------------------------------------------------------------------------

output "database_endpoint" {
  description = "RDS 資料庫端點"
  value       = module.rds.db_instance_endpoint
}

output "database_port" {
  description = "資料庫端口"
  value       = module.rds.db_instance_port
}

output "database_name" {
  description = "資料庫名稱"
  value       = module.rds.db_instance_database_name
}

output "database_credentials_secret_arn" {
  description = "資料庫認證 Secret ARN"
  value       = module.rds.db_credentials_secret_arn
  sensitive   = true
}

# -----------------------------------------------------------------------------
# ECS 服務資訊
# -----------------------------------------------------------------------------

output "ecs_cluster_name" {
  description = "ECS 叢集名稱"
  value       = module.ecs.ecs_cluster_name
}

output "ecs_service_name" {
  description = "ECS 服務名稱"
  value       = module.ecs.ecs_service_name
}

output "ecs_task_definition_arn" {
  description = "ECS 任務定義 ARN"
  value       = module.ecs.ecs_task_definition_arn
}

# -----------------------------------------------------------------------------
# 安全性資源
# -----------------------------------------------------------------------------

output "kms_key_id" {
  description = "KMS 金鑰 ID"
  value       = module.security.kms_key_id
}

output "app_secrets_arn" {
  description = "應用程式 Secrets Manager ARN"
  value       = module.security.app_secrets_arn
}

output "waf_web_acl_arn" {
  description = "WAF Web ACL ARN"
  value       = module.security.waf_web_acl_arn
}

output "ssl_certificate_arn" {
  description = "SSL 憑證 ARN"
  value       = module.security.ssl_certificate_arn
}

# -----------------------------------------------------------------------------
# 監控資源
# -----------------------------------------------------------------------------

output "cloudwatch_dashboard_url" {
  description = "CloudWatch Dashboard URL"
  value       = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#dashboards:name=${var.project_name}-${var.environment}-dashboard"
}

output "security_alerts_topic_arn" {
  description = "安全性警報 SNS 主題 ARN"
  value       = module.security.security_alerts_topic_arn
}

# -----------------------------------------------------------------------------
# DNS 資訊
# -----------------------------------------------------------------------------

output "domain_name" {
  description = "應用程式域名"
  value       = var.domain_name != null ? (var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}") : null
}

output "route53_record_name" {
  description = "Route 53 記錄名稱"
  value       = var.domain_name != null ? aws_route53_record.main[0].name : null
}

# -----------------------------------------------------------------------------
# 系統參數
# -----------------------------------------------------------------------------

output "ssm_parameters" {
  description = "SSM 參數列表"
  value = {
    database_endpoint = aws_ssm_parameter.database_endpoint.name
    database_port     = aws_ssm_parameter.database_port.name
    alb_dns_name      = aws_ssm_parameter.alb_dns_name.name
    application_url   = aws_ssm_parameter.application_url.name
  }
}

# -----------------------------------------------------------------------------
# 連接資訊摘要
# -----------------------------------------------------------------------------

output "connection_info" {
  description = "系統連接資訊摘要"
  value = {
    application_url      = var.domain_name != null ? "https://${var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}"}" : module.ecs.application_url
    health_check_url     = "${module.ecs.application_url}${var.alb_health_check_path}"
    alb_dns_name         = module.ecs.alb_dns_name
    database_endpoint    = module.rds.db_instance_endpoint
    database_port        = module.rds.db_instance_port
    database_name        = module.rds.db_instance_database_name
    cloudwatch_dashboard = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#dashboards:name=${var.project_name}-${var.environment}-dashboard"
  }
}

# -----------------------------------------------------------------------------
# 部署資訊
# -----------------------------------------------------------------------------

output "deployment_info" {
  description = "部署資訊摘要"
  value = {
    environment             = var.environment
    aws_region              = var.aws_region
    vpc_id                  = module.vpc.vpc_id
    ecs_cluster_name        = module.ecs.ecs_cluster_name
    ecs_service_name        = module.ecs.ecs_service_name
    database_instance_id    = module.rds.db_instance_id
    kms_key_id              = module.security.kms_key_id
    waf_enabled             = var.enable_waf
    ssl_enabled             = var.enable_ssl
    auto_scaling_enabled    = true
    monitoring_enabled      = true
    security_alerts_enabled = var.enable_security_alerts
  }
}

# -----------------------------------------------------------------------------
# 開發者快速存取資訊
# -----------------------------------------------------------------------------

output "developer_quick_access" {
  description = "開發者快速存取資訊"
  value = {
    application_url = var.domain_name != null ? "https://${var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}"}" : module.ecs.application_url
    api_base_url    = "${var.domain_name != null ? "https://${var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}"}" : module.ecs.application_url}/api"
    swagger_ui_url  = "${var.domain_name != null ? "https://${var.environment == "prod" ? var.domain_name : "${var.environment}.${var.domain_name}"}" : module.ecs.application_url}/swagger"
    health_check    = "${module.ecs.application_url}${var.alb_health_check_path}"
    logs_url        = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#logsV2:log-groups/log-group/${replace(module.ecs.cloudwatch_log_group_name, "/", "$252F")}"
    metrics_url     = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#dashboards:name=${var.project_name}-${var.environment}-dashboard"
  }
}

# =============================================================================
# 輸出說明
# =============================================================================
# 1. application_url: 應用程式的主要存取點
# 2. developer_quick_access: 開發者常用的 URL 集合
# 3. connection_info: 系統各元件的連接資訊
# 4. deployment_info: 部署配置的摘要
# 5. 所有敏感資訊都已標記為 sensitive
# 6. URL 會根據是否設定域名自動調整
# ============================================================================= 