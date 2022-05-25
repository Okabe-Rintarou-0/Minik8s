package apiObject

type NodeType string

const (
	ChoiceNode NodeType = "Choice"
	TaskNode   NodeType = "Task"
)

type Choice struct {
	Variable         string   `json:"variable"`
	NumericEquals    *float64 `json:"numericEquals,omitempty"`
	NumericNotEquals *float64 `json:"numericNotEquals,omitempty"`
	Next             *string  `json:"next"`
}

type Choices struct {
	Choices []Choice `json:"choices"`
}

type WorkflowNode struct {
	Type     NodeType `json:"type"`
	*Task    `json:",inline"`
	*Choices `json:",inline"`
}

type Task struct {
	Next *string `json:"next"`
}

type Workflow struct {
	Base    `json:",inline"`
	StartAt string                  `json:"startAt"`
	Nodes   map[string]WorkflowNode `json:"nodes"`
}
