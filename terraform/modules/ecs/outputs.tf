# =============================================================================
# Ultimate Security Intelligence Platform - ECS Module Outputs
# ECS 模組的輸出值，供其他模組使用
# =============================================================================

# -----------------------------------------------------------------------------
# ECS Cluster 輸出
# -----------------------------------------------------------------------------

output "ecs_cluster_id" {
  description = "ECS 叢集 ID"
  value       = aws_ecs_cluster.main.id
}

output "ecs_cluster_name" {
  description = "ECS 叢集名稱"
  value       = aws_ecs_cluster.main.name
}

output "ecs_cluster_arn" {
  description = "ECS 叢集 ARN"
  value       = aws_ecs_cluster.main.arn
}

# -----------------------------------------------------------------------------
# ECS Service 輸出
# -----------------------------------------------------------------------------

output "ecs_service_id" {
  description = "ECS 服務 ID"
  value       = aws_ecs_service.app.id
}

output "ecs_service_name" {
  description = "ECS 服務名稱"
  value       = aws_ecs_service.app.name
}

output "ecs_service_arn" {
  description = "ECS 服務 ARN"
  value       = aws_ecs_service.app.id
}

output "ecs_service_desired_count" {
  description = "ECS 服務所需任務數量"
  value       = aws_ecs_service.app.desired_count
}

# -----------------------------------------------------------------------------
# ECS Task Definition 輸出
# -----------------------------------------------------------------------------

output "ecs_task_definition_arn" {
  description = "ECS 任務定義 ARN"
  value       = aws_ecs_task_definition.app.arn
}

output "ecs_task_definition_family" {
  description = "ECS 任務定義 Family"
  value       = aws_ecs_task_definition.app.family
}

output "ecs_task_definition_revision" {
  description = "ECS 任務定義 Revision"
  value       = aws_ecs_task_definition.app.revision
}

# -----------------------------------------------------------------------------
# Application Load Balancer 輸出
# -----------------------------------------------------------------------------

output "alb_id" {
  description = "Application Load Balancer ID"
  value       = aws_lb.main.id
}

output "alb_arn" {
  description = "Application Load Balancer ARN"
  value       = aws_lb.main.arn
}

output "alb_dns_name" {
  description = "Application Load Balancer DNS 名稱"
  value       = aws_lb.main.dns_name
}

output "alb_zone_id" {
  description = "Application Load Balancer Zone ID"
  value       = aws_lb.main.zone_id
}

output "alb_hosted_zone_id" {
  description = "Application Load Balancer Hosted Zone ID"
  value       = aws_lb.main.zone_id
}

# -----------------------------------------------------------------------------
# ALB Target Group 輸出
# -----------------------------------------------------------------------------

output "alb_target_group_id" {
  description = "ALB Target Group ID"
  value       = aws_lb_target_group.app.id
}

output "alb_target_group_arn" {
  description = "ALB Target Group ARN"
  value       = aws_lb_target_group.app.arn
}

output "alb_target_group_name" {
  description = "ALB Target Group 名稱"
  value       = aws_lb_target_group.app.name
}

# -----------------------------------------------------------------------------
# ALB Listener 輸出
# -----------------------------------------------------------------------------

output "alb_listener_id" {
  description = "ALB Listener ID"
  value       = aws_lb_listener.app.id
}

output "alb_listener_arn" {
  description = "ALB Listener ARN"
  value       = aws_lb_listener.app.arn
}

# -----------------------------------------------------------------------------
# 安全群組輸出
# -----------------------------------------------------------------------------

output "alb_security_group_id" {
  description = "ALB 安全群組 ID"
  value       = aws_security_group.alb.id
}

output "alb_security_group_arn" {
  description = "ALB 安全群組 ARN"
  value       = aws_security_group.alb.arn
}

output "ecs_task_security_group_id" {
  description = "ECS 任務安全群組 ID"
  value       = aws_security_group.ecs_task.id
}

output "ecs_task_security_group_arn" {
  description = "ECS 任務安全群組 ARN"
  value       = aws_security_group.ecs_task.arn
}

# -----------------------------------------------------------------------------
# IAM Role 輸出
# -----------------------------------------------------------------------------

output "ecs_task_execution_role_arn" {
  description = "ECS 任務執行角色 ARN"
  value       = aws_iam_role.ecs_task_execution_role.arn
}

output "ecs_task_execution_role_name" {
  description = "ECS 任務執行角色名稱"
  value       = aws_iam_role.ecs_task_execution_role.name
}

output "ecs_task_role_arn" {
  description = "ECS 任務角色 ARN"
  value       = aws_iam_role.ecs_task_role.arn
}

output "ecs_task_role_name" {
  description = "ECS 任務角色名稱"
  value       = aws_iam_role.ecs_task_role.name
}

# -----------------------------------------------------------------------------
# CloudWatch 輸出
# -----------------------------------------------------------------------------

output "cloudwatch_log_group_name" {
  description = "ECS CloudWatch Log Group 名稱"
  value       = aws_cloudwatch_log_group.ecs.name
}

output "cloudwatch_log_group_arn" {
  description = "ECS CloudWatch Log Group ARN"
  value       = aws_cloudwatch_log_group.ecs.arn
}

# -----------------------------------------------------------------------------
# 自動擴展輸出
# -----------------------------------------------------------------------------

output "auto_scaling_target_resource_id" {
  description = "自動擴展目標資源 ID"
  value       = var.enable_auto_scaling ? aws_appautoscaling_target.ecs_target[0].resource_id : null
}

output "auto_scaling_cpu_policy_arn" {
  description = "CPU 自動擴展政策 ARN"
  value       = var.enable_auto_scaling ? aws_appautoscaling_policy.ecs_policy_cpu[0].arn : null
}

output "auto_scaling_memory_policy_arn" {
  description = "記憶體自動擴展政策 ARN"
  value       = var.enable_auto_scaling ? aws_appautoscaling_policy.ecs_policy_memory[0].arn : null
}

# -----------------------------------------------------------------------------
# 服務發現輸出
# -----------------------------------------------------------------------------

output "service_discovery_namespace_id" {
  description = "服務發現命名空間 ID"
  value       = var.enable_service_discovery ? aws_service_discovery_private_dns_namespace.main[0].id : null
}

output "service_discovery_namespace_name" {
  description = "服務發現命名空間名稱"
  value       = var.enable_service_discovery ? aws_service_discovery_private_dns_namespace.main[0].name : null
}

output "service_discovery_service_id" {
  description = "服務發現服務 ID"
  value       = var.enable_service_discovery ? aws_service_discovery_service.app[0].id : null
}

output "service_discovery_service_name" {
  description = "服務發現服務名稱"
  value       = var.enable_service_discovery ? aws_service_discovery_service.app[0].name : null
}

# -----------------------------------------------------------------------------
# CloudWatch 警報輸出
# -----------------------------------------------------------------------------

output "cloudwatch_alarm_cpu_id" {
  description = "CPU 使用率警報 ID"
  value       = aws_cloudwatch_metric_alarm.ecs_cpu_high.id
}

output "cloudwatch_alarm_memory_id" {
  description = "記憶體使用率警報 ID"
  value       = aws_cloudwatch_metric_alarm.ecs_memory_high.id
}

output "cloudwatch_alarm_healthy_hosts_id" {
  description = "健康主機數量警報 ID"
  value       = aws_cloudwatch_metric_alarm.alb_healthy_hosts.id
}

# -----------------------------------------------------------------------------
# 應用程式存取資訊
# -----------------------------------------------------------------------------

output "application_url" {
  description = "應用程式存取 URL"
  value       = "http://${aws_lb.main.dns_name}"
}

output "application_health_check_url" {
  description = "應用程式健康檢查 URL"
  value       = "http://${aws_lb.main.dns_name}${var.alb_health_check_path}"
}

# -----------------------------------------------------------------------------
# 完整的 ECS 配置摘要
# -----------------------------------------------------------------------------

output "ecs_summary" {
  description = "ECS 配置摘要"
  value = {
    cluster_name              = aws_ecs_cluster.main.name
    service_name              = aws_ecs_service.app.name
    task_definition_family    = aws_ecs_task_definition.app.family
    desired_count             = aws_ecs_service.app.desired_count
    task_cpu                  = var.task_cpu
    task_memory               = var.task_memory
    container_port            = var.container_port
    alb_dns_name              = aws_lb.main.dns_name
    application_url           = "http://${aws_lb.main.dns_name}"
    auto_scaling_enabled      = var.enable_auto_scaling
    service_discovery_enabled = var.enable_service_discovery
    min_capacity              = var.enable_auto_scaling ? var.min_capacity : null
    max_capacity              = var.enable_auto_scaling ? var.max_capacity : null
  }
}

# =============================================================================
# 輸出說明
# =============================================================================
# 1. ALB DNS 名稱可用於配置 DNS 記錄或 CDN
# 2. 安全群組 ID 可供其他資源使用以建立網路規則
# 3. IAM 角色 ARN 可用於其他服務的權限配置
# 4. 自動擴展資源 ID 可用於進一步的擴展政策配置
# 5. 服務發現資訊可用於微服務間通信
# 6. application_url 提供了直接的應用程式存取點
# ============================================================================= 