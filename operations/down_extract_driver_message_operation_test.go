package operations

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"ontology-mapping-go-lib/models/ontology"
	"ontology-mapping-go-lib/util"
	"testing"
)

var extractMessageOperation = DownExtractDriverOperation{}

func Test_should_extract_message_downmessage_as_valid_jmespath_expression(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("downmessage_sample.json")
	byteSream1, _ := ioutil.ReadFile("resources/downmessage_as_valid_jmespath_expression.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	commands := map[string]interface{}{
		"myDeviceCommand": data1,
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	outputMessage, _ := extractMessageOperation.ApplyDownOperation(&inputDownMessage, &downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/downmessage_as_valid_jmespath_expression_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Packet.Message = data
	assert.Equal(t, outputMessage, expectedOutputMessage)

}

func Test_should_throw_exception_downmessage_as_invalid_jmespath_expression(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("downmessage_sample.json")
	commands := map[string]interface{}{
		"myDeviceCommand": "{{ command.type }}",
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	_, err := extractMessageOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	assert.EqualError(t, err, "expected object for 'message' but returned value node or null")

}
func Test_should_throw_exception_downmessage_as_some_value(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("downmessage_sample.json")
	commands := map[string]interface{}{
		"myDeviceCommand": "somevalue",
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	_, err := extractMessageOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	assert.EqualError(t, err, "expected object but is a value node")

}
func Test_should_throw_exception_downmessage_as_empty(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("downmessage_sample.json")
	commands := map[string]interface{}{
		"myDeviceCommand": "",
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	_, err := extractMessageOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	assert.EqualError(t, err, "expected object but is a value node")

}
func Test_should_extract_default_jmesexpression_in_case_of_mismatch_message_downmessage(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("downmessage_sample.json")
	byteSream1, _ := ioutil.ReadFile("resources/downmessage_as_default_jmespath_expression.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	commands := map[string]interface{}{
		"default": data1,
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	outputMessage, _ := extractMessageOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/valid_jmesexpression_message_downmessage_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Packet.Message = data
	assert.Equal(t, outputMessage, expectedOutputMessage)

}
func Test_should_not_extract_any_command_in_case_of_message_mismatch_without_default_downmessage(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("downmessage_sample.json")
	byteSream1, _ := ioutil.ReadFile("resources/downmessage_without_default_jmespath_expression.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	commands := map[string]interface{}{
		"setTransmissionFrameStatusPeriod": data1,
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	outputMessage, _ := extractMessageOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then

	assert.Equal(t, *outputMessage, inputDownMessage)

}
func Test_should_extract_command_if_command_is_jmesexpression_downmessage(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("downmessage_sample.json")
	commands := map[string]interface{}{
		"myDeviceCommand": "{{ command.input }}",
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownExtractDriverMessage{Commands: commands}
	outputMessage, _ := extractMessageOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/downmessage_command_is_a_jmesexpression_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Packet.Message = data
	assert.Equal(t, outputMessage, expectedOutputMessage)

}
