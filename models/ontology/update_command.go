package ontology

type UpdateCommand struct {
	Id    string      `json:"id,omitempty"`
	Input interface{} `json:"input,omitempty"`
}
