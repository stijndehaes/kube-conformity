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
