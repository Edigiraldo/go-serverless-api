package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrMethodNotAllowed = "method not allowed"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func CreateUser(req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}

func GetUser(req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}

/**
 * Handles requests with unexpected methods
 *
 * @param req events.APIGatewayProxyRequest - request received
 *
 * @returns *events.APIGatewayProxyResponse - response
 */
func UnhandledMethod(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	errBody := ErrorBody{
		ErrorMsg: &ErrMethodNotAllowed,
	}

	return ApiResponse(http.StatusMethodNotAllowed, errBody)
}
