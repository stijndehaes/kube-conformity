package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"github.com/stijndehaes/kube-conformity/filters"
)

type RequestsFilledInRule struct {
	Name   string         `yaml:"name"`
	Filter filters.Filter `yaml:"filter"`
}

func (r RequestsFilledInRule) FindNonConformingPods(pods []v1.Pod) RuleResult {
	namespaceFiltered := r.Filter.FilterPods(pods)
	filteredList := []v1.Pod{}
	for _, pod := range namespaceFiltered {
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
