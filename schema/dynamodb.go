package schema

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDbWritable interface {
	ToDynamoDbMap() map[string]*dynamodb.AttributeValue
}
