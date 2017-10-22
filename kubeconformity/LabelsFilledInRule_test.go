package kubeconformity

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFilterOnLabelsFilledIn(t *testing.T) {
	rule := LabelsFilledInRule{[]string{"app"}}
	pod1 := newPodWithLabel("testing", "bar1", "test")
	pod2 := newPodWithLabel("testing", "bar2", "app")
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.findNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsFilledInMultipleLables(t *testing.T) {
	rule := LabelsFilledInRule{[]string{"app"}}
	pod1 := newPodWithLabels("testing", "bar1", []string{})
	pod2 := newPodWithLabels("testing", "bar2", []string{"app","environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.findNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsFilledInAllLabelsMatch(t *testing.T) {
	rule := LabelsFilledInRule{[]string{"app","environment"}}
	pod1 := newPodWithLabels("testing", "bar1", []string{})
	pod2 := newPodWithLabels("testing", "bar2", []string{"app","environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.findNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsFilledInOnlyOneLabelMatch(t *testing.T) {
	rule := LabelsFilledInRule{[]string{"app","environment"}}
	pod1 := newPodWithLabels("testing", "bar1", []string{"app","other"})
	pod2 := newPodWithLabels("testing", "bar2", []string{"app","environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.findNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
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