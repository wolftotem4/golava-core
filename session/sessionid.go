package session

import "github.com/google/uuid"

func NewSessionId() string {
	return uuid.New().String()
}
