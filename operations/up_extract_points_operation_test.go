package operations

import (
	"encoding/json"
	"io/ioutil"
	"ontology-mapping-go-lib/util"
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/models/flow"

var jmesPathOperation = UpExtractPointsOperation{}

func buildInputUpMessage(inputMessageFile string) flow.UpMessage {
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
		Points: make(map[string]flow.Point),
		Packet: &flow.MessagePacket{
			Type_:   "",
			Raw:     "",
			Message: nil,
		},
	}
	byteSream, _ := ioutil.ReadFile("resources/" + inputMessageFile)
	_ = json.Unmarshal(byteSream, &message.Packet.Message)
	return message
}

func Test_should_exclude_points_when_value_is_null(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("missing_point_value.json")
	humidity := ontology.JmesPathPoint{
		Value:     "{{packet.message.humidity}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"humidity": humidity,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)

	// Then
	assert.Equal(t, *outputUpMessage, inputUpMessage)
}

func Test_should_include_existing_points_from_input_message(t *testing.T) {
	// Given

	dummy := flow.Point{
		Type_:   "string",
		UnitId:  "D",
		Records: make([]flow.Record, 1),
	}
	inputPoints := map[string]flow.Point{
		"Dummy": dummy,
	}
	inputUpMessage := buildInputUpMessage("include_existing_points.json")
	inputUpMessage.Points = inputPoints
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.6,
		EventTime: inputUpMessage.Time,
	})
	expectedTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"Dummy":       dummy,
		"temperature": expectedTemperature,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_throw_exception_for_array_length_mismatch(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("length_mismatch_value_event_time.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.measures[?id == 'temperature'].value}}",
		EventTime: "{{packet.message.measures[?id == 'temperature'].time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// When & Then
	assert.EqualError(t, err, "there is a mismatch in cardinality for 'value' and 'eventTime' temperature")
}

func Test_should_throw_exception_for_only_value_is_array_not_event_time(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("only_value_is_array.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.measures[?id == 'temperature'].value}}",
		EventTime: "{{packet.message.measures[?id == 'temperature'].time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// When & Then
	assert.EqualError(t, err, "there is a mismatch in cardinality for 'value' and 'eventTime' temperature")
}

func Test_should_throw_exception_for_only_event_time_is_array_not_value(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("only_event_time_is_array.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// When & Then
	assert.EqualError(t, err, "there is a mismatch in cardinality for 'value' and 'eventTime' temperature")
}
func Test_should_throw_exception_for_event_time_is_null_and_value_not_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("event_time_is_null_and_value_not_null.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// When & Then
	assert.EqualError(t, err, "there is a mismatch in cardinality for 'value' and 'eventTime' temperature")
}

func Test_should_exclude_point_if_value_is_empty_array(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("empty_value_array.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_if_value_and_event_time_are_null(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("null_value_event_time.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_if_value_is_null_event_time_is_empty_array(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("null_value_empty_array_event_time.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	assert.Equal(t, *outputUpMessage, inputUpMessage)
}

func Test_should_exclude_point_if_value_is_empty_event_time_is_empty_array(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("value_empty_array_event_time_empty_array.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_if_value_is_empty_event_time_is_null(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("value_empty_event_time_null.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_throw_exception_if_value_is_not_null_event_time_empty_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("value_not_null_event_time_empty.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// When & Then
	assert.EqualError(t, err, "there is a mismatch in cardinality for 'value' and 'eventTime' temperature")
}
func Test_should_exclude_point_if_value_is_empty_event_time_is_not_null(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("value_is_empty_event_time_not_null.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_throw_exception_if_value_is_array_event_time_empty_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("value_array_event_time_empty.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// When & Then
	assert.EqualError(t, err, "there is a mismatch in cardinality for 'value' and 'eventTime' temperature")
}

func Test_should_return_transformed_message_value_event_time_numbers(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("value_event_time_numbers.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.6,
		EventTime: inputUpMessage.Time,
	})
	expectedTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"temperature": expectedTemperature,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_return_transformed_message_value_event_time_arrays_of_equal_length(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("value_event_time_equal_length_arrays.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	eventTimes, _ := time.Parse(time.RFC3339, "2020-02-06T09:14:05.688Z")
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     11.625,
		EventTime: eventTimes,
	})
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     11.625,
		EventTime: eventTimes,
	})
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     11.6875,
		EventTime: eventTimes,
	})
	expectedTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"temperature": expectedTemperature,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_throw_error_when_event_time_is_not_date_time(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("event_time_is_not_date_time.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// When & Then
	assert.Errorf(t, err, "parsing time")
}
func Test_should_throw_error_when_single_element_in_event_time_array_is_not_date_time(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("element_in_a_event_time_array_is_not_date_time.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// When & Then
	assert.Errorf(t, err, "parsing time")
}
func Test_should_include_point_when_value_is_array_of_equal_length_1_and_event_time_is_not_an_array(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("value_array_size_1_event_time_not_array.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.measures[?id == 'temperature'].value}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     11.625,
		EventTime: inputUpMessage.Time,
	})
	expectedTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"temperature": expectedTemperature,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_include_point_when_event_time_is_array_of_equal_length_1_and_value_is_not_an_array(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("event_time_array_1_value_not_array.json")
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature.value}}",
		EventTime: "{{packet.message.temperature.time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	parsedTime, _ := time.Parse(time.RFC3339, "2020-02-06T09:14:05.688Z")
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     11.625,
		EventTime: parsedTime,
	})
	expectedTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Cel",
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"temperature": expectedTemperature,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_include_point_when_event_time_is_not_array_and_coordinate_is_not_an_array(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinates_not_array_event_time_not_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   inputUpMessage.Time,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_exclude_point_when_event_time_is_null_and_coordinate_is_null(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinate_is_null_event_time_is_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_when_event_time_is_empty_and_coordinate_is_null(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinate_is_null_event_time_is_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_when_event_time_is_null_and_coordinate_is_empty(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinate_is_empty_event_time_is_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_when_event_time_is_empty_and_coordinate_is_empty(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinate_is_empty_event_time_is_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_when_event_time_is_not_null_and_coordinate_is_empty(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinate_is_empty_event_time_is_not_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_throw_exception_when_event_time_is_empty_and_coordinate_is_not_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_is_not_null_event_time_is_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}

func Test_should_include_point_when_event_time_is_array_and_size_1_coordinate_is_not_an_array(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinate_not_array_event_time_array_size_1.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	resTime, _ := time.Parse(time.RFC3339, "2020-01-01T10:00:05.000Z")
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   resTime,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_throw_exception_when_event_time_is_array_size_not_1_and_coordinate_is_not_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_not_array_event_time_is_array_size_not_1.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}

func Test_should_include_point_when_event_time_is_not_array_and_coordinate_is_array_size_1(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinate_is_array_size_1_event_time_is_not_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   inputUpMessage.Time,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_throw_exception_when_event_time_is_not_array_and_coordinate_is_array_size_not_1(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_is_array_size_not_1_event_time_not_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}

func Test_should_include_point_when_event_time_is_array_and_coordinate_is_array_equal_size(t *testing.T) {
	// Given

	inputUpMessage := buildInputUpMessage("coordinates_array_and_size_equals_event_time_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	resTime, _ := time.Parse(time.RFC3339, "2020-01-01T10:00:05.000Z")
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   resTime,
	})
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   resTime,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_throw_exception_when_event_time_is_array_and_coordinate_is_array_size_not_equal(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_array_event_time_array_size_not_equal.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}
func Test_should_exclude_point_when_event_time_is_array_and_coordinate_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_is_empty_event_time_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_throw_exception_when_event_time_is_null_and_coordinate_is_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_is_array_event_time_is_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}
func Test_should_throw_exception_when_event_time_is_null_and_coordinate_is_not_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinates_is_not_null_event_time_is_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}
func Test_should_throw_exception_when_event_time_is_empty_and_coordinate_is_not_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinates_is_not_empty_event_time_is_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}
func Test_should_exclude_point_when_event_time_is_not_empty_and_coordinate_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinates_is_empty_event_time_not_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_coordinate_when_lat_is_null_lng_is_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_null_lng_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_coordinate_when_lat_is_null_lng_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_null_lng_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_coordinate_when_lat_is_empty_lng_is_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_empty_lng_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_throw_exception_when_lat_is_empty_and_lng_is_not_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_empty_lng_not_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_is_not_null_and_lng_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_null_lng_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}

func Test_should_include_point_when_lat_is_not_array_and_lng_is_array_size_1(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_array_lng_array_size_1.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   inputUpMessage.Time,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_throw_exception_when_lat_is_not_array_and_lng_is_array_size_not_1(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_array_lng_array_size_not_1.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}
func Test_should_include_point_when_lat_is_array_size_1_and_lng_is_not_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_array_size_1_lng_not_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   inputUpMessage.Time,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_throw_exception_when_lat_is_array_size_not_1_and_lng_is_not_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_array_size_not_1_lng_is_not_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_is_array_and_lng_is_array_unequal_size(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_array_lng_array_unequal_size.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_is_null_and_lng_is_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_null_lng_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_is_array_and_lng_is_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_array_lng_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_is_null_and_lng_is_not_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_null_lng_not_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_is_not_null_and_lng_is_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_is_not_null_lng_is_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_is_not_empty_and_lng_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_is_not_empty_lng_is_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_is_empty_and_lng_is_not_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_is_empty_lng_is_not_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_exclude_point_when_lat_is_null_and_lng_is_null_alt_is_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_null_lng_null_alt_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_when_lat_is_null_and_lng_is_null_alt_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_null_lng_null_alt_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_when_lat_is_empty_and_lng_is_empty_alt_is_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_empty_lng_empty_alt_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}
func Test_should_exclude_point_when_lat_is_empty_and_lng_is_empty_alt_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_empty_lng_empty_alt_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Equal(t, *outputUpMessage, inputUpMessage)
}

func Test_should_throw_exception_when_lat_not_null_and_lng_not_null_alt_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_null_lng_not_null_alt_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_include_point_when_lat_not_null_and_lng_not_null_alt_is_not_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_null_lng_not_null_alt_not_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752, 2.1345677},
		EventTime:   inputUpMessage.Time,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_include_point_when_lat_not_null_and_lng_not_null_alt_is_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_null_lng_not_null_alt_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752, 2.1345677},
		EventTime:   inputUpMessage.Time,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_throw_exception_when_lat_not_null_and_lng_not_null_alt_is_array_size_not_1(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_null_lng_not_null_alt_array_size_1.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}

func Test_should_include_point_when_lat_array_and_lng_array_size_1_alt_is_not_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_array_and_lng_array_size_1_alt_not_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752, 2.1345677},
		EventTime:   inputUpMessage.Time,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_throw_exception_when_lat_array_and_lng_array_size_not_1_alt_is_not_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_array_and_lng_array_size_not_1_alt_is_not_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_include_point_when_lat_lng_is_array_and_alt_is_array_equal_size(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_lng_is_array_and_alt_is_array_equal_size.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	resTime, _ := time.Parse(time.RFC3339, "2020-01-01T10:00:05.000Z")
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752, 2.1345677},
		EventTime:   resTime,
	})
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752, 2.1345677},
		EventTime:   resTime,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_throw_exception_when_lat_lng_is_array_and_alt_is_array_unequal_size(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_array_and_lng_array_size_not_1_alt_is_not_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}

func Test_should_throw_exception_when_lat_array_and_lng_array_alt_is_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_array_and_lng_array_alt_is_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_not_null_and_lng_not_null_alt_is_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_null_and_lng_not_null_alt_is_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}
func Test_should_throw_exception_when_lat_not_empty_and_lng_not_empty_alt_is_empty(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("lat_not_empty_and_lng_not_empty_alt_is_empty.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields")
}

func Test_should_include_point_when_coordinate_is_not_null_value_is_not_null_event_time_is_not_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_is_not_null_value_is_not_null_event_time_is_not_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		Value:     "{{packet.message.value}}",
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   inputUpMessage.Time,
		Value:       23.0,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_include_point_when_coordinate_is_array_value_is_array_event_time_is_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_is_array_value_is_array_event_time_is_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		Value:     "{{packet.message.value}}",
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	resTime, _ := time.Parse(time.RFC3339, "2020-01-01T10:00:05.000Z")
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   resTime,
		Value:       23.0,
	})
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   resTime,
		Value:       23.0,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_throw_exception__when_coordinate_is_array_value_is_array_event_time_is_array_unequal_size(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_is_array_value_is_array_event_time_is_array_unequal_size.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
		Value:     "{{packet.message.value}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}
func Test_should_include_point_when_coordinate_with_alt_is_not_null_value_is_not_null_event_time_is_not_null(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_with_alt_is_not_null_value_is_not_null_event_time_is_not_null.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		Value:     "{{packet.message.value}}",
		EventTime: "{{time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752, 2.1345677},
		EventTime:   inputUpMessage.Time,
		Value:       23.0,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_include_point_when_coordinate_with_alt_is_array_value_is_array_event_time_is_array(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_with_alt_is_array_value_is_array_event_time_is_array.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		Value:     "{{packet.message.value}}",
		EventTime: "{{packet.message.time}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	resTime, _ := time.Parse(time.RFC3339, "2020-01-01T10:00:05.000Z")
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   resTime,
		Value:       23.0,
	})
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   resTime,
		Value:       23.0,
	})
	expectedCoordinates := flow.Point{
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"coordinates": expectedCoordinates,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_throw_exception_when_coordinate_with_alt_is_array_value_is_array_event_time_is_array_unequal_size(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("coordinate_with_alt_is_array_value_is_array_event_time_is_array_unequal_size.json")
	coordinates := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsAltitude | [0]}}"},
		EventTime: "{{packet.message.time}}",
		Value:     "{{packet.message.value}}",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates": coordinates,
	}}
	_, err := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then

	assert.Errorf(t, err, "there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields")
}

func Test_should_include_ontology_id_existing_points(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("ontology_id_existing_points.json")
	temperature := ontology.JmesPathPoint{
		OntologyId: "temperature:1:value",
		Value:      "{{packet.message.temperature}}",
		EventTime:  "{{time}}",
		Type_:      "double",
		UnitId:     "Cel",
	}

	// When
	var extractOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathOperation.ApplyUpOperation(&inputUpMessage, &extractOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.6,
		EventTime: inputUpMessage.Time,
	})
	expectedTemperature := flow.Point{
		OntologyId: "temperature:1:value",
		Type_:      "double",
		UnitId:     "Cel",
		Records:    expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"temperature": expectedTemperature,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
