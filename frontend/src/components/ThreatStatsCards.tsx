"use client";

interface ThreatStats {
  total_threats: number;
  high_severity_count: number;
  critical_severity_count: number;
  unique_sources: number;
  today_added: number;
}

interface ThreatStatsCardsProps {
  stats: ThreatStats | null;
  loading: boolean;
}

export function ThreatStatsCards({ stats, loading }: ThreatStatsCardsProps) {
  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="bg-white rounded-lg shadow p-6">
            <div className="animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-3/4 mb-3"></div>
              <div className="h-8 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="bg-white rounded-lg shadow p-6">
            <div className="text-center text-gray-500">
              <div className="text-sm">無資料</div>
              <div className="text-2xl font-bold">--</div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  const cards = [
    {
      title: "總威脅情報",
      value: stats.total_threats,
      icon: "🛡️",
      color: "blue",
      bgColor: "bg-blue-50",
      textColor: "text-blue-900",
      iconColor: "text-blue-600",
    },
    {
      title: "關鍵威脅",
      value: stats.critical_severity_count,
      icon: "🚨",
      color: "red",
      bgColor: "bg-red-50",
      textColor: "text-red-900",
      iconColor: "text-red-600",
    },
    {
      title: "高危威脅",
      value: stats.high_severity_count,
      icon: "⚠️",
      color: "orange",
      bgColor: "bg-orange-50",
      textColor: "text-orange-900",
      iconColor: "text-orange-600",
    },
    {
      title: "今日新增",
      value: stats.today_added,
      icon: "📊",
      color: "green",
      bgColor: "bg-green-50",
      textColor: "text-green-900",
      iconColor: "text-green-600",
    },
    {
      title: "情報來源",
      value: stats.unique_sources,
      icon: "🔗",
      color: "purple",
      bgColor: "bg-purple-50",
      textColor: "text-purple-900",
      iconColor: "text-purple-600",
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
      {cards.map((card, index) => (
        <div
          key={index}
          className={`${card.bgColor} rounded-lg shadow-sm border border-gray-100 p-6 transition-transform hover:scale-105`}
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600 mb-1">
                {card.title}
              </p>
              <p className={`text-2xl font-bold ${card.textColor}`}>
                {card.value.toLocaleString()}
              </p>
            </div>
            <div className={`text-2xl ${card.iconColor}`}>{card.icon}</div>
          </div>
        </div>
      ))}
    </div>
  );
}
