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

	deploymentFilter := DeploymentFilter{}
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

	podFilter := PodFilter{}
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
	filter := DeploymentFilter{}

	deployments := []v1beta1.Deployment{
		newDeployment("default", "name1", "uid1"),
	}

	filteredPods := filter.FilterDeployments(deployments)
	assert.Len(t, filteredPods, 1)
}

func TestPodFilter_FilterPods(t *testing.T) {
	filter := PodFilter{}

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

func TestFilter_FilterExcludeNamespace(t *testing.T) {
	filter := Filter{
		ExcludeNamespaces: []string{"default"},
	}

	pods :=[]v1.Pod{
		newPod("default", "name1", "uid1"),
		newPod("kube-system", "name2", "uid2"),
	}
	objects := convertPodsToObjects(pods)

	filteredObjects := filter.FilterExcludeNamespace(objects)
	assert.Len(t, filteredObjects, 1)
	assert.Equal(t, "name2", filteredObjects[0].GetName())
	assert.Equal(t, "kube-system", filteredObjects[0].GetNamespace())
}

func TestFilter_FilterExcludeNamespace_Empty(t *testing.T) {
	filter := Filter{
		ExcludeNamespaces: []string{},
	}

	objects := convertPodsToObjects([]v1.Pod{newPod("default", "name", "uid1")})

	filteredObjects := filter.FilterExcludeNamespace(objects)
	assert.Len(t, filteredObjects, 1)
}

func TestFilter_FilterExcludeAnnotations(t *testing.T) {
	filter := Filter{
		ExcludeAnnotations: map[string]string{"testkey1": "testvalue1"},
	}

	pods := []v1.Pod{
		newPodWithAnnotations("default", "name1", "uid1", map[string]string{"testkey1": "testvalue1"}),
		newPodWithAnnotations("default", "name2", "uid2", map[string]string{"testkey1": "testvalue2"}),
		newPodWithAnnotations("default", "name3", "uid3", map[string]string{"testkey2": "testvalue1"}),
	}
	objects := convertPodsToObjects(pods)

	filteredObjects := filter.FilterExcludeAnnotations(objects)
	assert.Len(t, filteredObjects, 2)
	assert.Equal(t, "name2", filteredObjects[0].GetName())
	assert.Equal(t, "default", filteredObjects[0].GetNamespace())
	assert.Equal(t, "name3", filteredObjects[1].GetName())
	assert.Equal(t, "default", filteredObjects[1].GetNamespace())
}

func TestFilter_FilterExcludeAnnotations_Empty(t *testing.T) {
	filter := Filter{
		ExcludeAnnotations: map[string]string{},
	}

	objects := convertPodsToObjects([]v1.Pod{newPod("default", "name", "uid1")})

	filteredObjects := filter.FilterExcludeAnnotations(objects)
	assert.Len(t, filteredObjects, 1)
}

func TestFilter_FilterExcludeLabels(t *testing.T) {
	filter := Filter{
		ExcludeLabels: map[string]string{"testkey1": "testvalue1"},
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", "uid1", map[string]string{"testkey1": "testvalue1"}),
		newPodWithLabels("default", "name2", "uid2", map[string]string{"testkey1": "testvalue2"}),
		newPodWithLabels("default", "name3", "uid3", map[string]string{"testkey2": "testvalue1"}),
	}
	objects := convertPodsToObjects(pods)

	filteredObjects := filter.FilterExcludeLabels(objects)
	assert.Len(t, filteredObjects, 2)
	assert.Equal(t, "name2", filteredObjects[0].GetName())
	assert.Equal(t, "default", filteredObjects[0].GetNamespace())
	assert.Equal(t, "name3", filteredObjects[1].GetName())
	assert.Equal(t, "default", filteredObjects[1].GetNamespace())
}

func TestFilter_FilterExcludeLabels_Empty(t *testing.T) {
	filter := Filter{
		ExcludeLabels: map[string]string{},
	}

	objects := convertPodsToObjects([]v1.Pod{newPod("default", "name", "uid1")})

	filteredObjects := filter.FilterExcludeLabels(objects)
	assert.Len(t, filteredObjects, 1)
}

func Test_convertPodsToObjects(t *testing.T) {
	pods :=[]v1.Pod{
		newPod("default", "name1", "uid1"),
		newPod("kube-system", "name2", "uid2"),
	}
	objects := convertPodsToObjects(pods)

	assert.Len(t, objects, 2)
	assert.Equal(t, objects[0], pods[0].GetObjectMeta())
	assert.Equal(t, objects[1], pods[1].GetObjectMeta())
	assert.NotEqual(t, objects[0], objects[1])
}

func Test_convertDeploymentsToObjects(t *testing.T) {
	deployments :=[]v1beta1.Deployment{
		newDeployment("default", "name1", "uid1"),
		newDeployment("kube-system", "name2", "uid2"),
	}
	objects := convertDeploymentsToObjects(deployments)

	assert.Len(t, objects, 2)
	assert.Equal(t, objects[0], deployments[0].GetObjectMeta())
	assert.Equal(t, objects[1], deployments[1].GetObjectMeta())
	assert.NotEqual(t, objects[0], objects[1])
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