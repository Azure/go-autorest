package autorest

import (
	"fmt"
	"reflect"
	"regexp"
)

// Constraint stroes constraint name, target field name
// Rule and chain validations.
type Constraint struct {

	// Target Field name for validation.
	Target string

	// Constraint name .e.g. Minlength, Maxlength, Pattern, etc.
	Name string

	// Rule for constraint .e.g. greater than 10, less than 5 etc.
	Rule interface{}

	// Chain Validations for struct type
	Chain *[]Constraint
}

// Validation stores parameter-wise validation.
type Validation struct {
	TargetValue interface{}
	Constraints []Constraint
}

const (
	empty            = "Empty"
	null             = "Null"
	readOnly         = "ReadOnly"
	pattern          = "Pattern"
	maxLength        = "MaxLength"
	minLength        = "MinLength"
	maxItems         = "MaxItems"
	minItems         = "MinItems"
	multipleOf       = "MultipleOf"
	uniqueItems      = "UniqueItems"
	inclusiveMaximum = "InclusiveMaximum"
	inclusiveMinimum = "InclusiveMinimum"
	exclusiveMaximum = "ExclusiveMaximum"
	exclusiveMinimum = "ExclusiveMinimum"
)

// Validate method validates constraints on parameter
// passed in validation array.
func Validate(m []Validation) error {
	for _, item := range m {
		v := reflect.ValueOf(item.TargetValue)
		for _, constraint := range item.Constraints {
			var err error
			switch v.Kind() {
			case reflect.Ptr:
				err = validatePtr(item.TargetValue, constraint)
			case reflect.String:
				err = validateString(item.TargetValue, constraint)
			case reflect.Struct:
				err = validateStruct(item.TargetValue, constraint)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				err = validateInt(item.TargetValue, constraint)
			case reflect.Float32, reflect.Float64:
				err = validateFloat(item.TargetValue, constraint)
			case reflect.Array, reflect.Slice, reflect.Map:
				err = validateArrayMap(item.TargetValue, constraint)
			default:
				err = createError(item.TargetValue, constraint, fmt.Sprintf("unknown type %v", v.Kind()))
			}

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func validateStruct(x interface{}, v Constraint) error {
	z := reflect.ValueOf(x)
	f := z.FieldByName(v.Target)
	if err := Validate([]Validation{
		{
			TargetValue: f.Interface(),
			Constraints: []Constraint{v},
		},
	}); err != nil {
		return err
	}
	return nil
}

func validatePtr(x interface{}, v Constraint) error {
	z := reflect.ValueOf(x)
	if v.Name == readOnly {
		if !z.IsNil() {
			return createError(x, v, "readonly parameter; must send as nil or empty in request")
		}
		return nil
	}
	if x == nil || z.IsNil() {
		return checkNil(x, v)
	}
	if v.Chain != nil {
		return Validate([]Validation{
			{
				TargetValue: z.Elem().Interface(),
				Constraints: *v.Chain,
			},
		})
	}
	return nil
}

func validateInt(x interface{}, v Constraint) error {
	i := reflect.ValueOf(x).Int()
	r, ok := v.Rule.(int)
	if !ok {
		return createError(x, v, fmt.Sprintf("rule must be integer value for %v constraint; got: %v", v.Name, v.Rule))
	}
	switch v.Name {
	case multipleOf:
		if i%int64(r) != 0 {
			return createError(x, v, fmt.Sprintf("value must be a multiple of %v", r))
		}
	case exclusiveMinimum:
		if i <= int64(r) {
			return createError(x, v, fmt.Sprintf("value must be greater than %v", r))
		}
	case exclusiveMaximum:
		if i >= int64(r) {
			return createError(x, v, fmt.Sprintf("value must be less than %v", r))
		}
	case inclusiveMinimum:
		if i < int64(r) {
			return createError(x, v, fmt.Sprintf("value must be greater than or equal to %v", r))
		}
	case inclusiveMaximum:
		if i > int64(r) {
			return createError(x, v, fmt.Sprintf("value must be less than or equal to %v", r))
		}
	default:
		return createError(x, v, fmt.Sprintf("constraint %s is not applicable for type integer", v.Name))
	}
	return nil
}

func validateFloat(x interface{}, v Constraint) error {
	f := reflect.ValueOf(x).Float()
	r, ok := v.Rule.(float64)
	if !ok {
		return createError(x, v, fmt.Sprintf("rule must be float value for %v constraint; got: %v", v.Name, v.Rule))
	}
	switch v.Name {
	case exclusiveMinimum:
		if f <= r {
			return createError(x, v, fmt.Sprintf("value must be greater than %v", r))
		}
	case exclusiveMaximum:
		if f >= r {
			return createError(x, v, fmt.Sprintf("value must be less than %v", r))
		}
	case inclusiveMinimum:
		if f < r {
			return createError(x, v, fmt.Sprintf("value must be greater than or equal to %v", r))
		}
	case inclusiveMaximum:
		if f > r {
			return createError(x, v, fmt.Sprintf("value must be less than or equal to %v", r))
		}
	default:
		return createError(x, v, fmt.Sprintf("constraint %s is not applicable for type float", v.Name))
	}
	return nil
}

func validateString(x interface{}, v Constraint) error {
	s := reflect.ValueOf(x).String()
	switch v.Name {
	case empty:
		if x == nil || len(s) == 0 {
			return checkEmpty(x, v)
		}
	case pattern:
		reg, err := regexp.Compile("^" + v.Rule.(string) + "$")
		if err != nil {
			return createError(x, v, err.Error())
		}
		if !reg.MatchString(s) {
			return createError(x, v, fmt.Sprintf("string '%s' doesn't match pattern %v", s, v.Rule))
		}
	case maxLength:
		if _, ok := v.Rule.(int); !ok {
			return createError(x, v, fmt.Sprintf("rule must be integer value for %v constraint; got: %v", v.Name, v.Rule))
		}
		if len(s) > v.Rule.(int) {
			return createError(x, v, fmt.Sprintf("string '%s' length must be less than %v", s, v.Rule))
		}
	case minLength:
		if _, ok := v.Rule.(int); !ok {
			return createError(x, v, fmt.Sprintf("rule must be integer value for %v constraint; got: %v", v.Name, v.Rule))
		}
		if len(s) < v.Rule.(int) {
			return createError(x, v, fmt.Sprintf("string '%s' length must be greater than %v", s, v.Rule))
		}
	default:
		return createError(x, v, fmt.Sprintf("constraint %s is not applicable to String type", v.Name))
	}

	if v.Chain != nil {
		return Validate([]Validation{
			{
				TargetValue: x,
				Constraints: *v.Chain,
			},
		})
	}
	return nil
}

func validateArrayMap(x interface{}, v Constraint) error {
	z := reflect.ValueOf(x)
	switch v.Name {
	case null:
		if z.IsNil() {
			return checkNil(x, v)
		}
	case empty:
		if x == nil || z.Len() == 0 {
			return checkEmpty(x, v)
		}
	case maxItems:
		if _, ok := v.Rule.(int); !ok {
			return createError(x, v, fmt.Sprintf("rule must be integer for %v constraint; got: %v", v.Name, v.Rule))
		}
		if z.Len() > v.Rule.(int) {
			return createError(x, v, fmt.Sprintf("maximum item limit is %v; got: %v", v.Rule, z.Len()))
		}
	case minItems:
		if _, ok := v.Rule.(int); !ok {
			return createError(x, v, fmt.Sprintf("rule must be integer for %v constraint; got: %v", v.Name, v.Rule))
		}
		if z.Len() < v.Rule.(int) {
			return createError(x, v, fmt.Sprintf("mininum item limit is %v; got: %v", v.Rule, z.Len()))
		}
	case uniqueItems:
		if z.Kind() == reflect.Array || z.Kind() == reflect.Slice {
			if !checkForUniqueInArray(x) {
				return createError(x, v, fmt.Sprintf("all items in parameter %v must be unique; got:%v", v.Target, z))
			}
		} else if z.Kind() == reflect.Map {
			if !checkForUniqueInMap(x) {
				return createError(x, v, fmt.Sprintf("all items in parameter %v must be unique; got:%v", v.Target, z))
			}
		} else {
			return createError(x, v, fmt.Sprintf("type must be array, slice or map for constraint %v; got: %v", v.Name, z.Kind()))
		}
	default:
		return createError(x, v, fmt.Sprintf("constraint %v is not applicable to array, slice and map type", v.Name))
	}

	if v.Chain != nil {
		return Validate([]Validation{
			{
				TargetValue: x,
				Constraints: *v.Chain,
			},
		})
	}
	return nil
}

func checkNil(x interface{}, v Constraint) error {
	if _, ok := v.Rule.(bool); !ok {
		return createError(x, v, fmt.Sprintf("rule must be bool value for %v constraint; got: %v", v.Name, v.Rule))
	}
	if v.Rule.(bool) {
		return createError(x, v, "value can not be null; required parameter")
	}
	return nil
}

func checkEmpty(x interface{}, v Constraint) error {
	if _, ok := v.Rule.(bool); !ok {
		return createError(x, v, fmt.Sprintf("rule must be bool value for %v constraint; got: %v", v.Name, v.Rule))
	}

	if v.Rule.(bool) {
		return createError(x, v, "value can not be null or empty; required parameter")
	}
	return nil
}

func checkForUniqueInArray(x interface{}) bool {
	z := reflect.ValueOf(x)
	if x == nil || z.Len() == 0 {
		return false
	}
	arrOfInterface := make([]interface{}, z.Len())

	for i := 0; i < z.Len(); i++ {
		arrOfInterface[i] = z.Index(i).Interface()
	}

	m := make(map[interface{}]bool)
	for _, val := range arrOfInterface {
		if m[val] {
			return false
		}
		m[val] = true
	}
	return true
}

func checkForUniqueInMap(x interface{}) bool {
	z := reflect.ValueOf(x)
	if x == nil || z.Len() == 0 {
		return false
	}
	mapOfInterface := make(map[interface{}]interface{}, z.Len())

	keys := z.MapKeys()
	for _, k := range keys {
		mapOfInterface[k.Interface()] = z.MapIndex(k).Interface()
	}

	m := make(map[interface{}]bool)
	for _, val := range mapOfInterface {
		if m[val] {
			return false
		}
		m[val] = true
	}
	return true
}

// ValidationError stroes detailed validation error
type ValidationError struct {
	Constraint  string
	Target      string
	TargetValue interface{}
	Details     string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("autorest/azure: validation error: Parameter=%v Value=%v Constraint=%v Details=%v",
		e.Target, e.TargetValue, e.Constraint, e.Details)
}

func createError(x interface{}, v Constraint, err string) ValidationError {
	return ValidationError{
		Constraint:  v.Name,
		Target:      v.Target,
		TargetValue: x,
		Details:     err,
	}
}
