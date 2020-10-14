package operations

import "ontology-mapping-go-lib/models/flow"

type OperationService struct {
	Factory OperationFactory
}

func (operationService *OperationService) ApplyUpOperations(message *flow.UpMessage, operations *OperationsUpSerDer) (*flow.UpMessage, error) {
	if message == nil {
		return nil, nil
	}
	var err error
	var handler OperationHandler
	var retMessage = new(flow.UpMessage)
	retMessage = message
	for _, operation := range operations.Operations {
		if retMessage == nil {
			return nil, nil
		}
		handler, err = operationService.Factory.BuildUp(operation)
		if err != nil {
			return nil, err
		}
		retMessage, err = handler.ApplyUpOperation(retMessage, &operation)
		if err != nil {
			return nil, err
		}
	}
	return retMessage, nil
}

func (operationService *OperationService) ApplyDownOperations(message *flow.DownMessage, operations *OperationsDownSerDer) (*flow.DownMessage, error) {
	if message == nil {
		return nil, nil
	}
	var err error
	var handler OperationHandler
	var retMessage = new(flow.DownMessage)
	retMessage = message
	for _, operation := range operations.Operations {
		if retMessage == nil {
			return nil, nil
		}
		handler, err = operationService.Factory.BuildDown(operation)
		if err != nil {
			return nil, err
		}
		retMessage, err = handler.ApplyDownOperation(retMessage, &operation)
		if err != nil {
			return nil, err
		}
	}
	return retMessage, nil
}
