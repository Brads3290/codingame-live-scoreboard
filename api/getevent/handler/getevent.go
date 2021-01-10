package handler

import (
	"codingame-live-scoreboard/codezone_util"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	// Path params will contain a 'guid' parameter

	return codezone_util.UnifyLambdaResponse(ctx, func() (sts int, resp interface{}, err error) {

		// Check that the url has a {guid}

		// Validate the guid by parsing it
		u, err := uuid.Parse()

	})
}
