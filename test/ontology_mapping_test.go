package test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"ontology-mapping-go-lib/models/flow"
	"ontology-mapping-go-lib/models/ontology"
	"ontology-mapping-go-lib/operations"
	"ontology-mapping-go-lib/util"
	"testing"
	"time"
)

var operationFactory = operations.OperationFactory{}
var oprServ = operations.OperationService{Factory: operationFactory}

func buildInputUpMessage(inputMessageFile string, inputPoints map[string]flow.Point) flow.UpMessage {
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
		Packet: &flow.MessagePacket{
			Type_:   "",
			Raw:     "",
			Message: nil,
		},
		Points: inputPoints,
	}
	byteSream, _ := ioutil.ReadFile("resources/" + inputMessageFile)
	_ = json.Unmarshal(byteSream, &message.Packet.Message)
	return message
}

func buildInputDownMessage(inputMessageFile string) flow.DownMessage {
	var eventTime, _ = time.Parse(time.RFC3339, "2020-01-01T10:00:00.000Z")
	byteSream, _ := ioutil.ReadFile("resources/" + inputMessageFile)
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	var message = flow.DownMessage{
		Id:   "00000000-000000-00000-000000000",
		Time: eventTime,
		SubAccount: &flow.Account{
			Id:      "sub1",
			RealmId: "realm1",
		},
		Origin: &flow.DownOrigin{
			Type_: "binder",
			Id:    "tpw",
			Time:  time.Now(),
		},
		Type_: "deviceDownlink",
		Command: &flow.Command{
			Id:    data.(map[string]interface{})["id"].(string),
			Input: data.(map[string]interface{})["input"],
		},
		Subscriber: &flow.Subscriber{
			Id:      "sub1",
			RealmId: "realm1",
		},
		Thing: &flow.Thing{
			Key: "lora:0102030405060708",
		},
		Packet: &flow.MessagePacket{
			Type_:   "",
			Raw:     "",
			Message: nil,
		},
	}
	return message
}

func Test_should_extract_elsys_points(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("elsys.json", nil)
	var operation operations.OperationsUpSerDer
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.temperature}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	humidity := ontology.JmesPathPoint{
		Value:     "{{packet.message.humidity}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "%RH",
	}
	light := ontology.JmesPathPoint{
		Value:     "{{packet.message.light}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "lx",
	}
	var upOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature": temperature,
		"humidity":    humidity,
		"light":       light,
	}}
	operation.Operations = append(operation.Operations, upOpr)
	// When
	outputUpMessage, _ := oprServ.ApplyUpOperations(&inputUpMessage, &operation)
	// Then
	expectedTemperature := flow.Point{
		Type_:  "double",
		UnitId: "Cel",
		Records: []flow.Record{{
			Value:     22.6,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedHumidity := flow.Point{
		Type_:  "double",
		UnitId: "%RH",
		Records: []flow.Record{{
			Value:     41.0,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedLight := flow.Point{
		Type_:  "double",
		UnitId: "lx",
		Records: []flow.Record{{
			Value:     39.0,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedPoints := map[string]flow.Point{
		"temperature": expectedTemperature,
		"humidity":    expectedHumidity,
		"light":       expectedLight,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_extract_sensing_labs_points(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("sensing_labs.json", map[string]flow.Point{})
	var operation operations.OperationsUpSerDer
	temperature := ontology.JmesPathPoint{
		Value:     "{{packet.message.measures[?id == 'temperature'].value}}",
		EventTime: "{{packet.message.measures[?id == 'temperature'].time}}",
		Type_:     "double",
		UnitId:    "Cel",
	}
	batteryCurrentLevel := ontology.JmesPathPoint{
		Value:     "{{packet.message.measures[?id == 'battery_current_level'].value}}",
		EventTime: "{{packet.message.measures[?id == 'battery_current_level'].time}}",
		Type_:     "double",
		UnitId:    "%RH",
	}
	var upOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"temperature":           temperature,
		"battery_current_level": batteryCurrentLevel,
	}}
	operation.Operations = append(operation.Operations, upOpr)
	// When
	outputUpMessage, _ := oprServ.ApplyUpOperations(&inputUpMessage, &operation)
	// Then
	resTime1, _ := time.Parse(time.RFC3339, "2020-02-06T09:14:05.688Z")
	resTime2, _ := time.Parse(time.RFC3339, "2020-02-06T09:15:45.688Z")
	resTime3, _ := time.Parse(time.RFC3339, "2020-02-06T09:17:25.688Z")
	resTime4, _ := time.Parse(time.RFC3339, "2020-02-06T09:18:15.688Z")
	expectedTemperature := flow.Point{
		Type_:  "double",
		UnitId: "Cel",
		Records: []flow.Record{{
			Value:     11.625,
			EventTime: resTime1,
		}, {
			Value:     11.625,
			EventTime: resTime2,
		}, {
			Value:     11.6875,
			EventTime: resTime3,
		}},
	}
	expectedBatteryCurrentLevel := flow.Point{
		Type_:  "double",
		UnitId: "%RH",
		Records: []flow.Record{{
			Value:     30.0,
			EventTime: resTime4,
		}},
	}
	expectedPoints := map[string]flow.Point{
		"temperature":           expectedTemperature,
		"battery_current_level": expectedBatteryCurrentLevel,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_extract_nke_points(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("nke.json", map[string]flow.Point{})
	var operation operations.OperationsUpSerDer
	energy := ontology.JmesPathPoint{
		Value:     "{{packet.message | to_array(@) | [?CommandID == 'ReportAttributes' && AttributeID == 'Attribute_0'].Data.TICFieldList.BBRHCJB | @[0]}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "W",
	}
	var upOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"energy": energy,
	}}
	operation.Operations = append(operation.Operations, upOpr)
	// When
	outputUpMessage, _ := oprServ.ApplyUpOperations(&inputUpMessage, &operation)
	// Then
	expectedEnergy := flow.Point{
		Type_:  "double",
		UnitId: "W",
		Records: []flow.Record{{
			Value:     123456789.0,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedPoints := map[string]flow.Point{
		"energy": expectedEnergy,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_extract_abeeway_points(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("abeeway.json", map[string]flow.Point{})
	var operation operations.OperationsUpSerDer
	coordinate := ontology.JmesPathPoint{
		Coordinates: []string{"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLongitude| [0]}}",
			"{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].gpsLatitude | [0]}}"},
		EventTime:  "{{@ | date_time_op(time, '+', packet.message.age, 's')}}",
		OntologyId: "Geolocation:1:coordinates",
	}
	temperatureMeasure := ontology.JmesPathPoint{
		OntologyId: "TemperatureMeasurement:1:measuredValue",
		Value:      "{{packet.message.temperatureMeasure}}",
		EventTime:  "{{time}}",
		Type_:      "double",
		UnitId:     "Cel",
	}
	rawPositionType := ontology.JmesPathPoint{
		Value:     "{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE'].rawPositionType | [0]}}",
		EventTime: "{{@ | date_time_op(time, '+', packet.message.age, 's')}}",
		Type_:     "string",
	}
	batteryStatus := ontology.JmesPathPoint{
		Value:     "{{packet.message.batteryStatus}}",
		EventTime: "{{time}}",
		Type_:     "string",
	}
	horizontalAccuracy := ontology.JmesPathPoint{
		OntologyId: "Geolocation:1:horizontalAccuracy",
		Value:      "{{packet.message | to_array(@) | [?messageType == 'POSITION_MESSAGE' && rawPositionType == 'GPS'].horizontalAccuracy | [0]}}",
		EventTime:  "{{@ | date_time_op(time, '+', packet.message.age, 's')}}",
		Type_:      "double",
		UnitId:     "m",
	}
	trackingMode := ontology.JmesPathPoint{
		Value:     "{{packet.message.trackingMode}}",
		EventTime: "{{time}}",
		Type_:     "string",
	}
	batteryLevel := ontology.JmesPathPoint{
		OntologyId: "PowerConfiguration:1:batteryPercentageRemaining",
		Value:      "{{packet.message.batteryLevel}}",
		EventTime:  "{{time}}",
		Type_:      "double",
		UnitId:     "%",
	}
	sosFlag := ontology.JmesPathPoint{
		Value:     "{{packet.message.sosFlag | to_boolean(@)}}",
		EventTime: "{{time}}",
		Type_:     "boolean",
	}
	dynamicMotionState := ontology.JmesPathPoint{
		Value:     "{{packet.message.dynamicMotionState}}",
		EventTime: "{{time}}",
		Type_:     "string",
	}
	var upOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"coordinates":        coordinate,
		"horizontalAccuracy": horizontalAccuracy,
		"rawPositionType":    rawPositionType,
		"trackingMode":       trackingMode,
		"batteryLevel":       batteryLevel,
		"batteryStatus":      batteryStatus,
		"temperatureMeasure": temperatureMeasure,
		"sosFlag":            sosFlag,
		"dynamicMotionState": dynamicMotionState,
	}}
	operation.Operations = append(operation.Operations, upOpr)
	// When
	outputUpMessage, _ := oprServ.ApplyUpOperations(&inputUpMessage, &operation)
	// Then
	resTime1, _ := time.Parse(time.RFC3339, "2020-01-01T10:00:05.000Z")
	expectedCoordinates := flow.Point{
		OntologyId: "Geolocation:1:coordinates",
		Records: []flow.Record{{
			Coordinates: []float64{7.0586624, 43.6618752},
			EventTime:   resTime1,
		}},
	}
	expectedHorizontalAccuracy := flow.Point{
		OntologyId: "Geolocation:1:horizontalAccuracy",
		Type_:      "double",
		UnitId:     "m",
		Records: []flow.Record{{
			Value:     19.0,
			EventTime: resTime1,
		}},
	}
	expectedRawPositionType := flow.Point{
		Type_: "string",
		Records: []flow.Record{{
			Value:     "GPS",
			EventTime: resTime1,
		}},
	}
	expectedTrackingMode := flow.Point{
		Type_: "string",
		Records: []flow.Record{{
			Value:     "PERMANENT_TRACKING",
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedBatteryLevel := flow.Point{
		OntologyId: "PowerConfiguration:1:batteryPercentageRemaining",
		Type_:      "double",
		UnitId:     "%",
		Records: []flow.Record{{
			Value:     95.0,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedBatteryStatus := flow.Point{
		Type_: "string",
		Records: []flow.Record{{
			Value:     "OPERATING",
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedTemperatureMeasure := flow.Point{
		OntologyId: "TemperatureMeasurement:1:measuredValue",
		Type_:      "double",
		UnitId:     "Cel",
		Records: []flow.Record{{
			Value:     25.8,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedSosFlag := flow.Point{
		Type_: "boolean",
		Records: []flow.Record{{
			Value:     true,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedDynamicMotionState := flow.Point{
		Type_: "string",
		Records: []flow.Record{{
			Value:     "STATIC",
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedPoints := map[string]flow.Point{
		"coordinates":        expectedCoordinates,
		"horizontalAccuracy": expectedHorizontalAccuracy,
		"rawPositionType":    expectedRawPositionType,
		"trackingMode":       expectedTrackingMode,
		"batteryLevel":       expectedBatteryLevel,
		"batteryStatus":      expectedBatteryStatus,
		"temperatureMeasure": expectedTemperatureMeasure,
		"sosFlag":            expectedSosFlag,
		"dynamicMotionState": expectedDynamicMotionState,
	}

	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}

func Test_should_extract_adeunis_downmessage(t *testing.T) {
	// Given
	var operation operations.OperationsDownSerDer
	inputDownMessage := buildInputDownMessage("adeunis.json")
	byteSream1, _ := ioutil.ReadFile("resources/adeunis_jmespath.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	commands := map[string]interface{}{
		"setTransmissionFrameStatusPeriod": data1,
	}
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	operation.Operations = append(operation.Operations, downOpr)
	//When
	outputMessage, _ := oprServ.ApplyDownOperations(&inputDownMessage, &operation)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/adeunis_expected_output.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Packet.Message = data
	assert.Equal(t, outputMessage, expectedOutputMessage)

}

func Test_should_extract_theoritical_downmessage(t *testing.T) {
	// Given
	var operation operations.OperationsDownSerDer
	inputDownMessage := buildInputDownMessage("theoritical.json")
	commands := map[string]interface{}{
		"myDeviceCommand": "{{ command | add_property(@.input, 'type', @.id) }}",
	}
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	operation.Operations = append(operation.Operations, downOpr)
	//When
	outputMessage, _ := oprServ.ApplyDownOperations(&inputDownMessage, &operation)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/theoritical_expected_output.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Packet.Message = data
	assert.Equal(t, outputMessage, expectedOutputMessage)

}
func Test_should_execute_both_jmespath_filter_operations(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("nke.json", map[string]flow.Point{})
	var operation operations.OperationsUpSerDer
	energy := ontology.JmesPathPoint{
		Value:     "{{packet.message | to_array(@) | [?CommandID == 'ReportAttributes' && AttributeID == 'Attribute_0'].Data.TICFieldList.BBRHCJB | @[0]}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "W",
	}
	var upOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"energy": energy,
	}}
	operation.Operations = append(operation.Operations, upOpr)
	upOpr = ontology.UpFilterOperation{KeepDeviceUplink: true}
	operation.Operations = append(operation.Operations, upOpr)
	// When
	outputUpMessage, _ := oprServ.ApplyUpOperations(&inputUpMessage, &operation)
	// Then
	expectedEnergy := flow.Point{
		Type_:  "double",
		UnitId: "W",
		Records: []flow.Record{{
			Value:     123456789.0,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedPoints := map[string]flow.Point{
		"energy": expectedEnergy,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_extract_final_message_after_jmespath_filter_point_operations(t *testing.T) {
	// Given
	inputUpMessage := buildInputUpMessage("nke.json", map[string]flow.Point{})
	var operation operations.OperationsUpSerDer
	energy := ontology.JmesPathPoint{
		Value:     "{{packet.message | to_array(@) | [?CommandID == 'ReportAttributes' && AttributeID == 'Attribute_0'].Data.TICFieldList.BBRHCJB | @[0]}}",
		EventTime: "{{time}}",
		Type_:     "double",
		UnitId:    "W",
	}
	var upOpr ontology.UpOperationInterface = ontology.UpExtractPoints{Points: map[string]ontology.JmesPathPoint{
		"energy": energy,
	}}
	operation.Operations = append(operation.Operations, upOpr)
	upOpr = ontology.UpFilterPointsOperation{Points: []string{"energy"}}
	operation.Operations = append(operation.Operations, upOpr)
	// When
	outputUpMessage, _ := oprServ.ApplyUpOperations(&inputUpMessage, &operation)
	// Then
	expectedEnergy := flow.Point{
		Type_:  "double",
		UnitId: "W",
		Records: []flow.Record{{
			Value:     123456789.0,
			EventTime: inputUpMessage.Time,
		}},
	}
	expectedPoints := map[string]flow.Point{
		"energy": expectedEnergy,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_update_sensing_labs_points(t *testing.T) {
	// Given
	resTime1, _ := time.Parse(time.RFC3339, "2020-02-06T09:14:05.688Z")
	resTime2, _ := time.Parse(time.RFC3339, "2020-02-06T09:15:45.688Z")
	resTime3, _ := time.Parse(time.RFC3339, "2020-02-06T09:17:25.688Z")
	resTime4, _ := time.Parse(time.RFC3339, "2020-02-06T09:18:15.688Z")
	inputTemperature := flow.Point{
		Type_:  "double",
		UnitId: "Cel",
		Records: []flow.Record{{
			Value:     11.625,
			EventTime: resTime1,
		}, {
			Value:     11.625,
			EventTime: resTime2,
		}, {
			Value:     11.6875,
			EventTime: resTime3,
		}},
	}
	inputBatteryCurrentLevel := flow.Point{
		Type_:  "double",
		UnitId: "%RH",
		Records: []flow.Record{{
			Value:     30.5656,
			EventTime: resTime4,
		}},
	}
	inputUpMessage := buildInputUpMessage("sensing_labs.json", map[string]flow.Point{
		"temperature":           inputTemperature,
		"battery_current_level": inputBatteryCurrentLevel,
	})
	var operation operations.OperationsUpSerDer
	temperature := ontology.JmesPathUpdatePoint{
		Value:  "{{@ | floor(@)}}",
		Type_:  "int64",
		UnitId: "Far",
	}
	batteryCurrentLevel := ontology.JmesPathUpdatePoint{
		Value: "{{@ | floor(@)}}",
		Type_: "int64",
	}
	var upOpr ontology.UpOperationInterface = ontology.UpUpdatePoints{Points: map[string]ontology.JmesPathUpdatePoint{
		"temperature":           temperature,
		"battery_current_level": batteryCurrentLevel,
	}}
	operation.Operations = append(operation.Operations, upOpr)
	// When
	outputUpMessage, _ := oprServ.ApplyUpOperations(&inputUpMessage, &operation)
	// Then
	expectedTemperature := flow.Point{
		Type_:  "int64",
		UnitId: "Far",
		Records: []flow.Record{{
			Value:     11.0,
			EventTime: resTime1,
		}, {
			Value:     11.0,
			EventTime: resTime2,
		}, {
			Value:     11.0,
			EventTime: resTime3,
		}},
	}
	expectedBatteryCurrentLevel := flow.Point{
		Type_:  "int64",
		UnitId: "%RH",
		Records: []flow.Record{{
			Value:     30.0,
			EventTime: resTime4,
		}},
	}
	expectedPoints := map[string]flow.Point{
		"temperature":           expectedTemperature,
		"battery_current_level": expectedBatteryCurrentLevel,
	}
	expectedOutputMessage := util.CopyUpMessage(&inputUpMessage)
	expectedOutputMessage.Points = expectedPoints
	assert.Equal(t, outputUpMessage, expectedOutputMessage)
}
func Test_should_update_commands(t *testing.T) {
	// Given
	var operation operations.OperationsDownSerDer
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	byteSream1, _ := ioutil.ReadFile("resources/update_command_id_input_jmespath.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	updateCommands := map[string]ontology.UpdateCommand{
		"myDeviceCommand": {Id: data1.(map[string]interface{})["id"].(string), Input: data1.(map[string]interface{})["input"]},
	}
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	operation.Operations = append(operation.Operations, downOpr)
	//When
	outputMessage, _ := oprServ.ApplyDownOperations(&inputDownMessage, &operation)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/update_command_id_input_jmespath_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Command.Id = data.(map[string]interface{})["id"].(string)
	expectedOutputMessage.Command.Input = data.(map[string]interface{})["input"]
	assert.Equal(t, outputMessage, expectedOutputMessage)
}
