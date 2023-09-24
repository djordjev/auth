package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB map[string]any

func (jsonField JSONB) Value() (driver.Value, error) {
	return json.Marshal(jsonField)
}

// Scan Unmarshal
func (jsonField *JSONB) Scan(value any) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("serialize jsonb: type assertion to []byte failed")
	}

	return json.Unmarshal(data, &jsonField)
}
