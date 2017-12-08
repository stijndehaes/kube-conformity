package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"github.com/stijndehaes/kube-conformity/filters"
	"fmt"
)

type PodRuleLimitsFilledIn struct {
	Name   string            `yaml:"name"`
	Filter filters.PodFilter `yaml:"filter"`
}

func (r PodRuleLimitsFilledIn) FindNonConformingPods(pods []v1.Pod) PodRuleResult {
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

	return PodRuleResult{
		Pods:     filteredList,
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
		return fmt.Errorf("Missing name for PodRuleLimitsFilledIn")
	}
	return nil
}
