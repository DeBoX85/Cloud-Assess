package stages

import (
	"fmt"
	"strconv"
	"strings"
)

type OptionSpec struct {
	Type        string
	Default     any
	Description string
}

var optionRegistry = map[string]map[string]OptionSpec{
	Plugin: {
		"target-regions": {
			Type:        "string",
			Default:     "",
			Description: "Comma-separated list of target regions for the scan",
		},
	},
}

// ParseParams preserves the pinned stage.key=value validation contract.
func ParseParams(params []string) (map[string]map[string]any, error) {
	options := make(map[string]map[string]any)
	for _, param := range params {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		stageKey, value, ok := strings.Cut(param, "=")
		if !ok {
			return nil, fmt.Errorf("stage param must be in the form stage.key=value: %s", param)
		}
		stage, key, ok := strings.Cut(stageKey, ".")
		if !ok || stage == "" || key == "" {
			return nil, fmt.Errorf("stage param must be in the form stage.key=value: %s", param)
		}

		stageSpecs, exists := optionRegistry[stage]
		if !exists {
			return nil, fmt.Errorf("unknown stage: %s", stage)
		}
		spec, exists := stageSpecs[key]
		if !exists {
			return nil, fmt.Errorf("unknown option %q for stage %q", key, stage)
		}

		parsed, err := parseOptionValue(value, spec.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid value for %s.%s: %w", stage, key, err)
		}
		if options[stage] == nil {
			options[stage] = make(map[string]any)
		}
		options[stage][key] = parsed
	}
	return options, nil
}

func (c *Config) ApplyParams(params []string) error {
	parsed, err := ParseParams(params)
	if err != nil {
		return err
	}
	for stage, values := range parsed {
		if err := c.SetOptions(stage, values); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) SetOptions(stage string, values map[string]any) error {
	stage = strings.ToLower(strings.TrimSpace(stage))
	if _, ok := defaults[stage]; !ok {
		return fmt.Errorf("unknown stage name: %s", stage)
	}
	if c.options == nil {
		c.options = make(map[string]map[string]any)
	}
	if c.options[stage] == nil {
		c.options[stage] = make(map[string]any)
	}
	for key, value := range values {
		c.options[stage][key] = value
	}
	return nil
}

func (c *Config) Options(stage string) map[string]any {
	stage = strings.ToLower(strings.TrimSpace(stage))
	values := c.options[stage]
	if values == nil {
		return nil
	}
	copyValues := make(map[string]any, len(values))
	for key, value := range values {
		copyValues[key] = value
	}
	return copyValues
}

func parseOptionValue(value, typeName string) (any, error) {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "\"'")

	switch typeName {
	case "bool":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("expected bool, got %q", value)
		}
		return parsed, nil
	case "int":
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("expected int, got %q", value)
		}
		return parsed, nil
	case "float64":
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("expected float64, got %q", value)
		}
		return parsed, nil
	case "string":
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeName)
	}
}
