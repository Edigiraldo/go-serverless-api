package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildUserUpdateExpression(t *testing.T) {
	type args struct {
		user User
	}
	type wants struct {
		expression string
	}
	type test struct {
		name  string
		args  args
		wants wants
	}

	tests := []test{
		{
			args: args{
				user: User{
					Email:       "test@example.com",
					FirstName:   "first name",
					LastName:    "last name",
					PhoneNumber: "0987654321",
				},
			},
			wants: wants{
				expression: "set first_name=:fn,last_name=:ln,phone_number=:pn",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expression := buildUserUpdateExpression(test.args.user)

			assert.Equal(t, test.wants.expression, expression)
		})
	}
}
