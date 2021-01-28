package main

import (
	deleteevent "codingame-live-scoreboard/api/deleteevent/handler"
	deleteround "codingame-live-scoreboard/api/deleteround/handler"
	getevent "codingame-live-scoreboard/api/getevent/handler"
	getevents "codingame-live-scoreboard/api/getevents/handler"
	putevent "codingame-live-scoreboard/api/putevent/handler"
	putround "codingame-live-scoreboard/api/putround/handler"
	stats "codingame-live-scoreboard/api/stats/handler"
	updateevent "codingame-live-scoreboard/api/updateevent/handler"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

func main() {

	//testGetEvents()

	//testPutEvent()
	//testPutRound("245eac65-0a53-4778-93a0-27aab09fff1a", "15368558b7207af754ef51f1dbc58d3f18a003d")

	//testDeleteEvent("245eac65-0a53-4778-93a0-27aab09fff1a")
	//testDeleteRound("245eac65-0a53-4778-93a0-27aab09fff1a" ,"15368558b7207af754ef51f1dbc58d3f18a003d")

	//testUpdateEvent("2f7798d3-8efa-4d73-9918-6bc7431dd764")

	testGetStats("0cbc045e-b127-4b2b-a033-015340c3f731")
}

func testPutEvent() {
	req := events.APIGatewayV2HTTPRequest{
		Body: `{"name": "Test Event 3"}`,
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

func testGetEvents() {
	req := events.APIGatewayV2HTTPRequest{}

	resp, err := getevents.Handle(context.TODO(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}

func testPutRound(eventId string, roundid string) {
	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"event_id": eventId,
		},

		Body: fmt.Sprintf(`{"round_id": "%s"}`, roundid),
	}

	resp, err := putround.Handle(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}

func testDeleteEvent(eventId string) {
	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"event_id": eventId,
		},
	}

	resp, err := deleteevent.Handle(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}

func testDeleteRound(eventId string, roundId string) {
	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"event_id": eventId,
			"round_id": roundId,
		},
	}

	resp, err := deleteround.Handle(context.TODO(), req)

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

func testGetStats(eventId string) {
	req := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"event_id": eventId,
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
			},
		},
	}

	resp, err := stats.Handle(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\v", resp)
}
