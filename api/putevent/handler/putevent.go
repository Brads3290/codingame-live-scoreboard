package putevent

import (
	"codingame-live-scoreboard/api"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/dbschema"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PUT /putevent
// Body: { "name": "Event_Name" }
func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return api.UnifyLambdaResponse(ctx, func() (sts int, resp interface{}, err error) {
		log := logrus.WithField(constants.API_LOGGER_FIELD, "putevent")

		// Get the name from the request body
		body := struct {
			Name string `json:"name"`
		}{}

		err = json.Unmarshal([]byte(request.Body), &body)
		if err != nil {
			log.Error("Failed to unmarshal JSON body:", err)
			return
		}

		// Create a new event ID
		u, err := uuid.NewRandom()
		if err != nil {
			log.Error("Failed to create new GUID:", err)
			return
		}

		// Create a new event
		model := &dbschema.EventModel{
			EventId:     u.String(),
			Name:        body.Name,
			LastUpdated: nil,
		}

		err = ddb.PutItemToDynamoDb(constants.DB_TABLE_EVENTS, model)
		if err != nil {
			log.Error("Failed to put item to DynamoDB:", err)
			return
		}

		return
	})
}
