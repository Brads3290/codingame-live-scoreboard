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
	"strings"
)

// GET /user/{user_id}
// GET /users

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return apishared.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {
		method := request.RequestContext.HTTP.Method
		path := request.RawPath
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}

		pathComponents := strings.Split(path, "/")

		switch method {
		case "GET":
			if pathComponents[0] == "user" {
				return getUser(ctx, request)
			} else if pathComponents[0] == "users" {
				return getUsers(ctx, request)
			} else {
				return 404, nil, errors.New("Unrecognised path")
			}
		case "PUT":

		default:
			return 405, nil, errors.New("Unsupported HTTP method: " + method)
		}
	})
}

// GET /user/{user_id}
func getUser(_ context.Context, req events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
	uid, ok := req.PathParameters["user_id"]
	if !ok {
		return 400, nil, errors.New("Missing user_id")
	}

	var u dbschema.UserModel
	u.UserId = uid

	err := ddb.PopulateItemFromDynamoDb(constants.DB_TABLE_USERS, &u)
	if err != nil {
		return 500, nil, err
	}

	return 200, u, nil
}

func getUsers(_ context.Context, _ events.APIGatewayV2HTTPRequest) (int, interface{}, error) {
	var users []dbschema.UserModel

	err := ddb.ScanItemsFromDynamoDb(constants.DB_TABLE_USERS, &users)
	if err != nil {
		return 500, nil, err
	}

	return 200, users, nil
}

func putUser(_ context.Context, req events.APIGatewayV2HTTPRequest) (int, interface{}, error) {

	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.Unmarshal([]byte(req.Body), &body)
	if err != nil {
		return 500, nil, err
	}

	userGuid, err := ddb.AddUserToDynamoDb(body.Username, body.Password)
	if err != nil {
		return 500, nil, err
	}

	result := struct {
		UserId string `json:"user_id"`
	}{userGuid}

	return 200, result, nil
}
