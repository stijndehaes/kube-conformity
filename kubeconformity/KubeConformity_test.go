package kubeconformity

import (
	"testing"
	"log"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
	"bytes"
	"github.com/stijndehaes/kube-conformity/config"
	"github.com/stijndehaes/kube-conformity/rules"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	kubeConformity := setup(t, pods)
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
	kubeConformity := setup(t, pods)
	kubeConformity.LogNonConformingPods()
	logOutput.String()
	assert.Equal(t, "rule name: \nrule reason: Labels: [app] are not filled in\nfoo_default()\n", logOutput.String())
}

func setup(t *testing.T, pods []v1.Pod) *KubeConformity {
	client := fake.NewSimpleClientset()

	for _, pod := range pods {
		if _, err := client.Core().Pods(pod.Namespace).Create(&pod); err != nil {
			t.Fatal(err)
		}
	}

	kubeConfig := config.KubeConformityConfig{
		LabelsFilledInRules: []rules.LabelsFilledInRule{
			{Labels: []string{"app"}},
		},
	}

	logOutput.Reset()

	return New(client, logger, kubeConfig)
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