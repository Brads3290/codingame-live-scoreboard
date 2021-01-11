package _boilerplate

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return api.UnifyLambdaResponse(ctx, func() (int, interface{}, error) {
		return 0, nil, nil
	})
}

func main() {
	lambda.Start(Handle)
}
