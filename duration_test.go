package null

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	durationJSON = []byte(`9203685806`)
)

func TestDurationFrom(t *testing.T) {
	i := DurationFrom(9203685806)
	assertDuration(t, i, "DurationFrom()")

	zero := DurationFrom(0)
	if !zero.Valid {
		t.Error("DurationFrom(0)", "is invalid, but should be valid")
	}
}

func TestDurationFromPtr(t *testing.T) {
	n := time.Duration(9203685806)
	iptr := &n
	i := DurationFromPtr(iptr)
	assertDuration(t, i, "DurationFromPtr()")

	null := DurationFromPtr(nil)
	assertNullDuration(t, null, "DurationFromPtr(nil)")
}

func TestUnmarshalDuration(t *testing.T) {
	var duration Duration
	err := json.Unmarshal(durationJSON, &duration)
	maybePanic(err)
	assertDuration(t, duration, "Duration json")

	var null Duration
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullDuration(t, null, "null json")

	var badType Duration
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullDuration(t, badType, "wrong type json")

	var invalid Duration
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullDuration(t, invalid, "invalid json")
}

func TestTextUnmarshalDuration(t *testing.T) {
	var i Duration
	err := i.UnmarshalText([]byte("9203685806"))
	maybePanic(err)
	assertDuration(t, i, "UnmarshalText() int64")

	var blank Duration
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullDuration(t, blank, "UnmarshalText() empty int64")
}

func TestMarshalDuration(t *testing.T) {
	i := DurationFrom(9203685806)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "9203685806", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewDuration(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalDurationText(t *testing.T) {
	i := DurationFrom(9203685806)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "9203685806", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewDuration(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestDurationPointer(t *testing.T) {
	i := DurationFrom(9203685806)
	ptr := i.Ptr()
	if *ptr != 9203685806 {
		t.Errorf("bad %s int64: %#v ≠ %d\n", "pointer", ptr, 9203685806)
	}

	null := NewDuration(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestDurationIsZero(t *testing.T) {
	i := DurationFrom(9203685806)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewDuration(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewDuration(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestDurationSetValid(t *testing.T) {
	change := NewDuration(0, false)
	assertNullDuration(t, change, "SetValid()")
	change.SetValid(9203685806)
	assertDuration(t, change, "SetValid()")
}

func TestDurationScan(t *testing.T) {
	var i Duration
	err := i.Scan(9203685806)
	maybePanic(err)
	assertDuration(t, i, "scanned int64")

	var null Duration
	err = null.Scan(nil)
	maybePanic(err)
	assertNullDuration(t, null, "scanned null")
}

func assertDuration(t *testing.T, i Duration, from string) {
	if i.Duration != 9203685806 {
		t.Errorf("bad %s int64: %d ≠ %d\n", from, i.Duration, 9203685806)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullDuration(t *testing.T, i Duration, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
