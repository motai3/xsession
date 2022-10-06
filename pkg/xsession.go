package xsession

import (
	"errors"
	"github.com/motai3/xsession/pkg/util/uuid"
)

var ErrorDisabled = errors.New("the session is disabled storage")

func NewSessionId() string {
	return uuid.NextStringId()
}
