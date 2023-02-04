package users

import (
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

func GetUser(email, tableName string, dynamoCli dynamodbiface.DynamoDBAPI) (*User, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}

	result, err := dynamoCli.GetItem(input)
	if err != nil {
		return &User{}, err
	}

	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func GetUsers(tableName string, dynamoCli dynamodbiface.DynamoDBAPI) error {
	return nil
}

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

func DeleteUser(tableName string, dynamoCli dynamodbiface.DynamoDBAPI) error {
	return nil
}

// builds "set first_name=:fn, last_name=:ln, phone_number=:pn" for non empty values
func buildUserUpdateExpression(user User) string {
	updateExpression := "set "
	toUpdate := make([]string, 3)
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
