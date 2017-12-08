package rules

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"gopkg.in/yaml.v2"
)

func TestFilterOnLabelsFilledIn(t *testing.T) {
	rule := PodRuleLabelsFilledIn{
		Name:   "app label",
		Labels: []string{"app"},
	}
	pod1 := newPodWithLabel("testing", "bar1", "test")
	pod2 := newPodWithLabel("testing", "bar2", "app")
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsFilledInMultipleLabels(t *testing.T) {
	rule := PodRuleLabelsFilledIn{
		Name:   "app label",
		Labels: []string{"app"},
	}
	pod1 := newPodWithLabels("testing", "bar1", []string{})
	pod2 := newPodWithLabels("testing", "bar2", []string{"app", "environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsFilledInAllLabelsMatch(t *testing.T) {
	rule := PodRuleLabelsFilledIn{
		Name:   "app label",
		Labels: []string{"app", "environment"},
	}
	pod1 := newPodWithLabels("testing", "bar1", []string{})
	pod2 := newPodWithLabels("testing", "bar2", []string{"app", "environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func TestFilterOnLabelsFilledInOnlyOneLabelMatch(t *testing.T) {
	rule := PodRuleLabelsFilledIn{
		Name:   "app label",
		Labels: []string{"app", "environment"},
	}
	pod1 := newPodWithLabels("testing", "bar1", []string{"app", "other"})
	pod2 := newPodWithLabels("testing", "bar2", []string{"app", "environment"})
	pods := []v1.Pod{
		pod1,
		pod2,
	}
	result := rule.FindNonConformingPods(pods)
	assert.Equal(t, 1, len(result.Pods))
	assert.Equal(t, pod1.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
	assert.NotEqual(t, pod2.ObjectMeta.Name, result.Pods[0].ObjectMeta.Name)
}

func TestPodRuleLabelsFilledIn_UnmarshalYAML_LabelsNotFilledIn(t *testing.T) {
	string := `
name: app label filled in`

	rule := PodRuleLabelsFilledIn{}

	err := yaml.Unmarshal([]byte(string), &rule)

	if err == nil {
		t.Fail()
	}
}

func TestPodRuleLabelsFilledIn_UnmarshalYAML_NameNotFilledIn(t *testing.T) {
	string := `
labels:
- app`

	rule := PodRuleLabelsFilledIn{}

	err := yaml.Unmarshal([]byte(string), &rule)

	if err == nil {
		t.Fail()
	}
}

func TestPodRuleLabelsFilledIn_UnmarshalYAML(t *testing.T) {
	string := `
name: app label filled in
labels:
- app`

	rule := PodRuleLabelsFilledIn{}

	err := yaml.Unmarshal([]byte(string), &rule)

	if err != nil {
		t.Fail()
	}
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
			Labels:    labelMap,
		},
	}
}
