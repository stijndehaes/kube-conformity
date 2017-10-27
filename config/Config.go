package config

import (
	"github.com/stijndehaes/kube-conformity/rules"
	"time"
	"fmt"
)

type Config struct {
	Interval              time.Duration                `yaml:"interval"`
	LabelsFilledInRules   []rules.LabelsFilledInRule   `yaml:"labels_filled_in_rules"`
	LimitsFilledInRules   []rules.LimitsFilledInRule   `yaml:"limits_filled_in_rules"`
	RequestsFilledInRules []rules.RequestsFilledInRule `yaml:"requests_filled_in_rules"`
	EmailConfig           EmailConfig                  `yaml:"email_config"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.Interval == 0 {
		return fmt.Errorf("Missing interval in config")
	}
	return nil
}
