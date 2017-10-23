package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestFilterOnLimits(t *testing.T) {
	rule := LimitsFilledInRule{}
	pod1 := newPodWithLimits("default", "foo", "", "")
	pod2 := newPodWithLimits("testing", "bar", "400m", "1.1Gi")
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func newPodWithLimits(namespace, name, limitCpu, limitMemory string) v1.Pod {
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