package transform

import (
	"encoding/json"
	"fmt"
	"github.com/tommed/ducto-dsl/model"
	"os"
)

// LoadProgram loads and parses a DSL program from disk.
func LoadProgram(path string) (*model.Program, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read program: %w", err)
	}
	var prog model.Program
	if err := json.Unmarshal(data, &prog); err != nil {
		return nil, fmt.Errorf("failed to parse program: %w", err)
	}
	return &prog, nil
}
