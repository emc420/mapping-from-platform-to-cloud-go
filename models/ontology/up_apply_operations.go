package ontology

type UpApplyOperations struct {
	Message    interface{}   `json:"message"`
	Operations []UpOperation `json:"operations"`
}
