package kubeconformity

import (
	"testing"
	"log"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	"bytes"
)

var logOutput = bytes.NewBuffer([]byte{})
var logger = log.New(logOutput, "", 0)

// TestCandidatesNamespaces tests that the list of pods available for
// termination can be restricted by namespaces.
func TestFindNonConformingPods(t *testing.T) {
	pods := []v1.Pod{
		newPodWithLabels("default", "foo", []string{}),
		newPodWithLabels("testing", "bar", []string{"app"}),
	}
	kubeConformity := setup(t, []string{}, pods)
	conformityResult, err := kubeConformity.EvaluateRules()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(conformityResult))
}

func TestLogNonConformingPodsResources(t *testing.T) {
	pods := []v1.Pod{
		newPodWithLabels("default", "foo", []string{}),
		newPodWithLabels("testing", "bar", []string{"app"}),
	}
	kubeConformity := setup(t, []string{}, pods)
	kubeConformity.LogNonConformingPods()
	logOutput.String()
	assert.Equal(t, "Labels are not filled in\nfoo_default()\n", logOutput.String())
}

func setup(t *testing.T, labels []string, pods []v1.Pod) *KubeConformity {
	client := fake.NewSimpleClientset()

	for _, pod := range pods {
		if _, err := client.Core().Pods(pod.Namespace).Create(&pod); err != nil {
			t.Fatal(err)
		}
	}

	rules := []Rule{
		LabelsFilledInRule{Labels: []string{"app"}},
	}

	logOutput.Reset()

	return New(client, logger, rules)
}
