package dto

import "errors"

// 通用錯誤
var (
	ErrInvalidIPAddress     = errors.New("invalid IP address")
	ErrInvalidDomain        = errors.New("invalid domain")
	ErrInvalidThreatType    = errors.New("invalid threat type")
	ErrInvalidSeverity      = errors.New("invalid severity level")
	ErrInvalidConfidence    = errors.New("invalid confidence score")
	ErrInvalidCountryCode   = errors.New("invalid country code")
	ErrInvalidASN           = errors.New("invalid ASN")
	ErrInvalidPagination    = errors.New("invalid pagination parameters")
	ErrInvalidSortBy        = errors.New("invalid sort by parameter")
	ErrInvalidSortOrder     = errors.New("invalid sort order")
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrInvalidUUID          = errors.New("invalid UUID")
	ErrInvalidSource        = errors.New("invalid source")
	ErrInvalidDescription   = errors.New("invalid description")
	ErrInvalidISP           = errors.New("invalid ISP")
	ErrInvalidTags          = errors.New("invalid tags")
	ErrInvalidMetadata      = errors.New("invalid metadata")
	ErrInvalidBulkOperation = errors.New("invalid bulk operation")
	ErrInvalidRole          = errors.New("invalid user role")
	ErrInvalidAPIKey        = errors.New("invalid API key")
	ErrInvalidUsername      = errors.New("invalid username")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidQuota         = errors.New("invalid quota")
	ErrInvalidUsage         = errors.New("invalid usage")
	ErrInvalidExpiration    = errors.New("invalid expiration")
	ErrInvalidCollectionInterval = errors.New("invalid collection interval")
	ErrInvalidJobStatus     = errors.New("invalid job status")
	ErrInvalidSourceConfig  = errors.New("invalid source configuration")
	
	// 認證相關錯誤
	ErrUsernameExists       = errors.New("username already exists")
	ErrEmailExists          = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserInactive         = errors.New("user is inactive")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidCurrentPassword = errors.New("invalid current password")
	
	// 威脅情報相關錯誤
	ErrThreatNotFound       = errors.New("threat not found")
) 