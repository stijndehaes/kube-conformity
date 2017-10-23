package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/stijndehaes/kube-conformity/kubeconformity"
)

func TestConstructRulesRequests(t *testing.T) {
	config := ConstructConfig()
	assert.Len(t, config.LabelsFilledInRules, 1)
	assert.Len(t, config.LimitsFilledInRules, 1)
	assert.Len(t, config.RequestsFilledInRules, 1)
}