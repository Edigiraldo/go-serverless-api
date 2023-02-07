package users

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type User struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

var (
	ErrUserNotFound = "user not found"
)

// Creates an user item or replaces an existing one in dynamoDB table
func CreateUser(user User, tableName string, dynamoCli dynamodbiface.DynamoDBAPI) error {

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
func GetUser(email, tableName string, dynamoCli dynamodbiface.DynamoDBAPI) (*User, error) {

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
			return &User{}, errors.New(ErrUserNotFound)
		}

		return &User{}, err
	}

	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

// Gets all users from database
func GetUsers(tableName string, dynamoCli dynamodbiface.DynamoDBAPI) ([]*User, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoCli.Scan(input)
	if err != nil {
		return nil, err
	}

	users := []*User{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Updates non empty user attributes in dynamoDB table
func UpdateUser(user User, tableName string, dynamoCli dynamodbiface.DynamoDBAPI) (*User, error) {

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
		return &User{}, err
	}

	updatedUserAttrs := User{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedUserAttrs)
	if err != nil {
		return &User{}, err
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
func buildUserUpdateExpression(user User) string {
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
