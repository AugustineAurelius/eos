package projectv2

import (
	"gopkg.in/yaml.v3"
)

type Spec struct {
	Domain Domain `yaml:"domain"`
}

type Domain struct {
	Entities map[string]Entity `yaml:"entities"`
}

type Entity struct {
	Fields map[string]FieldType `yaml:"fields"`
}

type FieldType struct {
	IsArray bool
	Type    string
}

func (ft *FieldType) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.SequenceNode {
		ft.IsArray = true
		if len(value.Content) > 0 {
			ft.Type = value.Content[0].Value
		}
		return nil
	}

	ft.IsArray = false
	ft.Type = value.Value
	return nil
}
