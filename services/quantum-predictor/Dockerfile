# 資安情報平台 - 量子預測服務 Dockerfile
# 多階段建置，最小化最終鏡像體積並提高安全性

# -----------------------------------------------------------------------------
# 第一階段：建置階段
# -----------------------------------------------------------------------------
FROM python:3.10-slim AS builder

# 設定工作目錄
WORKDIR /app

# 安裝編譯依賴
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential \
        gcc \
        git \
        && rm -rf /var/lib/apt/lists/*

# 複製需求檔案
COPY requirements.txt .

# 建立虛擬環境並安裝依賴
RUN python -m venv /venv && \
    /venv/bin/pip install --upgrade pip && \
    /venv/bin/pip install --no-cache-dir -r requirements.txt && \
    /venv/bin/pip install --no-cache-dir gunicorn

# -----------------------------------------------------------------------------
# 第二階段：最終執行階段
# -----------------------------------------------------------------------------
FROM python:3.10-slim

# 設定標籤
LABEL maintainer="資安情報平台開發團隊"
LABEL description="量子預測服務 - 用於資安威脅模式分析"

# 設定環境變數
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PATH="/venv/bin:$PATH" \
    USE_REAL_DEVICE=false \
    MODEL_DIR="/app/models" \
    PYTHONPATH="/app"

# 安裝運行時所需的系統依賴
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        tini \
        && rm -rf /var/lib/apt/lists/*

# 建立非特權使用者
RUN addgroup --gid 1000 appuser && \
    adduser --uid 1000 --gid 1000 --disabled-password --gecos "" appuser

# 從建置階段複製虛擬環境
COPY --from=builder /venv /venv

# 設定工作目錄
WORKDIR /app

# 複製應用程式檔案
COPY main.py circuit_designer.py ./
COPY models/ ./models/

# 建立必要的目錄
RUN mkdir -p /app/data /app/logs && \
    chown -R appuser:appuser /app

# 切換到非特權使用者
USER appuser

# 設定容器入口點為 tini 以便正確處理信號
ENTRYPOINT ["/usr/bin/tini", "--"]

# 設定容器執行指令
CMD ["python", "main.py"]

# 設定健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=15s --retries=3 \
    CMD python -c "import os; exit(0 if os.path.exists('/app/main.py') else 1)"

# 說明如何使用此容器的說明
# 
# 範例用法:
# 1. 基本執行:
#    docker run -it --rm security-intel/quantum-predictor
# 
# 2. 處理特定輸入檔案:
#    docker run -it --rm -v $(pwd)/data:/app/data security-intel/quantum-predictor \
#      python main.py --input /app/data/input.json --output /app/data/result.json
# 
# 3. 使用真實量子設備:
#    docker run -it --rm -e USE_REAL_DEVICE=true -e IBMQ_API_KEY=YOUR_API_KEY \
#      security-intel/quantum-predictor
