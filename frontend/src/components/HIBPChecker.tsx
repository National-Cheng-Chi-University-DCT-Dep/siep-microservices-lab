"use client";

import React, { useState } from "react";

interface HIBPCheckerProps {
  onCheckComplete?: (result: {
    breaches?: BreachResult[];
    pastes?: PasteResult[];
    pwned_count?: number;
    is_pwned?: boolean;
    domain?: string;
  }) => void;
}

interface BreachResult {
  Name: string;
  Title: string;
  Domain: string;
  BreachDate: string;
  PwnCount: number;
  Description: string;
  DataClasses: string[];
  IsVerified: boolean;
  IsSensitive: boolean;
}

interface PasteResult {
  Source: string;
  Id: string;
  Title: string;
  Date: string;
  EmailCount: number;
}

export function HIBPChecker({ onCheckComplete }: HIBPCheckerProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [domain, setDomain] = useState("");
  const [loading, setLoading] = useState(false);
  const [breachResults, setBreachResults] = useState<BreachResult[]>([]);
  const [pasteResults, setPasteResults] = useState<PasteResult[]>([]);
  const [passwordResult, setPasswordResult] = useState<{
    pwned_count: number;
    is_pwned: boolean;
  } | null>(null);
  const [domainResults, setDomainResults] = useState<{
    breaches: Record<string, string[]>;
    count: number;
  } | null>(null);
  const [error, setError] = useState<string | null>(null);

  const checkAccountBreaches = async () => {
    if (!email) {
      setError("請輸入電子郵件地址");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `/api/v1/hibp/account/${encodeURIComponent(email)}/breaches`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("access_token")}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.success) {
        setBreachResults(data.data.breaches || []);
        onCheckComplete?.(data.data);
      } else {
        setError(data.message || "檢查失敗");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "檢查失敗");
    } finally {
      setLoading(false);
    }
  };

  const checkAccountPastes = async () => {
    if (!email) {
      setError("請輸入電子郵件地址");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `/api/v1/hibp/account/${encodeURIComponent(email)}/pastes`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("access_token")}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.success) {
        setPasteResults(data.data.pastes || []);
        onCheckComplete?.(data.data);
      } else {
        setError(data.message || "檢查失敗");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "檢查失敗");
    } finally {
      setLoading(false);
    }
  };

  const checkPassword = async () => {
    if (!password) {
      setError("請輸入密碼");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `/api/v1/hibp/password/check?password=${encodeURIComponent(password)}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.success) {
        setPasswordResult(data.data);
        onCheckComplete?.(data.data);
      } else {
        setError(data.message || "檢查失敗");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "檢查失敗");
    } finally {
      setLoading(false);
    }
  };

  const checkDomainBreaches = async () => {
    if (!domain) {
      setError("請輸入域名");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `/api/v1/hibp/domain/${encodeURIComponent(domain)}/breaches`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("access_token")}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.success) {
        setDomainResults(data.data);
        onCheckComplete?.(data.data);
      } else {
        setError(data.message || "檢查失敗");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "檢查失敗");
    } finally {
      setLoading(false);
    }
  };

  const getSeverityColor = (pwnCount: number) => {
    if (pwnCount > 10000000) return "text-red-600";
    if (pwnCount > 1000000) return "text-orange-600";
    if (pwnCount > 100000) return "text-yellow-600";
    return "text-green-600";
  };

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">
          Have I Been Pwned 檢查器
        </h2>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {/* 帳戶泄露檢查 */}
        <div className="mb-8">
          <h3 className="text-lg font-semibold text-gray-700 mb-4">
            檢查帳戶泄露
          </h3>
          <div className="flex gap-4 mb-4">
            <input
              type="email"
              placeholder="輸入電子郵件地址"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={checkAccountBreaches}
              disabled={loading}
              className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? "檢查中..." : "檢查泄露"}
            </button>
            <button
              onClick={checkAccountPastes}
              disabled={loading}
              className="px-6 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50"
            >
              {loading ? "檢查中..." : "檢查 Paste"}
            </button>
          </div>

          {/* 泄露結果 */}
          {breachResults.length > 0 && (
            <div className="mt-4">
              <h4 className="font-semibold text-gray-700 mb-2">
                發現 {breachResults.length} 個泄露事件：
              </h4>
              <div className="space-y-3">
                {breachResults.map((breach, index) => (
                  <div
                    key={index}
                    className="border border-gray-200 rounded-lg p-4"
                  >
                    <div className="flex justify-between items-start">
                      <div>
                        <h5 className="font-semibold text-gray-800">
                          {breach.Title}
                        </h5>
                        <p className="text-sm text-gray-600">
                          域名: {breach.Domain}
                        </p>
                        <p className="text-sm text-gray-600">
                          泄露日期: {breach.BreachDate}
                        </p>
                        <p className="text-sm text-gray-600">
                          受影響帳戶: {breach.PwnCount.toLocaleString()}
                        </p>
                      </div>
                      <div className="text-right">
                        <span
                          className={`text-sm font-semibold ${getSeverityColor(
                            breach.PwnCount
                          )}`}
                        >
                          {breach.PwnCount > 10000000
                            ? "極高風險"
                            : breach.PwnCount > 1000000
                            ? "高風險"
                            : breach.PwnCount > 100000
                            ? "中風險"
                            : "低風險"}
                        </span>
                        {breach.IsVerified && (
                          <div className="text-xs text-green-600 mt-1">
                            已驗證
                          </div>
                        )}
                        {breach.IsSensitive && (
                          <div className="text-xs text-red-600 mt-1">
                            敏感數據
                          </div>
                        )}
                      </div>
                    </div>
                    <p className="text-sm text-gray-700 mt-2">
                      {breach.Description}
                    </p>
                    <div className="mt-2">
                      <span className="text-xs text-gray-500">
                        泄露數據類型:{" "}
                      </span>
                      {breach.DataClasses.map((dataClass, i) => (
                        <span
                          key={i}
                          className="inline-block bg-gray-100 text-gray-700 text-xs px-2 py-1 rounded mr-1 mb-1"
                        >
                          {dataClass}
                        </span>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Paste 結果 */}
          {pasteResults.length > 0 && (
            <div className="mt-4">
              <h4 className="font-semibold text-gray-700 mb-2">
                發現 {pasteResults.length} 個 Paste 記錄：
              </h4>
              <div className="space-y-3">
                {pasteResults.map((paste, index) => (
                  <div
                    key={index}
                    className="border border-gray-200 rounded-lg p-4"
                  >
                    <div className="flex justify-between items-start">
                      <div>
                        <h5 className="font-semibold text-gray-800">
                          {paste.Title}
                        </h5>
                        <p className="text-sm text-gray-600">
                          來源: {paste.Source}
                        </p>
                        <p className="text-sm text-gray-600">
                          日期: {new Date(paste.Date).toLocaleDateString()}
                        </p>
                        <p className="text-sm text-gray-600">
                          受影響郵箱: {paste.EmailCount.toLocaleString()}
                        </p>
                      </div>
                      <div className="text-right">
                        <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                          Paste
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* 密碼檢查 */}
        <div className="mb-8">
          <h3 className="text-lg font-semibold text-gray-700 mb-4">
            檢查密碼安全性
          </h3>
          <div className="flex gap-4 mb-4">
            <input
              type="password"
              placeholder="輸入密碼進行檢查"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={checkPassword}
              disabled={loading}
              className="px-6 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 disabled:opacity-50"
            >
              {loading ? "檢查中..." : "檢查密碼"}
            </button>
          </div>

          {passwordResult && (
            <div className="mt-4">
              <div
                className={`border rounded-lg p-4 ${
                  passwordResult.is_pwned
                    ? "border-red-200 bg-red-50"
                    : "border-green-200 bg-green-50"
                }`}
              >
                <div className="flex items-center">
                  <div
                    className={`w-4 h-4 rounded-full mr-3 ${
                      passwordResult.is_pwned ? "bg-red-500" : "bg-green-500"
                    }`}
                  ></div>
                  <div>
                    <h4
                      className={`font-semibold ${
                        passwordResult.is_pwned
                          ? "text-red-800"
                          : "text-green-800"
                      }`}
                    >
                      {passwordResult.is_pwned ? "密碼已被泄露" : "密碼安全"}
                    </h4>
                    <p
                      className={`text-sm ${
                        passwordResult.is_pwned
                          ? "text-red-700"
                          : "text-green-700"
                      }`}
                    >
                      {passwordResult.is_pwned
                        ? `此密碼在 ${passwordResult.pwned_count.toLocaleString()} 個泄露事件中被發現`
                        : "此密碼未在已知泄露事件中被發現"}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* 域名檢查 */}
        <div className="mb-8">
          <h3 className="text-lg font-semibold text-gray-700 mb-4">
            檢查域名泄露
          </h3>
          <div className="flex gap-4 mb-4">
            <input
              type="text"
              placeholder="輸入域名 (例如: example.com)"
              value={domain}
              onChange={(e) => setDomain(e.target.value)}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={checkDomainBreaches}
              disabled={loading}
              className="px-6 py-2 bg-orange-600 text-white rounded-md hover:bg-orange-700 disabled:opacity-50"
            >
              {loading ? "檢查中..." : "檢查域名"}
            </button>
          </div>

          {domainResults && (
            <div className="mt-4">
              <h4 className="font-semibold text-gray-700 mb-2">
                域名 {domain} 的泄露情況：
              </h4>
              <div className="border border-gray-200 rounded-lg p-4">
                <p className="text-sm text-gray-600 mb-2">
                  發現 {Object.keys(domainResults.breaches).length}{" "}
                  個受影響的郵箱地址
                </p>
                <div className="space-y-2">
                  {Object.entries(domainResults.breaches).map(
                    ([email, breaches]: [string, string[]]) => (
                      <div
                        key={email}
                        className="border-l-4 border-blue-500 pl-3"
                      >
                        <p className="font-medium text-gray-800">{email}</p>
                        <p className="text-sm text-gray-600">
                          涉及泄露事件:{" "}
                          {Array.isArray(breaches)
                            ? breaches.join(", ")
                            : breaches}
                        </p>
                      </div>
                    )
                  )}
                </div>
              </div>
            </div>
          )}
        </div>

        {/* 安全建議 */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h4 className="font-semibold text-blue-800 mb-2">安全建議</h4>
          <ul className="text-sm text-blue-700 space-y-1">
            <li>• 如果發現帳戶泄露，請立即更改密碼</li>
            <li>• 使用強密碼並啟用雙因素認證</li>
            <li>• 定期檢查帳戶安全狀態</li>
            <li>• 避免在多個網站使用相同密碼</li>
            <li>• 考慮使用密碼管理器</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
