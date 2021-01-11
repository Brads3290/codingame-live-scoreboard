package handler

import (
	"codingame-live-scoreboard/api"
	"codingame-live-scoreboard/api/apishared"
	"codingame-live-scoreboard/codingame"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/errors"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// The path params will have an event_id

	return api.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {
		log := logrus.WithField(constants.API_LOGGER_FIELD, "updateevent")

		// get the event_id path param and parse it to validate the guid
		eventIdStr, ok := request.PathParameters["event_id"]
		if !ok {
			return 404, nil, errors.New("missing event_id url parameter")
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			return 404, nil, errors.New("invalid event_id url parameter")
		}

		// Check if the event needs updating
		evt, err := ddb.GetEvent(eventId.String())
		if errors.IsNotFound(err) {
			return 404, nil, nil
		} else if err != nil {
			return 500, nil, err
		}

		if !evt.IsUpdating && (evt.LastUpdated == nil || time.Now().Sub(*evt.LastUpdated) > constants.MAX_EVENT_RECORD_AGE) {

			// Set Is_Updating to attempt to stop other requests from performing the same update
			err = ddb.SetEventUpdating(eventId.String(), true)
			if err != nil {
				return 500, nil, err
			}

			defer func() {
				err = ddb.SetEventUpdating(eventId.String(), false)
				if err != nil {
					log.Error("Failed to set event updating to false", err)
				}
			}()

			sd, err := codingame.GetCodinGameData(eventId.String())
			if err != nil {
				log.Error("Failed to get codingame data: ", err)
				return 500, nil, err
			}

			err = ddb.UpdateDynamoDbFromScoreData(eventId.String(), sd)
			if err != nil {
				log.Error("Failed to update dynamodb with codingame data: ", err)
				return 500, nil, err
			}
		}

		sb, err := apishared.GetApiScoreboardData(eventId.String())
		if err != nil {
			return 500, nil, err
		}

		return 200, sb, nil
	})
}
