package date

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/mocks"
	"reflect"
	"testing"
	"time"
)

func ExampleParseTime() {
	d, _ := ParseTime("2001-02-03T04:05:06Z")
	fmt.Println(d)
	// Output: 2001-02-03T04:05:06Z
}

func ExampleTime_MarshalBinary() {
	d, _ := ParseTime("2001-02-03T04:05:06Z")
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
	d, _ := ParseTime("2001-02-03T04:05:06Z")
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
	d, _ := ParseTime("2001-02-03T04:05:06Z")
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

func TestTimeString(t *testing.T) {
	d, _ := ParseTime("2001-02-03T04:05:06Z")
	if d.String() != "2001-02-03T04:05:06Z" {
		t.Errorf("date: String failed (%v)", d.String())
	}
}

func TestTimeStringReturnsEmptyStringForError(t *testing.T) {
	d := Time{Time: time.Date(20000, 01, 01, 01, 01, 01, 01, time.UTC)}
	if d.String() != "" {
		t.Errorf("date: Time#String failed empty string for an error")
	}
}

func TestTimeBinaryRoundTrip(t *testing.T) {
	d1, err := ParseTime("2001-02-03T04:05:06Z")
	t1, err := d1.MarshalBinary()
	if err != nil {
		t.Errorf("datetime: MarshalBinary failed (%v)", err)
	}

	d2 := Time{}
	err = d2.UnmarshalBinary(t1)
	if err != nil {
		t.Errorf("datetime: UnmarshalBinary failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("datetime: Round-trip Binary failed (%v, %v)", d1, d2)
	}
}

func TestTimeJSONRoundTrip(t *testing.T) {
	type s struct {
		Time Time `json:"datetime"`
	}
	var err error
	d1 := s{}
	d1.Time, err = ParseTime("2001-02-03T04:05:06Z")
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

func TestTimeTextRoundTrip(t *testing.T) {
	d1, err := ParseTime("2001-02-03T04:05:06Z")
	t1, err := d1.MarshalText()
	if err != nil {
		t.Errorf("datetime: MarshalText failed (%v)", err)
	}

	d2 := Time{}
	err = d2.UnmarshalText(t1)
	if err != nil {
		t.Errorf("datetime: UnmarshalText failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("datetime: Round-trip Text failed (%v, %v)", d1, d2)
	}
}

func TestTimeToTime(t *testing.T) {
	var d Time
	d, err := ParseTime("2001-02-03T04:05:06Z")
	if err != nil {
		t.Errorf("datetime: ParseTime failed (%v)", err)
	}
	var v interface{} = d.ToTime()
	switch v.(type) {
	case time.Time:
		return
	default:
		t.Errorf("datetime: ToTime failed to return a time.Time")
	}
}

func TestDate_ByUnmarshallingTime(t *testing.T) {
	t1 := Time{}
	t2 := Time{Time: time.Date(2001, time.February, 3, 4, 5, 6, 0, time.UTC)}

	r := mocks.NewResponseWithContent("2001-02-03T04:05:06Z")
	err := autorest.Respond(r,
		ByUnmarshallingTime(&t1),
		autorest.ByClosing())
	if err != nil {
		t.Errorf("date: ByUnmarshallingTime failed (%v)", err)
	}
	if !reflect.DeepEqual(t1, t2) {
		t.Errorf("date: ByUnmarshallingTime failed to properly unmarshall -- expected %v, received %v", t2, t1)
	}
}

func TestDate_ByUnmarshallingTimeFailsWithInvalidValues(t *testing.T) {
	t1 := Time{}

	r := mocks.NewResponseWithContent("Not a Time")
	err := autorest.Respond(r,
		ByUnmarshallingTime(&t1),
		autorest.ByClosing())
	if err == nil {
		t.Errorf("date: ByUnmarshallingTime failed to return an error for an invalid string")
	}
}

func TestDate_ByUnmarshallingTimeFailsWithInvalidReader(t *testing.T) {
	t1 := Time{}

	r := mocks.NewResponseWithContent("2001-02-03T04:05:06Z")
	r.Body.Close()
	err := autorest.Respond(r,
		ByUnmarshallingTime(&t1),
		autorest.ByClosing())
	if err == nil {
		t.Errorf("date: ByUnmarshallingTime failed to return an error for an invalid io.Reader")
	}
}
