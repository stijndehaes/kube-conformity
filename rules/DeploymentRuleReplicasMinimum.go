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

func (deploymentRuleReplicasMinimum DeploymentRuleReplicasMinimum) FindNonConformingDeployment(deployments []v1beta1.Deployment) DeploymentRuleResult {
	filteredDeployments := deploymentRuleReplicasMinimum.Filter.FilterDeployments(deployments)
	var nonConformingDeployments  []v1beta1.Deployment
	for _, deployment := range filteredDeployments {
		if *deployment.Spec.Replicas < deploymentRuleReplicasMinimum.MinimumReplicas {
			nonConformingDeployments = append(nonConformingDeployments, deployment)
		}
	}

	return DeploymentRuleResult{
		Deployments: nonConformingDeployments,
		Reason:      fmt.Sprintf("Replicas below the minimum: %v", deploymentRuleReplicasMinimum.MinimumReplicas),
		RuleName:    deploymentRuleReplicasMinimum.Name,
	}
}

func (deploymentRuleReplicasMinimum *DeploymentRuleReplicasMinimum) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain DeploymentRuleReplicasMinimum
	if err := unmarshal((*plain)(deploymentRuleReplicasMinimum)); err != nil {
		return err
	}
	if deploymentRuleReplicasMinimum.MinimumReplicas == 0 {
		return fmt.Errorf("Missing minimum replicas")
	}
	if deploymentRuleReplicasMinimum.Name == "" {
		deploymentRuleReplicasMinimum.Name = "Minimum replicas too low"
	}
	return nil
}
