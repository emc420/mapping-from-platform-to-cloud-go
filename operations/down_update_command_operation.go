package operations

import (
	"ontology-mapping-go-lib/models/flow"
)
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/util"

type DownUpdateCommandOperation struct {
}

func (downUpdateCommand *DownUpdateCommandOperation) ApplyUpOperation(message *flow.UpMessage, upOperation *ontology.UpOperationInterface) (*flow.UpMessage, error) {
	return nil, nil
}

func (downUpdateCommand *DownUpdateCommandOperation) ApplyDownOperation(message *flow.DownMessage, downOperation *ontology.DownOperationInterface) (*flow.DownMessage, error) {
	var retMessage = util.CopyDownMessage(message)
	jmesPathOperation := (*downOperation).(ontology.DownUpdateCommand)
	command := message.Command
	var err error
	if element, ok := (jmesPathOperation.Commands)[command.Id]; ok {
		if len(element.Id) > 0 {
			retMessage.Command.Id = element.Id
		}
		if element.Input != nil {
			retMessage.Command.Input, err = util.ExtractCommands(command, element.Input)
			if err != nil {
				return nil, err
			}
		}
		return retMessage, nil
	}
	if value, ok := (jmesPathOperation.Commands)["default"]; ok && len(command.Id) > 0 {
		if len(value.Id) > 0 {
			retMessage.Command.Id = value.Id
		}
		if value.Input != nil {
			retMessage.Command.Input, err = util.ExtractCommands(command, value.Input)
			if err != nil {
				return nil, err
			}
		}
	}
	return retMessage, nil
}
