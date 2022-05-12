package uidutil

import uuid "github.com/satori/go.uuid"

func New() string {
	return uuid.NewV4().String()
}
