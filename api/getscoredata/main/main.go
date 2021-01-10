package main

import (
	"codingame-live-scoreboard/codezone_util"
	"codingame-live-scoreboard/constants"
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"strconv"
	"time"
)

// Request made to this function, will contain query string params:
// 	- evt=<guid>
// Returned data:
// 	- List of players
//		- Name
//		- Total score in event
//		- Current rank in event
func handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return codezone_util.UnifyLambdaResponse(ctx, func() (sts int, resp interface{}, err error) {

		evtGuid, ok := request.QueryStringParameters["evt"]
		if !ok {
			return sts, resp, errors.New("no evt query specified")
		}

		// Get the event from the dynamodb
		evtData, err := codezone_util.GetItemFromDynamoDb(constants.DB_TABLE_EVENTS, "ID", evtGuid)
		if err != nil {
			return sts, resp, err
		}

		if evtData == nil {
			return 404, resp, nil
		}

		/*evtData contains:
		- ID: str
		- Name: str
		- LastUpdated: int
		*/

		lastUpdatedSecs, err := strconv.Atoi(evtData["lastUpdated"])
		if err != nil {
			return sts, resp, err
		}

		lastUpdated := time.Unix(int64(lastUpdatedSecs), 0)
		now := time.Now()

		if now.Sub(lastUpdated) > constants.MAX_EVENT_RECORD_AGE {
			_, err := codezone_util.GetCodinGameData(evtGuid)
			if err != nil {
				return sts, resp, err
			}

			// Update the DB
		} else {
			// Get the data from the DB
		}

		// Return the data to the client
		return
	})
}

func main() {
	lambda.Start(handle)
}
