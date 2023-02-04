package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func ApiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
	}

	return &resp, err
}
