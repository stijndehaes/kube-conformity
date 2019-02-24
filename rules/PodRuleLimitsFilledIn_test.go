package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/types"
)

func TestFilterOnLimits(t *testing.T) {
	rule := PodRuleLimitsFilledIn{}
	pod1 := newPodWithLimits("default", "foo", "uid1", "", "")
	pod2 := newPodWithLimits("testing", "bar", "uid2", "400m", "1.1Gi")
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func newPodWithLimits(namespace, name string, uid types.UID, limitCpu, limitMemory string) v1.Pod {
	resources := v1.ResourceRequirements{
		Limits:   make(v1.ResourceList),
	}
	if limitCpu != "" {
		resources.Limits[v1.ResourceCPU] = resource.MustParse(limitCpu)
	}
	if limitMemory != "" {
		resources.Limits[v1.ResourceMemory] = resource.MustParse(limitMemory)
	}
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
					Resources: resources,
				},
			},
		},
	}
}

func TestPodRuleLimitsFilledIn_UnmarshalYAML_NameNotFilledIn(t *testing.T) {
	yamlString := `
filter:
  namespaces: test`

	rule := PodRuleLimitsFilledIn{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err == nil {
		t.Fail()
	}
}

func TestPodRuleLimitsFilledIn_UnmarshalYAML(t *testing.T) {
	yamlString := `name: limits filled in`

	rule := PodRuleLimitsFilledIn{}

	err := yaml.Unmarshal([]byte(yamlString), &rule)

	if err != nil {
		t.Fail()
	}
}