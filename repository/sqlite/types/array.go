package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// todo copy structure to bool.go
type Array struct {
	driver.Valuer
	sql.Scanner
	Arr []interface{}
}

func (a Array) Value() (driver.Value, error) {
	return a.Arr, nil
}

func (a *Array) Scan(src interface{}) error {
	str := src.(string)
	err := json.Unmarshal([]byte(str), &a.Arr)
	if err != nil {
		return err
	}
	return nil
}

func (a Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Arr)
}

func (a *Array) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &a.Arr)
	if err != nil {
		return err
	}
	return nil
}
