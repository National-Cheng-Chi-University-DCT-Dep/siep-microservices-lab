案實踐藍圖：六階段實施計畫
Phase 0: 專案啟動與基礎設定 (約 1-2 天)
這是打地基的階段，目標是建立一個乾淨、標準化的開發起點。

目標：建立所有開發人員都能依循的統一專案結構與規範。

步驟：

建立 Git 儲存庫：在 GitHub 或 GitLab 上創建名為 security-intel-platform 的專案。

初始化專案結構：根據我們規劃的優化專案結構，手動創建所有頂層目錄 (backend, frontend, terraform, docs, .github/workflows 等)。

加入基礎設定檔：

在根目錄創建一個全面的 .gitignore 檔案（涵蓋 Go, Node.js, Terraform, OS 系統檔案）。

創建 LICENSE (MIT) 和 CODE_OF_CONDUCT.md。

將我們之前完成的 README.md 檔案放入專案根目錄。

初始化前後端專案：

cd backend 然後執行 go mod init <your-repo-path>/backend。

cd frontend 然後執行 npm create next-app@latest . 來建立 Next.js 專案。

主要交付成果：一個結構清晰、已設定好版本控制的空白專案儲存庫。

Phase 1: 本地開發 MVP (最小可行性產品) (約 1-2 週)
此階段的重點是讓核心功能在開發人員的電腦上順利運作。

目標：實現前後端與資料庫在本地環境的串接與核心功能。

步驟：

後端核心 API：

設計並實作第一個資料庫 migration 腳本，建立儲存情報的資料表。

撰寫第一個 API Endpoint (例如 GET /api/threats)，能從資料庫讀取資料。

撰寫第一個情報蒐集器的原型，能手動觸發從 AbuseIPDB 抓取資料並存入資料庫。

前端核心介面：

建立一個基本的儀表板頁面。

實作一個輸入框和按鈕，能呼叫後端的 API 並將結果顯示在一個簡單的表格中。

本地開發環境整合：

在 docker/ 目錄下完成 docker-compose.yml 的撰寫，確保 docker-compose up 可以一鍵啟動本地的 PostgreSQL 資料庫。

確保後端可以透過環境變數 (.env) 連接到此資料庫。

主要交付成果：

開發者 git clone 後，可以透過 docker-compose, make run, npm run dev 在本地完整運行一個具備「蒐集-儲存-展示」基本流程的應用程式。

Phase 2: 持續整合與容器化 (約 1 週)
讓程式碼的品質檢查和打包流程自動化。

目標：實現每次程式碼提交都能自動檢查、測試，並將後端打包成標準化的 Docker 鏡像。

步驟：

建立 CI 工作流程 (GitHub Actions)：

實作 .github/workflows/code-quality.yml，包含 Linter、單元測試、安全性掃描等。

確保每次 PR 和 push 都會觸發此流程，並在 GitHub 上看到檢查結果。

後端容器化 (Dockerfile)：

在 backend/ 目錄下撰寫一個多階段建置 (multi-stage build) 的 Dockerfile，以產生體積小且安全的正式鏡像。

設定 Drone CI：

架設 Drone CI 平台。

設定 .drone.yml，讓它在接收到 GitHub 的 Webhook 後，自動建置後端 Docker 鏡像，並將其推送到 AWS ECR。

主要交付成果：

一個全自動的 CI 流程：git push ➔ GitHub Actions 檢查 ➔ Drone CI 建置 ➔ Docker 鏡像出現在 ECR。

Phase 3: 雲端基礎設施即程式碼 (IaC) (約 1 週)
用程式碼打造應用的雲端家園。

AWS 環境建置的詳細步驟
整個 AWS 環境的建置可以分為兩大部分：「前置準備 (手動設定)」 和 「使用 Terraform 自動化建置」。

Part 1: 前置準備 (手動設定，在 Phase 3 開始前完成)
這部分是在執行任何自動化程式碼之前，必須先在 AWS Console 中手動完成的一次性設定。

註冊 AWS 帳戶

如果還沒有，請先註冊一個 AWS 帳戶。建議啟用免費方案 (Free Tier) 來降低初期成本。

設定帳戶安全與帳單警示

啟用 MFA (多因素驗證)：為您的「根使用者 (Root User)」啟用 MFA，這是最重要的安全第一步。

設定帳單警示 (Billing Alert)：在 AWS Budgets 中設定一個警示（例如：預算超過 $5 美元時發送郵件），避免費用超乎預期。

建立 IAM 使用者 (for Terraform & CI/CD)

原則：遵循最小權限原則，絕不使用根使用者進行日常操作。

操作：

前往 AWS IAM 服務。

建立一個新的使用者（例如 terraform-admin）。

賦予此使用者 AdministratorAccess 權限（在初期學習階段可以這樣設定，正式生產環境應收緊權限）。

產生一組存取金鑰 (Access Key ID 和 Secret Access Key)。

極度重要：將這組金鑰安全地儲存在您的密碼管理器中，並設定到您本地電腦的 AWS CLI (aws configure)。切勿將金鑰寫死在任何程式碼或提交到 Git。

建立 S3 儲存桶 (for Terraform Remote State)

目的：Terraform 需要一個地方來存放雲端資源的狀態檔案 (.tfstate)，以便團隊協作。

操作：

前往 AWS S3 服務。

手動建立一個全域唯一的 S3 儲存桶 (Bucket)，例如 my-security-intel-tfstate-2025。

(建議) 啟用「版本控制 (Versioning)」和「伺服器端加密 (Server-side encryption)」功能。

Part 2: 使用 Terraform 自動化建置 (Phase 3 核心)
完成前置準備後，接下來的所有 AWS 資源都應該用 Terraform 程式碼來建立。這就是 Phase 3 的核心任務。

初始化 Terraform 專案

在 terraform/environments/dev/ 目錄下，建立 backend.tf 檔案，設定 S3 儲存桶為遠端狀態後端。

執行 terraform init。

編寫網路層 (VPC)

在 terraform/modules/vpc 中，定義您專案的虛擬私有雲，包括：

VPC: 專案的獨立網路空間。

Subnets: 公開 (Public) 和私有 (Private) 子網路。

Internet Gateway & NAT Gateway: 讓服務能與外部網路溝通。

編寫資料庫層 (RDS)

在 terraform/modules/rds 中，定義 PostgreSQL 資料庫實例。

確保將其放置在私有子網路中以策安全。

透過 Terraform 與 AWS Secrets Manager 整合，自動產生並儲存資料庫密碼。

編寫運算層 (ECS on Fargate)

在 terraform/modules/ecs 中，定義運行後端 Docker 容器所需的一切：

ECS Cluster: 容器的叢集。

Task Definition: 容器的藍圖（使用哪個 Docker 鏡像、CPU/記憶體需求、環境變數等）。

ECS Service: 確保指定數量的容器持續運行，並設定負載平衡器 (ALB) 將流量導向容器。

編寫支援服務與安全性規則

IAM Roles: 為 ECS Task 建立專屬的角色，讓它有權限存取其他 AWS 服務（如 S3、Secrets Manager），而無需硬編碼金鑰。

Security Groups: 定義防火牆規則，例如「只允許負載平衡器存取後端服務的 8080 端口」、「只允許後端服務存取資料庫的 5432 端口」。

執行首次部署

在所有程式碼編寫完成後，回到 terraform/environments/dev/ 目錄。

執行 terraform plan 預覽將要建立的所有資源。

確認無誤後，執行 terraform apply。Terraform 會根據您的程式碼，在幾分鐘內自動在 AWS 上建置出完整的、互相連接的雲端環境。

目標：使用 Terraform 在 AWS 上建立所有必要的雲端資源。

步驟：

設定遠端狀態 (Remote State)：設定 Terraform 使用 AWS S3 bucket 來儲存 tfstate 檔案，以便團隊協作。

撰寫 Terraform 模組：

建立網路模組 (modules/vpc)：包含 VPC, Subnets, Internet Gateway 等。

建立資料庫模組 (modules/rds)：建立一個 RDS for PostgreSQL 實例。

建立運算層模組 (modules/ecs)：建立 ECS Cluster、Task Definition、Service 等，準備好運行我們的 Docker 容器。

部署開發環境：

在 terraform/environments/dev/ 中引用上述模組，並執行 terraform apply，手動部署第一套開發環境。

手動將資料庫密碼等 secrets 存入 AWS Secrets Manager。

主要交付成果：

一個完整、可重複部署的 AWS 環境，所有資源都由 Terraform 程式碼管理。

Phase 4: 持續部署 (CD) (約 1 週)
將自動化流程的最後一哩路打通，從程式碼到上線。

目標：實現 git push 後，應用程式能自動部署到 AWS。

步驟：

設定前端 CD：將 GitHub 專案與 Vercel 連結，實現前端的自動化部署。

設定後端 CD (AWS CodePipeline)：

建立 backend-deploy-pipeline，來源設為 ECR。

設定部署階段，當 ECR 有新鏡像時，自動更新在 Phase 3 建立的 ECS Service。

設定 IaC CD (AWS CodePipeline)：

建立 infra-terraform-pipeline，來源設為 GitHub 的 terraform/ 目錄。

設定 plan 和 apply 階段，並加入手動審批步驟。

主要交付成果：

一個完整的 DevSecOps 流程，開發者只需專注撰寫程式碼，後續流程完全自動化。

Phase 5: 進階功能與迭代優化 (持續進行)
在穩固的基礎上，豐富產品功能並進行優化。

目標：實現商業模式、增強使用者體驗，並確保系統穩定。

步驟：

實現核心商業邏輯：

開發分級訂閱制與權限管理功能。

整合 Stripe 支付 API。

導入實驗性功能：

整合 LLM API，開發報告生成功能。

探索與 IBM Quantum Lab 的串接。

監控與日誌 (Monitoring & Logging)：

導入 AWS CloudWatch Logs 和 Alarms。

建立儀表板監控系統健康狀況 (CPU、記憶體、API 延遲等)。

安全性強化與優化：

設定 AWS WAF (Web Application Firewall)。

定期進行資料庫效能調校、分析雲端費用並進行優化。

主要交付成果：一個功能完整、穩定、安全且可持續營運的線上產品。
