package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestKubeConformityConfig_UnmarshalYAML(t *testing.T) {
	test := ``
	config := Config{}
	err := yaml.Unmarshal([]byte(test), &config)
	if err != nil {
		t.Error(err)
	}
}

func TestKubeConformityConfig_UnmarshalYAML_PodRulesLabelsFilledIn(t *testing.T) {
	test := `
pod_rules_labels_filled_in:
- name: app label filled in
  labels:
  - app`

	config := Config{}

	yaml.Unmarshal([]byte(test), &config)
	assert.Len(t, config.PodRulesLabelsFilledIn, 1)
}

func TestKubeConformityConfig_UnmarshalYAML_PodRulesLimitsFilledIn(t *testing.T) {
	test := `
pod_rules_limits_filled_in:
- name: limits filled in
`

	config := Config{}

	yaml.Unmarshal([]byte(test), &config)
	assert.Len(t, config.PodRulesLimitsFilledIn, 1)
}

func TestKubeConformityConfig_UnmarshalYAML_PodRulesRequestsFilledIn(t *testing.T) {
	test := `
pod_rules_requests_filled_in:
- name: requests filled in`

	config := Config{}

	yaml.Unmarshal([]byte(test), &config)
	assert.Len(t, config.PodRulesRequestsFilledIn, 1)
}

func TestKubeConformityConfig_UnmarshalYAML_DeploymentRuleReplicasMinimum(t *testing.T) {
	test := `
deployment_rules_replicas_minimum:
- name: replicas minimum 1
  minimum_replicas: 2`

	config := Config{}

	yaml.Unmarshal([]byte(test), &config)
	assert.Len(t, config.DeploymentRuleReplicasMinimum, 1)
}

func TestKubeConformityConfig_UnmarshalYAML_Error(t *testing.T) {
	test := `random`

	config := Config{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err == nil {
		assert.Fail(t, "Should have failed")
	}
}