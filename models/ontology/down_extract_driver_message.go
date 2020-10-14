package ontology

type DownExtractDriverMessage struct {
	Commands map[string]interface{} `json:"commands"`
	DownOperation
}

func (extractDriver DownExtractDriverMessage) ValidDownOperation() string {
	return "extractDriverMessage"
}
