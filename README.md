<h1 align="center">資安情報平台 (Security Intelligence Platform)</h1>

<p align="center">
  一個創新、自動化且可擴展的開源資安威脅情報平台。
</p>

<p align="center">
  
  <a href="https://github.com/<your-repo>/blob/main/LICENSE"><img src="https://img.shields.io/github/license/<your-repo>/security-intel-platform?style=flat-square" alt="License"></a>
  <a href="https://github.com/<your-repo>/issues"><img src="https://img.shields.io/github/issues/<your-repo>/security-intel-platform?style=flat-square" alt="GitHub issues"></a>
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8.svg?style=flat-square" alt="Go Version">
  <img src="https://img.shields.io/badge/Next.js-14+-black.svg?style=flat-square&logo=next.js" alt="Next.js Version">
</p>

> 「資安情報與預測平台」專案的詳細文件，旨在提供一個清晰的概覽、技術細節和實施指南。

<br>

---

## 📖 目錄 (Table of Contents)

- [專案概述](#-專案概述-project-overview)
- [✨ 核心功能](#-核心功能-core-features)
- [🛠️ 技術棧](#️-技術棧-technology-stack)
- [🏗️ 優化的專案結構](#️-優化的專案結構-project-structure)
- [🚀 CI/CD 專業分工規劃](#-cicd-專業分工規劃)
- [⚡ 快速入門](#-快速入門-quick-start)
- [🤝 貢獻指南](#-貢獻指南-contributing)
- [📜 授權](#-授權-license)

---

## 📝 專案概述 (Project Overview)

「資安情報平台」是一個創新且自動化的資安威脅情報蒐集與分析平台。我們旨在建立一個易於部署、可擴展的解決方案，幫助個人、研究者乃至組織追蹤、理解並預測不斷演進的資安威 D 脅。本專案將實現多源情報整合、分級存取、多元支付選項，並實驗性地探索前沿技術如大型語言模型 (LLM) 報告生成及量子計算預測零日攻擊的潛力。

<details>
  <summary><strong>💡 補充知識：什麼是資安情報平台 (TIP)？</strong></summary>
  
  > 資安情報平台 (Threat Intelligence Platform, TIP) 是一種用於匯總、分析和分發威脅情報的技術解決方案。它的核心價值在於將來自不同來源（開源、商業、內部）的雜亂數據，轉化為具體、可操作的洞察，幫助組織：
  > 1.  **主動防禦**：在攻擊發生前識別潛在威脅。
  > 2.  **加速應變**：快速了解攻擊者的手法 (TTPs)，縮短調查時間。
  > 3.  **優化決策**：為資安投資和策略提供數據支持。
  >
  > 這個專案的目標就是打造一個現代化、自動化且可擴展的 TIP。
</details>

## ✨ 核心功能 (Core Features)

| 功能分類                                  | 描述                                                                                                                      |
| :---------------------------------------- | :------------------------------------------------------------------------------------------------------------------------ |
| **🤖 自動化情報蒐集**                     | 定期從多個公開來源（如 AbuseIPDB、Have I Been Pwned、NVD/CVE）抓取惡意 IP、漏洞資訊、數據泄露事件等，並進行清洗與正規化。 |
| **🔍 靈活的瀏覽與搜尋**                   | 提供直觀的儀表板與強大的搜尋篩選功能，快速定位特定威脅情報。                                                              |
| **🧩 階層式存取訂閱制提供免費層與付費層** | ，以 Freemium 模式滿足不同用戶需求，解鎖更即時、更深入的數據。                                                            |
| **💳 多元支付選項**                       | 支援傳統法幣 (Stripe) 與實驗性的穩定幣支付，提供彈性。                                                                    |
| **💳 多元支付選項**                       | 支援傳統法幣 (透過 Stripe) 與創新的穩定幣支付，提供全球用戶彈性的付款選擇。                                               |
| **✍️ 智慧報告生成 (實驗性)**              | 利用大型語言模型 (LLM) 自動生成情報摘要、分析報告或簡報草稿，提升分析效率。                                               |
| **🔮 量子預測 (概念驗證)**                | 探索在 IBM Quantum Lab 上運行量子演算法，預測潛在的零日攻擊模式。                                                         |
| **🔐 數據泄露監控**                       | 整合 Have I Been Pwned API，監控帳戶泄露、密碼安全性和域名安全狀態。                                                      |
| **🧩 可擴展與自動化**                     | 基於 Go + Next.js 的微服務友好架構，並透過 Terraform 與 CI/CD 流程實現完全自動化。                                        |

## 🛠️ 技術棧 (Technology Stack)

| 類別                  | 技術                                           | 選擇原因與說明                                                              |
| :-------------------- | :--------------------------------------------- | :-------------------------------------------------------------------------- |
| **後端 (Backend)**    | **Go (Golang)**                                | 高併發性能、靜態類型安全、部署簡單，非常適合網路服務與 API 開發。           |
|                       | Gin Gonic                                      | 一個輕量級、高效能的 Go Web 框架，路由和中介層設計直觀。                    |
| **前端 (Frontend)**   | **Next.js (React)**                            | 提供伺服器端渲染 (SSR) 和 App Router，有利於 SEO 和開發體驗。               |
|                       | Tailwind CSS                                   | Utility-First 的 CSS 框架，能快速建構客製化 UI，且易於維護。                |
| **資料庫 (Database)** | **PostgreSQL**                                 | 功能強大、穩定可靠的開源關聯式資料庫，支援複雜查詢和 JSONB 等多種資料類型。 |
|                       | AWS RDS / OCI Free DB                          | 託管式資料庫服務，簡化了維運、備份和擴展的複雜性。                          |
| **基礎設施 (IaC)**    | **Terraform**                                  | 用程式碼來定義和管理雲端資源，實現可重複、可追蹤的基礎設施部署。            |
| **CI/CD**             | **GitHub Actions, Drone CI, AWS CodePipeline** | 專業分工，結合各平台優勢，打造穩定、安全的自動化流程。                      |
|                       | Amazon ECR                                     | AWS 託管的 Docker 容器註冊表，安全可靠。                                    |
|                       | Vercel                                         | Next.js 的最佳部署平台，提供全球 CDN 和自動化的 CI/CD 流程。                |

HTTP REST API：統一的 JSON 回應格式，完整的 Swagger 文檔
gRPC API：高效能的二進制協議，支援串流
MQTT 通訊：即時威脅情報推送和訂閱
JWT 安全：無狀態認證，支援 Bearer Token
HIBP API 整合：完整的 Have I Been Pwned API v3 支援，包含帳戶泄露檢查、密碼安全驗證、域名監控
測試自動化：Robot Framework 覆蓋所有 API 端點

## 🏗️ 優化的專案結構 (Project Structure)

本專案採用前後端分離的 Monorepo 結構，以實現清晰的職責劃分和高效的開發流程。

````plaintext
security-intel-platform/
├── .github/                    # GitHub 相關配置 (CI/CD, Issue Templates)
├── backend/
├── quantum-service/            # IBM Quantum Python 模組
+ │   ├── circuits/               # 量子電路定義
+ │   │   └── pattern_prediction.py
+ │   ├── tests/                  # Python 測試
+ │   │   └── test_circuits.py
+ │   ├── main.py                 # 程式主進入點 (可被後端觸發)
+ │   ├── requirements.txt        # Python 依賴項 (qiskit, numpy)
+ │   └── Dockerfile              # 用於打包 Python 環境的 Dockerfile                  # Golang 後端服務 (Clean Architecture)
├── frontend/                   # Next.js 前端應用 (App Router)
├── terraform/                  # 基礎設施即程式碼 (IaC)
├── docs/                       # 專案文件 (架構, 指南)
├── docker/                     # 本地開發用的 Docker Compose
├── .drone.yml                  # Drone.io CI 配置
├── .gitignore
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── LICENSE
└── README.md                   # 本文件

## 系統架構圖
```mermaid
graph TD
    subgraph User Interaction
        U[User/Client] --> FE(Next.js on Vercel)
    end

    subgraph CI/CD Automation
        G[GitHub] --PR/Push--> GA[GitHub Actions]
        GA --Lint/Test--> G
        GA --Deploy--> FE
        GA --On Success--> D[Drone.io]
        D --Build & Push Image--> ECR[AWS ECR]
        ECR --New Image--> CP[AWS CodePipeline]
        G --Terraform Change--> CP_TF[AWS CodePipeline for TF]
    end

    subgraph AWS Cloud
        FE --> BE[API: Go on EC2/ECS]
        BE --> DB[(PostgreSQL on RDS)]

        CP --Deploy--> BE
        CP_TF --Apply--> AWS_INFRA[Managed AWS Resources]

        subgraph External Services
            BE --> S3[Stripe API]
            BE --> LLM[LLM API]
            BE --> QC[IBM Quantum API]
        end
    end

    subgraph Data Collection
        Collector[Cron: GitHub Actions] --> BE
    end
````

## 🔐 Have I Been Pwned (HIBP) API 整合

本平台完整整合了 Have I Been Pwned API v3，提供全面的數據泄露監控功能：

### 支援的 HIBP 功能

| 功能類別           | 端點                                          | 描述                                 |
| ------------------ | --------------------------------------------- | ------------------------------------ |
| **帳戶泄露檢查**   | `GET /api/v1/hibp/account/{account}/breaches` | 檢查特定帳戶是否在已知數據泄露事件中 |
| **Paste 記錄查詢** | `GET /api/v1/hibp/account/{account}/pastes`   | 查詢帳戶在 Paste 網站上的記錄        |
| **域名泄露監控**   | `GET /api/v1/hibp/domain/{domain}/breaches`   | 監控特定域名下的郵箱泄露情況         |
| **密碼安全檢查**   | `GET /api/v1/hibp/password/check`             | 檢查密碼是否已被泄露（無需認證）     |
| **泄露事件查詢**   | `GET /api/v1/hibp/breaches`                   | 獲取所有已知的數據泄露事件           |
| **最新泄露事件**   | `GET /api/v1/hibp/breach/latest`              | 獲取最新添加的泄露事件               |
| **Stealer Logs**   | `GET /api/v1/hibp/stealer/*`                  | 查詢竊取器日誌（需要高級訂閱）       |
| **訂閱狀態**       | `GET /api/v1/hibp/subscription/status`        | 查詢 HIBP 訂閱狀態和配額             |

### 前端功能

- **帳戶泄露檢查器**：輸入電子郵件地址檢查是否在泄露事件中
- **密碼安全驗證器**：檢查密碼是否已被泄露
- **域名監控器**：監控特定域名的安全狀態
- **實時安全建議**：根據檢查結果提供個性化安全建議

### 配置要求

在環境變數中設置 HIBP API 金鑰：

```bash
HIBP_API_KEY=your_hibp_api_key_here
```

### 使用範例

```bash
# 檢查帳戶泄露
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/hibp/account/user@example.com/breaches"

# 檢查密碼安全性
curl "http://localhost:8080/api/v1/hibp/password/check?password=yourpassword"

# 檢查域名泄露
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     "http://localhost:8080/api/v1/hibp/domain/example.com/breaches"
```

## 🚀 CI/CD 專業分工規劃

在本專案中，我們結合使用 GitHub Actions、Drone CI 和 AWS CodePipeline，讓它們各司其職，建立一個分工合作的專業級自動化流程。

### 職責分工總覽

| 工具                 | 在本專案中的主要職責                       | 觸發時機 / 核心優勢                                                                                                            |
| :------------------- | :----------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------- |
| **GitHub Actions**   | **程式碼品質檢查、前端部署、觸發下游流程** | **觸發**: `git push` 到任何分支, 創建 `PR` 時<br>**優勢**: 與 GitHub 無縫整合、快速回饋、生態系龐大、排程任務                  |
| **Drone CI**         | **後端應用程式的建置、打包、容器化**       | **觸發**: 在 GitHub Actions 檢查通過後，由 Webhook 觸發<br>**優勢**: 自我託管、Docker 原生、配置簡潔                           |
| **AWS CodePipeline** | **後端部署、基礎設施 (Terraform) 變更**    | **觸發**: ECR 有新鏡像時或 `terraform/` 目錄有變更時<br>**優勢**: 與 AWS IAM 深度整合、原生支援 AWS 服務、可視化流程與手動審批 |

quantum-service 目錄詳解
circuits/: 存放所有 Qiskit 量子電路的定義。將電路邏輯與主程式分離，有利於測試和管理。

tests/: 使用 pytest 框架編寫的單元測試。這些測試主要針對本地模擬器運行，以驗證電路建構和資料處理的正確性，而不會真的呼叫 IBM 的硬體。

main.py: 腳本的主進入點。它負責接收參數（例如，從 Go 後端傳來的情資數據），載入電路，連接到 IBM Quantum 服務，執行任務，並將結果寫回資料庫或日誌。

requirements.txt: 列出所有 Python 依賴項，例如 qiskit, qiskit-ibm-provider, numpy 等。

Dockerfile: 一個專為此 Python 環境設計的 Dockerfile。它會安裝所有依賴項，並將程式碼打包成一個可執行的容器。

後端如何與 Python 模組互動？
Go 後端 (backend) 不會直接執行 Python 程式碼。它會透過以下方式觸發這個打包好的 quantum-service 容器來執行任務：

準備數據: Go 後端將需要分析的情資數據整理好，存入資料庫中的一個特定任務表，狀態為「待處理」。

觸發執行: Go 後端可以透過 AWS SDK 觸發一個 AWS Batch 任務或 ECS Task，該任務就是運行我們打包好的 quantum-service Docker 鏡像。

執行與回寫: quantum-service 容器啟動後，從資料庫讀取「待處理」的任務，連接 IBM Quantum Cloud 執行計算，最後將結果寫回資料庫，並更新任務狀態為「已完成」。

這種非同步 (asynchronous) 的設計模式非常適合耗時較長的計算任務，不會阻塞主後端服務。

Part 2: 規劃量子計算的 CI/CD 流程
量子計算模組的 CI/CD 流程與 Web 應用不同，它更專注於程式碼的正確性、可重複性和打包。這個流程完全可以在 GitHub Actions 中完成。

我們可以在 .github/workflows/ 目錄下新增一個檔案 quantum-ci.yml。

quantum-ci.yml 工作流程詳解
觸發條件:

YAML

on:
push:
paths: - 'quantum-service/**'
pull_request:
paths: - 'quantum-service/**'
只有當 quantum-service/ 目錄下的程式碼有變動時，才觸發此工作流程，以節省資源。

Job 步驟:

設定 Python 環境:

YAML

- name: Set up Python
  uses: actions/setup-python@v4
  with:
  python-version: '3.10'
  安裝依賴項:

YAML

- name: Install dependencies
  run: |
  python -m pip install --upgrade pip
  pip install -r quantum-service/requirements.txt
  程式碼風格檢查 (Linting):

YAML

- name: Lint with flake8
  run: |
  pip install flake8
  flake8 quantum-service/ --count --show-source --statistics
  單元測試與模擬器測試 (Unit & Simulation Tests):

這是最關鍵的一步。測試不應該連接到真實的 IBM 量子電腦。

使用 pytest 和 Qiskit 內建的本地模擬器 (AerSimulator) 來驗證電路是否能成功運行並產出預期格式的結果。

YAML

- name: Test with pytest
  run: |
  pip install pytest
  pytest quantum-service/tests/
  建置並推送 Docker 鏡像 (CD 部分):

和 Drone CI 為後端做的事情一樣，當程式碼合併到 main 分支時，我們將 Python 環境打包成 Docker 鏡像並推送到 ECR。

YAML

- name: Build and push Docker image
  if: github.ref == 'refs/heads/main' && github.event_name == 'push'
  uses: docker/build-push-action@v4
  with:
  context: ./quantum-service
  push: true
  tags: your-aws-account-id.dkr.ecr.your-region.amazonaws.com/quantum-service:latest

### 各平台詳細任務規劃

<details>
  <summary><strong>點此展開詳細的 CI/CD 任務規劃</strong></summary>

#### 1. GitHub Actions：第一道防線與快速反應者

在 `.github/workflows/` 目錄下，設定以下工作流程：

- **`code-quality.yml` (在每次 `push` 和 `PR` 時觸發)**
  - **Linter 檢查**：對 `backend/` 運行 `golangci-lint`；對 `frontend/` 運行 `ESLint`。
  - **單元測試**：對 `backend/` 運行 `go test ./...`；對 `frontend/` 運行 `npm test`。
  - **安全性掃描**：使用 `Trivy` 或 `CodeQL` 掃描程式碼。
  - **Terraform 驗證**：對 `terraform/` 目錄運行 `terraform validate`。
- **`deploy-frontend.yml` (僅在 `push` 到 `main` 分支時觸發)**
  - **觸發 Vercel 部署**：利用 Vercel 官方整合自動部署前端。
- **`data-collection.yml` (`schedule` 定時觸發)**
  - **執行排程任務**：定時呼叫後端 API，啟動情報蒐集。
- **`trigger-downstream.yml` (`code-quality.yml` 成功後觸發)**
  - **觸發 Drone CI**：使用 Webhook 向 Drone CI 平台發送請求，啟動後端建置。

#### 2. Drone CI：後端建置與容器化專家

在根目錄的 `.drone.yml` 中，定義後端建置流程：

- **Pipeline 步驟**：
  1. **Clone**：複製最新的程式碼。
  2. **Build Binary & Docker Image**：編譯 Go 應用，並基於 `backend/Dockerfile` 建置 Docker 鏡像。
  3. **Push to ECR**：將標記好版本 (`Git SHA`) 的鏡像推送到 Amazon ECR。

#### 3. AWS CodePipeline：雲端部署總指揮官

在 AWS Console 中，設定兩個主要的 Pipeline：

- **後端應用程式部署 (`backend-deploy-pipeline`)**
  1. **Source**: 監聽 Amazon ECR 的新鏡像。
  2. **Deploy**: 使用 AWS CodeDeploy 將新鏡像安全地部署到 EC2/ECS。
- **基礎設施部署 (`infra-terraform-pipeline`)**
  1. **Source**: 監聽 `terraform/` 目錄的變更。
  2. **Build**: 使用 AWS CodeBuild 運行 `terraform plan`。
  3. **Approval**: 加入手動審批階段，人工確認變更。
  4. **Deploy**: 使用 AWS CodeBuild 運行 `terraform apply`。

</details>

### 總結流程

> 開發者推送程式碼 ➔ **GitHub Actions** (品質檢查) ➔ (若後端變更) **Drone CI** (建置鏡像) ➔ **AWS CodePipeline** (部署應用)

---

## ⚡ 快速入門 (Quick Start)

<details>
  <summary><strong>點此展開本地開發環境設置指南</strong></summary>

### 前提條件 (Prerequisites)

- Git, Docker & Docker Compose
- Go (>= 1.22)
- Node.js (>= 18)
- Terraform (>= 1.5)
- AWS CLI (已配置憑證)

### 本地開發環境設置 (Local Development Setup)

1.  **複製儲存庫**：

    ```bash
    git clone [https://github.com/](https://github.com/)<your-repo>/security-intel-platform.git
    cd security-intel-platform
    ```

2.  **配置環境變數**：

    > ⚠️ **安全警告**：切勿將包含 API Key、資料庫密碼等敏感資訊的 `.env` 檔案提交到 Git 儲存庫！請務必將其加入 `.gitignore` 檔案中。

    參考 `backend/.env.example` 和 `frontend/.env.example` (需自行創建) 來建立你的 `.env` 文件。

3.  **啟動依賴服務 (Docker)**：

    ```bash
    docker-compose -f docker/docker-compose.yml up -d
    ```

4.  **運行後端服務** (使用 Makefile)：

    ```bash
    cd backend
    make run
    ```

    > ℹ️ 後端預設會在 `http://localhost:8080` 運行。

5.  **運行前端應用**：
`bash
      cd frontend
      npm install
      npm run dev
      ` > ℹ️ 前端預設會在 `http://localhost:3000` 運行。
</details>

---

## 🤝 貢獻指南 (Contributing)

我們非常歡迎社群的貢獻！若您有興趣，請參考 `CONTRIBUTING.md` 文件了解詳細的貢獻流程，包括如何提出 Issue 和提交 Pull Request。

請務必遵守我們的 **行為準則 (`CODE_OF_CONDUCT.md`)**。

## 📜 授權 (License)

本專案根據 **MIT License** 發佈。詳情請參閱 `LICENSE` 文件。
