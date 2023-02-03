package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Edigiraldo/go-serverless-api/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynamoCli dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	if len(region) == 0 {
		log.Fatalf("region environment variable was not found")
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("could not create session: %v", err)
	}

	dynamoCli = dynamodb.New(awsSession)

	lambda.Start(handler)
	fmt.Print(region, awsSession, err)
	fmt.Println("Run go project")
}

const tablename = "users"

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case http.MethodGet:
		return handlers.GetUser(ctx, req, tablename, dynamoCli)
	case http.MethodPost:
		return handlers.CreateUser(ctx, req, tablename, dynamoCli)
	case http.MethodPut:
		return handlers.UpdateUser(ctx, req, tablename, dynamoCli)
	case http.MethodDelete:
		return handlers.DeleteUser(ctx, req, tablename, dynamoCli)
	default:
		return handlers.UnhandledMethod(ctx, req)
	}
}
