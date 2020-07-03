package scopes

import (
	"reflect"
	"strings"
)

// StructTag represents the struct tag that will be used for scope calculation.
const StructTag = "scope"

// Marshal takes a scope and an interface and marshals it according to the
// provided scope.
func Marshal(scope Scope, input interface{}) interface{} {
	inputType := reflect.TypeOf(input)
	value := reflect.ValueOf(input)
	if inputType.Kind() == reflect.Ptr {
		inputType = inputType.Elem()
		value = value.Elem()
	}

	if !value.CanInterface() {
		return input
	}

	if inputType.Kind() == reflect.Slice || inputType.Kind() == reflect.Array {
		output := []interface{}{}

		// Iterate through slice/array items and marshal them individually.
		for i := 0; i < value.Len(); i++ {
			item := Marshal(scope, value.Index(i).Interface())
			if item != nil {
				output = append(output, item)
			}
		}

		return output
	}

	if inputType.Kind() == reflect.Map {
		output := map[interface{}]interface{}{}

		iter := value.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			item := Marshal(scope, v.Interface())
			if item != nil {
				output[k] = item
			}
		}

		return output
	}

	if inputType.Kind() != reflect.Struct {
		return input
	}

	output := map[string]interface{}{}

	// Iterate through each field to check the scope.
	for i := 0; i < inputType.NumField(); i++ {
		// Get the scope from the field's struct tag.
		// If the struct tag is empty, the stringToScope function
		// defaults to unauthenticated/always included.
		fieldScope := stringToScope(inputType.Field(i).Tag.Get(StructTag))

		if fieldScope > scope {
			// If the required scope is higher than the provided one,
			// skip the field.
			continue
		}

		// Get the name from the json struct tag.
		name := strings.Split(inputType.Field(i).Tag.Get("json"), ",")[0]
		if name == "-" {
			// Ignore field.
			continue
		}
		if strings.Contains(inputType.Field(i).Tag.Get("json"), "omitempty") && isEmptyValue(value.Field(i)) {
			// Ignore if omitempty present and value is empty.
			continue
		}
		if name == "" {
			// If there is no json tag, use the field name.
			name = inputType.Field(i).Name
		}

		item := value.Field(i).Interface()
		item = Marshal(scope, item)

		// Add the value to the map.
		output[name] = item
	}

	return output
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
