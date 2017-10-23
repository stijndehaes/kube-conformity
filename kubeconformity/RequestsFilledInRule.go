package kubeconformity

import (
	"k8s.io/client-go/pkg/api/v1"
)

type RequestsFilledInRule struct {
	Name string `yaml:"name"`
}

func (r RequestsFilledInRule) findNonConformingPods(pods []v1.Pod) RuleResult {
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
	return RuleResult{
		Pods:     filteredList,
		Reason:   "Requests are not filled in",
		RuleName: r.Name,
	}
}
