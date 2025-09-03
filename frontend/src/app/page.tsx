"use client";

import { useState, useEffect } from "react";
import { ThreatIntelCard } from "@/components/ThreatIntelCard";
import { ThreatSearchForm } from "@/components/ThreatSearchForm";
import { ThreatStatsCards } from "@/components/ThreatStatsCards";
import { HIBPChecker } from "@/components/HIBPChecker";
import { ThreatCollector } from "@/components/ThreatCollector";

interface ThreatIntelligence {
  id: string;
  ip_address: string;
  domain?: string;
  threat_type: string;
  severity: string;
  confidence_score: number;
  description?: string;
  source: string;
  country_code?: string;
  created_at: string;
  updated_at: string;
}

interface ThreatStats {
  total_threats: number;
  high_severity_count: number;
  critical_severity_count: number;
  unique_sources: number;
  today_added: number;
}

export default function Dashboard() {
  const [threats, setThreats] = useState<ThreatIntelligence[]>([]);
  const [stats, setStats] = useState<ThreatStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [searchResults, setSearchResults] = useState<ThreatIntelligence[]>([]);
  const [isSearching, setIsSearching] = useState(false);

  useEffect(() => {
    loadInitialData();
  }, []);

  const loadInitialData = async () => {
    try {
      setLoading(true);

      // 並行載入威脅情報和統計資料
      const [threatsResponse, statsResponse] = await Promise.all([
        fetch(
          "/api/v1/threats?page=1&page_size=10&sort_by=created_at&sort_order=desc"
        ),
        fetch("/api/v1/threats/stats"),
      ]);

      if (threatsResponse.ok) {
        const threatsData = await threatsResponse.json();
        setThreats(threatsData.data?.items || []);
      }

      if (statsResponse.ok) {
        const statsData = await statsResponse.json();
        setStats(statsData.data || null);
      }
    } catch (error) {
      console.error("載入資料失敗:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (searchParams: {
    ip_address?: string;
    domain?: string;
    threat_type?: string;
    severity?: string;
    source?: string;
    country_code?: string;
  }) => {
    try {
      setIsSearching(true);

      const queryParams = new URLSearchParams();
      Object.entries(searchParams).forEach(([key, value]) => {
        if (value) {
          queryParams.append(key, value as string);
        }
      });

      const response = await fetch(`/api/v1/threats?${queryParams.toString()}`);

      if (response.ok) {
        const data = await response.json();
        setSearchResults(data.data?.items || []);
      } else {
        console.error("搜尋失敗");
        setSearchResults([]);
      }
    } catch (error) {
      console.error("搜尋錯誤:", error);
      setSearchResults([]);
    } finally {
      setIsSearching(false);
    }
  };

  const handleCollectionComplete = () => {
    // 收集完成後重新載入資料
    loadInitialData();
  };

  const displayThreats = searchResults.length > 0 ? searchResults : threats;

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 頁面標題 */}
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                資安威脅情報平台
              </h1>
              <p className="mt-2 text-gray-600">
                即時監控與分析威脅情報，保護您的數位資產
              </p>
            </div>
            <div className="flex items-center space-x-4">
              <div className="bg-green-100 text-green-800 px-3 py-1 rounded-full text-sm font-medium">
                系統正常
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* 統計卡片區域 */}
        <div className="mb-8">
          <ThreatStatsCards stats={stats} loading={loading} />
        </div>

        {/* 搜尋和收集器區域 */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
          {/* 威脅情報搜尋 */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">
                威脅情報搜尋
              </h2>
              <ThreatSearchForm onSearch={handleSearch} loading={isSearching} />
            </div>
          </div>

          {/* 威脅情報收集器 */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">
                威脅情報收集
              </h2>
              <ThreatCollector
                onCollectionComplete={handleCollectionComplete}
              />
            </div>
          </div>
        </div>

        {/* 威脅情報列表 */}
        <div className="bg-white rounded-lg shadow">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">
              {searchResults.length > 0 ? "搜尋結果" : "最新威脅情報"}
            </h2>
            {searchResults.length > 0 && (
              <button
                onClick={() => setSearchResults([])}
                className="mt-2 text-sm text-blue-600 hover:text-blue-800"
              >
                清除搜尋結果
              </button>
            )}
          </div>

          <div className="p-6">
            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
              </div>
            ) : displayThreats.length === 0 ? (
              <div className="text-center py-12">
                <div className="text-gray-500 text-lg mb-4">
                  {searchResults.length > 0
                    ? "沒有找到符合條件的威脅情報"
                    : "暫無威脅情報資料"}
                </div>
                <p className="text-gray-400">
                  {searchResults.length > 0
                    ? "請嘗試調整搜尋條件"
                    : "使用收集器開始收集威脅情報"}
                </p>
              </div>
            ) : (
              <div className="grid gap-4">
                {displayThreats.map((threat) => (
                  <ThreatIntelCard key={threat.id} threat={threat} />
                ))}
              </div>
            )}
          </div>
        </div>

        {/* HIBP 檢查器 */}
        <div className="mt-8">
          <HIBPChecker
            onCheckComplete={(result) => {
              console.log("HIBP 檢查完成:", result);
            }}
          />
        </div>
      </div>
    </div>
  );
}
