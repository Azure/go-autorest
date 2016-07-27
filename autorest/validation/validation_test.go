package validation

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCheckForUniqueInArrayTrue(t *testing.T) {
	if c := checkForUniqueInArray(reflect.ValueOf([]int{1, 2, 3})); !c {
		t.Fatalf("autorest/validation: checkForUniqueInArray failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInArrayFalse(t *testing.T) {
	if c := checkForUniqueInArray(reflect.ValueOf([]int{1, 2, 3, 3})); c {
		t.Fatalf("autorest/validation: checkForUniqueInArray failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInArrayEmpty(t *testing.T) {
	if c := checkForUniqueInArray(reflect.ValueOf([]int{})); c {
		t.Fatalf("autorest/validation: checkForUniqueInArray failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInMapTrue(t *testing.T) {
	if c := checkForUniqueInMap(reflect.ValueOf(map[string]int{"one": 1, "two": 2})); !c {
		t.Fatalf("autorest/validation: checkForUniqueInMap failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInMapFalse(t *testing.T) {
	if c := checkForUniqueInMap(reflect.ValueOf(map[int]string{1: "one", 2: "one"})); c {
		t.Fatalf("autorest/validation: checkForUniqueInMap failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInMapEmpty(t *testing.T) {
	if c := checkForUniqueInMap(reflect.ValueOf(map[int]string{})); c {
		t.Fatalf("autorest/validation: checkForUniqueInMap failed to check unique want: true; got: %v", c)
	}
}

func TestCheckEmpty_WithValueEmptyRuleTrue(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   true,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}

	if z := checkEmpty(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: checkEmpty failed to check Empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_WithValueNilRuleTrue(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   true,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}

	if z := checkEmpty(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: checkEmpty failed to check Empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_WithEmptyStringRuleFalse(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   false,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}

	if z := checkEmpty(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: checkEmpty failed to check Empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_WithNilRuleFalse(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   false,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := checkEmpty(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: checkEmpty failed to check Empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_IncorrectRule(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   10,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be bool value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := checkEmpty(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: checkEmpty failed to return error for incorrect rule \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_WithErrorArray(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   true,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}

	if z := checkEmpty(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: checkEmpty failed to check Empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckNil_WithNilValueRuleTrue(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "x",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{"x", MaxItems, 4, nil},
		},
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := checkNil(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: checkNil failed to return error for nil value \nexpect: nil;\ngot: %v", z)
	}
}

func TestCheckNil_WithNilValueRuleFalse(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "x",
		Name:   Null,
		Rule:   false,
		Chain: []Constraint{
			{"x", MaxItems, 4, nil},
		},
	}
	if z := checkNil(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: checkNil failed to return nil value \nexpect: nil;\ngot: %v", z)
	}
}

func TestCheckNil_IncorrectRule(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   Null,
		Rule:   10,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be bool value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := checkNil(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: checkNil failed to return error for incorrect rule \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_WithNilValueRuleTrue(t *testing.T) {
	var a []string
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   Null,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: a,
		Details:     "value can not be null; required parameter",
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error for null check \nexpect: %#v \ngot: %#v", e, z)
	}
}

func TestValidateArrayMap_WithNilValueRuleFalse(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "arr",
		Name:   Null,
		Rule:   false,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithValueRuleNullTrue(t *testing.T) {
	var x interface{} = []string{"1", "2"}
	c := Constraint{
		Target: "arr",
		Name:   Null,
		Rule:   false,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithEmptyValueRuleTrue(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "arr",
		Name:   Empty,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error for null check \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithEmptyValueRuleFalse(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "arr",
		Name:   Empty,
		Rule:   false,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithEmptyRuleEmptyTrue(t *testing.T) {
	var x interface{} = []string{"1", "2"}
	c := Constraint{
		Target: "arr",
		Name:   Empty,
		Rule:   false,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MaxItemsIncorrectRule(t *testing.T) {
	var x interface{} = []string{"1", "2"}
	c := Constraint{
		Target: "arr",
		Name:   MaxItems,
		Rule:   false,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error  \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MaxItemsNoError(t *testing.T) {
	var x interface{} = []string{"1", "2"}
	c := Constraint{
		Target: "arr",
		Name:   MaxItems,
		Rule:   2,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MaxItemsWithError(t *testing.T) {
	var x interface{} = []string{"1", "2", "3"}
	c := Constraint{
		Target: "arr",
		Name:   MaxItems,
		Rule:   2,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("maximum item limit is %v; got: 3", c.Rule),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_MaxItemsWithEmpty(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "arr",
		Name:   MaxItems,
		Rule:   2,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MinItemsIncorrectRule(t *testing.T) {
	var x interface{} = []int{1, 2}
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   false,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer for %v constraint; got: %v", c.Name, c.Rule),
	}

	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_MinItemsNoError1(t *testing.T) {
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   2,
		Chain:  nil,
	}

	if z := validateArrayMap(reflect.ValueOf([]int{1, 2}), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MinItemsNoError2(t *testing.T) {
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   2,
		Chain:  nil,
	}

	if z := validateArrayMap(reflect.ValueOf([]int{1, 2, 3}), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MinItemsWithError(t *testing.T) {
	var x interface{} = []int{1}
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   2,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("minimum item limit is %v; got: 1", c.Rule),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_MinItemsWithEmpty(t *testing.T) {
	var x interface{} = []int{}
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   2,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("minimum item limit is %v; got: 0", c.Rule),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MaxItemsIncorrectRule(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2"}
	c := Constraint{
		Target: "arr",
		Name:   MaxItems,
		Rule:   false,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MaxItemsNoError(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2"}
	c := Constraint{
		Target: "arr",
		Name:   MaxItems,
		Rule:   2,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MaxItemsWithError(t *testing.T) {
	a := map[int]string{1: "1", 2: "2", 3: "3"}
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   MaxItems,
		Rule:   2,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("maximum item limit is %v; got: %v", c.Rule, len(a)),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MaxItemsWithEmpty(t *testing.T) {
	a := map[int]string{}
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   MaxItems,
		Rule:   2,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MinItemsIncorrectRule(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2"}
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   false,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer for %v constraint; got: %v", c.Name, c.Rule),
	}

	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MinItemsNoError1(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2"}
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   2,
		Chain:  nil,
	}

	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MinItemsNoError2(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2", 3: "3"}
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   2,
		Chain:  nil,
	}

	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MinItemsWithError(t *testing.T) {
	a := map[int]string{1: "1"}
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   2,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("minimum item limit is %v; got: %v", c.Rule, len(a)),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MinItemsWithEmpty(t *testing.T) {
	a := map[int]string{}
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   MinItems,
		Rule:   2,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("minimum item limit is %v; got: %v", c.Rule, len(a)),
	}

	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MinItemsNil(t *testing.T) {
	var a map[int]float64
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   MinItems,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); e.Target != z.(Error).Target || e.Constraint != z.(Error).Constraint {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_UniqueItemsTrue(t *testing.T) {
	var x interface{} = map[float64]int{1.2: 1, 1.4: 2}
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_UniqueItemsFalse(t *testing.T) {
	var x interface{} = map[string]string{"1": "1", "2": "2", "3": "1"}
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); e.Target != z.(Error).Target ||
		e.Constraint != z.(Error).Constraint {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_UniqueItemsEmpty(t *testing.T) {
	// Consider Empty map as not unique returns false
	var x interface{} = map[int]float64{}
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); e.Target != z.(Error).Target || e.Constraint != z.(Error).Constraint {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_UniqueItemsNil(t *testing.T) {
	var a map[int]float64
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); e.Target != z.(Error).Target || e.Constraint != z.(Error).Constraint {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsTrue(t *testing.T) {
	var x interface{} = []int{1, 2}
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsFalse(t *testing.T) {
	var x interface{} = []string{"1", "2", "1"}
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); e.Target != z.(Error).Target || e.Constraint != z.(Error).Constraint {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsEmpty(t *testing.T) {
	// Consider Empty array as not unique returns false
	var x interface{} = []float64{}
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); e.Target != z.(Error).Target || e.Constraint != z.(Error).Constraint {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsNil(t *testing.T) {
	// Consider nil array as not unique returns false
	var a []float64
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}

	if z := validateArrayMap(reflect.ValueOf(x), c); e.Target != z.(Error).Target || e.Constraint != z.(Error).Constraint {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsInvalidType(t *testing.T) {
	var x interface{} = "hello"
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("type must be array, slice or map for constraint %v; got: %v", c.Name, reflect.ValueOf(x).Kind()),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsInvalidConstraint(t *testing.T) {
	var x interface{} = "hello"
	c := Constraint{
		Target: "str",
		Name:   "sdad",
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("constraint %v is not applicable to array, slice and map type", c.Name),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraint1(t *testing.T) {
	a := []int{1, 2, 3, 4}
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{"str", MaxItems, 3, nil},
		},
	}
	e := Error{
		Constraint:  (c.Chain)[0].Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("maximum item limit is %v; got: %v", (c.Chain)[0].Rule, len(a)),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraint2(t *testing.T) {
	a := []int{1, 2, 3, 4}
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   true,
		Chain: []Constraint{
			{"str", MaxItems, 3, nil},
		},
	}
	e := Error{
		Constraint:  (c.Chain)[0].Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("maximum item limit is %v; got: %v", (c.Chain)[0].Rule, len(a)),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraint3(t *testing.T) {
	var a []string
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{"str", MaxItems, 3, nil},
		},
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraint4(t *testing.T) {
	var x interface{} = []int{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   true,
		Chain: []Constraint{
			{"str", MaxItems, 3, nil},
		},
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraintNilNotRequired(t *testing.T) {
	var a []int
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   Null,
		Rule:   false,
		Chain: []Constraint{
			{"str", MaxItems, 3, nil},
		},
	}

	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_ValidateChainConstraintEmptyNotRequired(t *testing.T) {
	var x interface{} = map[string]int{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   false,
		Chain: []Constraint{
			{"str", MaxItems, 3, nil},
		},
	}

	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_ReadOnlyWithError(t *testing.T) {
	var x interface{} = []int{1, 2}
	c := Constraint{
		Target: "str",
		Name:   ReadOnly,
		Rule:   true,
		Chain: []Constraint{
			{"str", MaxItems, 3, nil},
		},
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("readonly parameter; must send as nil or Empty in request"),
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ReadOnlyWithoutError(t *testing.T) {
	var x interface{} = []int{}
	c := Constraint{
		Target: "str",
		Name:   ReadOnly,
		Rule:   true,
		Chain:  nil,
	}
	if z := validateArrayMap(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: nil\ngot: %v", z)
	}
}

func TestValidateString_ReadOnly(t *testing.T) {
	// Empty true means parameter is required but Empty returns error
	var x interface{} = "Hello Gopher"
	c := Constraint{
		Target: "str",
		Name:   ReadOnly,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("readonly parameter; must send as nil or empty in request"),
	}
	if z := validateString(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateString failed to return error for readOnly\nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_EmptyTrue(t *testing.T) {
	// Empty true means parameter is required but Empty returns error
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: "",
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateString(reflect.ValueOf(""), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateString failed to return error \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_EmptyFalse(t *testing.T) {
	// Empty false means parameter is not required and Empty return nil
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   Empty,
		Rule:   false,
		Chain:  nil,
	}
	if z := validateString(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_MaxLengthInvalid(t *testing.T) {
	// Empty true means parameter is required but Empty returns error
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   MaxLength,
		Rule:   4,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("string '%s' length must be less than %v", x, c.Rule),
	}
	if z := validateString(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateString failed to return error \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_MaxLengthValid(t *testing.T) {
	// Empty false means parameter is not required and Empty return nil
	c := Constraint{
		Target: "str",
		Name:   MaxLength,
		Rule:   7,
		Chain:  nil,
	}
	if z := validateString(reflect.ValueOf("Hello"), c); z != nil {
		t.Fatalf("autorest/validation: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_MaxLengthRuleInvalid(t *testing.T) {
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   MaxLength,
		Rule:   true, // must be int for maxLength
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateString(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateString_MinLengthInvalid(t *testing.T) {
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   MinLength,
		Rule:   10,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("string '%s' length must be greater than %v", x, c.Rule),
	}
	if z := validateString(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateString_MinLengthValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   MinLength,
		Rule:   2,
		Chain:  nil,
	}
	if z := validateString(reflect.ValueOf("Hello"), c); z != nil {
		t.Fatalf("autorest/validation: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_MinLengthRuleInvalid(t *testing.T) {
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   MinLength,
		Rule:   true, // must be int for minLength
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateString(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateString_PatternInvalidPattern(t *testing.T) {
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   Pattern,
		Rule:   "[[:alnum:",
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     "error parsing regexp: missing closing ]: `[[:alnum:$`",
	}

	if z := validateString(reflect.ValueOf(x), c); z.(Error).Details != e.Details {
		t.Fatalf("autorest/validation: validateString failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateString_PatternMatch1(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   Pattern,
		Rule:   "http://\\w+",
		Chain:  nil,
	}
	if z := validateString(reflect.ValueOf("http://masd"), c); z != nil {
		t.Fatalf("autorest/validation: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_PatternMatch2(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   Pattern,
		Rule:   "[a-zA-Z0-9]+",
		Chain:  nil,
	}
	if z := validateString(reflect.ValueOf("asdadad2323sad"), c); z != nil {
		t.Fatalf("autorest/validation: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_PatternNotMatch(t *testing.T) {
	var x interface{} = "asdad@@ad2323sad"
	c := Constraint{
		Target: "str",
		Name:   Pattern,
		Rule:   "[a-zA-Z0-9]+",
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("string '%v' doesn't match pattern %v", x, c.Rule),
	}
	if z := validateString(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateString_InvalidConstraint(t *testing.T) {
	var x interface{} = "asdad@@ad2323sad"
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   "[a-zA-Z0-9]+",
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("constraint %s is not applicable to string type", c.Name),
	}

	if z := validateString(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_InvalidConstraint(t *testing.T) {
	var x interface{} = 1.4
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   3.0,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("constraint %v is not applicable for type float", c.Name),
	}
	if z := validateFloat(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_InvalidRuleValue(t *testing.T) {
	var x interface{} = 1.4
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMinimum,
		Rule:   3,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be float value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateFloat(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateFloat failed to return nil \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_ExclusiveMinimumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMinimum,
		Rule:   1.0,
		Chain:  nil,
	}
	if z := validateFloat(reflect.ValueOf(1.42), c); z != nil {
		t.Fatalf("autorest/validation: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_ExclusiveMinimumConstraintInvalid(t *testing.T) {
	var x interface{} = 1.4
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMinimum,
		Rule:   1.5,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than %v", c.Rule),
	}
	if z := validateFloat(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_ExclusiveMinimumConstraintBoundary(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMinimum,
		Rule:   1.42,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than %v", c.Rule),
	}
	if z := validateFloat(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_exclusiveMaximumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMaximum,
		Rule:   2.0,
		Chain:  nil,
	}
	if z := validateFloat(reflect.ValueOf(1.42), c); z != nil {
		t.Fatalf("autorest/validation: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_exclusiveMaximumConstraintInvalid(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMaximum,
		Rule:   1.2,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than %v", c.Rule),
	}
	if z := validateFloat(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_exclusiveMaximumConstraintBoundary(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMaximum,
		Rule:   1.42,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than %v", c.Rule),
	}
	if z := validateFloat(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_inclusiveMaximumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   InclusiveMaximum,
		Rule:   2.0,
		Chain:  nil,
	}
	if z := validateFloat(reflect.ValueOf(1.42), c); z != nil {
		t.Fatalf("autorest/validation: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_inclusiveMaximumConstraintInvalid(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   InclusiveMaximum,
		Rule:   1.2,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than or equal to %v", c.Rule),
	}
	if z := validateFloat(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_inclusiveMaximumConstraintBoundary(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   InclusiveMaximum,
		Rule:   1.42,
		Chain:  nil,
	}
	if z := validateFloat(reflect.ValueOf(1.42), c); z != nil {
		t.Fatalf("autorest/validation: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_InclusiveMinimumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   InclusiveMinimum,
		Rule:   1.0,
		Chain:  nil,
	}
	if z := validateFloat(reflect.ValueOf(1.42), c); z != nil {
		t.Fatalf("autorest/validation: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_InclusiveMinimumConstraintInvalid(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   InclusiveMinimum,
		Rule:   1.5,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than or equal to %v", c.Rule),
	}
	if z := validateFloat(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_InclusiveMinimumConstraintBoundary(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   InclusiveMinimum,
		Rule:   1.42,
		Chain:  nil,
	}
	if z := validateFloat(reflect.ValueOf(1.42), c); z != nil {
		t.Fatalf("autorest/validation: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_InvalidConstraint(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   UniqueItems,
		Rule:   3,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("constraint %s is not applicable for type integer", c.Name),
	}
	if z := validateInt(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateInt failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateInt_InvalidRuleValue(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMinimum,
		Rule:   3.4,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateInt(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_ExclusiveMinimumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMinimum,
		Rule:   1,
		Chain:  nil,
	}
	if z := validateInt(reflect.ValueOf(3), c); z != nil {
		t.Fatalf("autorest/validation: valiateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_ExclusiveMinimumConstraintInvalid(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMinimum,
		Rule:   3,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than %v", c.Rule),
	}
	if z := validateInt(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_ExclusiveMinimumConstraintBoundary(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMinimum,
		Rule:   1,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than %v", c.Rule),
	}
	if z := validateInt(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_exclusiveMaximumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMaximum,
		Rule:   2,
		Chain:  nil,
	}

	if z := validateInt(reflect.ValueOf(1), c); z != nil {
		t.Fatalf("autorest/validation: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_exclusiveMaximumConstraintInvalid(t *testing.T) {
	var x interface{} = 2
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMaximum,
		Rule:   1,
		Chain:  nil,
	}

	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than %v", c.Rule),
	}
	if z := validateInt(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_exclusiveMaximumConstraintBoundary(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   ExclusiveMaximum,
		Rule:   1,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than %v", c.Rule),
	}

	if z := validateInt(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_inclusiveMaximumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   InclusiveMaximum,
		Rule:   2,
		Chain:  nil,
	}
	if z := validateInt(reflect.ValueOf(1), c); z != nil {
		t.Fatalf("autorest/validation: validateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_inclusiveMaximumConstraintInvalid(t *testing.T) {
	var x interface{} = 2
	c := Constraint{
		Target: "str",
		Name:   InclusiveMaximum,
		Rule:   1,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than or equal to %v", c.Rule),
	}
	if z := validateInt(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_inclusiveMaximumConstraintBoundary(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   InclusiveMaximum,
		Rule:   1,
		Chain:  nil,
	}
	if z := validateInt(reflect.ValueOf(1), c); z != nil {
		t.Fatalf("autorest/validation: validateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_InclusiveMinimumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   InclusiveMinimum,
		Rule:   1,
		Chain:  nil,
	}

	if z := validateInt(reflect.ValueOf(1), c); z != nil {
		t.Fatalf("autorest/validation: valiateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_InclusiveMinimumConstraintInvalid(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   InclusiveMinimum,
		Rule:   2,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than or equal to %v", c.Rule),
	}
	if z := validateInt(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_InclusiveMinimumConstraintBoundary(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   InclusiveMinimum,
		Rule:   1,
		Chain:  nil,
	}

	if z := validateInt(reflect.ValueOf(1), c); z != nil {
		t.Fatalf("autorest/validation: valiateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_MultipleOfWithoutError(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   MultipleOf,
		Rule:   10,
		Chain:  nil,
	}

	if z := validateInt(reflect.ValueOf(2300), c); z != nil {
		t.Fatalf("autorest/validation: valiateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_MultipleOfWithError(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   MultipleOf,
		Rule:   11,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: 2300,
		Details:     fmt.Sprintf("value must be a multiple of %v", c.Rule),
	}
	if z := validateInt(reflect.ValueOf(2300), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_NilTrue(t *testing.T) {
	var z *int
	var x interface{} = z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true, // Required property
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_NilFalse(t *testing.T) {
	var z *int
	var x interface{} = z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   false, // not required property
		Chain:  nil,
	}
	if z := validatePtr(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_NilReadonlyValid(t *testing.T) {
	var z *int
	var x interface{} = z
	c := Constraint{
		Target: "ptr",
		Name:   ReadOnly,
		Rule:   true,
		Chain:  nil,
	}
	if z := validatePtr(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_NilReadonlyInvalid(t *testing.T) {
	z := 10
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   ReadOnly,
		Rule:   true,
		Chain:  nil,
	}
	e := Error{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: z,
		Details:     "readonly parameter; must send as nil or Empty in request",
	}

	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_IntValid(t *testing.T) {
	z := 10
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   InclusiveMinimum,
		Rule:   3,
		Chain:  nil,
	}
	if z := validatePtr(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_IntInvalid(t *testing.T) {
	z := 10
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   InclusiveMinimum,
				Rule:   11,
				Chain:  nil,
			},
		},
	}
	e := Error{
		Constraint:  InclusiveMinimum,
		Target:      c.Target,
		TargetValue: z,
		Details:     "value must be greater than or equal to 11",
	}

	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}
func TestValidatePointer_IntInvalidConstraint(t *testing.T) {
	z := 10
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   MaxItems,
				Rule:   3,
				Chain:  nil,
			},
		},
	}

	e := Error{
		Constraint:  MaxItems,
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("constraint %v is not applicable for type integer", MaxItems),
	}
	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return correct error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_ValidInt64(t *testing.T) {
	z := int64(10)
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   InclusiveMinimum,
				Rule:   3,
				Chain:  nil,
			},
		}}
	if z := validatePtr(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_InvalidConstraintInt64(t *testing.T) {
	z := int64(10)
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   MaxItems,
				Rule:   3,
				Chain:  nil,
			},
		},
	}
	e := Error{
		Constraint:  MaxItems,
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("constraint %v is not applicable for type integer", MaxItems),
	}
	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_ValidFloat(t *testing.T) {
	z := 10.1
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   InclusiveMinimum,
				Rule:   3.0,
				Chain:  nil,
			}}}
	if z := validatePtr(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_InvalidFloat(t *testing.T) {
	z := 10.1
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   InclusiveMinimum,
				Rule:   12.0,
				Chain:  nil,
			}},
	}
	e := Error{
		Constraint:  InclusiveMinimum,
		Target:      c.Target,
		TargetValue: z,
		Details:     "value must be greater than or equal to 12",
	}
	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_InvalidConstraintFloat(t *testing.T) {
	z := 10.1
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   MaxItems,
				Rule:   3.0,
				Chain:  nil,
			}},
	}
	e := Error{
		Constraint:  MaxItems,
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("constraint %v is not applicable for type float", MaxItems),
	}
	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_StringValid(t *testing.T) {
	z := "hello"
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   Pattern,
				Rule:   "^[a-z]+$",
				Chain:  nil,
			}}}
	if z := validatePtr(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_StringInvalid(t *testing.T) {
	z := "hello"
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   MaxLength,
				Rule:   2,
				Chain:  nil,
			}}}

	e := Error{
		Constraint:  MaxLength,
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("string '%s' length must be less than 2", z),
	}
	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_ArrayValid(t *testing.T) {
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   UniqueItems,
				Rule:   "true",
				Chain:  nil,
			}}}
	if z := validatePtr(reflect.ValueOf(&[]string{"1", "2"}), c); z != nil {
		t.Fatalf("autorest/validation: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_ArrayInvalid(t *testing.T) {
	z := []string{"1", "2", "2"}
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{{
			Target: "ptr",
			Name:   UniqueItems,
			Rule:   true,
			Chain:  nil,
		}},
	}
	e := Error{
		Constraint:  UniqueItems,
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("all items in parameter %q must be unique; got:%v", c.Target, z),
	}
	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_MapValid(t *testing.T) {
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{
			{
				Target: "ptr",
				Name:   UniqueItems,
				Rule:   true,
				Chain:  nil,
			}}}
	if z := validatePtr(reflect.ValueOf(&map[interface{}]string{1: "1", "1": "2"}), c); z != nil {
		t.Fatalf("autorest/validation: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_MapInvalid(t *testing.T) {
	z := map[interface{}]string{1: "1", "1": "2", 1.3: "2"}
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   Null,
		Rule:   true,
		Chain: []Constraint{{
			Target: "ptr",
			Name:   UniqueItems,
			Rule:   true,
			Chain:  nil,
		}},
	}
	e := Error{
		Constraint:  UniqueItems,
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, z),
	}
	if z := validatePtr(reflect.ValueOf(x), c); e.Target != z.(Error).Target ||
		e.Constraint != z.(Error).Constraint || !reflect.DeepEqual(e.TargetValue, z.(Error).TargetValue) {
		t.Fatalf("autorest/validation: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

type Child struct {
	I string
}
type Product struct {
	C    *Child
	Str  *string
	Name string
	Arr  *[]string
	M    *map[string]string
	Num  *int32
}

type Sample struct {
	M    *map[string]*string
	Name string
}

func TestValidatePointer_StructWithError(t *testing.T) {
	s := "hello"
	var x interface{} = &Product{
		C:    &Child{"100"},
		Str:  &s,
		Name: "Gopher",
	}
	c := Constraint{
		"p", Null, "True",
		[]Constraint{
			{"C", Null, true,
				[]Constraint{
					{"I", MaxLength, 2, nil},
				}},
			{"Str", MaxLength, 2, nil},
			{"Name", MaxLength, 5, nil},
		},
	}
	e := Error{
		Constraint:  MaxLength,
		Target:      "I",
		TargetValue: "100",
		Details:     fmt.Sprintf("string '100' length must be less than 2"),
	}

	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validatePtr failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidatePointer_WithNilStruct(t *testing.T) {
	var p *Product
	var x interface{} = p
	c := Constraint{
		"p", Null, true,
		[]Constraint{
			{"C", Null, true,
				[]Constraint{
					{"I", Empty, true,
						[]Constraint{
							{"I", MaxLength, 5, nil},
						}},
				}},
			{"Str", MaxLength, 2, nil},
			{"Name", MaxLength, 5, nil},
		},
	}
	e := Error{
		Constraint:  Null,
		Target:      "p",
		TargetValue: p,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validatePtr(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validatePtr failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidatePointer_StructWithNoError(t *testing.T) {
	s := "hello"
	var x interface{} = &Product{
		C:    &Child{"100"},
		Str:  &s,
		Name: "Gopher",
	}
	c := Constraint{
		"p", Null, true,
		[]Constraint{
			{"C", Null, true,
				[]Constraint{
					{"I", Empty, true,
						[]Constraint{
							{"I", MaxLength, 5, nil},
						}},
				}},
		},
	}
	if z := validatePtr(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: validatePtr failed to nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidateStruct_FieldNotExist(t *testing.T) {
	s := "hello"
	var x interface{} = Product{
		C:    &Child{"100"},
		Str:  &s,
		Name: "Gopher",
	}
	c := Constraint{
		"C", Null, true,
		[]Constraint{
			{"Name", Empty, true, nil},
		},
	}

	s = "Name"
	e := Error{
		Constraint:  Empty,
		Target:      "Name",
		TargetValue: Child{"100"},
		Details:     fmt.Sprintf("field %q doesn't exist", s),
	}
	if z := validateStruct(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		fmt.Println(z)
		t.Fatalf("autorest/validation: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_WithChainConstraint(t *testing.T) {
	s := "hello"
	var x interface{} = Product{
		C:    &Child{"100"},
		Str:  &s,
		Name: "Gopher",
	}
	c := Constraint{
		"C", Null, true,
		[]Constraint{
			{"I", Empty, true,
				[]Constraint{
					{"I", MaxLength, 2, nil},
				}},
		},
	}
	e := Error{
		Constraint:  MaxLength,
		Target:      "I",
		TargetValue: "100",
		Details:     "string '100' length must be less than 2",
	}
	if z := validateStruct(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_WithoutChainConstraint(t *testing.T) {
	s := "hello"
	var x interface{} = Product{
		C:    &Child{""},
		Str:  &s,
		Name: "Gopher",
	}
	c := Constraint{"C", Null, true,
		[]Constraint{
			{"I", Empty, true, nil}, // throw error for Empty
		}}
	e := Error{
		Constraint:  Empty,
		Target:      "I",
		TargetValue: "",

		Details: fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateStruct(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_WithArrayNull(t *testing.T) {
	s := "hello"
	var x interface{} = Product{
		C:    &Child{""},
		Str:  &s,
		Name: "Gopher",
		Arr:  nil,
	}
	c := Constraint{"Arr", Null, true,
		[]Constraint{
			{"Arr", MaxItems, 4, nil},
			{"Arr", MinItems, 2, nil},
		},
	}

	e := Error{
		Constraint:  Null,
		Target:      "Arr",
		TargetValue: x.(Product).Arr,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validateStruct(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_WithArrayEmptyError(t *testing.T) {
	arr := []string{}
	var x interface{} = Product{
		Arr: &[]string{},
	}
	c := Constraint{
		"Arr", Null, true,
		[]Constraint{
			{"Arr", Empty, true, nil},
			{"Arr", MaxItems, 4, nil},
			{"Arr", MinItems, 2, nil},
		}}

	e := Error{
		Constraint:  Empty,
		Target:      "Arr",
		TargetValue: arr,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateStruct(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_WithArrayEmptyWithoutError(t *testing.T) {
	var x interface{} = Product{
		Arr: &[]string{},
	}
	c := Constraint{
		"Arr", Null, true,
		[]Constraint{
			{"Arr", Empty, false, nil},
			{"Arr", MaxItems, 4, nil},
		},
	}
	if z := validateStruct(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: validateStruct failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidateStruct_ArrayWithError(t *testing.T) {
	arr := []string{"1", "1"}
	var x interface{} = Product{
		Arr: &arr,
	}
	c := Constraint{
		"Arr", Null, true,
		[]Constraint{
			{"Arr", Empty, true, nil},
			{"Arr", MaxItems, 4, nil},
			{"Arr", UniqueItems, true, nil},
		},
	}
	s := "Arr"
	e := Error{
		Constraint:  UniqueItems,
		Target:      "Arr",
		TargetValue: arr,
		Details:     fmt.Sprintf("all items in parameter %q must be unique; got:%v", s, arr),
	}
	if z := validateStruct(reflect.ValueOf(x), c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_MapWithError(t *testing.T) {
	m := map[string]string{
		"a": "hello",
		"b": "hello",
	}
	var x interface{} = Product{
		M: &m,
	}
	c := Constraint{
		"M", Null, true,
		[]Constraint{
			{"M", Empty, true, nil},
			{"M", MaxItems, 4, nil},
			{"M", UniqueItems, true, nil},
		},
	}

	s := "M"
	e := Error{
		Constraint:  UniqueItems,
		Target:      "M",
		TargetValue: m,
		Details:     fmt.Sprintf("all items in parameter %q must be unique; got:%v", s, m),
	}

	if z := validateStruct(reflect.ValueOf(x), c); e.Constraint != z.(Error).Constraint ||
		e.Target != z.(Error).Target ||
		!reflect.DeepEqual(e.TargetValue, z.(Error).TargetValue) {
		t.Fatalf("autorest/validation: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}

}

func TestValidateStruct_MapWithNoError(t *testing.T) {
	m := map[string]string{}
	var x interface{} = Product{
		M: &m,
	}
	c := Constraint{
		"M", Null, true,
		[]Constraint{
			{"M", Empty, false, nil},
			{"M", MaxItems, 4, nil},
		},
	}
	if z := validateStruct(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: validateStruct failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidateStruct_MapNilNoError(t *testing.T) {
	var m map[string]string
	var x interface{} = Product{
		M: &m,
	}
	c := Constraint{
		"M", Null, false,
		[]Constraint{
			{"M", Empty, false, nil},
			{"M", MaxItems, 4, nil},
		},
	}
	if z := validateStruct(reflect.ValueOf(x), c); z != nil {
		t.Fatalf("autorest/validation: validateStruct failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_MapValidationWithError(t *testing.T) {
	var x1 interface{} = &Product{
		Arr: &[]string{"1", "2"},
		M:   &map[string]string{"a": "hello"},
	}
	s := "hello"
	var x2 interface{} = &Sample{
		M: &map[string]*string{"a": &s},
	}
	v := []Validation{
		{x1,
			[]Constraint{{"x1", Null, true,
				[]Constraint{
					{"Arr", Null, true,
						[]Constraint{
							{"Arr", Empty, true, nil},
							{"Arr", MaxItems, 4, nil},
							{"Arr", UniqueItems, true, nil},
						},
					},
					{"M", Null, false,
						[]Constraint{
							{"M", Empty, false, nil},
							{"M", MinItems, 1, nil},
							{"M", UniqueItems, true, nil},
						},
					},
				},
			}}},
		{x2,
			[]Constraint{
				{"x2", Null, true,
					[]Constraint{
						{"M", Null, false,
							[]Constraint{
								{"M", Empty, false, nil},
								{"M", MinItems, 2, nil},
								{"M", UniqueItems, true, nil},
							},
						},
					},
				},
				{"Name", Empty, true, nil},
			}},
	}

	e := Error{
		Constraint:  MinItems,
		Target:      "M",
		TargetValue: map[string]*string{"a": &s},
		Details:     fmt.Sprintf("minimum item limit is 2; got: 1"),
	}
	if z := Validate(v); e.Constraint != z.(Error).Constraint ||
		e.Target != z.(Error).Target ||
		!reflect.DeepEqual(e.TargetValue, z.(Error).TargetValue) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidate_MapValidationWithoutError(t *testing.T) {
	var x1 interface{} = &Product{
		Arr: &[]string{"1", "2"},
		M:   &map[string]string{"a": "hello"},
	}
	s := "hello"
	var x2 interface{} = &Sample{
		M: &map[string]*string{"a": &s},
	}
	v := []Validation{
		{x1,
			[]Constraint{{"x1", Null, true,
				[]Constraint{
					{"Arr", Null, true,
						[]Constraint{
							{"Arr", Empty, true, nil},
							{"Arr", MaxItems, 4, nil},
							{"Arr", UniqueItems, true, nil},
						},
					},
					{"M", Null, false,
						[]Constraint{
							{"M", Empty, false, nil},
							{"M", MinItems, 1, nil},
							{"M", UniqueItems, true, nil},
						},
					},
				},
			}}},
		{x2,
			[]Constraint{
				{"x2", Null, true,
					[]Constraint{
						{"M", Null, false,
							[]Constraint{
								{"M", Empty, false, nil},
								{"M", MinItems, 1, nil},
								{"M", UniqueItems, true, nil},
							},
						},
					},
				},
				{"Name", Empty, true, nil},
			}},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect:nil\ngot: %v", z)
	}
}

func TestValidate_UnknownType(t *testing.T) {
	var c chan int
	v := []Validation{
		{c,
			[]Constraint{{"c", Null, true, nil}}},
	}
	e := Error{
		Constraint:  Null,
		Target:      "c",
		TargetValue: c,
		Details:     fmt.Sprintf("unknown type chan"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidate_example1(t *testing.T) {
	var x1 interface{} = Product{
		Arr: &[]string{"1", "1"},
		M:   &map[string]string{"a": "hello"},
	}
	s := "hello"
	var x2 interface{} = Sample{
		M: &map[string]*string{"a": &s},
	}
	v := []Validation{
		{x1,
			[]Constraint{{"Arr", Null, true,
				[]Constraint{
					{"Arr", Empty, true, nil},
					{"Arr", MaxItems, 4, nil},
					{"Arr", UniqueItems, true, nil},
				}},
				{"M", Null, false,
					[]Constraint{
						{"M", Empty, false, nil},
						{"M", MinItems, 1, nil},
						{"M", UniqueItems, true, nil},
					},
				},
			}},
		{x2,
			[]Constraint{
				{"M", Null, false,
					[]Constraint{
						{"M", Empty, false, nil},
						{"M", MinItems, 1, nil},
						{"M", UniqueItems, true, nil},
					},
				},
				{"Name", Empty, true, nil},
			}},
	}
	s = "Arr"
	e := Error{
		Constraint:  UniqueItems,
		Target:      "Arr",
		TargetValue: []string{"1", "1"},
		Details:     fmt.Sprintf("all items in parameter %q must be unique; got:%v", s, []string{"1", "1"}),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidate_Int(t *testing.T) {
	n := int32(100)
	v := []Validation{
		{n,
			[]Constraint{
				{"n", MultipleOf, 10, nil},
				{"n", ExclusiveMinimum, 100, nil},
			},
		},
	}
	e := Error{
		Constraint:  ExclusiveMinimum,
		Target:      "n",
		TargetValue: n,
		Details:     fmt.Sprintf("value must be greater than 100"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidate_IntPointer(t *testing.T) {
	n := int32(100)
	p := &n
	v := []Validation{
		{p,
			[]Constraint{
				{"p", Null, true, []Constraint{
					{"p", ExclusiveMinimum, 100, nil},
				}},
			},
		},
	}
	e := Error{
		Constraint:  ExclusiveMinimum,
		Target:      "p",
		TargetValue: n,
		Details:     fmt.Sprintf("value must be greater than 100"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter
	p = nil
	v = []Validation{
		{p,
			[]Constraint{
				{"p", Null, true, []Constraint{
					{"p", ExclusiveMinimum, 100, nil},
				}},
			},
		},
	}
	e = Error{
		Constraint:  Null,
		Target:      "p",
		TargetValue: p,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	p = nil
	v = []Validation{
		{p,
			[]Constraint{
				{"p", Null, false, []Constraint{
					{"p", ExclusiveMinimum, 100, nil},
				}},
			},
		},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_IntStruct(t *testing.T) {
	n := int32(100)
	p := &Product{
		Num: &n,
	}

	v := []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Num", Null, true, []Constraint{
					{"Num", ExclusiveMinimum, 100, nil},
				}},
			},
		}}},
	}

	e := Error{
		Constraint:  ExclusiveMinimum,
		Target:      "Num",
		TargetValue: n,
		Details:     fmt.Sprintf("value must be greater than 100"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Num", Null, true, []Constraint{
					{"Num", ExclusiveMinimum, 100, nil},
				}},
			},
		}}},
	}

	e = Error{
		Constraint:  Null,
		Target:      "Num",
		TargetValue: p.Num,
		Details:     "value can not be null; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Num", Null, false, []Constraint{
					{"Num", ExclusiveMinimum, 100, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}

	// Parent not required
	p = nil
	v = []Validation{
		{p, []Constraint{{"p", Null, false,
			[]Constraint{
				{"Num", Null, false, []Constraint{
					{"Num", ExclusiveMinimum, 100, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_String(t *testing.T) {
	s := "hello"
	v := []Validation{
		{s,
			[]Constraint{
				{"s", Empty, true, nil},
				{"s", Empty, true,
					[]Constraint{{"s", MaxLength, 3, nil}}},
			},
		},
	}
	e := Error{
		Constraint:  MaxLength,
		Target:      "s",
		TargetValue: s,
		Details:     fmt.Sprintf("string '%s' length must be less than 3", s),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter
	s = ""
	v = []Validation{
		{s,
			[]Constraint{
				{"s", Empty, true, nil},
				{"s", Empty, true,
					[]Constraint{{"s", MaxLength, 3, nil}}},
			},
		},
	}
	e = Error{
		Constraint:  Empty,
		Target:      "s",
		TargetValue: s,
		Details:     "value can not be null or empty; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// not required paramter
	s = ""
	v = []Validation{
		{s,
			[]Constraint{
				{"s", Empty, false, nil},
				{"s", Empty, false,
					[]Constraint{{"s", MaxLength, 3, nil}}},
			},
		},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_StringStruct(t *testing.T) {
	s := "hello"
	p := &Product{
		Str: &s,
	}

	v := []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Str", Null, true, []Constraint{
					{"Str", Empty, true, nil},
					{"Str", MaxLength, 3, nil},
				}},
			},
		}}},
	}
	e := Error{
		Constraint:  MaxLength,
		Target:      "Str",
		TargetValue: s,
		Details:     fmt.Sprintf("string '%s' length must be less than 3", s),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be Empty
	s = ""
	p = &Product{
		Str: &s,
	}
	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Str", Null, true, []Constraint{
					{"Str", Empty, true, nil},
					{"Str", MaxLength, 3, nil},
				}},
			},
		}}},
	}

	e = Error{
		Constraint:  Empty,
		Target:      "Str",
		TargetValue: s,
		Details:     "value can not be null or empty; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be null
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Str", Null, true, []Constraint{
					{"Str", Empty, true, nil},
					{"Str", MaxLength, 3, nil},
				}},
			},
		}}},
	}

	e = Error{
		Constraint:  Null,
		Target:      "Str",
		TargetValue: p.Str,
		Details:     "value can not be null; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Str", Null, false, []Constraint{
					{"Str", Empty, true, nil},
					{"Str", MaxLength, 3, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}

	// Parent not required
	p = nil
	v = []Validation{
		{p, []Constraint{{"p", Null, false,
			[]Constraint{
				{"Str", Null, true, []Constraint{
					{"Str", Empty, true, nil},
					{"Str", MaxLength, 3, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_Array(t *testing.T) {
	s := []string{"hello"}
	v := []Validation{
		{s,
			[]Constraint{
				{"s", Null, true,
					[]Constraint{
						{"s", Empty, true, nil},
						{"s", MinItems, 2, nil},
					}},
			},
		},
	}
	e := Error{
		Constraint:  MinItems,
		Target:      "s",
		TargetValue: s,
		Details:     fmt.Sprintf("minimum item limit is 2; got: %v", len(s)),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Empty array
	v = []Validation{
		{[]string{},
			[]Constraint{
				{"s", Null, true,
					[]Constraint{
						{"s", Empty, true, nil},
						{"s", MinItems, 2, nil}}},
			},
		},
	}
	e = Error{
		Constraint:  Empty,
		Target:      "s",
		TargetValue: []string{},
		Details:     "value can not be null or empty; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// null array
	var s1 []string
	v = []Validation{
		{s1,
			[]Constraint{
				{"s1", Null, true,
					[]Constraint{
						{"s1", Empty, true, nil},
						{"s1", MinItems, 2, nil}}},
			},
		},
	}
	e = Error{
		Constraint:  Null,
		Target:      "s1",
		TargetValue: s1,
		Details:     "value can not be null; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// not required paramter
	v = []Validation{
		{s1,
			[]Constraint{
				{"s1", Null, false,
					[]Constraint{
						{"s1", Empty, true, nil},
						{"s1", MinItems, 2, nil}}},
			},
		},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_ArrayPointer(t *testing.T) {
	s := []string{"hello"}
	v := []Validation{
		{&s,
			[]Constraint{
				{"s", Null, true,
					[]Constraint{
						{"s", Empty, true, nil},
						{"s", MinItems, 2, nil},
					}},
			},
		},
	}
	e := Error{
		Constraint:  MinItems,
		Target:      "s",
		TargetValue: s,
		Details:     fmt.Sprintf("minimum item limit is 2; got: %v", len(s)),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Empty array
	v = []Validation{
		{&[]string{},
			[]Constraint{
				{"s", Null, true,
					[]Constraint{
						{"s", Empty, true, nil},
						{"s", MinItems, 2, nil}}},
			},
		},
	}
	e = Error{
		Constraint:  Empty,
		Target:      "s",
		TargetValue: []string{},
		Details:     "value can not be null or empty; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// null array
	var s1 *[]string
	v = []Validation{
		{s1,
			[]Constraint{
				{"s1", Null, true,
					[]Constraint{
						{"s1", Empty, true, nil},
						{"s1", MinItems, 2, nil}}},
			},
		},
	}
	e = Error{
		Constraint:  Null,
		Target:      "s1",
		TargetValue: s1,
		Details:     "value can not be null; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// not required paramter
	v = []Validation{
		{s1,
			[]Constraint{
				{"s1", Null, false,
					[]Constraint{
						{"s1", Empty, true, nil},
						{"s1", MinItems, 2, nil}}},
			},
		},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_ArrayInStruct(t *testing.T) {
	s := []string{"hello"}
	p := &Product{
		Arr: &s,
	}

	v := []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Arr", Null, true, []Constraint{
					{"Arr", Empty, true, nil},
					{"Arr", MinItems, 2, nil},
				}},
			},
		}}},
	}
	e := Error{
		Constraint:  MinItems,
		Target:      "Arr",
		TargetValue: s,
		Details:     fmt.Sprintf("minimum item limit is 2; got: %v", len(s)),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be Empty
	p = &Product{
		Arr: &[]string{},
	}
	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Arr", Null, true, []Constraint{
					{"Arr", Empty, true, nil},
					{"Arr", MinItems, 2, nil},
				}},
			},
		}}},
	}

	e = Error{
		Constraint:  Empty,
		Target:      "Arr",
		TargetValue: []string{},
		Details:     "value can not be null or empty; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be null
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Arr", Null, true, []Constraint{
					{"Arr", Empty, true, nil},
					{"Arr", MinItems, 2, nil},
				}},
			},
		}}},
	}

	e = Error{
		Constraint:  Null,
		Target:      "Arr",
		TargetValue: p.Arr,
		Details:     "value can not be null; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	v = []Validation{
		{&Product{}, []Constraint{{"p", Null, true,
			[]Constraint{
				{"Arr", Null, false, []Constraint{
					{"Arr", Empty, true, nil},
					{"Arr", MinItems, 2, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}

	// Parent not required
	p = nil
	v = []Validation{
		{p, []Constraint{{"p", Null, false,
			[]Constraint{
				{"Arr", Null, true, []Constraint{
					{"Arr", Empty, true, nil},
					{"Arr", MinItems, 2, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_StructInStruct(t *testing.T) {
	p := &Product{
		C: &Child{I: "hello"},
	}
	v := []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{{"C", Null, true,
				[]Constraint{{"I", MinLength, 7, nil}}},
			},
		}}},
	}
	e := Error{
		Constraint:  MinLength,
		Target:      "I",
		TargetValue: "hello",
		Details:     "string 'hello' length must be greater than 7",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be Empty
	p = &Product{
		C: &Child{I: ""},
	}

	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{{"C", Null, true,
				[]Constraint{{"I", Empty, true, nil}}},
			},
		}}},
	}

	e = Error{
		Constraint:  Empty,
		Target:      "I",
		TargetValue: "",
		Details:     "value can not be null or empty; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be null
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", Null, true,
			[]Constraint{{"C", Null, true,
				[]Constraint{{"I", Empty, true, nil}}},
			},
		}}},
	}

	e = Error{
		Constraint:  Null,
		Target:      "C",
		TargetValue: p.C,
		Details:     "value can not be null; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest/validation: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	v = []Validation{
		{&Product{}, []Constraint{{"p", Null, true,
			[]Constraint{{"C", Null, false,
				[]Constraint{{"I", Empty, true, nil}}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}

	// Parent not required
	p = nil
	v = []Validation{
		{p, []Constraint{{"p", Null, false,
			[]Constraint{{"C", Null, false,
				[]Constraint{{"I", Empty, true, nil}}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest/validation: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestError(t *testing.T) {
	e := Error{
		Constraint:  UniqueItems,
		Target:      "string",
		TargetValue: "go",
		Details:     "minimum length should be 5",
	}
	s := fmt.Sprintf("autorest/validation: validation error: Constraint=%v Parameter=%v Value=%#v Details=%v",
		e.Constraint, e.Target, e.TargetValue, e.Details)
	if !reflect.DeepEqual(e.Error(), s) {
		t.Fatalf("autorest/validation: Error failed to return coorect string \nexpect: %v\ngot: %v", s, e.Error())
	}
}
