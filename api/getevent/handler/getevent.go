package handler

import (
	"codingame-live-scoreboard/codezone_util"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/schema/dbschema"
	"codingame-live-scoreboard/schema/errors"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// Path params will contain a 'guid' parameter

	return codezone_util.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {
		log := logrus.WithField(constants.API_LOGGER_KEY, "getevent")

		// Check that the url has a {guid}
		guidStr, ok := request.PathParameters["guid"]
		if !ok {
			log.Warn("No guid in path:", request.RawPath)
			return 404, nil, nil
		}

		// Validate the guid by parsing it
		u, err := uuid.Parse(guidStr)
		if err != nil {
			log.Warn("Path guid not valid:", guidStr)
			return 404, nil, nil
		}

		// Create the item to populate
		var evt dbschema.EventModel
		evt.EventId = u.String()

		// Get the item from DynamoDB
		err = codezone_util.PopulateItemFromDynamoDb(constants.DB_TABLE_EVENTS, &evt)

		// If no match found, log but just return nothing
		if errors.IsNotFound(err) {
			log.Warn("No match found in DB for guid: ", u.String())
		}

		if err != nil {
			log.Error("Failed to get item from dynamoDB:", err)
			return 500, nil, err
		}

		return 200, evt, nil
	})
}
