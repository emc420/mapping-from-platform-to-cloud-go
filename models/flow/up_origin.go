package flow

import (
	"time"
)

type UpOrigin struct {
	Type_        UpOriginType `json:"type"`
	Id           string       `json:"id"`
	ConnectionId string       `json:"connectionId,omitempty"`
	Time         time.Time    `json:"time"`
}
