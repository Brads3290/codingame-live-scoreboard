package orm

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
)

func Unmarshal(m map[string]*dynamodb.AttributeValue, v interface{}) error {
	rv := reflect.ValueOf(v)

	// As we need to get rv.Elem(), we need to check that rv is Interface or Ptr
	if rv.Kind() != reflect.Ptr && rv.Kind() != reflect.Interface {
		return errors.New("v is not settable")
	}

	rv = rv.Elem()
	rt := rv.Type()

	// Iterate the fields and for each field, try to set it based on the corresponding key in
	// m.
	for i := 0; i < rt.NumField(); i++ {

		// Use the field name and struct tag to find which dynamoDB key it
		// should match to
		ddbName := rt.Field(i).Name
		if tag, ok := rt.Field(i).Tag.Lookup("ddb"); ok {
			dbt := newDdbTag(tag)

			// If the tag is "-", then the field
			// should be skipped
			if dbt.IsIgnored() {
				continue
			}

			// if we have a struct tag on a non-settable field, that's an error
			if !rv.Field(i).CanSet() {
				return errors.New("struct tag on non-settable field: " + rt.Field(i).Name)
			}

			ddbName = dbt.FieldName
		}

		// If the field is non-settable (but does not have a struct tag), simply skip it
		if !rv.Field(i).CanSet() {
			continue
		}

		// Can we find ddbName in the dynamodb map? If not, continue to the next field.
		if attrVal, ok := m[ddbName]; ok {

			// Found the key in the map, so set the field value
			err := setValueFromAttributeValue(attrVal, rv.Field(i), rt.Field(i).Type)
			if err != nil {
				return err
			}
		} else {
			// Explicit continue here in case we decide to add code to the loop,
			// which could introduce a bug otherwise.
			continue
		}
	}

	// Fields have been set, return no error
	return nil
}

func UnmarshalList(m []map[string]*dynamodb.AttributeValue, v interface{}) error {
	// Validate that v is a slice
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Slice {
		return errors.New("QueryItemsFromDynamoDb requires a pointer to a slice as argument 'v'")
	}

	// The underlying type of the slice
	ut := rt.Elem().Elem()

	// The slice value itself (dereferenced)
	rv := reflect.ValueOf(v).Elem()

	// Assign a new slice to rv, so that we replace any existing values and also
	// so if it's nil, we end up with an empty slice rather than nil
	newRv := reflect.MakeSlice(rt.Elem(), 0, 0)
	rv.Set(newRv)

	// Ensure ut is not a pointer by dereferencing it. We can convert back into a pointer
	// later if need be.
	dut := ut
	if dut.Kind() == reflect.Ptr {
		dut = ut.Elem()
	}

	// Iterate the list of items and unmarshal each of them into a new object of the
	// dereferenced list base type.
	for _, item := range m {

		// Create a new object of the underlying slice value
		// Note: reflect.New returns a pointer value
		newObj := reflect.New(dut)

		err := Unmarshal(item, newObj.Interface())
		if err != nil {
			return err
		}

		// If the original underlying type is not a pointer, newObjElem shouldn't be either.
		newObjElem := newObj
		if ut.Kind() != reflect.Ptr {
			newObjElem = newObjElem.Elem()
		}

		rv.Set(reflect.Append(rv, newObjElem))
	}

	return nil
}

func Marshal(v interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return marshalInner(v, kmAll)
}

func MarshalOnlyAttrs(v interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return marshalInner(v, kmNoKeys)
}

func MarshalOnlyKeys(v interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return marshalInner(v, kmOnlyKeys)
}

type keyMode int

const (
	kmNoKeys keyMode = iota
	kmOnlyKeys
	kmAll
)

func marshalInner(v interface{}, km keyMode) (map[string]*dynamodb.AttributeValue, error) {
	m := make(map[string]*dynamodb.AttributeValue)

	rv := reflect.ValueOf(v)

	// Deference rv
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	rt := rv.Type()

	// Iterate the fields in v and create keys in m based on either their names,
	// or their struct tag values
	for i := 0; i < rt.NumField(); i++ {

		// Work out what the key name should be using the struct tag/field name
		ddbName := rt.Field(i).Name
		var dbt ddbTag
		if tag, ok := rt.Field(i).Tag.Lookup("ddb"); ok {
			dbt = newDdbTag(tag)

			// If the tag is "-", skip this field
			if dbt.IsIgnored() {
				continue
			}

			ddbName = dbt.FieldName
		}

		// Skip depending on the keymode and whether or not this is a key
		if km == kmNoKeys && dbt.IsKey() {
			continue
		} else if km == kmOnlyKeys && !dbt.IsKey() {
			continue
		}

		// Get the value from the field
		fieldVal := rv.Field(i).Interface()

		// Convert the value into a *dynamodb.AttributeValue
		av, err := getAttributeValueFromValue(fieldVal)
		if err != nil {
			return nil, err
		}

		// Add it to the map
		m[ddbName] = av
	}

	return m, nil
}
