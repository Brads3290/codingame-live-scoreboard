package codezone_util

import (
	"codingame-live-scoreboard/schema/shared_utils"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type GenericWritable struct {
	data map[string]*dynamodb.AttributeValue
}

func (g *GenericWritable) ToDynamoDbMap() map[string]*dynamodb.AttributeValue {
	return g.data
}

func NewGenericWritable(keyVals ...interface{}) (*GenericWritable, error) {
	d, err := shared_utils.CreateDynamoDbKeyValueMap(keyVals)
	if err != nil {
		return nil, err
	}

	gw := &GenericWritable{
		data: d,
	}

	return gw, nil
}
