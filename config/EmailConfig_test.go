package config

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
	"github.com/stijndehaes/kube-conformity/rules"
)

func TestEmailConfig_UnmarshalYAML_FailMissingHost(t *testing.T) {
	test := `
to: test@gmail.com`

	config := EmailConfig{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err == nil {
		assert.Fail(t, "Expected error for missing host")
	}
}

func TestEmailConfig_UnmarshalYAML_FailMissingTo(t *testing.T) {
	test := `
host: 10.10.10.10`

	config := EmailConfig{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err == nil {
		assert.Fail(t, "Expected error for missing to")
	}
}

func TestEmailConfig_UnmarshalYAML_DefaultValues(t *testing.T) {
	test := `
host: 10.10.10.10
to: test@gmail.com`

	config := EmailConfig{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err != nil {
		assert.Fail(t, "UnMarshal should not fail")
	}

	assert.Equal(t, "mailtemplate.html", config.Template)
	assert.Equal(t, 24, config.Port)
	assert.Equal(t, "kube-conformity", config.Subject)
	assert.False(t, config.Enabled)
}

func TestEmailConfig_RenderTemplate(t *testing.T) {
	eConfig := DefaultEmailConfig
	eConfig.Enabled = true
	eConfig.Template = "../mailtemplate.html"

	template, err := eConfig.RenderTemplate([]rules.RuleResult{
		{
			Reason:   "A reason",
			RuleName: "A rule name",
		},
	})

	if err != nil {
		assert.Fail(t, "Template should render correctly")
	}
	assert.NotEqual(t, "", template)
}

func TestEmailConfig_ConstructEmailBody(t *testing.T) {
	eConfig := DefaultEmailConfig
	eConfig.Enabled = true
	eConfig.Template = "../mailtemplate.html"

	body, err := eConfig.ConstructEmailBody([]rules.RuleResult{
		{
			Reason:   "A reason",
			RuleName: "A rule name",
		},
	})

	if err != nil {
		assert.Fail(t, "Body should render correctly")
	}
	assert.NotEqual(t, "", body)
}

func TestEmailConfig_ConstructEmailBody_TemplateNoExist(t *testing.T) {
	eConfig := DefaultEmailConfig
	eConfig.Enabled = true
	eConfig.Template = "test.html"
	body, err := eConfig.ConstructEmailBody([]rules.RuleResult{{}})
	assert.NotEqual(t, nil, err, "Should fail because template does not exist")
	assert.Equal(t, []byte{}, body)
}
