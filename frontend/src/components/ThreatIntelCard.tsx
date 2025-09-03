"use client";

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

interface ThreatIntelCardProps {
  threat: ThreatIntelligence;
}

export function ThreatIntelCard({ threat }: ThreatIntelCardProps) {
  const getSeverityColor = (severity: string) => {
    switch (severity.toLowerCase()) {
      case "critical":
        return "bg-red-100 text-red-800 border-red-200";
      case "high":
        return "bg-orange-100 text-orange-800 border-orange-200";
      case "medium":
        return "bg-yellow-100 text-yellow-800 border-yellow-200";
      case "low":
        return "bg-green-100 text-green-800 border-green-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const getThreatTypeIcon = (threatType: string) => {
    switch (threatType.toLowerCase()) {
      case "malware":
        return "🦠";
      case "phishing":
        return "🎣";
      case "spam":
        return "📧";
      case "botnet":
        return "🤖";
      case "scanner":
        return "🔍";
      case "ddos":
        return "💥";
      case "bruteforce":
        return "🔨";
      default:
        return "⚠️";
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("zh-TW", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getConfidenceColor = (score: number) => {
    if (score >= 80) return "text-red-600";
    if (score >= 60) return "text-orange-600";
    if (score >= 40) return "text-yellow-600";
    return "text-green-600";
  };

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center space-x-3">
          <div className="text-2xl">
            {getThreatTypeIcon(threat.threat_type)}
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-900">
              {threat.ip_address}
            </h3>
            {threat.domain && (
              <p className="text-sm text-gray-600">域名: {threat.domain}</p>
            )}
          </div>
        </div>
        <div className="flex items-center space-x-2">
          <span
            className={`px-2 py-1 text-xs font-medium rounded-full border ${getSeverityColor(
              threat.severity
            )}`}
          >
            {threat.severity.toUpperCase()}
          </span>
          {threat.country_code && (
            <span className="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-800 rounded-full border border-blue-200">
              {threat.country_code}
            </span>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div>
          <p className="text-sm text-gray-600 mb-1">威脅類型</p>
          <p className="font-medium text-gray-900 capitalize">
            {threat.threat_type}
          </p>
        </div>
        <div>
          <p className="text-sm text-gray-600 mb-1">信心分數</p>
          <p
            className={`font-medium ${getConfidenceColor(
              threat.confidence_score
            )}`}
          >
            {threat.confidence_score}%
          </p>
        </div>
        <div>
          <p className="text-sm text-gray-600 mb-1">情報來源</p>
          <p className="font-medium text-gray-900">{threat.source}</p>
        </div>
        <div>
          <p className="text-sm text-gray-600 mb-1">發現時間</p>
          <p className="font-medium text-gray-900">
            {formatDate(threat.created_at)}
          </p>
        </div>
      </div>

      {threat.description && (
        <div className="mb-4">
          <p className="text-sm text-gray-600 mb-1">描述</p>
          <p className="text-gray-900 text-sm">{threat.description}</p>
        </div>
      )}

      <div className="flex items-center justify-between pt-4 border-t border-gray-200">
        <div className="flex items-center space-x-4">
          <button className="text-blue-600 hover:text-blue-800 text-sm font-medium">
            查看詳情
          </button>
          <button className="text-gray-600 hover:text-gray-800 text-sm font-medium">
            加入黑名單
          </button>
        </div>
        <div className="text-xs text-gray-500">ID: {threat.id.slice(-8)}</div>
      </div>
    </div>
  );
}
