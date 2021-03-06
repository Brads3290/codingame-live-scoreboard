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
			if _, ok := request.PathParameters["guid"]; ok {
				return getEvent(ctx, request)
			} else {
				return getEvents(ctx, request)
			}
		case "PUT":
			return putEvent(ctx, request)
		case "DELETE":
			return deleteEvent(ctx, request)
		default:
			return 405, nil, errors.New("Unsupported method: " + httpMethod)
		}
	})
}

func getEvent(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
	log := logrus.WithField(constants.API_LOGGER_FIELD, "getevent")
	log.Info(request)

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
	err = ddb.PopulateItemFromDynamoDb(constants.DB_TABLE_EVENTS, &evt)

	// If no match found, log but just return nothing
	if errors.IsNotFound(err) {
		log.Warn("No match found in DB for guid: ", u.String())
	}

	if err != nil {
		log.Error("Failed to get item from dynamoDB:", err)
		return 500, nil, err
	}

	return 200, evt, nil
}

func getEvents(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
	eventList, err := ddb.GetEvents()
	if err != nil {
		return 500, nil, err
	}

	return 200, eventList, nil
}

func putEvent(ctx context.Context, request events.APIGatewayV2HTTPRequest) (sts int, resp interface{}, err error) {
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
}

func deleteEvent(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
	// Get event_id and round_id from path param
	eventIdStr, ok := request.PathParameters["event_id"]
	if !ok {
		return 404, nil, errors.New("missing event_id")
	}

	// Validate eventIdStr
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		return 401, nil, errors.New("invalid event_id")
	}

	err = ddb.DeleteEvent(eventId.String())
	if err != nil {
		return 500, nil, err
	}

	return 200, nil, nil
}
