package rules

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
)

func TestDeploymentRuleReplicas_FindNonConformingDeployment(t *testing.T) {
	deployment1 := newDeploymentWithReplicas("default", "one", "uid1", int32(1))
	deployment2 := newDeploymentWithReplicas("default", "two", "uid2", int32(2))
	deployments := []appsv1.Deployment{deployment1, deployment2}

	rule := DeploymentRuleReplicasMinimum{
		MinimumReplicas: 2,
	}

	ruleResult := rule.FindNonConformingDeployment(deployments)
	assert.Equal(t, len(ruleResult.Deployments), 1)
	assert.Equal(t, ruleResult.Deployments[0].Name, "one")
}

func TestDeploymentRuleReplicas_UnmarshalYAML(t *testing.T) {
	yamlString := `
name: minimum replicas 2
minimum_replicas: 2`

	rule := DeploymentRuleReplicasMinimum{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err != nil {
		t.Fail()
	}
	assert.Equal(t, int32(2), rule.MinimumReplicas)
}

func TestDeploymentRuleReplicas_UnmarshalYAML_NameNotFilledIn(t *testing.T) {
	yamlString := `
minimum_replicas: 2`

	rule := DeploymentRuleReplicasMinimum{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err == nil {
		t.Fail()
	}
}

func TestDeploymentRuleReplicas_UnmarshalYAML_MinimumReplicasNotFilledIn(t *testing.T) {
	yamlString := `
name: minimum replicas 2`

	rule := DeploymentRuleReplicasMinimum{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err == nil {
		t.Fail()
	}
}

func newDeploymentWithReplicas(namespace, name string, uid types.UID, replicas int32) appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID:       uid,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
	}
}
