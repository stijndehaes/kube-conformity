package rules

import (
	"fmt"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	"github.com/stijndehaes/kube-conformity/filters"
)

type DeploymentRuleReplicasMinimum struct {
	Name            string                   `yaml:"name"`
	MinimumReplicas int32                    `yaml:"minimum_replicas"`
	Filter          filters.DeploymentFilter `yaml:"filter"`
}

func (d DeploymentRuleReplicasMinimum) FindNonConformingDeployment(deployments []v1beta1.Deployment) DeploymentRuleResult {
	filteredDeployments := d.Filter.FilterDeployments(deployments)
	var nonConformingDeployments  []v1beta1.Deployment
	for _, deployment := range filteredDeployments {
		if *deployment.Spec.Replicas < d.MinimumReplicas {
			nonConformingDeployments = append(nonConformingDeployments, deployment)
		}
	}

	return DeploymentRuleResult{
		Deployments: nonConformingDeployments,
		Reason:      fmt.Sprintf("Replicas below the minimum: %v", d.MinimumReplicas),
		RuleName:    d.Name,
	}
}

func (r *DeploymentRuleReplicasMinimum) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain DeploymentRuleReplicasMinimum
	if err := unmarshal((*plain)(r)); err != nil {
		return err
	}
	if r.MinimumReplicas == 0 {
		return fmt.Errorf("Missing minimum replicas")
	}
	if r.Name == "" {
		return fmt.Errorf("Missing name for DeploymentRuleReplicasMinimum")
	}
	return nil
}
