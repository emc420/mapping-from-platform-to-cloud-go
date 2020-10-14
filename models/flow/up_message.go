package flow

import (
	"time"
)

type UpMessage struct {
	Id         string           `json:"id,omitempty"`
	Time       time.Time        `json:"time"`
	Content    interface{}      `json:"content"`
	Type_      UpMessageType    `json:"type"`
	SubType    string           `json:"subType,omitempty"`
	Origin     *UpOrigin         `json:"origin"`
	SubAccount *Account          `json:"subAccount,omitempty"`
	Subscriber *Subscriber       `json:"subscriber"`
	Thing      *Thing            `json:"thing"`
	Points     map[string]Point `json:"points,omitempty"`
	Packet     *MessagePacket    `json:"packet,omitempty"`
}
