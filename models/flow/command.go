package flow

type Command struct {
	Id    string      `json:"id"`
	Input interface{} `json:"input,omitempty"`
}
