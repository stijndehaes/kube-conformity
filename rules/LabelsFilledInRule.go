package rules

import (
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
	"github.com/stijndehaes/kube-conformity/filters"
)

type LabelsFilledInRule struct {
	Name   string         `yaml:"name"`
	Labels []string       `yaml:"labels"`
	Filter filters.Filter `yaml:"filter"`
}

func (r LabelsFilledInRule) FindNonConformingPods(pods []v1.Pod) RuleResult {
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

	return RuleResult{
		Pods:     filteredList,
		Reason:   fmt.Sprintf("Labels: %v are not filled in", r.Labels),
		RuleName: r.Name,
	}
}

func (r *LabelsFilledInRule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain LabelsFilledInRule
	if err := unmarshal((*plain)(r)); err != nil {
		return err
	}
	if len(r.Labels) == 0 {
		return fmt.Errorf("Missing labels for LabelsFilledInRule")
	}
	if r.Name == "" {
		return fmt.Errorf("Missing name for LabelsFilledInRule")
	}
	return nil
}
