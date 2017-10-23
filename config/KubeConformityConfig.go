package config

import "github.com/stijndehaes/kube-conformity/rules"

type KubeConformityConfig struct {
	LabelsFilledInRules   []rules.LabelsFilledInRule`yaml:"labelsFilledInRules"`
	LimitsFilledInRules   []rules.LimitsFilledInRule   `yaml:"limitsFilledInRules"`
	RequestsFilledInRules []rules.RequestsFilledInRule `yaml:"requestsFilledInRules"`
}