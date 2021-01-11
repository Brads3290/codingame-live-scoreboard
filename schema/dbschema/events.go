package dbschema

import (
	"time"
)

type EventModel struct {
	EventId     string     `ddb:"Event_ID,key" json:"event_id"`
	Name        string     `ddb:"Name" json:"name"`
	LastUpdated *time.Time `ddb:"Last_Updated" json:"last_updated"`
}
