# =============================================================================
# Ultimate Security Intelligence Platform - Security Module Outputs
# 安全性模組的輸出值，供其他模組使用
# =============================================================================

# -----------------------------------------------------------------------------
# KMS 輸出
# -----------------------------------------------------------------------------

output "kms_key_id" {
  description = "KMS 金鑰 ID"
  value       = aws_kms_key.main.key_id
}

output "kms_key_arn" {
  description = "KMS 金鑰 ARN"
  value       = aws_kms_key.main.arn
}

output "kms_alias_name" {
  description = "KMS 金鑰別名"
  value       = aws_kms_alias.main.name
}

# -----------------------------------------------------------------------------
# Secrets Manager 輸出
# -----------------------------------------------------------------------------

output "app_secrets_arn" {
  description = "應用程式 Secrets Manager ARN"
  value       = aws_secretsmanager_secret.app_secrets.arn
}

output "app_secrets_name" {
  description = "應用程式 Secrets Manager 名稱"
  value       = aws_secretsmanager_secret.app_secrets.name
}

output "app_secrets_version_id" {
  description = "應用程式 Secrets 版本 ID"
  value       = aws_secretsmanager_secret_version.app_secrets.version_id
}

# -----------------------------------------------------------------------------
# IAM 角色輸出
# -----------------------------------------------------------------------------

output "ecs_task_execution_role_arn" {
  description = "ECS 任務執行角色 ARN"
  value       = var.create_ecs_task_execution_role ? aws_iam_role.ecs_task_execution_role[0].arn : null
}

output "ecs_task_execution_role_name" {
  description = "ECS 任務執行角色名稱"
  value       = var.create_ecs_task_execution_role ? aws_iam_role.ecs_task_execution_role[0].name : null
}

output "ecs_task_role_arn" {
  description = "ECS 任務角色 ARN"
  value       = var.create_ecs_task_role ? aws_iam_role.ecs_task_role[0].arn : null
}

output "ecs_task_role_name" {
  description = "ECS 任務角色名稱"
  value       = var.create_ecs_task_role ? aws_iam_role.ecs_task_role[0].name : null
}

# -----------------------------------------------------------------------------
# IAM 政策輸出
# -----------------------------------------------------------------------------

output "secrets_access_policy_arn" {
  description = "Secrets Manager 存取政策 ARN"
  value       = aws_iam_policy.secrets_access_policy.arn
}

output "s3_access_policy_arn" {
  description = "S3 存取政策 ARN"
  value       = aws_iam_policy.s3_access_policy.arn
}

output "cloudwatch_logs_policy_arn" {
  description = "CloudWatch Logs 政策 ARN"
  value       = aws_iam_policy.cloudwatch_logs_policy.arn
}

# -----------------------------------------------------------------------------
# 安全群組輸出
# -----------------------------------------------------------------------------

output "database_client_security_group_id" {
  description = "資料庫客戶端安全群組 ID"
  value       = aws_security_group.database_client.id
}

output "cache_client_security_group_id" {
  description = "快取客戶端安全群組 ID"
  value       = aws_security_group.cache_client.id
}

output "internal_services_security_group_id" {
  description = "內部服務安全群組 ID"
  value       = aws_security_group.internal_services.id
}

# -----------------------------------------------------------------------------
# WAF 輸出
# -----------------------------------------------------------------------------

output "waf_web_acl_id" {
  description = "WAF Web ACL ID"
  value       = var.enable_waf ? aws_wafv2_web_acl.main[0].id : null
}

output "waf_web_acl_arn" {
  description = "WAF Web ACL ARN"
  value       = var.enable_waf ? aws_wafv2_web_acl.main[0].arn : null
}

output "waf_web_acl_name" {
  description = "WAF Web ACL 名稱"
  value       = var.enable_waf ? aws_wafv2_web_acl.main[0].name : null
}

# -----------------------------------------------------------------------------
# SSL/TLS 憑證輸出
# -----------------------------------------------------------------------------

output "ssl_certificate_arn" {
  description = "SSL 憑證 ARN"
  value       = var.enable_ssl && var.domain_name != null ? aws_acm_certificate.main[0].arn : null
}

output "ssl_certificate_domain_name" {
  description = "SSL 憑證域名"
  value       = var.enable_ssl && var.domain_name != null ? aws_acm_certificate.main[0].domain_name : null
}

output "ssl_certificate_validation_options" {
  description = "SSL 憑證驗證選項"
  value       = var.enable_ssl && var.domain_name != null ? aws_acm_certificate.main[0].domain_validation_options : null
}

# -----------------------------------------------------------------------------
# SNS 主題輸出
# -----------------------------------------------------------------------------

output "security_alerts_topic_arn" {
  description = "安全性警報 SNS 主題 ARN"
  value       = var.enable_security_alerts ? aws_sns_topic.security_alerts[0].arn : null
}

output "security_alerts_topic_name" {
  description = "安全性警報 SNS 主題名稱"
  value       = var.enable_security_alerts ? aws_sns_topic.security_alerts[0].name : null
}

# -----------------------------------------------------------------------------
# CloudWatch 日誌群組輸出
# -----------------------------------------------------------------------------

output "waf_log_group_name" {
  description = "WAF CloudWatch 日誌群組名稱"
  value       = var.enable_waf ? aws_cloudwatch_log_group.waf_logs[0].name : null
}

output "vpc_flow_logs_group_name" {
  description = "VPC Flow Logs CloudWatch 日誌群組名稱"
  value       = var.enable_vpc_flow_logs ? aws_cloudwatch_log_group.vpc_flow_logs[0].name : null
}

# -----------------------------------------------------------------------------
# CloudWatch 警報輸出
# -----------------------------------------------------------------------------

output "waf_blocked_requests_alarm_id" {
  description = "WAF 封鎖請求警報 ID"
  value       = var.enable_waf ? aws_cloudwatch_metric_alarm.waf_blocked_requests[0].id : null
}

output "failed_logins_alarm_id" {
  description = "失敗登入警報 ID"
  value       = var.enable_security_monitoring ? aws_cloudwatch_metric_alarm.failed_logins[0].id : null
}

# -----------------------------------------------------------------------------
# VPC Flow Logs 輸出
# -----------------------------------------------------------------------------

output "vpc_flow_logs_id" {
  description = "VPC Flow Logs ID"
  value       = var.enable_vpc_flow_logs ? aws_flow_log.vpc_flow_logs[0].id : null
}

output "vpc_flow_logs_role_arn" {
  description = "VPC Flow Logs IAM 角色 ARN"
  value       = var.enable_vpc_flow_logs ? aws_iam_role.flow_logs_role[0].arn : null
}

# -----------------------------------------------------------------------------
# 安全性配置摘要
# -----------------------------------------------------------------------------

output "security_summary" {
  description = "安全性配置摘要"
  value = {
    kms_encryption_enabled      = true
    secrets_manager_enabled     = true
    waf_enabled                 = var.enable_waf
    ssl_enabled                 = var.enable_ssl
    security_alerts_enabled     = var.enable_security_alerts
    vpc_flow_logs_enabled       = var.enable_vpc_flow_logs
    security_monitoring_enabled = var.enable_security_monitoring
    iam_roles_created = {
      ecs_task_execution_role = var.create_ecs_task_execution_role
      ecs_task_role           = var.create_ecs_task_role
    }
    security_groups_count   = 3
    cloudwatch_alarms_count = (var.enable_waf ? 1 : 0) + (var.enable_security_monitoring ? 1 : 0)
  }
}

# -----------------------------------------------------------------------------
# 應用程式設定輸出（用於 ECS 任務定義）
# -----------------------------------------------------------------------------

output "ecs_secrets_configuration" {
  description = "ECS 任務定義的 Secrets 配置"
  value = [
    {
      name      = "JWT_SECRET"
      valueFrom = "${aws_secretsmanager_secret.app_secrets.arn}:jwt_secret::"
    },
    {
      name      = "ENCRYPTION_KEY"
      valueFrom = "${aws_secretsmanager_secret.app_secrets.arn}:encryption_key::"
    },
    {
      name      = "API_KEY"
      valueFrom = "${aws_secretsmanager_secret.app_secrets.arn}:api_key::"
    }
  ]
}

output "ecs_security_groups" {
  description = "ECS 任務建議使用的安全群組"
  value = [
    aws_security_group.database_client.id,
    aws_security_group.cache_client.id,
    aws_security_group.internal_services.id
  ]
}

# =============================================================================
# 輸出說明
# =============================================================================
# 1. KMS 金鑰用於加密所有敏感資料
# 2. Secrets Manager 安全地存儲應用程式機密
# 3. IAM 角色提供最小權限原則的存取控制
# 4. 安全群組提供網路層級的安全控制
# 5. WAF 提供應用程式層級的安全防護
# 6. CloudWatch 警報提供安全性事件監控
# 7. VPC Flow Logs 提供網路流量監控
# 8. ecs_secrets_configuration 可直接用於 ECS 任務定義
# ============================================================================= 