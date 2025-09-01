#!/bin/bash

# =============================================================================
# Ultimate Security Intelligence Platform - Status Check Script
# 開發環境狀態檢查腳本
# =============================================================================

# 顏色代碼
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日誌函數
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 檢查 Terraform 狀態
check_terraform_status() {
    log_info "檢查 Terraform 狀態..."
    
    if [ ! -d ".terraform" ]; then
        log_warning "Terraform 尚未初始化。請執行 'terraform init'。"
        return 1
    fi
    
    if terraform workspace show &>/dev/null; then
        WORKSPACE=$(terraform workspace show)
        log_success "Terraform 工作區: $WORKSPACE"
    else
        log_error "無法檢查 Terraform 工作區。"
        return 1
    fi
    
    return 0
}

# 檢查 AWS 資源狀態
check_aws_resources() {
    log_info "檢查 AWS 資源狀態..."
    
    # 檢查 VPC
    VPC_ID=$(terraform output -raw vpc_id 2>/dev/null || echo "")
    if [ ! -z "$VPC_ID" ]; then
        if aws ec2 describe-vpcs --vpc-ids "$VPC_ID" &>/dev/null; then
            log_success "VPC ($VPC_ID) 狀態正常"
        else
            log_error "VPC ($VPC_ID) 不存在或無法存取"
        fi
    else
        log_warning "無法獲取 VPC ID"
    fi
    
    # 檢查 RDS 實例
    DB_INSTANCE_ID=$(terraform output -raw ecs_cluster_name 2>/dev/null | sed 's/-cluster$//' | sed 's/$/-postgres/' || echo "")
    if [ ! -z "$DB_INSTANCE_ID" ]; then
        DB_STATUS=$(aws rds describe-db-instances --db-instance-identifier "$DB_INSTANCE_ID" --query 'DBInstances[0].DBInstanceStatus' --output text 2>/dev/null || echo "")
        if [ "$DB_STATUS" = "available" ]; then
            log_success "RDS 實例 ($DB_INSTANCE_ID) 狀態: $DB_STATUS"
        elif [ ! -z "$DB_STATUS" ]; then
            log_warning "RDS 實例 ($DB_INSTANCE_ID) 狀態: $DB_STATUS"
        else
            log_error "無法檢查 RDS 實例狀態"
        fi
    else
        log_warning "無法獲取 RDS 實例 ID"
    fi
    
    # 檢查 ECS 叢集
    CLUSTER_NAME=$(terraform output -raw ecs_cluster_name 2>/dev/null || echo "")
    if [ ! -z "$CLUSTER_NAME" ]; then
        CLUSTER_STATUS=$(aws ecs describe-clusters --clusters "$CLUSTER_NAME" --query 'clusters[0].status' --output text 2>/dev/null || echo "")
        if [ "$CLUSTER_STATUS" = "ACTIVE" ]; then
            log_success "ECS 叢集 ($CLUSTER_NAME) 狀態: $CLUSTER_STATUS"
            
            # 檢查 ECS 服務
            SERVICE_NAME=$(terraform output -raw ecs_service_name 2>/dev/null || echo "")
            if [ ! -z "$SERVICE_NAME" ]; then
                RUNNING_COUNT=$(aws ecs describe-services --cluster "$CLUSTER_NAME" --services "$SERVICE_NAME" --query 'services[0].runningCount' --output text 2>/dev/null || echo "0")
                DESIRED_COUNT=$(aws ecs describe-services --cluster "$CLUSTER_NAME" --services "$SERVICE_NAME" --query 'services[0].desiredCount' --output text 2>/dev/null || echo "0")
                
                if [ "$RUNNING_COUNT" = "$DESIRED_COUNT" ] && [ "$RUNNING_COUNT" != "0" ]; then
                    log_success "ECS 服務 ($SERVICE_NAME) 執行中: $RUNNING_COUNT/$DESIRED_COUNT 任務"
                else
                    log_warning "ECS 服務 ($SERVICE_NAME) 狀態異常: $RUNNING_COUNT/$DESIRED_COUNT 任務"
                fi
            fi
        elif [ ! -z "$CLUSTER_STATUS" ]; then
            log_warning "ECS 叢集 ($CLUSTER_NAME) 狀態: $CLUSTER_STATUS"
        else
            log_error "無法檢查 ECS 叢集狀態"
        fi
    else
        log_warning "無法獲取 ECS 叢集名稱"
    fi
    
    # 檢查 ALB
    ALB_DNS=$(terraform output -raw alb_dns_name 2>/dev/null || echo "")
    if [ ! -z "$ALB_DNS" ]; then
        if nslookup "$ALB_DNS" &>/dev/null; then
            log_success "ALB DNS ($ALB_DNS) 可解析"
        else
            log_warning "ALB DNS ($ALB_DNS) 無法解析"
        fi
    else
        log_warning "無法獲取 ALB DNS 名稱"
    fi
}

# 檢查應用程式健康狀態
check_application_health() {
    log_info "檢查應用程式健康狀態..."
    
    APP_URL=$(terraform output -raw application_url 2>/dev/null || echo "")
    if [ ! -z "$APP_URL" ]; then
        HEALTH_URL="${APP_URL}/health"
        
        log_info "檢查健康檢查端點: $HEALTH_URL"
        
        HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$HEALTH_URL" 2>/dev/null || echo "000")
        
        if [ "$HTTP_STATUS" = "200" ]; then
            log_success "應用程式健康檢查通過 (HTTP $HTTP_STATUS)"
        elif [ "$HTTP_STATUS" != "000" ]; then
            log_warning "應用程式健康檢查失敗 (HTTP $HTTP_STATUS)"
        else
            log_error "無法連接到應用程式"
        fi
    else
        log_warning "無法獲取應用程式 URL"
    fi
}

# 顯示重要資訊
show_important_info() {
    log_info "重要資訊摘要："
    echo
    
    # 應用程式 URL
    APP_URL=$(terraform output -raw application_url 2>/dev/null || echo "")
    if [ ! -z "$APP_URL" ]; then
        echo "🌐 應用程式 URL: $APP_URL"
    fi
    
    # API 基礎 URL
    if [ ! -z "$APP_URL" ]; then
        echo "🔧 API 基礎 URL: ${APP_URL}/api"
        echo "📚 Swagger UI: ${APP_URL}/swagger"
    fi
    
    # 資料庫端點
    DB_ENDPOINT=$(terraform output -raw database_endpoint 2>/dev/null || echo "")
    if [ ! -z "$DB_ENDPOINT" ]; then
        echo "🗃️  資料庫端點: $DB_ENDPOINT"
    fi
    
    # CloudWatch Dashboard
    DASHBOARD_URL=$(terraform output -raw cloudwatch_dashboard_url 2>/dev/null || echo "")
    if [ ! -z "$DASHBOARD_URL" ]; then
        echo "📊 監控儀表板: $DASHBOARD_URL"
    fi
    
    echo
}

# 顯示快速命令
show_quick_commands() {
    log_info "常用命令："
    echo
    echo "📋 查看所有輸出:"
    echo "   terraform output"
    echo
    echo "📊 查看 ECS 服務狀態:"
    echo "   aws ecs describe-services --cluster \$(terraform output -raw ecs_cluster_name) --services \$(terraform output -raw ecs_service_name)"
    echo
    echo "📝 查看應用程式日誌:"
    echo "   aws logs tail \$(terraform output -raw cloudwatch_log_group_name) --follow"
    echo
    echo "🔍 檢查資料庫狀態:"
    echo "   aws rds describe-db-instances --query 'DBInstances[?contains(DBInstanceIdentifier, \`security-intel-dev\`)].{ID:DBInstanceIdentifier,Status:DBInstanceStatus,Engine:Engine}'"
    echo
    echo "🗑️  銷毀環境:"
    echo "   ./destroy.sh"
    echo
}

# 主要函數
main() {
    echo "=== Security Intelligence Platform 狀態檢查 ==="
    echo
    
    # 檢查 Terraform 狀態
    if check_terraform_status; then
        # 檢查 AWS 資源
        check_aws_resources
        
        # 檢查應用程式健康狀態
        check_application_health
        
        echo
        # 顯示重要資訊
        show_important_info
        
        # 顯示快速命令
        show_quick_commands
    else
        log_error "Terraform 狀態檢查失敗，請先初始化並部署環境。"
        echo
        log_info "如需部署環境，請執行: ./deploy.sh"
    fi
}

# 執行主要函數
main "$@" 