package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConstructRulesRequests(t *testing.T) {
	configLocation = "config.yaml"
	config := ConstructConfig()
	assert.Len(t, config.LabelsFilledInRules, 1)
	assert.Len(t, config.LimitsFilledInRules, 1)
	assert.Len(t, config.RequestsFilledInRules, 1)
}