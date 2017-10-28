package config

import (
	"fmt"
	"bytes"
	"net/smtp"
	"github.com/stijndehaes/kube-conformity/rules"
	"html/template"
	"strconv"
	"encoding/base64"
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

func(e EmailConfig) GetMailHeaders() map[string]string {
	headers := make(map[string]string)
	headers["From"] = e.From
	headers["To"] = e.To
	headers["Subject"] = e.Subject + "!"
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

func (e EmailConfig) ConstructEmailBody(results []rules.RuleResult) ([]byte, error) {
	headers := ConstructHeadersString(e.GetMailHeaders())
	body, err := e.RenderTemplate(results)
	if err != nil {
		return []byte{}, err
	}
	return []byte(headers + "\n" + base64.StdEncoding.EncodeToString([]byte(body))), nil
}

func (e EmailConfig) SendMail(results []rules.RuleResult) error {
	msg, err := e.ConstructEmailBody(results)
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth(e.AuthIdentity, e.AuthUsername, e.AuthPassword, e.Host)
	err = smtp.SendMail(e.Host+":"+strconv.Itoa(e.Port), auth, e.From, []string{e.To}, msg)
	return err
}
