package operations

import (
	"errors"
	"ontology-mapping-go-lib/models/flow"
)
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/util"

type UpExtractPointsOperation struct {
}

func (extractPoints *UpExtractPointsOperation) ApplyUpOperation(message *flow.UpMessage, upOperation *ontology.UpOperationInterface) (*flow.UpMessage, error) {
	var newPoints = make(map[string]flow.Point)
	var retMessage = util.CopyUpMessage(message)
	if len(retMessage.Points) > 0 {
		newPoints = retMessage.Points
	}
	var jmesPathOperation = (*upOperation).(ontology.UpExtractPoints)
	var messageJson interface{} = message
	var err error
	for key, element := range jmesPathOperation.Points {
		var isValue = false
		var values interface{}
		if len(element.Value) > 0 {
			values, err = util.RetrieveValues(element.Value, &messageJson)
			if err != nil {
				return nil, err
			}
			isValue = true
		}
		var eventTime interface{}
		eventTime, err = util.RetrieveValues(element.EventTime, &messageJson)
		if err != nil {
			return nil, err
		}
		var longitude interface{}
		var latitude interface{}
		var altitude interface{}
		var isAltitude = false
		var isCoordinate = false
		if element.Coordinates != nil {
			isCoordinate = true
			if len(element.Coordinates) == 2 {
				longitude, err = util.RetrieveValues(element.Coordinates[0], &messageJson)
				latitude, err = util.RetrieveValues(element.Coordinates[1], &messageJson)
				if err != nil {
					return nil, err
				}
			} else if len(element.Coordinates) == 3 {
				longitude, err = util.RetrieveValues(element.Coordinates[0], &messageJson)
				latitude, err = util.RetrieveValues(element.Coordinates[1], &messageJson)
				altitude, err = util.RetrieveValues(element.Coordinates[2], &messageJson)
				isAltitude = true
				if err != nil {
					return nil, err
				}
			} else {
				err = errors.New("invalid 'coordinate' length, it must be 2 or 3 %s" + key)
				return nil, err
			}
		}
		var params = util.PointParams{Values: values, EventTime: eventTime, Longitude: longitude, Latitude: latitude, Altitude: altitude, IsAltitude: isAltitude, IsValue: isValue, IsCoordinate: isCoordinate}
		var records []flow.Record
		records, err = util.ExtractRecords(params, key)
		if err != nil {
			return nil, err
		}
		var assignpointType flow.PointType
		if &element.Type_ != nil {
			var pointType, _ = element.Type_.GetValue(element.Type_)
			assignpointType = flow.PointType(pointType)
		}
		if records != nil {
			newPoints[key] = flow.Point{OntologyId: element.OntologyId, UnitId: element.UnitId, Records: records, Type_: assignpointType}
		}
	}
	if len(newPoints) > 0 {
		retMessage.Points = newPoints
	}
	return retMessage, nil
}

func (extractPoints *UpExtractPointsOperation) ApplyDownOperation(message *flow.DownMessage, downOperation *ontology.DownOperationInterface) (*flow.DownMessage, error) {
	return nil, nil
}
