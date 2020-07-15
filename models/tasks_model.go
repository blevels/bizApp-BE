package models

import (
	"github.com/satori/go.uuid"
)

type Task struct {
	UUID      	uuid.UUID `json:"uuid" form:"-"`
	Task   		string    `json:"task,omitempty"`
	Status 		bool      `json:"status,omitempty"`
}