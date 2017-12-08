package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"github.com/stijndehaes/kube-conformity/filters"
	"fmt"
)

type PodRuleRequestsFilledIn struct {
	Name   string            `yaml:"name"`
	Filter filters.PodFilter `yaml:"filter"`
}

func (r PodRuleRequestsFilledIn) FindNonConformingPods(pods []v1.Pod) PodRuleResult {
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
	return PodRuleResult{
		Pods:     filteredList,
		Reason:   "Requests are not filled in",
		RuleName: r.Name,
	}
}

func (r *PodRuleRequestsFilledIn) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain PodRuleLimitsFilledIn
	if err := unmarshal((*plain)(r)); err != nil {
		return err
	}
	if r.Name == "" {
		return fmt.Errorf("Missing name for PodRuleRequestsFilledIn")
	}
	return nil
}
