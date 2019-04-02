package rules

import (
	"k8s.io/api/core/v1"
	"fmt"
	"github.com/stijndehaes/kube-conformity/filters"
)

type PodRuleLabelsFilledIn struct {
	Name   string            `yaml:"name"`
	Labels []string          `yaml:"labels"`
	Filter filters.PodFilter `yaml:"filter"`
}

func (r PodRuleLabelsFilledIn) FindNonConformingPods(pods []v1.Pod) PodRuleResult {
	filteredPods := r.Filter.FilterPods(pods)
	var nonConformingPods []v1.Pod
	for _, pod := range filteredPods {
		for _, label := range r.Labels {
			containsLabel := false
			for podLabelKey := range pod.ObjectMeta.Labels {
				if podLabelKey == label {
					containsLabel = true
				}
			}
			if !containsLabel {
				nonConformingPods = append(nonConformingPods, pod)
				break
			}
		}
	}

	return PodRuleResult{
		Pods:     nonConformingPods,
		Reason:   fmt.Sprintf("Labels: %v are not filled in", r.Labels),
		RuleName: r.Name,
	}
}

func (r *PodRuleLabelsFilledIn) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain PodRuleLabelsFilledIn
	if err := unmarshal((*plain)(r)); err != nil {
		return err
	}
	if len(r.Labels) == 0 {
		return fmt.Errorf("missing labels for PodRuleLabelsFilledIn")
	}
	if r.Name == "" {
		return fmt.Errorf("missing name for PodRuleLabelsFilledIn")
	}
	return nil
}
