package rules

import (
	"github.com/stijndehaes/kube-conformity/filters"
	"k8s.io/client-go/pkg/api/v1"
)

type PodRuleLimitsFilledIn struct {
	Name   string            `yaml:"name"`
	Filter filters.PodFilter `yaml:"filter"`
}

func (r PodRuleLimitsFilledIn) FindNonConformingPods(pods []v1.Pod) PodRuleResult {
	filteredPods := r.Filter.FilterPods(pods)
	var nonConformingPods []v1.Pod
	for _, pod := range filteredPods {
		var podNonConform = false

		for _, container := range pod.Spec.Containers {
			podNonConform = podNonConform || container.Resources.Limits.Cpu().IsZero()
			podNonConform = podNonConform || container.Resources.Limits.Memory().IsZero()
		}

		if podNonConform {
			nonConformingPods = append(nonConformingPods, pod)
		}
	}

	return PodRuleResult{
		Pods:     nonConformingPods,
		Reason:   "Limits are not filled in",
		RuleName: r.Name,
	}
}

func (r *PodRuleLimitsFilledIn) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain PodRuleLimitsFilledIn
	if err := unmarshal((*plain)(r)); err != nil {
		return err
	}
	if r.Name == "" {
		r.Name = "Pod resource limits are not filled in"
	}
	return nil
}
