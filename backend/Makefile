# 資安情報平台後端 Makefile

# 設定變數
BINARY_NAME=security-intel-platform
MAIN_PATH=cmd/server/main.go
PKG_PATH=./...

# 專案設定
APP_NAME = security-intelligence-platform
VERSION ?= latest
DOCKER_IMAGE = $(APP_NAME):$(VERSION)

# 開發工具
GOLANGCI_LINT_VERSION = v1.55.2

# 資料庫設定
DB_URL = "postgres://postgres:password@localhost:5432/security_intelligence?sslmode=disable"
MIGRATION_PATH = "./database/migrations"

# protobuf和gRPC設定
PROTO_PATH = api/proto
PROTO_OUT_PATH = api/proto

# 預設目標
.PHONY: help
help: ## 顯示幫助資訊
	@echo "可用命令："
	@echo "  開發工作流程:"
	@echo "    dev              - 快速開始開發環境"
	@echo "    run              - 執行應用程式"
	@echo "    watch            - 熱重載開發"
	@echo ""
	@echo "  建置和測試:"
	@echo "    build            - 建置應用程式"
	@echo "    test             - 執行測試"
	@echo "    coverage         - 測試覆蓋率"
	@echo ""
	@echo "  Proto和gRPC:"
	@echo "    proto            - 生成protobuf檔案"
	@echo "    install-proto-tools - 安裝protobuf工具"
	@echo ""
	@echo "  程式碼品質:"
	@echo "    fmt              - 格式化程式碼"
	@echo "    lint             - 程式碼檢查"
	@echo "    fix              - 修復問題"
	@echo ""
	@echo "  資料庫:"
	@echo "    migrate-up       - 執行資料庫遷移"
	@echo "    migrate-down     - 回滾資料庫遷移"
	@echo "    migrate-create   - 創建新遷移"
	@echo ""
	@echo "  Docker:"
	@echo "    docker-up        - 啟動開發環境"
	@echo "    docker-down      - 停止開發環境"
	@echo ""
	@echo "  工具:"
	@echo "    install-tools    - 安裝開發工具"
	@echo "    swagger          - 生成API文件"
	@echo "    help             - 顯示此幫助"

.DEFAULT_GOAL := help

# 開發相關
.PHONY: run
run: ## 啟動開發伺服器
	@echo "啟動開發伺服器..."
	go run $(MAIN_PATH)

.PHONY: build
build: ## 建置二進位檔案
	@echo "建置二進位檔案..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: clean
clean: ## 清理建置檔案
	@echo "清理建置檔案..."
	go clean
	rm -rf bin/

# 依賴管理
.PHONY: deps
deps: ## 安裝依賴套件
	@echo "安裝依賴套件..."
	go mod tidy
	go mod download

.PHONY: deps-update
deps-update: ## 更新依賴套件
	@echo "更新依賴套件..."
	go get -u ./...
	go mod tidy

# 測試相關
.PHONY: test
test: ## 執行測試
	@echo "執行測試..."
	go test $(PKG_PATH) -v

.PHONY: test-race
test-race: ## 執行競爭條件測試
	@echo "執行競爭條件測試..."
	go test $(PKG_PATH) -race -v

.PHONY: test-coverage
test-coverage: ## 生成測試覆蓋率報告
	@echo "生成測試覆蓋率報告..."
	go test $(PKG_PATH) -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆蓋率報告已生成: coverage.html"

.PHONY: benchmark
benchmark: ## 執行效能測試
	@echo "執行效能測試..."
	go test $(PKG_PATH) -bench=. -benchmem

# 程式碼品質
.PHONY: lint
lint: ## 執行程式碼檢查
	@echo "執行程式碼檢查..."
	golangci-lint run

.PHONY: fmt
fmt: ## 格式化程式碼
	@echo "格式化程式碼..."
	go fmt $(PKG_PATH)

.PHONY: vet
vet: ## 執行 go vet 檢查
	@echo "執行 go vet 檢查..."
	go vet $(PKG_PATH)

# 資料庫相關
.PHONY: migrate-up
migrate-up: ## 執行資料庫遷移 (up)
	@echo "執行資料庫遷移..."
	migrate -path $(MIGRATION_PATH) -database $(DB_URL) up

.PHONY: migrate-down
migrate-down: ## 執行資料庫遷移 (down)
	@echo "回滾資料庫遷移..."
	migrate -path $(MIGRATION_PATH) -database $(DB_URL) down

.PHONY: migrate-create
migrate-create: ## 建立新的遷移檔案 (使用: make migrate-create NAME=create_users_table)
	@echo "建立新的遷移檔案..."
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(NAME)

# API 文件
.PHONY: swagger
swagger: ## 生成 Swagger 文件
	@echo "生成 Swagger 文件..."
	swag init -g cmd/server/main.go -o docs/swagger

.PHONY: swagger-serve
swagger-serve: ## 啟動 Swagger UI 服務
	@echo "啟動 Swagger UI..."
	swagger serve -F=swagger docs/swagger/swagger.json

# Docker 相關
.PHONY: docker-build
docker-build: ## 建置 Docker 鏡像
	@echo "建置 Docker 鏡像..."
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: ## 運行 Docker 容器
	@echo "運行 Docker 容器..."
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

# 安全性
.PHONY: security-scan
security-scan: ## 執行安全性掃描
	@echo "執行安全性掃描..."
	gosec ./...

.PHONY: mod-verify
mod-verify: ## 驗證模組完整性
	@echo "驗證模組完整性..."
	go mod verify

# 開發工具安裝
.PHONY: install-tools
install-tools: ## 安裝開發工具
	@echo "安裝開發工具..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 全面檢查
.PHONY: check
check: fmt vet lint test ## 執行全面檢查

.PHONY: pre-commit
pre-commit: fmt vet lint test security-scan ## 提交前檢查

# 環境設定
.PHONY: env-setup
env-setup: ## 設定開發環境
	@echo "設定開發環境..."
	@if [ ! -f .env ]; then \
		echo "複製 .env.example 到 .env..."; \
		cp .env.example .env; \
		echo "請編輯 .env 檔案以設定您的環境變數"; \
	fi
	@echo "開發環境設定完成"

# 快速開始
.PHONY: quick-start
quick-start: deps install-tools env-setup ## 快速開始開發
	@echo "快速開始完成！"
	@echo "1. 編輯 .env 檔案設定資料庫連線"
	@echo "2. 執行 'make migrate-up' 來建立資料庫結構"
	@echo "3. 執行 'make run' 來啟動開發伺服器"

# 專案設定
APP_NAME = security-intelligence-platform
VERSION ?= latest
DOCKER_IMAGE = $(APP_NAME):$(VERSION)

# === 開發工作流程 ===

# 快速開始開發
.PHONY: dev
dev: docker-up migrate-up tidy run

# 安裝依賴
.PHONY: deps
deps:
	go mod download
	go mod verify

# 整理依賴
.PHONY: tidy
tidy:
	go mod tidy

# 執行應用程式
.PHONY: run
run:
	go run cmd/server/main.go

# 熱重載開發（需要安裝 air）
.PHONY: watch
watch:
	air

# === 建置和測試 ===

# 建置應用程式
.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(APP_NAME) cmd/server/main.go

# 交叉編譯
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux-amd64 cmd/server/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME)-darwin-amd64 cmd/server/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME)-windows-amd64.exe cmd/server/main.go

# 執行測試
.PHONY: test
test:
	go test -v -race -coverprofile=coverage.out ./...

# 測試覆蓋率
.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

# 基準測試
.PHONY: bench
bench:
	go test -bench=. -benchmem ./...

# === Proto 和 gRPC ===

# 安裝protobuf工具
.PHONY: install-proto-tools
install-proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成protobuf檔案
.PHONY: proto
proto:
	@echo "正在生成 protobuf 檔案..."
	protoc --go_out=$(PROTO_OUT_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_PATH) --go-grpc_opt=paths=source_relative \
		$(PROTO_PATH)/*.proto

# 清理生成的proto檔案
.PHONY: clean-proto
clean-proto:
	find $(PROTO_OUT_PATH) -name "*.pb.go" -delete

# === 程式碼品質 ===

# 程式碼格式化
.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w .

# 程式碼檢查
.PHONY: lint
lint:
	golangci-lint run

# 安裝linter
.PHONY: install-lint
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

# 修復可自動修復的問題
.PHONY: fix
fix:
	golangci-lint run --fix

# === 資料庫管理 ===

# 資料庫遷移 - 升級
.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATION_PATH) -database $(DB_URL) up

# 資料庫遷移 - 降級
.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATION_PATH) -database $(DB_URL) down

# 創建新的遷移檔案
.PHONY: migrate-create
migrate-create:
	@read -p "請輸入遷移檔案名稱: " name; \
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $$name

# 檢查遷移狀態
.PHONY: migrate-status
migrate-status:
	migrate -path $(MIGRATION_PATH) -database $(DB_URL) version

# 強制設定遷移版本
.PHONY: migrate-force
migrate-force:
	@read -p "請輸入版本號: " version; \
	migrate -path $(MIGRATION_PATH) -database $(DB_URL) force $$version

# === Docker 操作 ===

# 啟動開發環境
.PHONY: docker-up
docker-up:
	cd ../docker && docker-compose up -d

# 停止開發環境
.PHONY: docker-down
docker-down:
	cd ../docker && docker-compose down

# 重啟開發環境
.PHONY: docker-restart
docker-restart: docker-down docker-up

# 檢視日誌
.PHONY: docker-logs
docker-logs:
	cd ../docker && docker-compose logs -f

# 建置Docker映像
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# 推送Docker映像
.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE)

# === Swagger 文件 ===

# 生成Swagger文件
.PHONY: swagger
swagger:
	swag init -g cmd/server/main.go -o docs/swagger

# 安裝swagger工具
.PHONY: install-swagger
install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

# === 清理和維護 ===

# 清理建置檔案
.PHONY: clean
clean:
	rm -rf bin/
	rm -rf coverage.out coverage.html
	rm -rf docs/swagger

# 清理全部（包括依賴）
.PHONY: clean-all
clean-all: clean clean-proto
	go clean -modcache

# === 部署 ===

# 本地部署
.PHONY: deploy-local
deploy-local: build docker-build

# 生產環境部署
.PHONY: deploy-prod
deploy-prod: test build docker-build docker-push

# === 開發工具 ===

# 安裝開發工具
.PHONY: install-tools
install-tools: install-lint install-swagger install-proto-tools
	go install github.com/cosmtrek/air@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 檢查系統需求
.PHONY: check-requirements
check-requirements:
	@echo "檢查系統需求..."
	@command -v go >/dev/null 2>&1 || { echo "Go 未安裝"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "Docker 未安裝"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose 未安裝"; exit 1; }
	@command -v protoc >/dev/null 2>&1 || { echo "Protocol Buffers compiler 未安裝"; exit 1; }
	@echo "所有需求已滿足 ✓"

# === 測試資料 ===

# 載入測試資料
.PHONY: load-test-data
load-test-data:
	psql $(DB_URL) -f scripts/test_data.sql 