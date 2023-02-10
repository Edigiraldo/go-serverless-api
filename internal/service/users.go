package service

import (
	"errors"
	"strings"

	"github.com/Edigiraldo/go-serverless-api/internal/dtos"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrUserNotFound = "user not found"
)

// Creates an user item or replaces an existing one in dynamoDB table
func CreateUser(user dtos.User, tableName string, dynamoCli dynamodbiface.DynamoDBAPI) error {

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynamoCli.PutItem(input)

	return err
}

// Gets user by email
func GetUser(email, tableName string, dynamoCli dynamodbiface.DynamoDBAPI) (*dtos.User, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynamoCli.GetItem(input)
	if err != nil {
		var notFoundErr *dynamodb.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			return &dtos.User{}, errors.New(ErrUserNotFound)
		}

		return &dtos.User{}, err
	}

	user := dtos.User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return &dtos.User{}, err
	}

	return &user, nil
}

// Gets all users from database
func GetUsers(tableName string, dynamoCli dynamodbiface.DynamoDBAPI) ([]*dtos.User, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoCli.Scan(input)
	if err != nil {
		return nil, err
	}

	users := []*dtos.User{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Updates non empty user attributes in dynamoDB table
func UpdateUser(user dtos.User, tableName string, dynamoCli dynamodbiface.DynamoDBAPI) (*dtos.User, error) {

	updateExpression := buildUserUpdateExpression(user)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fn": {
				S: aws.String(user.FirstName),
			},
			":ln": {
				S: aws.String(user.LastName),
			},
			":pn": {
				S: aws.String(user.PhoneNumber),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(user.Email),
			},
		},
		UpdateExpression: aws.String(updateExpression),
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	result, err := dynamoCli.UpdateItem(input)
	if err != nil {
		return &dtos.User{}, err
	}

	updatedUserAttrs := dtos.User{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedUserAttrs)
	if err != nil {
		return &dtos.User{}, err
	}

	return &updatedUserAttrs, nil
}

// Deletes an user by email
func DeleteUser(email, tableName string, dynamoCli dynamodbiface.DynamoDBAPI) error {

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynamoCli.DeleteItem(input)
	return err
}

// Builds "set first_name=:fn, last_name=:ln, phone_number=:pn" to just update
// non empty values
func buildUserUpdateExpression(user dtos.User) string {
	updateExpression := "set "
	toUpdate := make([]string, 0)
	if user.FirstName != "" {
		toUpdate = append(toUpdate, "first_name=:fn")
	}
	if user.LastName != "" {
		toUpdate = append(toUpdate, "last_name=:ln")
	}
	if user.PhoneNumber != "" {
		toUpdate = append(toUpdate, "phone_number=:pn")
	}

	updateExpression += strings.Join(toUpdate, ",")

	return updateExpression
}
