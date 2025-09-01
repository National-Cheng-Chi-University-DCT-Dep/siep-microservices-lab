*** Settings ***
Documentation    認證API測試套件
Resource         ../config/test_config.robot
Suite Setup      Setup Test Environment
Suite Teardown   Teardown Test Environment
Test Setup       Setup Test Data
Test Teardown    Cleanup Test Data

*** Test Cases ***

Test User Registration Success
    [Documentation]    測試使用者註冊成功
    [Tags]    auth    register    positive
    
    # 準備測試資料
    ${random_email}=    Generate Random Email
    ${random_username}=    Generate Random String    12
    &{register_data}=    Create Dictionary
    ...    username=${random_username}
    ...    email=${random_email}
    ...    password=${TEST_USER_PASSWORD}
    ...    role=basic
    
    # 發送註冊請求
    ${response}=    POST On Session    api    /api/v1/auth/register    json=${register_data}    headers=${DEFAULT_HEADERS}
    
    # 驗證回應
    Verify Response Success    ${response}    201
    Verify Response Time    ${response}
    
    # 驗證回應資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Not Be Empty    ${data[0]['access_token']}
    Should Not Be Empty    ${data[0]['refresh_token']}
    Should Be Equal    ${data[0]['token_type']}    Bearer
    Should Be Equal    ${data[0]['user']['username']}    ${random_username}
    Should Be Equal    ${data[0]['user']['email']}    ${random_email}
    
    # 驗證資料庫記錄
    Verify Database Record Exists    users    username='${random_username}'

Test User Registration Duplicate Username
    [Documentation]    測試使用者註冊重複使用者名稱
    [Tags]    auth    register    negative
    
    # 先註冊一個使用者
    ${random_email1}=    Generate Random Email
    ${random_email2}=    Generate Random Email
    ${username}=    Generate Random String    12
    
    &{register_data1}=    Create Dictionary
    ...    username=${username}
    ...    email=${random_email1}
    ...    password=${TEST_USER_PASSWORD}
    
    ${response1}=    POST On Session    api    /api/v1/auth/register    json=${register_data1}    headers=${DEFAULT_HEADERS}
    Should Be Equal As Strings    ${response1.status_code}    201
    
    # 嘗試使用相同使用者名稱註冊
    &{register_data2}=    Create Dictionary
    ...    username=${username}
    ...    email=${random_email2}
    ...    password=${TEST_USER_PASSWORD}
    
    ${response2}=    POST On Session    api    /api/v1/auth/register    json=${register_data2}    headers=${DEFAULT_HEADERS}    expected_status=409
    
    # 驗證錯誤回應
    Verify Response Error    ${response2}    409
    ${error}=    Get Value From Json    ${response2.json()}    $.error
    Should Contain    ${error[0]['code']}    USERNAME_EXISTS

Test User Registration Invalid Email
    [Documentation]    測試使用者註冊無效郵箱格式
    [Tags]    auth    register    negative
    
    &{register_data}=    Create Dictionary
    ...    username=testuser123
    ...    email=invalid-email
    ...    password=${TEST_USER_PASSWORD}
    
    ${response}=    POST On Session    api    /api/v1/auth/register    json=${register_data}    headers=${DEFAULT_HEADERS}    expected_status=400
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    400

Test User Login Success
    [Documentation]    測試使用者登入成功
    [Tags]    auth    login    positive
    
    # 先註冊使用者
    ${random_email}=    Generate Random Email
    ${random_username}=    Generate Random String    12
    &{register_data}=    Create Dictionary
    ...    username=${random_username}
    ...    email=${random_email}
    ...    password=${TEST_USER_PASSWORD}
    
    ${register_response}=    POST On Session    api    /api/v1/auth/register    json=${register_data}    headers=${DEFAULT_HEADERS}
    Should Be Equal As Strings    ${register_response.status_code}    201
    
    # 測試登入
    &{login_data}=    Create Dictionary
    ...    username=${random_username}
    ...    password=${TEST_USER_PASSWORD}
    
    ${response}=    POST On Session    api    /api/v1/auth/login    json=${login_data}    headers=${DEFAULT_HEADERS}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證回應資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Not Be Empty    ${data[0]['access_token']}
    Should Not Be Empty    ${data[0]['refresh_token']}
    Should Be Equal    ${data[0]['token_type']}    Bearer
    Should Be Equal    ${data[0]['user']['username']}    ${random_username}

Test User Login With Email
    [Documentation]    測試使用郵箱登入
    [Tags]    auth    login    positive
    
    # 先註冊使用者
    ${random_email}=    Generate Random Email
    ${random_username}=    Generate Random String    12
    &{register_data}=    Create Dictionary
    ...    username=${random_username}
    ...    email=${random_email}
    ...    password=${TEST_USER_PASSWORD}
    
    ${register_response}=    POST On Session    api    /api/v1/auth/register    json=${register_data}    headers=${DEFAULT_HEADERS}
    Should Be Equal As Strings    ${register_response.status_code}    201
    
    # 使用郵箱登入
    &{login_data}=    Create Dictionary
    ...    username=${random_email}
    ...    password=${TEST_USER_PASSWORD}
    
    ${response}=    POST On Session    api    /api/v1/auth/login    json=${login_data}    headers=${DEFAULT_HEADERS}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Should Be Equal    ${response.json()['data']['user']['email']}    ${random_email}

Test User Login Invalid Credentials
    [Documentation]    測試使用者登入無效認證
    [Tags]    auth    login    negative
    
    &{login_data}=    Create Dictionary
    ...    username=nonexistentuser
    ...    password=wrongpassword
    
    ${response}=    POST On Session    api    /api/v1/auth/login    json=${login_data}    headers=${DEFAULT_HEADERS}    expected_status=401
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    401
    ${error}=    Get Value From Json    ${response.json()}    $.error
    Should Contain    ${error[0]['code']}    INVALID_CREDENTIALS

Test Token Refresh Success
    [Documentation]    測試令牌刷新成功
    [Tags]    auth    refresh    positive
    
    # 先登入取得令牌
    ${login_response}=    Login User
    ${refresh_token}=    Get Value From Json    ${login_response.json()}    $.data.refresh_token
    
    # 刷新令牌
    &{refresh_data}=    Create Dictionary    refresh_token=${refresh_token[0]}
    ${response}=    POST On Session    api    /api/v1/auth/refresh    json=${refresh_data}    headers=${DEFAULT_HEADERS}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證新令牌
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Not Be Empty    ${data[0]['access_token']}
    Should Not Be Empty    ${data[0]['refresh_token']}
    Should Be Equal    ${data[0]['token_type']}    Bearer

Test Token Refresh Invalid Token
    [Documentation]    測試令牌刷新無效令牌
    [Tags]    auth    refresh    negative
    
    &{refresh_data}=    Create Dictionary    refresh_token=invalid_token
    ${response}=    POST On Session    api    /api/v1/auth/refresh    json=${refresh_data}    headers=${DEFAULT_HEADERS}    expected_status=401
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    401

Test Get User Profile Success
    [Documentation]    測試取得使用者檔案成功
    [Tags]    auth    profile    positive
    
    # 登入使用者
    ${login_response}=    Login User
    
    # 取得使用者檔案
    ${response}=    GET Request With Auth    /api/v1/auth/profile
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證使用者資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Be Equal    ${data[0]['username']}    ${TEST_USER_USERNAME}
    Should Be Equal    ${data[0]['email']}    ${TEST_USER_EMAIL}

Test Get User Profile Unauthorized
    [Documentation]    測試未認證取得使用者檔案
    [Tags]    auth    profile    negative
    
    ${response}=    GET On Session    api    /api/v1/auth/profile    headers=${DEFAULT_HEADERS}    expected_status=401
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    401

Test Update User Profile Success
    [Documentation]    測試更新使用者檔案成功
    [Tags]    auth    profile    positive
    
    # 登入使用者
    ${login_response}=    Login User
    
    # 更新使用者檔案
    ${new_username}=    Generate Random String    12
    &{update_data}=    Create Dictionary    username=${new_username}
    
    ${response}=    PUT Request With Auth    /api/v1/auth/profile    ${update_data}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 驗證更新後的資料
    ${data}=    Get Value From Json    ${response.json()}    $.data
    Should Be Equal    ${data[0]['username']}    ${new_username}
    
    # 驗證資料庫記錄
    Verify Database Record Exists    users    username='${new_username}'

Test Change Password Success
    [Documentation]    測試修改密碼成功
    [Tags]    auth    password    positive
    
    # 登入使用者
    ${login_response}=    Login User
    
    # 修改密碼
    ${new_password}=    Set Variable    newpassword123
    &{password_data}=    Create Dictionary
    ...    current_password=${TEST_USER_PASSWORD}
    ...    new_password=${new_password}
    
    ${response}=    POST Request With Auth    /api/v1/auth/change-password    ${password_data}
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response}
    
    # 測試使用新密碼登入
    &{login_data}=    Create Dictionary
    ...    username=${TEST_USER_USERNAME}
    ...    password=${new_password}
    
    ${login_response}=    POST On Session    api    /api/v1/auth/login    json=${login_data}    headers=${DEFAULT_HEADERS}
    Should Be Equal As Strings    ${login_response.status_code}    200

Test Change Password Wrong Current Password
    [Documentation]    測試修改密碼錯誤的當前密碼
    [Tags]    auth    password    negative
    
    # 登入使用者
    ${login_response}=    Login User
    
    # 使用錯誤的當前密碼
    &{password_data}=    Create Dictionary
    ...    current_password=wrongpassword
    ...    new_password=newpassword123
    
    ${response}=    POST Request With Auth    /api/v1/auth/change-password    ${password_data}    expected_status=401
    
    # 驗證錯誤回應
    Verify Response Error    ${response}    401

Test Logout Success
    [Documentation]    測試登出成功
    [Tags]    auth    logout    positive
    
    # 登入使用者
    ${login_response}=    Login User
    
    # 登出
    ${response}=    POST Request With Auth    /api/v1/auth/logout
    
    # 驗證回應
    Verify Response Success    ${response}    200
    Verify Response Time    ${response} 