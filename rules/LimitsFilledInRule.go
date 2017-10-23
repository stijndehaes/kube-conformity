package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"github.com/stijndehaes/kube-conformity/filters"
)

type LimitsFilledInRule struct {
	Name   string         `yaml:"name"`
	Filter filters.Filter `yaml:"filter"`
}

func (r LimitsFilledInRule) FindNonConformingPods(pods []v1.Pod) RuleResult {
	namespaceFiltered := r.Filter.FilterPods(pods)
	filteredList := []v1.Pod{}
	for _, pod := range namespaceFiltered {
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
