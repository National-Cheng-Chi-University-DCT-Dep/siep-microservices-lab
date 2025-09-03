# =============================================================================
# Ultimate Security Intelligence Platform - Security Module
# 建立安全性相關資源：IAM 角色、安全群組、Secrets Manager 等
# =============================================================================

# -----------------------------------------------------------------------------
# Local Variables
# -----------------------------------------------------------------------------

locals {
  common_tags = merge(var.common_tags, {
    Module = "security"
  })
}

# -----------------------------------------------------------------------------
# KMS Keys for Encryption
# -----------------------------------------------------------------------------

resource "aws_kms_key" "main" {
  description             = "KMS key for ${var.project_name} ${var.environment}"
  deletion_window_in_days = var.kms_deletion_window_in_days
  enable_key_rotation     = true

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "Enable IAM User Permissions"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        }
        Action   = "kms:*"
        Resource = "*"
      },
      {
        Sid    = "Allow ECS Task Role"
        Effect = "Allow"
        Principal = {
          AWS = var.ecs_task_role_arn
        }
        Action = [
          "kms:Decrypt",
          "kms:DescribeKey"
        ]
        Resource = "*"
      }
    ]
  })

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-kms-key"
  })
}

resource "aws_kms_alias" "main" {
  name          = "alias/${var.project_name}-${var.environment}"
  target_key_id = aws_kms_key.main.key_id
}

# -----------------------------------------------------------------------------
# Data Sources
# -----------------------------------------------------------------------------

data "aws_caller_identity" "current" {}

# -----------------------------------------------------------------------------
# Secrets Manager Secrets
# -----------------------------------------------------------------------------

resource "aws_secretsmanager_secret" "app_secrets" {
  name        = "${var.project_name}-${var.environment}-app-secrets"
  description = "Application secrets for ${var.project_name} ${var.environment}"
  kms_key_id  = aws_kms_key.main.key_id

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-app-secrets"
  })
}

resource "aws_secretsmanager_secret_version" "app_secrets" {
  secret_id = aws_secretsmanager_secret.app_secrets.id
  secret_string = jsonencode({
    jwt_secret        = var.jwt_secret
    encryption_key    = var.encryption_key
    api_key          = var.api_key
    third_party_tokens = var.third_party_tokens
  })
}

# -----------------------------------------------------------------------------
# IAM Roles for Application Services
# -----------------------------------------------------------------------------

# ECS Task Execution Role (如果未提供)
resource "aws_iam_role" "ecs_task_execution_role" {
  count = var.create_ecs_task_execution_role ? 1 : 0

  name = "${var.project_name}-${var.environment}-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  count = var.create_ecs_task_execution_role ? 1 : 0

  role       = aws_iam_role.ecs_task_execution_role[0].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ECS Task Role (如果未提供)
resource "aws_iam_role" "ecs_task_role" {
  count = var.create_ecs_task_role ? 1 : 0

  name = "${var.project_name}-${var.environment}-ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = local.common_tags
}

# -----------------------------------------------------------------------------
# IAM Policies for Application Services
# -----------------------------------------------------------------------------

resource "aws_iam_policy" "secrets_access_policy" {
  name        = "${var.project_name}-${var.environment}-secrets-access"
  description = "Policy for accessing secrets manager"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ]
        Resource = [
          aws_secretsmanager_secret.app_secrets.arn,
          "${aws_secretsmanager_secret.app_secrets.arn}:*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "kms:Decrypt",
          "kms:DescribeKey"
        ]
        Resource = aws_kms_key.main.arn
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_policy" "s3_access_policy" {
  name        = "${var.project_name}-${var.environment}-s3-access"
  description = "Policy for accessing S3 buckets"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = [
          for bucket in var.s3_bucket_arns : "${bucket}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:ListBucket"
        ]
        Resource = var.s3_bucket_arns
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_policy" "cloudwatch_logs_policy" {
  name        = "${var.project_name}-${var.environment}-cloudwatch-logs"
  description = "Policy for CloudWatch Logs access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams"
        ]
        Resource = "arn:aws:logs:*:*:*"
      }
    ]
  })

  tags = local.common_tags
}

# -----------------------------------------------------------------------------
# Attach Policies to ECS Task Role
# -----------------------------------------------------------------------------

resource "aws_iam_role_policy_attachment" "ecs_task_secrets_policy" {
  count = var.create_ecs_task_role ? 1 : 0

  role       = aws_iam_role.ecs_task_role[0].name
  policy_arn = aws_iam_policy.secrets_access_policy.arn
}

resource "aws_iam_role_policy_attachment" "ecs_task_s3_policy" {
  count = var.create_ecs_task_role && length(var.s3_bucket_arns) > 0 ? 1 : 0

  role       = aws_iam_role.ecs_task_role[0].name
  policy_arn = aws_iam_policy.s3_access_policy.arn
}

resource "aws_iam_role_policy_attachment" "ecs_task_cloudwatch_policy" {
  count = var.create_ecs_task_role ? 1 : 0

  role       = aws_iam_role.ecs_task_role[0].name
  policy_arn = aws_iam_policy.cloudwatch_logs_policy.arn
}

# -----------------------------------------------------------------------------
# Security Groups
# -----------------------------------------------------------------------------

resource "aws_security_group" "database_client" {
  name_prefix = "${var.project_name}-${var.environment}-db-client-"
  vpc_id      = var.vpc_id
  description = "Security group for database clients"

  # 允許出站到資料庫
  egress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
    description = "PostgreSQL access to database"
  }

  # 允許出站到 Redis
  egress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
    description = "Redis access"
  }

  # 允許 HTTPS 出站
  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS outbound"
  }

  # 允許 HTTP 出站
  egress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP outbound"
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-db-client-sg"
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "cache_client" {
  name_prefix = "${var.project_name}-${var.environment}-cache-client-"
  vpc_id      = var.vpc_id
  description = "Security group for cache clients"

  # 允許出站到 Redis
  egress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
    description = "Redis access"
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-cache-client-sg"
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "internal_services" {
  name_prefix = "${var.project_name}-${var.environment}-internal-"
  vpc_id      = var.vpc_id
  description = "Security group for internal services communication"

  # 允許 VPC 內部通信
  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
    description = "Internal VPC communication"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-internal-sg"
  })

  lifecycle {
    create_before_destroy = true
  }
}

# -----------------------------------------------------------------------------
# WAF (Web Application Firewall) for ALB
# -----------------------------------------------------------------------------

resource "aws_wafv2_web_acl" "main" {
  count = var.enable_waf ? 1 : 0

  name  = "${var.project_name}-${var.environment}-waf"
  scope = "REGIONAL"

  default_action {
    allow {}
  }

  # Rate limiting rule
  rule {
    name     = "RateLimitRule"
    priority = 1

    action {
      block {}
    }

    statement {
      rate_based_statement {
        limit              = var.waf_rate_limit
        aggregate_key_type = "IP"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "RateLimitRule"
      sampled_requests_enabled   = true
    }
  }

  # AWS Managed Rules
  rule {
    name     = "AWSManagedRulesCommonRuleSet"
    priority = 2

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "CommonRuleSetMetric"
      sampled_requests_enabled   = true
    }
  }

  # SQL Injection Protection
  rule {
    name     = "AWSManagedRulesSQLiRuleSet"
    priority = 3

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesSQLiRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "SQLiRuleSetMetric"
      sampled_requests_enabled   = true
    }
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-waf"
  })
}

# -----------------------------------------------------------------------------
# CloudWatch Log Group for WAF
# -----------------------------------------------------------------------------

resource "aws_cloudwatch_log_group" "waf_logs" {
  count = var.enable_waf ? 1 : 0

  name              = "/aws/wafv2/${var.project_name}
  kms_key_id = aws_kms_key.main.arn

  kms_key_id = aws_kms_key.main.arn
-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-waf-logs"
  })
}

# -----------------------------------------------------------------------------
# Certificate Manager for SSL/TLS
# -----------------------------------------------------------------------------

resource "aws_acm_certificate" "main" {
  count = var.enable_ssl && var.domain_name != null ? 1 : 0

  domain_name       = var.domain_name
  validation_method = "DNS"

  subject_alternative_names = var.subject_alternative_names

  lifecycle {
    create_before_destroy = true
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-cert"
  })
}

# -----------------------------------------------------------------------------
# SNS Topic for Security Alerts
# -----------------------------------------------------------------------------

resource "aws_sns_topic" "security_alerts" {
  count = var.enable_security_alerts ? 1 : 0

  name         = "${var.project_name}-${var.environment}-security-alerts"
  display_name = "Security Alerts for ${var.project_name} ${var.environment}"

  kms_master_key_id = aws_kms_key.main.key_id

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-security-alerts"
  })
}

resource "aws_sns_topic_subscription" "security_alerts_email" {
  count = var.enable_security_alerts && var.security_alert_email != null ? 1 : 0

  topic_arn = aws_sns_topic.security_alerts[0].arn
  protocol  = "email"
  endpoint  = var.security_alert_email
}

# -----------------------------------------------------------------------------
# CloudWatch Alarms for Security Monitoring
# -----------------------------------------------------------------------------

resource "aws_cloudwatch_metric_alarm" "waf_blocked_requests" {
  count = var.enable_waf ? 1 : 0

  alarm_name          = "${var.project_name}-${var.environment}-waf-blocked-requests"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "BlockedRequests"
  namespace           = "AWS/WAFV2"
  period              = "300"
  statistic           = "Sum"
  threshold           = "100"
  alarm_description   = "This metric monitors WAF blocked requests"
  alarm_actions       = var.enable_security_alerts ? [aws_sns_topic.security_alerts[0].arn] : []

  dimensions = {
    WebACL = aws_wafv2_web_acl.main[0].name
    Region = var.aws_region
  }

  tags = local.common_tags
}

resource "aws_cloudwatch_metric_alarm" "failed_logins" {
  count = var.enable_security_monitoring ? 1 : 0

  alarm_name          = "${var.project_name}-${var.environment}-failed-logins"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "FailedLogins"
  namespace           = "CustomApp"
  period              = "300"
  statistic           = "Sum"
  threshold           = "50"
  alarm_description   = "This metric monitors failed login attempts"
  alarm_actions       = var.enable_security_alerts ? [aws_sns_topic.security_alerts[0].arn] : []

  tags = local.common_tags
}

# -----------------------------------------------------------------------------
# VPC Flow Logs for Network Monitoring
# -----------------------------------------------------------------------------

resource "aws_flow_log" "vpc_flow_logs" {
  count = var.enable_vpc_flow_logs ? 1 : 0

  iam_role_arn    = aws_iam_role.flow_logs_role[0].arn
  log_destination = aws_cloudwatch_log_group.vpc_flow_logs[0].arn
  traffic_type    = "ALL"
  vpc_id          = var.vpc_id

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-vpc-flow-logs"
  })
}

resource "aws_cloudwatch_log_group" "vpc_flow_logs" {
  count = var.enable_vpc_flow_logs ? 1 : 0

  name              = "/aws/vpc/flowlogs/${var.project_name}
  kms_key_id = aws_kms_key.main.arn

  kms_key_id = aws_kms_key.main.arn
-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-vpc-flow-logs"
  })
}

resource "aws_iam_role" "flow_logs_role" {
  count = var.enable_vpc_flow_logs ? 1 : 0

  name = "${var.project_name}-${var.environment}-flow-logs-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "vpc-flow-logs.amazonaws.com"
        }
      }
    ]
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy" "flow_logs_policy" {
  count = var.enable_vpc_flow_logs ? 1 : 0

  name = "${var.project_name}-${var.environment}-flow-logs-policy"
  role = aws_iam_role.flow_logs_role[0].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams"
        ]
        Resource = "*"
      }
    ]
  })
} 