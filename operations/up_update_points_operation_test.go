package operations

import (
	"ontology-mapping-go-lib/util"
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/models/flow"

var jmesPathUpdateOperation = UpUpdatePointsOperation{}

func buildInputUpUpdateMessage(inputPoints map[string]flow.Point) flow.UpMessage {
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

func Test_should_update_type_from_matched_points_from_update_points(t *testing.T) {
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
	temperature := ontology.JmesPathUpdatePoint{
		Type_: "string",
	}

	// When
	var updateOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathUpdateOperation.ApplyUpOperation(&inputUpMessage, &updateOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.6,
		EventTime: eventTime,
	})
	expectedTemperature := flow.Point{
		Type_:   "string",
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

func Test_should_update_unit_id_from_matched_points_from_update_points(t *testing.T) {
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
	temperature := ontology.JmesPathUpdatePoint{
		UnitId: "Far",
	}

	// When
	var updateOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathUpdateOperation.ApplyUpOperation(&inputUpMessage, &updateOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.6,
		EventTime: eventTime,
	})
	expectedTemperature := flow.Point{
		Type_:   "double",
		UnitId:  "Far",
		Records: expectedRecords,
	}
	expectedPoints := map[string]flow.Point{
		"temperature": expectedTemperature,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_update_value_from_matched_points_from_update_points(t *testing.T) {
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
	temperature := ontology.JmesPathUpdatePoint{
		Value: "{{floor(@)}}",
	}

	// When
	var updateOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathUpdateOperation.ApplyUpOperation(&inputUpMessage, &updateOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.0,
		EventTime: eventTime,
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

func Test_should_update_event_time_from_matched_points_from_update_points(t *testing.T) {
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
	temperature := ontology.JmesPathUpdatePoint{
		EventTime: "{{@ | date_time_op(@, '+', '5' , 's')}}",
	}

	// When
	var updateOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathUpdateOperation.ApplyUpOperation(&inputUpMessage, &updateOpr)
	// Then
	var expectedRecords []flow.Record
	var expectedEventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:05.000Z")
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.6,
		EventTime: expectedEventTime,
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
func Test_should_update_ontology_id_from_matched_points_from_update_points(t *testing.T) {
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
	temperature := ontology.JmesPathUpdatePoint{
		OntologyId: "temperatureMeasurement:1:temperature",
	}

	// When
	var updateOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathUpdateOperation.ApplyUpOperation(&inputUpMessage, &updateOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.6,
		EventTime: eventTime,
	})
	expectedTemperature := flow.Point{
		OntologyId: "temperatureMeasurement:1:temperature",
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

func Test_should_update_coordinate_from_matched_points_from_update_points(t *testing.T) {
	// Given
	var inputRecords []flow.Record
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	inputRecords = append(inputRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752},
		EventTime:   eventTime,
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
	temperature := ontology.JmesPathUpdatePoint{
		Coordinates: []string{"{{@ | floor(@)}}", "{{@ | floor(@)}}"},
	}

	// When
	var updateOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathUpdateOperation.ApplyUpOperation(&inputUpMessage, &updateOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0, 43.0},
		EventTime:   eventTime,
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

func Test_should_update_coordinate_altitude_from_matched_points_from_update_points(t *testing.T) {
	// Given
	var inputRecords []flow.Record
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	inputRecords = append(inputRecords, flow.Record{
		Coordinates: []float64{7.0586624, 43.6618752, 35.675858},
		EventTime:   eventTime,
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
	temperature := ontology.JmesPathUpdatePoint{
		Coordinates: []string{"{{@ | floor(@)}}", "{{@ | floor(@)}}", "{{@ | floor(@)}}"},
	}

	// When
	var updateOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"temperature": temperature,
	}}
	outputUpMessage, _ := jmesPathUpdateOperation.ApplyUpOperation(&inputUpMessage, &updateOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Coordinates: []float64{7.0, 43.0, 35.0},
		EventTime:   eventTime,
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

func Test_should_not_update_when_there_is_a_mismatch_points_from_update_points(t *testing.T) {
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
	humidity := ontology.JmesPathUpdatePoint{
		UnitId: "%RH",
	}

	// When
	var updateOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"humidity": humidity,
	}}
	outputUpMessage, _ := jmesPathUpdateOperation.ApplyUpOperation(&inputUpMessage, &updateOpr)
	// Then
	var expectedRecords []flow.Record
	expectedRecords = append(expectedRecords, flow.Record{
		Value:     22.6,
		EventTime: eventTime,
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
