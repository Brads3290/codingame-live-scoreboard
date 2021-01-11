package handler

import (
	"codingame-live-scoreboard/apishared"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/errors"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return apishared.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {

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
	})
}
