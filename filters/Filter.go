package filters

import (
	"k8s.io/client-go/pkg/api/v1"
)

type Filter struct {
	IncludeNamespaces  []string          `yaml:"include_namespaces"`
	ExcludeNamespaces  []string          `yaml:"exclude_namespaces"`
	ExcludeAnnotations map[string]string `yaml:"exclude_annotations"`
	ExcludeJobs        bool              `yaml:"exclude_jobs"`
	ExcludeLabels      map[string]string `yaml:"exclude_labels"`
}

func (f Filter) FilterPods(pods []v1.Pod) []v1.Pod {
	filteredPods := f.FilterIncludeNamespace(pods)
	filteredPods = f.FilterExcludeNamespace(filteredPods)
	filteredPods = f.FilterExcludeAnnotations(filteredPods)
	filteredPods = f.FilterExcludeLabels(filteredPods)
	filteredPods = f.FilterExcludeJobs(filteredPods)
	return filteredPods
}

func (f Filter) FilterIncludeNamespace(pods []v1.Pod) []v1.Pod {
	if len(f.IncludeNamespaces) == 0 {
		return pods
	}

	filteredPods := []v1.Pod{}

	for _, pod := range pods {
		include := false
		for _, namespace := range f.IncludeNamespaces {
			include = include || pod.Namespace == namespace
		}
		if include {
			filteredPods = append(filteredPods, pod)
		}
	}
	return filteredPods
}

func (f Filter) FilterExcludeNamespace(pods []v1.Pod) []v1.Pod {
	if len(f.ExcludeNamespaces) == 0 {
		return pods
	}

	filteredPods := []v1.Pod{}

	for _, pod := range pods {
		exclude := false
		for _, namespace := range f.ExcludeNamespaces {
			exclude = exclude || pod.Namespace == namespace
		}
		if !exclude {
			filteredPods = append(filteredPods, pod)
		}
	}
	return filteredPods
}

func (f Filter) FilterExcludeAnnotations(pods []v1.Pod) []v1.Pod {
	if len(f.ExcludeAnnotations) == 0 {
		return pods
	}

	filteredPods := []v1.Pod{}

	for _, pod := range pods {
		exclude := false
		for annotationsKey, annotationValue := range f.ExcludeAnnotations {
			if podAnnotationValue, exists := pod.Annotations[annotationsKey]; exists {
				exclude = exclude || podAnnotationValue == annotationValue
			}
		}
		if !exclude {
			filteredPods = append(filteredPods, pod)
		}
	}
	return filteredPods
}

func (f Filter) FilterExcludeLabels(pods []v1.Pod) []v1.Pod {
	if len(f.ExcludeLabels) == 0 {
		return pods
	}

	filteredPods := []v1.Pod{}

	for _, pod := range pods {
		exclude := false
		for labelKey, labelValue := range f.ExcludeLabels {
			if podLabelsValue, exists := pod.Labels[labelKey]; exists {
				exclude = exclude || podLabelsValue == labelValue
			}
		}
		if !exclude {
			filteredPods = append(filteredPods, pod)
		}
	}
	return filteredPods
}

func (f Filter) FilterExcludeJobs(pods []v1.Pod) []v1.Pod {
	if !f.ExcludeJobs {
		return pods
	}

	filteredPods := []v1.Pod{}

	for _, pod := range pods {
		if _, exists := pod.Labels["job-name"]; !exists {
			filteredPods = append(filteredPods, pod)
		}
	}
	return filteredPods
}
