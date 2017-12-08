package filters

import (
	"testing"
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/stretchr/testify/assert"
)

func newPod(namespace, name string) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
}

func newPodWithLabels(namespace, name string, labels map[string]string) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels:    labels,
		},
	}
}

func newPodWithAnnotations(namespace, name string, annotations map[string]string) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			Annotations: annotations,
		},
	}
}

func TestFilter_FilterIncludeNamespace(t *testing.T) {
	filter := PodFilter{
		IncludeNamespaces: []string{"default"},
	}

	pods := []v1.Pod{
		newPod("default", "name1"),
		newPod("kube-system", "name2"),
	}

	filteredPods := filter.FilterIncludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestFilter_FilterIncludeNamespace_Empty(t *testing.T) {
	filter := PodFilter{
		IncludeNamespaces: []string{},
	}

	pods := []v1.Pod{newPod("default", "name")}

	filteredPods := filter.FilterIncludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
}

func TestFilter_FilterExcludeNamespace(t *testing.T) {
	filter := PodFilter{
		ExcludeNamespaces: []string{"kube-system"},
	}

	pods := []v1.Pod{
		newPod("default", "name1"),
		newPod("kube-system", "name2"),
	}

	filteredPods := filter.FilterExcludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestFilter_FilterExcludeNamespace_Empty(t *testing.T) {
	filter := PodFilter{
		ExcludeNamespaces: []string{},
	}

	pods := []v1.Pod{
		newPod("default", "name1"),
	}

	filteredPods := filter.FilterIncludeNamespace(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name1", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
}

func TestFilter_FilterExcludeAnnotations(t *testing.T) {
	filter := PodFilter{
		ExcludeAnnotations: map[string]string{"testkey1": "testvalue1"},
	}

	pods := []v1.Pod{
		newPodWithAnnotations("default", "name1", map[string]string{"testkey1": "testvalue1"}),
		newPodWithAnnotations("default", "name2", map[string]string{"testkey1": "testvalue2"}),
		newPodWithAnnotations("default", "name3", map[string]string{"testkey2": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeAnnotations(pods)
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "name2", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
	assert.Equal(t, "name3", filteredPods[1].Name)
	assert.Equal(t, "default", filteredPods[1].Namespace)
}

func TestFilter_FilterExcludeAnnotations_Empty(t *testing.T) {
	filter := PodFilter{
		ExcludeAnnotations: map[string]string{},
	}

	pods := []v1.Pod{
		newPodWithAnnotations("default", "name1", map[string]string{"testkey1": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeAnnotations(pods)
	assert.Len(t, filteredPods, 1)
}

func TestFilter_FilterExcludeLabels(t *testing.T) {
	filter := PodFilter{
		ExcludeLabels: map[string]string{"testkey1": "testvalue1"},
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", map[string]string{"testkey1": "testvalue1"}),
		newPodWithLabels("default", "name2", map[string]string{"testkey1": "testvalue2"}),
		newPodWithLabels("default", "name3", map[string]string{"testkey2": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeLabels(pods)
	assert.Len(t, filteredPods, 2)
	assert.Equal(t, "name2", filteredPods[0].Name)
	assert.Equal(t, "default", filteredPods[0].Namespace)
	assert.Equal(t, "name3", filteredPods[1].Name)
	assert.Equal(t, "default", filteredPods[1].Namespace)
}

func TestFilter_FilterExcludeLabels_Empty(t *testing.T) {
	filter := PodFilter{
		ExcludeLabels: map[string]string{},
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", map[string]string{"testkey1": "testvalue1"}),
	}

	filteredPods := filter.FilterExcludeLabels(pods)
	assert.Len(t, filteredPods, 1)
}

func TestFilter_FilterExcludeJobs_true(t *testing.T) {
	filter := PodFilter{
		ExcludeJobs: true,
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", map[string]string{"job-name": "curator-212312"}),
		newPodWithLabels("default", "name2", map[string]string{}),
	}

	filteredPods := filter.FilterExcludeJobs(pods)
	assert.Len(t, filteredPods, 1)
	assert.Equal(t, "name2", filteredPods[0].Name)
}

func TestFilter_FilterExcludeJobs_false(t *testing.T) {
	filter := PodFilter{
		ExcludeJobs: false,
	}

	pods := []v1.Pod{
		newPodWithLabels("default", "name1", map[string]string{"job-name": "curator-212312"}),
		newPodWithLabels("default", "name2", map[string]string{}),
	}

	filteredPods := filter.FilterExcludeJobs(pods)
	assert.Len(t, filteredPods, 2)
}

func TestFilter_FilterPods(t *testing.T) {
	filter := PodFilter{}

	pods := []v1.Pod{
		newPod("default", "name1"),
	}

	filteredPods := filter.FilterPods(pods)
	assert.Len(t, filteredPods, 1)
}