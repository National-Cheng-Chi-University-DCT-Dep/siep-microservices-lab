# Security Intelligence Platform - 開發環境部署指南

這個目錄包含了 Security Intelligence Platform 開發環境的 Terraform 配置。

## 📋 目錄結構

```
terraform/environments/dev/
├── README.md                    # 本檔案
├── backend.tf                   # Terraform 後端配置
├── variables.tf                 # 變數定義
├── main.tf                      # 主要配置檔案
├── outputs.tf                   # 輸出定義
├── terraform.tfvars.example     # 配置範例檔案
├── deploy.sh                    # 自動部署腳本
└── destroy.sh                   # 自動銷毀腳本
```

## 🚀 快速開始

### 1. 前置準備

確保您已經安裝並配置了以下工具：

- [Terraform](https://terraform.io/downloads) (>= 1.0)
- [AWS CLI](https://aws.amazon.com/cli/) (>= 2.0)
- [jq](https://stedolan.github.io/jq/) (用於 JSON 處理)

### 2. 配置 AWS 憑證

```bash
aws configure
```

### 3. 準備配置檔案

```bash
# 複製範例配置檔案
cp terraform.tfvars.example terraform.tfvars

# 編輯配置檔案
vim terraform.tfvars
```

**重要：** 請務必修改 `terraform.tfvars` 中的敏感資訊，特別是：

- `jwt_secret`
- `encryption_key`
- `api_key`
- `third_party_tokens`
- `security_alert_email`

### 4. 部署環境

使用自動化腳本部署：

```bash
./deploy.sh
```

或手動部署：

```bash
# 初始化 Terraform
terraform init

# 驗證配置
terraform validate

# 規劃部署
terraform plan

# 應用部署
terraform apply
```

### 5. 存取應用程式

部署完成後，您可以通過以下方式存取：

```bash
# 查看所有輸出
terraform output

# 查看應用程式 URL
terraform output application_url

# 查看開發者快速存取資訊
terraform output developer_quick_access
```

## 🏗️ 架構概覽

這個 Terraform 配置會建立以下 AWS 資源：

### 網路層 (VPC Module)

- VPC 和子網路
- Internet Gateway 和 NAT Gateway
- 路由表
- 安全群組
- VPC Flow Logs

### 資料庫層 (RDS Module)

- PostgreSQL RDS 實例
- 資料庫子網路群組
- 資料庫安全群組
- 資料庫參數群組
- 自動備份配置

### 運算層 (ECS Module)

- ECS 叢集和服務
- Application Load Balancer
- ECS 任務定義
- 自動擴展配置
- CloudWatch 日誌

### 安全性層 (Security Module)

- KMS 金鑰（加密）
- Secrets Manager（機密管理）
- IAM 角色和政策
- WAF（Web Application Firewall）
- CloudWatch 警報

## 🔧 配置說明

### 環境變數

主要的配置變數包括：

| 變數名稱            | 描述           | 預設值           |
| ------------------- | -------------- | ---------------- |
| `project_name`      | 專案名稱       | `security-intel` |
| `environment`       | 環境名稱       | `dev`            |
| `aws_region`        | AWS 區域       | `ap-northeast-1` |
| `vpc_cidr`          | VPC CIDR 區塊  | `10.0.0.0/16`    |
| `db_instance_class` | RDS 實例類型   | `db.t3.micro`    |
| `ecs_task_cpu`      | ECS 任務 CPU   | `256`            |
| `ecs_task_memory`   | ECS 任務記憶體 | `512`            |

### 敏感資訊管理

所有敏感資訊都儲存在 AWS Secrets Manager 中：

- JWT 簽名密鑰
- 資料加密金鑰
- API 金鑰
- 第三方服務 Token

## 📊 監控和日誌

### CloudWatch Dashboard

部署完成後，您可以通過以下 URL 存取監控儀表板：

```bash
terraform output cloudwatch_dashboard_url
```

### 日誌查看

應用程式日誌儲存在 CloudWatch Logs 中：

```bash
# 查看 ECS 服務日誌
aws logs describe-log-groups --log-group-name-prefix "/ecs/security-intel-dev"

# 查看資料庫日誌
aws logs describe-log-groups --log-group-name-prefix "/aws/rds"
```

### 警報通知

系統會在以下情況發送警報：

- ECS 服務 CPU 或記憶體使用率過高
- RDS 資料庫 CPU 使用率過高
- RDS 資料庫連接數過多
- RDS 資料庫可用儲存空間不足
- WAF 封鎖大量請求
- 應用程式失敗登入次數過多

## 🔒 安全性考量

### 網路安全

- 資料庫部署在私有子網路
- 應用程式通過 ALB 存取
- WAF 提供 Web 應用程式防護
- 安全群組限制網路存取

### 資料安全

- 所有資料使用 KMS 加密
- 資料庫啟用加密
- Secrets Manager 管理敏感資訊
- VPC Flow Logs 記錄網路流量

### 存取控制

- IAM 角色使用最小權限原則
- ECS 任務使用專用角色
- 資料庫認證通過 Secrets Manager

## 🛠️ 疑難排解

### 常見問題

1. **部署失敗：S3 儲存桶不存在**

   ```bash
   # 手動建立 S3 儲存桶
   aws s3api create-bucket --bucket your-tfstate-bucket --region ap-northeast-1
   ```

2. **部署失敗：DynamoDB 表格不存在**

   ```bash
   # 手動建立 DynamoDB 表格
   aws dynamodb create-table --table-name your-tfstate-lock \
     --attribute-definitions AttributeName=LockID,AttributeType=S \
     --key-schema AttributeName=LockID,KeyType=HASH \
     --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
   ```

3. **應用程式無法存取資料庫**

   - 檢查安全群組設定
   - 檢查 Secrets Manager 中的資料庫認證
   - 檢查 ECS 任務日誌

4. **憑證驗證失敗**
   - 確保域名的 DNS 設定正確
   - 檢查 Route 53 記錄

### 除錯命令

```bash
# 檢查 Terraform 狀態
terraform show

# 檢查特定資源
terraform state show module.vpc.aws_vpc.main

# 檢查 ECS 服務狀態
aws ecs describe-services --cluster security-intel-dev-cluster --services security-intel-dev-app-service

# 檢查 RDS 實例狀態
aws rds describe-db-instances --db-instance-identifier security-intel-dev-postgres

# 檢查應用程式日誌
aws logs tail /ecs/security-intel-dev --follow
```

## 🗑️ 清理資源

### 使用自動化腳本

```bash
./destroy.sh
```

### 手動清理

```bash
terraform destroy
```

**警告：** 這會永久刪除所有資源，包括資料庫中的所有資料！

## 📞 支援

如果遇到問題，請檢查：

1. [Terraform 官方文件](https://terraform.io/docs)
2. [AWS Provider 文件](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
3. 專案的 GitHub Issues

## 📝 更新日誌

- **v1.0.0** - 初始版本，包含完整的開發環境配置
