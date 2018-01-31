package filters

import (
	"testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

func newDeployment(namespace, name string) v1beta1.Deployment {
	return v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
}

func newDeploymentWithLabels(namespace, name string, labels map[string]string) v1beta1.Deployment {
	return v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels:    labels,
		},
	}
}

func newDeploymentWithAnnotations(namespace, name string, annotations map[string]string) v1beta1.Deployment {
	return v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			Annotations: annotations,
		},
	}
}

func TestDeploymentFilter_FilterIncludeNamespace(t *testing.T) {
	filter := DeploymentFilter{
		IncludeNamespaces: []string{"default"},
	}

	deployments := []v1beta1.Deployment{
		newDeployment("default", "name1"),
		newDeployment("kube-system", "name2"),
	}

	filteredPods := filter.FilterIncludeNamespace(deployments)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestDeploymentFilter_FilterIncludeNamespace_Empty(t *testing.T) {
	filter := DeploymentFilter{
		IncludeNamespaces: []string{},
	}

	deployments := []v1beta1.Deployment{newDeployment("default", "name")}

	filteredPods := filter.FilterIncludeNamespace(deployments)
	assert.Len(t, filteredPods, 1)
}

func TestDeploymentFilter_FilterExcludeNamespace(t *testing.T) {
	filter := DeploymentFilter{
		ExcludeNamespaces: []string{"kube-system"},
	}

	deployments := []v1beta1.Deployment{
		newDeployment("default", "name1"),
		newDeployment("kube-system", "name2"),
	}

	filteredPods := filter.FilterExcludeNamespace(deployments)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestDeploymentFilter_FilterExcludeNamespace_Empty(t *testing.T) {
	filter := DeploymentFilter{
		ExcludeNamespaces: []string{},
	}

	deployments := []v1beta1.Deployment{
		newDeployment("default", "name1"),
	}

	filteredPods := filter.FilterIncludeNamespace(deployments)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestDeploymentFilter_FilterExcludeAnnotations(t *testing.T) {
	filter := DeploymentFilter{
		ExcludeAnnotations: map[string]string{"testkey1": "testvalue1"},
	}

	deployments := []v1beta1.Deployment{
		newDeploymentWithAnnotations("default", "name1", map[string]string{"testkey1": "testvalue1"}),
		newDeploymentWithAnnotations("default", "name2", map[string]string{"testkey1": "testvalue2"}),
		newDeploymentWithAnnotations("default", "name3", map[string]string{"testkey2": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeAnnotations(deployments)
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "name2", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
	assert.Equal(t, "name3", filteredPods[1].Name)
	assert.Equal(t, "default", filteredPods[1].Namespace)
}

func TestDeploymentFilter_FilterExcludeAnnotations_Empty(t *testing.T) {
	filter := DeploymentFilter{
		ExcludeAnnotations: map[string]string{},
	}

	deployments := []v1beta1.Deployment{
		newDeploymentWithAnnotations("default", "name1", map[string]string{"testkey1": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeAnnotations(deployments)
	assert.Len(t, filteredPods, 1)
}

func TestDeploymentFilter_FilterExcludeLabels(t *testing.T) {
	filter := DeploymentFilter{
		ExcludeLabels: map[string]string{"testkey1": "testvalue1"},
	}

	deployments := []v1beta1.Deployment{
		newDeploymentWithLabels("default", "name1", map[string]string{"testkey1": "testvalue1"}),
		newDeploymentWithLabels("default", "name2", map[string]string{"testkey1": "testvalue2"}),
		newDeploymentWithLabels("default", "name3", map[string]string{"testkey2": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeLabels(deployments)
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "name2", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
	assert.Equal(t, "name3", filteredPods[1].Name)
	assert.Equal(t, "default", filteredPods[1].Namespace)
}

func TestDeploymentFilter_FilterExcludeLabels_Empty(t *testing.T) {
	filter := DeploymentFilter{
		ExcludeLabels: map[string]string{},
	}

	deployments := []v1beta1.Deployment{
		newDeploymentWithLabels("default", "name1", map[string]string{"testkey1": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeLabels(deployments)
	assert.Len(t, filteredPods, 1)
}

func TestDeploymentFilter_FilterDeployments(t *testing.T) {
	filter := DeploymentFilter{}

	deployments := []v1beta1.Deployment{
		newDeployment("default", "name1"),
	}

	filteredPods := filter.FilterDeployments(deployments)
	assert.Len(t, filteredPods, 1)
}
