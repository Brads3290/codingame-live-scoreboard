package main

import (
	"codingame-live-scoreboard/api/getrounds/handler"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Handle)
}
