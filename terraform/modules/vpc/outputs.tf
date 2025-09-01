# =============================================================================
# Ultimate Security Intelligence Platform - VPC Module Outputs
# VPC 模組的輸出值，供其他模組使用
# =============================================================================

# -----------------------------------------------------------------------------
# VPC 輸出
# -----------------------------------------------------------------------------

output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "vpc_cidr_block" {
  description = "VPC CIDR 區塊"
  value       = aws_vpc.main.cidr_block
}

output "vpc_arn" {
  description = "VPC ARN"
  value       = aws_vpc.main.arn
}

# -----------------------------------------------------------------------------
# 子網路輸出
# -----------------------------------------------------------------------------

output "public_subnet_ids" {
  description = "公開子網路 ID 列表"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "私有子網路 ID 列表"
  value       = aws_subnet.private[*].id
}

output "public_subnet_cidrs" {
  description = "公開子網路 CIDR 區塊列表"
  value       = aws_subnet.public[*].cidr_block
}

output "private_subnet_cidrs" {
  description = "私有子網路 CIDR 區塊列表"
  value       = aws_subnet.private[*].cidr_block
}

output "availability_zones" {
  description = "使用的可用性區域列表"
  value       = aws_subnet.public[*].availability_zone
}

# -----------------------------------------------------------------------------
# 網路閘道輸出
# -----------------------------------------------------------------------------

output "internet_gateway_id" {
  description = "Internet Gateway ID"
  value       = aws_internet_gateway.main.id
}

output "nat_gateway_ids" {
  description = "NAT Gateway ID 列表"
  value       = aws_nat_gateway.main[*].id
}

output "nat_gateway_public_ips" {
  description = "NAT Gateway 公開 IP 列表"
  value       = aws_eip.nat[*].public_ip
}

# -----------------------------------------------------------------------------
# 路由表輸出
# -----------------------------------------------------------------------------

output "public_route_table_id" {
  description = "公開路由表 ID"
  value       = aws_route_table.public.id
}

output "private_route_table_ids" {
  description = "私有路由表 ID 列表"
  value       = aws_route_table.private[*].id
}

# -----------------------------------------------------------------------------
# VPC Endpoints 輸出
# -----------------------------------------------------------------------------

output "s3_vpc_endpoint_id" {
  description = "S3 VPC Endpoint ID"
  value       = var.enable_s3_endpoint ? aws_vpc_endpoint.s3[0].id : null
}

# -----------------------------------------------------------------------------
# 安全性相關輸出
# -----------------------------------------------------------------------------

output "default_security_group_id" {
  description = "預設安全群組 ID"
  value       = aws_vpc.main.default_security_group_id
}

output "vpc_flow_logs_log_group_name" {
  description = "VPC Flow Logs CloudWatch Log Group 名稱"
  value       = var.enable_vpc_flow_logs ? aws_cloudwatch_log_group.vpc_flow_logs.name : null
}

output "vpc_flow_logs_iam_role_arn" {
  description = "VPC Flow Logs IAM Role ARN"
  value       = aws_iam_role.vpc_flow_logs.arn
}

# -----------------------------------------------------------------------------
# Network ACLs 輸出
# -----------------------------------------------------------------------------

output "public_network_acl_id" {
  description = "公開網路 ACL ID"
  value       = var.enable_network_acls ? aws_network_acl.public[0].id : null
}

output "private_network_acl_id" {
  description = "私有網路 ACL ID"
  value       = var.enable_network_acls ? aws_network_acl.private[0].id : null
}

# -----------------------------------------------------------------------------
# 實用的組合輸出
# -----------------------------------------------------------------------------

output "public_subnet_route_table_associations" {
  description = "公開子網路路由表關聯"
  value = {
    for idx, subnet_id in aws_subnet.public[*].id :
    subnet_id => aws_route_table.public.id
  }
}

output "private_subnet_route_table_associations" {
  description = "私有子網路路由表關聯"
  value = {
    for idx, subnet_id in aws_subnet.private[*].id :
    subnet_id => aws_route_table.private[idx].id
  }
}

output "subnet_groups" {
  description = "子網路分組資訊"
  value = {
    public = {
      subnet_ids = aws_subnet.public[*].id
      cidrs      = aws_subnet.public[*].cidr_block
      azs        = aws_subnet.public[*].availability_zone
    }
    private = {
      subnet_ids = aws_subnet.private[*].id
      cidrs      = aws_subnet.private[*].cidr_block
      azs        = aws_subnet.private[*].availability_zone
    }
  }
}

# -----------------------------------------------------------------------------
# 網路配置摘要
# -----------------------------------------------------------------------------

output "network_summary" {
  description = "網路配置摘要"
  value = {
    vpc_id               = aws_vpc.main.id
    vpc_cidr             = aws_vpc.main.cidr_block
    public_subnets       = length(aws_subnet.public)
    private_subnets      = length(aws_subnet.private)
    availability_zones   = length(var.availability_zones)
    nat_gateways        = length(aws_nat_gateway.main)
    internet_gateway    = aws_internet_gateway.main.id
    vpc_flow_logs       = var.enable_vpc_flow_logs
    s3_endpoint         = var.enable_s3_endpoint
    network_acls        = var.enable_network_acls
  }
}

# =============================================================================
# 輸出說明
# =============================================================================
# 1. 這些輸出值將被其他模組（如 RDS、ECS）使用
# 2. 輸出值提供了完整的網路配置資訊
# 3. 組合輸出（如 subnet_groups）便於其他模組使用
# 4. network_summary 提供了快速的配置概覽
# ============================================================================= 