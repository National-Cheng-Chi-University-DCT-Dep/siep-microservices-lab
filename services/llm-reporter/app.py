#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
資安情報平台 - LLM 報告生成服務
此服務使用 Gradio 提供 API 端點，用於生成資安情報分析報告。
"""

import os
import json
import gradio as gr
import numpy as np
import pandas as pd
from datetime import datetime
from dotenv import load_dotenv
from typing import Dict, List, Any, Optional
import logging
from pathlib import Path

# 設置日誌
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[logging.StreamHandler()]
)
logger = logging.getLogger("llm-reporter")

# 載入環境變數
load_dotenv()

# 配置
MODEL_NAME = os.environ.get("MODEL_NAME", "mistralai/Mistral-7B-Instruct-v0.2")
USE_MOCK = os.environ.get("USE_MOCK", "true").lower() == "true"
API_KEY = os.environ.get("HF_API_TOKEN", None)
MAX_NEW_TOKENS = int(os.environ.get("MAX_NEW_TOKENS", "1024"))

# 模型與分詞器
model = None
tokenizer = None

def load_model():
    """
    載入 LLM 模型與分詞器
    如果 USE_MOCK 為 True，則不載入實際模型，以節省資源
    """
    global model, tokenizer
    
    if USE_MOCK:
        logger.info("使用模擬模式，不載入實際模型")
        return

    try:
        from transformers import AutoModelForCausalLM, AutoTokenizer
        
        logger.info(f"正在載入模型: {MODEL_NAME}")
        tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME, use_auth_token=API_KEY)
        model = AutoModelForCausalLM.from_pretrained(MODEL_NAME, use_auth_token=API_KEY)
        logger.info("模型載入完成")
        
    except Exception as e:
        logger.error(f"載入模型時發生錯誤: {str(e)}")
        raise

# 模擬生成報告的函數（開發測試用）
def generate_mock_report(threat_data, report_type, language):
    """生成模擬報告，用於開發和測試"""
    logger.info(f"使用模擬數據生成 {report_type} 報告，語言: {language}")
    
    # 報告模板
    templates = {
        "zh-Hant": {
            "summary": """# 資安威脅情報摘要報告

## 概述
在分析的資料中，我們發現共有 {threat_count} 個潛在威脅，其中包括 {high_count} 個高風險威脅。主要威脅類型為 {top_threat}。

## 高風險威脅
以下是識別出的前 3 個高風險威脅：
1. 惡意 IP {ip_1}：來自 {country_1}，涉及多次 {attack_type_1} 攻擊
2. 惡意 IP {ip_2}：來自 {country_2}，涉及 {attack_type_2}
3. 惡意 IP {ip_3}：來自 {country_3}，涉及 {attack_type_3}

## 建議
1. 立即將這些 IP 加入防火牆黑名單
2. 檢查系統是否有與這些 IP 的通訊記錄
3. 更新所有系統的安全補丁
""",
            "detailed": """# 資安威脅情報詳細分析報告

## 執行摘要
本報告基於 {date} 收集的資料，分析了 {threat_count} 個潛在威脅。整體安全風險等級評估為：**{risk_level}**。

## 威脅分類統計
| 類型 | 數量 | 百分比 |
|------|------|--------|
| 暴力破解 | {brute_force} | {brute_force_pct}% |
| DDoS | {ddos} | {ddos_pct}% |
| 資料外洩 | {data_leak} | {data_leak_pct}% |
| 惡意軟體 | {malware} | {malware_pct}% |
| 其他 | {others} | {others_pct}% |

## 地理分佈
攻擊來源前五名國家：
1. {country_1} ({country_1_pct}%)
2. {country_2} ({country_2_pct}%)
3. {country_3} ({country_3_pct}%)
4. {country_4} ({country_4_pct}%)
5. {country_5} ({country_5_pct}%)

## 詳細分析
在分析期間，我們觀察到多次來自 {top_country} 的協同攻擊，主要針對 {target_service} 服務。攻擊模式顯示這可能是一個有組織的攻擊者團體，採用了 {technique} 技術。

## 威脅時間趨勢
過去 24 小時內，攻擊高峰期出現在 {peak_time}，共記錄 {peak_count} 次攻擊嘗試。

## 建議行動
1. 立即部署 {defense_1} 防禦措施
2. 更新 {service} 至最新版本 {version}
3. 強化 {weak_point} 的存取控制
4. 設置額外監控於 {monitor_target}
5. 定期審查 {review_item} 日誌

## 附錄
完整的威脅 IP 列表已附加於報告末尾。請參考附件以獲取更多技術詳情。
"""
        },
        "en-US": {
            "summary": """# Security Threat Intelligence Summary Report

## Overview
In the analyzed data, we identified a total of {threat_count} potential threats, including {high_count} high-risk threats. The predominant threat type is {top_threat}.

## High Risk Threats
Here are the top 3 high-risk threats identified:
1. Malicious IP {ip_1}: From {country_1}, involved in multiple {attack_type_1} attacks
2. Malicious IP {ip_2}: From {country_2}, involved in {attack_type_2}
3. Malicious IP {ip_3}: From {country_3}, involved in {attack_type_3}

## Recommendations
1. Immediately add these IPs to your firewall blacklist
2. Check your systems for any communication with these IPs
3. Update all system security patches
"""
        }
    }
    
    # 選擇模板
    lang_key = "zh-Hant" if language == "繁體中文" else "en-US"
    template = templates.get(lang_key, {}).get(report_type, templates["zh-Hant"]["summary"])
    
    # 生成模擬數據
    mock_data = {
        "threat_count": np.random.randint(50, 500),
        "high_count": np.random.randint(5, 50),
        "top_threat": np.random.choice(["SQL注入", "跨站腳本攻擊", "DDoS", "憑證填充"]),
        "ip_1": f"192.168.{np.random.randint(1, 255)}.{np.random.randint(1, 255)}",
        "ip_2": f"10.0.{np.random.randint(1, 255)}.{np.random.randint(1, 255)}",
        "ip_3": f"172.16.{np.random.randint(1, 255)}.{np.random.randint(1, 255)}",
        "country_1": np.random.choice(["中國", "俄羅斯", "北韓"]),
        "country_2": np.random.choice(["美國", "巴西", "印度"]),
        "country_3": np.random.choice(["烏克蘭", "羅馬尼亞", "荷蘭"]),
        "attack_type_1": np.random.choice(["暴力破解", "DDoS", "SQL注入"]),
        "attack_type_2": np.random.choice(["勒索軟體", "釣魚", "後門"]),
        "attack_type_3": np.random.choice(["憑證填充", "中間人攻擊", "零日漏洞"]),
        "date": datetime.now().strftime("%Y-%m-%d"),
        "risk_level": np.random.choice(["低", "中", "高", "嚴重"]),
        "brute_force": np.random.randint(10, 100),
        "ddos": np.random.randint(5, 50),
        "data_leak": np.random.randint(1, 30),
        "malware": np.random.randint(10, 80),
        "country_1_pct": np.random.randint(20, 45),
        "country_2_pct": np.random.randint(15, 30),
        "country_3_pct": np.random.randint(10, 20),
        "country_4_pct": np.random.randint(5, 15),
        "country_5_pct": np.random.randint(1, 10),
        "top_country": np.random.choice(["中國", "俄羅斯", "北韓"]),
        "target_service": np.random.choice(["Web伺服器", "郵件系統", "資料庫"]),
        "technique": np.random.choice(["APT", "社交工程", "供應鏈攻擊"]),
        "peak_time": f"{np.random.randint(0, 24):02d}:00-{np.random.randint(0, 24):02d}:59",
        "peak_count": np.random.randint(50, 500),
        "defense_1": np.random.choice(["WAF", "入侵檢測系統", "進階終端防護"]),
        "service": np.random.choice(["Apache", "Nginx", "MySQL"]),
        "version": f"{np.random.randint(1, 10)}.{np.random.randint(0, 20)}.{np.random.randint(0, 99)}",
        "weak_point": np.random.choice(["管理介面", "API端點", "資料庫連線"]),
        "monitor_target": np.random.choice(["流量模式", "身份驗證嘗試", "敏感資料存取"]),
        "review_item": np.random.choice(["系統", "防火牆", "應用程式"])
    }
    
    # 計算其他值
    total = mock_data["brute_force"] + mock_data["ddos"] + mock_data["data_leak"] + mock_data["malware"]
    mock_data["others"] = np.random.randint(1, 20)
    new_total = total + mock_data["others"]
    
    mock_data["brute_force_pct"] = round(mock_data["brute_force"] / new_total * 100, 1)
    mock_data["ddos_pct"] = round(mock_data["ddos"] / new_total * 100, 1)
    mock_data["data_leak_pct"] = round(mock_data["data_leak"] / new_total * 100, 1)
    mock_data["malware_pct"] = round(mock_data["malware"] / new_total * 100, 1)
    mock_data["others_pct"] = round(mock_data["others"] / new_total * 100, 1)
    
    # 以模板格式化報告
    report = template.format(**mock_data)
    return report

# 使用 LLM 生成報告
def generate_report_with_llm(threat_data, report_type, language):
    """使用 LLM 生成報告"""
    if model is None or tokenizer is None:
        logger.warning("模型未載入，使用模擬報告替代")
        return generate_mock_report(threat_data, report_type, language)
    
    # 根據語言準備提示
    language_prompt = "以繁體中文撰寫" if language == "繁體中文" else "write in English"
    report_type_prompt = "摘要報告" if report_type == "summary" else "詳細分析報告"
    
    # 構建 prompt
    prompt = f"""你是一位專業的資安分析師，請根據以下資安威脅數據{language_prompt}生成一份資安情報{report_type_prompt}：

```json
{json.dumps(threat_data, ensure_ascii=False, indent=2)}
```

報告應該包含：
1. 威脅總覽與風險評估
2. 主要威脅類型分析
3. 地理來源分布
4. 具體的應對建議

請以 Markdown 格式輸出，保持專業客觀的語氣。
"""

    # 使用模型生成報告
    inputs = tokenizer(prompt, return_tensors="pt")
    outputs = model.generate(
        inputs["input_ids"],
        max_new_tokens=MAX_NEW_TOKENS,
        do_sample=True,
        temperature=0.7,
        top_p=0.9,
    )
    report = tokenizer.decode(outputs[0], skip_special_tokens=True)
    
    # 只返回生成的內容部分（去除原始提示）
    report = report[len(prompt):].strip()
    return report

# 處理上傳的 JSON 文件
def process_threat_file(file_obj):
    """處理上傳的威脅數據文件"""
    if file_obj is None:
        return None
    
    try:
        content = file_obj.read()
        if isinstance(content, bytes):
            content = content.decode('utf-8')
        return json.loads(content)
    except Exception as e:
        logger.error(f"處理檔案時發生錯誤: {str(e)}")
        return None

# 主要報告生成函數
def generate_security_report(
    threat_file: Optional[str] = None,
    threat_json: str = "",
    report_type: str = "summary",
    language: str = "繁體中文"
):
    """
    生成資安情報報告
    
    參數:
        threat_file: 上傳的威脅數據文件
        threat_json: JSON 格式的威脅數據文本
        report_type: 報告類型 (summary 或 detailed)
        language: 報告語言 (繁體中文 或 English)
    
    返回:
        生成的報告內容
    """
    logger.info(f"開始生成 {report_type} 報告，語言: {language}")
    
    # 獲取威脅數據
    threat_data = None
    
    # 從文件獲取
    if threat_file is not None:
        threat_data = process_threat_file(threat_file)
    
    # 如果沒有文件或文件處理失敗，嘗試從JSON文本獲取
    if threat_data is None and threat_json:
        try:
            threat_data = json.loads(threat_json)
        except Exception as e:
            logger.error(f"解析 JSON 時發生錯誤: {str(e)}")
            return "錯誤：無法解析威脅數據，請檢查 JSON 格式是否正確"
    
    # 如果仍然沒有數據，使用示例數據
    if threat_data is None:
        threat_data = {
            "sample_data": True,
            "threats": [
                {
                    "ip": "192.168.1.1",
                    "country": "Unknown",
                    "risk_score": 85,
                    "attack_type": "Brute Force",
                    "timestamp": "2025-11-01T12:34:56Z"
                }
            ]
        }
        logger.info("使用示例數據生成報告")
    
    # 生成報告
    try:
        if USE_MOCK:
            report = generate_mock_report(threat_data, report_type, language)
        else:
            report = generate_report_with_llm(threat_data, report_type, language)
        
        logger.info("報告生成成功")
        return report
    
    except Exception as e:
        logger.error(f"生成報告時發生錯誤: {str(e)}")
        return f"錯誤：生成報告時發生問題 - {str(e)}"

# 建立 Gradio 界面
def create_interface():
    """創建 Gradio Web 界面"""
    with gr.Blocks(title="資安情報報告生成器") as app:
        gr.Markdown("# 資安情報報告生成器")
        gr.Markdown("上傳資安威脅數據，自動生成專業的分析報告。")
        
        with gr.Row():
            with gr.Column(scale=1):
                threat_file = gr.File(label="上傳威脅數據文件（JSON格式）")
                threat_json = gr.Textbox(label="或直接輸入 JSON 格式的威脅數據", lines=10, placeholder="在此輸入 JSON 格式的威脅數據...")
                
                report_type = gr.Radio(
                    choices=["summary", "detailed"],
                    value="summary",
                    label="報告類型",
                    info="摘要報告適合快速了解情況，詳細報告包含完整分析"
                )
                
                language = gr.Radio(
                    choices=["繁體中文", "English"],
                    value="繁體中文",
                    label="報告語言"
                )
                
                generate_btn = gr.Button("生成報告", variant="primary")
            
            with gr.Column(scale=1):
                output = gr.Markdown(label="生成的報告")
        
        generate_btn.click(
            fn=generate_security_report,
            inputs=[threat_file, threat_json, report_type, language],
            outputs=output
        )
    
    return app

# API 端點定義
def create_api():
    """創建 Gradio API 端點"""
    return gr.Interface(
        fn=generate_security_report,
        inputs=[
            gr.File(label="威脅數據文件"),
            gr.Textbox(label="威脅數據JSON"),
            gr.Radio(
                choices=["summary", "detailed"],
                value="summary",
                label="報告類型"
            ),
            gr.Radio(
                choices=["繁體中文", "English"],
                value="繁體中文",
                label="報告語言"
            )
        ],
        outputs="text",
        title="資安情報報告生成 API",
        description="提供威脅數據，生成專業的資安分析報告",
        examples=[
            [
                None, 
                '{"threats":[{"ip":"192.168.1.1","risk_score":85,"attack_type":"Brute Force"}]}',
                "summary",
                "繁體中文"
            ]
        ]
    )

# 主程式
if __name__ == "__main__":
    # 只在非模擬模式下載入模型
    if not USE_MOCK:
        load_model()
    
    # 設定端口（如果在環境變數中提供）
    port = int(os.environ.get("PORT", 7860))
    
    # 創建 Web 界面和 API
    web_app = create_interface()
    api = create_api()
    
    # 啟動應用
    gr.mount_gradio_app(web_app, api, path="/api")
    web_app.launch(server_name="0.0.0.0", server_port=port)
