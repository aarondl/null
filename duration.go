package null

import (
	"bytes"
	"database/sql/driver"
	"time"
	"encoding/json"
	"strconv"
	"github.com/llaforge72/null/convert"
)

// Duration is a nullable time.Duration. It supports SQL and JSON serialization.
type Duration struct {
	Duration  time.Duration
	Valid bool
}

// NewDuration creates a new Duration.
func NewDuration(d time.Duration, valid bool) Duration {
	return Duration{
		Duration:  d,
		Valid: valid,
	}
}

// DurationFrom creates a new Duration that will always be valid.
func DurationFrom(d time.Duration) Duration {
	return NewDuration(d, true)
}

// DurationFromPtr creates a new Duration that will be null if t is nil.
func DurationFromPtr(d *time.Duration) Duration {
	if d == nil {
		return NewDuration(time.Minute * 0, false)
	}
	return NewDuration(*d, true)
}

// MarshalJSON implements json.Marshaler.
func (d Duration) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatInt(int64(d.Duration), 10)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Duration) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		d.Valid = false
		d.Duration = time.Minute * 0
		return nil
	}

	if err := json.Unmarshal(data, &d.Duration); err != nil {
		return err
	}

	d.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (d Duration) MarshalText() ([]byte, error) {
	if !d.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(d.Duration), 10)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Duration) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		d.Valid = false
		return nil
	}
	value, err := strconv.ParseInt(string(text), 10, 64)
	d.Valid = err == nil
	d.Duration = time.Duration(value)
	return err
}


// SetValid changes this Duration's value and sets it to be non-null.
func (d *Duration) SetValid(v time.Duration) {
	d.Duration = v
	d.Valid = true
}

// Ptr returns a pointer to this Duration's value, or a nil pointer if this Time is null.
func (d Duration) Ptr() *time.Duration {
	if !d.Valid {
		return nil
	}
	return &d.Duration
}

// IsZero returns true for invalid Duration's, for future omitempty support (Go 1.4?)
func (d Duration) IsZero() bool {
	return !d.Valid
}

// Scan implements the Scanner interface.
func (d *Duration) Scan(value interface{}) error {
	if value == nil {
		d.Duration, d.Valid = 0, false
		return nil
	}
	d.Valid = true
	return convert.ConvertAssign(&d.Duration, value)
}

// Value implements the driver Valuer interface.
func (d Duration) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Duration, nil
}
