"use client";

import { useState } from "react";

interface SearchParams {
  ip_address?: string;
  domain?: string;
  threat_type?: string;
  severity?: string;
  source?: string;
  country_code?: string;
}

interface ThreatSearchFormProps {
  onSearch: (params: SearchParams) => void;
  loading: boolean;
}

export function ThreatSearchForm({ onSearch, loading }: ThreatSearchFormProps) {
  const [searchParams, setSearchParams] = useState<SearchParams>({});

  const handleInputChange = (key: keyof SearchParams, value: string) => {
    setSearchParams((prev) => ({
      ...prev,
      [key]: value || undefined,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(searchParams);
  };

  const handleReset = () => {
    setSearchParams({});
    onSearch({});
  };

  const threatTypes = [
    { value: "", label: "全部類型" },
    { value: "malware", label: "惡意軟體" },
    { value: "phishing", label: "釣魚網站" },
    { value: "spam", label: "垃圾郵件" },
    { value: "botnet", label: "殭屍網路" },
    { value: "scanner", label: "掃描器" },
    { value: "ddos", label: "DDoS 攻擊" },
    { value: "bruteforce", label: "暴力破解" },
    { value: "other", label: "其他" },
  ];

  const severityLevels = [
    { value: "", label: "全部等級" },
    { value: "critical", label: "關鍵" },
    { value: "high", label: "高危" },
    { value: "medium", label: "中危" },
    { value: "low", label: "低危" },
  ];

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* IP 地址搜尋 */}
        <div>
          <label
            htmlFor="ip_address"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            IP 地址
          </label>
          <input
            type="text"
            id="ip_address"
            value={searchParams.ip_address || ""}
            onChange={(e) => handleInputChange("ip_address", e.target.value)}
            placeholder="例如: 192.168.1.1"
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* 域名搜尋 */}
        <div>
          <label
            htmlFor="domain"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            域名
          </label>
          <input
            type="text"
            id="domain"
            value={searchParams.domain || ""}
            onChange={(e) => handleInputChange("domain", e.target.value)}
            placeholder="例如: malicious.com"
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* 威脅類型 */}
        <div>
          <label
            htmlFor="threat_type"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            威脅類型
          </label>
          <select
            id="threat_type"
            value={searchParams.threat_type || ""}
            onChange={(e) => handleInputChange("threat_type", e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            {threatTypes.map((type) => (
              <option key={type.value} value={type.value}>
                {type.label}
              </option>
            ))}
          </select>
        </div>

        {/* 嚴重程度 */}
        <div>
          <label
            htmlFor="severity"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            嚴重程度
          </label>
          <select
            id="severity"
            value={searchParams.severity || ""}
            onChange={(e) => handleInputChange("severity", e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            {severityLevels.map((level) => (
              <option key={level.value} value={level.value}>
                {level.label}
              </option>
            ))}
          </select>
        </div>

        {/* 情報來源 */}
        <div>
          <label
            htmlFor="source"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            情報來源
          </label>
          <input
            type="text"
            id="source"
            value={searchParams.source || ""}
            onChange={(e) => handleInputChange("source", e.target.value)}
            placeholder="例如: AbuseIPDB"
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* 國家代碼 */}
        <div>
          <label
            htmlFor="country_code"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            國家代碼
          </label>
          <input
            type="text"
            id="country_code"
            value={searchParams.country_code || ""}
            onChange={(e) => handleInputChange("country_code", e.target.value)}
            placeholder="例如: US, CN, RU"
            maxLength={2}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
      </div>

      {/* 按鈕區域 */}
      <div className="flex space-x-3 pt-4">
        <button
          type="submit"
          disabled={loading}
          className="flex-1 bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {loading ? (
            <span className="flex items-center justify-center">
              <svg
                className="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              搜尋中...
            </span>
          ) : (
            "🔍 搜尋"
          )}
        </button>

        <button
          type="button"
          onClick={handleReset}
          className="px-6 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
        >
          清除
        </button>
      </div>
    </form>
  );
}
