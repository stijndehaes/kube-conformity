package rules

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
)

func TestStatefulSetRuleReplicas_FindNonConformingStatefulSet(t *testing.T) {
	statefulSet1 := newStatefulSetWithReplicas("default", "one", "uid1", int32(1))
	statefulSet2 := newStatefulSetWithReplicas("default", "two", "uid2", int32(2))
	statefulSets := []appsv1.StatefulSet{statefulSet1, statefulSet2}

	rule := StatefulSetRuleReplicasMinimum{
		MinimumReplicas: 2,
	}

	ruleResult := rule.FindNonConformingStatefulSet(statefulSets)
	assert.Equal(t, len(ruleResult.StatefulSets), 1)
	assert.Equal(t, ruleResult.StatefulSets[0].Name, "one")
}

func TestStatefulSetRuleReplicas_UnmarshalYAML(t *testing.T) {
	yamlString := `
name: minimum replicas 2
minimum_replicas: 2`

	rule := StatefulSetRuleReplicasMinimum{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err != nil {
		t.Fail()
	}
	assert.Equal(t, int32(2), rule.MinimumReplicas)
}

func TestStatefulSetRuleReplicas_UnmarshalYAML_NameNotFilledIn(t *testing.T) {
	yamlString := `
minimum_replicas: 2`

	rule := StatefulSetRuleReplicasMinimum{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err == nil {
		t.Fail()
	}
}

func TestStatefulSetRuleReplicas_UnmarshalYAML_MinimumReplicasNotFilledIn(t *testing.T) {
	yamlString := `
name: minimum replicas 2`

	rule := StatefulSetRuleReplicasMinimum{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err == nil {
		t.Fail()
	}
}

func newStatefulSetWithReplicas(namespace, name string, uid types.UID, replicas int32) appsv1.StatefulSet {
	return appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID:       uid,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
		},
	}
}
