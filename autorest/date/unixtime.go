package date

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"
)

// unixEpoch is the moment in time that should be treated as timestamp 0.
var unixEpoch = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

// UnixTime marshals and unmarshals a time that is represented as the number
// of seconds (ignoring skip-seconds) since the Unix Epoch.
type UnixTime time.Time

// Duration returns the time as a Duration since the UnixEpoch.
func (t UnixTime) Duration() time.Duration {
	return time.Time(t).Sub(unixEpoch)
}

// FromSeconds creates a UnixTime as a number of seconds from the UnixEpoch.
func FromSeconds(seconds float64) UnixTime {
	return UnixTime(UnixEpoch().Add(time.Duration(seconds * float64(time.Second))))
}

// FromNanoseconds creates a UnixTime as a number of nanoseconds from the UnixEpoch.
func FromNanoseconds(nanoseconds int64) UnixTime {
	return UnixTime(UnixEpoch().Add(time.Duration(nanoseconds)))
}

// UnixEpoch retreives the moment considered the Unix Epoch. I.e. The time represented by '0'
func UnixEpoch() time.Time {
	return unixEpoch
}

// MarshalJSON preserves the UnixTime as a JSON number conforming to Unix Timestamp requirements.
// (i.e. the number of seconds since midnight January 1st, 1970 not considering leap seconds.)
func (t UnixTime) MarshalJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	enc.Encode(float64(time.Time(t).Unix()))
	return buffer.Bytes(), nil
}

// UnmarshalJSON reconstitures a UnixTime saved as a JSON number of the number of seconds since
// midnight January 1st, 1970.
func (t *UnixTime) UnmarshalJSON(text []byte) error {
	dec := json.NewDecoder(bytes.NewReader(text))

	var secondsSinceEpoch float64
	if err := dec.Decode(&secondsSinceEpoch); err != nil {
		return err
	}

	*t = FromSeconds(secondsSinceEpoch)

	return nil
}

// MarshalText stores the number of seconds since the Unix Epoch as a textual floating point number.
func (t UnixTime) MarshalText() ([]byte, error) {
	cast := time.Time(t)
	return cast.MarshalText()
}

// UnmarshalText populates a UnixTime with a value stored textually as a floating point number of seconds since the Unix Epoch.
func (t *UnixTime) UnmarshalText(raw []byte) error {
	var unmarshaled time.Time

	if err := unmarshaled.UnmarshalText(raw); err != nil {
		return err
	}

	*t = UnixTime(unmarshaled)
	return nil
}

// MarshalBinary converts a UnixTime into a binary.LittleEndian float64 of nanoseconds since the epoch.
func (t UnixTime) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}

	payload := int64(t.Duration())

	if err := binary.Write(buf, binary.LittleEndian, &payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary converts a from a binary.LittleEndian float64 of nanoseconds since the epoch into a UnixTime.
func (t *UnixTime) UnmarshalBinary(raw []byte) error {
	var secondsSinceEpoch int64

	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, &secondsSinceEpoch); err != nil {
		return err
	}
	*t = FromNanoseconds(secondsSinceEpoch)
	return nil
}
