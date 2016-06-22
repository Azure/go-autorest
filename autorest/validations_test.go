package autorest

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCheckForUniqueInArrayTrue(t *testing.T) {
	if c := checkForUniqueInArray([]int{1, 2, 3}); !c {
		t.Fatalf("autorest: checkForUniqueInArray failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInArrayFalse(t *testing.T) {
	if c := checkForUniqueInArray([]int{1, 2, 3, 3}); c {
		t.Fatalf("autorest: checkForUniqueInArray failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInArrayEmpty(t *testing.T) {
	if c := checkForUniqueInArray([]int{}); c {
		t.Fatalf("autorest: checkForUniqueInArray failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInMapTrue(t *testing.T) {
	if c := checkForUniqueInMap(map[string]int{"one": 1, "two": 2}); !c {
		t.Fatalf("autorest: checkForUniqueInMap failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInMapFalse(t *testing.T) {
	if c := checkForUniqueInMap(map[int]string{1: "one", 2: "one"}); c {
		t.Fatalf("autorest: checkForUniqueInMap failed to check unique want: true; got: %v", c)
	}
}

func TestCheckForUniqueInMapEmpty(t *testing.T) {
	if c := checkForUniqueInMap(map[int]string{}); c {
		t.Fatalf("autorest: checkForUniqueInMap failed to check unique want: true; got: %v", c)
	}
}

func TestCheckEmpty_WithValueEmptyRuleTrue(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   true,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}

	if z := checkEmpty(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: checkEmpty failed to check empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_WithValueNilRuleTrue(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   true,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}

	if z := checkEmpty(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: checkEmpty failed to check empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_WithEmptyStringRuleFalse(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   false,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}

	if z := checkEmpty(x, c); z != nil {
		t.Fatalf("autorest: checkEmpty failed to check empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_WithNilRuleFalse(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   false,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := checkEmpty(x, c); z != nil {
		t.Fatalf("autorest: checkEmpty failed to check empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckEmpty_IncorrectRule(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   10,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be bool value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := checkEmpty(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: checkEmpty failed to return error for incorrect rule \nexpect: nil;\ngot: %v", z)
	}
}

func TestCheckEmpty_WithErrorArray(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   true,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}

	if z := checkEmpty(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: checkEmpty failed to check empty parameter \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestCheckNil_WithNilValueRuleTrue(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "x",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{"x", "MaxItems", 4, nil},
		},
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := checkNil(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: checkNil failed to return error for nil value \nexpect: nil;\ngot: %v", z)
	}
}

func TestCheckNil_WithNilValueRuleFalse(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "x",
		Name:   "Null",
		Rule:   false,
		Chain: &[]Constraint{
			{"x", "MaxItems", 4, nil},
		},
	}
	if z := checkNil(x, c); z != nil {
		t.Fatalf("autorest: checkNil failed to return nil value \nexpect: nil;\ngot: %v", z)
	}
}

func TestCheckNil_IncorrectRule(t *testing.T) {
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   "Null",
		Rule:   10,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be bool value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := checkNil(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: checkNil failed to return error for incorrect rule \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithNilValueRuleTrue(t *testing.T) {
	var a []string
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   "Null",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: a,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error for null check \nexpect: nil \ngot: %v", z)
	}
}

func TestValidateArrayMap_WithNilValueRuleFalse(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "arr",
		Name:   "Null",
		Rule:   false,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithValueRuleNullTrue(t *testing.T) {
	var x interface{} = []string{"1", "2"}
	c := Constraint{
		Target: "arr",
		Name:   "Null",
		Rule:   false,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithEmptyValueRuleTrue(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "arr",
		Name:   "Empty",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error for null check \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithEmptyValueRuleFalse(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "arr",
		Name:   "Empty",
		Rule:   false,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_WithEmptyRuleEmptyTrue(t *testing.T) {
	var x interface{} = []string{"1", "2"}
	c := Constraint{
		Target: "arr",
		Name:   "Empty",
		Rule:   false,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MaxItemsIncorrectRule(t *testing.T) {
	var x interface{} = []string{"1", "2"}
	c := Constraint{
		Target: "arr",
		Name:   "MaxItems",
		Rule:   false,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error  \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MaxItemsNoError(t *testing.T) {
	var x interface{} = []string{"1", "2"}
	c := Constraint{
		Target: "arr",
		Name:   "MaxItems",
		Rule:   2,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MaxItemsWithError(t *testing.T) {
	var x interface{} = []string{"1", "2", "3"}
	c := Constraint{
		Target: "arr",
		Name:   "MaxItems",
		Rule:   2,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("maximum item limit is %v; got: 3", c.Rule),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_MaxItemsWithEmpty(t *testing.T) {
	var x interface{} = []string{}
	c := Constraint{
		Target: "arr",
		Name:   "MaxItems",
		Rule:   2,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MinItemsIncorrectRule(t *testing.T) {
	var x interface{} = []int{1, 2}
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   false,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer for %v constraint; got: %v", c.Name, c.Rule),
	}

	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_MinItemsNoError1(t *testing.T) {
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   2,
		Chain:  nil,
	}

	if z := validateArrayMap([]int{1, 2}, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MinItemsNoError2(t *testing.T) {
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   2,
		Chain:  nil,
	}

	if z := validateArrayMap([]int{1, 2, 3}, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_MinItemsWithError(t *testing.T) {
	var x interface{} = []int{1}
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   2,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("mininum item limit is %v; got: 1", c.Rule),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_MinItemsWithEmpty(t *testing.T) {
	var x interface{} = []int{}
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   2,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("mininum item limit is %v; got: 0", c.Rule),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MaxItemsIncorrectRule(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2"}
	c := Constraint{
		Target: "arr",
		Name:   "MaxItems",
		Rule:   false,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MaxItemsNoError(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2"}
	c := Constraint{
		Target: "arr",
		Name:   "MaxItems",
		Rule:   2,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MaxItemsWithError(t *testing.T) {
	a := map[int]string{1: "1", 2: "2", 3: "3"}
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   "MaxItems",
		Rule:   2,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("maximum item limit is %v; got: %v", c.Rule, len(a)),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MaxItemsWithEmpty(t *testing.T) {
	a := map[int]string{}
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   "MaxItems",
		Rule:   2,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MinItemsIncorrectRule(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2"}
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   false,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer for %v constraint; got: %v", c.Name, c.Rule),
	}

	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MinItemsNoError1(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2"}
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   2,
		Chain:  nil,
	}

	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MinItemsNoError2(t *testing.T) {
	var x interface{} = map[int]string{1: "1", 2: "2", 3: "3"}
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   2,
		Chain:  nil,
	}

	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_MinItemsWithError(t *testing.T) {
	a := map[int]string{1: "1"}
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   2,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("mininum item limit is %v; got: %v", c.Rule, len(a)),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_MinItemsWithEmpty(t *testing.T) {
	a := map[int]string{}
	var x interface{} = a
	c := Constraint{
		Target: "arr",
		Name:   "MinItems",
		Rule:   2,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("mininum item limit is %v; got: %v", c.Rule, len(a)),
	}

	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_UniqueItemsTrue(t *testing.T) {
	var x interface{} = map[float64]int{1.2: 1, 1.4: 2}
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Map_UniqueItemsFalse(t *testing.T) {
	var x interface{} = map[string]string{"1": "1", "2": "2", "3": "1"}
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(x, c); e.Target != z.(ValidationError).Target ||
		e.Constraint != z.(ValidationError).Constraint {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_UniqueItemsEmpty(t *testing.T) {
	// Consider empty map as not unique returns false
	var x interface{} = map[int]float64{}
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(x, c); e.Target != z.(ValidationError).Target || e.Constraint != z.(ValidationError).Constraint {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Map_UniqueItemsNil(t *testing.T) {
	// Consider nil map as not unique returns false
	var x interface{} = map[int]float64{}
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(x, c); e.Target != z.(ValidationError).Target || e.Constraint != z.(ValidationError).Constraint {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsTrue(t *testing.T) {
	var x interface{} = []int{1, 2}
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsFalse(t *testing.T) {
	var x interface{} = []string{"1", "2", "1"}
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(x, c); e.Target != z.(ValidationError).Target || e.Constraint != z.(ValidationError).Constraint {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsEmpty(t *testing.T) {
	// Consider empty array as not unique returns false
	var x interface{} = []float64{}
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}
	if z := validateArrayMap(x, c); e.Target != z.(ValidationError).Target || e.Constraint != z.(ValidationError).Constraint {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsNil(t *testing.T) {
	// Consider nil array as not unique returns false
	var a []float64
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, x),
	}

	if z := validateArrayMap(x, c); e.Target != z.(ValidationError).Target || e.Constraint != z.(ValidationError).Constraint {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsInvalidType(t *testing.T) {
	var x interface{} = "hello"
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("type must be array, slice or map for constraint %v; got: %v", c.Name, reflect.ValueOf(x).Kind()),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_Array_UniqueItemsInvalidConstraint(t *testing.T) {
	var x interface{} = "hello"
	c := Constraint{
		Target: "str",
		Name:   "Abc",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("constraint %v is not applicable to array, slice and map type", c.Name),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraint1(t *testing.T) {
	a := []int{1, 2, 3, 4}
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{"str", "MaxItems", 3, nil},
		},
	}
	e := ValidationError{
		Constraint:  (*c.Chain)[0].Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("maximum item limit is %v; got: %v", (*c.Chain)[0].Rule, len(a)),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraint2(t *testing.T) {
	a := []int{1, 2, 3, 4}
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   true,
		Chain: &[]Constraint{
			{"str", "MaxItems", 3, nil},
		},
	}
	e := ValidationError{
		Constraint:  (*c.Chain)[0].Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("maximum item limit is %v; got: %v", (*c.Chain)[0].Rule, len(a)),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraint3(t *testing.T) {
	var a []string
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{"str", "MaxItems", 3, nil},
		},
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraint4(t *testing.T) {
	var x interface{} = []int{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   true,
		Chain: &[]Constraint{
			{"str", "MaxItems", 3, nil},
		},
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateArrayMap(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateArrayMap_ValidateChainConstraintNilNotRequired(t *testing.T) {
	var a []int
	var x interface{} = a
	c := Constraint{
		Target: "str",
		Name:   "Null",
		Rule:   false,
		Chain: &[]Constraint{
			{"str", "MaxItems", 3, nil},
		},
	}

	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateArrayMap_ValidateChainConstraintEmptyNotRequired(t *testing.T) {
	var x interface{} = map[string]int{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   false,
		Chain: &[]Constraint{
			{"str", "MaxItems", 3, nil},
		},
	}

	if z := validateArrayMap(x, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_EmptyTrue(t *testing.T) {
	// empty true means parameter is required but empty returns error
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateString(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateString failed to return error \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_EmptyFalse(t *testing.T) {
	// empty false means parameter is not required and empty return nil
	var x interface{}
	c := Constraint{
		Target: "str",
		Name:   "Empty",
		Rule:   false,
		Chain:  nil,
	}
	if z := validateString(x, c); z != nil {
		t.Fatalf("autorest: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_MaxLengthInvalid(t *testing.T) {
	// empty true means parameter is required but empty returns error
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   "MaxLength",
		Rule:   4,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("string '%s' length must be less than %v", x, c.Rule),
	}
	if z := validateString(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateString failed to return error \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_MaxLengthValid(t *testing.T) {
	// empty false means parameter is not required and empty return nil
	c := Constraint{
		Target: "str",
		Name:   "MaxLength",
		Rule:   7,
		Chain:  nil,
	}
	if z := validateString("Hello", c); z != nil {
		t.Fatalf("autorest: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_MaxLengthRuleInvalid(t *testing.T) {
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   "MaxLength",
		Rule:   true, // must be int for maxLength
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateString(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateString_MinLengthInvalid(t *testing.T) {
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   "MinLength",
		Rule:   10,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("string '%s' length must be greater than %v", x, c.Rule),
	}
	if z := validateString(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateString_MinLengthValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "MinLength",
		Rule:   2,
		Chain:  nil,
	}
	if z := validateString("Hello", c); z != nil {
		t.Fatalf("autorest: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_MinLengthRuleInvalid(t *testing.T) {
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   "MinLength",
		Rule:   true, // must be int for minLength
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateString(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateString_PatternInvalidPattern(t *testing.T) {
	var x interface{} = "Hello"
	c := Constraint{
		Target: "str",
		Name:   "Pattern",
		Rule:   "[[:alnum:",
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     "error parsing regexp: missing closing ]: `[[:alnum:$`",
	}

	if z := validateString(x, c); z.(ValidationError).Details != e.Details {
		t.Fatalf("autorest: validateString failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateString_PatternMatch1(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "Pattern",
		Rule:   "http://\\w+",
		Chain:  nil,
	}
	if z := validateString("http://masd", c); z != nil {
		t.Fatalf("autorest: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_PatternMatch2(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "Pattern",
		Rule:   "[a-zA-Z0-9]+",
		Chain:  nil,
	}
	if z := validateString("asdadad2323sad", c); z != nil {
		t.Fatalf("autorest: validateString failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateString_PatternNotMatch(t *testing.T) {
	var x interface{} = "asdad@@ad2323sad"
	c := Constraint{
		Target: "str",
		Name:   "Pattern",
		Rule:   "[a-zA-Z0-9]+",
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("string '%v' doesn't match pattern %v", x, c.Rule),
	}
	if z := validateString(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateString_InvalidConstraint(t *testing.T) {
	var x interface{} = "asdad@@ad2323sad"
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   "[a-zA-Z0-9]+",
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("constraint %s is not applicable to String type", c.Name),
	}

	if z := validateString(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateString failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_InvalidConstraint(t *testing.T) {
	var x interface{} = 1.4
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   3.0,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("constraint %s is not applicable for type float", c.Name),
	}
	if z := validateFloat(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_InvalidRuleValue(t *testing.T) {
	var x interface{} = 1.4
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMinimum",
		Rule:   3,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be float value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateFloat(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateFloat failed to return nil \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_exclusiveMinimumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMinimum",
		Rule:   1.0,
		Chain:  nil,
	}
	if z := validateFloat(1.42, c); z != nil {
		t.Fatalf("autorest: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_exclusiveMinimumConstraintInvalid(t *testing.T) {
	var x interface{} = 1.4
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMinimum",
		Rule:   1.5,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than %v", c.Rule),
	}
	if z := validateFloat(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_exclusiveMinimumConstraintBoundary(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMinimum",
		Rule:   1.42,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than %v", c.Rule),
	}
	if z := validateFloat(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_exclusiveMaximumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMaximum",
		Rule:   2.0,
		Chain:  nil,
	}
	if z := validateFloat(1.42, c); z != nil {
		t.Fatalf("autorest: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_exclusiveMaximumConstraintInvalid(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMaximum",
		Rule:   1.2,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than %v", c.Rule),
	}
	if z := validateFloat(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_exclusiveMaximumConstraintBoundary(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMaximum",
		Rule:   1.42,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than %v", c.Rule),
	}
	if z := validateFloat(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_inclusiveMaximumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMaximum",
		Rule:   2.0,
		Chain:  nil,
	}
	if z := validateFloat(1.42, c); z != nil {
		t.Fatalf("autorest: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_inclusiveMaximumConstraintInvalid(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMaximum",
		Rule:   1.2,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than or equal to %v", c.Rule),
	}
	if z := validateFloat(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_inclusiveMaximumConstraintBoundary(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMaximum",
		Rule:   1.42,
		Chain:  nil,
	}
	if z := validateFloat(1.42, c); z != nil {
		t.Fatalf("autorest: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_inclusiveMinimumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMinimum",
		Rule:   1.0,
		Chain:  nil,
	}
	if z := validateFloat(1.42, c); z != nil {
		t.Fatalf("autorest: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateFloat_inclusiveMinimumConstraintInvalid(t *testing.T) {
	var x interface{} = 1.42
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMinimum",
		Rule:   1.5,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than or equal to %v", c.Rule),
	}
	if z := validateFloat(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateFloat failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateFloat_inclusiveMinimumConstraintBoundary(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMinimum",
		Rule:   1.42,
		Chain:  nil,
	}
	if z := validateFloat(1.42, c); z != nil {
		t.Fatalf("autorest: valiateFloat failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_InvalidConstraint(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   "UniqueItems",
		Rule:   3,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("constraint %s is not applicable for type integer", c.Name),
	}
	if z := validateInt(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateInt failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateInt_InvalidRuleValue(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMinimum",
		Rule:   3.4,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("rule must be integer value for %v constraint; got: %v", c.Name, c.Rule),
	}
	if z := validateInt(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_exclusiveMinimumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMinimum",
		Rule:   1,
		Chain:  nil,
	}
	if z := validateInt(3, c); z != nil {
		t.Fatalf("autorest: valiateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_exclusiveMinimumConstraintInvalid(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMinimum",
		Rule:   3,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than %v", c.Rule),
	}
	if z := validateInt(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_exclusiveMinimumConstraintBoundary(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMinimum",
		Rule:   1,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than %v", c.Rule),
	}
	if z := validateInt(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateArrayMap failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_exclusiveMaximumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMaximum",
		Rule:   2,
		Chain:  nil,
	}

	if z := validateInt(1, c); z != nil {
		t.Fatalf("autorest: valiateArrayMap failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_exclusiveMaximumConstraintInvalid(t *testing.T) {
	var x interface{} = 2
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMaximum",
		Rule:   1,
		Chain:  nil,
	}

	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than %v", c.Rule),
	}
	if z := validateInt(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_exclusiveMaximumConstraintBoundary(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   "ExclusiveMaximum",
		Rule:   1,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than %v", c.Rule),
	}

	if z := validateInt(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_inclusiveMaximumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMaximum",
		Rule:   2,
		Chain:  nil,
	}
	if z := validateInt(1, c); z != nil {
		t.Fatalf("autorest: validateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_inclusiveMaximumConstraintInvalid(t *testing.T) {
	var x interface{} = 2
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMaximum",
		Rule:   1,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be less than or equal to %v", c.Rule),
	}
	if z := validateInt(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_inclusiveMaximumConstraintBoundary(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMaximum",
		Rule:   1,
		Chain:  nil,
	}
	if z := validateInt(1, c); z != nil {
		t.Fatalf("autorest: validateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_inclusiveMinimumConstraintValid(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMinimum",
		Rule:   1,
		Chain:  nil,
	}

	if z := validateInt(1, c); z != nil {
		t.Fatalf("autorest: valiateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_inclusiveMinimumConstraintInvalid(t *testing.T) {
	var x interface{} = 1
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMinimum",
		Rule:   2,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value must be greater than or equal to %v", c.Rule),
	}
	if z := validateInt(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidateInt_inclusiveMinimumConstraintBoundary(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "InclusiveMinimum",
		Rule:   1,
		Chain:  nil,
	}

	if z := validateInt(1, c); z != nil {
		t.Fatalf("autorest: valiateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_MultipleOfWithoutError(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "MultipleOf",
		Rule:   10,
		Chain:  nil,
	}

	if z := validateInt(2300, c); z != nil {
		t.Fatalf("autorest: valiateInt failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidateInt_MultipleOfWithError(t *testing.T) {
	c := Constraint{
		Target: "str",
		Name:   "MultipleOf",
		Rule:   11,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: 2300,
		Details:     fmt.Sprintf("value must be a multiple of %v", c.Rule),
	}
	if z := validateInt(2300, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiateInt failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_NilTrue(t *testing.T) {
	var z *int
	var x interface{} = z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true, // Required property
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_NilFalse(t *testing.T) {
	var z *int
	var x interface{} = z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   false, // not required property
		Chain:  nil,
	}
	if z := validatePtr(x, c); z != nil {
		t.Fatalf("autorest: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_NilReadonlyValid(t *testing.T) {
	var z *int
	var x interface{} = z
	c := Constraint{
		Target: "ptr",
		Name:   "ReadOnly",
		Rule:   true,
		Chain:  nil,
	}
	if z := validatePtr(x, c); z != nil {
		t.Fatalf("autorest: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_NilReadonlyInvalid(t *testing.T) {
	z := 10
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "ReadOnly",
		Rule:   true,
		Chain:  nil,
	}
	e := ValidationError{
		Constraint:  c.Name,
		Target:      c.Target,
		TargetValue: x,
		Details:     "readonly parameter; must send as nil or empty in request",
	}

	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_IntValid(t *testing.T) {
	z := 10
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "InclusiveMinimum",
		Rule:   3,
		Chain:  nil,
	}
	if z := validatePtr(x, c); z != nil {
		t.Fatalf("autorest: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_IntInvalid(t *testing.T) {
	z := 10
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "InclusiveMinimum",
				Rule:   11,
				Chain:  nil,
			},
		},
	}
	e := ValidationError{
		Constraint:  "InclusiveMinimum",
		Target:      c.Target,
		TargetValue: z,
		Details:     "value must be greater than or equal to 11",
	}

	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}
func TestValidatePointer_IntInvalidConstraint(t *testing.T) {
	z := 10
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "MaxItems",
				Rule:   3,
				Chain:  nil,
			},
		},
	}
	e := ValidationError{
		Constraint:  "MaxItems",
		Target:      c.Target,
		TargetValue: z,
		Details:     "constraint MaxItems is not applicable for type integer",
	}
	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_ValidInt64(t *testing.T) {
	z := int64(10)
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "InclusiveMinimum",
				Rule:   3,
				Chain:  nil,
			},
		}}
	if z := validatePtr(x, c); z != nil {
		t.Fatalf("autorest: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_InvalidConstraintInt64(t *testing.T) {
	z := int64(10)
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "MaxItems",
				Rule:   3,
				Chain:  nil,
			},
		},
	}
	e := ValidationError{
		Constraint:  "MaxItems",
		Target:      c.Target,
		TargetValue: z,
		Details:     "constraint MaxItems is not applicable for type integer",
	}
	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_ValidFloat(t *testing.T) {
	z := 10.1
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "InclusiveMinimum",
				Rule:   3.0,
				Chain:  nil,
			}}}
	if z := validatePtr(x, c); z != nil {
		t.Fatalf("autorest: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_InvalidFloat(t *testing.T) {
	z := 10.1
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "InclusiveMinimum",
				Rule:   12.0,
				Chain:  nil,
			}},
	}
	e := ValidationError{
		Constraint:  "InclusiveMinimum",
		Target:      c.Target,
		TargetValue: z,
		Details:     "value must be greater than or equal to 12",
	}
	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_InvalidConstraintFloat(t *testing.T) {
	z := 10.1
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "MaxItems",
				Rule:   3.0,
				Chain:  nil,
			}},
	}
	e := ValidationError{
		Constraint:  "MaxItems",
		Target:      c.Target,
		TargetValue: z,
		Details:     "constraint MaxItems is not applicable for type float",
	}
	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_StringValid(t *testing.T) {
	z := "hello"
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "Pattern",
				Rule:   "^[a-z]+$",
				Chain:  nil,
			}}}
	if z := validatePtr(x, c); z != nil {
		t.Fatalf("autorest: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_StringInvalid(t *testing.T) {
	z := "hello"
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "MaxLength",
				Rule:   2,
				Chain:  nil,
			}}}

	e := ValidationError{
		Constraint:  "MaxLength",
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("string '%s' length must be less than 2", z),
	}
	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_ArrayValid(t *testing.T) {
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "UniqueItems",
				Rule:   "true",
				Chain:  nil,
			}}}
	if z := validatePtr(&[]string{"1", "2"}, c); z != nil {
		t.Fatalf("autorest: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_ArrayInvalid(t *testing.T) {
	z := []string{"1", "2", "2"}
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{{
			Target: "ptr",
			Name:   "UniqueItems",
			Rule:   true,
			Chain:  nil,
		}},
	}
	e := ValidationError{
		Constraint:  "UniqueItems",
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, z),
	}
	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
	}
}

func TestValidatePointer_MapValid(t *testing.T) {
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{
			{
				Target: "ptr",
				Name:   "UniqueItems",
				Rule:   true,
				Chain:  nil,
			}}}
	if z := validatePtr(&map[interface{}]string{1: "1", "1": "2"}, c); z != nil {
		t.Fatalf("autorest: valiatePtr failed to return nil \nexpect: nil;\ngot: %v", z)
	}
}

func TestValidatePointer_MapInvalid(t *testing.T) {
	z := map[interface{}]string{1: "1", "1": "2", 1.3: "2"}
	var x interface{} = &z
	c := Constraint{
		Target: "ptr",
		Name:   "Null",
		Rule:   true,
		Chain: &[]Constraint{{
			Target: "ptr",
			Name:   "UniqueItems",
			Rule:   true,
			Chain:  nil,
		}},
	}
	e := ValidationError{
		Constraint:  "UniqueItems",
		Target:      c.Target,
		TargetValue: z,
		Details:     fmt.Sprintf("all items in parameter %v must be unique; got:%v", c.Target, z),
	}
	if z := validatePtr(x, c); e.Target != z.(ValidationError).Target ||
		e.Constraint != z.(ValidationError).Constraint || !reflect.DeepEqual(e.TargetValue, z.(ValidationError).TargetValue) {
		t.Fatalf("autorest: valiatePtr failed to return error \nexpect: %v;\ngot: %v", e, z)
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
		"p", "Null", "True",
		&[]Constraint{
			{"C", "Null", true,
				&[]Constraint{
					{"I", maxLength, 2, nil},
				}},
			{"Str", maxLength, 2, nil},
			{"Name", maxLength, 5, nil},
		},
	}
	e := ValidationError{
		Constraint:  "MaxLength",
		Target:      "I",
		TargetValue: "100",
		Details:     fmt.Sprintf("string '100' length must be less than 2"),
	}

	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validatePtr failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidatePointer_WithNilStruct(t *testing.T) {
	var p *Product
	var x interface{} = p
	c := Constraint{
		"p", "Null", true,
		&[]Constraint{
			{"C", "Null", true,
				&[]Constraint{
					{"I", "Empty", true,
						&[]Constraint{
							{"I", maxLength, 5, nil},
						}},
				}},
			{"Str", maxLength, 2, nil},
			{"Name", maxLength, 5, nil},
		},
	}
	e := ValidationError{
		Constraint:  "Null",
		Target:      "p",
		TargetValue: p,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validatePtr(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validatePtr failed to return error \nexpect: %v\ngot: %v", e, z)
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
		"p", "Null", true,
		&[]Constraint{
			{"C", "Null", true,
				&[]Constraint{
					{"I", "Empty", true,
						&[]Constraint{
							{"I", maxLength, 5, nil},
						}},
				}},
		},
	}
	if z := validatePtr(x, c); z != nil {
		t.Fatalf("autorest: validatePtr failed to nil \nexpect: nil\ngot: %v", z)
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
		"C", "Null", true,
		&[]Constraint{
			{"I", "Empty", true,
				&[]Constraint{
					{"I", maxLength, 2, nil},
				}},
		},
	}
	e := ValidationError{
		Constraint:  "MaxLength",
		Target:      "I",
		TargetValue: "100",
		Details:     fmt.Sprintf("string '100' length must be less than 2"),
	}
	if z := validateStruct(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_WithoutChainConstraint(t *testing.T) {
	s := "hello"
	var x interface{} = Product{
		C:    &Child{""},
		Str:  &s,
		Name: "Gopher",
	}
	c := Constraint{"C", "Null", true,
		&[]Constraint{
			{"I", "Empty", true, nil}, // throw error for empty
		}}
	e := ValidationError{
		Constraint:  "Empty",
		Target:      "I",
		TargetValue: "",

		Details: fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateStruct(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
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
	c := Constraint{"Arr", "Null", true,
		&[]Constraint{
			{"Arr", "MaxItems", 4, nil},
			{"Arr", "MinItems", 2, nil},
		},
	}

	e := ValidationError{
		Constraint:  "Null",
		Target:      "Arr",
		TargetValue: x.(Product).Arr,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := validateStruct(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_WithArrayEmptyError(t *testing.T) {
	arr := []string{}
	var x interface{} = Product{
		Arr: &[]string{},
	}
	c := Constraint{
		"Arr", "Null", true,
		&[]Constraint{
			{"Arr", "Empty", true, nil},
			{"Arr", "MaxItems", 4, nil},
			{"Arr", "MinItems", 2, nil},
		}}

	e := ValidationError{
		Constraint:  "Empty",
		Target:      "Arr",
		TargetValue: arr,
		Details:     fmt.Sprintf("value can not be null or empty; required parameter"),
	}
	if z := validateStruct(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidateStruct_WithArrayEmptyWithoutError(t *testing.T) {
	var x interface{} = Product{
		Arr: &[]string{},
	}
	c := Constraint{
		"Arr", "Null", true,
		&[]Constraint{
			{"Arr", "Empty", false, nil},
			{"Arr", "MaxItems", 4, nil},
		},
	}
	if z := validateStruct(x, c); z != nil {
		t.Fatalf("autorest: validateStruct failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidateStruct_ArrayWithError(t *testing.T) {
	arr := []string{"1", "1"}
	var x interface{} = Product{
		Arr: &arr,
	}
	c := Constraint{
		"Arr", "Null", true,
		&[]Constraint{
			{"Arr", "Empty", true, nil},
			{"Arr", "MaxItems", 4, nil},
			{"Arr", "UniqueItems", true, nil},
		},
	}
	e := ValidationError{
		Constraint:  "UniqueItems",
		Target:      "Arr",
		TargetValue: arr,
		Details:     fmt.Sprintf("all items in parameter Arr must be unique; got:%v", arr),
	}
	if z := validateStruct(x, c); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
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
		"M", "Null", true,
		&[]Constraint{
			{"M", "Empty", true, nil},
			{"M", "MaxItems", 4, nil},
			{"M", "UniqueItems", true, nil},
		},
	}

	e := ValidationError{
		Constraint:  "UniqueItems",
		Target:      "M",
		TargetValue: m,
		Details:     fmt.Sprintf("all items in parameter M must be unique; got:%v", m),
	}

	if z := validateStruct(x, c); e.Constraint != z.(ValidationError).Constraint ||
		e.Target != z.(ValidationError).Target ||
		!reflect.DeepEqual(e.TargetValue, z.(ValidationError).TargetValue) {
		t.Fatalf("autorest: validateStruct failed to return error \nexpect: %v\ngot: %v", e, z)
	}

}

func TestValidateStruct_MapWithNoError(t *testing.T) {
	m := map[string]string{}
	var x interface{} = Product{
		M: &m,
	}
	c := Constraint{
		"M", "Null", true,
		&[]Constraint{
			{"M", "Empty", false, nil},
			{"M", "MaxItems", 4, nil},
		},
	}
	if z := validateStruct(x, c); z != nil {
		t.Fatalf("autorest: validateStruct failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidateStruct_MapNilNoError(t *testing.T) {
	var m map[string]string
	var x interface{} = Product{
		M: &m,
	}
	c := Constraint{
		"M", "Null", false,
		&[]Constraint{
			{"M", "Empty", false, nil},
			{"M", "MaxItems", 4, nil},
		},
	}
	if z := validateStruct(x, c); z != nil {
		t.Fatalf("autorest: validateStruct failed to return nil \nexpect: nil\ngot: %v", z)
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
			[]Constraint{{"x1", "Null", true,
				&[]Constraint{
					{"Arr", "Null", true,
						&[]Constraint{
							{"Arr", "Empty", true, nil},
							{"Arr", "MaxItems", 4, nil},
							{"Arr", "UniqueItems", true, nil},
						},
					},
					{"M", "Null", false,
						&[]Constraint{
							{"M", "Empty", false, nil},
							{"M", "MinItems", 1, nil},
							{"M", "UniqueItems", true, nil},
						},
					},
				},
			}}},
		{x2,
			[]Constraint{
				{"x2", "Null", true,
					&[]Constraint{
						{"M", "Null", false,
							&[]Constraint{
								{"M", "Empty", false, nil},
								{"M", "MinItems", 2, nil},
								{"M", "UniqueItems", true, nil},
							},
						},
					},
				},
				{"Name", "Empty", true, nil},
			}},
	}

	e := ValidationError{
		Constraint:  "MinItems",
		Target:      "M",
		TargetValue: map[string]*string{"a": &s},
		Details:     fmt.Sprintf("mininum item limit is 2; got: 1"),
	}
	if z := Validate(v); e.Constraint != z.(ValidationError).Constraint ||
		e.Target != z.(ValidationError).Target ||
		!reflect.DeepEqual(e.TargetValue, z.(ValidationError).TargetValue) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
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
			[]Constraint{{"x1", "Null", true,
				&[]Constraint{
					{"Arr", "Null", true,
						&[]Constraint{
							{"Arr", "Empty", true, nil},
							{"Arr", "MaxItems", 4, nil},
							{"Arr", "UniqueItems", true, nil},
						},
					},
					{"M", "Null", false,
						&[]Constraint{
							{"M", "Empty", false, nil},
							{"M", "MinItems", 1, nil},
							{"M", "UniqueItems", true, nil},
						},
					},
				},
			}}},
		{x2,
			[]Constraint{
				{"x2", "Null", true,
					&[]Constraint{
						{"M", "Null", false,
							&[]Constraint{
								{"M", "Empty", false, nil},
								{"M", "MinItems", 1, nil},
								{"M", "UniqueItems", true, nil},
							},
						},
					},
				},
				{"Name", "Empty", true, nil},
			}},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect:nil\ngot: %v", z)
	}
}

func TestValidate_UnknownType(t *testing.T) {
	var c chan int
	v := []Validation{
		{c,
			[]Constraint{{"c", "Null", true, nil}}},
	}
	e := ValidationError{
		Constraint:  "Null",
		Target:      "c",
		TargetValue: c,
		Details:     fmt.Sprintf("unknown type chan"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
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
			[]Constraint{{"Arr", "Null", true,
				&[]Constraint{
					{"Arr", "Empty", true, nil},
					{"Arr", "MaxItems", 4, nil},
					{"Arr", "UniqueItems", true, nil},
				}},
				{"M", "Null", false,
					&[]Constraint{
						{"M", "Empty", false, nil},
						{"M", "MinItems", 1, nil},
						{"M", "UniqueItems", true, nil},
					},
				},
			}},
		{x2,
			[]Constraint{
				{"M", "Null", false,
					&[]Constraint{
						{"M", "Empty", false, nil},
						{"M", "MinItems", 1, nil},
						{"M", "UniqueItems", true, nil},
					},
				},
				{"Name", "Empty", true, nil},
			}},
	}
	e := ValidationError{
		Constraint:  "UniqueItems",
		Target:      "Arr",
		TargetValue: []string{"1", "1"},
		Details:     fmt.Sprintf("all items in parameter Arr must be unique; got:%v", []string{"1", "1"}),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidate_Int(t *testing.T) {
	n := int32(100)
	v := []Validation{
		{n,
			[]Constraint{
				{"n", "MultipleOf", 10, nil},
				{"n", "ExclusiveMinimum", 100, nil},
			},
		},
	}
	e := ValidationError{
		Constraint:  "ExclusiveMinimum",
		Target:      "n",
		TargetValue: n,
		Details:     fmt.Sprintf("value must be greater than 100"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}
}

func TestValidate_IntPointer(t *testing.T) {
	n := int32(100)
	p := &n
	v := []Validation{
		{p,
			[]Constraint{
				{"p", "Null", true, &[]Constraint{
					{"p", "ExclusiveMinimum", 100, nil},
				}},
			},
		},
	}
	e := ValidationError{
		Constraint:  "ExclusiveMinimum",
		Target:      "p",
		TargetValue: n,
		Details:     fmt.Sprintf("value must be greater than 100"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter
	p = nil
	v = []Validation{
		{p,
			[]Constraint{
				{"p", "Null", true, &[]Constraint{
					{"p", "ExclusiveMinimum", 100, nil},
				}},
			},
		},
	}
	e = ValidationError{
		Constraint:  "Null",
		Target:      "p",
		TargetValue: p,
		Details:     fmt.Sprintf("value can not be null; required parameter"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	p = nil
	v = []Validation{
		{p,
			[]Constraint{
				{"p", "Null", false, &[]Constraint{
					{"p", "ExclusiveMinimum", 100, nil},
				}},
			},
		},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_IntStruct(t *testing.T) {
	n := int32(100)
	p := &Product{
		Num: &n,
	}

	v := []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Num", "Null", true, &[]Constraint{
					{"Num", "ExclusiveMinimum", 100, nil},
				}},
			},
		}}},
	}

	e := ValidationError{
		Constraint:  "ExclusiveMinimum",
		Target:      "Num",
		TargetValue: n,
		Details:     fmt.Sprintf("value must be greater than 100"),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Num", "Null", true, &[]Constraint{
					{"Num", "ExclusiveMinimum", 100, nil},
				}},
			},
		}}},
	}

	e = ValidationError{
		Constraint:  "Null",
		Target:      "Num",
		TargetValue: p.Num,
		Details:     "value can not be null; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Num", "Null", false, &[]Constraint{
					{"Num", "ExclusiveMinimum", 100, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}

	// Parent not required
	p = nil
	v = []Validation{
		{p, []Constraint{{"p", "Null", false,
			&[]Constraint{
				{"Num", "Null", false, &[]Constraint{
					{"Num", "ExclusiveMinimum", 100, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_String(t *testing.T) {
	s := "hello"
	v := []Validation{
		{s,
			[]Constraint{
				{"s", "Empty", true, nil},
				{"s", "Empty", true,
					&[]Constraint{{"s", "MaxLength", 3, nil}}},
			},
		},
	}
	e := ValidationError{
		Constraint:  "MaxLength",
		Target:      "s",
		TargetValue: s,
		Details:     fmt.Sprintf("string '%s' length must be less than 3", s),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter
	s = ""
	v = []Validation{
		{s,
			[]Constraint{
				{"s", "Empty", true, nil},
				{"s", "Empty", true,
					&[]Constraint{{"s", "MaxLength", 3, nil}}},
			},
		},
	}
	e = ValidationError{
		Constraint:  "Empty",
		Target:      "s",
		TargetValue: s,
		Details:     "value can not be null or empty; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// not required paramter
	s = ""
	v = []Validation{
		{s,
			[]Constraint{
				{"s", "Empty", false, nil},
				{"s", "Empty", false,
					&[]Constraint{{"s", "MaxLength", 3, nil}}},
			},
		},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_StringStruct(t *testing.T) {
	s := "hello"
	p := &Product{
		Str: &s,
	}

	v := []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Str", "Null", true, &[]Constraint{
					{"Str", "Empty", true, nil},
					{"Str", "MaxLength", 3, nil},
				}},
			},
		}}},
	}
	e := ValidationError{
		Constraint:  "MaxLength",
		Target:      "Str",
		TargetValue: s,
		Details:     fmt.Sprintf("string '%s' length must be less than 3", s),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be empty
	s = ""
	p = &Product{
		Str: &s,
	}
	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Str", "Null", true, &[]Constraint{
					{"Str", "Empty", true, nil},
					{"Str", "MaxLength", 3, nil},
				}},
			},
		}}},
	}

	e = ValidationError{
		Constraint:  "Empty",
		Target:      "Str",
		TargetValue: s,
		Details:     "value can not be null or empty; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be null
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Str", "Null", true, &[]Constraint{
					{"Str", "Empty", true, nil},
					{"Str", "MaxLength", 3, nil},
				}},
			},
		}}},
	}

	e = ValidationError{
		Constraint:  "Null",
		Target:      "Str",
		TargetValue: p.Str,
		Details:     "value can not be null; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Str", "Null", false, &[]Constraint{
					{"Str", "Empty", true, nil},
					{"Str", "MaxLength", 3, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}

	// Parent not required
	p = nil
	v = []Validation{
		{p, []Constraint{{"p", "Null", false,
			&[]Constraint{
				{"Str", "Null", true, &[]Constraint{
					{"Str", "Empty", true, nil},
					{"Str", "MaxLength", 3, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_Array(t *testing.T) {
	s := []string{"hello"}
	v := []Validation{
		{s,
			[]Constraint{
				{"s", "Null", true,
					&[]Constraint{
						{"s", "Empty", true, nil},
						{"s", "MinItems", 2, nil},
					}},
			},
		},
	}
	e := ValidationError{
		Constraint:  "MinItems",
		Target:      "s",
		TargetValue: s,
		Details:     fmt.Sprintf("mininum item limit is 2; got: %v", len(s)),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// empty array
	v = []Validation{
		{[]string{},
			[]Constraint{
				{"s", "Null", true,
					&[]Constraint{
						{"s", "Empty", true, nil},
						{"s", "MinItems", 2, nil}}},
			},
		},
	}
	e = ValidationError{
		Constraint:  "Empty",
		Target:      "s",
		TargetValue: []string{},
		Details:     "value can not be null or empty; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// null array
	var s1 []string
	v = []Validation{
		{s1,
			[]Constraint{
				{"s1", "Null", true,
					&[]Constraint{
						{"s1", "Empty", true, nil},
						{"s1", "MinItems", 2, nil}}},
			},
		},
	}
	e = ValidationError{
		Constraint:  "Null",
		Target:      "s1",
		TargetValue: s1,
		Details:     "value can not be null; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// not required paramter
	v = []Validation{
		{s1,
			[]Constraint{
				{"s1", "Null", false,
					&[]Constraint{
						{"s1", "Empty", true, nil},
						{"s1", "MinItems", 2, nil}}},
			},
		},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_ArrayPointer(t *testing.T) {
	s := []string{"hello"}
	v := []Validation{
		{&s,
			[]Constraint{
				{"s", "Null", true,
					&[]Constraint{
						{"s", "Empty", true, nil},
						{"s", "MinItems", 2, nil},
					}},
			},
		},
	}
	e := ValidationError{
		Constraint:  "MinItems",
		Target:      "s",
		TargetValue: s,
		Details:     fmt.Sprintf("mininum item limit is 2; got: %v", len(s)),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// empty array
	v = []Validation{
		{&[]string{},
			[]Constraint{
				{"s", "Null", true,
					&[]Constraint{
						{"s", "Empty", true, nil},
						{"s", "MinItems", 2, nil}}},
			},
		},
	}
	e = ValidationError{
		Constraint:  "Empty",
		Target:      "s",
		TargetValue: []string{},
		Details:     "value can not be null or empty; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// null array
	var s1 *[]string
	v = []Validation{
		{s1,
			[]Constraint{
				{"s1", "Null", true,
					&[]Constraint{
						{"s1", "Empty", true, nil},
						{"s1", "MinItems", 2, nil}}},
			},
		},
	}
	e = ValidationError{
		Constraint:  "Null",
		Target:      "s1",
		TargetValue: s1,
		Details:     "value can not be null; required parameter",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// not required paramter
	v = []Validation{
		{s1,
			[]Constraint{
				{"s1", "Null", false,
					&[]Constraint{
						{"s1", "Empty", true, nil},
						{"s1", "MinItems", 2, nil}}},
			},
		},
	}
	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_ArrayInStruct(t *testing.T) {
	s := []string{"hello"}
	p := &Product{
		Arr: &s,
	}

	v := []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Arr", "Null", true, &[]Constraint{
					{"Arr", "Empty", true, nil},
					{"Arr", "MinItems", 2, nil},
				}},
			},
		}}},
	}
	e := ValidationError{
		Constraint:  "MinItems",
		Target:      "Arr",
		TargetValue: s,
		Details:     fmt.Sprintf("mininum item limit is 2; got: %v", len(s)),
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be empty
	p = &Product{
		Arr: &[]string{},
	}
	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Arr", "Null", true, &[]Constraint{
					{"Arr", "Empty", true, nil},
					{"Arr", "MinItems", 2, nil},
				}},
			},
		}}},
	}

	e = ValidationError{
		Constraint:  "Empty",
		Target:      "Arr",
		TargetValue: []string{},
		Details:     "value can not be null or empty; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be null
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Arr", "Null", true, &[]Constraint{
					{"Arr", "Empty", true, nil},
					{"Arr", "MinItems", 2, nil},
				}},
			},
		}}},
	}

	e = ValidationError{
		Constraint:  "Null",
		Target:      "Arr",
		TargetValue: p.Arr,
		Details:     "value can not be null; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	v = []Validation{
		{&Product{}, []Constraint{{"p", "Null", true,
			&[]Constraint{
				{"Arr", "Null", false, &[]Constraint{
					{"Arr", "Empty", true, nil},
					{"Arr", "MinItems", 2, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}

	// Parent not required
	p = nil
	v = []Validation{
		{p, []Constraint{{"p", "Null", false,
			&[]Constraint{
				{"Arr", "Null", true, &[]Constraint{
					{"Arr", "Empty", true, nil},
					{"Arr", "MinItems", 2, nil},
				}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestValidate_StructInStruct(t *testing.T) {
	p := &Product{
		C: &Child{I: "hello"},
	}
	v := []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{{"C", "Null", true,
				&[]Constraint{{"I", "MinLength", 7, nil}}},
			},
		}}},
	}
	e := ValidationError{
		Constraint:  "MinLength",
		Target:      "I",
		TargetValue: "hello",
		Details:     "string 'hello' length must be greater than 7",
	}
	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be empty
	p = &Product{
		C: &Child{I: ""},
	}

	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{{"C", "Null", true,
				&[]Constraint{{"I", "Empty", true, nil}}},
			},
		}}},
	}

	e = ValidationError{
		Constraint:  "Empty",
		Target:      "I",
		TargetValue: "",
		Details:     "value can not be null or empty; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// required paramter - can't be null
	p = &Product{}
	v = []Validation{
		{p, []Constraint{{"p", "Null", true,
			&[]Constraint{{"C", "Null", true,
				&[]Constraint{{"I", "Empty", true, nil}}},
			},
		}}},
	}

	e = ValidationError{
		Constraint:  "Null",
		Target:      "C",
		TargetValue: p.C,
		Details:     "value can not be null; required parameter",
	}

	if z := Validate(v); !reflect.DeepEqual(e, z) {
		t.Fatalf("autorest: Validate failed to return error \nexpect: %v\ngot: %v", e, z)
	}

	// Not required
	v = []Validation{
		{&Product{}, []Constraint{{"p", "Null", true,
			&[]Constraint{{"C", "Null", false,
				&[]Constraint{{"I", "Empty", true, nil}}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}

	// Parent not required
	p = nil
	v = []Validation{
		{p, []Constraint{{"p", "Null", false,
			&[]Constraint{{"C", "Null", false,
				&[]Constraint{{"I", "Empty", true, nil}}},
			},
		}}},
	}

	if z := Validate(v); z != nil {
		t.Fatalf("autorest: Validate failed to return nil \nexpect: nil\ngot: %v", z)
	}
}

func TestError(t *testing.T) {
	e := ValidationError{
		Constraint:  "UniqueItems",
		Target:      "string",
		TargetValue: "go",
		Details:     "minimum length should be 5",
	}
	s := fmt.Sprintf("autorest/azure: validation error: Parameter=%v Value=%v Constraint=%v Details=%v",
		e.Target, e.TargetValue, e.Constraint, e.Details)
	if !reflect.DeepEqual(e.Error(), s) {
		t.Fatalf("autorest: Error failed to return coorect string \nexpect: %v\ngot: %v", s, e.Error())
	}
}
