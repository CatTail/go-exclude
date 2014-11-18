package exclude

import (
	"reflect"
)

func Exclude(input, rules interface{}) {
	excludeInternal(reflect.ValueOf(input), reflect.ValueOf(rules))
}

// input and rules are always maps
func excludeInternal(input, rules reflect.Value) {
	// restore original value by Elm()
	if input.Type().Kind() == reflect.Interface {
		input = input.Elem()
	}
	if rules.Type().Kind() == reflect.Interface {
		rules = rules.Elem()
	}
	if rules.MapIndex(reflect.ValueOf("*")).IsValid() {
		for _, key := range input.MapKeys() {
			excludeInternal(input.MapIndex(key), rules.MapIndex(reflect.ValueOf("*")))
		}
	} else {
		for _, key := range input.MapKeys() {
			value := input.MapIndex(key)
			rule := rules.MapIndex(key)

			if rule.IsValid() {
				// restore original value by Elm()
				if value.Type().Kind() == reflect.Interface {
					value = value.Elem()
				}
				if rule.Type().Kind() == reflect.Interface {
					rule = rule.Elem()
				}

				switch rule.Type().Kind() {
				case reflect.Slice:
					for i := 0; i < rule.Len(); i++ {
						input.SetMapIndex(key, remove(value, find(value, rule.Index(i))))
					}
				case reflect.String:
					if rule.Interface() == "*" {
						input.SetMapIndex(key, reflect.Value{})
					}
				case reflect.Map:
					excludeInternal(value, rule)
				}
			}
		}
	}
}

func find(array, item reflect.Value) (index int) {
	index = -1
	for idx := 0; idx < array.Len(); idx++ {
		if array.Index(idx).Interface() == item.Interface() {
			index = idx
			break
		}
	}
	return
}

func remove(array reflect.Value, index int) reflect.Value {
	if index < 0 || index >= array.Len() {
		return array
	}
	return reflect.AppendSlice(array.Slice(0, index), array.Slice(index+1, array.Len()))
}
