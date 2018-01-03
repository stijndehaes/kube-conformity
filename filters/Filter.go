package filters

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

type Filter struct {
	IncludeNamespaces  []string          `yaml:"include_namespaces"`
	ExcludeNamespaces  []string          `yaml:"exclude_namespaces"`
	ExcludeAnnotations map[string]string `yaml:"exclude_annotations"`
	ExcludeLabels      map[string]string `yaml:"exclude_labels"`
}

type DeploymentFilter2 struct {
	Filter `yaml:",inline"`
}

type PodFilter2 struct {
	Filter          `yaml:",inline"`
	ExcludeJobs bool `yaml:"exclude_jobs"`
}

func (f Filter) FilterObjects(pods []metav1.Object) []metav1.Object {
	filteredPods := f.FilterIncludeNamespace(pods)
	filteredPods = f.FilterExcludeNamespace(filteredPods)
	filteredPods = f.FilterExcludeAnnotations(filteredPods)
	filteredPods = f.FilterExcludeLabels(filteredPods)
	return filteredPods
}

func (f PodFilter2) FilterPods(pods []v1.Pod) []v1.Pod {
	var objects []metav1.Object
	for _, pod := range pods {
		objects = append(objects, pod.GetObjectMeta())
	}
	filteredObjects := f.FilterObjects(objects)
	filteredObjects = f.FilterExcludeJobs(filteredObjects)
	var filteredPods []v1.Pod
	for _, pod := range pods {
		for _, object := range filteredObjects {
			if object.GetUID() == pod.GetUID() {
				filteredPods = append(filteredPods, pod)
			}
		}
	}
	return filteredPods
}

func (f DeploymentFilter2) FilterDeployments(deployments []v1beta1.Deployment) []v1beta1.Deployment {
	var objects []metav1.Object
	for _, deployment := range deployments {
		objects = append(objects, deployment.GetObjectMeta())
	}
	filteredObjects := f.FilterObjects(objects)
	var filteredDeployments []v1beta1.Deployment
	for _, deployment := range deployments {
		for _, object := range filteredObjects {
			if object.GetUID() == deployment.GetUID() {
				filteredDeployments = append(filteredDeployments, deployment)
			}
		}
	}
	return filteredDeployments
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

func (f PodFilter2) FilterExcludeJobs(objects []metav1.Object) []metav1.Object {
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
