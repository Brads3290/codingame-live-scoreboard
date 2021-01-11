package orm

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"strconv"
	"time"
)

func getAttributeValueFromValue(v interface{}) (*dynamodb.AttributeValue, error) {
	av := &dynamodb.AttributeValue{}
	vt := reflect.ValueOf(v)

	// Check if v can be nil. If yes, check if it is nil.
	if valueCanBeNil(vt) {

		// If v is nil, return a nil AttributeValue
		if vt.IsNil() {
			av.SetNULL(true)
			return av, nil
		}
	}

	// Dereference v
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	// The underlying value of v, after dereferencing
	vi := vt.Interface()

	// Depending on what vi is, call a different conversion method
	var err error
	switch vit := vi.(type) {
	case int:
		err = attributeValueFromInt(av, int64(vit))
	case string:
		err = attributeValueFromString(av, vit)
	case bool:
		err = attributeValueFromBool(av, vit)
	case time.Time:
		err = attributeValueFromTime(av, vit)
	default:
		return nil, errors.New("tried to convert unsupported type '" + vt.Type().Name() + "' to AttributeValue")
	}

	if err != nil {
		return nil, err
	}

	return av, nil
}

func attributeValueFromBool(av *dynamodb.AttributeValue, b bool) error {
	av.SetBOOL(b)
	return nil
}

func attributeValueFromInt(av *dynamodb.AttributeValue, i int64) error {
	av.SetN(strconv.FormatInt(i, 10))
	return nil
}

func attributeValueFromString(av *dynamodb.AttributeValue, s string) error {
	av.SetS(s)
	return nil
}

func attributeValueFromTime(av *dynamodb.AttributeValue, t time.Time) error {
	return attributeValueFromInt(av, t.Unix())
}
