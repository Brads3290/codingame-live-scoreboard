package main

import (
	getevent "codingame-live-scoreboard/api/getevent/handler"
	putevent "codingame-live-scoreboard/api/putevent/handler"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"time"
)

type T struct {
	Str string
	Tim *time.Time
}

func main() {
	//var t T
	//t.Str = "Hello world"
	//
	//rv := reflect.ValueOf(t)
	//rt := reflect.TypeOf(t)
	//
	//fmt.Println(rv, rt)

	testGetEvent()
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
