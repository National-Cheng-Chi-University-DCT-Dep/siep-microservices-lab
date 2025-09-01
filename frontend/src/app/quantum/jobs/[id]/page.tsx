'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowPathIcon, CheckCircleIcon, ExclamationCircleIcon, ClockIcon, PlayIcon } from '@heroicons/react/24/outline';
import JobStatusTracker from '@/components/JobStatusTracker';
import ProbabilityChart from '@/components/ProbabilityChart';

interface JobResult {
  prediction: number;
  probability: number;
  confidence: number;
  is_malicious: boolean;
  counts: Record<string, number>;
}

interface QuantumJob {
  id: string;
  title: string;
  description: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  created_at: string;
  updated_at: string;
  started_at?: string;
  completed_at?: string;
  execution_time_seconds?: number;
  quantum_backend?: string;
  is_simulation: boolean;
  confidence_score?: number;
  is_malicious?: boolean;
  results?: JobResult;
  error_message?: string;
  input_params_summary?: Record<string, unknown>;
  results_summary?: Record<string, unknown>;
}

export default function QuantumJobDetail({ params }: { params: { id: string } }) {
  const router = useRouter();
  const [job, setJob] = useState<QuantumJob | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  
  const jobId = params.id;
  
  // 輪詢任務狀態
  useEffect(() => {
    const fetchJobDetail = async () => {
      try {
        const response = await fetch(`/api/v1/quantum-jobs/${jobId}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        
        if (!response.ok) {
          throw new Error('無法獲取任務詳情');
        }
        
        const data = await response.json();
        setJob(data);
        setError(null);
      } catch (err: unknown) {
        const errorMessage = err instanceof Error ? err.message : '獲取任務詳情失敗';
        setError(errorMessage);
        console.error('獲取任務失敗:', err);
      } finally {
        setLoading(false);
      }
    };
    
    fetchJobDetail();
    
    // 如果任務尚未完成，設定輪詢
    const intervalId = setInterval(() => {
      if (job && (job.status === 'pending' || job.status === 'running')) {
        fetchJobDetail();
      }
    }, 5000); // 每5秒更新一次
    
    return () => clearInterval(intervalId);
  }, [jobId]);
  
  // 返回任務列表
  const handleBack = () => {
    router.push('/quantum/jobs');
  };
  
  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="flex flex-col items-center">
          <ArrowPathIcon className="h-12 w-12 text-indigo-500 animate-spin mb-4" />
          <h2 className="text-xl font-semibold text-gray-700 dark:text-gray-300">載入中...</h2>
        </div>
      </div>
    );
  }
  
  if (error || !job) {
    return (
      <div className="bg-white dark:bg-gray-900 shadow-lg rounded-lg p-6">
        <div className="flex items-center justify-center h-64 flex-col">
          <ExclamationCircleIcon className="h-12 w-12 text-red-500 mb-4" />
          <h2 className="text-xl font-semibold text-red-600 dark:text-red-400 mb-2">
            {error || '任務資訊無法載入'}
          </h2>
          <button
            onClick={handleBack}
            className="mt-4 px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
          >
            返回任務列表
          </button>
        </div>
      </div>
    );
  }
  
  // 格式化時間
  const formatDate = (dateString?: string) => {
    if (!dateString) return '尚未開始';
    return new Date(dateString).toLocaleString('zh-TW');
  };
  
  // 狀態顏色與圖示
  const getStatusInfo = (status: string) => {
    switch (status) {
      case 'pending':
        return {
          color: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300',
          icon: <ClockIcon className="h-5 w-5" />,
          text: '等待中'
        };
      case 'running':
        return {
          color: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300',
          icon: <PlayIcon className="h-5 w-5" />,
          text: '執行中'
        };
      case 'completed':
        return {
          color: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
          icon: <CheckCircleIcon className="h-5 w-5" />,
          text: '已完成'
        };
      case 'failed':
        return {
          color: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
          icon: <ExclamationCircleIcon className="h-5 w-5" />,
          text: '失敗'
        };
      default:
        return {
          color: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300',
          icon: <ClockIcon className="h-5 w-5" />,
          text: status
        };
    }
  };
  
  const statusInfo = getStatusInfo(job.status);
  
  return (
    <div className="bg-white dark:bg-gray-900 shadow-lg rounded-lg p-6 max-w-4xl mx-auto">
      <div className="mb-4 flex justify-between items-center">
        <button
          onClick={handleBack}
          className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200"
        >
          &larr; 返回任務列表
        </button>
        
        <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium ${statusInfo.color}`}>
          {statusInfo.icon}
          <span className="ml-1">{statusInfo.text}</span>
        </span>
      </div>
      
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">{job.title}</h1>
      
      {job.description && (
        <p className="text-gray-600 dark:text-gray-400 mb-6">{job.description}</p>
      )}
      
      {/* 任務狀態追蹤器 */}
      <div className="mb-8">
        <JobStatusTracker 
          status={job.status}
          createdAt={job.created_at}
          startedAt={job.started_at}
          completedAt={job.completed_at}
        />
      </div>
      
      {/* 任務詳情區塊 */}
      <div className="grid md:grid-cols-2 gap-6 mb-8">
        <div className="bg-gray-50 dark:bg-gray-800 p-4 rounded-md">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">任務資訊</h3>
          
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-sm text-gray-500 dark:text-gray-400">任務 ID</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">{job.id}</span>
            </div>
            
            <div className="flex justify-between">
              <span className="text-sm text-gray-500 dark:text-gray-400">建立時間</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">
                {formatDate(job.created_at)}
              </span>
            </div>
            
            <div className="flex justify-between">
              <span className="text-sm text-gray-500 dark:text-gray-400">開始時間</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">
                {formatDate(job.started_at)}
              </span>
            </div>
            
            <div className="flex justify-between">
              <span className="text-sm text-gray-500 dark:text-gray-400">完成時間</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">
                {formatDate(job.completed_at)}
              </span>
            </div>
            
            {job.execution_time_seconds !== undefined && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-500 dark:text-gray-400">執行時間</span>
                <span className="text-sm font-medium text-gray-900 dark:text-white">
                  {job.execution_time_seconds} 秒
                </span>
              </div>
            )}
            
            <div className="flex justify-between">
              <span className="text-sm text-gray-500 dark:text-gray-400">量子後端</span>
              <span className="text-sm font-medium text-gray-900 dark:text-white">
                {job.quantum_backend || '尚未指派'} 
                {job.is_simulation && ' (模擬器)'}
              </span>
            </div>
          </div>
        </div>
        
        <div className="bg-gray-50 dark:bg-gray-800 p-4 rounded-md">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">分析參數</h3>
          
          {job.input_params_summary ? (
            <div className="space-y-3">
              {Object.entries(job.input_params_summary).map(([key, value]) => (
                <div key={key} className="flex justify-between">
                  <span className="text-sm text-gray-500 dark:text-gray-400">
                    {key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())}
                  </span>
                  <span className="text-sm font-medium text-gray-900 dark:text-white">
                    {Array.isArray(value) ? value.join(', ') : String(value)}
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-4 text-gray-500 dark:text-gray-400">
              <p>無可用參數資訊</p>
            </div>
          )}
        </div>
      </div>
      
      {/* 分析結果區塊 */}
      {job.status === 'completed' && job.results && (
        <div className="bg-gray-50 dark:bg-gray-800 p-6 rounded-md mb-8">
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-6">分析結果</h3>
          
          <div className="grid md:grid-cols-2 gap-6">
            <div>
              <div className="mb-6">
                <div className="flex items-center mb-2">
                  <span className="text-sm text-gray-500 dark:text-gray-400">威脅判定</span>
                  <span className={`ml-auto px-3 py-1 rounded-full text-xs font-medium ${
                    job.is_malicious 
                      ? 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'
                      : 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'
                  }`}>
                    {job.is_malicious ? '惡意威脅' : '良性活動'}
                  </span>
                </div>
                
                <div className="flex items-center mb-2">
                  <span className="text-sm text-gray-500 dark:text-gray-400">信心分數</span>
                  <div className="ml-auto flex items-center">
                    <div className="w-40 bg-gray-200 dark:bg-gray-700 rounded-full h-2.5 mr-2">
                      <div
                        className={`h-2.5 rounded-full ${
                          job.confidence_score && job.confidence_score > 75 
                            ? 'w-[75%] bg-green-500'
                            : job.confidence_score && job.confidence_score > 50
                              ? 'w-[50%] bg-yellow-500'
                              : 'w-[25%] bg-red-500'
                        }`}
                      />
                    </div>
                    <span className="text-sm font-medium text-gray-900 dark:text-white">
                      {job.confidence_score?.toFixed(1) || 0}%
                    </span>
                  </div>
                </div>
                
                <div className="flex items-center">
                  <span className="text-sm text-gray-500 dark:text-gray-400">概率分布</span>
                  <span className="ml-auto text-sm font-medium text-gray-900 dark:text-white">
                    {(job.results.probability * 100).toFixed(1)}%
                  </span>
                </div>
              </div>
              
              {job.results_summary && (
                <div className="bg-white dark:bg-gray-700 p-4 rounded shadow-sm">
                  <h4 className="text-md font-medium text-gray-800 dark:text-white mb-2">摘要發現</h4>
                  <ul className="space-y-1 text-sm text-gray-600 dark:text-gray-300 list-disc list-inside">
                    {Object.entries(job.results_summary).map(([key, value]) => (
                      <li key={key}>
                        {key}: {String(value)}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
            
            <div>
              <ProbabilityChart 
                counts={job.results.counts} 
                isMalicious={!!job.is_malicious}
                confidenceScore={job.confidence_score || 0}
              />
            </div>
          </div>
        </div>
      )}
      
      {/* 錯誤訊息 */}
      {job.status === 'failed' && job.error_message && (
        <div className="bg-red-50 dark:bg-red-900/30 border-l-4 border-red-500 p-4 mb-6">
          <div className="flex">
            <ExclamationCircleIcon className="h-6 w-6 text-red-500 mr-3" />
            <div>
              <h3 className="text-md font-medium text-red-800 dark:text-red-300 mb-1">任務執行失敗</h3>
              <p className="text-sm text-red-700 dark:text-red-400">{job.error_message}</p>
            </div>
          </div>
        </div>
      )}
      
      {/* 按鈕區域 */}
      <div className="flex justify-end mt-8">
        {job.status === 'failed' && (
          <button
            onClick={() => router.push(`/quantum/jobs/new?clone=${job.id}`)}
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 mr-4"
          >
            以此任務為基礎重新提交
          </button>
        )}
        
        <button
          onClick={handleBack}
          className="px-4 py-2 bg-gray-200 text-gray-700 dark:bg-gray-700 dark:text-gray-200 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
        >
          返回任務列表
        </button>
      </div>
    </div>
  );
}
