#!/bin/bash

# =============================================================================
# Ultimate Security Intelligence Platform - Development Environment Deployment
# 開發環境 Terraform 部署腳本
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

# 檢查配置檔案
check_configuration() {
    log_info "檢查配置檔案..."
    
    if [ ! -f "terraform.tfvars" ]; then
        log_warning "terraform.tfvars 檔案不存在。"
        if [ -f "terraform.tfvars.example" ]; then
            log_info "複製 terraform.tfvars.example 為 terraform.tfvars..."
            cp terraform.tfvars.example terraform.tfvars
            log_warning "請編輯 terraform.tfvars 檔案，設定正確的變數值。"
            log_warning "特別是敏感資訊（如 jwt_secret、encryption_key 等）。"
            read -p "是否要繼續部署？ (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log_info "部署已取消。請先配置 terraform.tfvars 檔案。"
                exit 0
            fi
        else
            log_error "terraform.tfvars.example 檔案不存在。"
            exit 1
        fi
    fi
    
    log_success "配置檔案檢查完成。"
}

# 檢查 S3 儲存桶
check_s3_bucket() {
    log_info "檢查 Terraform 狀態儲存桶..."
    
    BUCKET_NAME=$(grep -E "^\s*bucket\s*=" backend.tf | sed 's/.*=\s*"\(.*\)".*/\1/')
    
    if [ -z "$BUCKET_NAME" ]; then
        log_error "無法從 backend.tf 找到 S3 儲存桶名稱。"
        exit 1
    fi
    
    if ! aws s3api head-bucket --bucket "$BUCKET_NAME" 2>/dev/null; then
        log_warning "S3 儲存桶 '$BUCKET_NAME' 不存在。"
        read -p "是否要建立此儲存桶？ (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            log_info "建立 S3 儲存桶 '$BUCKET_NAME'..."
            aws s3api create-bucket --bucket "$BUCKET_NAME" --region ap-northeast-1 --create-bucket-configuration LocationConstraint=ap-northeast-1
            aws s3api put-bucket-versioning --bucket "$BUCKET_NAME" --versioning-configuration Status=Enabled
            aws s3api put-bucket-encryption --bucket "$BUCKET_NAME" --server-side-encryption-configuration '{
                "Rules": [
                    {
                        "ApplyServerSideEncryptionByDefault": {
                            "SSEAlgorithm": "AES256"
                        }
                    }
                ]
            }'
            log_success "S3 儲存桶建立完成。"
        else
            log_error "需要 S3 儲存桶才能繼續部署。"
            exit 1
        fi
    else
        log_success "S3 儲存桶 '$BUCKET_NAME' 存在。"
    fi
}

# 檢查 DynamoDB 表格
check_dynamodb_table() {
    log_info "檢查 DynamoDB 鎖定表格..."
    
    TABLE_NAME=$(grep -E "^\s*dynamodb_table\s*=" backend.tf | sed 's/.*=\s*"\(.*\)".*/\1/')
    
    if [ -z "$TABLE_NAME" ]; then
        log_error "無法從 backend.tf 找到 DynamoDB 表格名稱。"
        exit 1
    fi
    
    if ! aws dynamodb describe-table --table-name "$TABLE_NAME" &>/dev/null; then
        log_warning "DynamoDB 表格 '$TABLE_NAME' 不存在。"
        read -p "是否要建立此表格？ (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            log_info "建立 DynamoDB 表格 '$TABLE_NAME'..."
            aws dynamodb create-table \
                --table-name "$TABLE_NAME" \
                --attribute-definitions AttributeName=LockID,AttributeType=S \
                --key-schema AttributeName=LockID,KeyType=HASH \
                --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
            
            log_info "等待 DynamoDB 表格建立完成..."
            aws dynamodb wait table-exists --table-name "$TABLE_NAME"
            log_success "DynamoDB 表格建立完成。"
        else
            log_error "需要 DynamoDB 表格才能繼續部署。"
            exit 1
        fi
    else
        log_success "DynamoDB 表格 '$TABLE_NAME' 存在。"
    fi
}

# 初始化 Terraform
terraform_init() {
    log_info "初始化 Terraform..."
    terraform init
    log_success "Terraform 初始化完成。"
}

# 驗證 Terraform 配置
terraform_validate() {
    log_info "驗證 Terraform 配置..."
    terraform validate
    log_success "Terraform 配置驗證通過。"
}

# 格式化 Terraform 檔案
terraform_fmt() {
    log_info "格式化 Terraform 檔案..."
    terraform fmt -recursive
    log_success "Terraform 檔案格式化完成。"
}

# 規劃 Terraform 部署
terraform_plan() {
    log_info "規劃 Terraform 部署..."
    terraform plan -out=tfplan
    log_success "Terraform 規劃完成。"
}

# 應用 Terraform 部署
terraform_apply() {
    log_info "應用 Terraform 部署..."
    
    if [ -f "tfplan" ]; then
        terraform apply tfplan
    else
        log_warning "找不到 tfplan 檔案，直接應用..."
        terraform apply
    fi
    
    log_success "Terraform 部署完成。"
}

# 清理暫存檔案
cleanup() {
    log_info "清理暫存檔案..."
    rm -f tfplan
    log_success "清理完成。"
}

# 顯示部署結果
show_outputs() {
    log_info "顯示部署結果..."
    terraform output
    log_success "部署結果顯示完成。"
}

# 主要函數
main() {
    log_info "開始部署 Security Intelligence Platform 開發環境..."
    
    # 檢查先決條件
    check_prerequisites
    check_configuration
    check_s3_bucket
    check_dynamodb_table
    
    # Terraform 操作
    terraform_init
    terraform_validate
    terraform_fmt
    terraform_plan
    
    # 確認部署
    echo
    log_warning "即將開始部署 AWS 資源。這可能需要 10-15 分鐘。"
    read -p "是否要繼續部署？ (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        terraform_apply
        show_outputs
        cleanup
        
        echo
        log_success "=== 部署完成 ==="
        log_info "您現在可以存取以下服務："
        echo
        terraform output -json | jq -r '.developer_quick_access.value | to_entries[] | "  \(.key): \(.value)"'
        echo
        log_info "請查看 CloudWatch Dashboard 以監控系統狀態。"
        log_info "使用 'terraform destroy' 命令可以刪除所有資源。"
    else
        log_info "部署已取消。"
        cleanup
    fi
}

# 錯誤處理
trap 'log_error "部署過程中發生錯誤。"; cleanup; exit 1' ERR

# 執行主要函數
main "$@" 