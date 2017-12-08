package rules

import (
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/magiconair/properties/assert"
)


func TestDeploymentRuleReplicas_FindNonConformingDeployment(t *testing.T) {
	deployment1 := newDeploymentWithReplicas("default", "one", int32(1))
	deployment2 := newDeploymentWithReplicas("default", "two", int32(2))
	deployments := []v1beta1.Deployment{deployment1, deployment2}

	rule := DeploymentRuleReplicas{
		MinimumReplicas: 2,
	}

	ruleResult := rule.FindNonConformingDeployment(deployments)
	assert.Equal(t, len(ruleResult.Deployments), 1)
	assert.Equal(t, ruleResult.Deployments[0].Name, "one")
}

func TestDeploymentRuleReplicas_UnmarshalYAML(t *testing.T) {
	string := `
name: minimum replicas 2
minimum_replicas: 2`

	rule := DeploymentRuleReplicas{}

	err := yaml.Unmarshal([]byte(string), &rule)

	if err != nil {
		t.Fail()
	}
	assert.Equal(t, int32(2), rule.MinimumReplicas)
}

func TestDeploymentRuleReplicas_UnmarshalYAML_NameNotFilledIn(t *testing.T) {
	string := `
minimum_replicas: 2`

	rule := DeploymentRuleReplicas{}

	err := yaml.Unmarshal([]byte(string), &rule)

	if err == nil {
		t.Fail()
	}
}

func TestDeploymentRuleReplicas_UnmarshalYAML_MinimumReplicasNotFilledIn(t *testing.T) {
	string := `
name: minimum replicas 2`

	rule := DeploymentRuleReplicas{}

	err := yaml.Unmarshal([]byte(string), &rule)

	if err == nil {
		t.Fail()
	}
}


func newDeploymentWithReplicas(namespace, name string, replicas int32) v1beta1.Deployment {
	return v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
		},
	}
}
