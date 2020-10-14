package ontology

type UpUpdatePoints struct {
	Points map[string]JmesPathUpdatePoint `json:"points"`
	UpOperation
}

func (updatePoints UpUpdatePoints) ValidUpOperation() string {
	return "updatePoints"
}
