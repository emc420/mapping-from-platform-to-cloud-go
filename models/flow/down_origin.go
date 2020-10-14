package flow

import (
	"time"
)

type DownOrigin struct {
	Type_        DownOriginType `json:"type"`
	Id           string         `json:"id"`
	ConnectionId string         `json:"connectionId,omitempty"`
	Time         time.Time      `json:"time"`
}
