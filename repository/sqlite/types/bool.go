package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

var BoolPtr *bool = nil

// Bool is a bool representation of an int64 value.
//
// Can replace int64 as data type for models where SQLite's int
// is mapped to int64 by sqlc. Useful when representing booleans
// as 0s and 1s in non-null int columns in SQLite, and using
// sqlc-generated models as domain types.
// todo see https://pkg.go.dev/database/sql/driver#Valuer
type Bool struct {
	Int64 int64
}

func BoolFrom(v bool) Bool {
	if v {
		return Bool{Int64: 1}
	}
	return Bool{Int64: 0}
}

// Value implements database/sql/driver/Valuer.
//
// Error is always nil.
func (p Bool) Value() (driver.Value, error) {
	fmt.Println("Value()", p.Int64)
	return p.Int64 > 0, nil
}

func (p *Bool) Scan(src interface{}) error {
	v, ok := src.(int64)
	if !ok {
		return fmt.Errorf("expected int64, got %T", src)
	}
	*p = Bool{Int64: v}
	return nil
}

func (p Bool) MarshalJSON() ([]byte, error) {
	fmt.Println("MarshalJSON()")
	v, _ := p.Value()
	return []byte(strconv.FormatBool(v.(bool))), nil
}

func (p *Bool) UnmarshalJSON(data []byte) error {
	v := string(data)
	switch v {
	case "", "nil", "false":
		p.Int64 = 0
	case "true":
		p.Int64 = 1
	default:
		return fmt.Errorf("unknown boolean value %q", v)
	}
	return nil
}

func (p Bool) Equal(other bool) bool {
	v, _ := p.Value()
	return v == other
}
