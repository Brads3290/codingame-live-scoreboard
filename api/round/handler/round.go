package handler

import (
	"codingame-live-scoreboard/apishared"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/dbschema"
	"codingame-live-scoreboard/schema/errors"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return apishared.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {
		httpMethod := strings.ToUpper(request.RequestContext.HTTP.Method)

		switch httpMethod {
		case "GET":
			return getRounds(ctx, request)
		case "PUT":
			return putRound(ctx, request)
		case "DELETE":
			return deleteRound(ctx, request)
		default:
			return 405, nil, errors.New("Unsupported method: " + httpMethod)
		}
	})
}

func getRounds(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
	log := logrus.WithField(constants.API_LOGGER_FIELD, "getrounds")

	// Check that the url has a {guid}
	guidStr, ok := request.PathParameters["event_id"]
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

	rounds, err := ddb.GetAllRoundsForEvent(u.String())
	if err != nil {
		return 500, nil, err
	}

	return 200, rounds, nil
}

func deleteRound(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
	// Get event_id and round_id from path param
	eventIdStr, ok := request.PathParameters["event_id"]
	if !ok {
		return 404, nil, errors.New("missing event_id")
	}

	roundId, ok := request.PathParameters["round_id"]
	if !ok {
		return 404, nil, errors.New("missing round_id")
	}

	// Validate eventIdStr
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		return 401, nil, errors.New("invalid event_id")
	}

	err = ddb.DeleteRound(eventId.String(), roundId)
	if err != nil {
		return 500, nil, err
	}

	return 200, nil, nil
}

func putRound(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
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

	// Set/overwrite eventId; the round also starts active
	r.EventId = eventId.String()
	r.Active = true

	// Write to the database
	err = ddb.PutItemToDynamoDb(constants.DB_TABLE_ROUNDS, r)
	if err != nil {
		return 500, nil, err
	}

	return 200, nil, nil
}
