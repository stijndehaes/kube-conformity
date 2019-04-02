package rules

import (
	"fmt"
	"github.com/stijndehaes/kube-conformity/filters"
	appsv1 "k8s.io/api/apps/v1"
)

type StatefulSetRuleReplicasMinimum struct {
	Name            string                    `yaml:"name"`
	MinimumReplicas int32                     `yaml:"minimum_replicas"`
	Filter          filters.StatefulsetFilter `yaml:"filter"`
}

func (statefulSetRuleReplicasMinimum StatefulSetRuleReplicasMinimum) FindNonConformingStatefulSet(statefulSets []appsv1.StatefulSet) StatefulSetRuleResult {
	filteredStatefulsets := statefulSetRuleReplicasMinimum.Filter.FilterStatefulSets(statefulSets)
	var nonConformingStatefulSets []appsv1.StatefulSet
	for _, statefulset := range filteredStatefulsets {
		if *statefulset.Spec.Replicas < statefulSetRuleReplicasMinimum.MinimumReplicas {
			nonConformingStatefulSets = append(nonConformingStatefulSets, statefulset)
		}
	}

	return StatefulSetRuleResult{
		StatefulSets: nonConformingStatefulSets,
		Reason:       fmt.Sprintf("StatefulSet replicas below the minimum: %v", statefulSetRuleReplicasMinimum.MinimumReplicas),
		RuleName:     statefulSetRuleReplicasMinimum.Name,
	}
}

func (statefulSetRuleReplicasMinimum *StatefulSetRuleReplicasMinimum) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain StatefulSetRuleReplicasMinimum
	if err := unmarshal((*plain)(statefulSetRuleReplicasMinimum)); err != nil {
		return err
	}
	if statefulSetRuleReplicasMinimum.MinimumReplicas == 0 {
		return fmt.Errorf("missing minimum replicas")
	}
	if statefulSetRuleReplicasMinimum.Name == "" {
		return fmt.Errorf("missing name for StatefulSetRuleReplicasMinimum")
	}
	return nil
}
