"""
量子預測器模組
提供威脅情報的量子分析功能
"""
import json
import time
import logging
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
import numpy as np

try:
    from qiskit import QuantumCircuit, execute, Aer, IBMQ
    from qiskit.providers.ibmq import least_busy
    from qiskit.visualization import plot_histogram
    QISKIT_AVAILABLE = True
except ImportError:
    QISKIT_AVAILABLE = False
    logging.warning("Qiskit not available, using mock quantum operations")

# 設定日誌
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@dataclass
class ThreatData:
    """威脅數據結構"""
    ip_address: str
    threat_type: str
    risk_score: int
    country: str
    attack_type: str
    timestamp: str


@dataclass
class AnalysisParams:
    """分析參數結構"""
    use_simulator: bool = True
    shots: int = 1024
    confidence_threshold: float = 0.7


class QuantumPredictor:
    """量子預測器類"""
    
    def __init__(self):
        """初始化量子預測器"""
        self.backend = None
        self.circuit = None
        self.provider = None
        self._setup_backend()
    
    def _setup_backend(self):
        """設定量子後端"""
        if not QISKIT_AVAILABLE:
            logger.warning("Using mock quantum backend")
            self.backend = MockBackend()
            return
        
        try:
            # 使用本地模擬器作為預設
            self.backend = Aer.get_backend('qasm_simulator')
            logger.info("Using local simulator backend")
        except Exception as e:
            logger.error(f"Failed to setup backend: {e}")
            self.backend = MockBackend()
    
    def load_account(self, api_key: str) -> bool:
        """載入 IBM Quantum 帳戶"""
        if not QISKIT_AVAILABLE:
            logger.warning("Qiskit not available, skipping account loading")
            return False
        
        try:
            IBMQ.load_account(api_key)
            self.provider = IBMQ.get_provider(hub='ibm-q')
            logger.info("Successfully loaded IBM Quantum account")
            return True
        except Exception as e:
            logger.error(f"Failed to load IBM Quantum account: {e}")
            return False
    
    def preprocess_data(self, input_data: Dict[str, Any]) -> List[float]:
        """預處理輸入數據"""
        threats = input_data.get('threats', [])
        if not threats:
            return [0.0] * 4
        
        # 提取特徵
        features = []
        for threat in threats:
            # 風險分數標準化 (0-100 -> 0-1)
            risk_score = threat.get('risk_score', 0) / 100.0
            
            # 威脅類型編碼
            threat_type = threat.get('threat_type', 'unknown')
            threat_encoding = self._encode_threat_type(threat_type)
            
            # 國家編碼
            country = threat.get('country', 'unknown')
            country_encoding = self._encode_country(country)
            
            # 攻擊類型編碼
            attack_type = threat.get('attack_type', 'unknown')
            attack_encoding = self._encode_attack_type(attack_type)
            
            features.extend([risk_score, threat_encoding, country_encoding, attack_encoding])
        
        # 如果特徵不足，用零填充
        while len(features) < 4:
            features.append(0.0)
        
        # 只取前4個特徵
        return features[:4]
    
    def _encode_threat_type(self, threat_type: str) -> float:
        """編碼威脅類型"""
        encoding_map = {
            'malware': 0.2,
            'ddos': 0.4,
            'phishing': 0.6,
            'brute_force': 0.8,
            'sql_injection': 1.0
        }
        return encoding_map.get(threat_type.lower(), 0.0)
    
    def _encode_country(self, country: str) -> float:
        """編碼國家"""
        # 簡化的國家風險編碼
        high_risk_countries = ['CN', 'RU', 'KP', 'IR']
        medium_risk_countries = ['BR', 'IN', 'PK', 'NG']
        
        if country.upper() in high_risk_countries:
            return 1.0
        elif country.upper() in medium_risk_countries:
            return 0.5
        else:
            return 0.0
    
    def _encode_attack_type(self, attack_type: str) -> float:
        """編碼攻擊類型"""
        encoding_map = {
            'ddos': 0.2,
            'brute_force': 0.4,
            'sql_injection': 0.6,
            'xss': 0.8,
            'rce': 1.0
        }
        return encoding_map.get(attack_type.lower(), 0.0)
    
    def create_circuit(self, num_qubits: int) -> QuantumCircuit:
        """創建量子電路"""
        if not QISKIT_AVAILABLE:
            return MockCircuit(num_qubits)
        
        circuit = QuantumCircuit(num_qubits, 1)
        
        # 簡單的變分量子分類器
        for i in range(num_qubits):
            circuit.h(i)  # Hadamard gate
        
        # 添加一些旋轉門來創建可訓練的參數
        for i in range(num_qubits):
            circuit.rz(0.5, i)  # 旋轉Z門
        
        # 再次應用 Hadamard 門
        for i in range(num_qubits):
            circuit.h(i)
        
        # 測量
        circuit.measure_all()
        
        return circuit
    
    def run_quantum_analysis(self, input_data: Dict[str, Any]) -> Dict[str, Any]:
        """運行量子分析"""
        start_time = time.time()
        
        # 預處理數據
        features = self.preprocess_data(input_data)
        
        # 創建電路
        self.circuit = self.create_circuit(len(features))
        
        # 執行量子電路
        if QISKIT_AVAILABLE and self.backend and not isinstance(self.backend, MockBackend):
            job = execute(self.circuit, self.backend, shots=1024)
            result = job.result()
            counts = result.get_counts(self.circuit)
        else:
            # 使用模擬結果
            counts = self._simulate_counts(features)
        
        # 解釋結果
        interpretation = self.interpret_results(counts)
        
        # 計算執行時間
        execution_time = int(time.time() - start_time)
        
        return {
            **interpretation,
            'counts': counts,
            'backend': str(self.backend),
            'execution_time': execution_time,
            'timestamp': time.strftime('%Y-%m-%dT%H:%M:%SZ')
        }
    
    def _simulate_counts(self, features: List[float]) -> Dict[str, int]:
        """模擬量子測量結果"""
        # 基於特徵計算惡意概率
        avg_risk = sum(features) / len(features)
        malicious_prob = min(avg_risk, 1.0)
        
        # 模擬測量結果
        total_shots = 1024
        malicious_shots = int(total_shots * malicious_prob)
        benign_shots = total_shots - malicious_shots
        
        return {
            '0000': benign_shots,
            '0001': malicious_shots
        }
    
    def interpret_results(self, counts: Dict[str, int]) -> Dict[str, Any]:
        """解釋量子測量結果"""
        total_shots = sum(counts.values())
        if total_shots == 0:
            return {
                'prediction': 0,
                'probability': 0.0,
                'confidence': 0.0,
                'is_malicious': False
            }
        
        # 計算惡意概率
        malicious_states = ['0001', '0010', '0011', '0100', '0101', '0110', '0111',
                           '1000', '1001', '1010', '1011', '1100', '1101', '1110', '1111']
        
        malicious_count = sum(counts.get(state, 0) for state in malicious_states)
        probability = malicious_count / total_shots
        
        # 計算置信度
        confidence = self.calculate_confidence(counts)
        
        # 做出預測
        prediction = 1 if probability > 0.5 else 0
        is_malicious = probability > 0.5
        
        return {
            'prediction': prediction,
            'probability': probability,
            'confidence': confidence,
            'is_malicious': is_malicious
        }
    
    def calculate_confidence(self, counts: Dict[str, int]) -> float:
        """計算置信度"""
        total_shots = sum(counts.values())
        if total_shots == 0:
            return 0.0
        
        # 計算測量結果的標準差作為置信度指標
        values = list(counts.values())
        mean = np.mean(values)
        std = np.std(values)
        
        # 標準化置信度到 0-100
        confidence = min(100.0, max(0.0, (1 - std / mean) * 100))
        return confidence
    
    def analyze_threats(self, input_data: Dict[str, Any]) -> Dict[str, Any]:
        """分析威脅數據的主函數"""
        try:
            # 驗證輸入
            if not self.validate_input(input_data):
                raise ValueError("Invalid input data")
            
            # 運行量子分析
            result = self.run_quantum_analysis(input_data)
            
            return result
            
        except Exception as e:
            logger.error(f"Analysis failed: {e}")
            return {
                'prediction': 0,
                'probability': 0.0,
                'confidence': 0.0,
                'is_malicious': False,
                'error': str(e),
                'timestamp': time.strftime('%Y-%m-%dT%H:%M:%SZ')
            }
    
    def validate_input(self, input_data: Dict[str, Any]) -> bool:
        """驗證輸入數據"""
        if not isinstance(input_data, dict):
            return False
        
        if 'threats' not in input_data:
            return False
        
        threats = input_data.get('threats', [])
        if not isinstance(threats, list):
            return False
        
        # 檢查每個威脅數據
        for threat in threats:
            if not isinstance(threat, dict):
                return False
            
            required_fields = ['ip_address', 'threat_type', 'risk_score']
            for field in required_fields:
                if field not in threat:
                    return False
        
        return True
    
    def save_results(self, results: Dict[str, Any], filename: str):
        """保存結果到文件"""
        try:
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(results, f, indent=2, ensure_ascii=False)
            logger.info(f"Results saved to {filename}")
        except Exception as e:
            logger.error(f"Failed to save results: {e}")
    
    def load_results(self, filename: str) -> Dict[str, Any]:
        """從文件載入結果"""
        try:
            with open(filename, 'r', encoding='utf-8') as f:
                return json.load(f)
        except Exception as e:
            logger.error(f"Failed to load results: {e}")
            return {}


class MockBackend:
    """模擬量子後端"""
    def __str__(self):
        return "mock_backend"


class MockCircuit:
    """模擬量子電路"""
    def __init__(self, num_qubits: int):
        self.num_qubits = num_qubits
    
    def measure_all(self):
        """模擬測量操作"""
        pass


# 主函數用於測試
if __name__ == "__main__":
    # 測試量子預測器
    predictor = QuantumPredictor()
    
    test_data = {
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
            "use_simulator": True,
            "shots": 1024,
            "confidence_threshold": 0.7
        }
    }
    
    result = predictor.analyze_threats(test_data)
    print("分析結果:", json.dumps(result, indent=2, ensure_ascii=False))
