"""
量子預測器簡化測試
不依賴於 Qiskit 的具體實現
"""
import pytest
import json
import tempfile
import os
from quantum_predictor import QuantumPredictor


class TestQuantumPredictorSimple:
    """量子預測器簡化測試類"""

    def setup_method(self):
        """每個測試方法前的設置"""
        self.predictor = QuantumPredictor()
        self.sample_input = {
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

    def test_initialization(self):
        """測試預測器初始化"""
        assert self.predictor is not None
        assert hasattr(self.predictor, 'backend')
        assert hasattr(self.predictor, 'circuit')

    def test_preprocess_data(self):
        """測試數據預處理"""
        processed_data = self.predictor.preprocess_data(self.sample_input)
        
        assert isinstance(processed_data, list)
        assert len(processed_data) > 0
        assert all(isinstance(item, (int, float)) for item in processed_data)

    def test_create_circuit(self):
        """測試創建量子電路"""
        circuit = self.predictor.create_circuit(4)
        assert circuit is not None

    def test_run_quantum_analysis(self):
        """測試運行量子分析"""
        result = self.predictor.run_quantum_analysis(self.sample_input)
        
        assert isinstance(result, dict)
        assert 'prediction' in result
        assert 'probability' in result
        assert 'confidence' in result
        assert 'is_malicious' in result

    def test_save_results(self):
        """測試保存結果"""
        test_results = {
            "prediction": 1,
            "probability": 0.85,
            "confidence": 75.5,
            "is_malicious": True
        }
        
        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            temp_file = f.name
        
        try:
            self.predictor.save_results(test_results, temp_file)
            
            with open(temp_file, 'r') as f:
                saved_data = json.load(f)
            
            assert saved_data == test_results
        finally:
            os.unlink(temp_file)

    def test_load_results(self):
        """測試載入結果"""
        test_results = {
            "prediction": 1,
            "probability": 0.85,
            "confidence": 75.5,
            "is_malicious": True
        }
        
        with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
            json.dump(test_results, f)
            temp_file = f.name
        
        try:
            loaded_results = self.predictor.load_results(temp_file)
            assert loaded_results == test_results
        finally:
            os.unlink(temp_file)

    def test_validate_input(self):
        """測試輸入驗證"""
        # 測試有效輸入
        assert self.predictor.validate_input(self.sample_input) is True
        
        # 測試無效輸入
        invalid_input = {"invalid": "data"}
        assert self.predictor.validate_input(invalid_input) is False

    def test_calculate_confidence(self):
        """測試置信度計算"""
        counts = {'0000': 500, '0001': 524}
        confidence = self.predictor.calculate_confidence(counts)
        
        assert isinstance(confidence, float)
        assert 0 <= confidence <= 100

    def test_interpret_results(self):
        """測試結果解釋"""
        counts = {'0000': 500, '0001': 524}
        interpretation = self.predictor.interpret_results(counts)
        
        assert isinstance(interpretation, dict)
        assert 'prediction' in interpretation
        assert 'is_malicious' in interpretation

    def test_analyze_threats(self):
        """測試威脅分析主函數"""
        result = self.predictor.analyze_threats(self.sample_input)
        
        assert isinstance(result, dict)
        assert 'prediction' in result
        assert 'probability' in result
        assert 'confidence' in result
        assert 'is_malicious' in result
        assert 'backend' in result
        assert 'execution_time' in result

    def test_encode_functions(self):
        """測試編碼函數"""
        # 測試威脅類型編碼
        threat_encoding = self.predictor._encode_threat_type('malware')
        assert isinstance(threat_encoding, float)
        assert 0 <= threat_encoding <= 1
        
        # 測試國家編碼
        country_encoding = self.predictor._encode_country('US')
        assert isinstance(country_encoding, float)
        assert 0 <= country_encoding <= 1
        
        # 測試攻擊類型編碼
        attack_encoding = self.predictor._encode_attack_type('ddos')
        assert isinstance(attack_encoding, float)
        assert 0 <= attack_encoding <= 1

    def test_simulate_counts(self):
        """測試模擬計數"""
        features = [0.5, 0.3, 0.7, 0.2]
        counts = self.predictor._simulate_counts(features)
        
        assert isinstance(counts, dict)
        assert '0000' in counts
        assert '0001' in counts
        assert sum(counts.values()) == 1024


class TestQuantumPredictorIntegration:
    """量子預測器整合測試類"""

    @pytest.fixture
    def predictor(self):
        """創建預測器實例"""
        return QuantumPredictor()

    def test_end_to_end_analysis(self, predictor):
        """端到端分析測試"""
        input_data = {
            "threats": [
                {
                    "ip_address": "10.0.0.1",
                    "threat_type": "ddos",
                    "risk_score": 90,
                    "country": "CN",
                    "attack_type": "ddos",
                    "timestamp": "2024-01-15T11:00:00Z"
                }
            ],
            "analysis_params": {
                "use_simulator": True,
                "shots": 100,
                "confidence_threshold": 0.6
            }
        }
        
        result = predictor.analyze_threats(input_data)
        
        assert isinstance(result, dict)
        assert 'prediction' in result
        assert 'probability' in result
        assert 'confidence' in result
        assert 'is_malicious' in result
        assert 'backend' in result
        assert 'execution_time' in result

    def test_error_handling(self, predictor):
        """測試錯誤處理"""
        # 測試無效輸入
        invalid_input = {"invalid": "data"}
        result = predictor.analyze_threats(invalid_input)
        
        assert isinstance(result, dict)
        assert 'error' in result
        assert result['prediction'] == 0
        assert result['is_malicious'] is False


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
