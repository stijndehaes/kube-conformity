package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
	"github.com/stijndehaes/kube-conformity/filters"
)

type PodRuleReadinessProbeFilledIn struct {
	Name   string            `yaml:"name"`
	Filter filters.PodFilter `yaml:"filter"`
}

func (r PodRuleReadinessProbeFilledIn) FindNonConformingPods(pods []v1.Pod) PodRuleResult {
	filteredPods := r.Filter.FilterPods(pods)
	var nonConformingPods []v1.Pod
	for _, pod := range filteredPods {
		for _, container := range pod.Spec.Containers {
			if container.ReadinessProbe == nil {
				nonConformingPods = append(nonConformingPods, pod)
				break
			}
		}
	}

	return PodRuleResult{
		Pods:     nonConformingPods,
		Reason:   fmt.Sprintf("ReadinessProbes are not filled in"),
		RuleName: r.Name,
	}
}

func (r *PodRuleReadinessProbeFilledIn) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain PodRuleReadinessProbeFilledIn
	if err := unmarshal((*plain)(r)); err != nil {
		return err
	}
	if r.Name == "" {
		return fmt.Errorf("Missing name for PodRuleReadinessProbeFilledIn")
	}
	return nil
}
