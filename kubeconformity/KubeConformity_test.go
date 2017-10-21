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
	pods := []v1.Pod{
		newPod("default", "foo"),
		newPodWithAllFilledIn("testing", "bar", []string{}),
	}
	kubeConformity := setup(t, []string{}, pods)
	conformityResult, err := kubeConformity.FindNonConformingPods()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(conformityResult.ResourceProblems))
	assert.Equal(t, 0, len(conformityResult.LabelProblems))
}

func TestLogNonConformingPodsResources(t *testing.T) {
	pods := []v1.Pod{
		newPodWithAllFilledIn("default", "foo1", []string{}),
		newPod("default", "foo2"),
	}
	kubeConformity := setup(t, []string{}, pods)
	kubeConformity.LogNonConformingPods()
	logOutput.String()
	assert.Equal(t, "Found pods without resources set\nfoo2_default()\n", logOutput.String())
}

func TestLogNonConformingPodsLabels(t *testing.T) {
	pods := []v1.Pod{
		newPodWithAllFilledIn("default", "foo1", []string{"app"}),
		newPodWithAllFilledIn("default", "foo2", []string{"test"}),
	}
	kubeConformity := setup(t, []string{"app"}, pods)
	kubeConformity.LogNonConformingPods()
	logOutput.String()
	assert.Equal(t, "Found pods with label problems\nfoo2_default()\n", logOutput.String())
}

func TestFilterOnResources(t *testing.T) {
	pod1 := newPod("default", "foo")
	pod2 := newPodWithAllFilledIn("testing", "bar", []string{})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	pods, err := filterOnResources(pods)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1.ObjectMeta.Name, pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, pods[0].ObjectMeta.Name)
}

func TestFilterOnResourcesRequestsFilledIn(t *testing.T) {
	pod1 := newPodWithRequestAndLimit("default", "foo", "", "", "400m", "1.1Gi", []string{})
	pod2 := newPodWithAllFilledIn("testing", "bar", []string{})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	pods, err := filterOnResources(pods)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1.ObjectMeta.Name, pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, pods[0].ObjectMeta.Name)
}

func TestFilterOnResourcesLimitsFilledIn(t *testing.T) {
	pod1 := newPodWithRequestAndLimit("default", "foo", "400m", "1.1Gi", "", "", []string{})
	pod2 := newPodWithAllFilledIn("testing", "bar", []string{})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	pods, err := filterOnResources(pods)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1.ObjectMeta.Name, pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, pods[0].ObjectMeta.Name)
}

func TestFilterOnLabels(t *testing.T) {
	pod1 := newPodWithLabel("testing", "bar1", "test")
	pod2 := newPodWithLabel("testing", "bar2", "app")
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	pods, err := filterOnLabels(pods, []string{"app"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1.ObjectMeta.Name, pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsMultipleLables(t *testing.T) {
	pod1 := newPod("testing", "bar1")
	pod2 := newPodWithLabels("testing", "bar2", []string{"app","environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	pods, err := filterOnLabels(pods, []string{"app","environment"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1.ObjectMeta.Name, pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsAllLabelsMatch(t *testing.T) {
	pod1 := newPod("testing", "bar1")
	pod2 := newPodWithLabels("testing", "bar2", []string{"app","environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	pods, err := filterOnLabels(pods, []string{"app","environment"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1.ObjectMeta.Name, pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsOnlyOneLabelMatch(t *testing.T) {
	pod1 := newPodWithLabels("testing", "bar1", []string{"app","other"})
	pod2 := newPodWithLabels("testing", "bar2", []string{"app","environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	pods, err := filterOnLabels(pods, []string{"app","environment"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(pods))
	assert.Equal(t, pod1.ObjectMeta.Name, pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, pods[0].ObjectMeta.Name)
}

func newPodWithLabel(namespace, name, label string) v1.Pod {
	return newPodWithLabels(namespace, name, []string{label})
}

func newPodWithLabels(namespace, name string, labels []string) v1.Pod {
	labelMap := make(map[string]string)
	for _, label := range labels {
		labelMap[label] = "randomString"
	}
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: labelMap,
		},
	}
}

func newPod(namespace, name string) v1.Pod {
	return newPodWithRequestAndLimit(namespace, name, "", "", "", "", []string{})
}

func newPodWithAllFilledIn(namespace, name string, labels []string) v1.Pod {
	return newPodWithRequestAndLimit(namespace, name, "100m", "1.1Gi", "100m", "1.1Gi", labels)
}

func newPodWithRequestAndLimit(namespace, name, limitCpu, limitMemory, requestCpu, requestMemory string, labels []string) v1.Pod {
	labelMap := make(map[string]string)
	for _, label := range labels {
		labelMap[label] = "randomString"
	}
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
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: 	namespace,
			Name:      	name,
			Labels: 	labelMap,
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

func setup(t *testing.T, labels []string, pods []v1.Pod) *KubeConformity {
	client := fake.NewSimpleClientset()

	for _, pod := range pods {
		if _, err := client.Core().Pods(pod.Namespace).Create(&pod); err != nil {
			t.Fatal(err)
		}
	}

	logOutput.Reset()

	return New(client, logger, labels)
}
