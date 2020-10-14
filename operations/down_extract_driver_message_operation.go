package operations

import (
	"ontology-mapping-go-lib/models/flow"
)
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/util"

type DownExtractDriverOperation struct {
}

func (downExtractDriver *DownExtractDriverOperation) ApplyUpOperation(message *flow.UpMessage, upOperation *ontology.UpOperationInterface) (*flow.UpMessage, error) {
	return nil, nil
}

func (downExtractDriver *DownExtractDriverOperation) ApplyDownOperation(message *flow.DownMessage, downOperation *ontology.DownOperationInterface) (*flow.DownMessage, error) {

	var messageJson interface{} = message
	var command = message.Command
	var resultJson interface{}
	var err error
	jmesPathOperation := (*downOperation).(ontology.DownExtractDriverMessage)
	if element, ok := jmesPathOperation.Commands[command.Id]; ok {
		resultJson, err = util.ExtractMessage(messageJson, element)
		if err != nil {
			return nil, err
		}
	}
	if resultJson == nil {
		value, ok := jmesPathOperation.Commands["default"]
		if ok {
			resultJson, err = util.ExtractMessage(messageJson, value)
			if err != nil {
				return nil, err
			}
		}
	}
	var retMessage = util.CopyDownMessage(message)
	retMessage.Packet.Message = resultJson
	return retMessage, nil
}
