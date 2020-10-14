package ontology

type JmesPathUpdatePoint struct {
	OntologyId  string            `json:"ontologyId,omitempty"`
	Value       string            `json:"value,omitempty"`
	Coordinates []string          `json:"coordinates,omitempty"`
	EventTime   string            `json:"eventTime,omitempty"`
	Type_       JmesPathPointType `json:"type,omitempty"`
	UnitId      string            `json:"unitId,omitempty"`
}
