package date_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/marstr/go-autorest/autorest/date"
)

func ExampleUnixTime_MarshalJSON() {
	epoch := date.UnixTime(date.UnixEpoch())
	text, _ := json.Marshal(epoch)
	fmt.Print(string(text))
	// Output: 0
}

func ExampleUnixTime_UnmarshalSON() {
	var myTime date.UnixTime
	json.Unmarshal([]byte("1.3e2"), &myTime)
	fmt.Printf("%v", time.Time(myTime))
	// Output: 1970-01-01 00:02:10 +0000 UTC
}

func TestUnixTime_MarshalJSON(t *testing.T) {
	testCases := []time.Time{
		date.UnixEpoch().Add(-1 * time.Second),                   // One second befote the Unix Epoch
		time.Date(2017, time.April, 14, 20, 27, 47, 0, time.UTC), // The time this test was written
	}

	for _, tc := range testCases {
		t.Run("", func(subT *testing.T) {
			target := date.UnixTime(tc)

			expected := string([]byte(fmt.Sprintf("%d", tc.Unix())))
			if actual, err := json.Marshal(target); err != nil {
				subT.Error(err)
			} else if expected != string(actual) {
				subT.Logf("got: \t%s\nwant:\t%s", string(actual), expected)
				subT.Fail()
			} else {
				subT.Logf("passed with value: %s", string(actual))
			}
		})
	}
}

func TestUnixTime_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		text     string
		expected time.Time
	}{
		{"1", date.UnixEpoch().Add(time.Second)},
		{"0", date.UnixEpoch()},
		{"1492203742", time.Date(2017, time.April, 14, 21, 02, 22, 0, time.UTC)}, // The time this test was written
		{"-1", time.Date(1969, time.December, 31, 23, 59, 59, 0, time.UTC)},
		{"1.5", date.UnixEpoch().Add(1500 * time.Millisecond)},
		{"0e1", date.UnixEpoch()}, // See http://json.org for 'number' format definition.
		{"1.3e+2", date.UnixEpoch().Add(130 * time.Second)},
		{"1.6E-10", date.UnixEpoch()}, // This is so small, it should get truncated into the UnixEpoch
		{"2E-6", date.UnixEpoch().Add(2 * time.Microsecond)},
		{"1.289345e9", date.UnixEpoch().Add(1289345000 * time.Second)},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(subT *testing.T) {
			var rehydrated date.UnixTime
			if err := json.Unmarshal([]byte(tc.text), &rehydrated); err != nil {
				subT.Error(err)
				return
			}

			if time.Time(rehydrated) != tc.expected {
				subT.Logf("\ngot: \t%v\nwant:\t%v\ndiff:\t%v", time.Time(rehydrated), tc.expected, time.Time(rehydrated).Sub(tc.expected))
				subT.Fail()
			} else {
				subT.Logf("rehydrated matched expected '%v'", tc.expected)
			}
		})
	}
}

func TestUnixTime_JSONRoundTrip(t *testing.T) {
	testCases := []time.Time{
		date.UnixEpoch(),
		time.Date(2005, time.November, 5, 0, 0, 0, 0, time.UTC), // The day V for Vendetta (film) was released.
		date.UnixEpoch().Add(-6 * time.Second),
		date.UnixEpoch().Add(800 & time.Hour),
	}

	for _, tc := range testCases {
		t.Run(tc.String(), func(subT *testing.T) {
			subject := date.UnixTime(tc)
			var marshaled []byte
			if temp, err := json.Marshal(subject); err == nil {
				marshaled = temp
			} else {
				t.Error(err)
				return
			}

			var unmarshaled date.UnixTime
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Error(err)
				return
			} else if time.Time(subject) != time.Time(unmarshaled) {
				t.Logf("round trip failed for: %v", time.Time(subject))
				t.Fail()
			}
		})
	}
}

func TestUnixTime_MarshalBinary(t *testing.T) {
	testCases := []struct {
		expected int64
		subject  time.Time
	}{
		{0, date.UnixEpoch()},
		{-15 * int64(time.Second), date.UnixEpoch().Add(-15 * time.Second)},
		{54, date.UnixEpoch().Add(54 * time.Nanosecond)},
	}

	for _, tc := range testCases {
		t.Run("", func(subT *testing.T) {
			var marshaled []byte

			if temp, err := date.UnixTime(tc.subject).MarshalBinary(); err == nil {
				marshaled = temp
			} else {
				subT.Error(err)
				return
			}

			var unmarshaled int64
			if err := binary.Read(bytes.NewReader(marshaled), binary.LittleEndian, &unmarshaled); err != nil {
				subT.Error(err)
				return
			}

			if unmarshaled != tc.expected {
				subT.Logf("\ngot: \t%d\nwant:\t%d", unmarshaled, tc.expected)
				subT.Fail()
			}
		})
	}
}

func TestUnixTime_BinaryRoundTrip(t *testing.T) {
	testCases := []time.Time{
		date.UnixEpoch(),
		date.UnixEpoch().Add(800 * time.Minute),
		date.UnixEpoch().Add(7 * time.Hour),
		date.UnixEpoch().Add(-1 * time.Nanosecond),
	}

	for _, tc := range testCases {
		t.Run(tc.String(), func(subT *testing.T) {
			original := date.UnixTime(tc)
			var marshaled []byte

			if temp, err := original.MarshalBinary(); err == nil {
				marshaled = temp
			} else {
				subT.Error(err)
				return
			}

			var traveled date.UnixTime
			if err := traveled.UnmarshalBinary(marshaled); err != nil {
				subT.Error(err)
				return
			}

			if traveled != original {
				subT.Logf("\ngot: \t%s\nwant:\t%s", time.Time(original).String(), time.Time(traveled).String())
				subT.Fail()
			}
		})
	}
}

func TestUnixTime_MarshalText(t *testing.T) {
	testCases := []time.Time{
		date.UnixEpoch(),
		date.UnixEpoch().Add(45 * time.Second),
		date.UnixEpoch().Add(time.Nanosecond),
		date.UnixEpoch().Add(-100000 * time.Second),
	}

	for _, tc := range testCases {
		expected, _ := tc.MarshalText()
		t.Run("", func(subT *testing.T) {
			var marshaled []byte

			if temp, err := date.UnixTime(tc).MarshalText(); err == nil {
				marshaled = temp
			} else {
				subT.Error(err)
				return
			}

			if string(marshaled) != string(expected) {
				subT.Logf("\ngot: \t%s\nwant:\t%s", string(marshaled), string(expected))
				subT.Fail()
			}
		})
	}
}

func TestUnixTime_TextRoundTrip(t *testing.T) {
	testCases := []time.Time{
		date.UnixEpoch(),
		date.UnixEpoch().Add(-1 * time.Nanosecond),
		date.UnixEpoch().Add(1 * time.Nanosecond),
		time.Date(2017, time.April, 17, 21, 00, 00, 00, time.UTC),
	}

	for _, tc := range testCases {
		t.Run(tc.String(), func(subT *testing.T) {
			unixTC := date.UnixTime(tc)

			var marshaled []byte

			if temp, err := unixTC.MarshalText(); err == nil {
				marshaled = temp
			} else {
				subT.Error(err)
				return
			}

			var unmarshaled date.UnixTime
			if err := unmarshaled.UnmarshalText(marshaled); err != nil {
				subT.Error(err)
				return
			}

			if unmarshaled != unixTC {
				t.Logf("\ngot: \t%s\nwant:\t%s", time.Time(unmarshaled).String(), tc.String())
				t.Fail()
			}
		})
	}
}
