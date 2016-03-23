package date

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func ExampleParseTimeRfc1123() {
	d, _ := ParseTime(rfc1123, "Mon, 02 Jan 2006 15:04:05 MST")
	fmt.Println(d)
	// Output: 2006-01-02 15:04:05 +0000 MST
}

func ExampleTimeRfc1123_MarshalBinary() {
	ti, _ := ParseTime(rfc1123, "Mon, 02 Jan 2006 15:04:05 MST")
	d := TimeRfc1123{ti}
	b, _ := d.MarshalBinary()
	fmt.Println(string(b))
	// Output: Mon, 02 Jan 2006 15:04:05 MST
}

func ExampleTimeRfc1123_UnmarshalBinary() {
	d := TimeRfc1123{}
	t := "Mon, 02 Jan 2006 15:04:05 MST"
	err := d.UnmarshalBinary([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: Mon, 02 Jan 2006 15:04:05 MST
}

func ExampleTimeRfc1123_MarshalJSON() {
	ti, _ := ParseTime(rfc1123, "Mon, 02 Jan 2006 15:04:05 MST")
	d := TimeRfc1123{ti}
	j, _ := json.Marshal(d)
	fmt.Println(string(j))
	// Output: "Mon, 02 Jan 2006 15:04:05 MST"
}

func TestTimeRfc1123MarshalJSONInvalid(t *testing.T) {
	ti := time.Date(20000, 01, 01, 00, 00, 00, 00, time.UTC)
	d := TimeRfc1123{ti}
	_, err := json.Marshal(d)
	if err == nil {
		t.Errorf("date: TimeRfc1123#Marshal failed for invalid date")
	}
}

func ExampleTimeRfc1123_UnmarshalJSON() {
	var d struct {
		Time TimeRfc1123 `json:"datetime"`
	}
	j := `{
    "datetime" : "Mon, 02 Jan 2006 15:04:05 MST"
  }`

	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d.Time)
	// Output: Mon, 02 Jan 2006 15:04:05 MST
}

func ExampleTimeRfc1123_MarshalText() {
	ti, _ := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	d := TimeRfc1123{ti}
	t, _ := d.MarshalText()
	fmt.Println(string(t))
	// Output: Sat, 03 Feb 2001 04:05:06 UTC
}

func ExampleTimeRfc1123_UnmarshalText() {
	d := TimeRfc1123{}
	t := "Sat, 03 Feb 2001 04:05:06 UTC"

	err := d.UnmarshalText([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: Sat, 03 Feb 2001 04:05:06 UTC
}

func TestUnmarshalJSONforInvalidDateRfc1123(t *testing.T) {
	dt := `"Mon, 02 Jan 2000000 15:05 MST"`
	d := TimeRfc1123{}
	err := d.UnmarshalJSON([]byte(dt))
	if err == nil {
		t.Errorf("date: TimeRfc1123#Unmarshal failed for invalid date")
	}
}

func TestUnmarshalTextforInvalidDateRfc1123(t *testing.T) {
	dt := "Mon, 02 Jan 2000000 15:05 MST"
	d := TimeRfc1123{}
	err := d.UnmarshalText([]byte(dt))
	if err == nil {
		t.Errorf("date: TimeRfc1123#Unmarshal failed for invalid date")
	}
}

func TestTimeStringRfc1123(t *testing.T) {
	ti, _ := ParseTime(rfc1123, "Mon, 02 Jan 2006 15:04:05 MST")
	d := TimeRfc1123{ti}
	if d.String() != "Mon, 02 Jan 2006 15:04:05 MST" {
		t.Errorf("date: TimeRfc1123#String failed (%v)", d.String())
	}
}

func TestTimeStringReturnsEmptyStringForErrorRfc1123(t *testing.T) {
	d := TimeRfc1123{Time: time.Date(20000, 01, 01, 01, 01, 01, 01, time.UTC)}
	if d.String() != "" {
		t.Errorf("date: TimeRfc1123#String failed empty string for an error")
	}
}

func TestTimeBinaryRoundTripRfc1123(t *testing.T) {
	ti, err := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	d1 := TimeRfc1123{ti}
	t1, err := d1.MarshalBinary()
	if err != nil {
		t.Errorf("date: TimeRfc1123#MarshalBinary failed (%v)", err)
	}

	d2 := TimeRfc1123{}
	err = d2.UnmarshalBinary(t1)
	if err != nil {
		t.Errorf("date: TimeRfc1123#UnmarshalBinary failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip Binary failed (%v, %v)", d1, d2)
	}
}

func TestTimeJSONRoundTripRfc1123(t *testing.T) {
	type s struct {
		Time TimeRfc1123 `json:"datetime"`
	}
	var err error
	ti, err := ParseTime(rfc1123, "Mon, 02 Jan 2006 15:04:05 MST")
	d1 := s{Time: TimeRfc1123{ti}}
	j, err := json.Marshal(d1)
	if err != nil {
		t.Errorf("date: TimeRfc1123#MarshalJSON failed (%v)", err)
	}

	d2 := s{}
	err = json.Unmarshal(j, &d2)
	if err != nil {
		t.Errorf("date: TimeRfc1123#UnmarshalJSON failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip JSON failed (%v, %v)", d1, d2)
	}
}

func TestTimeTextRoundTripRfc1123(t *testing.T) {
	ti, err := ParseTime(rfc1123, "Mon, 02 Jan 2006 15:04:05 MST")
	d1 := TimeRfc1123{Time: ti}
	t1, err := d1.MarshalText()
	if err != nil {
		t.Errorf("date: TimeRfc1123#MarshalText failed (%v)", err)
	}

	d2 := TimeRfc1123{}
	err = d2.UnmarshalText(t1)
	if err != nil {
		t.Errorf("date: TimeRfc1123#UnmarshalText failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip Text failed (%v, %v)", d1, d2)
	}
}

func TestTimeToTimeRfc1123(t *testing.T) {
	ti, err := ParseTime(rfc1123, "Mon, 02 Jan 2006 15:04:05 MST")
	d := TimeRfc1123{ti}
	if err != nil {
		t.Errorf("date: TimeRfc1123#ParseTime failed (%v)", err)
	}
	var v interface{} = d.ToTime()
	switch v.(type) {
	case time.Time:
		return
	default:
		t.Errorf("date: TimeRfc1123#ToTime failed to return a time.Time")
	}
}
