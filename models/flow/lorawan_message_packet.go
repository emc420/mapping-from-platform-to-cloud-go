package flow

type LorawanMessagePacket struct {
	Type_   string             `json:"type"`
	Raw     string             `json:"raw,omitempty"`
	Message interface{}        `json:"message,omitempty"`
	Meta    *LorawanPacketMeta `json:"meta"`
}
