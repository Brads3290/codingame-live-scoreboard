package ddb

import (
	"codingame-live-scoreboard/ddb/orm"
	"codingame-live-scoreboard/schema/errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"math"
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
	processedKey := orm.CreateKeyValuesFromMap(keyVals)

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

	err = orm.Unmarshal(res.Item, v)
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
		attrVals[nv] = orm.CreateAttributeValueFromValue(kv)

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
		attrVals[nv] = orm.CreateAttributeValueFromValue(kv)

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

	err = orm.UnmarshalList(qo.Items, v)
	if err != nil {
		return err
	}

	return nil
}

func ScanItemsFromDynamoDb(tbl string, v interface{}) error {
	var sii dynamodb.ScanInput
	sii.SetTableName(tbl)

	sio, err := dynamodbClient.Scan(&sii)
	if err != nil {
		return err
	}

	err = orm.UnmarshalList(sio.Items, v)
	if err != nil {
		return err
	}

	return nil
}

func PopulateItemFromDynamoDb(tbl string, v interface{}) error {

	// Get the keys out of the interface
	keysAv, err := orm.MarshalOnlyKeys(v)
	if err != nil {
		return err
	}

	// Check that each key has a non-empty value
	for k, kv := range keysAv {
		if !orm.AttributeHasNonEmptyValue(kv) {
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

	err = orm.Unmarshal(gio.Item, v)
	if err != nil {
		return err
	}

	return nil
}

func PutItemToDynamoDb(tableName string, v interface{}) error {
	avMap, err := orm.Marshal(v)
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

func BatchPutItemsToDynamoDb(tableName string, v interface{}) error {
	return batchWriteItemsToDynamoDbInternal(tableName, v, bwtPut)
}

func BatchDeleteItemsFromDynamoDb(tableName string, v interface{}) error {
	return batchWriteItemsToDynamoDbInternal(tableName, v, bwtDelete)
}

type batchWriteType int

const (
	bwtPut batchWriteType = iota
	bwtDelete
)

func batchWriteItemsToDynamoDbInternal(tableName string, v interface{}, bwt batchWriteType) error {
	// v must be a slice
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Slice {
		return errors.New("v must be a slice")
	}

	rv := reflect.ValueOf(v)

	// Make sure slice has values
	if rv.Len() == 0 {
		return nil
	}

	// Iterate the slice and create the WriteRequests

	requestItems := make([]*dynamodb.WriteRequest, 0)
	for i := 0; i < rv.Len(); i++ {
		vi := rv.Index(i)

		var wr dynamodb.WriteRequest
		if bwt == bwtPut {
			attrs, err := orm.Marshal(vi.Interface())
			if err != nil {
				return err
			}

			var pr dynamodb.PutRequest
			pr.SetItem(attrs)
			wr.SetPutRequest(&pr)
		} else if bwt == bwtDelete {
			attrs, err := orm.MarshalOnlyKeys(vi.Interface())
			if err != nil {
				return err
			}

			var dr dynamodb.DeleteRequest
			dr.SetKey(attrs)
			wr.SetDeleteRequest(&dr)
		}

		requestItems = append(requestItems, &wr)
	}

	const batchSize = 25
	for i := 0; i < len(requestItems); i += batchSize {
		from := i
		to := int(math.Min(float64(from+batchSize), float64(len(requestItems))))

		requestTables := map[string][]*dynamodb.WriteRequest{
			tableName: requestItems[from:to],
		}

		var bwii dynamodb.BatchWriteItemInput
		bwii.SetRequestItems(requestTables)

		_, err := dynamodbClient.BatchWriteItem(&bwii)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateItemInDynamoDb(tableName string, v interface{}, keyVals map[string]interface{}) error {
	keys := orm.CreateKeyValuesFromMap(keyVals)
	attrs, err := orm.MarshalOnlyAttrs(v)
	if err != nil {
		return err
	}

	var uii dynamodb.UpdateItemInput
	uii.SetTableName(tableName)
	uii.SetKey(keys)
	uii.SetExpressionAttributeValues(attrs)

	_, err = dynamodbClient.UpdateItem(&uii)
	if err != nil {
		return err
	}

	return nil
}

func UpdateItemAttrsInDynamoDb(tableName string, keyVals map[string]interface{}, attrVals map[string]interface{}) error {
	keys := orm.CreateKeyValuesFromMap(keyVals)

	exprAttrNames := make(map[string]*string)
	exprAttrVals := make(map[string]*dynamodb.AttributeValue)

	// Iterate the keys and add them to the attribute name/value lists, and construct the key expression list
	updateExpressionList := make([]string, len(attrVals))
	i := 0
	for k, kv := range attrVals {
		nk := "#ATTR_" + strconv.Itoa(i)
		exprAttrNames[nk] = &k

		nv := ":VAL_" + strconv.Itoa(i)
		exprAttrVals[nv] = orm.CreateAttributeValueFromValue(kv)

		updateExpressionList[i] = fmt.Sprintf("%s = %s", nk, nv)
		i++
	}

	var uii dynamodb.UpdateItemInput
	uii.SetTableName(tableName)
	uii.SetKey(keys)
	uii.SetExpressionAttributeNames(exprAttrNames)
	uii.SetExpressionAttributeValues(exprAttrVals)
	uii.SetUpdateExpression("SET " + strings.Join(updateExpressionList, ", "))

	_, err := dynamodbClient.UpdateItem(&uii)
	if err != nil {
		return err
	}

	return nil
}

func DeleteItemFromDynamoDb(tableName string, v interface{}) error {
	keys, err := orm.MarshalOnlyKeys(v)
	if err != nil {
		return err
	}

	var dii dynamodb.DeleteItemInput
	dii.SetKey(keys)
	dii.SetTableName(tableName)

	_, err = dynamodbClient.DeleteItem(&dii)
	if err != nil {
		return err
	}

	return nil
}
