#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
資安情報平台 - 量子模式預測服務

此服務使用量子電路模型對威脅數據進行預測分析，
可以作為獨立的微服務運行，或被後端 API 調用。

作者：資安情報平台開發團隊
版本：1.0.0
"""

import os
import json
import numpy as np
import logging
from pathlib import Path
from datetime import datetime
import argparse
import sys
from typing import Dict, List, Any, Optional, Tuple

# 量子計算相關庫
from qiskit import QuantumCircuit, Aer, transpile, assemble
from qiskit.visualization import plot_histogram
from qiskit_ibm_provider import IBMProvider
from qiskit.circuit import Parameter

# 設置日誌
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[logging.StreamHandler()]
)
logger = logging.getLogger("quantum-predictor")

# 全局配置
MODEL_DIR = os.environ.get("MODEL_DIR", os.path.join(os.path.dirname(__file__), "models"))
DEFAULT_MODEL = os.environ.get("DEFAULT_MODEL", "quantum_model_params.json")
USE_REAL_DEVICE = os.environ.get("USE_REAL_DEVICE", "false").lower() == "true"
IBM_PROVIDER = os.environ.get("IBMQ_PROVIDER", "ibm-q")
IBM_HUB = os.environ.get("IBMQ_HUB", "ibm-q-research")

class QuantumThreatClassifier:
    """量子威脅分類器類"""
    
    def __init__(self, model_path: str = None):
        """
        初始化量子威脅分類器
        
        參數:
            model_path: 模型參數文件路徑，如果為 None 則使用默認模型
        """
        if model_path is None:
            model_path = os.path.join(MODEL_DIR, DEFAULT_MODEL)
            
        # 加載模型參數
        self.load_model(model_path)
        
        # 初始化 IBM Quantum 提供者
        self.provider = None
        self.backend = None
        self._initialize_quantum_backend()
        
    def load_model(self, model_path: str) -> None:
        """
        從文件中加載模型參數
        
        參數:
            model_path: 模型參數文件路徑
        """
        try:
            if not os.path.exists(model_path):
                logger.warning(f"模型文件 {model_path} 不存在，使用默認參數")
                # 使用默認參數
                self.num_qubits = 4
                self.num_layers = 2
                self.parameters = np.random.rand(self.num_layers * 2 * self.num_qubits) * 2 * np.pi
                self.feature_mean = np.zeros(4)
                self.feature_scale = np.ones(4)
                self.threshold = 0.5
                return
                
            logger.info(f"從 {model_path} 加載模型參數")
            with open(model_path, 'r') as f:
                model_params = json.load(f)
                
            self.num_qubits = model_params.get('num_qubits', 4)
            self.num_layers = model_params.get('num_layers', 2)
            self.parameters = np.array(model_params.get('parameters', []))
            
            # 加載特徵縮放參數
            scaler_params = model_params.get('feature_scaler', {})
            self.feature_mean = np.array(scaler_params.get('mean', [0, 0, 0, 0]))
            self.feature_scale = np.array(scaler_params.get('scale', [1, 1, 1, 1]))
            
            self.threshold = model_params.get('threshold', 0.5)
            logger.info(f"模型加載成功: {self.num_qubits} 量子比特, {self.num_layers} 層")
            
        except Exception as e:
            logger.error(f"加載模型失敗: {str(e)}")
            raise
            
    def _initialize_quantum_backend(self) -> None:
        """初始化量子後端"""
        if USE_REAL_DEVICE:
            try:
                # 嘗試獲取 IBM Quantum API Token
                token = os.environ.get("IBMQ_API_KEY")
                if token:
                    IBMProvider.save_account(token, overwrite=True)
                
                # 加載提供者
                self.provider = IBMProvider()
                logger.info("成功連接到 IBM Quantum")
                
                # 獲取可用的後端
                backends = self.provider.backends()
                if backends:
                    # 選擇最少排隊的後端
                    least_busy = min(backends, key=lambda b: b.status().pending_jobs)
                    self.backend = least_busy
                    logger.info(f"選擇後端: {self.backend.name()}")
                else:
                    logger.warning("沒有可用的 IBM Quantum 後端，使用模擬器")
                    self.backend = Aer.get_backend('qasm_simulator')
            except Exception as e:
                logger.error(f"初始化 IBM Quantum 後端失敗: {str(e)}")
                logger.info("使用本地模擬器")
                self.backend = Aer.get_backend('qasm_simulator')
        else:
            logger.info("使用本地模擬器")
            self.backend = Aer.get_backend('qasm_simulator')
    
    def feature_map(self, features: np.ndarray) -> QuantumCircuit:
        """
        創建特徵映射電路
        
        參數:
            features: 輸入特徵
            
        返回:
            特徵映射量子電路
        """
        # 確保特徵數量與量子比特數量相符
        if len(features) > self.num_qubits:
            features = features[:self.num_qubits]
        
        qc = QuantumCircuit(self.num_qubits)
        
        # 首先將所有量子比特置於疊加態
        for q in range(self.num_qubits):
            qc.h(q)
        
        # 根據特徵值進行旋轉
        for i, feature in enumerate(features):
            if i < self.num_qubits:
                qc.rz(feature, i)
        
        # 加入糾纏操作
        for q in range(self.num_qubits-1):
            qc.cx(q, q+1)
        
        # 再次旋轉
        for i, feature in enumerate(features):
            if i < self.num_qubits:
                qc.ry(feature, i)
        
        return qc
        
    def variational_circuit(self, parameters: np.ndarray) -> QuantumCircuit:
        """
        創建變分電路
        
        參數:
            parameters: 變分參數
            
        返回:
            變分量子電路
        """
        qc = QuantumCircuit(self.num_qubits)
        
        # 將參數重塑為層數和每層參數
        num_layers = len(parameters) // (2 * self.num_qubits)
        params = parameters.reshape(num_layers, 2 * self.num_qubits)
        
        # 實現變分形式
        for layer in range(num_layers):
            # Rx 和 Rz 旋轉
            for q in range(self.num_qubits):
                qc.rx(params[layer][q], q)
                qc.rz(params[layer][self.num_qubits + q], q)
            
            # 糾纏操作
            for q in range(self.num_qubits-1):
                qc.cx(q, q+1)
            if self.num_qubits > 1:  # 添加循環連接
                qc.cx(self.num_qubits-1, 0)
        
        return qc
        
    def create_classifier_circuit(self, features: np.ndarray) -> QuantumCircuit:
        """
        創建完整的分類器電路
        
        參數:
            features: 輸入特徵
            
        返回:
            完整的量子分類器電路
        """
        # 特徵映射
        qc = self.feature_map(features)
        
        # 變分電路
        var_circuit = self.variational_circuit(self.parameters)
        qc = qc.compose(var_circuit)
        
        # 測量
        qc.measure_all()
        
        return qc
    
    def preprocess_features(self, features: np.ndarray) -> np.ndarray:
        """
        預處理特徵
        
        參數:
            features: 原始特徵
            
        返回:
            預處理後的特徵
        """
        if len(features) > len(self.feature_mean):
            features = features[:len(self.feature_mean)]
        elif len(features) < len(self.feature_mean):
            # 如果特徵維度不足，用 0 填充
            padding = np.zeros(len(self.feature_mean) - len(features))
            features = np.concatenate([features, padding])
        
        # 標準化特徵
        return (features - self.feature_mean) / self.feature_scale
        
    def predict(self, features: np.ndarray) -> Tuple[float, Dict[str, int]]:
        """
        使用量子電路對輸入特徵進行預測
        
        參數:
            features: 輸入特徵
            
        返回:
            預測概率和測量結果
        """
        # 預處理特徵
        preprocessed_features = self.preprocess_features(features)
        
        # 創建分類器電路
        qc = self.create_classifier_circuit(preprocessed_features)
        
        # 運行電路
        compiled_circuit = transpile(qc, self.backend)
        job = self.backend.run(compiled_circuit, shots=1024)
        result = job.result()
        counts = result.get_counts()
        
        # 計算預測概率
        prob_class_1 = 0
        for bitstring, count in counts.items():
            # 使用第一個量子比特作為輸出
            if bitstring[0] == '1':
                prob_class_1 += count / 1024
                
        return prob_class_1, counts
        
    def classify(self, features: np.ndarray) -> Dict[str, Any]:
        """
        對輸入特徵進行分類並返回結果
        
        參數:
            features: 輸入特徵
            
        返回:
            包含分類結果的字典
        """
        prob, counts = self.predict(features)
        
        # 應用閾值進行分類
        prediction = 1 if prob >= self.threshold else 0
        
        # 計算信心得分 (0-100)
        confidence = abs(prob - 0.5) * 2 * 100
        
        return {
            "prediction": int(prediction),
            "probability": float(prob),
            "confidence": float(confidence),
            "threshold": float(self.threshold),
            "counts": counts,
            "is_malicious": bool(prediction == 1),
            "timestamp": datetime.now().isoformat(),
            "backend": self.backend.name() if hasattr(self.backend, 'name') else "simulator"
        }

def process_file(file_path: str) -> Dict[str, Any]:
    """
    處理包含威脅數據的 JSON 文件
    
    參數:
        file_path: JSON 文件路徑
        
    返回:
        處理結果
    """
    try:
        with open(file_path, 'r') as f:
            data = json.load(f)
            
        return process_data(data)
    except Exception as e:
        logger.error(f"處理文件 {file_path} 失敗: {str(e)}")
        return {
            "error": str(e),
            "status": "failed",
            "timestamp": datetime.now().isoformat()
        }
        
def process_data(data: Dict[str, Any]) -> Dict[str, Any]:
    """
    處理威脅數據
    
    參數:
        data: 威脅數據字典
        
    返回:
        處理結果
    """
    try:
        # 初始化分類器
        classifier = QuantumThreatClassifier()
        
        # 提取特徵
        features = extract_features(data)
        
        # 進行分類
        result = classifier.classify(features)
        
        # 添加元數據
        result["input_data_summary"] = {
            "num_threats": len(data.get("threats", [])),
            "data_timestamp": data.get("timestamp", "unknown"),
            "features_used": features.tolist()
        }
        
        result["status"] = "success"
        return result
    except Exception as e:
        logger.error(f"處理數據失敗: {str(e)}")
        return {
            "error": str(e),
            "status": "failed",
            "timestamp": datetime.now().isoformat()
        }

def extract_features(data: Dict[str, Any]) -> np.ndarray:
    """
    從威脅數據中提取特徵
    
    參數:
        data: 威脅數據字典
        
    返回:
        特徵向量
    """
    threats = data.get("threats", [])
    
    if not threats:
        logger.warning("沒有找到威脅數據，使用零特徵向量")
        return np.zeros(4)
        
    # 提取特徵
    risk_scores = [t.get("risk_score", 0) for t in threats if "risk_score" in t]
    brute_force_count = sum(1 for t in threats if t.get("attack_type", "").lower() == "brute force")
    ddos_count = sum(1 for t in threats if t.get("attack_type", "").lower() == "ddos")
    unique_countries = len(set(t.get("country", "unknown") for t in threats if "country" in t))
    
    # 合成特徵向量
    avg_risk = np.mean(risk_scores) if risk_scores else 0
    max_risk = np.max(risk_scores) if risk_scores else 0
    
    features = np.array([
        avg_risk / 100.0 if avg_risk else 0,                # 平均風險分數 (標準化到 0-1)
        max_risk / 100.0 if max_risk else 0,                # 最高風險分數 (標準化到 0-1)
        min(brute_force_count + ddos_count, 100) / 100.0,   # 攻擊類型計數 (標準化到 0-1)
        min(unique_countries, 50) / 50.0                    # 不同國家數 (標準化到 0-1)
    ])
    
    return features

def save_result(result: Dict[str, Any], output_path: str) -> None:
    """
    將結果保存到文件
    
    參數:
        result: 處理結果
        output_path: 輸出文件路徑
    """
    try:
        with open(output_path, 'w') as f:
            json.dump(result, f, indent=2)
        logger.info(f"結果已保存到 {output_path}")
    except Exception as e:
        logger.error(f"保存結果失敗: {str(e)}")

def main():
    """主函數"""
    parser = argparse.ArgumentParser(description="資安情報平台量子模式預測服務")
    parser.add_argument("--input", "-i", type=str, help="輸入 JSON 文件路徑")
    parser.add_argument("--output", "-o", type=str, help="輸出 JSON 文件路徑")
    parser.add_argument("--model", "-m", type=str, help="模型參數文件路徑")
    parser.add_argument("--real-device", "-r", action="store_true", help="使用真實量子設備")
    
    args = parser.parse_args()
    
    # 檢查是否提供了輸入文件
    if not args.input:
        print("錯誤：必須提供輸入 JSON 文件路徑")
        parser.print_help()
        sys.exit(1)
        
    # 設置輸出文件路徑
    output_path = args.output if args.output else args.input.replace(".json", "_result.json")
    
    # 設置環境變數
    if args.real_device:
        os.environ["USE_REAL_DEVICE"] = "true"
        
    if args.model:
        os.environ["DEFAULT_MODEL"] = args.model
        
    # 處理文件
    result = process_file(args.input)
    
    # 保存結果
    save_result(result, output_path)
    
    # 輸出結果摘要
    if result.get("status") == "success":
        prediction = "惡意" if result.get("is_malicious", False) else "良性"
        confidence = result.get("confidence", 0)
        print(f"預測結果: {prediction} (信心度: {confidence:.2f}%)")
    else:
        print(f"處理失敗: {result.get('error', '未知錯誤')}")
        
    return 0

if __name__ == "__main__":
    sys.exit(main())
