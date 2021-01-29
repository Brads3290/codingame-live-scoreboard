package handler

import (
	"codingame-live-scoreboard/apishared"
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"strings"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// The path params should contain an event_id

	return apishared.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {
		httpMethod := strings.ToUpper(request.RequestContext.HTTP.Method)

		switch httpMethod {
		case "GET":
			return getScoreboard(ctx, request)
		default:
			return 405, nil, errors.New("Unsupported method: " + httpMethod)
		}
	})
}

func getScoreboard(ctx context.Context, request events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
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
}
