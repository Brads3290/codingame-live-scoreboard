package main

import (
	"codingame-live-scoreboard/putevent/handler"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Handle)
}
