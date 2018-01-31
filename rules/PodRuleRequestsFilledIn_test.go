package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"gopkg.in/yaml.v2"
)

func TestFilterOnRequestsFilledIn(t *testing.T) {
	rule := PodRuleRequestsFilledIn{}
	pod1 := newPodWithRequests("default", "foo", "", "")
	pod2 := newPodWithRequests("testing", "bar", "400m", "1.1Gi")
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func newPodWithRequests(namespace, name, requestCpu, requestMemory string) v1.Pod {
	resources := v1.ResourceRequirements{
		Requests:   make(v1.ResourceList),
	}
	if requestCpu != "" {
		resources.Requests[v1.ResourceCPU] = resource.MustParse(requestCpu)
	}
	if requestMemory != "" {
		resources.Requests[v1.ResourceMemory] = resource.MustParse(requestMemory)
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

func TestPodRuleRequestsFilledIn_UnmarshalYAML_NameNotFilledIn(t *testing.T) {
	string := `
filter:
  namespaces: test`

	rule := PodRuleRequestsFilledIn{}

	err := yaml.Unmarshal([]byte(string), &rule)

	if err == nil {
		t.Fail()
	}
}

func TestPodRuleRequestsFilledIn_UnmarshalYAML(t *testing.T) {
	string := `name: requests filled in`

	rule := PodRuleRequestsFilledIn{}

	err := yaml.Unmarshal([]byte(string), &rule)

	if err != nil {
		t.Fail()
	}
}