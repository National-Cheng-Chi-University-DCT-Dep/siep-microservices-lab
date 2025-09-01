# =============================================================================
# Ultimate Security Intelligence Platform - RDS Module Outputs
# RDS 模組的輸出值，供其他模組使用
# =============================================================================

# -----------------------------------------------------------------------------
# 資料庫實例輸出
# -----------------------------------------------------------------------------

output "db_instance_id" {
  description = "RDS 實例 ID"
  value       = aws_db_instance.main.id
}

output "db_instance_arn" {
  description = "RDS 實例 ARN"
  value       = aws_db_instance.main.arn
}

output "db_instance_endpoint" {
  description = "RDS 實例連接端點"
  value       = aws_db_instance.main.endpoint
}

output "db_instance_hosted_zone_id" {
  description = "RDS 實例 Hosted Zone ID"
  value       = aws_db_instance.main.hosted_zone_id
}

output "db_instance_port" {
  description = "資料庫連接埠"
  value       = aws_db_instance.main.port
}

output "db_instance_status" {
  description = "RDS 實例狀態"
  value       = aws_db_instance.main.status
}

# -----------------------------------------------------------------------------
# 資料庫認證輸出
# -----------------------------------------------------------------------------

output "db_instance_username" {
  description = "資料庫主使用者名稱"
  value       = aws_db_instance.main.username
  sensitive   = true
}

output "db_instance_database_name" {
  description = "預設資料庫名稱"
  value       = aws_db_instance.main.db_name
}

# -----------------------------------------------------------------------------
# Secrets Manager 輸出
# -----------------------------------------------------------------------------

output "db_credentials_secret_arn" {
  description = "資料庫認證 Secret ARN"
  value       = aws_secretsmanager_secret.db_credentials.arn
}

output "db_credentials_secret_name" {
  description = "資料庫認證 Secret 名稱"
  value       = aws_secretsmanager_secret.db_credentials.name
}

output "db_credentials_secret_id" {
  description = "資料庫認證 Secret ID"
  value       = aws_secretsmanager_secret.db_credentials.id
}

# -----------------------------------------------------------------------------
# 安全性輸出
# -----------------------------------------------------------------------------

output "db_security_group_id" {
  description = "RDS 安全群組 ID"
  value       = aws_security_group.rds.id
}

output "db_security_group_arn" {
  description = "RDS 安全群組 ARN"
  value       = aws_security_group.rds.arn
}

# -----------------------------------------------------------------------------
# 子網路群組輸出
# -----------------------------------------------------------------------------

output "db_subnet_group_id" {
  description = "DB 子網路群組 ID"
  value       = aws_db_subnet_group.main.id
}

output "db_subnet_group_arn" {
  description = "DB 子網路群組 ARN"
  value       = aws_db_subnet_group.main.arn
}

# -----------------------------------------------------------------------------
# 參數群組輸出
# -----------------------------------------------------------------------------

output "db_parameter_group_id" {
  description = "DB 參數群組 ID"
  value       = aws_db_parameter_group.main.id
}

output "db_parameter_group_arn" {
  description = "DB 參數群組 ARN"
  value       = aws_db_parameter_group.main.arn
}

# -----------------------------------------------------------------------------
# 加密金鑰輸出
# -----------------------------------------------------------------------------

output "kms_key_id" {
  description = "RDS 加密 KMS 金鑰 ID"
  value       = aws_kms_key.rds.key_id
}

output "kms_key_arn" {
  description = "RDS 加密 KMS 金鑰 ARN"
  value       = aws_kms_key.rds.arn
}

output "kms_alias_name" {
  description = "RDS 加密 KMS 金鑰別名"
  value       = aws_kms_alias.rds.name
}

# -----------------------------------------------------------------------------
# 監控輸出
# -----------------------------------------------------------------------------

output "cloudwatch_log_group_name" {
  description = "PostgreSQL CloudWatch Log Group 名稱"
  value       = aws_cloudwatch_log_group.postgresql.name
}

output "cloudwatch_log_group_arn" {
  description = "PostgreSQL CloudWatch Log Group ARN"
  value       = aws_cloudwatch_log_group.postgresql.arn
}

output "rds_monitoring_role_arn" {
  description = "RDS 監控 IAM Role ARN"
  value       = aws_iam_role.rds_monitoring.arn
}

# -----------------------------------------------------------------------------
# 讀取副本輸出（如果啟用）
# -----------------------------------------------------------------------------

output "read_replica_id" {
  description = "讀取副本實例 ID"
  value       = var.create_read_replica ? aws_db_instance.read_replica[0].id : null
}

output "read_replica_endpoint" {
  description = "讀取副本連接端點"
  value       = var.create_read_replica ? aws_db_instance.read_replica[0].endpoint : null
}

# -----------------------------------------------------------------------------
# 快照輸出（如果建立）
# -----------------------------------------------------------------------------

output "manual_snapshot_id" {
  description = "手動快照 ID"
  value       = var.create_manual_snapshot ? aws_db_snapshot.manual_snapshot[0].id : null
}

# -----------------------------------------------------------------------------
# CloudWatch 警報輸出
# -----------------------------------------------------------------------------

output "cloudwatch_alarm_cpu_id" {
  description = "CPU 使用率警報 ID"
  value       = aws_cloudwatch_metric_alarm.database_cpu.id
}

output "cloudwatch_alarm_connections_id" {
  description = "連接數警報 ID"
  value       = aws_cloudwatch_metric_alarm.database_connections.id
}

output "cloudwatch_alarm_storage_id" {
  description = "可用儲存空間警報 ID"
  value       = aws_cloudwatch_metric_alarm.database_free_storage.id
}

# -----------------------------------------------------------------------------
# 連接字串輸出（用於應用程式配置）
# -----------------------------------------------------------------------------

output "connection_info" {
  description = "資料庫連接資訊"
  value = {
    host     = aws_db_instance.main.endpoint
    port     = aws_db_instance.main.port
    database = aws_db_instance.main.db_name
    username = aws_db_instance.main.username
    # 注意：密碼不在此輸出，應從 Secrets Manager 獲取
  }
  sensitive = true
}

# -----------------------------------------------------------------------------
# 完整的資料庫配置摘要
# -----------------------------------------------------------------------------

output "database_summary" {
  description = "資料庫配置摘要"
  value = {
    instance_id       = aws_db_instance.main.id
    engine            = aws_db_instance.main.engine
    engine_version    = aws_db_instance.main.engine_version
    instance_class    = aws_db_instance.main.instance_class
    allocated_storage = aws_db_instance.main.allocated_storage
    storage_encrypted = aws_db_instance.main.storage_encrypted
    multi_az         = aws_db_instance.main.multi_az
    backup_retention = aws_db_instance.main.backup_retention_period
    maintenance_window = aws_db_instance.main.maintenance_window
    backup_window     = aws_db_instance.main.backup_window
    monitoring_enabled = aws_db_instance.main.monitoring_interval > 0
    performance_insights = aws_db_instance.main.performance_insights_enabled
    deletion_protection = aws_db_instance.main.deletion_protection
    read_replica_count = var.create_read_replica ? 1 : 0
  }
}

# =============================================================================
# 輸出說明
# =============================================================================
# 1. 敏感資訊（如密碼）標記為 sensitive，不會在 plan/apply 輸出中顯示
# 2. 連接資訊可用於配置應用程式的環境變數
# 3. Secrets Manager ARN 可用於應用程式動態獲取認證
# 4. 安全群組 ID 可供其他資源使用以建立網路規則
# 5. summary 提供了完整的配置概覽
# ============================================================================= 