"use client";

import { useState } from "react";

interface ThreatCollectorProps {
  onCollectionComplete: () => void;
}

interface CollectionResult {
  success: boolean;
  message: string;
  data?: {
    ip_address?: string;
    total?: number;
    successful?: number;
    failed?: number;
    successful_ips?: string[];
    failed_ips?: string[];
  };
}

export function ThreatCollector({
  onCollectionComplete,
}: ThreatCollectorProps) {
  const [mode, setMode] = useState<"single" | "bulk">("single");
  const [singleIP, setSingleIP] = useState("");
  const [bulkIPs, setBulkIPs] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<CollectionResult | null>(null);

  const collectSingleIP = async () => {
    if (!singleIP.trim()) return;

    try {
      setLoading(true);
      setResult(null);

      const response = await fetch("/api/v1/collector/ip", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          ip_address: singleIP.trim(),
        }),
      });

      const data = await response.json();
      setResult(data);

      if (data.success) {
        setSingleIP("");
        onCollectionComplete();
      }
    } catch {
      setResult({
        success: false,
        message: "網路連接錯誤",
      });
    } finally {
      setLoading(false);
    }
  };

  const collectBulkIPs = async () => {
    const ips = bulkIPs
      .split("\n")
      .map((ip) => ip.trim())
      .filter((ip) => ip.length > 0);

    if (ips.length === 0) return;

    try {
      setLoading(true);
      setResult(null);

      const response = await fetch("/api/v1/collector/bulk-ip", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          ip_addresses: ips,
        }),
      });

      const data = await response.json();
      setResult(data);

      if (data.success) {
        setBulkIPs("");
        onCollectionComplete();
      }
    } catch {
      setResult({
        success: false,
        message: "網路連接錯誤",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (mode === "single") {
      collectSingleIP();
    } else {
      collectBulkIPs();
    }
  };

  return (
    <div className="space-y-4">
      {/* 模式選擇 */}
      <div className="flex space-x-1 bg-gray-100 p-1 rounded-lg">
        <button
          type="button"
          onClick={() => setMode("single")}
          className={`flex-1 px-3 py-2 text-sm font-medium rounded-md transition-colors ${
            mode === "single"
              ? "bg-white text-blue-600 shadow-sm"
              : "text-gray-500 hover:text-gray-700"
          }`}
        >
          單一 IP
        </button>
        <button
          type="button"
          onClick={() => setMode("bulk")}
          className={`flex-1 px-3 py-2 text-sm font-medium rounded-md transition-colors ${
            mode === "bulk"
              ? "bg-white text-blue-600 shadow-sm"
              : "text-gray-500 hover:text-gray-700"
          }`}
        >
          批量 IP
        </button>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        {mode === "single" ? (
          /* 單一 IP 收集 */
          <div>
            <label
              htmlFor="single-ip"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              IP 地址
            </label>
            <input
              type="text"
              id="single-ip"
              value={singleIP}
              onChange={(e) => setSingleIP(e.target.value)}
              placeholder="例如: 192.168.1.1"
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              disabled={loading}
            />
          </div>
        ) : (
          /* 批量 IP 收集 */
          <div>
            <label
              htmlFor="bulk-ips"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              IP 地址列表 (每行一個)
            </label>
            <textarea
              id="bulk-ips"
              value={bulkIPs}
              onChange={(e) => setBulkIPs(e.target.value)}
              placeholder="192.168.1.1&#10;10.0.0.1&#10;8.8.8.8"
              rows={6}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              disabled={loading}
            />
            <p className="mt-1 text-xs text-gray-500">最多 50 個 IP 地址</p>
          </div>
        )}

        <button
          type="submit"
          disabled={
            loading || (mode === "single" ? !singleIP.trim() : !bulkIPs.trim())
          }
          className="w-full bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
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
              收集中...
            </span>
          ) : (
            `🔄 開始收集${mode === "bulk" ? "批量" : ""}威脅情報`
          )}
        </button>
      </form>

      {/* 結果顯示 */}
      {result && (
        <div
          className={`p-4 rounded-md ${
            result.success
              ? "bg-green-50 border border-green-200"
              : "bg-red-50 border border-red-200"
          }`}
        >
          <div className="flex items-start">
            <div className="flex-shrink-0">
              {result.success ? (
                <svg
                  className="h-5 w-5 text-green-400"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                    clipRule="evenodd"
                  />
                </svg>
              ) : (
                <svg
                  className="h-5 w-5 text-red-400"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"
                  />
                </svg>
              )}
            </div>
            <div className="ml-3">
              <h3
                className={`text-sm font-medium ${
                  result.success ? "text-green-800" : "text-red-800"
                }`}
              >
                {result.success ? "收集成功" : "收集失敗"}
              </h3>
              <div
                className={`mt-2 text-sm ${
                  result.success ? "text-green-700" : "text-red-700"
                }`}
              >
                <p>{result.message}</p>
                {result.data && mode === "bulk" && (
                  <div className="mt-2 space-y-1">
                    <p>總數: {result.data.total}</p>
                    <p>成功: {result.data.successful}</p>
                    <p>失敗: {result.data.failed}</p>
                    {result.data.failed_ips &&
                      result.data.failed_ips.length > 0 && (
                        <p className="text-xs">
                          失敗的 IP: {result.data.failed_ips.join(", ")}
                        </p>
                      )}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* 說明文字 */}
      <div className="text-xs text-gray-500 space-y-1">
        <p>• 使用 AbuseIPDB API 收集威脅情報</p>
        <p>• 只收集信心分數 ≥ 25% 的威脅</p>
        <p>• 批量收集會自動限制頻率以避免 API 限制</p>
      </div>
    </div>
  );
}
