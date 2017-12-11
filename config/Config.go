package config

import (
	"github.com/stijndehaes/kube-conformity/rules"
	"time"
	"fmt"
)

type Config struct {
	Interval                      time.Duration                         `yaml:"interval"`
	PodRulesLabelsFilledIn        []rules.PodRuleLabelsFilledIn         `yaml:"pod_rules_labels_filled_in"`
	PodRulesLimitsFilledIn        []rules.PodRuleLimitsFilledIn         `yaml:"pod_rules_limits_filled_in"`
	PodRulesRequestsFilledIn      []rules.PodRuleRequestsFilledIn       `yaml:"pod_rules_requests_filled_in"`
	DeploymentRuleReplicasMinimum []rules.DeploymentRuleReplicasMinimum `yaml:"deployment_rules_replicas_minimum"`
	EmailConfig                   EmailConfig                           `yaml:"email_config"`
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
