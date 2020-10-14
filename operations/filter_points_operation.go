package operations

import (
	"ontology-mapping-go-lib/models/flow"
)
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/util"

type FilterPointsOperation struct {
}

func (filterPointsOperation *FilterPointsOperation) ApplyUpOperation(message *flow.UpMessage, upOperation *ontology.UpOperationInterface) (*flow.UpMessage, error) {
	var retMessage = util.CopyUpMessage(message)
	var points []string
	var upFilterPointsOperation = (*upOperation).(ontology.UpFilterPointsOperation)
	var pointsUpMessage = message.Points
	if pointsUpMessage != nil {
		points = upFilterPointsOperation.Points
		var newPoints = make(map[string]flow.Point)
		if len(points) > 0 {
			for key, element := range pointsUpMessage {
				if util.Contains(points, key) {
					newPoints[key] = element
				}
			}
		}
		retMessage.Points = newPoints
		return retMessage, nil
	}
	return retMessage, nil
}

func (filterPointsOperation *FilterPointsOperation) ApplyDownOperation(message *flow.DownMessage, downOperation *ontology.DownOperationInterface) (*flow.DownMessage, error) {
	return nil, nil
}
