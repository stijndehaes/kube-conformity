package config

import (
	"fmt"
	"bytes"
	"net/smtp"
	"github.com/stijndehaes/kube-conformity/rules"
	"html/template"
	"strconv"
	"encoding/base64"
	"os"
)

var (
	DefaultEmailConfig = EmailConfig{
		Enabled:  false,
		Subject:  "kube-conformity",
		Template: "mailtemplate.html",
		Port:     24,
		From:     "no-reply@kube-conformity.com",
	}
)

type EmailConfig struct {
	Enabled      bool   `yaml:"enabled,omitempty"`
	To           string `yaml:"to,omitempty"`
	From         string `yaml:"from,omitempty"`
	Host         string `yaml:"host,omitempty"`
	Port         int    `yaml:"port"`
	Subject      string `yaml:"subject"`
	AuthUsername string `yaml:"auth_username,omitempty"`
	AuthPassword string `yaml:"auth_password,omitempty"`
	AuthIdentity string `yaml:"auth_identity,omitempty"`
	Template     string `yaml:"template,omitempty"`
}

func (emailConfig *EmailConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*emailConfig = DefaultEmailConfig
	type plain EmailConfig
	if err := unmarshal((*plain)(emailConfig)); err != nil {
		return err
	}
	if emailConfig.To == "" {
		return fmt.Errorf("missing to address in email config")
	}
	if emailConfig.Host == "" {
		return fmt.Errorf("missing host in email config")
	}
	authPassword := os.Getenv("CONFORMITY_EMAIL_AUTH_PASSWORD")
	if authPassword != "" {
		emailConfig.AuthPassword = authPassword
	}
	return nil
}

func (emailConfig EmailConfig) RenderTemplate(podRuleResults []rules.PodRuleResult, deploymentResults []rules.DeploymentRuleResult) (string, error) {
	templateData := struct {
		PodRuleResults []rules.PodRuleResult
		DeploymentRuleResults []rules.DeploymentRuleResult
	}{
		PodRuleResults: podRuleResults,
		DeploymentRuleResults: deploymentResults,
	}
	t, err := template.ParseFiles(emailConfig.Template)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func(emailConfig EmailConfig) GetMailHeaders() map[string]string {
	headers := make(map[string]string)
	headers["From"] = emailConfig.From
	headers["To"] = emailConfig.To
	headers["Subject"] = emailConfig.Subject + "!"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""
	headers["Content-Transfer-Encoding"] = "base64"
	return headers
}

func ConstructHeadersString(headers map[string]string) string {
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return message
}

func (emailConfig EmailConfig) ConstructEmailBody(podRuleResults []rules.PodRuleResult, deploymentResults []rules.DeploymentRuleResult) ([]byte, error) {
	headers := ConstructHeadersString(emailConfig.GetMailHeaders())
	body, err := emailConfig.RenderTemplate(podRuleResults, deploymentResults)
	if err != nil {
		return []byte{}, err
	}
	return []byte(headers + "\n" + base64.StdEncoding.EncodeToString([]byte(body))), nil
}

func (emailConfig EmailConfig) SendMail(podRuleResults []rules.PodRuleResult, deploymentRuleResults []rules.DeploymentRuleResult) error {
	msg, err := emailConfig.ConstructEmailBody(podRuleResults, deploymentRuleResults)
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth(emailConfig.AuthIdentity, emailConfig.AuthUsername, emailConfig.AuthPassword, emailConfig.Host)
	err = smtp.SendMail(emailConfig.Host+":"+strconv.Itoa(emailConfig.Port), auth, emailConfig.From, []string{emailConfig.To}, msg)
	return err
}
