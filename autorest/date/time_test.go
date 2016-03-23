package date

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func ExampleParseTime() {
	d, _ := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	fmt.Println(d)
	// Output: 2001-02-03 04:05:06 +0000 UTC
}

func ExampleTime_MarshalBinary() {
	ti, _ := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	d := Time{ti}
	t, _ := d.MarshalBinary()
	fmt.Println(string(t))
	// Output: 2001-02-03T04:05:06Z
}

func ExampleTime_UnmarshalBinary() {
	d := Time{}
	t := "2001-02-03T04:05:06Z"

	err := d.UnmarshalBinary([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: 2001-02-03T04:05:06Z
}

func ExampleTime_MarshalJSON() {
	d, _ := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	j, _ := json.Marshal(d)
	fmt.Println(string(j))
	// Output: "2001-02-03T04:05:06Z"
}

func ExampleTime_UnmarshalJSON() {
	var d struct {
		Time Time `json:"datetime"`
	}
	j := `{
    "datetime" : "2001-02-03T04:05:06Z"
  }`

	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d.Time)
	// Output: 2001-02-03T04:05:06Z
}

func ExampleTime_MarshalText() {
	d, _ := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	t, _ := d.MarshalText()
	fmt.Println(string(t))
	// Output: 2001-02-03T04:05:06Z
}

func ExampleTime_UnmarshalText() {
	d := Time{}
	t := "2001-02-03T04:05:06Z"

	err := d.UnmarshalText([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: 2001-02-03T04:05:06Z
}

func TestUnmarshalTextforInvalidDate(t *testing.T) {
	d := Time{}
	dt := "2001-02-03T04:05:06AAA"

	err := d.UnmarshalText([]byte(dt))
	if err == nil {
		t.Errorf("date: Time#Unmarshal failed for invalid date")
	}
}

func TestUnmarshalJSONforInvalidDate(t *testing.T) {
	d := Time{}
	dt := `"2001-02-03T04:05:06AAA"`

	err := d.UnmarshalJSON([]byte(dt))
	if err == nil {
		t.Errorf("date: Time#Unmarshal failed for invalid date")
	}
}

func TestTimeString(t *testing.T) {
	ti, _ := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	d := Time{ti}
	if d.String() != "2001-02-03T04:05:06Z" {
		t.Errorf("date: Time#String failed (%v)", d.String())
	}
}

func TestTimeStringReturnsEmptyStringForError(t *testing.T) {
	d := Time{Time: time.Date(20000, 01, 01, 01, 01, 01, 01, time.UTC)}
	if d.String() != "" {
		t.Errorf("date: Time#String failed empty string for an error")
	}
}

func TestTimeBinaryRoundTrip(t *testing.T) {
	ti, err := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	d1 := Time{ti}
	t1, err := d1.MarshalBinary()
	if err != nil {
		t.Errorf("date: Time#MarshalBinary failed (%v)", err)
	}

	d2 := Time{}
	err = d2.UnmarshalBinary(t1)
	if err != nil {
		t.Errorf("date: Time#UnmarshalBinary failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date:Round-trip Binary failed (%v, %v)", d1, d2)
	}
}

func TestTimeJSONRoundTrip(t *testing.T) {
	type s struct {
		Time Time `json:"datetime"`
	}
	var err error
	ti, err := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	d1 := s{Time: Time{ti}}
	j, err := json.Marshal(d1)
	if err != nil {
		t.Errorf("date: Time#MarshalJSON failed (%v)", err)
	}

	d2 := s{}
	err = json.Unmarshal(j, &d2)
	if err != nil {
		t.Errorf("date: Time#UnmarshalJSON failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip JSON failed (%v, %v)", d1, d2)
	}
}

func TestTimeTextRoundTrip(t *testing.T) {
	ti, err := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	d1 := Time{Time: ti}
	t1, err := d1.MarshalText()
	if err != nil {
		t.Errorf("date: Time#MarshalText failed (%v)", err)
	}

	d2 := Time{}
	err = d2.UnmarshalText(t1)
	if err != nil {
		t.Errorf("date: Time#UnmarshalText failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip Text failed (%v, %v)", d1, d2)
	}
}

func TestTimeToTime(t *testing.T) {
	ti, err := ParseTime(rfc3339, "2001-02-03T04:05:06Z")
	d := Time{ti}
	if err != nil {
		t.Errorf("date: Time#ParseTime failed (%v)", err)
	}
	var v interface{} = d.ToTime()
	switch v.(type) {
	case time.Time:
		return
	default:
		t.Errorf("date: Time#ToTime failed to return a time.Time")
	}
}
