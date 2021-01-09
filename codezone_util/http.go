package codezone_util

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type lambdaResponseData struct {
	Success bool `json:"success"`
	Error string `json:"error"`
	Data interface{} `json:"data"`
}

func UnifyLambdaResponse(ctx context.Context, f func () (int, interface{}, error)) (events.APIGatewayV2HTTPResponse, error) {
	sts, dataResp, err := f()

	responseData := lambdaResponseData{
		Success: true,
		Data:    dataResp,
	}

	if err != nil {
		responseData.Success = false
		responseData.Error = err.Error()

		if sts == 0 {
			sts = 500
		}
	}

	var resp events.APIGatewayV2HTTPResponse

	b, err := json.Marshal(responseData)
	if err != nil {
		sts = 500

		responseData = lambdaResponseData{
			Success: false,
			Data: nil,
			Error: err.Error(),
		}

		b, err = json.Marshal(responseData)
		if err != nil {
			panic(err) // Should not happen  :/
		}
	}

	if sts == 0 {
		sts = 200
	}

	resp.Body = string(b)
	resp.StatusCode = sts

	return resp, nil
}
