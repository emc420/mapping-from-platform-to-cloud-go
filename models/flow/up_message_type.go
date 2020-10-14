package flow

type UpMessageType string

const (
	DEVICE_UPLINK_UpMessageType        UpMessageType = "deviceUplink"
	DEVICE_DOWNLINK_SENT_UpMessageType UpMessageType = "deviceDownlinkSent"
	DEVICE_LOCATION_UpMessageType      UpMessageType = "deviceLocation"
	DEVICE_NOTIFICATION_UpMessageType  UpMessageType = "deviceNotification"
)
