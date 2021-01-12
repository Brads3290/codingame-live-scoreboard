package handler

import (
	"codingame-live-scoreboard/apishared"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/ddb"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Expects /{event_id}
func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return apishared.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {
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
	})
}
