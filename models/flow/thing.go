package flow

type Thing struct {
	Key         string     `json:"key"`
	Model       *ModuleSpec `json:"model,omitempty"`
	Application *ModuleSpec `json:"application,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}
