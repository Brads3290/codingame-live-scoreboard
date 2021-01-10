package shared_utils

import (
	"codingame-live-scoreboard/schema/errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
	"time"
)

func CreateDynamoDbKeyValueMap(keyVals ...interface{}) (map[string]*dynamodb.AttributeValue, error) {
	if len(keyVals)%2 != 0 {
		return nil, errors.New("keyVals must be provided as pairs of key/value")
	}

	processedKey := make(map[string]*dynamodb.AttributeValue)

	keyMap := make(map[string]interface{})
	for i := 0; i < len(keyVals); i += 2 {
		switch kt := keyVals[i].(type) {
		case string:
			keyMap[kt] = keyVals[i+1]
		default:
			return nil, errors.New(fmt.Sprintf("key value at position %v is not a string", i))
		}
	}

	for k, v := range keyMap {
		var a dynamodb.AttributeValue

		switch vt := v.(type) {
		case string:
			a.SetS(vt)
			break
		case bool:
			a.SetBOOL(vt)
			break
		case int:
			a.SetS(strconv.Itoa(vt))
			break
		case *time.Time:
			if vt == nil {
				a.SetNULL(true)
			} else {
				t := vt.Unix()
				i := strconv.FormatInt(t, 10)
				a.SetS(i)
			}
		case time.Time:
			t := vt.Unix()
			i := strconv.FormatInt(t, 10)
			a.SetS(i)
		default:
			return nil, errors.New("unsupported partKey type")
		}

		processedKey[k] = &a
	}

	return processedKey, nil
}
