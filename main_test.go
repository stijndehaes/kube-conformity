package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	log "github.com/sirupsen/logrus"
)

func TestConstructConfig(t *testing.T) {
	configLocation = "config.yaml"
	config, err := ConstructConfig()
	assert.Nil(t, err)
	assert.Len(t, config.LabelsFilledInRules, 1)
	assert.Len(t, config.LimitsFilledInRules, 1)
	assert.Len(t, config.RequestsFilledInRules, 1)
}

func TestConfigureLogging(t *testing.T) {
	debug = false
	jsonLogging = false

	ConfigureLogging()
	assert.Equal(t, &log.TextFormatter{},log.StandardLogger().Formatter)
	assert.Equal(t, log.InfoLevel,log.StandardLogger().Level)

	debug = true
	jsonLogging = true

	ConfigureLogging()
	assert.Equal(t, &log.JSONFormatter{},log.StandardLogger().Formatter)
	assert.Equal(t, log.DebugLevel,log.StandardLogger().Level)
}

func TestConstructConfig_InvalidLocation(t *testing.T) {
	configLocation = "invalid.yaml"
	_, err := ConstructConfig()
	assert.NotNil(t, err)
}