package rules

import (
	"fmt"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

type DeploymentRuleReplicasMinimum struct {
	Name            string            `yaml:"name"`
	MinimumReplicas int32             `yaml:"minimum_replicas"`
}

func (d DeploymentRuleReplicasMinimum) FindNonConformingDeployment(deployments []v1beta1.Deployment) DeploymentRuleResult {
	filteredList := []v1beta1.Deployment{}
	for _, deployment := range deployments {
		if *deployment.Spec.Replicas < d.MinimumReplicas {
			filteredList = append(filteredList, deployment)
		}
	}

	return DeploymentRuleResult{
		Deployments:     filteredList,
		Reason:   fmt.Sprintf("Replicas below the minimum: %v", d.MinimumReplicas),
		RuleName: d.Name,
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
