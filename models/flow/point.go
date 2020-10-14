package flow

type Point struct {
	OntologyId string    `json:"ontologyId,omitempty"`
	Type_      PointType `json:"type,omitempty"`
	UnitId     string    `json:"unitId,omitempty"`
	Records    []Record  `json:"records"`
}
