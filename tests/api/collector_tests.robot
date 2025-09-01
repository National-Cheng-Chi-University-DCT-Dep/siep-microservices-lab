*** Settings ***
Documentation    收集器API測試套件
Resource         ../config/test_config.robot
Suite Setup      Suite Setup With Auth
Suite Teardown   Teardown Test Environment
Test Setup       Setup Test Data
Test Teardown    Cleanup Test Data

*** Keywords ***
Suite Setup With Auth
    [Documentation]    設定測試環境並登入
    Setup Test Environment
    Login User

*** Test Cases ***

Test Collect Single IP Success
    [Documentation]    測試單一IP收集成功
    [Tags]    collector    ip    single    positive
    
    # 使用已知惡意IP進行測試（這是一個測試用的IP）
    ${test_ip}=    Set Variable    118.25.6.39
    
    # 發送收集請求
    &{collect_data}=    Create Dictionary    ip_address=${test_ip}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ip    ${collect_data}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}    10    # 給收集器更多時間
    
    # 驗證回應資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Not Be Empty    ${data[0]['threat_intelligence']}
    Should Be Equal    ${data[0]['threat_intelligence']['ip_address']}    ${test_ip}
    Should Not Be Empty    ${data[0]['threat_intelligence']['source']}
    
    # 驗證威脅情報已保存到資料庫
    Verify Database Record Exists    threat_intelligence    ip_address='${test_ip}'

Test Collect Single IP Invalid Format
    [Documentation]    測試單一IP收集無效格式
    [Tags]    collector    ip    single    negative
    
    # 使用無效IP格式
    ${invalid_ip}=    Set Variable    256.256.256.256
    
    &{collect_data}=    Create Dictionary    ip_address=${invalid_ip}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ip    ${collect_data}    expected_status=400
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    400

Test Collect Single IP Not Found
    [Documentation]    測試收集不在威脅資料庫中的IP
    [Tags]    collector    ip    single    not_found
    
    # 使用Google DNS IP（通常不會在威脅資料庫中）
    ${clean_ip}=    Set Variable    8.8.8.8
    
    &{collect_data}=    Create Dictionary    ip_address=${clean_ip}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ip    ${collect_data}
    
    # 驗證回應（可能返回404或200但沒有威脅資料）
    Run Keyword If    ${response.status_code} == 200
    ...    Should Be Empty    ${response.json()['data']['threat_intelligence']}
    ...    ELSE IF    ${response.status_code} == 404
    ...    Verify Response Error    ${response}    404

Test Collect Multiple IPs Success
    [Documentation]    測試多個IP收集成功
    [Tags]    collector    ip    batch    positive
    
    # 準備多個測試IP
    @{test_ips}=    Create List    118.25.6.39    198.71.233.103
    
    # 發送批量收集請求
    &{collect_data}=    Create Dictionary    ip_addresses=${test_ips}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ips    ${collect_data}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}    15    # 批量收集需要更多時間
    
    # 驗證回應資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Not Be Empty    ${data[0]['results']}
    Should Be True    ${data[0]['total_processed']} >= 1
    
    # 檢查至少有一個成功的結果
    ${results}=    Get Value From Json    ${response.json()}    $.data.results
    ${success_count}=    Set Variable    0
    FOR    ${result}    IN    @{results[0]}
        Run Keyword If    '${result['status']}' == 'success'
        ...    Set Variable    ${success_count + 1}
    END
    Should Be True    ${success_count} > 0

Test Collect Multiple IPs With Mixed Results
    [Documentation]    測試多個IP收集混合結果
    [Tags]    collector    ip    batch    mixed
    
    # 準備混合IP列表（有效和無效）
    @{mixed_ips}=    Create List    118.25.6.39    8.8.8.8    256.256.256.256
    
    # 發送批量收集請求
    &{collect_data}=    Create Dictionary    ip_addresses=${mixed_ips}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ips    ${collect_data}
    
    # 驗證回應（部分成功的情況下仍返回200）
    Should Be True    ${response.status_code} in [200, 207]    # 200或207（部分成功）
    
    # 驗證結果包含不同狀態
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Be Equal As Integers    ${data[0]['total_processed']}    3
    Should Be True    ${data[0]['failed_count']} > 0

Test Collect IPs Empty List
    [Documentation]    測試收集空IP列表
    [Tags]    collector    ip    batch    negative
    
    # 發送空列表
    @{empty_ips}=    Create List
    &{collect_data}=    Create Dictionary    ip_addresses=${empty_ips}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ips    ${collect_data}    expected_status=400
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    400

Test Collect IPs Too Many IPs
    [Documentation]    測試收集過多IP（超過限制）
    [Tags]    collector    ip    batch    limit    negative
    
    # 生成超過限制的IP列表（假設限制為100）
    @{too_many_ips}=    Create List
    FOR    ${i}    IN RANGE    101
        ${ip}=    Set Variable    192.168.1.${i}
        Append To List    ${too_many_ips}    ${ip}
    END
    
    &{collect_data}=    Create Dictionary    ip_addresses=${too_many_ips}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ips    ${collect_data}    expected_status=400
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    400

Test Collector Rate Limiting
    [Documentation]    測試收集器速率限制
    [Tags]    collector    rate_limit    negative
    
    # 快速發送多個請求以觸發速率限制
    ${test_ip}=    Set Variable    118.25.6.39
    &{collect_data}=    Create Dictionary    ip_address=${test_ip}
    
    # 發送多個快速請求
    FOR    ${i}    IN RANGE    10
        ${response}=    POST Request With Auth    /api/v1/collector/collect-ip    ${collect_data}    expected_status=any
        # 如果觸發速率限制，應該返回429
        Exit For Loop If    ${response.status_code} == 429
        Sleep    0.1s
    END
    
    # 檢查是否有觸發速率限制的回應
    Run Keyword If    ${response.status_code} == 429
    ...    Log    Rate limiting is working correctly
    ...    ELSE
    ...    Log    Rate limiting may not be enabled or limit is high

Test Collector Source Information
    [Documentation]    測試收集器來源資訊
    [Tags]    collector    source    positive
    
    # 收集IP資訊
    ${test_ip}=    Set Variable    118.25.6.39
    &{collect_data}=    Create Dictionary    ip_address=${test_ip}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ip    ${collect_data}
    
    # 驗證回應成功
    Verify Response Success    ${response}    200
    
    # 驗證來源資訊
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Run Keyword If    '${data[0]['threat_intelligence']}' != 'None'
    ...    Run Keywords
    ...    Should Not Be Empty    ${data[0]['threat_intelligence']['source']}
    ...    AND    Should Contain Any    ${data[0]['threat_intelligence']['source']}    abuseipdb    virustotal    otx

Test Collector Confidence Score
    [Documentation]    測試收集器信心分數
    [Tags]    collector    confidence    positive
    
    # 收集IP資訊
    ${test_ip}=    Set Variable    118.25.6.39
    &{collect_data}=    Create Dictionary    ip_address=${test_ip}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ip    ${collect_data}
    
    # 驗證回應成功
    Verify Response Success    ${response}    200
    
    # 驗證信心分數
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Run Keyword If    '${data[0]['threat_intelligence']}' != 'None'
    ...    Run Keywords
    ...    Should Be True    ${data[0]['threat_intelligence']['confidence_score']} >= 0
    ...    AND    Should Be True    ${data[0]['threat_intelligence']['confidence_score']} <= 100

Test Collector Threat Classification
    [Documentation]    測試收集器威脅分類
    [Tags]    collector    classification    positive
    
    # 收集IP資訊
    ${test_ip}=    Set Variable    118.25.6.39
    &{collect_data}=    Create Dictionary    ip_address=${test_ip}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ip    ${collect_data}
    
    # 驗證回應成功
    Verify Response Success    ${response}    200
    
    # 驗證威脅分類
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Run Keyword If    '${data[0]['threat_intelligence']}' != 'None'
    ...    Run Keywords
    ...    Should Not Be Empty    ${data[0]['threat_intelligence']['threat_type']}
    ...    AND    Should Contain Any    ${data[0]['threat_intelligence']['threat_type']}    malware    phishing    spam    botnet    scanner
    ...    AND    Should Not Be Empty    ${data[0]['threat_intelligence']['severity']}
    ...    AND    Should Contain Any    ${data[0]['threat_intelligence']['severity']}    low    medium    high    critical

Test Collector Metadata Collection
    [Documentation]    測試收集器元資料收集
    [Tags]    collector    metadata    positive
    
    # 收集IP資訊
    ${test_ip}=    Set Variable    118.25.6.39
    &{collect_data}=    Create Dictionary    ip_address=${test_ip}
    ${response}=    POST Request With Auth    /api/v1/collector/collect-ip    ${collect_data}
    
    # 驗證回應成功
    Verify Response Success    ${response}    200
    
    # 驗證元資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Run Keyword If    '${data[0]['threat_intelligence']}' != 'None'
    ...    Should Not Be Empty    ${data[0]['threat_intelligence']['metadata']}

Test Unauthorized Collector Access
    [Documentation]    測試未授權收集器存取
    [Tags]    collector    security    negative
    
    # 清除認證標頭
    &{no_auth_headers}=    Create Dictionary    Content-Type=application/json    Accept=application/json
    
    ${test_ip}=    Set Variable    118.25.6.39
    &{collect_data}=    Create Dictionary    ip_address=${test_ip}
    
    ${response}=    POST On Session    api    /api/v1/collector/collect-ip    json=${collect_data}    headers=${no_auth_headers}    expected_status=401
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    401 