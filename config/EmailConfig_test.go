package config

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
	"github.com/stijndehaes/kube-conformity/rules"
	"os"
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
	assert.Equal(t, "no-reply@kube-conformity.com", config.From)
	assert.False(t, config.Enabled)
}

func TestEmailConfig_UnmarshalYAML_AuthPasswordEnvironment(t *testing.T) {
	test := `
host: 10.10.10.10
to: test@gmail.com`

	config := EmailConfig{}
	os.Setenv("CONFORMITY_EMAIL_AUTH_PASSWORD", "secret")
	err := yaml.Unmarshal([]byte(test), &config)


	if err != nil {
		assert.Fail(t, "UnMarshal should not fail")
	}

	assert.Equal(t, "secret", config.AuthPassword)
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

func TestEmailConfig_UnmarshalYAML_Error(t *testing.T) {
	test := `random`

	config := EmailConfig{}

	err := yaml.Unmarshal([]byte(test), &config)

	if err == nil {
		assert.Fail(t, "Should have failed")
	}
}

func TestEmailConfig_GetMailHeaders(t *testing.T) {
	eConfig := DefaultEmailConfig
	eConfig.To = "test@mail.com"
	headers := eConfig.GetMailHeaders()

	assert.Equal(t, "no-reply@kube-conformity.com", headers["From"])
	assert.Equal(t, "test@mail.com", headers["To"])
	assert.Equal(t, "kube-conformity!", headers["Subject"])
	assert.Equal(t, "1.0", headers["MIME-Version"])
	assert.Equal(t, "text/html; charset=\"utf-8\"", headers["Content-Type"])
	assert.Equal(t, "base64", headers["Content-Transfer-Encoding"])
}

func TestConstructHeadersString(t *testing.T) {
	headers := make(map[string]string)
	headers["Test1"] = "test1"
	headers["Test2"] = "test2"

	string := ConstructHeadersString(headers)

	assert.True(t, "Test1: test1\r\nTest2: test2\r\n" == string || "Test2: test2\r\nTest1: test1\r\n" == string)
}