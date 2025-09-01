package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB 自定義 JSONB 類型，用於儲存結構化資料
type JSONB map[string]interface{}

// Value 實現 driver.Valuer 接口
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 實現 sql.Scanner 接口
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(data, j)
}
