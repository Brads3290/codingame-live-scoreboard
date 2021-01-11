package ddbmarshal

import "reflect"

func valueCanBeNil(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr || v.Kind() == reflect.Chan || v.Kind() == reflect.Func ||
		v.Kind() == reflect.Interface || v.Kind() == reflect.Map || v.Kind() == reflect.Slice
}
