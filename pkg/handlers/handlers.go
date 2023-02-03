package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func GetUser(ctx context.Context, req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}

func CreateUser(ctx context.Context, req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}

func UpdateUser(ctx context.Context, req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}

func DeleteUser(ctx context.Context, req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}

func UnhandledMethod(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return nil, nil
}
