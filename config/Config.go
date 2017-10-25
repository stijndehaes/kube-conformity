package config

import (
	"github.com/stijndehaes/kube-conformity/rules"
	"time"
	"fmt"
	"net/smtp"
	"log"
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
	AuthUsername string `yaml:"auth_username,omitempty"`
	AuthPassword string `yaml:"auth_password,omitempty"`
	AuthIdentity string `yaml:"auth_identity,omitempty"`
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

func (c *EmailConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain EmailConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	if !c.Enabled {
		return nil
	}
	if c.To == "" {
		return fmt.Errorf("missing to address in email config")
	}
	if c.Host == "" {
		return fmt.Errorf("missing host in email config")
	}
	if c.Port == 0 {
		c.Port = 24
	}
	return nil
}

func (e EmailConfig) sendMail(message []byte) {
	auth := smtp.PlainAuth(
		e.AuthIdentity,
		e.AuthUsername,
		e.AuthPassword,
		e.Host,
	)
	err := smtp.SendMail(
		e.Host+":"+string(e.Port),
		auth,
		e.From,
		[]string{e.To},
		message,
	)
	if err != nil {
		log.Fatal(err)
	}
}
