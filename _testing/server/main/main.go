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
	"github.com/gabriel-vasile/mimetype"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	// This is a local testing server. It has two components:
	// 1. Serving static files
	// 2. Serving the API

	http.HandleFunc("/", handleStatic)

	// List API endpoints
	http.HandleFunc("/api", createHandler(deleteevent.Handle))
	http.HandleFunc("/api", createHandler(deleteround.Handle))
	http.HandleFunc("/api", createHandler(getevent.Handle))
	http.HandleFunc("/api", createHandler(getevents.Handle))
	http.HandleFunc("/api", createHandler(putevent.Handle))
	http.HandleFunc("/api", createHandler(putround.Handle))
	http.HandleFunc("/api", createHandler(stats.Handle))
	http.HandleFunc("/api", createHandler(updateevent.Handle))

	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)
}

func createHandler(apiHandler func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Within here, transform request into events.APIGatewayV2HTTPRequest
		// When we get a return value, write it to the writer
	}
}

func handleApi(writer http.ResponseWriter, request *http.Request) {

}

// HandleStatic retrieves static files based on the path, and returns them as-is.
// This function assumes that the working directory is the root of the project
func handleStatic(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			writer.WriteHeader(500)
			_, err := writer.Write([]byte(fmt.Sprint(r)))
			if err != nil {
				fmt.Println("Error writing response: ", err)
			}
		}
	}()

	fmt.Printf("%s %s\n", request.Method, request.URL.String())

	// Only supports get requests
	if strings.ToUpper(request.Method) != "GET" {
		panic("Unsupported HTTP Method")
	}

	path := request.URL.Path
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(wd, "client", path)

	// If the file does not exist, return 404
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			writer.WriteHeader(404)
			return
		} else {
			panic(err)
		}
	}

	// Otherwise, return the file
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Get the mime type from the filename
	m, err := mimetype.DetectFile(filePath)
	if err != nil {
		panic(err)
	}

	writer.Header().Set("Content-Type", m.String())

	writer.WriteHeader(200)

	_, err = writer.Write(b)
	if err != nil {
		fmt.Println("ERROR: Failed to write to response: ", err)
		return
	}

	return
}
