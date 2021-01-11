package orm

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"strconv"
	"time"
)

func setValueFromAttributeValue(av *dynamodb.AttributeValue, v reflect.Value, t reflect.Type) error {

	// First check if av is null and, if so, set v to nil.
	// If v cannot be nil, error
	if av.NULL != nil && *av.NULL {

		// If v cannot be nil, we have a problem
		if !typeCanBeNil(t) {
			return errors.New("unable to assign nil to field of type '" + t.String() + "'")
		}

		// av is NULL and v can be nil, so we assign nil to v
		v.Set(reflect.Zero(t))
		return nil
	}

	// t may be a pointer, so get the underlying type
	ut := t
	if t.Kind() == reflect.Ptr {
		ut = t.Elem()
	}

	// Now we know av is not NULL, so set v in a different way depending on the type
	// of vi (which is just the underlying interface of v)
	var err error
	var val interface{}
	switch true {
	case reflect.TypeOf((*int)(nil)).Elem().AssignableTo(ut):
		var i int64
		i, err = intFromAttributeValue(av)
		val = int(i)
	case reflect.TypeOf((*string)(nil)).Elem().AssignableTo(ut):
		val, err = stringFromAttributeValue(av)
	case reflect.TypeOf((*bool)(nil)).Elem().AssignableTo(ut):
		val, err = boolFromAttributeValue(av)
	case reflect.TypeOf((*time.Time)(nil)).Elem().AssignableTo(ut):
		val, err = timeFromAttributeValue(av)
	default:
		return errors.New("unable to assign value to field of unsupported type: '" + t.String() + "'")
	}

	if err != nil {
		return err
	}

	vVal := reflect.ValueOf(val)

	// Sanity check: If val is nil, and t cannot hold a nil value, that's an error.
	if valueCanBeNil(vVal) && vVal.IsNil() && !typeCanBeNil(t) {
		return errors.New("cannot assign null value to type '" + t.String() + "'")
	}

	// We have the value, assign it to v. If t is a pointer, we need to assign the underlying value
	// unless val is also a pointer.
	//
	// Also, if val is a pointer and t is not, deference val before assigning
	if t.Kind() == reflect.Ptr && reflect.TypeOf(val).Kind() != reflect.Ptr {
		v.Elem().Set(vVal)
	} else if t.Kind() != reflect.Ptr && reflect.TypeOf(val).Kind() == reflect.Ptr {
		v.Set(vVal.Elem())
	} else {
		v.Set(vVal)
	}

	return nil
}

func boolFromAttributeValue(av *dynamodb.AttributeValue) (bool, error) {
	if av.BOOL == nil {
		return false, errors.New("bool value of ddb field is not set")
	}

	return *av.BOOL, nil
}

func intFromAttributeValue(av *dynamodb.AttributeValue) (int64, error) {
	if av.N == nil {
		return 0, errors.New("int value of ddb field is not set")
	}

	i, err := strconv.ParseInt(*av.N, 10, 64)
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

func timeFromAttributeValue(av *dynamodb.AttributeValue) (*time.Time, error) {

	// If the attribute has no value, we return nil
	if !AttributeHasNonEmptyValue(av) {
		return nil, nil
	}

	if av.N == nil {
		return nil, errors.New("time value of ddb field is not set")
	}

	i, err := strconv.ParseInt(*av.N, 10, 64)
	if err != nil {
		return nil, errors.New("time value of ddb field is not valid")
	}

	t := time.Unix(i, 0)
	return &t, nil
}
