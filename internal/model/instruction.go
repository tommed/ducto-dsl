package model

type Instruction struct {
	Op string `json:"op"`

	// Common fields
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`

	// Copy
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`

	// Delete
	Regex bool `json:"regex,omitempty"`

	// Conditional
	Condition map[string]interface{} `json:"condition,omitempty"`
	Not       bool                   `json:"not,omitempty"`
	Then      []Instruction          `json:"then,omitempty"`

	// Others will be added as needed
}

type Program struct {
	OnError      string        `json:"on_error,omitempty"`
	Instructions []Instruction `json:"instructions"`
}
