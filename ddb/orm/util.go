package orm

import (
	"codingame-live-scoreboard/schema/errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
)

func valueCanBeNil(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr || v.Kind() == reflect.Chan || v.Kind() == reflect.Func ||
		v.Kind() == reflect.Interface || v.Kind() == reflect.Map || v.Kind() == reflect.Slice
}

func typeCanBeNil(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr || t.Kind() == reflect.Chan || t.Kind() == reflect.Func ||
		t.Kind() == reflect.Interface || t.Kind() == reflect.Map || t.Kind() == reflect.Slice
}

func CreateKeyValuesFromList(keyVals []interface{}) map[string]*dynamodb.AttributeValue {
	if len(keyVals)%2 != 0 {
		panic(errors.New("keyVals must be provided as pairs of key/value"))
	}

	keyMap := make(map[string]interface{})
	for i := 0; i < len(keyVals); i += 2 {
		switch kt := keyVals[i].(type) {
		case string:
			keyMap[kt] = keyVals[i+1]
		default:
			panic(errors.New(fmt.Sprintf("key value at position %v is not a string", i)))
		}
	}

	return CreateKeyValuesFromMap(keyMap)
}

func CreateKeyValuesFromMap(keyMap map[string]interface{}) map[string]*dynamodb.AttributeValue {
	processedKey := make(map[string]*dynamodb.AttributeValue)

	for k, v := range keyMap {
		processedKey[k] = CreateAttributeValueFromValue(v)
	}

	return processedKey
}

// CreateAttributeValueFromValue is an exported wrapper for getAttributeValueFromValue
func CreateAttributeValueFromValue(v interface{}) *dynamodb.AttributeValue {
	av, err := getAttributeValueFromValue(v)
	if err != nil {
		panic(err)
	}

	return av
}

// SetAttributeValueFromValue is an exported wrapper for setAttributeValueFromValue
func SetValueFromAttributeValue(av *dynamodb.AttributeValue, v reflect.Value, t reflect.Type) error {
	return setValueFromAttributeValue(av, v, t)
}

func AttributeHasNonEmptyValue(av *dynamodb.AttributeValue) bool {
	if av.BOOL != nil {
		return true
	}

	if av.N != nil && len(*av.N) > 0 {
		return true
	}

	if av.S != nil && len(*av.S) > 0 {
		return true
	}

	return false
}
