package mocks

import v1 "core/api/v1"

func User() *v1.User {
	return &v1.User{
		Id:          "user_id",
		Email:       "user@example.com",
		PhoneNumber: "1234567890",
		FirstName:   "John",
		LastName:    "Doe",
	}
}
