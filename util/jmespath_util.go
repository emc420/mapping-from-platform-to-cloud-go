package util

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"
)
import "ontology-mapping-go-lib/models/flow"
import "git.int.actility.com/Thingpark-X/go-jmespath"
import "strings"

func RetrieveValues(jmesExpression string, message *interface{}) (interface{}, error) {
	if jmesExpression == "" {
		return nil, nil
	}
	if strings.Contains(jmesExpression, "{{") && strings.Contains(jmesExpression, "}}") {
		var str = strings.Replace(strings.Replace(jmesExpression, "{{", "", -1), "}}", "", -1)
		var searchResult, err = jmespath.Search(str, *message)
		if err != nil {
			return nil, err
		}
		if searchResult != nil {
			return searchResult, nil
		}
	}
	return nil, nil
}

func ExtractRecords(pointParams PointParams, point string) ([]flow.Record, error) {
	var recordList []flow.Record
	var value = toArray(pointParams.Values)
	var lng = toArray(pointParams.Longitude)
	var lat = toArray(pointParams.Latitude)
	var alt = toArray(pointParams.Altitude)
	var eventTime = toArray(pointParams.EventTime)
	var err error = nil
	_, err = checkCardinality(pointParams, value, lng, lat, alt, eventTime, point)
	if err == nil {
		var validTime time.Time
		var la, lg, al float64
		if len(value) > 0 && len(lng) > 0 && len(alt) > 0 {
			for i := 0; i < len(eventTime); i++ {
				validTime, err = toTime(eventTime[i])
				la, err = toDouble(lat[i])
				lg, err = toDouble(lng[i])
				al, err = toDouble(alt[i])
				if err != nil {
					return nil, err
				}
				var record = flow.Record{Value: value[i], EventTime: validTime, Coordinates: []float64{lg, la, al}}
				recordList = append(recordList, record)
			}
		} else if len(value) > 0 && len(lng) > 0 {
			for i := 0; i < len(eventTime); i++ {
				validTime, err = toTime(eventTime[i])
				la, err = toDouble(lat[i])
				lg, err = toDouble(lng[i])
				if err != nil {
					return nil, err
				}
				var record = flow.Record{Value: value[i], EventTime: validTime, Coordinates: []float64{lg, la}}
				recordList = append(recordList, record)
			}
		} else if len(lng) > 0 && len(alt) > 0 {
			for i := 0; i < len(eventTime); i++ {
				validTime, err = toTime(eventTime[i])
				la, err = toDouble(lat[i])
				lg, err = toDouble(lng[i])
				al, err = toDouble(alt[i])
				if err != nil {
					return nil, err
				}
				var record = flow.Record{EventTime: validTime, Coordinates: []float64{lg, la, al}}
				recordList = append(recordList, record)
			}
		} else if len(lng) > 0 {
			for i := 0; i < len(eventTime); i++ {
				validTime, err = toTime(eventTime[i])
				la, err = toDouble(lat[i])
				lg, err = toDouble(lng[i])
				if err != nil {
					return nil, err
				}
				var record = flow.Record{EventTime: validTime, Coordinates: []float64{lg, la}}
				recordList = append(recordList, record)
			}
		} else if len(value) > 0 {
			for i := 0; i < len(eventTime); i++ {
				validTime, err = toTime(eventTime[i])
				if err != nil {
					return nil, err
				}
				var record = flow.Record{Value: value[i], EventTime: validTime}
				recordList = append(recordList, record)
			}
		}
		return recordList, err
	} else {
		return nil, err
	}

}

func ExtractMessage(message interface{}, operation interface{}) (interface{}, error) {
	var err error
	if reflect.TypeOf(operation).Kind() == reflect.String &&
		strings.Contains(operation.(string), "{{") && strings.Contains(operation.(string), "}}") {
		var retrievedJmesValue interface{}
		retrievedJmesValue, err = RetrieveValues(operation.(string), &message)
		if err != nil {
			return nil, err
		}
		if retrievedJmesValue != nil && reflect.TypeOf(retrievedJmesValue).Kind() == reflect.Map {
			return retrievedJmesValue, nil
		} else {
			return nil, errors.New("expected object for 'message' but returned value node or null")
		}
	}
	if reflect.TypeOf(operation).Kind() != reflect.Map {
		return nil, errors.New("expected object but is a value node")
	}
	return extractMessageRecursion(message, operation)
}

func ExtractCommands(message interface{}, operation interface{}) (interface{}, error) {
	var err error
	if reflect.TypeOf(operation).Kind() == reflect.String &&
		strings.Contains(operation.(string), "{{") && strings.Contains(operation.(string), "}}") {
		var retrievedJmesValue interface{}
		retrievedJmesValue, err = RetrieveValues(operation.(string), &message)
		if err != nil {
			return nil, err
		}
		if retrievedJmesValue != nil {
			return retrievedJmesValue, nil
		} else {
			return nil, errors.New("retrieved value is null or not a map")
		}
	}
	if reflect.TypeOf(operation).Kind() != reflect.Map {
		return operation, nil
	}
	return extractMessageRecursion(message, operation)
}

func extractMessageRecursion(message interface{}, operation interface{}) (interface{}, error) {
	var returnJson = make(map[string]interface{})
	fields := operation.(map[string]interface{})
	var err error
	for key, element := range fields {
		if reflect.TypeOf(element).Kind() == reflect.Map {
			returnJson[key], err = extractMessageRecursion(message, element)
			if err != nil {
				return nil, err
			}
		} else if strings.Contains(element.(string), "{{") && strings.Contains(element.(string), "}}") {
			var value interface{}
			value, err = RetrieveValues(element.(string), &message)
			if err != nil {
				return nil, err
			}
			if value != nil {
				returnJson[key] = value
			} else {
				return nil, errors.New("nothing could be extracted from the jmes expression" + key)
			}
		} else {
			returnJson[key] = element
		}
	}
	return returnJson, err
}

func checkCardinality(pointParams PointParams, value []interface{}, lng []interface{}, lat []interface{}, alt []interface{}, eventTime []interface{}, point string) (int, error) {
	if pointParams.IsValue {
		if len(value) > 0 && len(value) != len(eventTime) {
			return 0, errors.New("there is a mismatch in cardinality for 'value' and 'eventTime' " + point)
		}
	}
	if pointParams.IsCoordinate {
		if len(lng) > 0 && len(lng) != len(eventTime) {
			return 0, errors.New("there is a mismatch in cardinality between 'coordinates' and 'eventTime' fields " + point)

		} else if len(lng) != len(lat) || (pointParams.IsAltitude && len(lng) != len(alt)) {
			return 0, errors.New("there is a mismatch in cardinality between 'latitude' and 'longitude' and 'altitude' fields " + point)
		}
	}
	if pointParams.IsValue && pointParams.IsCoordinate {
		if len(value) != len(lng) {
			return 0, errors.New("there is a mismatch in cardinality for 'value' and 'coordinates' " + point)
		}
	}
	return 0, nil
}

func toArray(point interface{}) []interface{} {
	var result []interface{}
	if point == nil {
		return nil
	} else if reflect.TypeOf(point).Kind() == reflect.Array || reflect.TypeOf(point).Kind() == reflect.Slice {
		return point.([]interface{})
	} else {
		return append(result, point)
	}

}

func Contains(types []string, subType string) bool {
	for _, a := range types {
		if a == subType {
			return true
		}
	}
	return false
}

func toTime(param interface{}) (time.Time, error) {
	if reflect.TypeOf(param).Kind() == reflect.String {
		resTime, err := time.Parse(time.RFC3339, param.(string))
		if err != nil {
			return time.Time{}, err
		}
		return resTime, nil
	} else if reflect.TypeOf(param).Kind() == reflect.TypeOf(time.Time{}).Kind() {
		return param.(time.Time), nil
	} else {
		return time.Time{}, errors.New("error while converting interface to time.Time")
	}
}
func toDouble(param interface{}) (float64, error) {
	if reflect.TypeOf(param).Kind() == reflect.Float64 {
		return param.(float64), nil
	} else if reflect.TypeOf(param).Kind() == reflect.String {
		val, err := strconv.ParseFloat(param.(string), 64)
		if err != nil {
			return 0, err
		}
		return val, nil
	} else {
		return 0, errors.New("error while converting interface to double")
	}
}

func CopyUpMessage(message *flow.UpMessage) *flow.UpMessage {
	var copiedUpMessage = new(flow.UpMessage)
	copiedUpMessage.Id = message.Id
	copiedUpMessage.Type_ = message.Type_
	copiedUpMessage.Time = message.Time
	copiedUpMessage.Points = clonePoints(message.Points)
	copiedUpMessage.Packet = clonePacket(message.Packet)
	if message.Content != nil {
		copiedUpMessage.Content = cloneInterface(&message.Content)
	} else {
		copiedUpMessage.Content = nil
	}
	copiedUpMessage.SubAccount = cloneSubAccount(message.SubAccount)
	if message.Origin!=nil{
		copiedUpMessage.Origin = &flow.UpOrigin{
			Type_:        message.Origin.Type_,
			Id:           message.Origin.Id,
			ConnectionId: message.Origin.ConnectionId,
			Time:         message.Origin.Time,
		}
	}
	copiedUpMessage.SubType = message.SubType
	copiedUpMessage.Thing = cloneThing(message.Thing)
	if message.Subscriber!=nil{

		copiedUpMessage.Subscriber = &flow.Subscriber{
			Id:      message.Subscriber.Id,
			RealmId: message.Subscriber.RealmId,
		}
	}
	return copiedUpMessage
}
func clonePacket(packet *flow.MessagePacket) *flow.MessagePacket {
	var message interface{}
	if packet!=nil{
		if packet.Message != nil {
			message = cloneInterface(packet.Message)
		} else {
			message = nil
		}
		return &flow.MessagePacket{
			Type_:   packet.Type_,
			Raw:     packet.Raw,
			Message: message,
		}
	}else{
		return nil
	}
}

func cloneSubAccount(subAccount *flow.Account) *flow.Account {
	if subAccount!=nil{
		return &flow.Account{
			Id:      subAccount.Id,
			RealmId: subAccount.RealmId,
		}
	}else{
		return nil
	}
}
func cloneInterface(inter interface{}) interface{} {
	if inter!=nil{
		data, _ := json.Marshal(inter)
		var copiedObj interface{}
		_ = json.Unmarshal(data, &copiedObj)
		return copiedObj
	}else{
		return nil
	}
}
func clonePoints(points map[string]flow.Point) map[string]flow.Point {
	newPoints := make(map[string]flow.Point, len(points))
	for k, v := range points {
		var point = new(flow.Point)
		point.OntologyId = v.OntologyId
		point.UnitId = v.UnitId
		point.Type_ = v.Type_
		point.Records = cloneRecord(v.Records)
		newPoints[k] = *point
	}
	return newPoints
}

func cloneRecord(records []flow.Record) []flow.Record {
	var newRecords []flow.Record
	for _, j := range records {
		var record = new(flow.Record)
		if j.Value != nil {
			record.Value = cloneInterface(j.Value)
		} else {
			record.Value = nil
		}
		record.EventTime = j.EventTime
		if len(j.Coordinates) == 3 {
			record.Coordinates = []float64{j.Coordinates[0], j.Coordinates[1], j.Coordinates[2]}
		} else if len(j.Coordinates) == 3 {
			record.Coordinates = []float64{j.Coordinates[0], j.Coordinates[1]}
		} else {
			record.Coordinates = nil
		}
		newRecords = append(newRecords, *record)
	}
	return newRecords
}
func cloneThing(thing *flow.Thing) *flow.Thing {
	if thing!=nil{
		var newThing = new(flow.Thing)
		newThing.Key = thing.Key
		if thing.Model!=nil{
			newThing.Model = &flow.ModuleSpec{
				ProducerId: thing.Model.ProducerId,
				ModuleId:   thing.Model.ModuleId,
				Version:    thing.Model.Version,
			}
		}
		if thing.Application!=nil{
			newThing.Application = &flow.ModuleSpec{
				ProducerId: thing.Application.ProducerId,
				ModuleId:   thing.Application.ModuleId,
				Version:    thing.Application.Version,
			}
		}
		var tags []string
		for _, j := range thing.Tags {
			tags = append(tags, j)
		}
		newThing.Tags = tags
		return newThing
	}else{
		return nil
	}
}
func cloneCommand(command *flow.Command) *flow.Command{
	if command!=nil{
		var input interface{}
		if command.Input != nil {
			input = cloneInterface(&command.Input)
		} else {
			input = nil
		}
		return &flow.Command{
			Id:    command.Id,
			Input: input,
		}
	}else{
		return nil
	}

}
func CopyDownMessage(message *flow.DownMessage) *flow.DownMessage {
	var copiedDownMessage = new(flow.DownMessage)
	copiedDownMessage.Time = message.Time
	if message.Subscriber!=nil{
		copiedDownMessage.Subscriber = &flow.Subscriber{
			Id:      message.Subscriber.Id,
			RealmId: message.Subscriber.RealmId,
		}
	}
	copiedDownMessage.Id = message.Id
	copiedDownMessage.Thing = cloneThing(message.Thing)
	copiedDownMessage.Type_ = message.Type_
	if message.Origin!=nil{
		copiedDownMessage.Origin = &flow.DownOrigin{
			Type_:        message.Origin.Type_,
			Id:           message.Origin.Id,
			ConnectionId: message.Origin.ConnectionId,
			Time:         message.Origin.Time,
		}
	}
	copiedDownMessage.SubAccount = cloneSubAccount(message.SubAccount)
	if message.Content != nil {
		copiedDownMessage.Content = cloneInterface(message.Content)
	}
	copiedDownMessage.Packet = clonePacket(message.Packet)
	copiedDownMessage.Command = cloneCommand(message.Command)
	return copiedDownMessage
}
