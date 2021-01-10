package main

import (
	"codingame-live-scoreboard/api/putevent/handler"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

func main() {
	req := events.APIGatewayV2HTTPRequest{
		Body: `{"name": "Test Event"}`,
	}

	resp, err := handler.Handle(context.TODO(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}
