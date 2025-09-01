-- 測試資料腳本
-- 插入一些範例威脅情報用於測試

INSERT INTO threat_intelligence (
    id,
    ip_address,
    domain,
    threat_type,
    severity,
    confidence_score,
    description,
    source,
    country_code,
    tags,
    metadata,
    created_at,
    updated_at
) VALUES 
(
    gen_random_uuid(),
    '192.168.1.100',
    'malicious.example.com',
    'malware',
    'high',
    85,
    '檢測到惡意軟體活動，IP 地址涉及 Trojan 下載',
    'TestData',
    'US',
    ARRAY['test', 'malware', 'trojan'],
    '{"test": true, "source_confidence": "high", "last_activity": "2024-12-01"}',
    NOW() - INTERVAL '1 hour',
    NOW() - INTERVAL '1 hour'
),
(
    gen_random_uuid(),
    '10.0.0.50',
    NULL,
    'scanner',
    'medium',
    65,
    '檢測到端口掃描活動',
    'TestData',
    'CN',
    ARRAY['test', 'scanner', 'port-scan'],
    '{"test": true, "scan_type": "port", "ports_scanned": [22, 80, 443]}',
    NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '2 hours'
),
(
    gen_random_uuid(),
    '203.0.113.15',
    'phishing.bad-site.com',
    'phishing',
    'critical',
    95,
    '確認的釣魚網站，偽造知名銀行登入頁面',
    'TestData',
    'RU',
    ARRAY['test', 'phishing', 'banking'],
    '{"test": true, "target_bank": "example_bank", "active": true}',
    NOW() - INTERVAL '30 minutes',
    NOW() - INTERVAL '30 minutes'
),
(
    gen_random_uuid(),
    '198.51.100.25',
    NULL,
    'ddos',
    'high',
    78,
    'DDoS 攻擊來源 IP',
    'TestData',
    'KR',
    ARRAY['test', 'ddos', 'botnet'],
    '{"test": true, "attack_type": "udp_flood", "volume_gbps": 5.2}',
    NOW() - INTERVAL '3 hours',
    NOW() - INTERVAL '3 hours'
),
(
    gen_random_uuid(),
    '172.16.0.75',
    'spam.sender.net',
    'spam',
    'low',
    35,
    '垃圾郵件發送者',
    'TestData',
    'BR',
    ARRAY['test', 'spam', 'email'],
    '{"test": true, "spam_type": "commercial", "volume_per_day": 1000}',
    NOW() - INTERVAL '4 hours',
    NOW() - INTERVAL '4 hours'
);

-- 檢查插入的資料
SELECT 
    ip_address,
    threat_type,
    severity,
    confidence_score,
    source,
    country_code,
    created_at
FROM threat_intelligence 
WHERE source = 'TestData'
ORDER BY created_at DESC; 