package kubeconformity

import (
	"k8s.io/client-go/pkg/api/v1"
	"fmt"
)

type LabelsFilledInRule struct {
	Labels []string
}

func (r LabelsFilledInRule) findNonConformingPods(pods []v1.Pod) RuleResult {
	filteredList := []v1.Pod{}
	if len(r.Labels) == 0 {
		return RuleResult{filteredList, fmt.Sprintf("Labels: %v are not filled in", r.Labels)}
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

	return RuleResult{filteredList, fmt.Sprintf("Labels: %v are not filled in", r.Labels)}
}
