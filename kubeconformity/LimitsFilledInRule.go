package kubeconformity

import (
	"k8s.io/client-go/pkg/api/v1"
)

type LimitsFilledInRule struct {
	Name string `yaml:"name"`
}

func (r LimitsFilledInRule) findNonConformingPods(pods []v1.Pod) RuleResult {
	filteredList := []v1.Pod{}
	for _, pod := range pods {
		var podNonConform = false

		for _, container := range pod.Spec.Containers {
			podNonConform = podNonConform || container.Resources.Limits.Cpu().IsZero()
			podNonConform = podNonConform || container.Resources.Limits.Memory().IsZero()
		}

		if podNonConform {
			filteredList = append(filteredList, pod)
		}
	}

	return RuleResult{
		Pods:     filteredList,
		Reason:   "Limits are not filled in",
		RuleName: r.Name,
	}
}
