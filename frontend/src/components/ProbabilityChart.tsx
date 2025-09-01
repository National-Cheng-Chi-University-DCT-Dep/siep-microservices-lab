import React, { useEffect, useRef } from 'react';
import Chart from 'chart.js/auto';

interface ProbabilityChartProps {
  counts: Record<string, number>;
  isMalicious: boolean;
  confidenceScore: number;
}

const ProbabilityChart: React.FC<ProbabilityChartProps> = ({ counts, isMalicious, confidenceScore }) => {
  const chartRef = useRef<HTMLCanvasElement>(null);
  const chartInstance = useRef<Chart | null>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    // 如果已經有圖表實例，則先銷毀它
    if (chartInstance.current) {
      chartInstance.current.destroy();
    }

    // 解析數據
    const labels: string[] = [];
    const data: number[] = [];
    let total = 0;

    Object.entries(counts).forEach(([key, value]) => {
      labels.push(key);
      data.push(value);
      total += value;
    });

    // 計算百分比
    const dataPercentage = data.map(value => (value / total * 100).toFixed(1));

    // 設置顏色
    const baseColors = isMalicious 
      ? ['rgba(239, 68, 68, 0.7)', 'rgba(249, 115, 22, 0.7)', 'rgba(234, 179, 8, 0.7)'] // 紅色系
      : ['rgba(34, 197, 94, 0.7)', 'rgba(16, 185, 129, 0.7)', 'rgba(6, 182, 212, 0.7)']; // 綠色系
    
    const borderColors = isMalicious
      ? ['rgba(220, 38, 38, 1)', 'rgba(234, 88, 12, 1)', 'rgba(202, 138, 4, 1)']
      : ['rgba(22, 163, 74, 1)', 'rgba(5, 150, 105, 1)', 'rgba(8, 145, 178, 1)'];

    // 確保顏色數量足夠
    const backgroundColor = Array(labels.length).fill('').map((_, i) => 
      baseColors[i % baseColors.length]
    );
    
    const borderColor = Array(labels.length).fill('').map((_, i) => 
      borderColors[i % borderColors.length]
    );

    // 創建圖表
    const ctx = chartRef.current.getContext('2d');
    if (ctx) {
      chartInstance.current = new Chart(ctx, {
        type: 'bar',
        data: {
          labels,
          datasets: [
            {
              label: '觀測機率 (%)',
              data: dataPercentage,
              backgroundColor,
              borderColor,
              borderWidth: 1,
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            tooltip: {
              callbacks: {
                label: function(context) {
                  const value = context.raw;
                  const count = data[context.dataIndex];
                  return `${value}% (${count} counts)`;
                }
              }
            },
            legend: {
              display: false
            },
            title: {
              display: true,
              text: '量子測量結果分布',
              color: isMalicious ? '#dc2626' : '#16a34a',
              font: {
                size: 16
              }
            },
            subtitle: {
              display: true,
              text: `信心分數: ${confidenceScore.toFixed(1)}%`,
              color: '#6b7280',
              font: {
                size: 14,
                weight: 'normal'
              },
              padding: {
                bottom: 10
              }
            }
          },
          scales: {
            y: {
              beginAtZero: true,
              title: {
                display: true,
                text: '機率 (%)',
              },
              ticks: {
                callback: function(value) {
                  return value + '%';
                }
              }
            },
            x: {
              title: {
                display: true,
                text: '量子測量結果',
              }
            }
          }
        }
      });
    }

    return () => {
      if (chartInstance.current) {
        chartInstance.current.destroy();
      }
    };
  }, [counts, isMalicious, confidenceScore]);

  return (
    <div className="w-full h-72">
      <canvas ref={chartRef}></canvas>
    </div>
  );
};

export default ProbabilityChart;
