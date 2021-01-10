package _boilerplate

func Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return codezone_util.UnifyLambdaResponse(ctx, func() (sts int, resp interface{}, err error) {
		return
	})
}

func main() {
	lambda.Start(Handle)
}
