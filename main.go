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
	headers := map[string]string{
		"content-type":    "application/json",
		"accept":          "*/*",
		"sec-fetch-site":  "same-origin",
		"accept-language": "en-US,en;q=0.9",
		"sec-fetch-mode":  "cors",
		"origin":          "https://google.com",
		"user-agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 16_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Mobile/15E148 Safari/604.1",
		"referer":         "https://google.com/",
		"sec-fetch-dest":  "empty",
	}

	req, err := http.NewRequest("GET", url.URL, nil)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "unable to create request",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, err
	}

	for k, h := range headers {
		req.Header.Add(k, h)
	}

	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "unable to make request to URL",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, err
	}
	// Create a client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "error making request",
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

	response, err := json.Marshal(string(body))
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
