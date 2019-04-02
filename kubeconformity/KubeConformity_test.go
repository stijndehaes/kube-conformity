package kubeconformity

import (
	"testing"
	"log"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/api/core/v1"
	"bytes"
	"github.com/stijndehaes/kube-conformity/config"
	"github.com/stijndehaes/kube-conformity/rules"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

var logOutput = bytes.NewBuffer([]byte{})
var logger = log.New(logOutput, "", 0)

// TestCandidatesNamespaces tests that the list of pods available for
// termination can be restricted by namespaces.
func TestKubeConformity_EvaluatePodRules(t *testing.T) {
	kubeConfig := config.Config{
		PodRulesLabelsFilledIn: []rules.PodRuleLabelsFilledIn{
			{Labels: []string{"app"}},
		},
		PodRulesLimitsFilledIn:   []rules.PodRuleLimitsFilledIn{{}},
		PodRulesRequestsFilledIn: []rules.PodRuleRequestsFilledIn{{}},
	}
	pods := []v1.Pod{
		newPodWithLabels("default", "foo", "uid1", []string{}),
		newPodWithLabels("testing", "bar", "uid2", []string{"app"}),
	}
	kubeConformity := setup(t, pods, nil, nil, kubeConfig)
	conformityResult := kubeConformity.EvaluatePodRules()
	assert.Equal(t, 3, len(conformityResult))
}

func TestKubeConformity_EvaluateDeploymentRulesRules(t *testing.T) {
	kubeConfig := config.Config{
		DeploymentRuleReplicasMinimum: []rules.DeploymentRuleReplicasMinimum{{
			MinimumReplicas: 2,
		}},
	}
	deployments := []appsv1.Deployment{
		newDeployment("default", "foo", "uid1", 1),
		newDeployment("testing", "bar", "uid2", 2),
	}
	kubeConformity := setup(t, nil, deployments, nil, kubeConfig)
	conformityResult := kubeConformity.EvaluateDeploymentRules()
	assert.Equal(t, 1, len(conformityResult))
}

func TestKubeConformity_EvaluateStatefulSetRulesRules(t *testing.T) {
	kubeConfig := config.Config{
		StatefulSetRuleReplicasMinimum: []rules.StatefulSetRuleReplicasMinimum{{
			MinimumReplicas: 2,
		}},
	}
	statefulSets := []appsv1.StatefulSet{
		newStatefulSet("default", "foo", "uid1", 1),
		newStatefulSet("testing", "bar", "uid2", 2),
	}
	kubeConformity := setup(t, nil, nil, statefulSets, kubeConfig)
	conformityResult := kubeConformity.EvaluateStatefulSetRules()
	assert.Equal(t, 1, len(conformityResult))
}

func TestKubeConformity_LogNonConforming_Pods(t *testing.T) {
	kubeConfig := config.Config{
		PodRulesLabelsFilledIn: []rules.PodRuleLabelsFilledIn{
			{Labels: []string{"app"}},
		},
	}
	pods := []v1.Pod{
		newPodWithLabels("default", "foo", "uid1", []string{}),
		newPodWithLabels("testing", "bar", "uid2", []string{"app"}),
	}
	kubeConformity := setup(t, pods, nil, nil, kubeConfig)
	kubeConformity.LogNonConforming()
	logOutput.String()
	assert.Equal(t, "Presenting Pod rule results\nrule name: \nrule reason: Labels: [app] are not filled in\nfoo_default\n", logOutput.String())
}

func TestKubeConformity_LogNonConforming_Deployments(t *testing.T) {
	kubeConfig := config.Config{
		DeploymentRuleReplicasMinimum: []rules.DeploymentRuleReplicasMinimum{{
			MinimumReplicas: 2,
		}},
	}
	deployments := []appsv1.Deployment{
		newDeployment("default", "foo", "uid1", 1),
		newDeployment("testing", "bar", "uid2", 2),
	}
	kubeConformity := setup(t, nil, deployments, nil, kubeConfig)
	kubeConformity.LogNonConforming()
	logOutput.String()
	assert.Equal(t, "Presenting Deployment rule results\nrule name: \nrule reason: Deployment replicas below the minimum: 2\nfoo_default\n", logOutput.String())
}

func TestKubeConformity_LogNonConforming_StatefulSets(t *testing.T) {
	kubeConfig := config.Config{
		StatefulSetRuleReplicasMinimum: []rules.StatefulSetRuleReplicasMinimum{{
			MinimumReplicas: 2,
		}},
	}
	statefulSets := []appsv1.StatefulSet{
		newStatefulSet("default", "foo", "uid1", 1),
		newStatefulSet("testing", "bar", "uid2", 2),
	}
	kubeConformity := setup(t, nil, nil, statefulSets, kubeConfig)
	kubeConformity.LogNonConforming()
	logOutput.String()
	assert.Equal(t, "Presenting StatefulSet rule results\nrule name: \nrule reason: StatefulSet replicas below the minimum: 2\nfoo_default\n", logOutput.String())
}

func setup(t *testing.T, pods []v1.Pod, deployments []appsv1.Deployment, statefulSets []appsv1.StatefulSet, kubeConfig config.Config) *KubeConformity {
	client := fake.NewSimpleClientset()

	for _, pod := range pods {
		if _, err := client.Core().Pods(pod.Namespace).Create(&pod); err != nil {
			t.Fatal(err)
		}
	}
	for _, deployment := range deployments {
		if _, err := client.AppsV1().Deployments(deployment.Namespace).Create(&deployment); err != nil {
			t.Fatal(err)
		}
	}
	for _, statefulSet := range statefulSets {
		if _, err := client.AppsV1().StatefulSets(statefulSet.Namespace).Create(&statefulSet); err != nil {
			t.Fatal(err)
		}
	}
	logOutput.Reset()

	return New(client, logger, kubeConfig)
}

func newPodWithLabels(namespace, name string, uid types.UID, labels []string) v1.Pod {
	labelMap := make(map[string]string)
	for _, label := range labels {
		labelMap[label] = "randomString"
	}
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID: 	   uid,
			Labels:    labelMap,
		},
	}
}

func newDeployment(namespace, name string, uid types.UID, replicas int32) appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID:       uid,
		},
		Spec:appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
	}
}

func newStatefulSet(namespace, name string, uid types.UID, replicas int32) appsv1.StatefulSet {
	return appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID:       uid,
		},
		Spec:appsv1.StatefulSetSpec{
			Replicas: &replicas,
		},
	}
}
