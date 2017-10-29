package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConstructConfig(t *testing.T) {
	configLocation = "config.yaml"
	config, err := ConstructConfig()
	assert.Nil(t, err)
	assert.Len(t, config.LabelsFilledInRules, 1)
	assert.Len(t, config.LimitsFilledInRules, 1)
	assert.Len(t, config.RequestsFilledInRules, 1)
}

func TestConstructConfig_InvalidLocation(t *testing.T) {
	configLocation = "invalid.yaml"
	_, err := ConstructConfig()
	assert.NotNil(t, err)
}