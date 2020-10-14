package ontology

type UpExtractPoints struct {
	Points map[string]JmesPathPoint `json:"points"`
	UpOperation
}

func (extractPoints UpExtractPoints) ValidUpOperation() string {
	return "extractPoints"
}
