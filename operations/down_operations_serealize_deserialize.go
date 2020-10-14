package operations

import (
	"encoding/json"
	"errors"
	"ontology-mapping-go-lib/models/ontology"
)

type OperationsDownSerDer struct {
	Operations    []ontology.DownOperationInterface `json:"-"`
	RawOperations []json.RawMessage                 `json:"operations"`
}

func (opr *OperationsDownSerDer) UnmarshalJSON(b []byte) error {
	type operations OperationsDownSerDer
	err := json.Unmarshal(b, (*operations)(opr))
	if err != nil {
		return err
	}

	for _, raw := range opr.RawOperations {
		var operation ontology.DownOperation
		err = json.Unmarshal(raw, &operation)
		if err != nil {
			return err
		}
		var i ontology.DownOperationInterface
		switch operation.Op {
		case "extractDriverMessage":
			i = &ontology.DownExtractDriverMessage{}
		case "updateCommand":
			i = &ontology.DownUpdateCommand{}
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
func (opr *OperationsDownSerDer) MarshalJSON() ([]byte, error) {

	type operations OperationsDownSerDer
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
