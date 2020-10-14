package ontology

type DownApplyOperations struct {
	Message        interface{}     `json:"message"`
	OperationsDown []DownOperation `json:"operationsDown"`
}
