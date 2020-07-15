package models

import (
	"github.com/satori/go.uuid"
)

type User struct {
	UUID      uuid.UUID `json:"uuid" form:"-"`
	UserName  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
	Email	  string `json:"email"`
	Token     string `json:"token"`
	Role
}