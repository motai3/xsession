package pkg

import "xsession/pkg/util/uuid"

func NewSessionId() string {
	return uuid.NextStringId()
}
