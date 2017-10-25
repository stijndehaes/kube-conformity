package config

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestKubeConformityConfig_UnmarshalYAML_MissingInterval(t *testing.T) {
	test := `
limits_filled_in_rules:
- name: limits filled in`

	config := Config{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err == nil {
		t.Fail()
	}
}
func TestKubeConformityConfig_UnmarshalYAML(t *testing.T) {
	test := `
interval: 1h`

	config := Config{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err != nil {
		t.Error(err)
	}
	dur, _ := time.ParseDuration("1h")
	assert.Equal(t, dur, config.Interval)
}

func TestKubeConformityConfig_UnmarshalYAML_Labels(t *testing.T) {
	test := `
interval: 1h
labels_filled_in_rules:
- name: app label filled in
  labels:
  - app`

	config := Config{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err != nil {
		t.Error(err)
	}
	assert.Len(t, config.LabelsFilledInRules, 1)
}

func TestKubeConformityConfig_UnmarshalYAML_Limits(t *testing.T) {
	test := `
interval: 1h
limits_filled_in_rules:
- name: limits filled in
`

	config := Config{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err != nil {
		t.Error(err)
	}
	assert.Len(t, config.LimitsFilledInRules, 1)
}

func TestKubeConformityConfig_UnmarshalYAML_Requests(t *testing.T) {
	test := `
interval: 1h
requests_filled_in_rules:
- name: requests filled in`

	config := Config{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err != nil {
		t.Error(err)
	}
	assert.Len(t, config.RequestsFilledInRules, 1)
}
