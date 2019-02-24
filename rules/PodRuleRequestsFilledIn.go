package rules

import (
	"k8s.io/api/core/v1"
	"github.com/stijndehaes/kube-conformity/filters"
	"fmt"
)

type PodRuleRequestsFilledIn struct {
	Name   string            `yaml:"name"`
	Filter filters.PodFilter `yaml:"filter"`
}

func (r PodRuleRequestsFilledIn) FindNonConformingPods(pods []v1.Pod) PodRuleResult {
	filteredPods := r.Filter.FilterPods(pods)
	var nonConformingPods []v1.Pod
	for _, pod := range filteredPods {
		var podNonConform = false

		for _, container := range pod.Spec.Containers {
			podNonConform = podNonConform || container.Resources.Requests.Cpu().IsZero()
			podNonConform = podNonConform || container.Resources.Requests.Memory().IsZero()
		}

		if podNonConform {
			nonConformingPods = append(nonConformingPods, pod)
		}
	}
	return PodRuleResult{
		Pods:     nonConformingPods,
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
