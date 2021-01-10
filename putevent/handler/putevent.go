package handler

import (
	"codingame-live-scoreboard/codezone_util"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/schema"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

// PUT /putevent
// Body: { "name": "Event_Name" }
func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return codezone_util.UnifyLambdaResponse(ctx, func() (sts int, resp interface{}, err error) {

		// Get the name from the request body
		body := struct {
			Name string `json:"name"`
		}{}

		err = json.Unmarshal([]byte(request.Body), &body)
		if err != nil {
			return
		}

		// Create a new event ID
		u, err := uuid.NewRandom()
		if err != nil {
			return
		}

		// Create a new event
		model := &schema.EventModel{
			EventId:     u.String(),
			Name:        body.Name,
			LastUpdated: nil,
		}

		err = codezone_util.PutItemToDynamoDb(constants.DB_TABLE_EVENTS, model)
		if err != nil {
			return
		}

		return
	})
}
