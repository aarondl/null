package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/volatiletech/null/convert"
)

// Float32 is a nullable float32.
type Float32 struct {
	Float32 float32
	Valid   bool
}

// NewFloat32 creates a new Float32
func NewFloat32(f float32, valid bool) Float32 {
	return Float32{
		Float32: f,
		Valid:   valid,
	}
}

// Float32From creates a new Float32 that will always be valid.
func Float32From(f float32) Float32 {
	return NewFloat32(f, true)
}

// Float32FromPtr creates a new Float32 that be null if f is nil.
func Float32FromPtr(f *float32) Float32 {
	if f == nil {
		return NewFloat32(0, false)
	}
	return NewFloat32(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Float32) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		f.Valid = false
		f.Float32 = 0
		return nil
	}

	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	var r float64
	switch x := v.(type) {
	case float64:
		r = v.(float64)
	case string:
		str := string(x)
		if len(str) == 0 {
			f.Valid = false
			return nil
		}

		r, err = strconv.ParseFloat(str, 32)
	case nil:
		f.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int32", reflect.TypeOf(v).Name())
	}

	if r > math.MaxFloat32 {
		return fmt.Errorf("json: %f overflows max float32 value", r)
	}

	f.Float32 = float32(r)
	f.Valid = (err == nil) && (f.Float32 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (f *Float32) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		f.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseFloat(string(text), 32)
	f.Valid = err == nil
	if f.Valid {
		f.Float32 = float32(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (f Float32) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (f Float32) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// SetValid changes this Float32's value and also sets it to be non-null.
func (f *Float32) SetValid(n float32) {
	f.Float32 = n
	f.Valid = true
}

// Ptr returns a pointer to this Float32's value, or a nil pointer if this Float32 is null.
func (f Float32) Ptr() *float32 {
	if !f.Valid {
		return nil
	}
	return &f.Float32
}

// IsZero returns true for invalid Float32s, for future omitempty support (Go 1.4?)
func (f Float32) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (f *Float32) Scan(value interface{}) error {
	if value == nil {
		f.Float32, f.Valid = 0, false
		return nil
	}
	f.Valid = true
	return convert.ConvertAssign(&f.Float32, value)
}

// Value implements the driver Valuer interface.
func (f Float32) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return float64(f.Float32), nil
}

// Randomize for sqlboiler
func (f *Float32) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		f.Float32 = 0
		f.Valid = false
	} else {
		f.Float32 = float32(nextInt()%10)/10.0 + float32(nextInt()%10)
		f.Valid = true
	}
}
