package date

import (
	"time"
)

const (
	rfc3339JSON = `"` + time.RFC3339Nano + `"`
	rfc3339     = time.RFC3339Nano
)

// Time defines a type similar to time.Time but assumes a layout of RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
type Time struct {
	time.Time
}

// MarshalBinary preserves the Time as a byte array conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (d Time) MarshalBinary() ([]byte, error) {
	return d.Time.MarshalText()
}

// UnmarshalBinary reconstitutes a Time saved as a byte array conforming to RFC3339 date-time
// (i.e., 2006-01-02T15:04:05Z).
func (d *Time) UnmarshalBinary(data []byte) error {
	return d.UnmarshalText(data)
}

// MarshalJSON preserves the Time as a JSON string conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (d Time) MarshalJSON() (json []byte, err error) {
	return d.Time.MarshalJSON()
}

// UnmarshalJSON reconstitutes the Time from a JSON string conforming to RFC3339 date-time
// (i.e., 2006-01-02T15:04:05Z).
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = ParseTime(rfc3339JSON, string(data))
	if err != nil {
		return err
	}
	return nil
}

// MarshalText preserves the Time as a byte array conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (t Time) MarshalText() (text []byte, err error) {
	return t.Time.MarshalText()
}

// UnmarshalText reconstitutes a Time saved as a byte array conforming to RFC3339 date-time
// (i.e., 2006-01-02T15:04:05Z).
func (t *Time) UnmarshalText(data []byte) (err error) {
	t.Time, err = ParseTime(rfc3339, string(data))
	if err != nil {
		return err
	}
	return nil
}

// String returns the Time formatted as an RFC3339 date-time string (i.e.,
// 2006-01-02T15:04:05Z).
func (t Time) String() string {
	// Note: time.Time.String does not return an RFC3339 compliant string, time.Time.MarshalText does.
	b, err := t.MarshalText()
	if err != nil {
		return ""
	}
	return string(b)
}

// ToTime returns a Time as a time.Time
func (d Time) ToTime() time.Time {
	return d.Time
}
