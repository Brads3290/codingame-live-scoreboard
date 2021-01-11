package handler

import (
	"codingame-live-scoreboard/codezone_util"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/schema/dbschema"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"strings"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// request path parameters will contain event_id, with the guid of the event to add the round tp

	return codezone_util.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {

		// get the event_id path param and parse it to validate the guid
		eventIdStr, ok := request.PathParameters["event_id"]
		if !ok {
			return 404, nil, errors.New("missing event_id url parameter")
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			return 404, nil, errors.New("invalid event_id url parameter")
		}

		// Parse the JSON body into a RoundModel object, which we will write
		// to the database with the given event ID (if the eventId is specified
		// in the body, it will be overwritten)
		var r dbschema.RoundModel

		err = json.Unmarshal([]byte(request.Body), &r)
		if err != nil {
			return 401, nil, err
		}

		// Validate r.roundId
		if len(strings.Trim(r.RoundId, " \t\n\r")) == 0 {
			return 401, nil, errors.New("round_id must be specified")
		}

		// Set/overwrite eventId
		r.EventId = eventId.String()

		// Write to the database
		err = codezone_util.PutItemToDynamoDb(constants.DB_TABLE_ROUNDS, r)
		if err != nil {
			return 500, nil, err
		}

		return 200, nil, nil
	})
}
