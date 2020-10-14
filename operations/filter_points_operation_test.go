package operations

import (
	"github.com/stretchr/testify/assert"
	"ontology-mapping-go-lib/models/flow"
	"ontology-mapping-go-lib/models/ontology"
	"ontology-mapping-go-lib/util"
	"testing"
	"time"
)

var filterPointOperation = FilterPointsOperation{}

func buildInputFilterPointUpMessage(inputPoints map[string]flow.Point) flow.UpMessage {
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	var message = flow.UpMessage{
		Id:    "00000000-000000-00000-000000000",
		Time:  eventTime,
		Type_: "deviceUplink",
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
		Points: inputPoints,
	}
	return message
}
func Test_should_be_unchanged_message_when_points_filter_is_empty_message_points_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputFilterPointUpMessage(make(map[string]flow.Point))
	points := []string{}
	var upFilterPoints ontology.UpOperationInterface = ontology.UpFilterPointsOperation{Points: points}

	// When
	outputUpMessage, _ := filterPointOperation.ApplyUpOperation(&inputUpMessage, &upFilterPoints)

	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}

func Test_should_return_empty_points_message_when_points_filter_is_empty_message_points_not_empty(t *testing.T) {
	// Given
	var inputRecords []flow.Record
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	inputRecords = append(inputRecords, flow.Record{
		Value:     22.6,
		EventTime: eventTime,
	})
	inputTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: inputRecords,
	}
	inputPoints := map[string]flow.Point{
		"temperature": inputTemperature,
	}
	inputUpMessage := buildInputUpUpdateMessage(inputPoints)
	points := []string{}
	var upFilterPoints ontology.UpOperationInterface = ontology.UpFilterPointsOperation{Points: points}

	// When
	outputUpMessage, _ := filterPointOperation.ApplyUpOperation(&inputUpMessage, &upFilterPoints)
	// Then
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = make(map[string]flow.Point)

	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_be_unchanged_message_when_points_filter_is_not_empty_message_points_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputFilterPointUpMessage(make(map[string]flow.Point))
	points := []string{"temperature", "humidity"}
	var upFilterPoints ontology.UpOperationInterface = ontology.UpFilterPointsOperation{Points: points}

	// When
	outputUpMessage, _ := filterPointOperation.ApplyUpOperation(&inputUpMessage, &upFilterPoints)

	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}

func Test_should_be_unchanged_message_when_points_filter_is_not_empty_message_points_not_empty_match(t *testing.T) {
	// Given
	var inputRecords []flow.Record
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	inputRecords = append(inputRecords, flow.Record{
		Value:     22.6,
		EventTime: eventTime,
	})
	inputTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: inputRecords,
	}
	inputPoints := map[string]flow.Point{
		"temperature": inputTemperature,
	}
	inputUpMessage := buildInputUpUpdateMessage(inputPoints)
	points := []string{"temperature", "humidity"}
	var upFilterPoints ontology.UpOperationInterface = ontology.UpFilterPointsOperation{Points: points}

	// When
	outputUpMessage, _ := filterPointOperation.ApplyUpOperation(&inputUpMessage, &upFilterPoints)
	// Then
	assert.Equal(t, *outputUpMessage, inputUpMessage)
}

func Test_should_return_empty_points_when_points_filter_is_not_empty_message_points_not_empty_not_match(t *testing.T) {
	// Given
	var inputRecords []flow.Record
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	inputRecords = append(inputRecords, flow.Record{
		Value:     22.6,
		EventTime: eventTime,
	})
	inputTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: inputRecords,
	}
	inputPoints := map[string]flow.Point{
		"temperature": inputTemperature,
	}
	inputUpMessage := buildInputUpUpdateMessage(inputPoints)
	points := []string{"humidity"}
	var upFilterPoints ontology.UpOperationInterface = ontology.UpFilterPointsOperation{Points: points}

	// When
	outputUpMessage, _ := filterPointOperation.ApplyUpOperation(&inputUpMessage, &upFilterPoints)
	// Then
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = make(map[string]flow.Point)

	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_return_filtered_points_when_points_filter_has_multiple_entries_message_points_not_empty_selective_match(t *testing.T) {
	// Given
	var inputRecordsTemperature []flow.Record
	var inputRecordsBatteryLevel []flow.Record
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	inputRecordsTemperature = append(inputRecordsTemperature, flow.Record{
		Value:     22.6,
		EventTime: eventTime,
	})
	inputRecordsBatteryLevel = append(inputRecordsBatteryLevel, flow.Record{
		Value:     30,
		EventTime: eventTime,
	})
	inputTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: inputRecordsTemperature,
	}
	inputBatteryCurrentLevel := flow.Point{
		Type_:   "double",
		UnitId:  "%RH",
		Records: inputRecordsBatteryLevel,
	}
	inputPoints := map[string]flow.Point{
		"temperature":           inputTemperature,
		"battery_current_level": inputBatteryCurrentLevel,
	}
	inputUpMessage := buildInputUpUpdateMessage(inputPoints)
	points := []string{"battery_current_level", "coordinates", "batteryLevel", "batteryStatus"}
	var upFilterPoints ontology.UpOperationInterface = ontology.UpFilterPointsOperation{Points: points}

	// When
	outputUpMessage, _ := filterPointOperation.ApplyUpOperation(&inputUpMessage, &upFilterPoints)
	// Then
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = map[string]flow.Point{
		"battery_current_level": inputBatteryCurrentLevel,
	}

	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
