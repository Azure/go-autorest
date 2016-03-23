package date

import (
	"strings"
	"time"
)

// ParseTime to parse Time string to specified format.
func ParseTime(format string, t string) (d time.Time, err error) {
	return parseTime(format, t)
}

//parseTime parses Time string after converting it to uppercase.
func parseTime(format string, t string) (time.Time, error) {
	d, err := time.Parse(format, strings.ToUpper(t))
	return d, err
}
