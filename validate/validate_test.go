package validate

import (
	"fmt"
	"testing"
)

type User struct {
	FirstName string
	Email     string
}

func TestFooBarBaz(t *testing.T) {
	user := User{
		Email:     "gggg.com",
		FirstName: "bear",
	}
	errors, ok := Validate(user, Fields{
		"Email":     Rules(Email),
		"FirstName": Rules(Min(10), Max(20)),
	})
	if !ok {
		fmt.Println(errors)
	}

	// Validate(user, Fields{"Email": "email|required|max:10|min:2"})
}
