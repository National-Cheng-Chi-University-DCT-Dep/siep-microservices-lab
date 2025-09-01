import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowPathIcon } from '@heroicons/react/24/outline';

interface FormData {
  title: string;
  description: string;
  priority: number;
  inputParams: {
    dataSources: string[];
    threatType: string;
    timeWindow: string;
    useSimulator: boolean;
  };
  tags: string[];
}

const QuantumJobForm: React.FC = () => {
  const router = useRouter();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const [formData, setFormData] = useState<FormData>({
    title: '',
    description: '',
    priority: 5,
    inputParams: {
      dataSources: ['hibp', 'abuseipdb'],
      threatType: 'malware',
      timeWindow: '24h',
      useSimulator: true
    },
    tags: []
  });
  
  const [tagInput, setTagInput] = useState('');
  
  const threatTypes = [
    { id: 'malware', name: '惡意軟體' },
    { id: 'ddos', name: 'DDoS 攻擊' },
    { id: 'brute_force', name: '暴力破解' },
    { id: 'phishing', name: '釣魚攻擊' },
    { id: 'ransomware', name: '勒索軟體' },
    { id: 'zero_day', name: '零日漏洞' }
  ];
  
  const timeWindows = [
    { id: '6h', name: '最近6小時' },
    { id: '12h', name: '最近12小時' },
    { id: '24h', name: '最近24小時' },
    { id: '48h', name: '最近48小時' },
    { id: '7d', name: '最近7天' },
    { id: '30d', name: '最近30天' }
  ];
  
  const dataSources = [
    { id: 'hibp', name: 'Have I Been Pwned', description: '資料外洩檢測' },
    { id: 'abuseipdb', name: 'AbuseIPDB', description: '惡意IP資料庫' },
    { id: 'internal', name: '內部威脅情報', description: '來自本平台的威脅資料' },
    { id: 'virus_total', name: 'VirusTotal', description: '病毒與惡意軟體檢測' },
    { id: 'mitre', name: 'MITRE ATT&CK', description: '攻擊技術知識庫' }
  ];
  
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    
    if (name.includes('.')) {
      const [parent, child] = name.split('.');
      setFormData(prev => ({
        ...prev,
        [parent]: {
          ...prev[parent as keyof FormData],
          [child]: value
        }
      }));
    } else {
      setFormData(prev => ({
        ...prev,
        [name]: value
      }));
    }
  };
  
  const handleDataSourceChange = (id: string) => {
    setFormData(prev => {
      const currentSources = prev.inputParams.dataSources;
      const newSources = currentSources.includes(id)
        ? currentSources.filter(source => source !== id)
        : [...currentSources, id];
        
      return {
        ...prev,
        inputParams: {
          ...prev.inputParams,
          dataSources: newSources
        }
      };
    });
  };
  
  const handleSimulatorToggle = () => {
    setFormData(prev => ({
      ...prev,
      inputParams: {
        ...prev.inputParams,
        useSimulator: !prev.inputParams.useSimulator
      }
    }));
  };
  
  const addTag = () => {
    if (tagInput.trim() && !formData.tags.includes(tagInput.trim())) {
      setFormData(prev => ({
        ...prev,
        tags: [...prev.tags, tagInput.trim()]
      }));
      setTagInput('');
    }
  };
  
  const removeTag = (tag: string) => {
    setFormData(prev => ({
      ...prev,
      tags: prev.tags.filter(t => t !== tag)
    }));
  };
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);
    
    try {
      // 將表單資料轉換為 API 所需的格式
      const apiData = {
        title: formData.title,
        description: formData.description,
        priority: parseInt(formData.priority.toString()),
        input_params: {
          data_sources: formData.inputParams.dataSources,
          threat_type: formData.inputParams.threatType,
          time_window: formData.inputParams.timeWindow,
          use_simulator: formData.inputParams.useSimulator
        },
        tags: formData.tags
      };
      
      // 呼叫 API 提交任務
      const response = await fetch('/api/v1/quantum-jobs', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(apiData)
      });
      
      if (!response.ok) {
        throw new Error('提交任務失敗');
      }
      
      const result = await response.json();
      
      // 提交成功後導航到結果頁面
      router.push(`/quantum/jobs/${result.job_id}`);
      
    } catch (err: any) {
      console.error('提交任務失敗:', err);
      setError(err.message || '提交任務失敗，請稍後再試');
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <div className="bg-white dark:bg-gray-900 shadow-lg rounded-lg p-6 max-w-4xl mx-auto">
      <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">提交量子分析任務</h2>
      
      {error && (
        <div className="bg-red-50 dark:bg-red-900/30 border-l-4 border-red-500 p-4 mb-6">
          <p className="text-red-700 dark:text-red-400">{error}</p>
        </div>
      )}
      
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* 基本資訊區塊 */}
        <div className="bg-gray-50 dark:bg-gray-800 p-4 rounded-md">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">基本資訊</h3>
          
          <div className="mb-4">
            <label htmlFor="title" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              任務標題 *
            </label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              required
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
              placeholder="輸入任務標題"
            />
          </div>
          
          <div className="mb-4">
            <label htmlFor="description" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              任務描述
            </label>
            <textarea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleChange}
              rows={3}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
              placeholder="描述此次分析任務的目的和期望"
            ></textarea>
          </div>
          
          <div className="mb-4">
            <label htmlFor="priority" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              優先級 (1-10)
            </label>
            <input
              type="range"
              id="priority"
              name="priority"
              min="1"
              max="10"
              value={formData.priority}
              onChange={handleChange}
              className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer"
            />
            <div className="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-1">
              <span>低 (1)</span>
              <span>中 (5)</span>
              <span>高 (10)</span>
            </div>
            <div className="text-center mt-1">
              <span className="text-sm font-medium">目前: {formData.priority}</span>
            </div>
          </div>
          
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              標籤
            </label>
            <div className="flex items-center">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                className="flex-grow px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-l-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
                placeholder="新增標籤"
                onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
              />
              <button
                type="button"
                onClick={addTag}
                className="px-4 py-2 bg-indigo-600 text-white rounded-r-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                新增
              </button>
            </div>
            {formData.tags.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-2">
                {formData.tags.map((tag) => (
                  <div
                    key={tag}
                    className="px-2 py-1 bg-indigo-100 dark:bg-indigo-900 text-indigo-800 dark:text-indigo-200 rounded-full text-sm flex items-center"
                  >
                    {tag}
                    <button
                      type="button"
                      onClick={() => removeTag(tag)}
                      className="ml-1 text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-200 focus:outline-none"
                    >
                      &times;
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
        
        {/* 分析參數區塊 */}
        <div className="bg-gray-50 dark:bg-gray-800 p-4 rounded-md">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">分析參數</h3>
          
          <div className="mb-4">
            <label htmlFor="threatType" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              威脅類型 *
            </label>
            <select
              id="threatType"
              name="inputParams.threatType"
              value={formData.inputParams.threatType}
              onChange={handleChange}
              required
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
            >
              {threatTypes.map((type) => (
                <option key={type.id} value={type.id}>
                  {type.name}
                </option>
              ))}
            </select>
          </div>
          
          <div className="mb-4">
            <label htmlFor="timeWindow" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              時間窗口 *
            </label>
            <select
              id="timeWindow"
              name="inputParams.timeWindow"
              value={formData.inputParams.timeWindow}
              onChange={handleChange}
              required
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
            >
              {timeWindows.map((window) => (
                <option key={window.id} value={window.id}>
                  {window.name}
                </option>
              ))}
            </select>
          </div>
          
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              資料來源 *
            </label>
            <div className="space-y-2">
              {dataSources.map((source) => (
                <div key={source.id} className="flex items-center">
                  <input
                    id={`source-${source.id}`}
                    type="checkbox"
                    checked={formData.inputParams.dataSources.includes(source.id)}
                    onChange={() => handleDataSourceChange(source.id)}
                    className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 dark:border-gray-600 rounded"
                  />
                  <label htmlFor={`source-${source.id}`} className="ml-2 text-sm text-gray-700 dark:text-gray-300">
                    <span className="font-medium">{source.name}</span>
                    <span className="ml-1 text-gray-500 dark:text-gray-400">({source.description})</span>
                  </label>
                </div>
              ))}
            </div>
            {formData.inputParams.dataSources.length === 0 && (
              <p className="text-sm text-red-600 dark:text-red-400 mt-1">請至少選擇一個資料來源</p>
            )}
          </div>
          
          <div className="mt-4">
            <div className="flex items-center">
              <input
                id="useSimulator"
                type="checkbox"
                checked={formData.inputParams.useSimulator}
                onChange={handleSimulatorToggle}
                className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 dark:border-gray-600 rounded"
              />
              <label htmlFor="useSimulator" className="ml-2 text-sm text-gray-700 dark:text-gray-300">
                使用模擬器（不使用實際量子設備）
              </label>
            </div>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1 ml-6">
              模擬器執行較快，但可能缺少真實量子設備的某些特性。真實量子設備可能需要較長的排隊等待時間。
            </p>
          </div>
        </div>
        
        {/* 提交按鈕 */}
        <div className="flex justify-end">
          <button
            type="submit"
            disabled={isSubmitting || formData.inputParams.dataSources.length === 0}
            className={`
              px-6 py-3 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 
              focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500
              flex items-center
              ${(isSubmitting || formData.inputParams.dataSources.length === 0) ? 'opacity-50 cursor-not-allowed' : ''}
            `}
          >
            {isSubmitting ? (
              <>
                <ArrowPathIcon className="h-5 w-5 mr-2 animate-spin" />
                提交中...
              </>
            ) : (
              '提交量子分析任務'
            )}
          </button>
        </div>
      </form>
    </div>
  );
};

export default QuantumJobForm;
