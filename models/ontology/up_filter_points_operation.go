package ontology

type UpFilterPointsOperation struct {
	Points []string `json:"points,omitempty"`
	UpOperation
}

func (filterPoints UpFilterPointsOperation) ValidUpOperation() string {
	return "filterPoints"
}
