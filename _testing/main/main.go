package main

import (
	getevent "codingame-live-scoreboard/api/getevent/handler"
	putevent "codingame-live-scoreboard/api/putevent/handler"
	putround "codingame-live-scoreboard/api/putround/handler"
	updateevent "codingame-live-scoreboard/api/updateevent/handler"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

func main() {
	//testPutRound("round_1")
	//testPutRound("round_2")
	//testPutRound("round_3")

	//testPutRound("15368558b7207af754ef51f1dbc58d3f18a003d")
	testUpdateEvent("3d2183f5-9238-4959-95dd-79d9b088a17f")
}

func testPutEvent() {
	req := events.APIGatewayV2HTTPRequest{
		Body: `{"name": "Test Event"}`,
	}

	resp, err := putevent.Handle(context.TODO(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}

func testGetEvent() {
	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"guid": "3d2183f5-9238-4959-95dd-79d9b088a17f",
		},
	}

	resp, err := getevent.Handle(context.TODO(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}

func testPutRound(roundid string) {
	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"event_id": "3d2183f5-9238-4959-95dd-79d9b088a17f",
		},

		Body: fmt.Sprintf(`{"round_id": "%s"}`, roundid),
	}

	resp, err := putround.Handle(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}

func testUpdateEvent(eventId string) {
	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"event_id": eventId,
		},
	}

	resp, err := updateevent.Handle(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}
