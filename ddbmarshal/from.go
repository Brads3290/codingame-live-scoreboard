package ddbmarshal

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"strconv"
	"time"
)

func setValueFromAttributeValue(av *dynamodb.AttributeValue, v reflect.Value) error {
	vi := v.Interface()

	// First check if av is null and, if so, set v to nil.
	// If v cannot be nil, error
	if av.NULL != nil && *av.NULL {

		// If v cannot be nil, we have a problem
		if !valueCanBeNil(v) {
			return errors.New("unable to assign nil to field of type '" + v.Type().Name() + "'")
		}

		// av is NULL and v can be nil, so we assign nil to v
		v.Set(reflect.ValueOf(nil))
		return nil
	}

	// Now we know av is not NULL, so set v in a different way depending on the type
	// of vi (which is just the underlying interface of v)
	var err error
	var val interface{}
	switch vi.(type) {
	case int:
		var i int64
		i, err = intFromAttributeValue(av)
		val = int(i)
	case string:
		val, err = stringFromAttributeValue(av)
	case bool:
		val, err = boolFromAttributeValue(av)
	case time.Time:
		val, err = timeFromAttributeValue(av)
	default:
		return errors.New("unable to assign value to field of unsupported type: '" + v.Type().Name() + "'")
	}

	if err != nil {
		return err
	}

	// We have the value, assign it to v
	v.Set(reflect.ValueOf(val))

	return nil
}

func boolFromAttributeValue(av *dynamodb.AttributeValue) (bool, error) {
	if av.BOOL == nil {
		return false, errors.New("bool value of ddb field is not set")
	}

	return *av.BOOL, nil
}

func intFromAttributeValue(av *dynamodb.AttributeValue) (int64, error) {
	if av.S == nil {
		return 0, errors.New("int value of ddb field is not set")
	}

	i, err := strconv.ParseInt(*av.S, 10, 64)
	if err != nil {
		return 0, errors.New("int value of ddb field is not valid")
	}

	return i, nil
}

func stringFromAttributeValue(av *dynamodb.AttributeValue) (string, error) {
	if av.S == nil {
		return "", errors.New("string value of ddb field is not set")
	}

	return *av.S, nil
}

func timeFromAttributeValue(av *dynamodb.AttributeValue) (time.Time, error) {
	if av.S == nil {
		return time.Time{}, errors.New("time value of ddb field is not set")
	}

	i, err := strconv.ParseInt(*av.S, 10, 64)
	if err != nil {
		return time.Time{}, errors.New("time value of ddb field is not valid")
	}

	return time.Unix(i, 0), nil
}
