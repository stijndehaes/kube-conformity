package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

type PodRuleResult struct {
	Pods     []v1.Pod
	Reason   string
	RuleName string
}

type DeploymentRuleResult struct {
	Deployments []v1beta1.Deployment
	Reason      string
	RuleName    string
}
