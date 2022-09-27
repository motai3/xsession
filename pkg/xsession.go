package xsession

import (
	"errors"
	"xsession/pkg/util/uuid"
)

var ErrorDisabled = errors.New("the session is disabled storage")

func NewSessionId() string {
	return uuid.NextStringId()
}
