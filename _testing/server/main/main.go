package main

import (
	event "codingame-live-scoreboard/api/event/handler"
	round "codingame-live-scoreboard/api/round/handler"
	scoreboard "codingame-live-scoreboard/api/scoreboard/handler"
	stats "codingame-live-scoreboard/api/stats/handler"
	update "codingame-live-scoreboard/api/update/handler"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var handlerMap = make(map[string]func(writer http.ResponseWriter, request *http.Request))

var mimeMap = map[string]string{
	".css":  "text/css",
	".html": "text/html",
	".js":   "application/javascript",
}

func main() {

	// This is a local testing server. It has two components:
	// 1. Serving static files
	// 2. Serving the API

	// List API endpoints
	handlerMap["/api/stats"] = createHandler(stats.Handle, regexp.MustCompile(`/api/stats/(?P<event_id>[\w-]+)`))
	handlerMap["/api/scoreboard"] = createHandler(scoreboard.Handle, regexp.MustCompile(`/api/scoreboard/(?P<event_id>[\w-]+)`))
	handlerMap["/api/event"] = createHandler(event.Handle, regexp.MustCompile(`/api/event/(?P<event_id>[\w-]+)`))
	handlerMap["/api/update"] = createHandler(update.Handle, regexp.MustCompile(`/api/update/(?P<event_id>[\w-]+)`))
	handlerMap["/api/round"] = createHandler(round.Handle, regexp.MustCompile(`/api/round/(?P<event_id>[\w-]+)`), regexp.MustCompile(`/api/round/(?P<event_id>[\w-]+)/(?P<round_id>[\w-]+)`))

	http.HandleFunc("/", handleStatic)

	err := http.ListenAndServe(":8080", nil)

	fmt.Println(err)
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

	// First, check if the path matches any in handlerMap, otherwise serve from static.
	path := request.URL.Path
	for k, v := range handlerMap {
		if strings.Index(path, k) == 0 {
			v(writer, request)
			return
		}
	}

	fmt.Printf("%s %s\n", request.Method, request.URL.String())

	// Only supports get requests
	if strings.ToUpper(request.Method) != "GET" {
		panic("Unsupported HTTP Method")
	}

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
	var mimeType string
	ext := filepath.Ext(filePath)
	if mt, ok := mimeMap[ext]; ok {
		mimeType = mt
	} else {
		mimeType = "text/plain"
	}

	writer.Header().Set("Content-Type", mimeType)

	writer.WriteHeader(200)

	_, err = writer.Write(b)
	if err != nil {
		fmt.Println("ERROR: Failed to write to response: ", err)
		return
	}

	return
}

func createHandler(apiHandler func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error), re ...*regexp.Regexp) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				writer.WriteHeader(500)
				_, err := writer.Write([]byte(fmt.Sprint(r)))
				if err != nil {
					fmt.Println("Error writing response: ", err)
				}
			}
		}()

		// Within here, transform request into events.APIGatewayV2HTTPRequest
		// When we get a return value, write it to the write

		defer request.Body.Close()

		// Get the body for the awsReqData
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		awsReqData := events.APIGatewayV2HTTPRequest{
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method:   request.Method,
					Path:     request.URL.Path,
					Protocol: request.URL.Scheme,
				},
			},
			QueryStringParameters: convertQueryStringParams(request.URL.Query()),
			Body:                  string(b),
			RawPath:               request.URL.Path,
			RawQueryString:        request.URL.RawQuery,
			PathParameters:        getPathParameters(request.URL.Path, re...),
		}

		awsRespData, err := apiHandler(context.TODO(), awsReqData)
		if err != nil {
			panic(err)
		}

		for k, v := range awsRespData.Headers {
			writer.Header().Set(k, v)
		}

		writer.WriteHeader(awsRespData.StatusCode)

		_, err = writer.Write([]byte(awsRespData.Body))
		if err != nil {
			panic(err)
		}

		return
	}
}

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	if match == nil {
		return nil
	}

	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}

func convertQueryStringParams(m map[string][]string) map[string]string {
	ret := make(map[string]string)

	for k, v := range m {
		ret[k] = strings.Join(v, ",")
	}

	return ret
}

func getPathParameters(p string, regexes ...*regexp.Regexp) map[string]string {

	m := make(map[string]string)
	for _, r := range regexes {
		matches := reSubMatchMap(r, p)

		if matches != nil {
			for k, v := range matches {
				m[k] = v
			}
		}
	}

	return m
}
