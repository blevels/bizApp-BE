package models

import (
	"github.com/satori/go.uuid"
)

type Permission struct {
	UUID      	uuid.UUID `json:"uuid" form:"-"`
	Permission	string    `json:"role,omitempty"`
}