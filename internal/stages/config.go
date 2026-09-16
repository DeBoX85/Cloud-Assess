package stages

import (
	"fmt"
	"sort"
	"strings"
)

const (
	Graph                   = "graph"
	Advisor                 = "advisor"
	Defender                = "defender"
	DefenderRecommendations = "defender-recommendations"
	Arc                     = "arc"
	Policy                  = "policy"
	Cost                    = "cost"
	Diagnostics             = "diagnostics"
	Plugin                  = "plugin"
)

var defaults = map[string]bool{
	Graph:                   true,
	Diagnostics:             true,
	Advisor:                 true,
	Defender:                true,
	DefenderRecommendations: false,
	Arc:                     false,
	Policy:                  false,
	Cost:                    false,
	Plugin:                  false,
}

type Config struct {
	enabled map[string]bool
	options map[string]map[string]any
}

func NewDefault() *Config {
	c := &Config{
		enabled: map[string]bool{},
		options: map[string]map[string]any{},
	}
	for name, enabled := range defaults {
		c.enabled[name] = enabled
	}
	return c
}

func (c *Config) IsEnabled(name string) bool {
	return c.enabled[strings.ToLower(strings.TrimSpace(name))]
}

func (c *Config) Set(name string, enabled bool) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if _, ok := defaults[name]; !ok {
		return fmt.Errorf("unknown stage name: %s", name)
	}
	c.enabled[name] = enabled
	return nil
}

func (c *Config) Apply(values []string) error {
	for _, value := range values {
		for _, token := range strings.Split(value, ",") {
			token = strings.ToLower(strings.TrimSpace(token))
			if token == "" {
				continue
			}
			enabled := true
			if strings.HasPrefix(token, "-") {
				enabled = false
				token = strings.TrimPrefix(token, "-")
			}
			if err := c.Set(token, enabled); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) Validate() error {
	if !c.IsEnabled(Graph) {
		return fmt.Errorf("graph stage is mandatory for regular scans and cannot be disabled")
	}
	return nil
}

func (c *Config) EnabledStages() []string {
	var out []string
	for name, enabled := range c.enabled {
		if enabled {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}
