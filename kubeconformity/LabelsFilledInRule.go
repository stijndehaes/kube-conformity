package kubeconformity

import (
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
)

type LabelsFilledInRule struct {
	Name   string   `yaml:"name"`
	Labels []string `yaml:"labels"`
}

func (r LabelsFilledInRule) findNonConformingPods(pods []v1.Pod) RuleResult {
	filteredList := []v1.Pod{}
	if len(r.Labels) == 0 {
		return RuleResult{
			Pods:     filteredList,
			Reason:   fmt.Sprintf("Labels: %v are not filled in", r.Labels),
			RuleName: r.Name,
		}
	}
	for _, pod := range pods {
		for _, label := range r.Labels {
			containsLabel := false
			for podLabelKey, _ := range pod.ObjectMeta.Labels {
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
