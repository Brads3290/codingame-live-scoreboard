package main

import (
	getevent "codingame-live-scoreboard/api/getevent/handler"
	putevent "codingame-live-scoreboard/api/putevent/handler"
	putround "codingame-live-scoreboard/api/putround/handler"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/dbschema"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

func main() {
	//testPutRound("round_1")
	//testPutRound("round_2")
	//testPutRound("round_3")

	var rms []dbschema.RoundModel

	err := ddb.QueryItemsFromDynamoDbWithFilter(constants.DB_TABLE_ROUNDS, &rms, map[string]interface{}{
		"Event_ID": "3d2183f5-9238-4959-95dd-79d9b088a17f",
	}, map[string]interface{}{
		"Is_Active": true,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rms)
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
