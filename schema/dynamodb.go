package schema

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDbReadable interface {
	ToDynamoDbMap() map[string]*dynamodb.AttributeValue
}
