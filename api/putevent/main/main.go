package main

import (
	"codingame-live-scoreboard/api/putevent/handler"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(putevent.Handle)
}