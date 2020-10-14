package operations

import (
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/models/flow"

var filterOperation = FilterOperation{}

func buildInputUpMessageFilter(type_ flow.UpMessageType) flow.UpMessage {
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	var message = flow.UpMessage{
		Id:      "00000000-000000-00000-000000000",
		Time:    eventTime,
		Content: new(interface{}),
		Type_:   type_,
		Origin: &flow.UpOrigin{
			Type_: "binder",
			Id:    "tpw",
			Time:  time.Now(),
		},
		SubAccount: &flow.Account{
			Id:      "sub1",
			RealmId: "realm1",
		},
		Subscriber: &flow.Subscriber{
			Id:      "sub1",
			RealmId: "realm1",
		},
		Thing: &flow.Thing{
			Key: "lora:0102030405060708",
		},
	}
	return message
}

func Test_should_include_message_for_keep_device_downlink_sent_true_type_device_downlink_sent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceDownlinkSent: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Equal(t, *outputMessage, inputUpMessage)
}

func Test_should_exclude_message_for_keep_device_downlink_sent_true_type_not_device_downlink_sent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceDownlinkSent: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}

func Test_should_exclude_message_for_keep_device_downlink_sent_false_type_device_downlink_sent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceDownlinkSent: false}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}

func Test_exclude_message_for_keep_device_downlink_sent_false_type_not_device_downlink_sent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceDownlinkSent: false}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_include_message_for_keep_device_uplink_true_type_device_uplink(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceUplink: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Equal(t, *outputMessage, inputUpMessage)
}
func Test_should_exclude_message_for_keep_device_uplink_true_type_not_device_uplink(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceUplink: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}

func Test_should_exclude_message_for_keep_device_uplink_false_type_device_uplink(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceUplink: false}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_uplink_false_type_not_device_uplink(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceUplink: false}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}

func Test_should_include_message_for_keel_device_Location_true_type_device_location(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceLocation")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceLocation: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Equal(t, *outputMessage, inputUpMessage)
}

func Test_should_exclude_message_for_keep_device_location_true_type_not_device_location(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceLocation: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_location_false_type_device_location(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceLocation")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceLocation: false}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_location_false_type_not_device_location(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceLocation: false}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_notification_false_message_type_device_notification(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceNotification")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: false}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_notification_false_message_type_not_device_notification(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: false}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_include_message_for_keep_device_notification_true_sub_type_null_message_sub_type_null_device_notification_present(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceNotification")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Equal(t, *outputMessage, inputUpMessage)
}
func Test_should_exclude_message_for_keep_device_notification_true_sub_type_null_message_sub_type_null_device_notification_absent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_include_message_for_keep_device_notification_true_sub_type_null_message_sub_type_not_null_device_notification_present(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceNotification")
	inputUpMessage.SubType = "subType1"
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then
	assert.Equal(t, *outputMessage, inputUpMessage)
}
func Test_should_exclude_message_for_keep_device_notification_true_sub_type_null_message_sub_type_not_null_device_notification_absent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceDownlinkSent")
	inputUpMessage.SubType = "subType1"
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_include_message_for_keep_device_notification_true_sub_type_empty_message_sub_type_null_device_notification_present(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceNotification")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: *new([]string)}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Equal(t, *outputMessage, inputUpMessage)
}
func Test_should_exclude_message_for_keep_device_notification_true_sub_type_empty_message_sub_type_null_device_notification_absent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: *new([]string)}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_notification_true_sub_type_empty_message_sub_type_not_null_device_notification_absent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	inputUpMessage.SubType = "subType1"
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: *new([]string)}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then

	assert.Nil(t, outputMessage)
}

func Test_should_exclude_message_for_keep_device_notification_true_sub_type_empty_message_sub_type_not_null_device_notification_present(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceNotification")
	inputUpMessage.SubType = "subType1"
	var emptyStringArray = []string{}
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: emptyStringArray}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then
	assert.Nil(t, outputMessage)
}
func Test_should_include_message_for_keep_device_notification_true_sub_type_not_empty_message_sub_type_not_null_and_match_device_notification_present(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceNotification")
	inputUpMessage.SubType = "subType1"
	var emptyStringArray = []string{"subType1", "subType2"}
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: emptyStringArray}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then
	assert.Equal(t, *outputMessage, inputUpMessage)
}
func Test_should_exclude_message_for_keep_deviceNotification_true_sub_type_not_empty_message_sub_type_not_null_and_match_device_notification_absent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	inputUpMessage.SubType = "subType1"
	var emptyStringArray = []string{"subType1", "subType2"}
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: emptyStringArray}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then
	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_notification_true_sub_type_not_empty_message_sub_type_not_null_and_not_match_device_notification_present(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceNotification")
	inputUpMessage.SubType = "subType3"
	var emptyStringArray = []string{"subType1", "subType2"}
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: emptyStringArray}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then
	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_notification_true_sub_type_not_empty_message_sub_type_not_null_and_not_match_device_notification_absent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	inputUpMessage.SubType = "subType3"
	var emptyStringArray = []string{"subType1", "subType2"}
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: emptyStringArray}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then
	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_notification_true_sub_type_not_empty_message_sub_type_null_device_notification_present(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceNotification")
	var emptyStringArray = []string{"subType1", "subType2"}
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: emptyStringArray}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then
	assert.Nil(t, outputMessage)
}
func Test_should_exclude_message_for_keep_device_notification_true_sub_type_not_empty_message_sub_type_null_device_notification_absent(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessageFilter("deviceUplink")
	var emptyStringArray = []string{"subType1", "subType2"}
	var upFilterOperation ontology.UpOperationInterface = ontology.UpFilterOperation{KeepDeviceNotification: true, KeepDeviceNotificationSubTypes: emptyStringArray}
	// When
	outputMessage, _ := filterOperation.ApplyUpOperation(&inputUpMessage, &upFilterOperation)
	// Then
	assert.Nil(t, outputMessage)
}
