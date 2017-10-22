package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/stijndehaes/kube-conformity/kubeconformity"
)

func TestConstructRulesRequests(t *testing.T) {
	requests = true
	limits = false
	labels = []string{}

	rules := ConstructRules()
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, reflect.TypeOf(kubeconformity.RequestsFilledInRule{}), reflect.TypeOf(rules[0]))
}

func TestConstructRulesLimits(t *testing.T) {
	requests = false
	limits = true
	labels = []string{}

	rules := ConstructRules()
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, reflect.TypeOf(kubeconformity.LimitsFilledInRule{}), reflect.TypeOf(rules[0]))
}

func TestConstructRulesLables(t *testing.T) {
	requests = false
	limits = false
	labels = []string{"app"}

	rules := ConstructRules()
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, reflect.TypeOf(kubeconformity.LabelsFilledInRule{}), reflect.TypeOf(rules[0]))
}