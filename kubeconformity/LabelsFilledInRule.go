package kubeconformity

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
)

type LabelsFilledInRule struct {
	Labels []string
}

func (r LabelsFilledInRule) findNonConformingPods(client kubernetes.Interface) (RuleResult, error) {
	podList, err := client.CoreV1().Pods(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return RuleResult{}, err
	}
	filteredPods := filterOnLabelsFilledIn(podList.Items, r.Labels)
	return RuleResult{
		Pods:   filteredPods,
		Reason: fmt.Sprintf("Labels: %v are not filled in", r.Labels),
	}, nil
}

func filterOnLabelsFilledIn(pods []v1.Pod, labels []string) ([]v1.Pod) {

	filteredList := []v1.Pod{}
	if len(labels) == 0 {
		return filteredList
	}
	for _, pod := range pods {
		for _, label := range labels {
			containsLabel := false
			for podLabelKey, _ := range pod.ObjectMeta.Labels {
				if podLabelKey == label {
					containsLabel = true
				}
			}
			if !containsLabel {
				filteredList = append(filteredList, pod)
				break
			}
		}
	}

	return filteredList
}
