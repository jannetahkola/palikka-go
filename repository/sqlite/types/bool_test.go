package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type targetType struct {
	Foo   string `json:"foo"`
	Valid Bool   `json:"valid"`
}

// todo test this
type targetEmbeddedType struct {
	Bar    string     `json:"bar"`
	Target targetType `json:"target"`
}

func Test_BoolFrom(t *testing.T) {
	assert.Equal(t, int64(1), BoolFrom(true).Int64)
	assert.Equal(t, int64(0), BoolFrom(false).Int64)
}

func TestBool_Value(t *testing.T) {
	bt := &Bool{Int64: 1}
	vt, _ := bt.Value()
	assert.True(t, vt.(bool))

	bf := &Bool{Int64: 0}
	vf, _ := bf.Value()
	assert.False(t, vf.(bool))
}

func TestBool_Scan_From_Int64(t *testing.T) {
	var bt Bool
	err := bt.Scan(int64(1))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), bt.Int64)

	var bf Bool
	err = bf.Scan(int64(0))
	assert.NoError(t, err)
	assert.Equal(t, int64(0), bf.Int64)
}

func TestBool_Scan_From_Int32(t *testing.T) {
	var b Bool
	err := b.Scan(int32(1))
	assert.Error(t, err)
}

func TestBool_Scan_From_String(t *testing.T) {
	var b Bool
	err := b.Scan("int64(1)")
	assert.Error(t, err)
}

func TestBool_UnmarshalJSON_To_Boolean_True(t *testing.T) {
	jsonString := `{"foo": "bar", "valid": true}`
	var s targetType
	if err := json.Unmarshal([]byte(jsonString), &s); err != nil {
		t.Fatal(err)
	}

	v, _ := s.Valid.Value()
	assert.Equal(t, "bar", s.Foo)
	assert.Equal(t, true, v.(bool))
}

func TestBool_UnmarshalJSON_To_Boolean_False(t *testing.T) {
	jsonStrings := []string{
		`{"foo": "bar", "valid": false}`,
		`{"foo": "bar"}`,
	}
	for _, jsonString := range jsonStrings {
		var s targetType
		if err := json.Unmarshal([]byte(jsonString), &s); err != nil {
			t.Fatal(err)
		}

		v, _ := s.Valid.Value()
		assert.Equal(t, "bar", s.Foo)
		assert.Equal(t, false, v.(bool))
	}
}

func TestBool_Equal(t *testing.T) {
	assert.True(t, BoolFrom(true).Equal(true))
	assert.False(t, BoolFrom(false).Equal(true))
}

// todo test marshalling
