package ontology

type UpFilterOperation struct {
	KeepDeviceUplink               bool     `json:"keepDeviceUplink,omitempty"`
	KeepDeviceDownlinkSent         bool     `json:"keepDeviceDownlinkSent,omitempty"`
	KeepDeviceLocation             bool     `json:"keepDeviceLocation,omitempty"`
	KeepDeviceNotification         bool     `json:"keepDeviceNotification,omitempty"`
	KeepDeviceNotificationSubTypes []string `json:"keepDeviceNotificationSubTypes,omitempty"`
	UpOperation
}

func (filter UpFilterOperation) ValidUpOperation() string {
	return "filter"
}
