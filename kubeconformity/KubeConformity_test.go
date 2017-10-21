package kubeconformity

import (
	"testing"
	"log"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"bytes"
)

var logOutput = bytes.NewBuffer([]byte{})
var logger = log.New(logOutput, "", 0)

// TestCandidatesNamespaces tests that the list of pods available for
// termination can be restricted by namespaces.
func TestFindNonConformingPods(t *testing.T) {
	kubeConformity := setup(t)
	pods, err := kubeConformity.FindNonConformingPods()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
}

func TestFilterConformingPodsNothingFilledIn(t *testing.T) {
	pod1 := newPod("default", "foo")
	pods := []v1.Pod{
		*pod1,
		*newPodWithRequestAndLimit("testing", "bar", "400m", "1.1Gi", "400m", "1.1Gi"),
	}
	pods, err := filterConformingPods(pods)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1, pods[0])
}

func TestFilterConformingPodsRequestsFilledIn(t *testing.T) {
	pod1 := newPodWithRequestAndLimit("default", "foo", "", "", "400m", "1.1Gi")
	resource := pod1.Spec.Containers[0].Resources
	logger.Fatal(resource)
	pods := []v1.Pod{
		*pod1,
		*newPodWithRequestAndLimit("testing", "bar", "400m", "1.1Gi", "400m", "1.1Gi"),
	}
	pods, err := filterConformingPods(pods)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1, pods[0])
}

func newPod(namespace, name string) *v1.Pod {
	return newPodWithRequestAndLimit(namespace, name, "", "", "", "")
}

func newPodWithRequestAndLimit(namespace, name, limitCpu, limitMemory, requestCpu, requestMemory string) *v1.Pod {
	resources := v1.ResourceRequirements{
		Limits:   make(v1.ResourceList),
		Requests: make(v1.ResourceList),
	}
	if limitCpu != "" {
		resources.Limits[v1.ResourceCPU] = resource.MustParse(limitCpu)
	}
	if limitMemory != "" {
		resources.Limits[v1.ResourceMemory] = resource.MustParse(limitMemory)
	}
	if requestCpu != "" {
		resources.Requests[v1.ResourceCPU] = resource.MustParse(requestCpu)
	}
	if requestMemory != "" {
		resources.Requests[v1.ResourceMemory] = resource.MustParse(requestMemory)
	}
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
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

func setup(t *testing.T) *KubeConformity {
	pods := []v1.Pod{
		*newPod("default", "foo"),
		*newPodWithRequestAndLimit("testing", "bar" , "400m", "1.1Gi", "400m", "1.1Gi"),
	}

	client := fake.NewSimpleClientset()

	for _, pod := range pods {
		if _, err := client.Core().Pods(pod.Namespace).Create(&pod); err != nil {
			t.Fatal(err)
		}
	}

	logOutput.Reset()

	return New(client, logger)
}
