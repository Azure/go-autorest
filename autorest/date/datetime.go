package date

import (
	"time"
)

// Defines a type similar to time.Time but assumes a layout of RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
type DateTime struct {
	time.Time
}

// Create a new DateTime from the passed string.
func ParseDateTime(date string) (d DateTime, err error) {
	d = DateTime{}
	d.Time, err = time.Parse(time.RFC3339, date)
	return d, err
}

// Preserves the DateTime as a byte array conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (d DateTime) MarshalBinary() ([]byte, error) {
	return d.Time.MarshalText()
}

// Reconstitutes a DateTime saved as a byte array conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (d *DateTime) UnmarshalBinary(data []byte) error {
	return d.Time.UnmarshalText(data)
}

// Preserves the DateTime as a JSON string conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (d DateTime) MarshalJSON() (json []byte, err error) {
	return d.Time.MarshalJSON()
}

// Reconstitutes the DateTime from a JSON string conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (d *DateTime) UnmarshalJSON(data []byte) (err error) {
	return d.Time.UnmarshalJSON(data)
}

// Preserves the DateTime as a byte array conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (d DateTime) MarshalText() (text []byte, err error) {
	return d.Time.MarshalText()
}

// Reconstitutes a DateTime saved as a byte array conforming to RFC3339 date-time (i.e.,
// 2006-01-02T15:04:05Z).
func (d *DateTime) UnmarshalText(data []byte) (err error) {
	return d.Time.UnmarshalText(data)
}

// Returns the DateTime formatted as an RFC3339 date-time string (i.e., 2006-01-02T15:04:05Z).
func (d DateTime) String() string {
	// Note: time.Time.String does not return an RFC3339 compliant string, time.Time.MarshalText does.
	b, err := d.Time.MarshalText()
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

// Returns a DateTime as a time.Time
func (d DateTime) ToTime() time.Time {
	return d.Time
}
