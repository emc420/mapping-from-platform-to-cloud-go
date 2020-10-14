package ontology

type DownUpdateCommand struct {
	Commands map[string]UpdateCommand `json:"commands"`
	DownOperation
}

func (updateCommand DownUpdateCommand) ValidDownOperation() string {
	return "updateCommand"
}
