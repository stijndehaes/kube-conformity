package filters

import "k8s.io/client-go/pkg/apis/apps/v1beta1"

type DeploymentFilter struct {
	IncludeNamespaces  []string          `yaml:"include_namespaces"`
	ExcludeNamespaces  []string          `yaml:"exclude_namespaces"`
	ExcludeAnnotations map[string]string `yaml:"exclude_annotations"`
	ExcludeLabels      map[string]string `yaml:"exclude_labels"`
}

func (f DeploymentFilter) FilterDeployments(deployments []v1beta1.Deployment) []v1beta1.Deployment {
	filteredDeployments := f.FilterIncludeNamespace(deployments)
	filteredDeployments = f.FilterExcludeNamespace(filteredDeployments)
	filteredDeployments = f.FilterExcludeAnnotations(filteredDeployments)
	filteredDeployments = f.FilterExcludeLabels(filteredDeployments)
	return filteredDeployments
}

func (f DeploymentFilter) FilterIncludeNamespace(deployments []v1beta1.Deployment) []v1beta1.Deployment {
	if len(f.IncludeNamespaces) == 0 {
		return deployments
	}

	var filteredDeployments []v1beta1.Deployment

	for _, deployment := range deployments {
		include := false
		for _, namespace := range f.IncludeNamespaces {
			include = include || deployment.Namespace == namespace
		}
		if include {
			filteredDeployments = append(filteredDeployments, deployment)
		}
	}
	return filteredDeployments
}

func (f DeploymentFilter) FilterExcludeNamespace(deployments []v1beta1.Deployment) []v1beta1.Deployment {
	if len(f.ExcludeNamespaces) == 0 {
		return deployments
	}

	var filteredDeployments []v1beta1.Deployment

	for _, deployment := range deployments {
		exclude := false
		for _, namespace := range f.ExcludeNamespaces {
			exclude = exclude || deployment.Namespace == namespace
		}
		if !exclude {
			filteredDeployments = append(filteredDeployments, deployment)
		}
	}
	return filteredDeployments
}

func (f DeploymentFilter) FilterExcludeAnnotations(deployments []v1beta1.Deployment) []v1beta1.Deployment {
	if len(f.ExcludeAnnotations) == 0 {
		return deployments
	}

	var filteredDeployments []v1beta1.Deployment

	for _, deployment := range deployments {
		exclude := false
		for annotationsKey, annotationValue := range f.ExcludeAnnotations {
			if deploymentAnnotationValue, exists := deployment.Annotations[annotationsKey]; exists {
				exclude = exclude || deploymentAnnotationValue == annotationValue
			}
		}
		if !exclude {
			filteredDeployments = append(filteredDeployments, deployment)
		}
	}
	return filteredDeployments
}

func (f DeploymentFilter) FilterExcludeLabels(deployments []v1beta1.Deployment) []v1beta1.Deployment {
	if len(f.ExcludeLabels) == 0 {
		return deployments
	}

	var filteredDeployments []v1beta1.Deployment

	for _, deployment := range deployments {
		exclude := false
		for labelKey, labelValue := range f.ExcludeLabels {
			if DeploymentLabelsValue, exists := deployment.Labels[labelKey]; exists {
				exclude = exclude || DeploymentLabelsValue == labelValue
			}
		}
		if !exclude {
			filteredDeployments = append(filteredDeployments, deployment)
		}
	}
	return filteredDeployments
}