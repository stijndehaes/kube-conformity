package filters

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Filter struct {
	IncludeNamespaces  []string          `yaml:"include_namespaces"`
	ExcludeNamespaces  []string          `yaml:"exclude_namespaces"`
	ExcludeAnnotations map[string]string `yaml:"exclude_annotations"`
	ExcludeLabels      map[string]string `yaml:"exclude_labels"`
}

type DeploymentFilter struct {
	Filter `yaml:",inline"`
}

type StatefulsetFilter struct {
	Filter `yaml:",inline"`
}

type PodFilter struct {
	Filter           `yaml:",inline"`
	ExcludeJobs bool `yaml:"exclude_jobs"`
}

func (f Filter) FilterObjects(objects []metav1.Object) []metav1.Object {
	filteredObjects := f.FilterIncludeNamespace(objects)
	filteredObjects = f.FilterExcludeNamespace(filteredObjects)
	filteredObjects = f.FilterExcludeAnnotations(filteredObjects)
	filteredObjects = f.FilterExcludeLabels(filteredObjects)
	return filteredObjects
}

func convertPodsToObjects(pods []apiv1.Pod) []metav1.Object {
	var objects []metav1.Object
	for idx := range pods {
		objects = append(objects, pods[idx].GetObjectMeta())
	}
	return objects
}

func convertDeploymentsToObjects(deployments []appsv1.Deployment) []metav1.Object {
	var objects []metav1.Object
	for idx := range deployments {
		objects = append(objects, deployments[idx].GetObjectMeta())
	}
	return objects
}

func convertStatefulSetToObjects(statefulSets []appsv1.StatefulSet) []metav1.Object {
	var objects []metav1.Object
	for idx := range statefulSets {
		objects = append(objects, statefulSets[idx].GetObjectMeta())
	}
	return objects
}


func (f PodFilter) FilterPods(pods []apiv1.Pod) []apiv1.Pod {
	objects := convertPodsToObjects(pods)
	filteredObjects := f.FilterObjects(objects)
	filteredObjects = f.FilterExcludeJobs(filteredObjects)
	var filteredPods []apiv1.Pod
	for _, pod := range pods {
		for _, object := range filteredObjects {
			if object.GetUID() == pod.GetUID() {
				filteredPods = append(filteredPods, pod)
			}
		}
	}
	return filteredPods
}

func (f DeploymentFilter) FilterDeployments(deployments []appsv1.Deployment) []appsv1.Deployment {
	objects := convertDeploymentsToObjects(deployments)
	filteredObjects := f.FilterObjects(objects)
	var filteredDeployments []appsv1.Deployment
	for _, deployment := range deployments {
		for _, object := range filteredObjects {
			if object.GetUID() == deployment.GetUID() {
				filteredDeployments = append(filteredDeployments, deployment)
			}
		}
	}
	return filteredDeployments
}

func (f StatefulsetFilter) FilterStatefulSets(statefulSets []appsv1.StatefulSet) []appsv1.StatefulSet {
	objects := convertStatefulSetToObjects(statefulSets)
	filteredObjects := f.FilterObjects(objects)
	var filteredStatefulsets []appsv1.StatefulSet
	for _, statefulset := range statefulSets {
		for _, object := range filteredObjects {
			if object.GetUID() == statefulset.GetUID() {
				filteredStatefulsets = append(filteredStatefulsets, statefulset)
			}
		}
	}
	return filteredStatefulsets
}

func (f Filter) FilterIncludeNamespace(objects []metav1.Object) []metav1.Object {
	if len(f.IncludeNamespaces) == 0 {
		return objects
	}

	var filteredObjects []metav1.Object

	for _, object := range objects {
		include := false
		for _, namespace := range f.IncludeNamespaces {
			include = include || object.GetNamespace() == namespace
		}
		if include {
			filteredObjects = append(filteredObjects, object)
		}
	}
	return filteredObjects
}

func (f Filter) FilterExcludeNamespace(objects []metav1.Object) []metav1.Object {
	if len(f.ExcludeNamespaces) == 0 {
		return objects
	}

	var filteredObjects []metav1.Object

	for _, object := range objects {
		exclude := false
		for _, namespace := range f.ExcludeNamespaces {
			exclude = exclude || object.GetNamespace() == namespace
		}
		if !exclude {
			filteredObjects = append(filteredObjects, object)
		}
	}
	return filteredObjects
}

func (f Filter) FilterExcludeAnnotations(objects []metav1.Object) []metav1.Object {
	if len(f.ExcludeAnnotations) == 0 {
		return objects
	}

	var filteredObjects []metav1.Object

	for _, object := range objects {
		exclude := false
		for annotationsKey, annotationValue := range f.ExcludeAnnotations {
			if podAnnotationValue, exists := object.GetAnnotations()[annotationsKey]; exists {
				exclude = exclude || podAnnotationValue == annotationValue
			}
		}
		if !exclude {
			filteredObjects = append(filteredObjects, object)
		}
	}
	return filteredObjects
}

func (f Filter) FilterExcludeLabels(objects []metav1.Object) []metav1.Object {
	if len(f.ExcludeLabels) == 0 {
		return objects
	}

	var filteredObjects []metav1.Object

	for _, object := range objects {
		exclude := false
		for labelKey, labelValue := range f.ExcludeLabels {
			if podLabelsValue, exists := object.GetLabels()[labelKey]; exists {
				exclude = exclude || podLabelsValue == labelValue
			}
		}
		if !exclude {
			filteredObjects = append(filteredObjects, object)
		}
	}
	return filteredObjects
}

func (f PodFilter) FilterExcludeJobs(objects []metav1.Object) []metav1.Object {
	if !f.ExcludeJobs {
		return objects
	}

	var filteredObjects []metav1.Object

	for _, object := range objects {
		if _, exists := object.GetLabels()["job-name"]; !exists {
			filteredObjects = append(filteredObjects, object)
		}
	}
	return filteredObjects
}
