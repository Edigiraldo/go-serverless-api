package handlers

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/Edigiraldo/go-serverless-api/pkg/users"
	"github.com/Edigiraldo/go-serverless-api/pkg/validators"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrMethodNotAllowed    = "method not allowed"
	ErrInvalidUserFormat   = "the user format is invalid"
	ErrInternalServerError = "there was a server error"
	ErrCreatingUser        = "there was an error while creating the user"
	ErrInvalidEmailFormat  = "the email format is invalid"
	ErrUpdatingUser        = "there was an error while updating the user"
	ErrDeletingUser        = "there was an error while deleting the user"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func CreateUser(req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (res *events.APIGatewayProxyResponse, err error) {
	body := []byte(req.Body)

	user := users.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		errBody := ErrorBody{
			ErrorMsg: &ErrInvalidUserFormat,
		}

		return ApiResponse(http.StatusBadRequest, errBody)
	}

	if err = users.CreateUser(user, tablename, dynamoCli); err != nil {
		errBody := ErrorBody{
			ErrorMsg: &ErrCreatingUser,
		}
		log.Print(err)

		return ApiResponse(http.StatusInternalServerError, errBody)
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		errBody := ErrorBody{
			ErrorMsg: &ErrInternalServerError,
		}
		log.Print(err)

		return ApiResponse(http.StatusInternalServerError, errBody)
	}

	return ApiResponse(http.StatusCreated, jsonUser)
}

func GetUser(req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if !validators.IsEmailValid(email) {
		errBody := ErrorBody{
			ErrorMsg: &ErrInvalidEmailFormat,
		}

		return ApiResponse(http.StatusBadRequest, errBody)
	}

	user, err := users.GetUser(email, tablename, dynamoCli)
	if err != nil {
		errBody := ErrorBody{
			ErrorMsg: nil,
		}
		if err.Error() == users.ErrUserNotFound {
			errBody.ErrorMsg = &users.ErrUserNotFound

			return ApiResponse(http.StatusNotFound, errBody)
		} else {
			errBody.ErrorMsg = &ErrCreatingUser
			log.Print(err)

			return ApiResponse(http.StatusInternalServerError, errBody)
		}
	}

	jsonUser, err := json.Marshal(*user)
	if err != nil {
		errBody := ErrorBody{
			ErrorMsg: &ErrInternalServerError,
		}
		log.Print(err)

		return ApiResponse(http.StatusInternalServerError, errBody)
	}

	return ApiResponse(http.StatusOK, jsonUser)
}

func UpdateUser(req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	body := []byte(req.Body)

	user := users.User{}
	err := json.Unmarshal(body, &user)
	if err != nil {
		errBody := ErrorBody{
			ErrorMsg: &ErrInvalidUserFormat,
		}

		return ApiResponse(http.StatusBadRequest, errBody)
	}

	updatedUser, err := users.UpdateUser(user, tablename, dynamoCli)
	if err != nil {
		errBody := ErrorBody{
			ErrorMsg: &ErrUpdatingUser,
		}
		log.Print(err)

		return ApiResponse(http.StatusInternalServerError, errBody)
	}

	jsonUser, err := json.Marshal(*updatedUser)
	if err != nil {
		errBody := ErrorBody{
			ErrorMsg: &ErrInternalServerError,
		}
		log.Print(err)

		return ApiResponse(http.StatusInternalServerError, errBody)
	}

	return ApiResponse(http.StatusCreated, jsonUser)
}

func DeleteUser(req events.APIGatewayProxyRequest, tablename string, dynamoCli dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if !validators.IsEmailValid(email) {
		errBody := ErrorBody{
			ErrorMsg: &ErrInvalidEmailFormat,
		}

		return ApiResponse(http.StatusBadRequest, errBody)
	}

	if err := users.DeleteUser(email, tablename, dynamoCli); err != nil {
		errBody := ErrorBody{
			ErrorMsg: nil,
		}
		if err.Error() == users.ErrUserNotFound {
			errBody.ErrorMsg = &users.ErrUserNotFound

			return ApiResponse(http.StatusNoContent, nil)
		} else {
			errBody.ErrorMsg = &ErrDeletingUser
			log.Print(err)

			return ApiResponse(http.StatusInternalServerError, errBody)
		}
	}

	return ApiResponse(http.StatusNoContent, nil)
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
