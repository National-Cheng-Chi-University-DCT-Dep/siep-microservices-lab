# Terraform 基礎設施

本目錄包含資安情報平台的 AWS 基礎設施程式碼，使用 Terraform 進行管理。

## 目錄結構

```
terraform/
├── environments/          # 環境特定配置
│   ├── dev/              # 開發環境
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── prod/             # 生產環境
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
├── modules/              # 可重複使用的模組
│   ├── vpc/              # VPC 模組
│   ├── rds/              # RDS 模組
│   └── ecs/              # ECS 模組
└── README.md
```

## 前置準備

### 1. 安裝必要工具

```bash
# 安裝 Terraform (版本 >= 1.5)
# macOS (使用 Homebrew)
brew install terraform

# 或下載二進位檔案
# https://www.terraform.io/downloads.html

# 安裝 AWS CLI
brew install awscli

# 或使用 pip
pip install awscli
```

### 2. 設定 AWS 憑證

```bash
# 設定 AWS 憑證
aws configure

# 或使用環境變數
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_DEFAULT_REGION="ap-northeast-1"
```

### 3. 建立 S3 Bucket (遠端狀態存儲)

在使用 Terraform 之前，需要手動建立 S3 bucket 來存儲 Terraform 狀態檔案：

```bash
# 建立 S3 bucket
aws s3 mb s3://security-intel-tfstate-2024 --region ap-northeast-1

# 啟用版本控制
aws s3api put-bucket-versioning \
    --bucket security-intel-tfstate-2024 \
    --versioning-configuration Status=Enabled

# 啟用伺服器端加密
aws s3api put-bucket-encryption \
    --bucket security-intel-tfstate-2024 \
    --server-side-encryption-configuration '{
        "Rules": [
            {
                "ApplyServerSideEncryptionByDefault": {
                    "SSEAlgorithm": "AES256"
                }
            }
        ]
    }'
```

### 4. 建立 DynamoDB 表格 (狀態鎖定)

```bash
# 建立 DynamoDB 表格用於狀態鎖定
aws dynamodb create-table \
    --table-name terraform-state-lock \
    --attribute-definitions AttributeName=LockID,AttributeType=S \
    --key-schema AttributeName=LockID,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --region ap-northeast-1
```

## 部署指南

### 開發環境部署

```bash
# 進入開發環境目錄
cd environments/dev

# 初始化 Terraform
terraform init

# 檢查計畫
terraform plan

# 應用變更
terraform apply

# 查看輸出
terraform output
```

### 生產環境部署

```bash
# 進入生產環境目錄
cd environments/prod

# 初始化 Terraform
terraform init

# 檢查計畫
terraform plan

# 應用變更 (需要確認)
terraform apply
```

## 模組說明

### VPC 模組

建立虛擬私有雲 (VPC) 和相關網路資源：

- VPC
- 公開和私有子網路
- Internet Gateway
- NAT Gateway
- 路由表

### RDS 模組

建立 PostgreSQL 資料庫：

- RDS 實例
- 子網路群組
- 安全群組
- 密碼管理 (AWS Secrets Manager)

### ECS 模組

建立容器運行環境：

- ECS 叢集
- ECS 服務
- 應用程式負載平衡器 (ALB)
- ECR 儲存庫
- IAM 角色和政策
- CloudWatch 日誌

## 環境變數

可以透過 `.tfvars` 檔案自訂變數：

```hcl
# dev.tfvars
aws_region        = "ap-northeast-1"
vpc_cidr         = "10.0.0.0/16"
db_instance_class = "db.t3.micro"
container_port   = 8080
```

使用方式：

```bash
terraform plan -var-file="dev.tfvars"
terraform apply -var-file="dev.tfvars"
```

## 常用命令

### 基本操作

```bash
# 初始化
terraform init

# 格式化程式碼
terraform fmt

# 驗證配置
terraform validate

# 計畫變更
terraform plan

# 應用變更
terraform apply

# 銷毀資源
terraform destroy
```

### 狀態管理

```bash
# 查看狀態
terraform state list

# 查看特定資源
terraform state show <resource>

# 匯入現有資源
terraform import <resource> <id>

# 移除資源 (不銷毀)
terraform state rm <resource>
```

### 輸出和檢查

```bash
# 查看輸出
terraform output

# 查看特定輸出
terraform output <output_name>

# 以 JSON 格式輸出
terraform output -json
```

## 最佳實踐

### 1. 版本控制

- 使用 Git 管理 Terraform 程式碼
- 不要提交 `.tfstate` 檔案
- 使用 `.tfvars` 檔案管理環境特定變數

### 2. 安全性

- 使用 AWS Secrets Manager 管理敏感資料
- 限制 IAM 權限 (最小權限原則)
- 啟用 S3 bucket 加密和版本控制

### 3. 結構化

- 使用模組化設計
- 分離環境配置
- 使用一致的命名規則

### 4. 測試

- 在開發環境測試變更
- 使用 `terraform plan` 檢查變更
- 定期執行 `terraform validate`

## 監控和日誌

部署後，可以透過以下方式監控基礎設施：

### CloudWatch

```bash
# 查看 ECS 服務日誌
aws logs describe-log-groups --log-group-name-prefix="/ecs/security-intel"

# 查看最近的日誌
aws logs tail /ecs/security-intel-platform-dev --follow
```

### 應用程式健康檢查

```bash
# 取得負載平衡器 DNS
terraform output load_balancer_dns

# 檢查應用程式狀態
curl http://$(terraform output -raw load_balancer_dns)/api/v1/health
```

## 疑難排解

### 常見問題

1. **初始化失敗**

   - 檢查 AWS 憑證設定
   - 確認 S3 bucket 存在且可存取

2. **計畫失敗**

   - 檢查 IAM 權限
   - 確認資源名稱唯一性

3. **應用失敗**
   - 檢查配額限制
   - 查看 CloudWatch 日誌

### 除錯技巧

```bash
# 啟用詳細日誌
export TF_LOG=DEBUG

# 檢查 AWS 憑證
aws sts get-caller-identity

# 檢查區域設定
aws configure get region
```

## 成本優化

### 開發環境

- 使用 `t3.micro` 或 `t3.small` 實例
- 啟用 ECS Fargate Spot
- 設定自動縮放政策

### 生產環境

- 使用預留實例
- 啟用 S3 生命週期政策
- 定期審查未使用資源

## 升級和維護

### 升級 Terraform

```bash
# 檢查目前版本
terraform version

# 升級後重新初始化
terraform init -upgrade
```

### 升級 AWS Provider

```bash
# 更新 provider 版本
terraform init -upgrade

# 檢查變更
terraform plan
```

## 備份和災難恢復

### 狀態檔案備份

- S3 自動版本控制
- 定期備份 `.tfstate` 檔案
- 測試恢復程序

### 資料庫備份

- 啟用 RDS 自動備份
- 定期建立快照
- 測試恢復程序

## 聯絡資訊

如有問題或建議，請聯絡：

- 專案維護者：[Your Name]
- 問題回報：[GitHub Issues]
- 文件：[Project Documentation]
