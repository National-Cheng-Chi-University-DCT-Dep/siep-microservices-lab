gantt
title 資安情報平台 - 免費增值極致派整合計畫
dateFormat YYYY-MM-DD
axisFormat %m/%d
section Phase 1: 基礎建設與本地設定
初始化 Git Repo 與 Monorepo 架構 :done, 2025-09-01, 1d
開發 "Hello World" 前後端服務 :done, 2025-09-02, 1d
建立本地 Docker Compose 環境 :active, 2025-09-03, 1d
撰寫後端 Dockerfile :2025-09-04, 1d
確認本地前後端串接 :2025-09-05, 1d
設定本地 .env 環境變數 :2025-09-06, 1d
完成 .gitignore 設定 :2025-09-07, 1d
撰寫初步的 README.md :2025-09-08, 1d
建立 GitHub Issue/PR 模板 :2025-09-09, 1d
完成 Phase 1 內部審核 :2025-09-10, 1d

    section Phase 2: 前端部署 (Vercel)
    建立 Vercel 帳號並連結 GitHub     :2025-09-11, 1d
    建立 Vercel 專案並指向 frontend/  :2025-09-12, 1d
    完成首次部署並驗證上線            :2025-09-13, 1d
    設定 Vercel 環境變數 (API URL 佔位) :2025-09-14, 1d
    測試 Vercel 預覽部署功能         :2025-09-15, 1d
    建立基本 UI 佈局 (Layout)        :2025-09-16, 2d
    完成主要頁面的基本路由            :2025-09-18, 1d
    串接一個公開免費 API 測試        :2025-09-19, 1d
    設計響應式 (RWD) 基礎         :2025-09-20, 1d
    完成 Phase 2 UI/UX 審核        :2025-09-21, 1d

    section Phase 3: 資料庫與認證 (Supabase)
    建立 Supabase 專案             :2025-09-22, 1d
    設計並建立核心資料表結構          :2025-09-23, 2d
    啟用 Supabase Auth 並設定提供者  :2025-09-25, 1d
    取得 DB 連線字串與 API 金鑰      :2025-09-26, 1d
    將金鑰儲存於 Vercel 環境變數     :2025-09-27, 1d
    在前端建立登入/註冊 UI 元件       :2025-09-28, 2d
    整合 Supabase.js Client       :2025-09-30, 1d
    測試前端使用者註冊流程            :2025-10-01, 1d
    實作前端受保護路由              :2025-10-02, 1d
    完成 Phase 3 認證流程測試        :2025-10-03, 1d

    section Phase 4: 後端核心邏輯開發 (本地)
    整合 Go 的 PostgreSQL 驅動     :2025-10-04, 1d
    連接本地 Go 後端至遠端 Supabase :2025-10-05, 1d
    建立符合 DB Schema 的 Go struct :2025-10-06, 1d
    實作 JWT 驗證中介層            :2025-10-07, 2d
    開發核心功能的 CRUD API         :2025-10-09, 2d
    撰寫核心邏輯的單元測試            :2025-10-11, 1d
    優化 Dockerfile (多階段建置)    :2025-10-12, 1d
    實作 API 回應的統一格式         :2025-10-13, 1d
    加入 Gin 的 Logger 與 Recovery  :2025-10-14, 1d
    完成 Phase 4 本地功能驗證        :2025-10-15, 1d

    section Phase 5: 後端部署與整合 (Render)
    建立 Render 帳號並連結 GitHub     :2025-10-16, 1d
    建立 Render Web Service         :2025-10-17, 1d
    設定 Render 環境變數 (DB 連線等) :2025-10-18, 1d
    完成後端首次部署並驗證          :2025-10-19, 1d
    更新 Vercel 的 API URL         :2025-10-20, 1d
    在後端設定 CORS 策略            :2025-10-21, 1d
    重新部署後端並測試 CORS         :2025-10-22, 1d
    測試完整串接 (FE->BE->DB)        :2025-10-23, 1d
    壓力測試休眠喚醒時間              :2025-10-24, 1d
    完成 Phase 5 整合測試           :2025-10-25, 1d

    section Phase 6: 核心功能整合 (HIBP)
    取得 HIBP API 金鑰              :2025-10-26, 1d
    將金鑰儲存於 Render Secrets       :2025-10-27, 1d
    開發後端 HIBP API 服務層       :2025-10-28, 2d
    建立後端 HIBP 相關路由          :2025-10-30, 1d
    實作簡易後端快取機制              :2025-10-31, 1d
    建立前端 HIBP 查詢 UI          :2025-11-01, 2d
    串接前端 UI 與後端 HIBP API    :2025-11-03, 1d
    設計結果的視覺化呈現              :2025-11-04, 1d
    部署前後端更新                  :2025-11-05, 1d
    完成 Phase 6 端對端功能測試      :2025-11-06, 1d

    section Phase 7: 自動化排程任務 (GitHub Actions)
    在後端建立情報蒐集的主要邏輯      :2025-11-07, 2d
    建立一個用金鑰保護的觸發 API      :2025-11-09, 1d
    在 GitHub Repo 設定金鑰 Secret  :2025-11-10, 1d
    撰寫 collector.yml 工作流程檔案   :2025-11-11, 1d
    設定 cron 排程規則              :2025-11-12, 1d
    撰寫 curl 觸發後端的腳本         :2025-11-13, 1d
    手動觸發 workflow 進行測試       :2025-11-14, 1d
    檢查後端日誌確認任務執行          :2025-11-15, 1d
    觀察排程自動執行狀況              :2025-11-16, 1d
    完成 Phase 7 自動化流程驗證      :2025-11-17, 1d

    section Phase 8: 實驗性功能 (LLM 整合)
    選擇輕量級 LLM 模型             :2025-11-18, 1d
    建立 Hugging Face Space 專案    :2025-11-19, 1d
    用 Gradio/Streamlit 撰寫推理腳本 :2025-11-20, 2d
    部署 Space 並取得 API 端點      :2025-11-22, 1d
    在後端建立呼叫 LLM 的服務       :2025-11-23, 1d
    建立後端報告生成 API              :2025-11-24, 1d
    建立前端觸發報告生成的 UI         :2025-11-25, 1d
    串接前端 UI 與後端報告 API       :2025-11-26, 1d
    部署所有更新                    :2025-11-27, 1d
    完成 Phase 8 LLM 功能端對端測試  :2025-11-28, 1d

    section Phase 9: 實驗性功能 (量子服務)
    完善 quantum-service Dockerfile :2025-11-29, 1d
    確認量子任務的輸入與輸出          :2025-11-30, 1d
    設計觸發機制 (手動 GitHub Action) :2025-12-01, 1d
    建立 quantum-runner.yml 工作流程 :2025-12-02, 1d
    設定 workflow_dispatch 手動觸發  :2025-12-03, 1d
    撰寫 workflow 腳本來建置並運行 Docker :2025-12-04, 2d
    測試將結果寫回 Supabase DB      :2025-12-06, 1d
    建立前端頁面來展示量子計算結果    :2025-12-07, 1d
    手動觸發並驗證完整流程            :2025-12-08, 1d
    完成 Phase 9 量子 POC 驗證       :2025-12-09, 1d

    section Phase 10: 最終化、監控與文件化
    進行全面的程式碼審查 (Code Review) :2025-12-10, 1d
    設定 UptimeRobot 進行服務監控    :2025-12-11, 1d
    啟用 Vercel Analytics           :2025-12-12, 1d
    全面檢視與改善 UI/UX             :2025-12-13, 2d
    撰寫使用者操作手冊                :2025-12-15, 1d
    更新 README.md 的架構與部署連結 :2025-12-16, 1d
    清理程式碼與註解                  :2025-12-17, 1d
    檢查所有 Secrets 與 Keys 的安全性 :2025-12-18, 1d
    凍結 v1.0 功能並建立 Tag        :2025-12-19, 1d
    v1.0 專案正式上線！              :2025-12-20, 1d

```mermaid
graph TD
    subgraph A. 開發實驗室 (Kaggle / Lightning AI)
        A1[💡 點子發想] --> A2[🧪 在 Notebook 中<br>測試不同 LLM/量子演算法];
        A2 --> A3[🔧 微調模型<br>與優化程式碼];
        A3 --> A4[✅ 找到最終版的<br>模型與程式];
    end

    subgraph B. 版本控制 (GitHub)
        A4 --> B1[將最終程式碼<br>Push 到 Repo];
    end

    subgraph C. 線上服務站 (Hugging Face Spaces)
        B1 --> C1[Hugging Face 從 GitHub<br>自動拉取程式碼];
        C1 --> C2[🚀 部署成一個<br>帶有 API 的線上服務];
    end

    subgraph D. 主應用 (Vercel + Render)
        D1[Go 後端 on Render] -- 呼叫 API --> C2;
    end

```

aggle Notebooks： 首選的實驗場！ 每週有 30 小時的免費 GPU/TPU 時間，非常慷慨。你可以在這裡盡情地測試不同的 LLM 模型，看看哪個生成的報告品質最好，或是運行你的 quantum-service 腳本來驗證演算法。

Lightning AI Studio / AWS SageMaker Studio Lab： 專業級的開發環境！ 它們提供更像 IDE 的體驗，適合比較複雜的專案。當你需要持久化的儲存和更完整的開發工作流時，它們就是絕佳選擇。你可以在這裡把模型程式碼寫得更工程化、更完美。

Paperspace Gradient： 備用的強力工作站。 當其他平台的 GPU 類型不符合你的需求時，可以來這裡找找看。
