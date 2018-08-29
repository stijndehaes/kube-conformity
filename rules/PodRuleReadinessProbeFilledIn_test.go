package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
)

func TestPodRuleReadinessProbeFilledIn_UnmarshalYAML(t *testing.T) {
	yamlString := `
name: readinessprobe not filled in`

	rule := PodRuleReadinessProbeFilledIn{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err != nil {
		t.Fail()
	}
}

func TestPodRuleReadinessProbeFilledIn_UnmarshalYAML_NameNotFilledIn(t *testing.T) {
	yamlString := `
filter:
  namespaces: test`

	rule := PodRuleReadinessProbeFilledIn{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err == nil {
		t.Fail()
	}
}

func TestPodRuleReadinessProbeFilledIn_FindNonConformingPods(t *testing.T) {
	rule := PodRuleReadinessProbeFilledIn{}
	pods := []v1.Pod{
		newPodWithoutReadinessProbe("default", "name1", "uid1"),
		newPodWithReadinessProbe("default", "name2", "uid2"),
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, "name1", result.Pods[0].Name)
}

func newPodWithoutReadinessProbe(namespace, name string, uid types.UID) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: 	namespace,
			Name:      	name,
			UID:		uid,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:      "container",
				},
			},
		},
	}
}

func newPodWithReadinessProbe(namespace, name string, uid types.UID) v1.Pod {
	pod := newPodWithoutReadinessProbe(namespace, name, uid)

	pod.Spec.Containers[0].ReadinessProbe = &v1.Probe{}
	return pod
}