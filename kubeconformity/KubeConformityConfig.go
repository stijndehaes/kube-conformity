package kubeconformity

type KubeConformityConfig struct {
	LabelsFilledInRules   []LabelsFilledInRule   `yaml:"labelsFilledInRules"`
	LimitsFilledInRules   []LimitsFilledInRule   `yaml:"limitsFilledInRules"`
	RequestsFilledInRules []RequestsFilledInRule `yaml:"requestsFilledInRules"`
}