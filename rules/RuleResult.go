package rules

import (
	apiv1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
)

type PodRuleResult struct {
	Pods     []apiv1.Pod
	Reason   string
	RuleName string
}

type DeploymentRuleResult struct {
	Deployments []appsv1.Deployment
	Reason      string
	RuleName    string
}
