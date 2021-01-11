package handler

import (
	"codingame-live-scoreboard/api"
	"codingame-live-scoreboard/api/apishared"
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// The path params should contain an event_id

	return api.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {

		// Check for event_id and validate it by parsing it to a guid
		eventIdStr, ok := request.PathParameters["event_id"]
		if !ok {
			return 404, nil, errors.New("missing event_id")
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			return 404, nil, errors.New("invalid event_id")
		}

		sb, err := apishared.GetApiScoreboardData(eventId.String())
		if err != nil {
			return 500, nil, err
		}

		return 200, sb, nil
	})
}
