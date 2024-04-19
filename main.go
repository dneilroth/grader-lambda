package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type urlReq struct {
	URL string `json:"url"`
}

func handler(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// Log the body for debugging purposes
	fmt.Println("Received event: ", event)

	if event.Body == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "empty body",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var url urlReq
	// Read the body of the request
	err := json.Unmarshal([]byte(event.Body), &url)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "unable to unmarshal",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, err
	}

	resp, err := http.Get(url.URL)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "unable to make request to URL",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "unable to read response body",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, err
	}

	response, err := json.Marshal(body)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error creating response JSON",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(response),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
