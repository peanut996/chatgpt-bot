package model

import "fmt"

type User struct {
	UserName string `json:"username,omitempty"`

	FirstName string `json:"first_name,omitempty"`

	LastName string `json:"last_name,omitempty"`
}

func (u *User) String() string {
	if u.UserName != "" {
		return fmt.Sprintf("@%s", u.UserName)
	}
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func NewUser(username, firstName, lastName string) *User {
	return &User{
		UserName:  username,
		FirstName: firstName,
		LastName:  lastName,
	}
}
