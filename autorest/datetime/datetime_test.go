package datetime

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func ExampleSimpleParse() {
	fmt.Println(Parse("2001-02-03T04:05:06Z"))
	// Output: 2001-02-03T04:05:06Z <nil>
}

func ExampleMarshalBinary() {
	d, _ := Parse("2001-02-03T04:05:06Z")
	t, _ := d.MarshalBinary()
	fmt.Println(string(t))
	// Output: 2001-02-03T04:05:06Z
}

func ExampleUnmarshalBinary() {
	d := DateTime{}
	t := "2001-02-03T04:05:06Z"

	err := d.UnmarshalBinary([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: 2001-02-03T04:05:06Z
}

func ExampleMarshalJSON() {
	d, _ := Parse("2001-02-03T04:05:06Z")
	j, _ := json.Marshal(d)
	fmt.Println(string(j))
	// Output: "2001-02-03T04:05:06Z"
}

func ExampleUnmarshalJSON() {
	var d struct {
		DateTime DateTime `json:"datetime"`
	}
	j := `{
    "datetime" : "2001-02-03T04:05:06Z"
  }`

	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d.DateTime)
	// Output: 2001-02-03T04:05:06Z
}

func ExampleMarshalText() {
	d, _ := Parse("2001-02-03T04:05:06Z")
	t, _ := d.MarshalText()
	fmt.Println(string(t))
	// Output: 2001-02-03T04:05:06Z
}

func ExampleUnmarshalText() {
	d := DateTime{}
	t := "2001-02-03T04:05:06Z"

	err := d.UnmarshalText([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: 2001-02-03T04:05:06Z
}

func ExampleString() {
	d, _ := Parse("2001-02-03T04:05:06Z")
	fmt.Printf("Date is %s", d)
	// Output: Date is 2001-02-03T04:05:06Z
}

func TestBinaryRoundTrip(t *testing.T) {
	d1, err := Parse("2001-02-03T04:05:06Z")
	t1, err := d1.MarshalBinary()
	if err != nil {
		t.Errorf("datetime: MarshalBinary failed (%v)", err)
	}

	d2 := DateTime{}
	err = d2.UnmarshalBinary(t1)
	if err != nil {
		t.Errorf("datetime: UnmarshalBinary failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("datetime: Round-trip Binary failed (%v, %v)", d1, d2)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	type s struct {
		DateTime DateTime `json:"datetime"`
	}
	var err error
	d1 := s{}
	d1.DateTime, err = Parse("2001-02-03T04:05:06Z")
	j, err := json.Marshal(d1)
	if err != nil {
		t.Errorf("datetime: MarshalJSON failed (%v)", err)
	}

	d2 := s{}
	err = json.Unmarshal(j, &d2)
	if err != nil {
		t.Errorf("datetime: UnmarshalJSON failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("datetime: Round-trip JSON failed (%v, %v)", d1, d2)
	}
}

func TestTextRoundTrip(t *testing.T) {
	d1, err := Parse("2001-02-03T04:05:06Z")
	t1, err := d1.MarshalText()
	if err != nil {
		t.Errorf("datetime: MarshalText failed (%v)", err)
	}

	d2 := DateTime{}
	err = d2.UnmarshalText(t1)
	if err != nil {
		t.Errorf("datetime: UnmarshalText failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("datetime: Round-trip Text failed (%v, %v)", d1, d2)
	}
}

func TestToTime(t *testing.T) {
	var d DateTime
	d, err := Parse("2001-02-03T04:05:06Z")
	if err != nil {
		t.Errorf("datetime: Parse failed (%v)", err)
	}
	var v interface{} = d.ToTime()
	switch v.(type) {
	case time.Time:
		return
	default:
		t.Errorf("datetime: ToTime failed to return a time.Time")
	}
}
