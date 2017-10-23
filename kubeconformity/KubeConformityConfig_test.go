package kubeconformity

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
)

func TestKubeConformityConfig_UnmarshalYAML_labels(t *testing.T) {
	test := `
labelsFilledInRules:
- name: app label filled in
  labels:
  - app`

	config := KubeConformityConfig{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err != nil {
		t.Error(err)
	}
	assert.Len(t, config.LabelsFilledInRules, 1)
}

func TestKubeConformityConfig_UnmarshalYAML_limits(t *testing.T) {
	test := `
limitsFilledInRules:
- name: limits filled in
`

	config := KubeConformityConfig{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err != nil {
		t.Error(err)
	}
	assert.Len(t, config.LimitsFilledInRules, 1)
}

func TestKubeConformityConfig_UnmarshalYAML_requests(t *testing.T) {
	test := `
requestsFilledInRules:
- name: requests filled in`

	config := KubeConformityConfig{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err != nil {
		t.Error(err)
	}
	assert.Len(t, config.RequestsFilledInRules, 1)
}
