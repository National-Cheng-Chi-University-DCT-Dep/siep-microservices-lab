#!/bin/bash

# 量子預測服務 Docker 測試腳本

set -e

echo "🧪 開始測試量子預測服務 Docker 鏡像..."

# 設定變數
IMAGE_NAME="ghcr.io/lipeichen/Ultimate-Security-Intelligence-Platform/quantum-predictor:latest"
TEST_INPUT="test_input.json"
TEST_OUTPUT="test_output.json"

# 創建測試輸入文件
cat > $TEST_INPUT << EOF
{
  "threats": [
    {
      "ip_address": "192.168.1.100",
      "threat_type": "malware",
      "risk_score": 85,
      "country": "US",
      "attack_type": "brute_force",
      "timestamp": "2024-01-15T10:30:00Z"
    },
    {
      "ip_address": "10.0.0.50",
      "threat_type": "ddos",
      "risk_score": 92,
      "country": "CN",
      "attack_type": "ddos",
      "timestamp": "2024-01-15T11:00:00Z"
    }
  ],
  "analysis_params": {
    "use_simulator": true,
    "shots": 1024,
    "confidence_threshold": 0.7
  }
}
EOF

echo "📝 創建測試輸入文件: $TEST_INPUT"

# 拉取 Docker 鏡像
echo "🐳 拉取 Docker 鏡像: $IMAGE_NAME"
docker pull $IMAGE_NAME

# 測試鏡像是否包含必要的文件
echo "🔍 檢查鏡像內容..."
docker run --rm $IMAGE_NAME ls -la /app/

# 測試幫助命令
echo "❓ 測試幫助命令..."
docker run --rm $IMAGE_NAME python main.py --help

# 運行實際測試
echo "🚀 運行量子預測測試..."
docker run --rm \
  -v "$(pwd)/$TEST_INPUT:/app/input.json" \
  -v "$(pwd)/$TEST_OUTPUT:/app/output.json" \
  $IMAGE_NAME \
  python main.py --input /app/input.json --output /app/output.json

# 檢查輸出
if [ -f "$TEST_OUTPUT" ]; then
    echo "✅ 測試成功！輸出文件已生成:"
    cat $TEST_OUTPUT
else
    echo "❌ 測試失敗！未找到輸出文件"
    exit 1
fi

# 清理測試文件
echo "🧹 清理測試文件..."
rm -f $TEST_INPUT $TEST_OUTPUT

echo "🎉 所有測試通過！Docker 鏡像工作正常。"
