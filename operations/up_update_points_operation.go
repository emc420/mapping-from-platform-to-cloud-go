package operations

import (
	"errors"
	"ontology-mapping-go-lib/models/flow"
	"strings"
	"time"
)
import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/util"

type UpUpdatePointsOperation struct {
}

func (updatePointsOperation *UpUpdatePointsOperation) ApplyUpOperation(message *flow.UpMessage, upOperation *ontology.UpOperationInterface) (*flow.UpMessage, error) {
	var retMessage = util.CopyUpMessage(message)
	var existingPoints = message.Points
	var newPoints = make(map[string]flow.Point)
	var jmesPathOperation = (*upOperation).(ontology.UpUpdatePoints)
	var err error
	for key, element := range existingPoints {
		var jmesPathPoints = jmesPathOperation.Points
		value, ok := jmesPathPoints[key]
		if ok {
			var isValue = false
			var isAltitude = false
			var isCoordinate = false
			var records = element.Records
			var values []interface{}
			var latitudes []interface{}
			var longitudes []interface{}
			var altitudes []interface{}
			var eventTimes []interface{}
			for i := 0; i < len(records); i++ {
				if records[i].Value != nil {
					var retVal interface{}
					retVal, err = getValue(records[i].Value, value.Value)
					if err != nil {
						return nil, err
					}
					values = append(values, retVal)
					isValue = true
				}
				var retTime interface{}
				retTime, err = getEventTime(records[i].EventTime, value.EventTime)
				if err != nil {
					return nil, err
				}
				eventTimes = append(eventTimes, retTime)
				if records[i].Coordinates != nil && len(records[i].Coordinates) > 0 {
					isCoordinate = true
					var lng interface{}
					var lat interface{}
					var alt interface{}
					if value.Coordinates != nil {
						var resLng interface{}
						var resLat interface{}
						var resAlt interface{}
						if len(value.Coordinates) == 2 {
							lng = records[i].Coordinates[0]
							lat = records[i].Coordinates[1]
							resLng, err = util.RetrieveValues(value.Coordinates[0], &lng)
							resLat, err = util.RetrieveValues(value.Coordinates[1], &lat)
							if err != nil {
								return nil, err
							}
							longitudes = append(longitudes, resLng)
							latitudes = append(latitudes, resLat)
						} else if len(value.Coordinates) == 3 {
							lng = records[i].Coordinates[0]
							lat = records[i].Coordinates[1]
							alt = records[i].Coordinates[2]
							resLng, err = util.RetrieveValues(value.Coordinates[0], &lng)
							resLat, err = util.RetrieveValues(value.Coordinates[1], &lat)
							resAlt, err = util.RetrieveValues(value.Coordinates[2], &alt)
							if err != nil {
								return nil, err
							}
							longitudes = append(longitudes, resLng)
							latitudes = append(latitudes, resLat)
							altitudes = append(altitudes, resAlt)
							isAltitude = true
						} else {
							err = errors.New("invalid 'coordinate' length, it must be 2 or 3 %s" + key)
							return nil, err
						}
					} else {
						lng = records[i].Coordinates[0]
						longitudes = append(longitudes, lng)
						lat = records[i].Coordinates[1]
						latitudes = append(latitudes, lat)
						if len(records[i].Coordinates) == 3 {
							alt = records[i].Coordinates[2]
							altitudes = append(altitudes, alt)
							isAltitude = true
						}
					}
				}
			}
			var params = util.PointParams{Values: values, EventTime: eventTimes, Longitude: longitudes, Latitude: latitudes, Altitude: altitudes, IsAltitude: isAltitude, IsValue: isValue, IsCoordinate: isCoordinate}
			var newRecords []flow.Record
			newRecords, err = util.ExtractRecords(params, key)
			if err != nil {
				return nil, err
			}
			var assignpointType flow.PointType
			var pointType string
			if len(value.Type_) > 0 {
				pointType, _ = value.Type_.GetValue(value.Type_)
				assignpointType = flow.PointType(pointType)
			} else {
				pointType, _ = element.Type_.GetValue(element.Type_)
				assignpointType = flow.PointType(pointType)
			}
			var ontologyId string
			var unitId string
			if len(value.OntologyId) > 0 {
				ontologyId = value.OntologyId
			} else {
				ontologyId = element.OntologyId
			}
			if len(value.UnitId) > 0 {
				unitId = value.UnitId
			} else {
				unitId = element.UnitId
			}
			if newRecords != nil {
				newPoints[key] = flow.Point{OntologyId: ontologyId, UnitId: unitId, Records: newRecords, Type_: assignpointType}
			}

		} else {
			newPoints[key] = element
		}

	}
	if len(newPoints) > 0 {
		retMessage.Points = newPoints
	}
	return retMessage, nil
}

func (updatePointsOperation *UpUpdatePointsOperation) ApplyDownOperation(message *flow.DownMessage, downOperation *ontology.DownOperationInterface) (*flow.DownMessage, error) {
	return nil, nil
}

func getValue(recordVal interface{}, jmesPathValue string) (interface{}, error) {
	if &jmesPathValue != nil && strings.Contains(jmesPathValue, "{{") && strings.Contains(jmesPathValue, "}}") {
		return util.RetrieveValues(jmesPathValue, &recordVal)
	} else if &jmesPathValue != nil && len(jmesPathValue) > 0 {
		var retVal interface{} = jmesPathValue
		return retVal, nil
	} else {
		return recordVal, nil
	}
}
func getEventTime(recordTime time.Time, jmesPathEventTime string) (interface{}, error) {
	var recTime interface{} = recordTime
	if &jmesPathEventTime != nil && strings.Contains(jmesPathEventTime, "{{") && strings.Contains(jmesPathEventTime, "}}") {
		return util.RetrieveValues(jmesPathEventTime, &recTime)
	} else if &jmesPathEventTime != nil && len(jmesPathEventTime) > 0 {
		recTime = jmesPathEventTime
		return recTime, nil
	} else {
		recTime = recordTime
		return recordTime, nil
	}
}
