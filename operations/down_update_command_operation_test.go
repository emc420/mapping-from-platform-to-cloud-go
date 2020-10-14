package operations

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"ontology-mapping-go-lib/models/flow"
	"ontology-mapping-go-lib/models/ontology"
	"ontology-mapping-go-lib/util"
	"testing"
	"time"
)

var updateCommandOperation = DownUpdateCommandOperation{}

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

func Test_should_update_command_id_when_there_is_a_match(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	updateCommands := map[string]ontology.UpdateCommand{
		"myDeviceCommand": {Id: "newCommandId"},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	outputMessage, _ := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/update_command_id_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Command.Id = data.(map[string]interface{})["id"].(string)
	expectedOutputMessage.Command.Input = data.(map[string]interface{})["input"]
	assert.Equal(t, outputMessage, expectedOutputMessage)
}

func Test_should_update_command_input_when_there_is_a_match(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	byteSream1, _ := ioutil.ReadFile("resources/update_command_input_jmespath.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	updateCommands := map[string]ontology.UpdateCommand{
		"myDeviceCommand": {Input: data1.(map[string]interface{})},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	outputMessage, _ := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/update_command_input_jmespath_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Command.Id = data.(map[string]interface{})["id"].(string)
	expectedOutputMessage.Command.Input = data.(map[string]interface{})["input"]
	assert.Equal(t, outputMessage, expectedOutputMessage)
}

func Test_should_update_command_id_input_when_there_is_a_match(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	byteSream1, _ := ioutil.ReadFile("resources/update_command_id_input_jmespath.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	updateCommands := map[string]ontology.UpdateCommand{
		"myDeviceCommand": {Id: data1.(map[string]interface{})["id"].(string), Input: data1.(map[string]interface{})["input"]},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	outputMessage, _ := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/update_command_id_input_jmespath_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Command.Id = data.(map[string]interface{})["id"].(string)
	expectedOutputMessage.Command.Input = data.(map[string]interface{})["input"]
	assert.Equal(t, outputMessage, expectedOutputMessage)
}

func Test_should_do_nothing_command_id_is_a_mismatch(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	byteSream1, _ := ioutil.ReadFile("resources/update_command_id_input_jmespath.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	updateCommands := map[string]ontology.UpdateCommand{
		"misMatchCommand": {Id: data1.(map[string]interface{})["id"].(string), Input: data1.(map[string]interface{})["input"]},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	outputMessage, _ := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then

	assert.Equal(t, *outputMessage, inputDownMessage)
}

func Test_should_update_command_id_when_there_is_a_default_mismatch(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	updateCommands := map[string]ontology.UpdateCommand{
		"default": {Id: "default"},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	outputMessage, _ := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/update_command_id_default_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Command.Id = data.(map[string]interface{})["id"].(string)
	expectedOutputMessage.Command.Input = data.(map[string]interface{})["input"]
	assert.Equal(t, outputMessage, expectedOutputMessage)
}

func Test_should_update_command_input_when_there_is_a_default_mismatch(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	updateCommands := map[string]ontology.UpdateCommand{
		"default": {Input: "{{ @.input.prop1}}"},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	outputMessage, _ := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/update_command_default_input_jmespath_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Command.Id = data.(map[string]interface{})["id"].(string)
	expectedOutputMessage.Command.Input = data.(map[string]interface{})["input"]
	assert.Equal(t, outputMessage, expectedOutputMessage)
}
func Test_should_update_command_id_input_when_there_is_a_default_mismatch(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	byteSream1, _ := ioutil.ReadFile("resources/update_command_id_input_default_jmespath.json")
	var data1 interface{}
	_ = json.Unmarshal(byteSream1, &data1)
	updateCommands := map[string]ontology.UpdateCommand{
		"default": {Id: data1.(map[string]interface{})["id"].(string), Input: data1.(map[string]interface{})["input"]},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	outputMessage, _ := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/update_command_id_input_default_jmespath_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Command.Id = data.(map[string]interface{})["id"].(string)
	expectedOutputMessage.Command.Input = data.(map[string]interface{})["input"]
	assert.Equal(t, outputMessage, expectedOutputMessage)
}

func Test_should_update_command_id_when_there_is_a_default_match(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	updateCommands := map[string]ontology.UpdateCommand{
		"default":         {Id: "default"},
		"myDeviceCommand": {Id: "newCommandId"},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	outputMessage, _ := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	byteSream, _ := ioutil.ReadFile("resources/update_command_id_expected.json")
	var data interface{}
	_ = json.Unmarshal(byteSream, &data)
	expectedOutputMessage := util.CopyDownMessage(&inputDownMessage)
	expectedOutputMessage.Command.Id = data.(map[string]interface{})["id"].(string)
	expectedOutputMessage.Command.Input = data.(map[string]interface{})["input"]
	assert.Equal(t, outputMessage, expectedOutputMessage)
}

func Test_should_throw_exception_when_failed_to_extract_jmesexpression(t *testing.T) {
	// Given
	inputDownMessage := buildInputDownMessage("update_command_id.json")
	updateCommands := map[string]ontology.UpdateCommand{
		"default": {Input: "{{ @.input.prop3}}"},
	}
	//When
	var downOpr ontology.DownOperationInterface = ontology.DownUpdateCommand{Commands: updateCommands}
	_, err := updateCommandOperation.ApplyDownOperation(&inputDownMessage,
		&downOpr)
	//Then
	assert.EqualError(t, err, "retrieved value is null or not a map")
}
