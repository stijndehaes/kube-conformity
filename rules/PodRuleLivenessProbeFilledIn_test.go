package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
)

func TestPodRuleLivenessProbeFilledIn_UnmarshalYAML(t *testing.T) {
	yamlString := `
name: liveness probe not filled in`

	rule := PodRuleLivenessProbeFilledIn{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err != nil {
		t.Fail()
	}
}

func TestPodRuleLivenessProbeFilledIn_UnmarshalYAML_DefaultName(t *testing.T) {
	yamlString := `
filter:
  namespaces: test`

	rule := PodRuleLivenessProbeFilledIn{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err != nil {
		t.Fail()
	}
	assert.Equal(t, "Liveness probe not filled in", rule.Name)
}

func TestPodRuleLivenessProbeFilledIn_FindNonConformingPods(t *testing.T) {
	rule := PodRuleLivenessProbeFilledIn{}
	pods := []v1.Pod{
		newPodWithoutLivenessProbe("default", "name1", "uid1"),
		newPodWithLivenessProbe("default", "name2", "uid2"),
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, "name1", result.Pods[0].Name)
}

func newPodWithoutLivenessProbe(namespace, name string, uid types.UID) v1.Pod {
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

func newPodWithLivenessProbe(namespace, name string, uid types.UID) v1.Pod {
	pod := newPodWithoutLivenessProbe(namespace, name, uid)

	pod.Spec.Containers[0].LivenessProbe = &v1.Probe{}
	return pod
}