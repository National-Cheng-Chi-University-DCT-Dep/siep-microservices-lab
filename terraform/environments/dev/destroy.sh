#!/bin/bash

# =============================================================================
# Ultimate Security Intelligence Platform - Development Environment Destroy
# 開發環境 Terraform 銷毀腳本
# =============================================================================

set -e  # 遇到錯誤時停止執行

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

# 檢查必要的工具
check_prerequisites() {
    log_info "檢查必要工具..."
    
    # 檢查 Terraform
    if ! command -v terraform &> /dev/null; then
        log_error "Terraform 未安裝。請先安裝 Terraform。"
        exit 1
    fi
    
    # 檢查 AWS CLI
    if ! command -v aws &> /dev/null; then
        log_error "AWS CLI 未安裝。請先安裝 AWS CLI。"
        exit 1
    fi
    
    # 檢查 AWS 憑證
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS 憑證未配置。請執行 'aws configure' 配置憑證。"
        exit 1
    fi
    
    log_success "所有必要工具都已安裝並配置完成。"
}

# 檢查 Terraform 狀態
check_terraform_state() {
    log_info "檢查 Terraform 狀態..."
    
    if [ ! -f ".terraform/terraform.tfstate" ] && [ ! -f "terraform.tfstate" ]; then
        if ! terraform init &>/dev/null; then
            log_error "無法初始化 Terraform。請確保配置正確。"
            exit 1
        fi
    fi
    
    # 檢查是否有資源需要銷毀
    if terraform show -json | jq -e '.values.root_module.resources | length > 0' &>/dev/null; then
        log_info "發現需要銷毀的資源。"
    else
        log_info "沒有發現需要銷毀的資源。"
        read -p "是否要繼續執行銷毀操作？ (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "銷毀操作已取消。"
            exit 0
        fi
    fi
    
    log_success "Terraform 狀態檢查完成。"
}

# 顯示將要銷毀的資源
show_destroy_plan() {
    log_info "顯示將要銷毀的資源..."
    terraform plan -destroy
    log_warning "以上資源將被永久刪除！"
}

# 確認銷毀操作
confirm_destroy() {
    echo
    log_warning "⚠️  警告：這將永久刪除所有 AWS 資源！"
    log_warning "⚠️  包括資料庫、儲存桶、網路設備等所有資源！"
    log_warning "⚠️  此操作無法復原！"
    echo
    
    # 第一次確認
    read -p "您確定要刪除所有資源嗎？請輸入 'yes' 確認: " confirm1
    if [ "$confirm1" != "yes" ]; then
        log_info "銷毀操作已取消。"
        exit 0
    fi
    
    # 第二次確認
    echo
    log_warning "最後確認：這將刪除所有數據，包括資料庫中的所有記錄！"
    read -p "請再次輸入 'DELETE' 確認刪除: " confirm2
    if [ "$confirm2" != "DELETE" ]; then
        log_info "銷毀操作已取消。"
        exit 0
    fi
    
    log_info "確認完成，準備開始銷毀資源..."
}

# 備份重要資料（可選）
backup_data() {
    read -p "是否要在銷毀前備份資料庫？ (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "正在建立資料庫備份..."
        
        # 獲取資料庫實例 ID
        DB_INSTANCE_ID=$(terraform output -raw db_instance_id 2>/dev/null || echo "")
        
        if [ ! -z "$DB_INSTANCE_ID" ]; then
            SNAPSHOT_ID="${DB_INSTANCE_ID}-final-backup-$(date +%Y%m%d-%H%M%S)"
            
            log_info "建立資料庫快照: $SNAPSHOT_ID"
            aws rds create-db-snapshot \
                --db-instance-identifier "$DB_INSTANCE_ID" \
                --db-snapshot-identifier "$SNAPSHOT_ID"
            
            log_info "等待快照建立完成（這可能需要幾分鐘）..."
            aws rds wait db-snapshot-completed --db-snapshot-identifier "$SNAPSHOT_ID"
            
            log_success "資料庫備份完成: $SNAPSHOT_ID"
        else
            log_warning "無法找到資料庫實例 ID，跳過備份。"
        fi
    else
        log_info "跳過資料庫備份。"
    fi
}

# 執行 Terraform 銷毀
terraform_destroy() {
    log_info "開始銷毀 Terraform 資源..."
    log_info "這個過程可能需要 10-20 分鐘，請耐心等候..."
    
    # 使用 -auto-approve 跳過互動式確認，因為我們已經做過確認了
    terraform destroy -auto-approve
    
    log_success "Terraform 銷毀完成。"
}

# 清理本地檔案
cleanup_local_files() {
    log_info "清理本地暫存檔案..."
    
    # 清理 Terraform 暫存檔案
    rm -f tfplan
    rm -f terraform.tfstate.backup
    
    # 詢問是否清理 .terraform 目錄
    read -p "是否要清理 .terraform 目錄？ (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf .terraform
        log_info ".terraform 目錄已清理。"
    fi
    
    log_success "本地檔案清理完成。"
}

# 檢查剩餘資源
check_remaining_resources() {
    log_info "檢查是否有剩餘的 AWS 資源..."
    
    # 檢查是否還有 EC2 實例
    INSTANCES=$(aws ec2 describe-instances --filters "Name=tag:Project,Values=security-intel" "Name=instance-state-name,Values=running,pending,stopping,stopped" --query 'Reservations[*].Instances[*].InstanceId' --output text)
    if [ ! -z "$INSTANCES" ]; then
        log_warning "發現剩餘的 EC2 實例: $INSTANCES"
    fi
    
    # 檢查是否還有 RDS 實例
    RDS_INSTANCES=$(aws rds describe-db-instances --query 'DBInstances[?contains(DBInstanceIdentifier, `security-intel`)].DBInstanceIdentifier' --output text)
    if [ ! -z "$RDS_INSTANCES" ]; then
        log_warning "發現剩餘的 RDS 實例: $RDS_INSTANCES"
    fi
    
    # 檢查是否還有負載平衡器
    LOAD_BALANCERS=$(aws elbv2 describe-load-balancers --query 'LoadBalancers[?contains(LoadBalancerName, `security-intel`)].LoadBalancerName' --output text)
    if [ ! -z "$LOAD_BALANCERS" ]; then
        log_warning "發現剩餘的負載平衡器: $LOAD_BALANCERS"
    fi
    
    log_info "剩餘資源檢查完成。"
}

# 顯示費用節省資訊
show_cost_savings() {
    log_success "=== 銷毀完成 ==="
    echo
    log_info "所有 AWS 資源已被刪除，您不會再產生相關費用。"
    log_info "如果您建立了資料庫快照，請注意快照會產生少量的儲存費用。"
    log_info "您可以在 AWS 控制台中手動刪除不需要的快照。"
    echo
    log_info "如果需要重新部署，請執行 './deploy.sh' 腳本。"
}

# 主要函數
main() {
    log_info "開始銷毀 Security Intelligence Platform 開發環境..."
    
    # 檢查先決條件
    check_prerequisites
    check_terraform_state
    
    # 顯示銷毀計劃
    show_destroy_plan
    
    # 確認銷毀操作
    confirm_destroy
    
    # 備份資料（可選）
    backup_data
    
    # 執行銷毀
    terraform_destroy
    
    # 清理本地檔案
    cleanup_local_files
    
    # 檢查剩餘資源
    check_remaining_resources
    
    # 顯示完成資訊
    show_cost_savings
}

# 錯誤處理
trap 'log_error "銷毀過程中發生錯誤。請檢查 AWS 控制台確認資源狀態。"; exit 1' ERR

# 執行主要函數
main "$@" 