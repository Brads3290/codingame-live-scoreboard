package dbschema

import (
	"codingame-live-scoreboard/schema/shared_utils"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
)

type EventModel struct {
	EventId     string     `ddb:"Event_ID,key" json:"event_id"`
	Name        string     `ddb:"Name" json:"name"`
	LastUpdated *time.Time `ddb:"Last_Updated" json:"last_updated"`
}

func (e *EventModel) ToDynamoDbMap() map[string]*dynamodb.AttributeValue {
	m, err := shared_utils.CreateKeyValuesFromList(
		"Event_ID", e.EventId,
		"Name", e.Name,
		"Last_Updated", e.LastUpdated,
	)

	if err != nil {
		panic(err)
	}

	return m
}
