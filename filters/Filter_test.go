package filters

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestDeploymentFilter_UnmarshalYAML(t *testing.T) {
	test := `
include_namespaces:
- include_namespace
exclude_namespaces:
- exclude_namespace
exclude_annotations:
  annotationKey: annotationValue
exclude_labels:
  labelKey: labelValue`

	deploymentFilter := DeploymentFilter2{}
	err := yaml.Unmarshal([]byte(test), &deploymentFilter)

	if err != nil {
		assert.Fail(t, "UnMarshal should not fail")
	}

	assert.Len(t, deploymentFilter.Filter.IncludeNamespaces, 1)
	assert.Equal(t, "include_namespace", deploymentFilter.IncludeNamespaces[0])
	assert.Len(t, deploymentFilter.ExcludeNamespaces, 1)
	assert.Equal(t, "exclude_namespace", deploymentFilter.ExcludeNamespaces[0])
	assert.Len(t, deploymentFilter.ExcludeAnnotations, 1)
	assert.Equal(t, "annotationValue", deploymentFilter.ExcludeAnnotations["annotationKey"])
	assert.Len(t, deploymentFilter.ExcludeLabels, 1)
	assert.Equal(t, "labelValue", deploymentFilter.ExcludeLabels["labelKey"])
}

func TestPodFilter_UnmarshalYAML(t *testing.T) {
	test := `
include_namespaces:
- include_namespace
exclude_namespaces:
- exclude_namespace
exclude_annotations:
  annotationKey: annotationValue
exclude_labels:
  labelKey: labelValue
exclude_jobs: true`

	podFilter := PodFilter2{}
	err := yaml.Unmarshal([]byte(test), &podFilter)

	if err != nil {
		assert.Fail(t, "UnMarshal should not fail")
	}

	assert.Len(t, podFilter.Filter.IncludeNamespaces, 1)
	assert.Equal(t, "include_namespace", podFilter.IncludeNamespaces[0])
	assert.Len(t, podFilter.ExcludeNamespaces, 1)
	assert.Equal(t, "exclude_namespace", podFilter.ExcludeNamespaces[0])
	assert.Len(t, podFilter.ExcludeAnnotations, 1)
	assert.Equal(t, "annotationValue", podFilter.ExcludeAnnotations["annotationKey"])
	assert.Len(t, podFilter.ExcludeLabels, 1)
	assert.Equal(t, "labelValue", podFilter.ExcludeLabels["labelKey"])
	assert.Equal(t, true, podFilter.ExcludeJobs)
}

func TestDeploymentFilter_FilterDeployments(t *testing.T) {
	filter := DeploymentFilter2{}

	deployments := []v1beta1.Deployment{
		newDeployment("default", "name1", "uid1"),
	}

	filteredPods := filter.FilterDeployments(deployments)
	assert.Len(t, filteredPods, 1)
}

func TestPodFilter_FilterPods(t *testing.T) {
	filter := PodFilter2{}

	pods := []v1.Pod{
		newPod("default", "name1", "uid1"),
	}

	filteredPods := filter.FilterPods(pods)
	assert.Len(t, filteredPods, 1)
}

func TestFilter_FilterIncludeNamespace(t *testing.T) {
	filter := Filter{
		IncludeNamespaces: []string{"default"},
	}

	pods :=[]v1.Pod{
		newPod("default", "name1", "uid1"),
		newPod("kube-system", "name2", "uid2"),
	}
	objects := convertPodsToObjects(pods)

	filteredObjects := filter.FilterIncludeNamespace(objects)
	assert.Len(t, filteredObjects, 1)
	assert.Equal(t, "name1", filteredObjects[0].GetName())
	assert.Equal(t, "default", filteredObjects[0].GetNamespace())
}

func TestFilter_FilterIncludeNamespace_Empty(t *testing.T) {
	filter := Filter{
		IncludeNamespaces: []string{},
	}

	objects := convertPodsToObjects([]v1.Pod{newPod("default", "name", "uid1")})

	filteredObjects := filter.FilterIncludeNamespace(objects)
	assert.Len(t, filteredObjects, 1)
}

func convertPodsToObjects(pods []v1.Pod) []metav1.Object {
	var objects []metav1.Object
	for _, pod := range pods {
		objects = append(objects, pod.GetObjectMeta())
	}
	return objects
}

func newDeployment(namespace, name string, uid types.UID) v1beta1.Deployment {
	return v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID: uid,
		},
	}
}

func newDeploymentWithLabels(namespace, name string, uid types.UID, labels map[string]string) v1beta1.Deployment {
	return v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID: uid,
			Labels:    labels,
		},
	}
}

func newDeploymentWithAnnotations(namespace, name string, uid types.UID, annotations map[string]string) v1beta1.Deployment {
	return v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			UID: uid,
			Annotations: annotations,
		},
	}
}

func newPod(namespace, name string, uid types.UID) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID: uid,
		},
	}
}

func newPodWithLabels(namespace, name string, uid types.UID, labels map[string]string) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			UID: uid,
			Labels:    labels,
		},
	}
}

func newPodWithAnnotations(namespace, name string, uid types.UID, annotations map[string]string) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			UID: uid,
			Annotations: annotations,
		},
	}
}