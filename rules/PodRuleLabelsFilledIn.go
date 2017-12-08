package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
	"github.com/stijndehaes/kube-conformity/filters"
)

type PodRuleLabelsFilledIn struct {
	Name   string            `yaml:"name"`
	Labels []string          `yaml:"labels"`
	Filter filters.PodFilter `yaml:"filter"`
}

func (r PodRuleLabelsFilledIn) FindNonConformingPods(pods []v1.Pod) PodRuleResult {
	namespaceFiltered := r.Filter.FilterPods(pods)
	filteredList := []v1.Pod{}
	for _, pod := range namespaceFiltered {
		for _, label := range r.Labels {
			containsLabel := false
			for podLabelKey := range pod.ObjectMeta.Labels {
				if podLabelKey == label {
					containsLabel = true
				}
			}
			if !containsLabel {
				filteredList = append(filteredList, pod)
				break
			}
		}
	}

	return PodRuleResult{
		Pods:     filteredList,
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
		return fmt.Errorf("Missing labels for PodRuleLabelsFilledIn")
	}
	if r.Name == "" {
		return fmt.Errorf("Missing name for PodRuleLabelsFilledIn")
	}
	return nil
}
