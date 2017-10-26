package config

import (
	"github.com/stijndehaes/kube-conformity/rules"
	"time"
	"fmt"
	"net/smtp"
	"strconv"
	"bytes"
	"html/template"
)

type Config struct {
	Interval              time.Duration                `yaml:"interval"`
	LabelsFilledInRules   []rules.LabelsFilledInRule   `yaml:"labels_filled_in_rules"`
	LimitsFilledInRules   []rules.LimitsFilledInRule   `yaml:"limits_filled_in_rules"`
	RequestsFilledInRules []rules.RequestsFilledInRule `yaml:"requests_filled_in_rules"`
	EmailConfig           EmailConfig                  `yaml:"email_config"`
}

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

var (
	DefaultEmailConfig = EmailConfig{
		Enabled: false,
		Subject: "kube-conformity",
		Template: "mailtemplate.html",
		Port: 24,
	}
)

func (c *EmailConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultEmailConfig
	type plain EmailConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if c.To == "" {
		return fmt.Errorf("missing to address in email config")
	}
	if c.Host == "" {
		return fmt.Errorf("missing host in email config")
	}
	return nil
}

func (e EmailConfig) RenderTemplate(results []rules.RuleResult) (string, error) {
	templateData := struct {
		RuleResults []rules.RuleResult
	}{
		RuleResults: results,
	}
	t, err := template.ParseFiles(e.Template)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (e EmailConfig) SendMail(results []rules.RuleResult) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + e.Subject + "!\n"
	body, err := e.RenderTemplate(results)
	if err != nil {
		return err
	}
	msg := []byte(subject + mime + "\n" + body)
	auth := smtp.PlainAuth(
		e.AuthIdentity,
		e.AuthUsername,
		e.AuthPassword,
		e.Host,
	)
	err = smtp.SendMail(
		e.Host+":"+strconv.Itoa(e.Port),
		auth,
		e.From,
		[]string{e.To},
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}
