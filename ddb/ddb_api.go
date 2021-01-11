package ddb

import (
	"codingame-live-scoreboard/ddb/ddbmarshal"
	"codingame-live-scoreboard/schema/errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"strconv"
	"strings"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var dynamodbClient = dynamodb.New(sess)

// GetItemFromDynamoDb retrieves a single item from a given table based on key/value pairs given as variadic arguments.
func GetItemFromDynamoDb(tbl string, v interface{}, keyVals map[string]interface{}) error {
	processedKey := ddbmarshal.CreateKeyValuesFromMap(keyVals)

	consistentRead := false
	gii := &dynamodb.GetItemInput{
		ConsistentRead: &consistentRead,
		Key:            processedKey,
		TableName:      &tbl,
	}

	res, err := dynamodbClient.GetItem(gii)
	if err != nil {
		return err
	}

	if res.Item == nil {
		return errors.NewItemNotFound("item not found")
	}

	err = ddbmarshal.Unmarshal(res.Item, v)
	if err != nil {
		return err
	}

	return nil
}

func QueryItemsFromDynamoDb(tbl string, v interface{}, keyVals map[string]interface{}) error {
	return queryItemsFromDynamoDbInternal(tbl, v, keyVals, make(map[string]interface{}))
}

func QueryItemsFromDynamoDbWithFilter(tbl string, v interface{}, keyVals map[string]interface{}, filterVals map[string]interface{}) error {
	return queryItemsFromDynamoDbInternal(tbl, v, keyVals, filterVals)
}

func queryItemsFromDynamoDbInternal(tbl string, v interface{}, keyVals map[string]interface{}, filterVals map[string]interface{}) error {

	// Validate that v is a slice
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Slice {
		return errors.New("QueryItemsFromDynamoDb requires a pointer to a slice as argument 'v'")
	}

	ut := rt.Elem().Elem()
	rv := reflect.ValueOf(v).Elem()

	var qi dynamodb.QueryInput
	qi.SetTableName(tbl)

	// Set up the expression attribute names
	attrNames := make(map[string]*string)
	attrVals := make(map[string]*dynamodb.AttributeValue)

	// Iterate the keys and add them to the attribute name/value lists, and construct the key expression list
	keyConditionExpressionList := make([]string, len(keyVals))
	i := 0
	for k, kv := range keyVals {
		nk := "#KEY_" + strconv.Itoa(i)
		attrNames[nk] = &k

		nv := ":VAL_" + strconv.Itoa(i)
		attrVals[nv] = ddbmarshal.CreateAttributeValueFromValue(kv)

		keyConditionExpressionList[i] = fmt.Sprintf("%s = %s", nk, nv)
		i++
	}

	// Iterate the filters and add them to the attribute name/value lists, and construct the attribute expression list
	filterExpressionList := make([]string, len(filterVals))
	j := 0
	for k, kv := range filterVals {
		nk := "#ATTR_" + strconv.Itoa(j)
		attrNames[nk] = &k

		nv := ":VAL_" + strconv.Itoa(i+j)
		attrVals[nv] = ddbmarshal.CreateAttributeValueFromValue(kv)

		filterExpressionList[j] = fmt.Sprintf("%s = %s", nk, nv)
		j++
	}

	qi.SetExpressionAttributeNames(attrNames)
	qi.SetExpressionAttributeValues(attrVals)
	qi.SetKeyConditionExpression(strings.Join(keyConditionExpressionList, " AND "))

	if len(filterExpressionList) > 0 {
		qi.SetFilterExpression(strings.Join(filterExpressionList, " AND "))
	}

	qo, err := dynamodbClient.Query(&qi)
	if err != nil {
		return err
	}

	// Ensure ut is deferenced
	dut := ut
	if dut.Kind() == reflect.Ptr {
		dut = ut.Elem()
	}

	for _, item := range qo.Items {

		// Create a new object of the underlying slice value
		// Note: reflect.New returns a pointer value
		newObj := reflect.New(dut)

		err = ddbmarshal.Unmarshal(item, newObj.Interface())
		if err != nil {
			return err
		}

		newObjElem := newObj
		if ut.Kind() != reflect.Ptr {
			newObjElem = newObjElem.Elem()
		}

		rv.Set(reflect.Append(rv, newObjElem))
	}

	return nil
}

func PopulateItemFromDynamoDb(tbl string, v interface{}) error {

	// Get the keys out of the interface
	keysAv, err := ddbmarshal.MarshalOnlyKeys(v)
	if err != nil {
		return err
	}

	// Check that each key has a non-empty value
	for k, kv := range keysAv {
		if !ddbmarshal.AttributeHasNonEmptyValue(kv) {
			return errors.New("missing key: " + k)
		}
	}

	var gii dynamodb.GetItemInput
	gii.SetConsistentRead(false)
	gii.SetKey(keysAv)
	gii.SetTableName(tbl)

	gio, err := dynamodbClient.GetItem(&gii)
	if err != nil {
		return err
	}

	if gio.Item == nil {
		return errors.NewItemNotFound("item not found")
	}

	err = ddbmarshal.Unmarshal(gio.Item, v)
	if err != nil {
		return err
	}

	return nil
}

func PutItemToDynamoDb(tableName string, v interface{}) error {
	avMap, err := ddbmarshal.Marshal(v)
	if err != nil {
		return err
	}

	var pii dynamodb.PutItemInput
	pii.SetTableName(tableName)
	pii.SetItem(avMap)

	_, err = dynamodbClient.PutItem(&pii)
	if err != nil {
		return err
	}

	return nil
}

func UpdateItemInDynamoDb(tableName string, v interface{}, keyVals map[string]interface{}) error {
	keys := ddbmarshal.CreateKeyValuesFromMap(keyVals)

	attrValuesToWrite, err := ddbmarshal.MarshalNoKeys(v)
	if err != nil {
		return err
	}

	var uii dynamodb.UpdateItemInput
	uii.SetTableName(tableName)
	uii.SetKey(keys)
	uii.SetExpressionAttributeValues(attrValuesToWrite)

	_, err = dynamodbClient.UpdateItem(&uii)
	if err != nil {
		return err
	}

	return nil
}
