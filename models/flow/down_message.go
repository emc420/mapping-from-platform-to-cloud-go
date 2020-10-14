package flow

import (
	"time"
)

type DownMessage struct {
	Id         string          `json:"id,omitempty"`
	Time       time.Time       `json:"time"`
	Type_      DownMessageType `json:"type"`
	Content    interface{}     `json:"content"`
	Origin     *DownOrigin      `json:"origin"`
	Command    *Command         `json:"command,omitempty"`
	SubAccount *Account         `json:"subAccount,omitempty"`
	Subscriber *Subscriber      `json:"subscriber"`
	Thing      *Thing           `json:"thing"`
	Packet     *MessagePacket   `json:"packet,omitempty"`
}
