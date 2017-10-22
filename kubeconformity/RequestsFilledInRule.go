package kubeconformity

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

type RequestsFilledInRule struct {
}

func (r RequestsFilledInRule) findNonConformingPods(client kubernetes.Interface) (RuleResult, error) {
	podList, err := client.CoreV1().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return RuleResult{}, err
	}
	filteredPods := filterOnRequestsFilledIn(podList.Items)
	return RuleResult{
		Pods:   filteredPods,
		Reason: "Requests are not filled in",
	}, nil
}

func filterOnRequestsFilledIn(pods []v1.Pod) ([]v1.Pod) {

	filteredList := []v1.Pod{}
	for _, pod := range pods {
		var podNonConform = false

		for _, container := range pod.Spec.Containers {
			podNonConform = podNonConform || container.Resources.Requests.Cpu().IsZero()
			podNonConform = podNonConform || container.Resources.Requests.Memory().IsZero()
		}

		if podNonConform {
			filteredList = append(filteredList, pod)
		}
	}

	return filteredList
}
