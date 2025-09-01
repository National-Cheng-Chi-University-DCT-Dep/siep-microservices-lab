import React from 'react';
import { CheckIcon, ClockIcon, ArrowPathIcon, ExclamationTriangleIcon } from '@heroicons/react/24/solid';

interface JobStatusTrackerProps {
  status: 'pending' | 'running' | 'completed' | 'failed';
  createdAt: string;
  startedAt?: string;
  completedAt?: string;
}

const JobStatusTracker: React.FC<JobStatusTrackerProps> = ({
  status,
  createdAt,
  startedAt,
  completedAt
}) => {
  // 計算每個狀態的樣式
  const getStepStatus = (step: string) => {
    switch (step) {
      case 'created':
        return 'completed'; // 創建步驟總是完成的
      
      case 'started':
        if (status === 'running' || status === 'completed' || status === 'failed') {
          return 'completed';
        }
        return 'upcoming';
      
      case 'completed':
        if (status === 'completed') {
          return 'completed';
        } else if (status === 'failed') {
          return 'failed';
        } else if (status === 'running') {
          return 'current';
        }
        return 'upcoming';
      
      default:
        return 'upcoming';
    }
  };

  // 獲取圖標
  const getStepIcon = (step: string) => {
    const stepStatus = getStepStatus(step);
    
    if (step === 'completed' && status === 'failed') {
      return <ExclamationTriangleIcon className="h-6 w-6 text-red-500" />;
    }
    
    switch (stepStatus) {
      case 'completed':
        return <CheckIcon className="h-6 w-6 text-green-500" />;
      case 'current':
        return <ArrowPathIcon className="h-6 w-6 text-blue-500 animate-spin" />;
      case 'failed':
        return <ExclamationTriangleIcon className="h-6 w-6 text-red-500" />;
      default:
        return <ClockIcon className="h-6 w-6 text-gray-400" />;
    }
  };

  // 獲取圖標背景色
  const getStepIconBackground = (step: string) => {
    const stepStatus = getStepStatus(step);
    
    if (step === 'completed' && status === 'failed') {
      return 'bg-red-100 dark:bg-red-900/30';
    }
    
    switch (stepStatus) {
      case 'completed':
        return 'bg-green-100 dark:bg-green-900/30';
      case 'current':
        return 'bg-blue-100 dark:bg-blue-900/30';
      case 'failed':
        return 'bg-red-100 dark:bg-red-900/30';
      default:
        return 'bg-gray-100 dark:bg-gray-800';
    }
  };

  // 獲取連接線的顏色
  const getConnectorColor = (step: string) => {
    const stepStatus = getStepStatus(step);
    
    if (stepStatus === 'completed') {
      return 'bg-green-500 dark:bg-green-400';
    } else if (stepStatus === 'current') {
      return 'bg-blue-500 dark:bg-blue-400';
    }
    
    return 'bg-gray-300 dark:bg-gray-600';
  };

  // 格式化時間
  const formatTime = (dateString?: string) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleTimeString('zh-TW', { 
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  };
  
  // 計算自動時間估計 (適用於尚未完成的任務)
  const getEstimatedTime = (step: string) => {
    if (step === 'completed' && !completedAt) {
      if (startedAt) {
        const startTime = new Date(startedAt).getTime();
        const now = new Date().getTime();
        const elapsedSeconds = Math.floor((now - startTime) / 1000);
        
        if (status === 'running') {
          // 如果正在執行中，顯示已執行的時間
          return `執行中 (${elapsedSeconds} 秒)`;
        }
      }
      return '待完成';
    }
    return '';
  };

  const steps = [
    {
      id: 'created',
      name: '任務建立',
      time: formatTime(createdAt)
    },
    {
      id: 'started',
      name: '開始執行',
      time: startedAt ? formatTime(startedAt) : '等待中...'
    },
    {
      id: 'completed',
      name: status === 'failed' ? '執行失敗' : '完成分析',
      time: completedAt ? formatTime(completedAt) : getEstimatedTime('completed')
    }
  ];

  return (
    <div className="py-4">
      <ol className="relative flex items-center justify-between w-full">
        {steps.map((step, index) => (
          <React.Fragment key={step.id}>
            <li className="flex flex-col items-center">
              <div className={`flex items-center justify-center w-10 h-10 rounded-full ${getStepIconBackground(step.id)}`}>
                {getStepIcon(step.id)}
              </div>
              <div className="mt-2 text-center">
                <h3 className="text-sm font-medium text-gray-900 dark:text-white">
                  {step.name}
                </h3>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  {step.time}
                </p>
              </div>
            </li>
            
            {index < steps.length - 1 && (
              <div className="flex-grow mx-2">
                <div className={`h-1 ${getConnectorColor(steps[index + 1].id)}`}></div>
              </div>
            )}
          </React.Fragment>
        ))}
      </ol>
    </div>
  );
};

export default JobStatusTracker;
