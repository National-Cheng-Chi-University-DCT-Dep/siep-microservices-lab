*** Settings ***
Documentation    威脅情報API測試套件
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

Test Create Threat Intelligence Success
    [Documentation]    測試建立威脅情報成功
    [Tags]    threat    create    positive
    
    # 準備測試資料
    ${threat_data}=    Generate Test Threat Data
    
    # 發送建立請求
    ${response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
    
    # 驗證回應
    Verify Response Success    ${response}    201
    Verify Response Time    ${response}
    
    # 驗證回應資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Not Be Empty    ${data[0]['id']}
    Verify Threat Intelligence Data    ${data[0]}    ${threat_data}
    
    # 驗證資料庫記錄
    ${threat_id}=    Set Variable    ${data[0]['id']}
    Verify Database Record Exists    threat_intelligence    id='${threat_id}'

Test Create Threat Intelligence Invalid Data
    [Documentation]    測試建立威脅情報無效資料
    [Tags]    threat    create    negative
    
    # 準備無效資料（缺少必填欄位）
    &{invalid_data}=    Create Dictionary
    ...    description=Invalid threat data
    
    ${response}=    POST Request With Auth    /api/v1/threat-intelligence    ${invalid_data}    expected_status=400
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    400

Test Get Threat Intelligence Success
    [Documentation]    測試取得威脅情報成功
    [Tags]    threat    get    positive
    
    # 先建立威脅情報
    ${threat_data}=    Generate Test Threat Data
    ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
    Should Be Equal As Strings    ${create_response.status_code}    201
    
    ${threat_id}=    Get Value From Json    ${create_response.json()}    $.data.id
    
    # 取得威脅情報
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence/${threat_id[0]}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Be Equal    ${data[0]['id']}    ${threat_id[0]}
    Verify Threat Intelligence Data    ${data[0]}    ${threat_data}

Test Get Threat Intelligence Not Found
    [Documentation]    測試取得不存在的威脅情報
    [Tags]    threat    get    negative
    
    ${fake_id}=    Set Variable    123e4567-e89b-12d3-a456-426614174000
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence/${fake_id}    expected_status=404
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    404

Test Update Threat Intelligence Success
    [Documentation]    測試更新威脅情報成功
    [Tags]    threat    update    positive
    
    # 先建立威脅情報
    ${threat_data}=    Generate Test Threat Data
    ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
    Should Be Equal As Strings    ${create_response.status_code}    201
    
    ${threat_id}=    Get Value From Json    ${create_response.json()}    $.data.id
    
    # 更新威脅情報
    &{update_data}=    Create Dictionary
    ...    severity=critical
    ...    description=Updated threat description
    ...    confidence_score=${95}
    
    ${response}=    PUT Request With Auth    /api/v1/threat-intelligence/${threat_id[0]}    ${update_data}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證更新後的資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Be Equal    ${data[0]['severity']}    critical
    Should Be Equal    ${data[0]['description']}    Updated threat description
    Should Be Equal As Integers    ${data[0]['confidence_score']}    95

Test Delete Threat Intelligence Success
    [Documentation]    測試刪除威脅情報成功
    [Tags]    threat    delete    positive
    
    # 先建立威脅情報
    ${threat_data}=    Generate Test Threat Data
    ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
    Should Be Equal As Strings    ${create_response.status_code}    201
    
    ${threat_id}=    Get Value From Json    ${create_response.json()}    $.data.id
    
    # 刪除威脅情報
    ${response}=    DELETE Request With Auth    /api/v1/threat-intelligence/${threat_id[0]}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證資料庫記錄已刪除
    Verify Database Record Not Exists    threat_intelligence    id='${threat_id[0]}' AND deleted_at IS NULL

Test List Threat Intelligence Success
    [Documentation]    測試列出威脅情報成功
    [Tags]    threat    list    positive
    
    # 建立多筆威脅情報
    FOR    ${i}    IN RANGE    3
        ${threat_data}=    Generate Test Threat Data
        Set To Dictionary    ${threat_data}    ip_address    192.168.1.${i}
        ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
        Should Be Equal As Strings    ${create_response.status_code}    201
    END
    
    # 列出威脅情報
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證分頁資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    ${threats}=    Get Value From Json    ${response.json()}    $.data.threats
    ${pagination}=    Get Value From Json    ${response.json()}    $.data.pagination
    
    Should Not Be Empty    ${threats[0]}
    Should Be True    ${pagination[0]['total_records']} >= 3

Test List Threat Intelligence With Filters
    [Documentation]    測試使用篩選器列出威脅情報
    [Tags]    threat    list    filter    positive
    
    # 建立不同類型的威脅情報
    ${malware_data}=    Generate Test Threat Data
    Set To Dictionary    ${malware_data}    threat_type    malware    severity    high
    ${phishing_data}=    Generate Test Threat Data
    Set To Dictionary    ${phishing_data}    threat_type    phishing    severity    medium    ip_address    192.168.1.101
    
    ${create_response1}=    POST Request With Auth    /api/v1/threat-intelligence    ${malware_data}
    ${create_response2}=    POST Request With Auth    /api/v1/threat-intelligence    ${phishing_data}
    Should Be Equal As Strings    ${create_response1.status_code}    201
    Should Be Equal As Strings    ${create_response2.status_code}    201
    
    # 測試威脅類型篩選
    &{params}=    Create Dictionary    threat_type=malware
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence    ${params}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    ${threats}=    Get Value From Json    ${response.json()}    $.data.threats
    
    # 驗證所有結果都是malware類型
    FOR    ${threat}    IN    @{threats[0]}
        Should Be Equal    ${threat['threat_type']}    malware
    END

Test List Threat Intelligence With Pagination
    [Documentation]    測試分頁列出威脅情報
    [Tags]    threat    list    pagination    positive
    
    # 建立多筆威脅情報
    FOR    ${i}    IN RANGE    5
        ${threat_data}=    Generate Test Threat Data
        Set To Dictionary    ${threat_data}    ip_address    192.168.2.${i}
        ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
        Should Be Equal As Strings    ${create_response.status_code}    201
    END
    
    # 測試第一頁
    &{params}=    Create Dictionary    page=1    page_size=2
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence    ${params}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    ${threats}=    Get Value From Json    ${response.json()}    $.data.threats
    ${pagination}=    Get Value From Json    ${response.json()}    $.data.pagination
    
    Should Be True    len(${threats[0]}) <= 2
    Should Be Equal As Integers    ${pagination[0]['current_page']}    1
    Should Be Equal As Integers    ${pagination[0]['page_size']}    2

Test Search Threat Intelligence Success
    [Documentation]    測試搜尋威脅情報成功
    [Tags]    threat    search    positive
    
    # 建立威脅情報
    ${threat_data}=    Generate Test Threat Data
    Set To Dictionary    ${threat_data}    description    Unique test malware description
    ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
    Should Be Equal As Strings    ${create_response.status_code}    201
    
    # 搜尋威脅情報
    &{search_params}=    Create Dictionary    query=Unique test malware
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence/search    ${search_params}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證搜尋結果
    ${threats}=    Get Value From Json    ${response.json()}    $.data.threats
    Should Not Be Empty    ${threats[0]}
    Should Contain    ${threats[0][0]['description']}    Unique test malware

Test Search By IP Address
    [Documentation]    測試按IP地址搜尋威脅情報
    [Tags]    threat    search    ip    positive
    
    # 建立威脅情報
    ${test_ip}=    Set Variable    10.0.0.100
    ${threat_data}=    Generate Test Threat Data
    Set To Dictionary    ${threat_data}    ip_address    ${test_ip}
    ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
    Should Be Equal As Strings    ${create_response.status_code}    201
    
    # 按IP搜尋
    &{search_params}=    Create Dictionary    ip_address=${test_ip}
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence/search    ${search_params}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    ${threats}=    Get Value From Json    ${response.json()}    $.data.threats
    Should Not Be Empty    ${threats[0]}
    Should Be Equal    ${threats[0][0]['ip_address']}    ${test_ip}

Test Search By Domain
    [Documentation]    測試按域名搜尋威脅情報
    [Tags]    threat    search    domain    positive
    
    # 建立威脅情報
    ${test_domain}=    Set Variable    evil.test.com
    ${threat_data}=    Generate Test Threat Data
    Set To Dictionary    ${threat_data}    domain    ${test_domain}
    ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
    Should Be Equal As Strings    ${create_response.status_code}    201
    
    # 按域名搜尋
    &{search_params}=    Create Dictionary    domain=${test_domain}
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence/search    ${search_params}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    ${threats}=    Get Value From Json    ${response.json()}    $.data.threats
    Should Not Be Empty    ${threats[0]}
    Should Be Equal    ${threats[0][0]['domain']}    ${test_domain}

Test Get Threat Intelligence Statistics
    [Documentation]    測試取得威脅情報統計
    [Tags]    threat    statistics    positive
    
    # 建立不同類型的威脅情報
    FOR    ${i}    IN RANGE    3
        ${threat_data}=    Generate Test Threat Data
        Set To Dictionary    ${threat_data}    threat_type    malware    ip_address    192.168.3.${i}
        ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
        Should Be Equal As Strings    ${create_response.status_code}    201
    END
    
    FOR    ${i}    IN RANGE    2
        ${threat_data}=    Generate Test Threat Data
        Set To Dictionary    ${threat_data}    threat_type    phishing    ip_address    192.168.4.${i}
        ${create_response}=    POST Request With Auth    /api/v1/threat-intelligence    ${threat_data}
        Should Be Equal As Strings    ${create_response.status_code}    201
    END
    
    # 取得統計資料
    ${response}=    GET Request With Auth    /api/v1/threat-intelligence/statistics
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證統計資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Be True    ${data[0]['total_threats']} >= 5
    Should Not Be Empty    ${data[0]['threat_type_stats']}
    Should Not Be Empty    ${data[0]['source_stats']}

Test Batch Create Threat Intelligence
    [Documentation]    測試批量建立威脅情報
    [Tags]    threat    batch    create    positive
    
    # 準備批量資料
    @{threats_list}=    Create List
    FOR    ${i}    IN RANGE    3
        ${threat_data}=    Generate Test Threat Data
        Set To Dictionary    ${threat_data}    ip_address    192.168.5.${i}
        Append To List    ${threats_list}    ${threat_data}
    END
    
    &{batch_data}=    Create Dictionary    threats=${threats_list}
    
    # 發送批量建立請求
    ${response}=    POST Request With Auth    /api/v1/threat-intelligence/batch    ${batch_data}
    
    # 驗證回應
    Verify Response Success    ${response}    201
    Verify Response Time    ${response}
    
    # 驗證批量建立結果
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Be Equal As Integers    ${data[0]['total_created']}    3
    Should Be Equal As Integers    ${data[0]['failed_count']}    0

Test Unauthorized Access
    [Documentation]    測試未授權存取
    [Tags]    threat    security    negative
    
    # 清除認證標頭
    &{no_auth_headers}=    Create Dictionary    Content-Type=application/json    Accept=application/json
    
    ${response}=    GET On Session    api    /api/v1/threat-intelligence    headers=${no_auth_headers}    expected_status=401
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    401 