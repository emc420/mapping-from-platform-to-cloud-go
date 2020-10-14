package flow

type MessagePacket struct {
	Type_   string      `json:"type"`
	Raw     string      `json:"raw,omitempty"`
	Message interface{} `json:"message,omitempty"`
}
