package operations

import "ontology-mapping-go-lib/models/ontology"
import "ontology-mapping-go-lib/models/flow"

type OperationHandler interface {
	ApplyUpOperation(message *flow.UpMessage, upOperation *ontology.UpOperationInterface) (*flow.UpMessage, error)
	ApplyDownOperation(message *flow.DownMessage, downOperation *ontology.DownOperationInterface) (*flow.DownMessage, error)
}
