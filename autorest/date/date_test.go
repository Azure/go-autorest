package date

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func ExampleSimpleParse() {
	fmt.Println(Parse("2001-02-03"))
	// Output: 2001-02-03 <nil>
}

func ExampleMarshalBinary() {
	d, _ := Parse("2001-02-03")
	t, _ := d.MarshalBinary()
	fmt.Println(string(t))
	// Output: 2001-02-03
}

func ExampleUnmarshalBinary() {
	d := Date{}
	t := "2001-02-03"

	err := d.UnmarshalBinary([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: 2001-02-03
}

func ExampleMarshalJSON() {
	d, _ := Parse("2001-02-03")
	j, _ := json.Marshal(d)
	fmt.Println(string(j))
	// Output: "2001-02-03"
}

func ExampleUnmarshalJSON() {
	var d struct {
		Date Date `json:"date"`
	}
	j := `{
    "date" : "2001-02-03"
  }`

	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d.Date)
	// Output: 2001-02-03
}

func ExampleMarshalText() {
	d, _ := Parse("2001-02-03")
	t, _ := d.MarshalText()
	fmt.Println(string(t))
	// Output: 2001-02-03
}

func ExampleUnmarshalText() {
	d := Date{}
	t := "2001-02-03"

	err := d.UnmarshalText([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: 2001-02-03
}

func ExampleString() {
	d, _ := Parse("2001-02-03")
	fmt.Printf("Date is %s", d)
	// Output: Date is 2001-02-03
}

func TestBinaryRoundTrip(t *testing.T) {
	d1, err := Parse("2001-02-03")
	t1, err := d1.MarshalBinary()
	if err != nil {
		t.Errorf("date: MarshalBinary failed (%v)", err)
	}

	d2 := Date{}
	err = d2.UnmarshalBinary(t1)
	if err != nil {
		t.Errorf("date: UnmarshalBinary failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip Binary failed (%v, %v)", d1, d2)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	type s struct {
		Date Date `json:"date"`
	}
	var err error
	d1 := s{}
	d1.Date, err = Parse("2001-02-03")
	j, err := json.Marshal(d1)
	if err != nil {
		t.Errorf("date: MarshalJSON failed (%v)", err)
	}

	d2 := s{}
	err = json.Unmarshal(j, &d2)
	if err != nil {
		t.Errorf("date: UnmarshalJSON failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip JSON failed (%v, %v)", d1, d2)
	}
}

func TestTextRoundTrip(t *testing.T) {
	d1, err := Parse("2001-02-03")
	t1, err := d1.MarshalText()
	if err != nil {
		t.Errorf("date: MarshalText failed (%v)", err)
	}

	d2 := Date{}
	err = d2.UnmarshalText(t1)
	if err != nil {
		t.Errorf("date: UnmarshalText failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip Text failed (%v, %v)", d1, d2)
	}
}

func TestToTime(t *testing.T) {
	var d Date
	d, err := Parse("2001-02-03")
	if err != nil {
		t.Errorf("date: Parse failed (%v)", err)
	}
	var v interface{} = d.ToTime()
	switch v.(type) {
	case time.Time:
		return
	default:
		t.Errorf("date: ToTime failed to return a time.Time")
	}
}
