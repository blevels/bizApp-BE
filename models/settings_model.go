package models

import (
	"github.com/satori/go.uuid"
)

type Setting struct {
	UUID      uuid.UUID `json:"uuid" form:"-"`
	Whitelist struct{
		Domain	string `json:"domain"`
	} `json:"whitelist"`
	Sources struct {
		Source  string `json:"source"`
	} `json:"sources"`
}