package operations

import (
	"encoding/json"
	"errors"
	"ontology-mapping-go-lib/models/ontology"
)

type OperationsUpSerDer struct {
	Operations    []ontology.UpOperationInterface `json:"-"`
	RawOperations []json.RawMessage               `json:"operations"`
}

func (opr *OperationsUpSerDer) UnmarshalJSON(b []byte) error {
	type operations OperationsUpSerDer
	err := json.Unmarshal(b, (*operations)(opr))
	if err != nil {
		return err
	}

	for _, raw := range opr.RawOperations {
		var operation ontology.UpOperation
		err = json.Unmarshal(raw, &operation)
		if err != nil {
			return err
		}
		var i ontology.UpOperationInterface
		switch operation.Op {
		case "extractPoints":
			i = &ontology.UpExtractPoints{}
		case "updatePoints":
			i = &ontology.UpUpdatePoints{}
		case "filter":
			i = &ontology.UpFilterOperation{}
		case "filterPoints":
			i = &ontology.UpFilterPointsOperation{}
		default:
			return errors.New("unknown operation type")
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		opr.Operations = append(opr.Operations, i)
	}
	return nil
}
func (opr *OperationsUpSerDer) MarshalJSON() ([]byte, error) {

	type operations OperationsUpSerDer
	if opr.Operations != nil {
		for _, v := range opr.Operations {
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			opr.RawOperations = append(opr.RawOperations, b)
		}
	}
	return json.Marshal((*operations)(opr))
}
