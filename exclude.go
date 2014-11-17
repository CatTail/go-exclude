package exclude

import (
	"reflect"
)

func Exclude(data, excludes interface{}) {
	excludeInternal(reflect.ValueOf(data), reflect.ValueOf(excludes))
}

func excludeInternal(data, excludes reflect.Value) {
	if data.Type().Kind() == reflect.Interface {
		data = data.Elem()
	}
	if excludes.Type().Kind() == reflect.Interface {
		excludes = excludes.Elem()
	}
	if excludes.MapIndex(reflect.ValueOf("*")).IsValid() {
		for _, key := range data.MapKeys() {
			excludeInternal(data.MapIndex(key), excludes.MapIndex(reflect.ValueOf("*")))
		}
	} else {
		for _, key := range data.MapKeys() {
			value := data.MapIndex(key)
			exclude := excludes.MapIndex(key)
			if exclude.IsValid() {
				if exclude.Type().Kind() == reflect.Interface {
					exclude = exclude.Elem()
				}
				switch exclude.Type().Kind() {
				case reflect.Slice:
					for i := 0; i < exclude.Len(); i++ {
						data.SetMapIndex(key, remove(value, find(value, exclude.Index(i))))
					}
				case reflect.String:
					if exclude.Interface() == "*" {
						data.SetMapIndex(key, reflect.ValueOf(make(map[string][]string)))
					}
				case reflect.Map:
					excludeInternal(value, exclude)
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
	if index < 0 || index > array.Len() {
		return array
	}
	return reflect.AppendSlice(array.Slice(0, index), array.Slice(index+1, array.Len()))
}
