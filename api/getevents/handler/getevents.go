package handler

import (
	"codingame-live-scoreboard/apishared"
	"codingame-live-scoreboard/ddb"
	"context"
	"github.com/aws/aws-lambda-go/events"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return apishared.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {

		eventList, err := ddb.GetEvents()
		if err != nil {
			return 500, nil, err
		}

		return 200, eventList, nil
	})
}
