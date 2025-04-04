package transform

type Program struct {
	Version      int           `json:"version" mapstructure:"version"`
	OnError      string        `json:"on_error,omitempty" mapstructure:"on_error"`
	Instructions []Instruction `json:"instructions" mapstructure:"instructions"`
}

type Instruction struct {
	Op string `json:"op" mapstructure:"op"`

	// Common fields
	Key   string      `json:"key,omitempty" mapstructure:"key"`
	Value interface{} `json:"value,omitempty" mapstructure:"value"`

	// Copy
	From string `json:"from,omitempty" mapstructure:"from"`
	To   string `json:"to,omitempty" mapstructure:"to"`

	// Delete
	Regex bool `json:"regex,omitempty" mapstructure:"regex"`

	// Conditional
	Condition map[string]interface{} `json:"condition,omitempty" mapstructure:"condition"`
	Not       bool                   `json:"not,omitempty" mapstructure:"not"`
	Then      []Instruction          `json:"then,omitempty" mapstructure:"then"`
	IfNotSet  bool                   `json:"if_not_set,omitempty" mapstructure:"if_not_set"`
}
