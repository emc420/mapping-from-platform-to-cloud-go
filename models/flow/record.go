package flow

import (
	"time"
)

type Record struct {
	Value       interface{} `json:"value,omitempty"`
	Coordinates []float64   `json:"coordinates,omitempty"`
	EventTime   time.Time   `json:"eventTime"`
}
