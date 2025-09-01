# 量子預測服務 (Quantum Predictor Service)

這是資安情報平台的量子預測服務，使用量子計算來分析威脅情報數據。

## 功能特色

- 🔬 **量子模式識別**: 使用變分量子分類器 (VQC) 進行威脅模式分析
- 🚀 **非同步處理**: 支援長時間運行的量子計算任務
- 🔄 **多後端支援**: 支援 IBM Quantum 真實設備和本地模擬器
- 📊 **視覺化結果**: 提供詳細的量子測量結果和機率分布圖表
- 🐳 **容器化部署**: 完整的 Docker 支援

## 快速開始

### 使用 Docker (推薦)

```bash
# 從 GitHub Container Registry 拉取鏡像
docker pull ghcr.io/lipeichen/Ultimate-Security-Intelligence-Platform/quantum-predictor:latest

# 運行容器
docker run -it --rm \
  -v $(pwd)/input.json:/app/input.json \
  -v $(pwd)/output.json:/app/output.json \
  ghcr.io/lipeichen/Ultimate-Security-Intelligence-Platform/quantum-predictor:latest \
  python main.py --input /app/input.json --output /app/output.json
```

### 本地開發

```bash
# 安裝依賴
pip install -r requirements.txt

# 設定環境變數
export IBMQ_API_KEY="your_ibm_quantum_api_key"
export USE_REAL_DEVICE="false"  # 使用模擬器

# 運行服務
python main.py --input input.json --output result.json
```

## 輸入格式

服務接受 JSON 格式的輸入文件：

```json
{
  "threats": [
    {
      "ip_address": "192.168.1.100",
      "threat_type": "malware",
      "risk_score": 85,
      "country": "US",
      "attack_type": "brute_force",
      "timestamp": "2024-01-15T10:30:00Z"
    }
  ],
  "analysis_params": {
    "use_simulator": true,
    "shots": 1024,
    "confidence_threshold": 0.7
  }
}
```

## 輸出格式

服務產生 JSON 格式的結果：

```json
{
  "prediction": 1,
  "probability": 0.85,
  "confidence": 75.5,
  "is_malicious": true,
  "counts": {
    "0000": 123,
    "0001": 456,
    "0010": 234,
    "0011": 211
  },
  "backend": "ibmq_qasm_simulator",
  "execution_time": 45,
  "timestamp": "2024-01-15T10:35:00Z"
}
```

## 環境變數

| 變數名 | 描述 | 預設值 |
|--------|------|--------|
| `IBMQ_API_KEY` | IBM Quantum API 金鑰 | - |
| `USE_REAL_DEVICE` | 是否使用真實量子設備 | `false` |
| `MODEL_DIR` | 模型文件目錄 | `/app/models` |
| `DEFAULT_MODEL` | 預設模型文件名 | `quantum_model_params.json` |

## API 端點

### 提交量子任務

```http
POST /api/v1/quantum-jobs
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "DDoS攻擊模式分析",
  "description": "分析近期網路流量中的 DDoS 攻擊模式",
  "priority": 5,
  "input_params": {
    "data_sources": ["hibp", "abuseipdb"],
    "threat_type": "ddos",
    "time_window": "24h",
    "use_simulator": true
  },
  "tags": ["ddos", "network", "analysis"]
}
```

### 查詢任務狀態

```http
GET /api/v1/quantum-jobs/{job_id}
Authorization: Bearer <token>
```

### 列出任務

```http
GET /api/v1/quantum-jobs?status=completed&page=1&page_size=20
Authorization: Bearer <token>
```

## 開發指南

### 本地開發環境

```bash
# 克隆專案
git clone https://github.com/lipeichen/Ultimate-Security-Intelligence-Platform.git
cd Ultimate-Security-Intelligence-Platform/services/quantum-predictor

# 建立虛擬環境
python -m venv venv
source venv/bin/activate  # Linux/Mac
# 或
venv\Scripts\activate  # Windows

# 安裝依賴
pip install -r requirements.txt

# 運行測試
pytest tests/

# 運行服務
python main.py --help
```

### 建置 Docker 鏡像

```bash
# 建置鏡像
docker build -t quantum-predictor .

# 運行容器
docker run -it --rm quantum-predictor python main.py --help
```

### 推送到 GitHub Container Registry

```bash
# 登入 GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_USERNAME --password-stdin

# 標籤鏡像
docker tag quantum-predictor ghcr.io/lipeichen/Ultimate-Security-Intelligence-Platform/quantum-predictor:latest

# 推送鏡像
docker push ghcr.io/lipeichen/Ultimate-Security-Intelligence-Platform/quantum-predictor:latest
```

## 量子電路設計

服務使用變分量子分類器 (VQC) 進行威脅分析：

1. **特徵映射**: 將威脅數據映射到量子態
2. **變分電路**: 使用可訓練的量子電路進行分類
3. **測量**: 測量量子態獲得分類結果

### 電路結構

```
┌───┐     ┌─────────┐     ┌───┐
│ H │─────│ RZ(θ₁)  │─────│ H │
├───┤     ├─────────┤     ├───┤
│ H │─────│ RZ(θ₂)  │─────│ H │
├───┤     ├─────────┤     ├───┤
│ H │─────│ RZ(θ₃)  │─────│ H │
├───┤     ├─────────┤     ├───┤
│ H │─────│ RZ(θ₄)  │─────│ H │
└───┘     └─────────┘     └───┘
```

## 故障排除

### 常見問題

1. **IBM Quantum 連接失敗**
   - 檢查 `IBMQ_API_KEY` 是否正確設定
   - 確認網路連接正常

2. **Docker 鏡像拉取失敗**
   - 檢查網路連接
   - 確認鏡像標籤正確

3. **量子計算超時**
   - 增加 `shots` 參數
   - 使用模擬器而非真實設備

### 日誌查看

```bash
# 查看容器日誌
docker logs <container_id>

# 查看服務日誌
tail -f /var/log/quantum-predictor.log
```

## 貢獻指南

1. Fork 專案
2. 建立功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交變更 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 開啟 Pull Request

## 授權

本專案採用 MIT 授權條款 - 詳見 [LICENSE](../LICENSE) 文件。

## 支援

如有問題或建議，請開啟 [GitHub Issue](https://github.com/lipeichen/Ultimate-Security-Intelligence-Platform/issues)。
