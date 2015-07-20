package date

import (
	"fmt"
	"time"
)

const (
	RFC3339FullDate = "2006-01-02"
	dateFormat      = "%4d-%02d-%02d"
	jsonFormat      = `"%4d-%02d-%02d"`
)

// Defines a type similar to time.Time but assumes a layout of RFC3339 full-date (i.e., 2006-01-02).
type Date struct {
	time.Time
}

// Create a new Date from the passed string.
func Parse(date string) (d Date, err error) {
	d = Date{}
	d.Time, err = time.Parse(RFC3339FullDate, date)
	return d, err
}

// Preserves the Date as a byte array conforming to RFC3339 full-date (i.e., 2006-01-02).
func (d Date) MarshalBinary() ([]byte, error) {
	return d.MarshalText()
}

// Reconstitutes a Date saved as a byte array conforming to RFC3339 full-date (i.e., 2006-01-02).
func (d *Date) UnmarshalBinary(data []byte) error {
	return d.UnmarshalText(data)
}

// Preserves the Date as a JSON string conforming to RFC3339 full-date (i.e., 2006-01-02).
func (d Date) MarshalJSON() (json []byte, err error) {
	return []byte(fmt.Sprintf(jsonFormat, d.Year(), d.Month(), d.Day())), nil
}

// Reconstitutes the Date from a JSON string conforming to RFC3339 full-date (i.e., 2006-01-02).
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	if data[0] == '"' {
		data = data[1 : len(data)-1]
	}
	d.Time, err = time.Parse(RFC3339FullDate, string(data))
	if err != nil {
		return err
	}
	return nil
}

// Preserves the Date as a byte array conforming to RFC3339 full-date (i.e., 2006-01-02).
func (d Date) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf(dateFormat, d.Year(), d.Month(), d.Day())), nil
}

// Reconstitutes a Date saved as a byte array conforming to RFC3339 full-date (i.e., 2006-01-02).
func (d *Date) UnmarshalText(data []byte) (err error) {
	d.Time, err = time.Parse(RFC3339FullDate, string(data))
	if err != nil {
		return err
	}
	return nil
}

// Returns the Date formatted as an RFC3339 full-date string (i.e., 2006-01-02).
func (d Date) String() string {
	return fmt.Sprintf(dateFormat, d.Year(), d.Month(), d.Day())
}

// Returns a Date as a time.Time
func (d Date) ToTime() time.Time {
	return d.Time
}
