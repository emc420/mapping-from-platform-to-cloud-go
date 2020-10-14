package operations

import (
	"errors"
	"ontology-mapping-go-lib/models/ontology"
)

type OperationFactory struct {
}

func (operationFactory *OperationFactory) BuildUp(operation ontology.UpOperationInterface) (OperationHandler, error) {
	switch operation.(type) {
	case ontology.UpExtractPoints:
		return &UpExtractPointsOperation{}, nil
	case ontology.UpUpdatePoints:
		return &UpUpdatePointsOperation{}, nil
	case ontology.UpFilterOperation:
		return &FilterOperation{}, nil
	case ontology.UpFilterPointsOperation:
		return &FilterPointsOperation{}, nil
	default:
		return nil, errors.New("unknown up Operation")
	}
}

func (operationFactory *OperationFactory) BuildDown(operation ontology.DownOperationInterface) (OperationHandler, error) {
	switch operation.(type) {
	case ontology.DownExtractDriverMessage:
		return &DownExtractDriverOperation{}, nil
	case ontology.DownUpdateCommand:
		return &DownUpdateCommandOperation{}, nil
	default:
		return nil, errors.New("unknown up Operation")
	}

}
