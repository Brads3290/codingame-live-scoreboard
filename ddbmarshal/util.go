package ddbmarshal

import "reflect"

func valueCanBeNil(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr || v.Kind() == reflect.Chan || v.Kind() == reflect.Func ||
		v.Kind() == reflect.Interface || v.Kind() == reflect.Map || v.Kind() == reflect.Slice
}

func typeCanBeNil(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr || t.Kind() == reflect.Chan || t.Kind() == reflect.Func ||
		t.Kind() == reflect.Interface || t.Kind() == reflect.Map || t.Kind() == reflect.Slice
}
