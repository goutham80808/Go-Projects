package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RuleSet struct {
	FileType  string      `yaml:"file_type"`
	HasHeader bool        `yaml:"has_header"`
	Fields    []FieldRule `yaml:"fields"`
}

type FieldRule struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Rules Rules  `yaml:"rules"`
}

type Rules struct {
	Required *bool    `yaml:"required"`
	Min      *int     `yaml:"min"`
	Max      *int     `yaml:"max"`
	Regex    *string  `yaml:"regex"`
	Enum     []string `yaml:"enum"`
}

// Load reads a YAML file from the given path and unmarshals it into a RuleSet struct.
func Load(path string) (*RuleSet, error) {
	// Read the file content
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rs RuleSet

	err = yaml.Unmarshal(data, &rs)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}
