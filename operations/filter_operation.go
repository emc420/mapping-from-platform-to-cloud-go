package operations

import (
	"ontology-mapping-go-lib/models/flow"
)
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/util"

type FilterOperation struct {
}

func (filterOperation *FilterOperation) ApplyUpOperation(message *flow.UpMessage, upOperation *ontology.UpOperationInterface) (*flow.UpMessage, error) {

	var upFilterOperation = (*upOperation).(ontology.UpFilterOperation)
	if isFilterPresent(upFilterOperation.KeepDeviceDownlinkSent, string(message.Type_), "deviceDownlinkSent") {
		return message, nil
	}
	if isFilterPresent(upFilterOperation.KeepDeviceUplink, string(message.Type_), "deviceUplink") {
		return message, nil
	}
	if isFilterPresent(upFilterOperation.KeepDeviceLocation, string(message.Type_), "deviceLocation") {
		return message, nil
	}
	if isFilterPresent(upFilterOperation.KeepDeviceNotification, string(message.Type_), "deviceNotification") &&
		isSubtypePresent(upFilterOperation.KeepDeviceNotificationSubTypes, message.SubType) {
		return message, nil
	}
	return nil, nil
}

func (filterOperation *FilterOperation) ApplyDownOperation(message *flow.DownMessage, downOperation *ontology.DownOperationInterface) (*flow.DownMessage, error) {
	return nil, nil
}

func isFilterPresent(filterVal bool, messageType string, typeString string) bool {
	return filterVal && messageType == typeString
}

func isSubtypePresent(subTypes []string, subType string) bool {
	if subTypes == nil {
		return true
	}
	if len(subTypes) == 0 && len(subType) == 0 {
		return true
	}
	if len(subType) > 0 && util.Contains(subTypes, subType) {
		return true
	}
	return false
}
