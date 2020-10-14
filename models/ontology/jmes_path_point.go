package ontology

type JmesPathPoint struct {
	OntologyId  string            `json:"ontologyId,omitempty"`
	Value       string            `json:"value,omitempty"`
	Coordinates []string          `json:"coordinates,omitempty"`
	EventTime   string            `json:"eventTime"`
	Type_       JmesPathPointType `json:"type,omitempty"`
	UnitId      string            `json:"unitId,omitempty"`
}
