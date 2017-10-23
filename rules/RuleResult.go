package rules

import "k8s.io/client-go/pkg/api/v1"

type RuleResult struct {
	Pods     []v1.Pod
	Reason   string
	RuleName string
}