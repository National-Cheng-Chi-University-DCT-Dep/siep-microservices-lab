*** Settings ***
Documentation    測試配置檔案
Library          RequestsLibrary
Library          JSONLibrary
Library          DatabaseLibrary
Library          Collections
Library          String
Library          DateTime
Library          OperatingSystem

*** Variables ***
# API 基本設定
${BASE_URL}                http://localhost:8080
${API_VERSION}             v1
${API_BASE_URL}           ${BASE_URL}/api/${API_VERSION}

# 認證相關
${AUTH_ENDPOINT}          ${API_BASE_URL}/auth
${LOGIN_ENDPOINT}         ${AUTH_ENDPOINT}/login
${REGISTER_ENDPOINT}      ${AUTH_ENDPOINT}/register
${REFRESH_ENDPOINT}       ${AUTH_ENDPOINT}/refresh
${PROFILE_ENDPOINT}       ${AUTH_ENDPOINT}/profile

# 威脅情報相關
${THREAT_ENDPOINT}        ${API_BASE_URL}/threat-intelligence
${THREAT_SEARCH_ENDPOINT} ${THREAT_ENDPOINT}/search
${THREAT_STATS_ENDPOINT}  ${THREAT_ENDPOINT}/statistics

# 收集器相關
${COLLECTOR_ENDPOINT}     ${API_BASE_URL}/collector
${COLLECT_IP_ENDPOINT}    ${COLLECTOR_ENDPOINT}/collect-ip
${COLLECT_IPS_ENDPOINT}   ${COLLECTOR_ENDPOINT}/collect-ips

# 健康檢查
${HEALTH_ENDPOINT}        ${BASE_URL}/health

# 資料庫設定
${DB_HOST}                localhost
${DB_PORT}                5432
${DB_NAME}                security_intelligence
${DB_USER}                postgres
${DB_PASSWORD}            password
${DB_CONNECTION_STRING}   postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}

# 測試用戶資料
${TEST_USER_USERNAME}     testuser
${TEST_USER_EMAIL}        testuser@example.com
${TEST_USER_PASSWORD}     testpassword123
${TEST_ADMIN_USERNAME}    admin
${TEST_ADMIN_EMAIL}       admin@example.com
${TEST_ADMIN_PASSWORD}    adminpassword123

# 測試威脅情報資料
${TEST_IP_ADDRESS}        192.168.1.100
${TEST_DOMAIN}            malicious.example.com
${TEST_URL}               http://malicious.example.com/malware
${TEST_FILE_HASH}         e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
${TEST_THREAT_TYPE}       malware
${TEST_SEVERITY}          high
${TEST_SOURCE}            test_source

# HTTP 標頭
&{DEFAULT_HEADERS}        Content-Type=application/json    Accept=application/json
&{AUTH_HEADERS}           Content-Type=application/json    Accept=application/json

# 回應時間限制（秒）
${RESPONSE_TIME_LIMIT}    5

# 分頁設定
${DEFAULT_PAGE_SIZE}      20
${MAX_PAGE_SIZE}          100

*** Keywords ***
# === 設定和清理 ===

Setup Test Environment
    [Documentation]    設定測試環境
    Create Session    api    ${BASE_URL}
    Connect To Database    psycopg2    ${DB_NAME}    ${DB_USER}    ${DB_PASSWORD}    ${DB_HOST}    ${DB_PORT}

Teardown Test Environment
    [Documentation]    清理測試環境
    Delete All Sessions
    Disconnect From Database

Setup Test Data
    [Documentation]    設定測試資料
    # 清理現有測試資料
    Cleanup Test Data
    # 創建測試用戶
    Create Test User

Cleanup Test Data
    [Documentation]    清理測試資料
    Execute Sql String    DELETE FROM threat_intelligence WHERE source = '${TEST_SOURCE}';
    Execute Sql String    DELETE FROM users WHERE username IN ('${TEST_USER_USERNAME}', '${TEST_ADMIN_USERNAME}');

Create Test User
    [Documentation]    創建測試用戶
    &{user_data}=    Create Dictionary
    ...    username=${TEST_USER_USERNAME}
    ...    email=${TEST_USER_EMAIL}
    ...    password=${TEST_USER_PASSWORD}
    ...    role=basic
    
    ${response}=    POST On Session    api    /api/v1/auth/register    json=${user_data}    headers=${DEFAULT_HEADERS}
    Should Be Equal As Strings    ${response.status_code}    201

# === 認證相關 ===

Login User
    [Arguments]    ${username}=${TEST_USER_USERNAME}    ${password}=${TEST_USER_PASSWORD}
    [Documentation]    使用者登入
    &{login_data}=    Create Dictionary    username=${username}    password=${password}
    ${response}=    POST On Session    api    /api/v1/auth/login    json=${login_data}    headers=${DEFAULT_HEADERS}
    Should Be Equal As Strings    ${response.status_code}    200
    ${token}=    Get Value From Json    ${response.json()}    $.data.access_token
    Set Suite Variable    ${ACCESS_TOKEN}    ${token[0]}
    Set To Dictionary    ${AUTH_HEADERS}    Authorization=Bearer ${ACCESS_TOKEN}
    [Return]    ${response}

Get Auth Headers
    [Documentation]    取得認證標頭
    [Return]    ${AUTH_HEADERS}

# === HTTP 請求輔助 ===

GET Request With Auth
    [Arguments]    ${endpoint}    ${params}=${EMPTY}
    [Documentation]    發送帶認證的GET請求
    ${response}=    GET On Session    api    ${endpoint}    params=${params}    headers=${AUTH_HEADERS}
    [Return]    ${response}

POST Request With Auth
    [Arguments]    ${endpoint}    ${data}=${EMPTY}
    [Documentation]    發送帶認證的POST請求
    ${response}=    POST On Session    api    ${endpoint}    json=${data}    headers=${AUTH_HEADERS}
    [Return]    ${response}

PUT Request With Auth
    [Arguments]    ${endpoint}    ${data}=${EMPTY}
    [Documentation]    發送帶認證的PUT請求
    ${response}=    PUT On Session    api    ${endpoint}    json=${data}    headers=${AUTH_HEADERS}
    [Return]    ${response}

DELETE Request With Auth
    [Arguments]    ${endpoint}
    [Documentation]    發送帶認證的DELETE請求
    ${response}=    DELETE On Session    api    ${endpoint}    headers=${AUTH_HEADERS}
    [Return]    ${response}

# === 回應驗證 ===

Verify Response Success
    [Arguments]    ${response}    ${expected_status}=200
    [Documentation]    驗證回應成功
    Should Be Equal As Strings    ${response.status_code}    ${expected_status}
    ${success}=    Get Value From Json    ${response.json()}    $.success
    Should Be True    ${success[0]}

Verify Response Error
    [Arguments]    ${response}    ${expected_status}=400
    [Documentation]    驗證回應錯誤
    Should Be Equal As Strings    ${response.status_code}    ${expected_status}
    ${success}=    Get Value From Json    ${response.json()}    $.success
    Should Be Equal    ${success[0]}    ${False}

Verify Response Time
    [Arguments]    ${response}    ${max_time}=${RESPONSE_TIME_LIMIT}
    [Documentation]    驗證回應時間
    Should Be True    ${response.elapsed.total_seconds()} < ${max_time}

# === 資料驗證 ===

Generate Test Threat Data
    [Documentation]    生成測試威脅情報資料
    &{threat_data}=    Create Dictionary
    ...    ip_address=${TEST_IP_ADDRESS}
    ...    domain=${TEST_DOMAIN}
    ...    url=${TEST_URL}
    ...    file_hash=${TEST_FILE_HASH}
    ...    threat_type=${TEST_THREAT_TYPE}
    ...    severity=${TEST_SEVERITY}
    ...    description=Test threat intelligence data
    ...    tags=@{['test', 'automation']}
    ...    confidence_score=${85}
    ...    source=${TEST_SOURCE}
    [Return]    ${threat_data}

Verify Threat Intelligence Data
    [Arguments]    ${threat_data}    ${expected_data}
    [Documentation]    驗證威脅情報資料
    Should Be Equal    ${threat_data['ip_address']}    ${expected_data['ip_address']}
    Should Be Equal    ${threat_data['domain']}    ${expected_data['domain']}
    Should Be Equal    ${threat_data['threat_type']}    ${expected_data['threat_type']}
    Should Be Equal    ${threat_data['severity']}    ${expected_data['severity']}

# === 資料庫驗證 ===

Verify Database Record Exists
    [Arguments]    ${table}    ${condition}
    [Documentation]    驗證資料庫記錄存在
    ${count}=    Row Count    SELECT * FROM ${table} WHERE ${condition}
    Should Be True    ${count} > 0

Verify Database Record Not Exists
    [Arguments]    ${table}    ${condition}
    [Documentation]    驗證資料庫記錄不存在
    ${count}=    Row Count    SELECT * FROM ${table} WHERE ${condition}
    Should Be Equal As Integers    ${count}    0

Get Database Record
    [Arguments]    ${query}
    [Documentation]    取得資料庫記錄
    ${result}=    Query    ${query}
    [Return]    ${result}

# === 工具函數 ===

Generate Random String
    [Arguments]    ${length}=10
    [Documentation]    生成隨機字串
    ${random_string}=    Generate Random String    ${length}    [LETTERS]
    [Return]    ${random_string}

Generate Random Email
    [Documentation]    生成隨機郵箱
    ${random_part}=    Generate Random String    8
    ${email}=    Set Variable    ${random_part}@example.com
    [Return]    ${email}

Get Current Timestamp
    [Documentation]    取得當前時間戳
    ${timestamp}=    Get Current Date    result_format=%Y-%m-%d %H:%M:%S
    [Return]    ${timestamp}

Wait For Server Ready
    [Arguments]    ${timeout}=30
    [Documentation]    等待伺服器就緒
    FOR    ${i}    IN RANGE    ${timeout}
        ${status}    ${response}=    Run Keyword And Ignore Error    GET On Session    api    /health
        Run Keyword If    '${status}' == 'PASS' and ${response.status_code} == 200    RETURN
        Sleep    1s
    END
    Fail    Server not ready after ${timeout} seconds 